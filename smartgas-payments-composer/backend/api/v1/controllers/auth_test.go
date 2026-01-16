package controllers_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"smartgas-payment/internal/dto"
	"smartgas-payment/internal/injectors"
	"smartgas-payment/internal/lang"
	"smartgas-payment/internal/models"
	"smartgas-payment/internal/repository"
	"smartgas-payment/internal/schemas"
	"smartgas-payment/internal/utils"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

const (
	adminEmail    = "payments@smartgas.com"
	nonAdminEmail = "payments-nonadmin@smartgas.com"
)

type authControllerTestSuite struct {
	suite.Suite
	testRequest        *utils.TestRequest
	userMockRepository *repository.MockUserRepository
	userId             uuid.UUID
}

func (suite *authControllerTestSuite) SetupSuite() {
	setup, _ := injectors.InitializeServerWithMocks()

	suite.testRequest = &utils.TestRequest{
		Router: setup.Router,
	}

	suite.userMockRepository = setup.UserRepositoryMock

	suite.userId = uuid.New()
	user := &models.User{
		ID:        suite.userId,
		Email:     adminEmail,
		Password:  "password",
		FirstName: "Smart",
		LastName:  "Gas",
		IsAdmin:   utils.BoolAddr(true),
	}
	user.HashPassword()

	nonAdminUser := &models.User{
		ID:        uuid.New(),
		Email:     nonAdminEmail,
		Password:  "nonadmin-password",
		FirstName: "Smart",
		LastName:  "Gas",
		IsAdmin:   utils.BoolAddr(true),
	}

	nonAdminUser.HashPassword()

	suite.userMockRepository.On("GetUserByEmail", adminEmail).Return(user, nil)
	suite.userMockRepository.On("GetUserByEmail", nonAdminUser).Return(nonAdminUser, nil)
	suite.userMockRepository.On("GetUserByEmail", "internal@error.com").Return(nil, errors.New(lang.InternalServerError))
	suite.userMockRepository.On("GetUserByEmail", mock.AnythingOfType("string")).Return(nil, gorm.ErrRecordNotFound)

	suite.userMockRepository.On("GetUserByID", suite.userId).Return(user, nil)
	suite.userMockRepository.On("GetUserByID", mock.AnythingOfType("uuid.UUID")).Return(nil, gorm.ErrRecordNotFound)

}

func (suite *authControllerTestSuite) TestLogin() {
	url := "/api/v1/auth/login"

	validate := validator.New()

	err := validate.Struct(dto.AuthRequestBody{})
	badRequestExpected := utils.MapValidatorError[dto.AuthRequestBody](err)
	testcases := []struct {
		Name               string
		Url                string
		Body               any
		ExpectedStatusCode int
		ExpectedResponse   any
	}{
		{
			Name:               "AuthControllerTest_LoginBadRequest",
			Url:                url,
			Body:               dto.AuthRequestBody{},
			ExpectedStatusCode: http.StatusBadRequest,
			ExpectedResponse:   badRequestExpected,
		},
		{
			Name: "AuthControllerTest_LoginInternalServerError",
			Url:  url,
			Body: dto.AuthRequestBody{
				Email:    "internal@error.com",
				Password: "anything",
			},
			ExpectedStatusCode: http.StatusInternalServerError,
			ExpectedResponse:   dto.GeneralMessage{Detail: lang.InternalServerError},
		},
		{
			Name: "AuthControllerTest_LoginUserEmailIncorrect",
			Url:  url,
			Body: dto.AuthRequestBody{
				Email:    "notfoundmail@gmail.com",
				Password: "anything",
			},
			ExpectedStatusCode: http.StatusNotFound,
			ExpectedResponse:   dto.GeneralMessage{Detail: lang.UserOrPasswordIncorrect},
		},
		{
			Name:               "AuthControllerTest_LoginUserPasswordIncorrect",
			Url:                url,
			ExpectedStatusCode: http.StatusNotFound,
			ExpectedResponse:   dto.GeneralMessage{Detail: lang.UserOrPasswordIncorrect},
			Body: dto.AuthRequestBody{
				Email:    adminEmail,
				Password: "anything",
			},
		},
		{
			Name:               "AuthControllerTest_LoginUserNotAdmin",
			Url:                url,
			ExpectedStatusCode: http.StatusNotFound,
			ExpectedResponse:   dto.GeneralMessage{Detail: lang.UserOrPasswordIncorrect},
			Body: dto.AuthRequestBody{
				Email:    nonAdminEmail,
				Password: "nonadmin-password",
			},
		},
		{
			Name:               "AuthControllerTest_LoginAdmin",
			Url:                url,
			ExpectedStatusCode: http.StatusOK,
			Body: dto.AuthRequestBody{
				Email:    adminEmail,
				Password: "password",
			},
		},
	}

	t := suite.T()
	for _, tc := range testcases {
		t.Run(tc.Name, func(t *testing.T) {
			suite.userMockRepository.Test(t)

			res := suite.testRequest.Post(tc.Url, tc.Body)

			if tc.ExpectedStatusCode != 0 {
				suite.Equal(tc.ExpectedStatusCode, res.Code, utils.PrintExpectedValues(tc.ExpectedStatusCode, res.Code))
			}

			if tc.ExpectedResponse != nil {
				expected, _ := json.Marshal(tc.ExpectedResponse)
				suite.Equal(string(expected), res.Body.String(), utils.PrintExpectedValues(string(expected), res.Body.String()))
			}
		})
	}
}

