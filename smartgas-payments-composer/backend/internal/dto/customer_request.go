package dto

type CustomerDeleteCard struct {
	CardID string `json:"card_id" validate:"required" binding:"required"`
}
