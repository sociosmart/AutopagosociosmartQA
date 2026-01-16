package schemas

import (
	"time"
)

type SetGasPump struct {
	Status  uint8  `json:"Estatus"`
	Message string `json:"Mensaje"`
}

type AcumPoints struct {
	TransID          string `json:"N_Transaccion"`
	CrePermission    string `json:"Cod_Gasolinero"`
	GasPumpNumber    string `json:"PosicionCarga"`
	Date             string `json:"Fecha"`
	Time             string `json:"Hora"`
	ProductID        string `json:"Id_Producto"`
	Amount           string `json:"Monto"`
	TotalLiter       string `json:"Cantidad"`
	PhoneNumber      string `json:"N_Cliente"`
	Start            string `json:"EmpiezaDespacho"`
	Status           string `json:"Estatus"`
	ClientRegistered string `json:"N_ClienteRegistrado"`
}

func (p *AcumPoints) FillDateTime() {
	loc, _ := time.LoadLocation("America/Mazatlan")
	now := time.Now().In(loc)

	p.Time = now.Format("15:04:05")
	p.Date = now.Format("2006-01-02")
}
