package dto

import "github.com/google/uuid"

type GasStationListResponse struct {
	ID            uuid.UUID `json:"id" example:"23ae8c18-4d7a-41a3-a148-8ae2d0a75690"`
	ExternalID    string    `json:"external_id" example:"13"`
	Name          string    `json:"name" example:"Guerrero"`
	Ip            string    `json:"ip" example:"192.168.100.100"`
	CrePermission string    `json:"cre_permission"`
	Active        bool      `json:"active" example:"true"`
}

type GasStationGetResponse struct {
	Name          string `json:"name" example:"Guerrero"`
	Ip            string `json:"ip" example:"192.168.100.100"`
	Active        bool   `json:"active" example:"true"`
	ExternalID    string `json:"external_id" example:"13"`
	CrePermission string `json:"cre_permission"`
}

type GasStationCreateResponse struct {
	ID uuid.UUID `json:"id" example:"23ae8c18-4d7a-41a3-a148-8ae2d0a75690"`
}

type GasStationListAllResponse struct {
	ID   uuid.UUID `json:"id" example:"23ae8c18-4d7a-41a3-a148-8ae2d0a75690"`
	Name string    `json:"name" example:"Guerrero"`
	Ip   string    `json:"ip" example:"192.168.100.100"`
}
