package dto

import (
	"github.com/google/uuid"
)

type GasPumpListResponse struct {
	ID           uuid.UUID `json:"id"             example:"23ae8c18-4d7a-41a3-a148-8ae2d0a75690"`
	Active       bool      `json:"active"         example:"true"`
	Number       string    `json:"number"`
	ExternalID   string    `json:"external_id"`
	RegularPrice float64   `json:"regular_price"`
	PremiumPrice float64   `json:"premium_price"`
	DieselPrice  float64   `json:"diesel_price"`
	GasStationID uuid.UUID `json:"gas_station_id"`
}

type GasPumpGetResponse struct {
	Active       bool    `json:"active"        example:"true"`
	Number       string  `json:"number"`
	ExternalID   string  `json:"external_id"`
	RegularPrice float64 `json:"regular_price"`
	PremiumPrice float64 `json:"premium_price"`
	DieselPrice  float64 `json:"diesel_price"`
	GasStation   struct {
		ID         uuid.UUID `json:"id"`
		Name       string    `json:"name"`
		ExternalID string    `json:"external_id"`
	} `json:"gas_station"`
}

type GasPumpGetDetailForCustomerResponse struct {
	Number       string  `json:"number"`
	RegularPrice float64 `json:"regular_price"`
	PremiumPrice float64 `json:"premium_price"`
	DieselPrice  float64 `json:"diesel_price"`
	GasStation   struct {
		Name          string `json:"name"`
		Street        string `json:"street"`
		ZipCode       string `json:"zip_code"`
		City          string `json:"city"`
		State         string `json:"state"`
		Neighborhood  string `json:"neighborhood"`
		OutsideNumber string `json:"outside_number"`
	} `json:"gas_station"`
	DiscountType string `json:"discount_type"`
	Campaign     *struct {
		Name     string  `json:"name"`
		Discount float64 `json:"discount"`
	} `json:"campaign"`
}

type GasPumpCreateResponse struct {
	ID uuid.UUID `json:"id" example:"23ae8c18-4d7a-41a3-a148-8ae2d0a75690"`
}
