package dto

type LevelListResponse struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	Discount   float64 `json:"discount"`
	MinCharges int     `json:"min_charges"`
	MinAmount  float64 `json:"min_amount"`
	Active     bool    `json:"active"`
}

type CustomerLevelListResponse struct {
	ID            string `json:"id"`
	LevelID       string `json:"elegibility_level_id"`
	CustomerID    string `json:"customer_id"`
	ValidityMonth int    `json:"validity_month"`
	ValidityYear  int    `json:"validity_year"`
}

type LevelCreateResponse struct {
	ID string `json:"id"`
}

type LevelListAllResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type CustomerLevelCreateResponse struct {
	ID string `json:"id"`
}
