package middlewares_test

import (
	"encoding/json"
	"errors"
	"fmt"
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

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type testAuthMiddlewareSuit struct {
	suite.Suite
	userRepositoryMock *repository.MockUserRepository
	testRequest        *utils.TestRequest
	id                 uuid.UUID
	noAdminId          uuid.UUID
	notFoundId         uuid.UUID
	serverErrorId      uuid.UUID
}

func (suite *testAuthMiddlewareSuit) SetupSuite() {
	setup, err := injectors.InitializeServerWithMocks()

	suite.Nil(err, "Expected nil on init server mocking")

	suite.userRepositoryMock = setup.UserRepositoryMock
	suite.testRequest = &utils.TestRequest{
		Router: setup.Router,
	}

	suite.id = uuid.New()
	suite.noAdminId = uuid.New()
	suite.notFoundId = uuid.New()
	suite.serverErrorId = uuid.New()

	userLogin := &models.User{
		ID:        suite.id,
		FirstName: "Smart",
		Email:     "mesterlum@hotmail.com",
		LastName:  "Gas",
		Active:    utils.BoolAddr(true),
		IsAdmin:   utils.BoolAddr(true),
	}

	noAdminLogin := userLogin
	noAdminLogin.ID = suite.noAdminId
	noAdminLogin.IsAdmin = utils.BoolAddr(false)

	userLogin.HashPassword()

	suite.userRepositoryMock.On("GetUserByID", suite.id).Return(userLogin, nil)
	suite.userRepositoryMock.On("GetUserByID", suite.noAdminId).Return(noAdminLogin, nil)
	suite.userRepositoryMock.On("GetUserByID", suite.notFoundId).Return(nil, gorm.ErrRecordNotFound)
	suite.userRepositoryMock.On("GetUserByID", suite.serverErrorId).Return(nil, errors.New(lang.InternalServerError))
}

func (suite *testAuthMiddlewareSuit) TestMe() {
	url := "/api/v1/users/me"
	testcases := []struct {
		Name               string
		Url                string
		ExpectedStatusCode int
		ExpectedErr        any
		ExpectedData       any
		Pre                func()
		Post               func()
	}{
		{
			Name:               "TestAuthMiddleware_NoHeaderToken",
			Url:                url,
			ExpectedErr:        dto.GeneralMessage{Detail: lang.NoAuthorizationHeader},
			ExpectedStatusCode: http.StatusUnauthorized,
			Pre:                func() {},
			Post:               func() {},
		},
		{
			Name:               "TestAuthMiddleware_NoUUIDSetted",
			Url:                url,
			ExpectedStatusCode: http.StatusUnauthorized,
			Pre: func() {
				suite.testRequest.SetBearerToken("Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c")
			},
			Post: func() {
				suite.testRequest.SetBearerToken("")
			},
		},
		{
			Name:               "TestAuthMiddleware_TokenMalformed",
			Url:                url,
			ExpectedErr:        dto.GeneralMessage{Detail: lang.AuthorizationHeaderMalformed},
			ExpectedStatusCode: http.StatusUnauthorized,
			Pre: func() {
				suite.testRequest.SetBearerToken("Bearer")
			},
			Post: func() {
				suite.testRequest.SetBearerToken("")
			},
		},
		{
			Name:               "TestAuthMiddleware_TokenInternalServerError",
			Url:                url,
			ExpectedErr:        dto.GeneralMessage{Detail: lang.InternalServerError},
			ExpectedStatusCode: http.StatusInternalServerError,
			Pre: func() {
				claims := schemas.JwtClaims{
					Sub: suite.serverErrorId,
				}

				token, _ := claims.ClaimToken()
				suite.testRequest.SetBearerToken("Bearer " + token)
			},
			Post: func() {
				suite.testRequest.SetBearerToken("")
			},
		},
		{
			Name:               "TestAuthMiddleware_TokenMalformedNoToken",
			Url:                url,
			ExpectedErr:        dto.GeneralMessage{Detail: lang.AuthorizationHeaderMalformed},
			ExpectedStatusCode: http.StatusUnauthorized,
			Pre: func() {
				suite.testRequest.SetBearerToken("Bearer")
			},
			Post: func() {
				suite.testRequest.SetBearerToken("")
			},
		},
		{
			Name:               "TestAuthMiddleware_TokenIncorrect",
			Url:                url,
			ExpectedStatusCode: http.StatusUnauthorized,
			Pre: func() {
				suite.testRequest.SetBearerToken("Bearer asdadsad")
			},
			Post: func() {
				suite.testRequest.SetBearerToken("")
			},
		},
		{
			Name:               "TestAuthMiddleware_TokenWithUserUnathorized",
			Url:                url,
			ExpectedStatusCode: http.StatusUnauthorized,
			Pre: func() {
				user := schemas.JwtClaims{
					Sub: uuid.New(),
				}

				token, _ := user.ClaimToken()
				suite.testRequest.SetBearerToken("Bearer" + token)
			},
			Post: func() {
				suite.testRequest.SetBearerToken("")
			},
		},
		{
			Name:               "TestAuthMiddleware_TokenExpired",
			Url:                url,
			ExpectedStatusCode: http.StatusUnauthorized,
			//ExpectedErr:        dto.GeneralMessage{Detail: jwt.ErrTokenExpired.Error()},
			Pre: func() {
				user := schemas.JwtClaims{
					Sub: uuid.New(),
				}
				user.ExpiresAt = jwt.NewNumericDate(time.Now().Add(-1 * time.Minute))

				token, _ := user.ClaimToken()

				suite.testRequest.SetBearerToken("Bearer " + token)

			},
			Post: func() {
				suite.testRequest.SetBearerToken("")
			},
		},
		//{
		//Name:               "TestAuthMiddleware_NoAdmin",
		//Url:                url,
		//ExpectedStatusCode: http.StatusUnauthorized,
		//ExpectedErr:        dto.GeneralMessage{Detail: lang.NoAdminPermissions},
		//Pre: func() {
		//user := schemas.JwtClaims{
		//Sub: suite.noAdminId,
		//}

		//token, _ := user.ClaimToken()

		//suite.testRequest.SetBearerToken("Bearer " + token)

		//},
		//Post: func() {
		//suite.testRequest.SetBearerToken("")
		//},
		//},
		{
			Name:               "TestAuthMiddleware_NotFoundUser",
			Url:                url,
			ExpectedStatusCode: http.StatusUnauthorized,
			ExpectedErr:        dto.GeneralMessage{Detail: lang.NotFoundRecord},
			Pre: func() {
				user := schemas.JwtClaims{
					Sub: suite.notFoundId,
				}

				token, _ := user.ClaimToken()

				suite.testRequest.SetBearerToken("Bearer " + token)

			},
			Post: func() {
				suite.testRequest.SetBearerToken("")
			},
		},
	}

	t := suite.Suite.T()

	for i := range testcases {
		tc := testcases[i]

		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			suite.userRepositoryMock.Test(t)
			tc.Pre()

			res := suite.testRequest.Get(tc.Url, nil)

			if tc.ExpectedStatusCode != 0 {
				suite.Equal(tc.ExpectedStatusCode, res.Code, fmt.Sprintf("Expected status code %v, got %v", tc.ExpectedStatusCode, res.Code))
			}

			if tc.ExpectedErr != nil {
				data, _ := json.Marshal(tc.ExpectedErr)
				suite.Equal(string(data), res.Body.String(), fmt.Sprintf("Expected body %v, got %v", string(data), res.Body.String()))
			}

			tc.Post()
		})
	}

}

func TestAuthMiddleware(t *testing.T) {
	suite.Run(t, new(testAuthMiddlewareSuit))
}
