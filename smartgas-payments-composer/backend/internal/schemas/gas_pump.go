package schemas

type GasPump struct {
	ExternalID   string   `json:"Cve_Id"`
	Number       string   `json:"Bomba"`
	RegularPrice *float64 `json:"Ppregular,omitempty"`
	PremiumPrice *float64 `json:"Ppremium,omitempty"`
	DieselPrice  *float64 `json:"Ppdiesel,omitempty"`
	Active       *bool    `json:"EstatusBomba"`
}
