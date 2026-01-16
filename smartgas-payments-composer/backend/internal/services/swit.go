package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"smartgas-payment/config"
	"smartgas-payment/internal/dto"
	"smartgas-payment/internal/schemas"
	"smartgas-payment/internal/utils"
	"time"

	"github.com/jinzhu/copier"
)

const (
	timeout = 60 // Timeout in seconds for swit
)

var ErrProccesingPayment = errors.New("CVV error or not funds")

type ReserveFundsOpts struct {
	CustomerID  string  `json:"customerId"`
	SourceID    string  `json:"sourceId"`
	Cvv         string  `json:"cvv"`
	Last4       string  `json:"cardLastDigits"`
	Amount      float32 `json:"amount"`
	Description string  `json:"description"`
	Capture     bool    `json:"capture"`
}

//go:generate mockery --name SwitService --filename=mock_swit.go --inpackage=true
type SwitService interface {
	CreateCustomer(*schemas.Customer) (string, error)
	ListCardsByCustomer(string) ([]schemas.SwitSource, error)
	ReserveFunds(ReserveFundsOpts) (string, error)
	CancelFundReservation(string) error
	ConfirmFundReservation(string, float32) error
	DeleteCard(string, string) error
}

type switService struct {
	cfg config.Config
}

func ProvideSwitService(cfg config.Config) *switService {
	return &switService{
		cfg: cfg,
	}
}

func (ss *switService) addHeaders(r *http.Request) {
	r.Header.Add("business", ss.cfg.SwitBusiness)
	r.Header.Add("token", ss.cfg.SwitToken)
	r.Header.Add("x-api-key", ss.cfg.SwitApiKey)
	r.Header.Add("Content-Type", "application/json")
}

func (ss *switService) DeleteCard(cusID string, cardID string) error {
	client := &http.Client{
		Timeout: time.Second * timeout,
	}

	req, err := http.NewRequest(
		"DELETE",
		ss.cfg.SwitBaseUrl+"/customers/"+cusID+"/cards/"+cardID,
		nil,
	)
	if err != nil {
		return err
	}

	ss.addHeaders(req)

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	var data dto.SwitGeneralResponse

	json.NewDecoder(res.Body).Decode(&data)

	if data.Status != "Success" {
		return errors.New("Error deleting card")
	}

	return nil
}

func (ss *switService) ConfirmFundReservation(transID string, amount float32) error {
	client := &http.Client{
		Timeout: time.Second * timeout,
	}

	payload := struct {
		Amount float32 `json:"amount"`
	}{
		Amount: amount,
	}
	payloadData, _ := json.Marshal(payload)

	payloadRequest := bytes.NewReader(payloadData)

	req, err := http.NewRequest("PUT", ss.cfg.SwitBaseUrl+"/payments/"+transID, payloadRequest)
	if err != nil {
		return err
	}

	ss.addHeaders(req)

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	var data dto.SwitGeneralResponse

	json.NewDecoder(res.Body).Decode(&data)

	if data.Status != "Success" {
		return errors.New("Error confirming")
	}

	return nil
}

func (ss *switService) CancelFundReservation(transID string) error {
	client := &http.Client{
		Timeout: time.Second * timeout,
	}

	req, err := http.NewRequest("PUT", ss.cfg.SwitBaseUrl+"/payments/"+transID+"/cancel", nil)
	if err != nil {
		return err
	}

	ss.addHeaders(req)

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	var data dto.SwitGeneralResponse

	json.NewDecoder(res.Body).Decode(&data)

	if data.Status != "Success" {
		return errors.New("Error canceling: " + res.Status)
	}

	return nil
}

func (ss *switService) ReserveFunds(opts ReserveFundsOpts) (string, error) {
	client := &http.Client{
		Timeout: time.Second * timeout,
	}

	if opts.Description == "" {
		opts.Description = "Automatic charge"
	}

	payloadData, _ := json.Marshal(opts)

	payloadRequest := bytes.NewReader(payloadData)

	req, err := http.NewRequest("POST", ss.cfg.SwitBaseUrl+"/payments", payloadRequest)
	if err != nil {
		return "", err
	}

	ss.addHeaders(req)

	response, err := client.Do(req)
	if err != nil {
		return "", err
	}

	var data dto.SwitGeneralResponse

	json.NewDecoder(response.Body).Decode(&data)

	if data.Status != "Success" {
		return "", ErrProccesingPayment
	}

	result, _ := data.Result.(map[string]any)

	transId := result["transactionId"].(string)

	return transId, nil
}

func (ss *switService) CreateCustomer(cus *schemas.Customer) (string, error) {
	client := &http.Client{
		Timeout: time.Second * timeout,
	}

	payload := struct {
		Email     string `json:"email"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}{
		Email:     cus.Email,
		FirstName: cus.FirstName,
		LastName:  cus.FirstLastName + " " + cus.SecondLastName,
	}

	copier.Copy(&payload, cus)

	payloadData, _ := json.Marshal(payload)

	payloadRequest := bytes.NewReader(payloadData)

	req, err := http.NewRequest("POST", ss.cfg.SwitBaseUrl+"/customers", payloadRequest)
	if err != nil {
		return "", err
	}

	ss.addHeaders(req)

	response, err := client.Do(req)
	if err != nil {
		return "", err
	}

	if response.StatusCode != 200 {
		return "", errors.New("Error creating customer in swit - " + response.Status)
	}

	var data dto.SwitGeneralResponse

	json.NewDecoder(response.Body).Decode(&data)

	if data.Status != "Success" {
		return "", errors.New("Error creating new customer")
	}

	id, _ := data.Result.(string)

	return id, nil
}

func (ss *switService) ListCardsByCustomer(cusID string) ([]schemas.SwitSource, error) {
	client := &http.Client{
		Timeout: time.Second * timeout,
	}

	req, err := http.NewRequest("GET", ss.cfg.SwitBaseUrl+"/customers/"+cusID+"/cards", nil)
	if err != nil {
		return nil, err
	}

	ss.addHeaders(req)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, errors.New("Eror listing customer cards - " + res.Status)
	}

	sources := make([]schemas.SwitSource, 0)

	var data dto.SwitGeneralResponse

	json.NewDecoder(res.Body).Decode(&data)

	if data.Status != "Success" {
		return nil, errors.New("Error listing customer cards")
	}

	for _, s := range data.Result.([]any) {
		v := utils.Transform[schemas.SwitSource](s)

		sources = append(sources, *v)
	}

	return sources, nil
}