func (suite *authControllerTestSuite) TestRrefreshLogin() {
	url := "/api/v1/auth/refresh-token"
	validate := validator.New()
	err := validate.Struct(dto.JwtRefreshRequest{})

	badRequestExpected := utils.MapValidatorError[dto.JwtRefreshRequest](err)

	claims := &schemas.JwtClaims{
		Sub: suite.userId,
	}

	validRefreshToken, _ := claims.ClaimRefreshToken()
	// Simulating that is a valid token but with wrong signature
	invalidRefreshToken, _ := claims.ClaimToken()

	// Expired token
	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(-1 * time.Minute))
	expiredRefreshToken, _ := claims.ClaimRefreshToken()
	testcases := []struct {
		Name               string
		Url                string
		Body               any
		ExpectedStatusCode int
		ExpectedResponse   any
	}{
		{
			Name:               "AuthControllerTest_RefreshTokenBadRequest",
			Url:                url,
			Body:               dto.JwtRefreshRequest{},
			ExpectedStatusCode: http.StatusBadRequest,
			ExpectedResponse:   badRequestExpected,
		},
		{
			Name:               "AuthControllerTest_RefreshTokenExpiredToken",
			Url:                url,
			Body:               dto.JwtRefreshRequest{RefreshToken: expiredRefreshToken},
			ExpectedStatusCode: http.StatusUnauthorized,
		},
		{
			Name:               "AuthControllerTest_RefreshInvalidToken",
			Url:                url,
			Body:               dto.JwtRefreshRequest{RefreshToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c"},
			ExpectedStatusCode: http.StatusUnauthorized,
		},
		{
			Name:               "AuthControllerTest_RefreshTokenInvalidSignature",
			Url:                url,
			Body:               dto.JwtRefreshRequest{RefreshToken: invalidRefreshToken},
			ExpectedStatusCode: http.StatusUnauthorized,
		},
		{
			Name:               "AuthControllerTest_RefreshValidToken",
			Url:                url,
			Body:               dto.JwtRefreshRequest{RefreshToken: validRefreshToken},
			ExpectedStatusCode: http.StatusOK,
		},
	}

	t := suite.Suite.T()
	for i := range testcases {
		tc := testcases[i]
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			suite.userMockRepository.Test(t)

			res := suite.testRequest.Post(tc.Url, tc.Body)

			if tc.ExpectedStatusCode != 0 {
				suite.Equal(tc.ExpectedStatusCode, res.Code, utils.PrintExpectedValues(tc.ExpectedStatusCode, res.Code))
			}

			if tc.ExpectedResponse != nil {
				expected, _ := json.Marshal(tc.ExpectedResponse)
				suite.Equal(string(expected), res.Body.String(), utils.PrintExpectedValues(string(expected), res.Body.String()))
			}
		})
	}
}

func TestAuthController(t *testing.T) {
	suite.Run(t, new(authControllerTestSuite))
}
