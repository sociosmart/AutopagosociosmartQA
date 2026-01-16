package dto

import "github.com/google/uuid"

type PermissionListAllResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type GroupListAllResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
