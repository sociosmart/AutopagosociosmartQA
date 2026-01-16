package dto

type GeneralMessage struct {
	Detail string `json:"detail"`
}

type BadRequestMessage struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}
