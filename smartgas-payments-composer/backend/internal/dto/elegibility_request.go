package dto

type EelegibilityLevelCreateRequest struct {
	Name       string   `json:"name"        validate:"required"       binding:"required"`
	Discount   *float64 `json:"discount"    validate:"required,gte=0" binding:"required,gte=0"`
	MinAmount  *float64 `json:"min_amount"  validate:"required,gte=0" binding:"required,gte=0"`
	MinCharges *int     `json:"min_charges" validate:"required,gte=0" binding:"required,gte=0"`
	Active     *bool    `json:"active"      validate:"omitempty"      binding:"omitempty"      example:"true"`
}

type ElegibilityLevelUpdatesRequest struct {
	Name       string   `json:"name"        validate:"omitempty"       binding:"omitempty"`
	Discount   *float64 `json:"discount"    validate:"omitempty,gte=0" binding:"omitempty,gte=0"`
	MinAmount  *float64 `json:"min_amount"  validate:"omitempty,gte=0" binding:"omitempty,gte=0"`
	MinCharges *int     `json:"min_charges" validate:"omitempty,gte=0" binding:"omitempty,gte=0"`
	Active     *bool    `json:"active"      validate:"omitempty"       binding:"omitempty"       example:"true"`
}

type ElegibilityLevelUpdatePathRequest struct {
	ID string `json:"id" uri:"id" binding:"required,uuid4" validate:"required,uuid4" example:"23ae8c18-4d7a-41a3-a148-8ae2d0a75690"`
}

type ElegibilityCreateCustomerLevelRequest struct {
	LevelID       string `json:"elegibility_level_id" binding:"required,uuid4"             validate:"required,uuid4"             example:"23ae8c18-4d7a-41a3-a148-8ae2d0a75690"`
	CustomerID    string `json:"customer_id"          binding:"required,uuid4"             validate:"required,uuid4"             example:"23ae8c18-4d7a-41a3-a148-8ae2d0a75690"`
	ValidityMonth int    `json:"validity_month"       binding:"required,gte=1,lte=12"      validate:"required,gte=1,lte=12"`
	ValidityYear  int    `json:"validity_year"        binding:"required,gte=2015,lte=2030" validate:"required,gte=2015,lte=2030"`
}

type ElegibilityCustomerLevelUpdatePathRequest struct {
	ID string `json:"id" uri:"id" binding:"required,uuid4" validate:"required,uuid4" example:"23ae8c18-4d7a-41a3-a148-8ae2d0a75690"`
}
type ElegibilityUpdateCustomerLevelRequest struct {
	LevelID       *string `json:"elegibility_level_id" binding:"omitempty,uuid4"             validate:"omitempty,uuid4"             example:"23ae8c18-4d7a-41a3-a148-8ae2d0a75690"`
	CustomerID    *string `json:"customer_id"          binding:"omitempty,uuid4"             validate:"omitempty,uuid4"             example:"23ae8c18-4d7a-41a3-a148-8ae2d0a75690"`
	ValidityMonth *int    `json:"validity_month"       binding:"omitempty,gte=1,lte=12"      validate:"omitempty,gte=1,lte=12"`
	ValidityYear  *int    `json:"validity_year"        binding:"omitempty,gte=2015,lte=2030" validate:"omitempty,gte=2015,lte=2030"`
}
