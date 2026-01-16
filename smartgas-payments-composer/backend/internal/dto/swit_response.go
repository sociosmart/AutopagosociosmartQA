package dto

type SwitGeneralResponse struct {
	Status string `json:"status"`
	Result any    `json:"result"`
	Errors any    `json:"errors"`
}
