package schemas

type PaymentMethod struct {
	ID   string `json:"id"`
	Card struct {
		Last4 string `json:"last_4"`
		Brand string `json:"brand"`
	} `json:"card"`
	IsLastUsed bool `json:"is_last_used"`
}
