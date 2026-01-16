package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"smartgas-payment/config"
	"smartgas-payment/internal/models"
	"smartgas-payment/internal/schemas"
	"smartgas-payment/internal/utils"
	"strconv"
	"time"

	"github.com/google/uuid"
)

var UnauthorizedEmployeeErr = errors.New("Employee unauthorized")

type ReportTransactionOpts struct {
	Ip     string
	Number string
}

type SetGasPumpOptions struct {
	Number    string
	Ip        string
	Amount    float32
	FuelType  string
	PaymentID uuid.UUID
	Discount  float64
}

type ResponseAccumPoints struct {
	Amount float32
	Id     string
}

type ValidateEmployeeOpts struct {
	ExternalGasStationID string
	EmployeeID           string
	EmployeeNIP          string
	GasPumpCRE           string
}

//go:generate mockery --name SocioSmartService --filename=mock_socio_smart.go --inpackage=true
type SocioSmartService interface {
	GetGasStations() ([]schemas.GasStation, error)
	GetGasPumpsByCrePermission(string) ([]schemas.GasPump, error)
	SetGasPump(SetGasPumpOptions) (*schemas.SetGasPump, error)
	AccumPoints(*models.Payment) (*ResponseAccumPoints, error)
	ReportTransaction(opts ReportTransactionOpts) error
	ValidateEmployee(opts ValidateEmployeeOpts) (bool, error)
}

type socioSmartService struct {
	config config.Config
}

func ProvideSocioSmartService(config config.Config) *socioSmartService {
	return &socioSmartService{
		config: config,
	}
}

func (ss *socioSmartService) ValidateEmployee(opts ValidateEmployeeOpts) (bool, error) {
	url := fmt.Sprintf(
		"%s/rest/operacion?&usuario=%s&nip=%s&pv=%s&Cre=%s&Validacion=true",
		ss.config.SocioSmartUrl,
		opts.EmployeeID,
		opts.EmployeeNIP,
		opts.ExternalGasStationID,
		opts.GasPumpCRE,
	)

	client := &http.Client{
		Timeout: time.Second * 2,
	}
	response, err := client.Get(url)
	if err != nil {
		return false, err
	}

	if response.StatusCode != 200 {
		return false, errors.New("Internal error in smartgas, got " + response.Status)
	}

	var data []map[string]string

	json.NewDecoder(response.Body).Decode(&data)

	if data[0]["Estatus"] != "1" {
		return false, nil
	}

	return true, nil
}

func (ss *socioSmartService) GetGasStations() ([]schemas.GasStation, error) {
	url := fmt.Sprintf("%s/rest/operacion", ss.config.SocioSmartUrl)

	client := &http.Client{
		Timeout: time.Second * 2,
	}

	response, err := client.Get(url)
	if err != nil {
		return []schemas.GasStation{}, err
	}

	gasStations := make([]schemas.GasStation, 0)

	json.NewDecoder(response.Body).Decode(&gasStations)

	return gasStations, nil
}

func (ss *socioSmartService) GetGasPumpsByCrePermission(
	crePermission string,
) ([]schemas.GasPump, error) {
	url := fmt.Sprintf(
		"%s/rest/operacion?cre=%v&autopago=true&Todos=true",
		ss.config.SocioSmartUrl,
		crePermission,
	)

	client := &http.Client{
		Timeout: time.Second * 2,
	}

	response, err := client.Get(url)
	if err != nil {
		return []schemas.GasPump{}, err
	}

	gasPumps := make([]schemas.GasPump, 0)

	var data []map[string]any

	json.NewDecoder(response.Body).Decode(&data)

	for _, obj := range data {
		gasPump := utils.Transform[schemas.GasPump](obj)

		dieselPrice, err := strconv.ParseFloat(obj["Ppdiesel"].(string), 64)
		if err != nil {
			// TODO: Log in sentry
		}

		regularPrice, err := strconv.ParseFloat(obj["Ppregular"].(string), 64)
		if err != nil {
			// TODO: Log in sentry
		}
		premiumPrice, err := strconv.ParseFloat(obj["Ppremium"].(string), 64)
		if err != nil {
			// TODO: Log in sentry
		}

		isActive := false
		if obj["EstatusBomba"] == "1" {
			isActive = true
		}

		gasPump.DieselPrice = &dieselPrice
		gasPump.RegularPrice = &regularPrice
		gasPump.PremiumPrice = &premiumPrice
		gasPump.Active = &isActive
		gasPumps = append(gasPumps, *gasPump)
	}

	return gasPumps, nil
}

