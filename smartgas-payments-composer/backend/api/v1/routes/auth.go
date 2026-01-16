package routes

import (
	"smartgas-payment/api/v1/controllers"

	"github.com/gin-gonic/gin"
)

type AuthRoutes struct {
	controller controllers.IAUthController
}

func ProvideAuthRoutes(controller controllers.IAUthController) *AuthRoutes {
	return &AuthRoutes{
		controller: controller,
	}
}

func (ar *AuthRoutes) Setup(rg *gin.RouterGroup) {
	r := rg.Group("/auth")

	r.POST("/login", ar.controller.Login)
	r.POST("/refresh-token", ar.controller.RefreshToken)
}
