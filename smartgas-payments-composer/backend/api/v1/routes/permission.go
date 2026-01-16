package routes

import (
	"smartgas-payment/api/v1/controllers"
	"smartgas-payment/internal/middlewares"

	"github.com/gin-gonic/gin"
)

type PermissionRoutes struct {
	authMiddleware *middlewares.AuthMiddleware
	controller     controllers.PermissionController
}

func ProvidePermissionRoutes(
	authMiddleware *middlewares.AuthMiddleware,
	controller controllers.PermissionController,
) *PermissionRoutes {
	return &PermissionRoutes{
		authMiddleware: authMiddleware,
		controller:     controller,
	}
}

func (pr *PermissionRoutes) Setup(group *gin.RouterGroup) {
	router := group.Group("/permissions")

	router.Use(pr.authMiddleware.Middleware(middlewares.DefaultAuthmiddlewareOptions()))

	router.GET("/all", pr.controller.ListAll)
	router.GET("/all-groups", pr.controller.ListAllGroups)
}
