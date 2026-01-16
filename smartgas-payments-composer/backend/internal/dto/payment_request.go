package dto

type CreatePaymentIntentRequest struct {
	FuelType        string  `json:"fuel_type"        validate:"required,oneof=regular premium diesel"             binding:"required,oneof=regular premium diesel"`
	Amount          float32 `json:"amount"           validate:"required_if=ChargeType by_total,omitempty,gte=10"  binding:"required_if=ChargeType by_total,omitempty,gte=10"`
	TotalLiter      float32 `json:"total_liter"      validate:"required_if=ChargeType by_liter,omitempty,gte=0.5" binding:"required_if=ChargeType by_liter,omitempty,gte=0.5"`
	ChargeType      string  `json:"charge_type"      validate:"required,oneof=by_liter by_total"                  binding:"required,oneof=by_liter by_total"`
	GasPumpID       string  `json:"gas_pump_id"      validate:"required,uuid4"                                    binding:"required,uuid4"                                    example:"23ae8c18-4d7a-41a3-a148-8ae2d0a75690"`
	PaymentProvider string  `json:"payment_provider" validate:"required,oneof=stripe swit debit"                  binding:"required,oneof=stripe swit debit"`
	SourceID        string  `json:"source_id"                                                                     binding:"required_if=PaymentProvider swit"`
	Last4           string  `json:"last_4"                                                                        binding:"required_if=PaymentProvider swit"`
	Cvv             string  `json:"cvv"                                                                           binding:"required_if=PaymentProvider swit"`
}

type CreatePaymentIntentOperationRequest struct {
	FuelType           string  `json:"fuel_type"            validate:"required,oneof=regular premium diesel" binding:"required,oneof=regular premium diesel"`
	ChargeType         string  `json:"charge_type"          validate:"required,oneof=customer card_key"      binding:"required,oneof=customer card_key"`
	Amount             float32 `json:"amount"               validate:"required,gte=10"                       binding:"required,gte=10"`
	PumpNumber         string  `json:"pump_number"          validate:"required,len=2"                        binding:"required,len=2"`
	ExternalCustomerID string  `json:"external_customer_id" validate:"required_if=ChargeType customer"       binding:"required_if=ChargeType customer"`
	CardKey            string  `json:"card_key"             validate:"required_if=ChargeType card_key"       binding:"required_if=ChargeType card_key"`
}

type AddEventPathRequest struct {
	ID string `json:"id" uri:"id" binding:"required" validate:"required" example:"AP_23ae8c18-4d7a-41a3-a148-8ae2d0a75690"`
}

type AddEventBodyRequest struct {
	Type          string  `json:"type"           binding:"required,oneof=serving serving_paused served" validate:"required,oneof=serving serving_paused served"`
	AmountCharged float32 `json:"amount_charged" binding:"required_if=Type served"                      validate:"required_if=Type served"`
}

type PaymentDetailCustomerPath struct {
	ID string `json:"id" uri:"id" binding:"required,uuid4" validate:"required,uuid4" example:"23ae8c18-4d7a-41a3-a148-8ae2d0a75690"`
}

type PaymentDetailCustomerPathWebsocket struct {
	ID string `json:"id" uri:"id" binding:"required,uuid4" validate:"required,uuid4" example:"23ae8c18-4d7a-41a3-a148-8ae2d0a75690"`
}

type SignInvoiceRequest struct {
	Rfc           string `json:"rfc"            binding:"required" validate:"required"`
	Email         string `json:"email"          binding:"required" validate:"required"`
	RazonSocial   string `json:"razon_social"   binding:"required" validate:"required"`
	CP            string `json:"cp"             binding:"required" validate:"required"`
	UsoCFDI       string `json:"uso_cfdi"       binding:"required" validate:"required"`
	RegimenFiscal string `json:"regimen_fiscal" binding:"required" validate:"required"`
}

type ResendInvoiceRequest struct {
	Email string `json:"email" binding:"required,email" validate:"required,email"`
}

type DoPaymentActionRequest struct {
	Action string `json:"action" binding:"required,oneof=refund preset" validate:"required,oneof=refund preset"`
}

type DoPaymentActionRequestPath struct {
	ID string `json:"id" uri:"id" binding:"required,uuid4" validate:"required,uuid4" example:"23ae8c18-4d7a-41a3-a148-8ae2d0a75690"`
}
