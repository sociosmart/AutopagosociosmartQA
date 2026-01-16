package schemas

import (
	"fmt"
	"strings"
)

type Customer struct {
	ExternalID     string `json:"ID"`
	FirstName      string `json:"Nombre"`
	FirstLastName  string `json:"Ap_Paterno"`
	SecondLastName string `json:"Ap_Materno"`
	PhoneNumber    string `json:"Num_Celular"`
	Status         string `json:"Estado"`
	Email          string `json:"Correo"`
	SwitCustomerID string `json:"TokenSwit"`
}

type SwitSource struct {
	ID         string `json:"card_id"`
	Last4      string `json:"last4"`
	Brand      string `json:"brand"`
	IsLastUsed bool   `json:"isLastUsed"`
}

func (c Customer) Fullname() string {
	return strings.TrimSpace(
		fmt.Sprintf("%v %v %v", c.FirstName, c.FirstLastName, c.SecondLastName),
	)
}
