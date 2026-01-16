package routes

import (
	"smartgas-payment/api/v1/controllers"
	"smartgas-payment/internal/middlewares"

	"github.com/gin-gonic/gin"
)

type SettingRoutes struct {
	authMiddleware *middlewares.AuthMiddleware
	controller     controllers.SettingController
}

func ProvideSettingRoutes(
	authMiddleware *middlewares.AuthMiddleware,
	settingController controllers.SettingController,
) *SettingRoutes {
	return &SettingRoutes{
		authMiddleware: authMiddleware,
		controller:     settingController,
	}
}

func (sr *SettingRoutes) Setup(group *gin.RouterGroup) {
	router := group.Group("/settings")

	router.Use(sr.authMiddleware.Middleware(middlewares.DefaultAuthmiddlewareOptions()))

	router.GET("", sr.controller.List)
	router.POST("", sr.controller.CreateOrUpdate)
	//router.GET("/all-groups", sr.controller.ListAllGroups)
}
