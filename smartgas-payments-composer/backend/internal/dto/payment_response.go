package dto

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type PaymentWebsocketNotification struct {
	Status string `json:"status"`
}

type PaymentCrateIntentResponse struct {
	ClientSecret string    `json:"client_secret,omitempty" description:"Client secret used to pay the rerquested charge of fuel"`
	Amount       float64   `json:"amount"                  description:"the amount that is gonna be charged"`
	TotalLiter   float64   `json:"total_liter"             description:"The total liter that are gonna be charged"`
	ID           uuid.UUID `json:"id"`
}

type PaymentCrateIntentOperationResponse struct {
	Amount     float64   `json:"amount"      description:"the amount that is gonna be charged"`
	TotalLiter float64   `json:"total_liter" description:"The total liter that are gonna be charged"`
	ID         uuid.UUID `json:"id"`
}

type PaymentListResponse struct {
	ID                    uuid.UUID `json:"id"`
	ExternalTransactionID string    `json:"external_transaction_id"`
	PaymentProvider       string    `json:"payment_provider"`
	Amount                float32   `json:"amount"`
	TotalLiter            float32   `json:"total_liter"`
	Price                 float64   `json:"price"`
	ChargeType            string    `json:"charge_type"`
	FuelType              string    `json:"fuel_type"`
	RefundedAmount        float32   `json:"refunded_amount"`
	RealAmountReported    float32   `json:"real_amount_reported"`
	DiscountPerLiter      float64   `json:"discount_per_liter"`
	ChargeFee             float32   `json:"charge_fee"`
	GMPoints              float32   `json:"gm_points"`
	Status                string    `json:"status"`
	CreatedAt             time.Time `json:"created_at"`
	Customer              struct {
		ID             uuid.UUID `json:"id"`
		FirstName      string    `json:"first_name"`
		FirstLastName  string    `json:"first_last_name"`
		SecondLastName string    `json:"second_last_name"`
	} `json:"customer"`
	GasPump struct {
		ID         uuid.UUID `json:"id"`
		Number     string    `json:"number"`
		GasStation struct {
			ID   uuid.UUID `json:"id"`
			Name string    `json:"name"`
		} `json:"gas_station"`
	} `json:"gas_pump"`
	Events []struct {
		Type      string    `json:"type"`
		CreatedAt time.Time `json:"created_at"`
	} `json:"events"`
}

type PaymentDetailCustomer struct {
	Amount              float32   `json:"amount"`
	TotalLiter          float32   `json:"total_liter"`
	Price               float64   `json:"price"`
	FuelType            string    `json:"fuel_type"`
	CreatedAt           time.Time `json:"created_at"`
	RefundedAmount      float32   `json:"refunded_amount"`
	RealAmountReported  float32   `json:"real_amount_reported"`
	ChargeFee           float32   `json:"charge_fee"`
	GMPoints            float32   `json:"gm_points"`
	RealDiscountApplied float64   `json:"real_discount_applied"`
	GasPump             struct {
		Number     string `json:"number"`
		GasStation struct {
			Name string `json:"name"`
		} `json:"gas_station"`
	} `json:"gas_pump"`
	Events []struct {
		Type      string    `json:"type"`
		CreatedAt time.Time `json:"created_at"`
	} `json:"events"`
}

func (plr *PaymentListResponse) MarshalJSON() ([]byte, error) {
	if plr.Events == nil {
		plr.Events = make([]struct {
			Type      string    "json:\"type\""
			CreatedAt time.Time "json:\"created_at\""
		}, 0)
	}

	return json.Marshal(*plr)
}

type GetPaymentProviderResponse struct {
	PaymentProvider string `json:"payment_provider"`
	Business        string `json:"business,omitempty"`
	Token           string `json:"token,omitempty"`
	CustomerID      string `json:"customer_id,omitempty"`
}

type SignInvoiceResponse struct {
	UUID string `json:"uuid"`
}

type GetInvoicePDFResponse struct {
	UrlPDF string `json:"urlPDF"`
}
