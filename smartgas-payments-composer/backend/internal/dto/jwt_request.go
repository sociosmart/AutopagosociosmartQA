package dto

type JwtRefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required,jwt" validate:"required,jwt"`
}

type JwtAuthorizationHeader struct {
	Authorization string `header:"Authorization" binding:"required"`
}
