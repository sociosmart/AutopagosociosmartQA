package dto

type AuthRequestBody struct {
	Email    string `json:"email" form:"email" binding:"required,email" validate:"required,email" example:"hello@world.com"`
	Password string `json:"password" form:"password" binding:"required" validate:"required" example:"yourpassword"`
}
