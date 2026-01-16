package dto

type CampaignCreateRequest struct {
	Name         string   `json:"name"         validate:"required"                              binding:"required"`
	Discount     *float64 `json:"discount"     validate:"required,gt=0"                         binding:"required,gt=0"`
	ValidFromStr string   `json:"valid_from"   validate:"required,datetime=2006-01-02 15:04:05" binding:"required,datetime=2006-01-02 15:04:05"`
	ValidToStr   string   `json:"valid_to"     validate:"required,datetime=2006-01-02 15:04:05" binding:"required,datetime=2006-01-02 15:04:05"`
	Active       *bool    `json:"active"       validate:"omitempty"                             binding:"omitempty"                             example:"true"`
	GasStations  *[]struct {
		ID string `json:"id" binding:"required,uuid4"`
	} `json:"gas_stations"                                                  binding:"omitempty,dive"`
}

type CampaignUpdatePathRequest struct {
	ID string `json:"id" uri:"id" binding:"required,uuid4" validate:"required,uuid4" example:"23ae8c18-4d7a-41a3-a148-8ae2d0a75690"`
}

type CampaignDetailPathRequest struct {
	ID string `json:"id" uri:"id" binding:"required,uuid4" validate:"required,uuid4" example:"23ae8c18-4d7a-41a3-a148-8ae2d0a75690"`
}

type CampaignUpdateRequest struct {
	Name               *string  `json:"name"         validate:"omitempty"                              binding:"omitempty"`
	Discount           *float64 `json:"discount"     validate:"omitempty,gt=0"                         binding:"omitempty,gt=0"`
	ValidFromStrUpdate *string  `json:"valid_from"   validate:"omitempty,datetime=2006-01-02 15:04:05" binding:"omitempty,datetime=2006-01-02 15:04:05"`
	ValidToStrUpdate   *string  `json:"valid_to"     validate:"omitempty,datetime=2006-01-02 15:04:05" binding:"omitempty,datetime=2006-01-02 15:04:05"`
	Active             *bool    `json:"active"       validate:"omitempty"                              binding:"omitempty"                              example:"true"`
	GasStations        *[]struct {
		ID string `json:"id" binding:"required,uuid4"`
	} `json:"gas_stations"                                                   binding:"omitempty,dive"`
}
