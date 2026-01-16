package middlewares

import (
	"errors"
	"net/http"
	"smartgas-payment/internal/dto"
	"smartgas-payment/internal/enums"
	"smartgas-payment/internal/lang"
	"smartgas-payment/internal/repository"
	"smartgas-payment/internal/utils"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthMiddlewareOptions struct {
	OnlyAdmin          bool
	RequiredPermission enums.Permission
}

func DefaultAuthmiddlewareOptions() AuthMiddlewareOptions {
	return AuthMiddlewareOptions{
		OnlyAdmin: true,
	}
}

type AuthMiddleware struct {
	userRepository repository.UserRepository
}

func ProvideAuthMiddleware(userRepository repository.UserRepository) *AuthMiddleware {
	return &AuthMiddleware{
		userRepository: userRepository,
	}
}

func (am *AuthMiddleware) Middleware(opts AuthMiddlewareOptions) gin.HandlerFunc {
	return func(c *gin.Context) {
		var jwtAuthHeader dto.JwtAuthorizationHeader

		if err := c.ShouldBindHeader(&jwtAuthHeader); err != nil {
			c.JSON(http.StatusUnauthorized, dto.GeneralMessage{Detail: lang.NoAuthorizationHeader})
			c.Abort()
			return
		}

		splittedHeader := strings.Split(jwtAuthHeader.Authorization, " ")

		if len(splittedHeader) < 2 || splittedHeader[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, dto.GeneralMessage{Detail: lang.AuthorizationHeaderMalformed})
			c.Abort()
			return
		}

		claims, err := utils.ParseJwtToken(splittedHeader[1], false)

		if err != nil {
			c.JSON(http.StatusUnauthorized, dto.GeneralMessage{Detail: err.Error()})
			c.Abort()
			return
		}

		user, err := am.userRepository.GetUserByID(claims.Sub)

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusUnauthorized, dto.GeneralMessage{Detail: lang.NotFoundRecord})
				c.Abort()
				return
			}
			// Logging error for when swit list cards fails
			opts := &utils.TrackErrorOpts{
				Tags: map[string]string{"scope": "auth_middleware"},
			}
			utils.TrackError(c, err, opts)
			c.JSON(http.StatusInternalServerError, dto.GeneralMessage{Detail: lang.InternalServerError})
			c.Abort()
			return
		}

		if opts.OnlyAdmin {
			if !*user.IsAdmin {
				c.JSON(http.StatusUnauthorized, dto.GeneralMessage{Detail: lang.NoAdminPermissions})
				c.Abort()
				return
			}
		} else if opts.RequiredPermission != "" {
			if !*user.IsAdmin && !utils.RequiredPermissionInUser(opts.RequiredPermission, user) {
				c.JSON(http.StatusUnauthorized, dto.GeneralMessage{Detail: lang.NotEnoughPermissions + ": " + string(opts.RequiredPermission)})
				c.Abort()
				return
			}
		}

		c.Set("user", user)

		c.Next()
	}
}
