package dto

type SettingCreateBody struct {
	Name  string `json:"name" form:"name" binding:"required" validate:"required" example:"stripe"`
	Value string `json:"value" form:"value" binding:"required" validate:"required" example:"any_value"`
}
