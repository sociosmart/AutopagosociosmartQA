package services

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"smartgas-payment/config"
	"smartgas-payment/internal/dto"
	"smartgas-payment/internal/models"
	"smartgas-payment/internal/repository"
	"strconv"
	"time"

	"gorm.io/gorm"
)

var (
	ErrData             = errors.New("Error with the given data")
	ErrNotDocumentFound = errors.New("Not documents found")
)

const (
	iva      = 1.16
	version  = "4.0"
	currency = "MXN"
)

type SignInvoiceOpts struct {
	Params  dto.SignInvoiceRequest
	Payment *models.Payment
}

//go:generate mockery --name InvoicingService --filename=mock_invoicing.go --inpackage=true
type InvoicingService interface {
	SignInvoice(SignInvoiceOpts) (string, error)
	ResendInvoice(string, string) error
	GetInvoicePDF(string) (string, error)
	GetIeps(string) (float64, error)
}

type invoicingService struct {
	config       config.Config
	settingsRepo repository.SettingRepository
}

func ProvideInvoicingService(
	config config.Config,
	settingsRepo repository.SettingRepository,
) *invoicingService {
	return &invoicingService{
		config:       config,
		settingsRepo: settingsRepo,
	}
}

func (is *invoicingService) GetIeps(fuelType string) (float64, error) {
	// Setup pump in swit
	setting, err := is.settingsRepo.GetByName("ieps_" + fuelType)
	if err != nil {
		var csmErr error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			csmErr = errors.New("ieps setting not setted: " + "ieps_" + fuelType)
		} else {
			csmErr = err
		}

		return 0, csmErr
	}

	ieps, err := strconv.ParseFloat(setting.Value, 64)
	if err != nil {
		return 0, errors.New(
			"Impossile to parse saved value in database to number - " + err.Error(),
		)
	}

	return ieps, nil
}

