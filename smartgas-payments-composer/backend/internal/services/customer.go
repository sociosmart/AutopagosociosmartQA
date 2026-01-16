package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"smartgas-payment/config"
	"smartgas-payment/internal/lang"
	"smartgas-payment/internal/schemas"
	"smartgas-payment/internal/utils"
)

//go:generate mockery --name CustomerService --filename=mock_customer.go --inpackage=true
type CustomerService interface {
	Verify(string) (*schemas.Customer, error)
}

func ProvideCustomerService(cfg config.Config) *customerService {
	return &customerService{
		config: cfg,
	}
}

type customerService struct {
	config config.Config
}

func (cs *customerService) Verify(token string) (*schemas.Customer, error) {
	body, _ := json.Marshal(map[string]any{"Token": token})
	url := fmt.Sprintf("%v/rest/clientes?Verifica", cs.config.SocioSmartUrl)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))

	if err != nil {
		// TODO: Logging error in sentry
		return nil, errors.New(lang.InternalServerError)
	}

	respData := make([]map[string]any, 0)

	json.NewDecoder(resp.Body).Decode(&respData)

	if len(respData) < 1 {
		// TODO: Logging sentry error here
		return nil, errors.New(lang.InternalServerError)
	}

	data := respData[0]

	if data["status"] == "error" {
		// TODO: Logging in sentry error
		resultErr, ok := data["result"].(map[string]any)
		if !ok {
			return nil, errors.New(lang.InternalServerError)
		}

		if resultErr["error_id"] == "401" {
			return nil, errors.New(lang.InvalidOrExpiredToken)
		}
		return nil, errors.New(lang.InternalServerError)
	}

	cus := utils.Transform[schemas.Customer](data)

	return cus, nil
}
