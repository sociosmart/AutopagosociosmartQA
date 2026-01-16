package routes

import (
	"smartgas-payment/api/v1/controllers"
	"smartgas-payment/internal/enums"
	"smartgas-payment/internal/middlewares"

	"github.com/gin-gonic/gin"
)

type GasStationRoutes struct {
	controller     controllers.GasStationController
	authMiddleware *middlewares.AuthMiddleware
}

func ProvideGasStationRoutes(controller controllers.GasStationController, authMiddleware *middlewares.AuthMiddleware) *GasStationRoutes {
	return &GasStationRoutes{
		controller:     controller,
		authMiddleware: authMiddleware,
	}
}

func (gs *GasStationRoutes) Setup(group *gin.RouterGroup) {
	router := group.Group("/gas-stations")

	viewOpts := middlewares.AuthMiddlewareOptions{
		RequiredPermission: enums.ViewGasStations,
	}

	//addOpts := middlewares.AuthMiddlewareOptions{
	//RequiredPermission: enums.AddGasStation,
	//}

	//editOpts := middlewares.AuthMiddlewareOptions{
	//RequiredPermission: enums.EditGasStation,
	//}

	adminOpts := middlewares.DefaultAuthmiddlewareOptions()

	router.GET("", gs.authMiddleware.Middleware(viewOpts), gs.controller.List)
	router.GET("/:id", gs.authMiddleware.Middleware(viewOpts), gs.controller.Get)
	router.PUT("/:id", gs.authMiddleware.Middleware(adminOpts), gs.controller.Update)
	router.POST("", gs.authMiddleware.Middleware(adminOpts), gs.controller.Create)
	router.GET("/all", gs.authMiddleware.Middleware(viewOpts), gs.controller.ListAll)
}
