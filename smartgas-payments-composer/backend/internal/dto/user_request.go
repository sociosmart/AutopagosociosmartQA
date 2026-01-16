package dto

type UserCreateRequest struct {
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=3,max=100"`
	IsAdmin     *bool  `json:"is_admin" binding:"required"`
	Active      *bool  `json:"active" binding:"required"`
	FirstName   string `json:"first_name" binding:"required,min=2,max=100"`
	LastName    string `json:"last_name" binding:"required,min=2,max=100"`
	Permissions []struct {
		ID string `json:"id" binding:"required,uuid4"`
	} `json:"permissions" binding:"required,dive"`
	Groups []struct {
		ID string `json:"id" binding:"required,uuid4"`
	} `json:"groups" binding:"required,dive"`
	GasStations *[]struct {
		ID string `json:"id" binding:"required,uuid4"`
	} `json:"gas_stations" binding:"omitempty,dive"`
}

type UserUpdatePathRequest struct {
	ID string `json:"id" uri:"id" binding:"required,uuid4" validate:"required,uuid4" example:"23ae8c18-4d7a-41a3-a148-8ae2d0a75690"`
}

type UserGetDetailPathRequest struct {
	ID string `json:"id" uri:"id" binding:"required,uuid4" validate:"required,uuid4" example:"23ae8c18-4d7a-41a3-a148-8ae2d0a75690"`
}

type UserUpdateRequest struct {
	Password    *string `json:"password" binding:"omitempty,min=3,max=100"`
	IsAdmin     *bool   `json:"is_admin" binding:"omitempty"`
	Active      *bool   `json:"active" binding:"omitempty"`
	FirstName   *string `json:"first_name" binding:"omitempty,min=2,max=100"`
	LastName    *string `json:"last_name" binding:"omitempty,min=2,max=100"`
	Permissions *[]struct {
		ID string `json:"id" binding:"required,uuid4"`
	} `json:"permissions" binding:"omitempty,dive"`
	Groups *[]struct {
		ID string `json:"id" binding:"required,uuid4"`
	} `json:"groups" binding:"omitempty,dive"`
	GasStations *[]struct {
		ID string `json:"id" binding:"required,uuid4"`
	} `json:"gas_stations" binding:"omitempty,dive"`
}
