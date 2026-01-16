package dto

type GasPumpGetPathRequest struct {
	ID string `json:"id" uri:"id" binding:"required,uuid4" validate:"required,uuid4" example:"23ae8c18-4d7a-41a3-a148-8ae2d0a75690"`
}

type GasPumpCreateRequest struct {
	ExternalID   string  `json:"external_id" binding:"required" validate:"required" examplea:"11"`
	RegularPrice float64 `json:"regular_price" binding:"omitempty,gte=0" validate:"omitempty,gte=0"`
	PremiumPrice float64 `json:"premium_price" binding:"omitempty,gte=0" validate:"omitempty,gte=0"`
	DieselPrice  float64 `json:"diesel_price" binding:"omitempty,gte=0" validate:"omitempty,gte=0"`
	Number       string  `json:"number" binding:"required,min=2" validate:"required,min=2" example:"01, 02, 03..."`
	GasStationID string  `json:"gas_station_id" binding:"required,uuid4" validate:"required,uuid4" example:"23ae8c18-4d7a-41a3-a148-8ae2d0a75690"`
	Active       *bool   `json:"active" binding:"omitempty" validate:"omitempty" example:"true"`
}

// Redundant, but necesary if scaling
type GasPumpUpdatePathRequest struct {
	ID string `json:"id" uri:"id" binding:"required,uuid4" validate:"required,uuid4" example:"23ae8c18-4d7a-41a3-a148-8ae2d0a75690"`
}

type GasPumpGetDetailForCustomerPathRequest struct {
	ID string `json:"id" uri:"id" binding:"required,uuid4" validate:"required,uuid4" example:"23ae8c18-4d7a-41a3-a148-8ae2d0a75690"`
}

type GasPumpUpdateRequest struct {
	ExternalID   string   `json:"external_id" binding:"omitempty" validate:"omitempty" example:"12"`
	RegularPrice *float64 `json:"regular_price" binding:"omitempty,gte=0" validate:"omitempty,gte=0"`
	PremiumPrice *float64 `json:"premium_price" binding:"omitempty,gte=0" validate:"omitempty,gte=0"`
	DieselPrice  *float64 `json:"diesel_price" binding:"omitempty,gte=0" validate:"omitempty,gte=0"`
	Number       *string  `json:"number" binding:"omitempty,min=2" validate:"omitempty,min=2" example:"01, 02, 03..."`
	GasStationID *string  `json:"gas_station_id" binding:"omitempty,uuid4" validate:"omitempty,uuid4" example:"23ae8c18-4d7a-41a3-a148-8ae2d0a75690"`
	Active       *bool    `json:"active" binding:"omitempty" validate:"omitempty" example:"true"`
}
