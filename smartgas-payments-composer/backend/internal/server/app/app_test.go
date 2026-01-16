package app_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"smartgas-payment/internal/dto"
	"smartgas-payment/internal/injectors"
	"smartgas-payment/internal/lang"
	"smartgas-payment/internal/utils"
	"testing"

	"github.com/stretchr/testify/suite"
)

type appTestSuite struct {
	suite.Suite
	testRequest *utils.TestRequest
}

func (suite *appTestSuite) SetupSuite() {
	setup, _ := injectors.InitializeServerWithMocks()

	suite.testRequest = &utils.TestRequest{
		Router: setup.Router,
	}
}

func (suite *appTestSuite) TestHealthcheck() {
	url := "/healthcheck"
	testcases := []struct {
		Name               string
		Url                string
		ExpectedStatusCode int
		ExpectedResponse   any
	}{
		{
			Name:               "TestHealthCheck",
			Url:                url,
			ExpectedStatusCode: http.StatusOK,
			ExpectedResponse:   dto.GeneralMessage{Detail: lang.Healthy},
		},
	}

	t := suite.Suite.T()

	for i := range testcases {

		tc := testcases[i]
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			res := suite.testRequest.Get(tc.Url, nil)

			if tc.ExpectedStatusCode != 0 {
				suite.Equal(tc.ExpectedStatusCode, res.Code, fmt.Sprintf("Expected status code %v, got %v", tc.ExpectedStatusCode, res.Code))
			}

			if tc.ExpectedResponse != nil {
				data, _ := json.Marshal(tc.ExpectedResponse)
				suite.Equal(string(data), res.Body.String(), fmt.Sprintf("Expected body %v, got %v", string(data), res.Body.String()))
			}
		})
	}

}

func TestServer(t *testing.T) {
	suite.Run(t, new(appTestSuite))
}