func (ss *socioSmartService) ReportTransaction(opts ReportTransactionOpts) error {
	// TODO: Check with german possible responses from GG
	// {"uno": "uno"}
	number, _ := strconv.Atoi(opts.Number)
	url := fmt.Sprintf(
		`http://%s:4346/Ticket/Combustible?json={"validarSerie":true,"promotor":"12345","serie":"630C574332F22E","bomba": %v,"tipoPagoCV":6,"clienteID":0,"cuentaID":0,"lealtadTarjeta":"","lealtadTipo":0,"productos":"","tipoVenta":3}`,
		opts.Ip,
		number,
	)

	body, _ := json.Marshal(map[string]any{})

	client := &http.Client{
		Timeout: time.Second * 15,
	}

	_, err := client.Post(url, "", bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	return nil
}

func (ss *socioSmartService) SetGasPump(opts SetGasPumpOptions) (*schemas.SetGasPump, error) {
	fuelTypeNumber := 0

	if opts.FuelType == "premium" {
		fuelTypeNumber = 1
	} else if opts.FuelType == "diesel" {
		fuelTypeNumber = 2
	}

	number, _ := strconv.Atoi(opts.Number)

	url := fmt.Sprintf(
		`http://%s:4346/General/Prefijar?json={"validarSerie":true,"clave":"1123","serie":"630C574332F22E","bomba":%v,"combustible":%v,"cantidad":%v,"tipo":"P","TipoVentaMonedero":0,"preautorizar": "%s","descuento": %v}`,
		opts.Ip,
		number,
		fuelTypeNumber,
		opts.Amount,
		fmt.Sprintf("AP_%s", opts.PaymentID),
		opts.Discount,
	)

	client := &http.Client{
		Timeout: time.Second * 15,
	}

	res, err := client.Post(url, "", nil)
	if err != nil {
		return nil, err
	}

	var responseData schemas.SetGasPump

	json.NewDecoder(res.Body).Decode(&responseData)

	// TODO: Return ID

	return &responseData, nil
}

func (ss *socioSmartService) AccumPoints(payment *models.Payment) (*ResponseAccumPoints, error) {
	// fmt.Println("VECES LLAMADAS")
	url := fmt.Sprintf("%s/rest/operacion?user=65411&bd=2&campana=0", ss.config.SocioSmartUrl)

	client := &http.Client{
		Timeout: time.Second * 30,
	}

	pointsSchema := schemas.AcumPoints{
		TransID:          "GA_" + payment.ID.String(),
		CrePermission:    payment.GasPump.GasStation.CrePermission,
		GasPumpNumber:    payment.GasPump.Number,
		Amount:           fmt.Sprintf("%v", payment.RealAmountReported),
		TotalLiter:       fmt.Sprintf("%v", payment.RealAmountReported/float32(payment.Price)),
		PhoneNumber:      payment.Customer.PhoneNumber,
		Start:            "0",
		Status:           "1",
		ClientRegistered: "0",
	}

	// TODO: Change code whenm required

	if payment.FuelType == "regular" {
		pointsSchema.ProductID = "444"
	} else if payment.FuelType == "premium" {
		pointsSchema.ProductID = "445"
	} else if payment.FuelType == "diesel" {
		pointsSchema.ProductID = "32"
	}

	pointsSchema.FillDateTime()

	values := make([]map[string]any, 0)

	data := make(map[string]any)
	data["0"] = pointsSchema

	values = append(values, data)

	payloadData, _ := json.Marshal(values)

	payloadRequest := bytes.NewReader(payloadData)

	// fmt.Println("PaYLOAD", string(payloadData))

	req, err := http.NewRequest("POST", url, payloadRequest)
	if err != nil {
		// TODO: Report to sentry
		return nil, err
	}

	res, err := client.Do(req)
	if err != nil {
		// TODO: Report to sentry
		return nil, err
	}

	var resData []map[string]any

	json.NewDecoder(res.Body).Decode(&resData)

	if len(resData) < 1 {
		return nil, errors.New("No data to read")
	}

	result := resData[0]["result"].(map[string]any)

	if result["Estatus"].(string) != "1" {
		return nil, errors.New("Error posting points")
	}

	pointsStr, ok := result["PuntosAcumulados"].(string)
	var points float32

	toReturn := ResponseAccumPoints{}

	if ok {
		p, _ := strconv.ParseFloat(pointsStr, 32)
		points = float32(p)
	}

	id, ok := result["Folio"].(string)

	folio := ""

	if ok {
		folio = id
	}

	toReturn.Amount = points
	toReturn.Id = folio

	return &toReturn, nil
}
