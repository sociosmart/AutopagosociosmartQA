package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"smartgas-payment/config"
	"time"
)

var (
	DebitErrUnauthorized         = errors.New("Unauthorized")
	DebitErrNotFound             = errors.New("Customer not found")
	DebitUnsufficientFunds       = errors.New("Not enough funds")
	DebitInternalServerError     = errors.New("Internal Server Error")
	DebitGiftCardInUse           = errors.New("Gift card in use")
	DebitValidationError         = errors.New("Data Validation error")
	DebitGivenAmountGreaterError = errors.New(
		"Given amount is greater than reserved",
	)
	DebitPaymentAlreadyConfirmedOrCanceledError = errors.New(
		"Given payment is already confirmed or canceled",
	)
)

type DebitReserveFundsOpts struct {
	Amount              float32 `json:"amount"`
	ExternalCustomerID  string  `json:"external_customer_id,omitempty"`
	ExternalLegalNameID string  `json:"external_legal_name_id"`
	CardKey             string  `json:"card_key,omitempty"`
}

//go:generate mockery --name DebitService --filename=mock_debit.go --inpackage=true
type DebitService interface {
	ReserveFunds(DebitReserveFundsOpts) (string, error)
	CancelReservation(string) error
	PaymentConfirmation(string, float32) error
}

type debitService struct {
	config config.Config
}

func (db *debitService) addHeaders(r *http.Request) {
	r.Header.Add("X-APP-KEY", db.config.DebitAppKey)
	r.Header.Add("X-API-KEY", db.config.DebitApiKey)
	r.Header.Add("Content-Type", "application/json")
}

func (db *debitService) PaymentConfirmation(id string, amount float32) error {
	client := http.Client{
		Timeout: time.Second * 10,
	}
	url := fmt.Sprintf("%s/api/v1/payments/confirmation/%s", db.config.DebitBaseUrl, id)
	body, _ := json.Marshal(struct {
		Amount float32 `json:"amount"`
	}{Amount: amount})

	req, _ := http.NewRequest("POST", url, bytes.NewReader(body))

	db.addHeaders(req)
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode == 404 {
		return DebitErrNotFound
	}

	if res.StatusCode == 401 {
		return DebitErrUnauthorized
	}

	if res.StatusCode == 406 {
		return DebitGivenAmountGreaterError
	}

	if res.StatusCode == 412 {
		return DebitPaymentAlreadyConfirmedOrCanceledError
	}

	if res.StatusCode == 500 {
		return DebitInternalServerError
	}

	if res.StatusCode == 422 {
		return DebitValidationError
	}

	return nil
}

func (db *debitService) CancelReservation(id string) error {
	client := http.Client{
		Timeout: time.Second * 10,
	}
	url := fmt.Sprintf("%s/api/v1/payments/cancelation/%s", db.config.DebitBaseUrl, id)
	body, _ := json.Marshal(struct{}{})

	req, _ := http.NewRequest("PUT", url, bytes.NewReader(body))

	db.addHeaders(req)
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode == 404 {
		return DebitErrNotFound
	}

	if res.StatusCode == 401 {
		return DebitErrUnauthorized
	}

	if res.StatusCode == 500 {
		return DebitInternalServerError
	}

	if res.StatusCode == 412 {
		return DebitPaymentAlreadyConfirmedOrCanceledError
	}

	if res.StatusCode == 422 {
		return DebitValidationError
	}

	return nil
}

func (db *debitService) ReserveFunds(opts DebitReserveFundsOpts) (string, error) {
	client := http.Client{
		Timeout: time.Second * 10,
	}
	url := fmt.Sprintf("%s/api/v1/payments", db.config.DebitBaseUrl)
	body, _ := json.Marshal(opts)

	req, _ := http.NewRequest("POST", url, bytes.NewReader(body))

	db.addHeaders(req)
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}

	if res.StatusCode == 401 {
		return "", DebitErrUnauthorized
	}

	if res.StatusCode == 409 {
		return "", DebitGiftCardInUse
	}

	if res.StatusCode == 406 {
		return "", DebitUnsufficientFunds
	}

	if res.StatusCode == 404 {
		return "", DebitErrNotFound
	}

	if res.StatusCode == 500 {
		return "", DebitInternalServerError
	}

	if res.StatusCode == 404 {
		return "", DebitErrNotFound
	}

	if res.StatusCode == 422 {
		return "", DebitValidationError
	}

	var response struct {
		ID string `json:"id"`
	}

	json.NewDecoder(res.Body).Decode(&response)

	return response.ID, nil
}

func ProvideDebitService(config config.Config) *debitService {
	return &debitService{
		config: config,
	}
}
