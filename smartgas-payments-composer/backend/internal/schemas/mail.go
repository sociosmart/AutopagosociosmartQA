package schemas

import (
	"fmt"
	"smartgas-payment/internal/models"
	"strings"
)

type FuelRequest struct {
	CustomerName   string
	TransactionID  string
	CustomerID     string
	Date           string
	FuelType       string
	GasPump        string
	GasStation     string
	Amount         float32
	RefundedAmount float32
	TotalLiter     float32
	RealAmount     float32
	Error          bool
}

func (fr *FuelRequest) FillData(payment *models.Payment) {

	fr.CustomerName = fmt.Sprintf("%v %v %v", payment.Customer.FirstName, payment.Customer.FirstLastName, payment.Customer.SecondLastName)
	fr.TransactionID = payment.ID.String()
	fr.CustomerID = payment.Customer.PhoneNumber
	fr.Date = payment.CreatedAt.Format("01-02-2006")
	fr.FuelType = strings.Title(payment.FuelType)
	fr.GasPump = payment.GasPump.Number
	fr.GasStation = payment.GasPump.GasStation.Name
	fr.Amount = payment.Amount
	fr.RefundedAmount = payment.RefundedAmount
	fr.RealAmount = payment.RealAmountReported
}