func (is *invoicingService) SignInvoice(opts SignInvoiceOpts) (string, error) {
	loc, _ := time.LoadLocation("America/Mazatlan")
	now := time.Now().In(loc)
	dateStr := now.Format("2006-01-02T15:04:05")

	iepsPrice, err := is.GetIeps(opts.Payment.FuelType)
	if err != nil {
		return "", err
	}

	total := float64(opts.Payment.RealAmountReported)
	totalLiter := total / opts.Payment.Price
	iepsTotal := totalLiter * iepsPrice

	totalNoIeps := total - iepsTotal

	base := totalNoIeps / iva
	subtotal := base + iepsTotal

	realPricePerLiter := subtotal / totalLiter

	ivaImp := base * .16

	// fmt.Println("total", total, "subtotal", subtotal, "base", base, "ieps", iepsTotal, "real reported", total, "price x liter", opts.Payment.Price, "total litersd", totalLiter)

	payload := map[string]any{
		"Version":             "4.0",
		"Serie":               "Serie",
		"Folio":               opts.Payment.ID.String(),
		"Fecha":               dateStr,
		"Sello":               "",
		"FormaPago":           "04", // TODO: CHECK IF NEEDS TO PASS PAYMENT METHOD
		"FormaPagoSpecified":  true,
		"NoCertificado":       is.config.Invoicing.CertificateNumber,
		"Certificado":         "",
		"CondicionesDePago":   "Una exhibicion",
		"SubTotal":            fmt.Sprintf("%.2f", subtotal),
		"Moneda":              "MXN",
		"TipoCambio":          1,
		"TipoCambioSpecified": true,
		"Total":               fmt.Sprintf("%.2f", total),
		"TipoDeComprobante":   "I",
		"Exportacion":         "01",
		"MetodoPago":          "PUE",
		"MetodoPagoSpecified": true,
		"LugarExpedicion":     is.config.Invoicing.CpEmitter,
		"Emisor": map[string]any{
			"Rfc":           is.config.Invoicing.Rfc,
			"Nombre":        is.config.Invoicing.Name,
			"RegimenFiscal": is.config.Invoicing.FiscalName,
		},
		"Receptor": map[string]any{
			"Rfc":                       opts.Params.Rfc,
			"Nombre":                    opts.Params.RazonSocial,
			"DomicilioFiscalReceptor":   opts.Params.CP,
			"ResidenciaFiscalSpecified": false,
			"RegimenFiscalReceptor":     opts.Params.RegimenFiscal,
			"UsoCFDI":                   opts.Params.UsoCFDI,
		},
		"Impuestos": map[string]any{
			"Traslados": []map[string]any{
				{
					"Base":       fmt.Sprintf("%.2f", base),
					"Importe":    fmt.Sprintf("%.2f", ivaImp),
					"Impuesto":   "002",
					"TasaOCuota": "0.160000",
					"TipoFactor": "Tasa",
				},
			},
			"TotalImpuestosTrasladados": fmt.Sprintf("%.2f", ivaImp),
		},
		"Conceptos": []map[string]any{
			{
				"ClaveProdServ":    "15101514",
				"NoIdentificacion": opts.Payment.GasPump.GasStation.CrePermission,
				"Cantidad":         fmt.Sprintf("%.2f", totalLiter),
				"ClaveUnidad":      "LTR",
				"Unidad":           "Litro",
				"Descripcion": fmt.Sprintf(
					"%s - %s",
					opts.Payment.GasPump.GasStation.CrePermission,
					opts.Payment.FuelType,
				),
				"ValorUnitario": fmt.Sprintf("%.2f", realPricePerLiter),
				"Importe":       fmt.Sprintf("%.2f", subtotal),
				"ObjetoImp":     "02",
				"Impuestos": map[string]any{
					"Traslados": []map[string]any{
						{
							"Base":       fmt.Sprintf("%.2f", base),
							"Impuesto":   "002",
							"TipoFactor": "Tasa",
							"TasaOCuota": "0.160000",
							"Importe":    fmt.Sprintf("%.2f", ivaImp),
						},
					},
				},
			},
		},
	}

	payloadBytes, _ := json.Marshal(payload)
	// fmt.Println(string(payloadBytes))

	client := http.Client{
		Timeout: time.Second * 30,
	}

	url := fmt.Sprintf("%s/v4/cfdi33/issue/json/v1", is.config.ConectiaUrl)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", err
	}

	req.Header.Add("Content-Type", "application/jsontoxml")
	req.Header.Add("Authorization", "Bearer "+is.config.ConectiaToken)
	req.Header.Add("email", opts.Params.Email)

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}

	var data map[string]any

	json.NewDecoder(res.Body).Decode(&data)

	status := data["status"]

	if status != "success" {
		message := data["message"].(string)
		return message, ErrData
	}

	type XMML struct {
		XMLName xml.Name `xml:"TimbreFiscalDigital"`
		ID      string   `xml:"UUID,attr"`
	}

	data2 := data["data"].(map[string]any)
	var dataToDecode XMML
	xmlData := data2["tfd"].(string)

	xml.Unmarshal([]byte(xmlData), &dataToDecode)

	return dataToDecode.ID, nil
}

func (is *invoicingService) ResendInvoice(id string, to string) error {
	body := map[string]any{
		"uuid": id,
		"to":   to,
	}

	payloadBytes, _ := json.Marshal(body)
	client := http.Client{
		Timeout: time.Second * 30,
	}

	url := fmt.Sprintf("%s/comprobante/resendemail", is.config.ConectiaUrlApi)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+is.config.ConectiaToken)

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	var data map[string]any

	json.NewDecoder(res.Body).Decode(&data)

	status := data["status"]

	if status != "success" {
		return ErrNotDocumentFound
	}

	return nil
}

func (is *invoicingService) GetInvoicePDF(id string) (string, error) {
	client := http.Client{
		Timeout: time.Second * 30,
	}

	url := fmt.Sprintf("%s/datawarehouse/v1/live/%s", is.config.ConectiaUrlApi, id)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("Authorization", "Bearer "+is.config.ConectiaToken)

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}

	var data map[string]any

	json.NewDecoder(res.Body).Decode(&data)

	status := data["status"]

	if status != "success" {
		return "", ErrNotDocumentFound
	}

	iData, ok := data["data"].(map[string]any)

	if !ok {
		return "", errors.New("Unable to unparse records")
	}

	records, ok := iData["records"].([]any)

	if !ok {
		return "", errors.New("Unable to unparse records")
	}

	if len(records) == 0 {
		return "", errors.New("No data")
	}

	record := records[len(records)-1].(map[string]any)

	return record["urlPDF"].(string), nil
}
