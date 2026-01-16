package dto

type CustomerAuthorizationHeader struct {
	Authorization string `header:"Authorization" binding:"required"`
}
