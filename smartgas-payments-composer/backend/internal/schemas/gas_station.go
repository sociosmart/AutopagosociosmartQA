package schemas

type GasStation struct {
	ExternalID    string `json:"Cve_PuntoDeVenta"`
	Name          string `json:"NombreComercial"`
	CrePermission string `json:"Num_PermisoCRE"`
	Ip            string `json:"Vpn,omitempty"`
	Street        string `json:"Calle"`
	ZipCode       string `json:"CP"`
	City          string `json:"Ciudad"`
	State         string `json:"Estado"`
	Neighborhood  string `json:"Colonia"`
	OutsideNumber string `json:"Num_Exterior"`
	Latitude      string `json:"Latitud"`
	Longitude     string `json:"Longitud"`
	LegalNameID   string `json:"Fk_RazonSocial"`
}
