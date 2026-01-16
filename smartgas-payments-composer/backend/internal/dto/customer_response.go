package dto

import "smartgas-payment/internal/schemas"

type ListCustomerPaymentMethodResponse []schemas.PaymentMethod

type ListAllCustomersResponse struct {
	ID             string `json:"id"`
	FirstName      string `json:"first_name"`
	FirstLastName  string `json:"first_last_name"`
	SecondLastName string `json:"second_last_name"`
}

type CustomerLevelAssignedResponse struct {
	ID            string  `json:"id"`
	Name          string  `json:"name"`
	Discount      float64 `json:"discount"`
	LevelsEnabled bool    `json:"levels_enabled"`
}
