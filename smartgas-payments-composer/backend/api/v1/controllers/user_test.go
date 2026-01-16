package controllers_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"smartgas-payment/internal/injectors"
	"smartgas-payment/internal/models"
	"smartgas-payment/internal/repository"
	"smartgas-payment/internal/schemas"
	"smartgas-payment/internal/utils"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

type testControllerSuite struct {
	suite.Suite
	userRepositoryMock *repository.MockUserRepository
	testRequest        *utils.TestRequest
	id                 uuid.UUID
}

func (suite *testControllerSuite) SetupSuite() {
	setup, err := injectors.InitializeServerWithMocks()

	suite.Nil(err, "Expected nil on init server mocking")

	suite.userRepositoryMock = setup.UserRepositoryMock
	suite.testRequest = &utils.TestRequest{
		Router: setup.Router,
	}

	suite.id = uuid.New()
	userLogin := &models.User{
		ID:        suite.id,
		FirstName: "Smart",
		Email:     "mesterlum@hotmail.com",
		LastName:  "Gas",
		Active:    utils.BoolAddr(true),
		IsAdmin:   utils.BoolAddr(true),
	}

	userLogin.HashPassword()

	suite.userRepositoryMock.On("GetUserByID", suite.id).Return(userLogin, nil)
}

func (suite *testControllerSuite) TestMe() {
	url := "/api/v1/users/me"
	testcases := []struct {
		Name               string
		Url                string
		ExpectedStatusCode int
		ExpectedData       any
		Pre                func()
		Post               func()
	}{
		{
			Name:               "TestMeController_GetMe",
			Url:                url,
			ExpectedStatusCode: http.StatusOK,
			Pre: func() {
				user := schemas.JwtClaims{
					Sub: suite.id,
				}

				token, _ := user.ClaimToken()

				suite.testRequest.SetBearerToken("Bearer " + token)
			},
			Post: func() {
				suite.testRequest.SetBearerToken("")
			},
		},
	}
	for i := range testcases {
		tc := testcases[i]

		t := suite.Suite.T()
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			suite.userRepositoryMock.Test(t)
			tc.Pre()

			res := suite.testRequest.Get(tc.Url, nil)

			if tc.ExpectedStatusCode != 0 {
				suite.Equal(tc.ExpectedStatusCode, res.Code, fmt.Sprintf("Expected status code %v, got %v", tc.ExpectedStatusCode, res.Code))
			}

			if tc.ExpectedData != nil {
				data, _ := json.Marshal(tc.ExpectedData)
				suite.Equal(string(data), res.Body.String(), fmt.Sprintf("Expected body %v, got %v", string(data), res.Body.String()))
			}

			tc.Post()
		})
	}

}

func TestMeController(t *testing.T) {
	suite.Run(t, new(testControllerSuite))
}
