package dto

type GasStationGetPathRequest struct {
	ID string `json:"id" uri:"id" binding:"required,uuid4" validate:"required,uuid4" example:"23ae8c18-4d7a-41a3-a148-8ae2d0a75690"`
}

type GasStationCreateRequest struct {
	ExternalID    string `json:"external_id" binding:"required" validate:"required" examplea:"11"`
	Name          string `json:"name" binding:"required,min=3,max=255" validate:"required,min=3,max=255" example:"Guerrero"`
	Ip            string `json:"ip" binding:"required,ipv4" validate:"required,ipv4" example:"192.168.100.100"`
	CrePermission string `json:"cre_permission" binding:"required" validate:"required" example:"PL/01/01..."`
	Active        *bool  `json:"active" binding:"omitempty" validate:"omitempty" example:"true"`
}

// Redundant, but necesary if scaling
type GasStationUpdatePathRequest struct {
	ID string `json:"id" uri:"id" binding:"required,uuid4" validate:"required,uuid4" example:"23ae8c18-4d7a-41a3-a148-8ae2d0a75690"`
}

type GasStationUpdateRequest struct {
	ExternalID    string `json:"external_id" binding:"omitempty" validate:"omitempty" example:"12"`
	Name          string `json:"name" binding:"omitempty,min=3,max=255" validate:"omitempty,min=3,max=255" example:"Guerrero"`
	Ip            string `json:"ip" binding:"omitempty,ipv4" validate:"omitempty,ipv4" example:"192.168.100.100"`
	CrePermission string `json:"cre_permission" binding:"omitempty" validate:"omitempty" example:"PL/01/01..."`
	Active        *bool  `json:"active" binding:"omitempty" validate:"omitempty" example:"true"`
}
