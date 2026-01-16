package controllers

import (
	"errors"
	"net/http"
	_ "smartgas-payment/docs"
	"smartgas-payment/internal/dto"
	"smartgas-payment/internal/lang"
	"smartgas-payment/internal/repository"
	"smartgas-payment/internal/schemas"
	"smartgas-payment/internal/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IAUthController interface {
	Login(*gin.Context)
	RefreshToken(*gin.Context)
}

type AuthController struct {
	userRepository repository.UserRepository
}

func ProvideAuthController(userRepository repository.UserRepository) *AuthController {
	return &AuthController{
		userRepository: userRepository,
	}
}

// @BasePath /api/v1
// @Summary Login
// @Description Authorize users
// @Tags Authorization
// @Produce json
// @Accept json
// @Router /api/v1/auth/login [POST]
// @Param Auth body dto.AuthRequestBody true "Auth Body"
// @Success 200 {object} dto.JwtResponse "Access token & Refresh Token"
// @Failure 404 {object} dto.GeneralMessage "User or password incorrect"
// @Failure 500 {object} dto.GeneralMessage "Internal server error"
// @Failure 400 {array} dto.BadRequestMessage "Bad Request, failed on body, form, query..."
func (ac *AuthController) Login(c *gin.Context) {
	var body dto.AuthRequestBody

	if err := c.ShouldBind(&body); err != nil {
		c.JSON(http.StatusBadRequest, utils.MapValidatorError[dto.AuthRequestBody](err))
		return
	}

	user, err := ac.userRepository.GetUserByEmail(body.Email)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, dto.GeneralMessage{Detail: lang.UserOrPasswordIncorrect})
			return
		}
		// Logging Error in sengtry
		utils.TrackError(c, err, nil)
		c.JSON(http.StatusInternalServerError, dto.GeneralMessage{Detail: lang.InternalServerError})
		return
	}

	ok := user.ComparePassword(body.Password)

	if !ok {
		c.JSON(http.StatusNotFound, dto.GeneralMessage{Detail: lang.UserOrPasswordIncorrect})
		return
	}

	tokenClaim := &schemas.JwtClaims{
		Sub: user.ID,
	}

	accessToken, err := tokenClaim.ClaimToken()

	if err != nil {
		// Logging Error in sengtry
		utils.TrackError(c, err, nil)
		c.JSON(http.StatusInternalServerError, dto.GeneralMessage{Detail: lang.InternalServerError})

		return
	}

	refreshToken, err := tokenClaim.ClaimRefreshToken()

	if err != nil {
		// Logging Error in sengtry
		utils.TrackError(c, err, nil)
		c.JSON(http.StatusInternalServerError, dto.GeneralMessage{Detail: lang.InternalServerError})
		return
	}

	c.JSON(http.StatusOK, &dto.JwtResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})

}

// @Summary Refresh Token
// @Description Authorize users
// @Tags Authorization
// @Produce json
// @Accept json
// @Router /api/v1/auth/refresh-token [POST]
// @Param Body body dto.JwtRefreshRequest true "Auth Body"
// @Success 200 {object} dto.JwtResponse "Access token & Refresh Token"
// @Failure 400 {array} dto.BadRequestMessage "Bad Request, failed on body, form, query..."
// @Failure 401 {object} dto.GeneralMessage "Unauthorized, token malformed, token invalid, expired, etc"
func (au *AuthController) RefreshToken(c *gin.Context) {

	var body dto.JwtRefreshRequest

	if err := c.ShouldBind(&body); err != nil {
		c.JSON(http.StatusBadRequest, utils.MapValidatorError[dto.JwtRefreshRequest](err))
		return
	}

	claims, err := utils.ParseJwtToken(body.RefreshToken, true)

	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.GeneralMessage{Detail: err.Error()})
		return
	}

	accessToken, _ := claims.ClaimToken()
	refreshToken, _ := claims.ClaimRefreshToken()

	c.JSON(http.StatusOK, &dto.JwtResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})

}
