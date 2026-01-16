package dto

type PaginateRequest struct {
	Page  int `json:"page" form:"page" binding:"omitempty,gte=1" validate:"gte=1" minimum:"1"`
	Limit int `json:"limit" form:"limit" binding:"omitempty,gte=1,lte=100" validate:"gte=1,lte=100" minimum:"1" maximum:"100"`
}
