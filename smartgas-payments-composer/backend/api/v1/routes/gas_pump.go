package routes

import (
	"smartgas-payment/api/v1/controllers"
	"smartgas-payment/internal/enums"
	"smartgas-payment/internal/middlewares"

	"github.com/gin-gonic/gin"
)

type GasPumpRoutes struct {
	controller         controllers.GasPumpController
	authMiddleware     *middlewares.AuthMiddleware
	customerMiddleware *middlewares.CustomerAuthMiddleware
}

func ProvideGasPumpRoutes(
	controller controllers.GasPumpController,
	authMiddleware *middlewares.AuthMiddleware,
	customerMiddleware *middlewares.CustomerAuthMiddleware,
) *GasPumpRoutes {
	return &GasPumpRoutes{
		controller:         controller,
		authMiddleware:     authMiddleware,
		customerMiddleware: customerMiddleware,
	}
}

func (gp *GasPumpRoutes) Setup(group *gin.RouterGroup) {
	router := group.Group("/gas-pumps")

	viewOpts := middlewares.AuthMiddlewareOptions{
		RequiredPermission: enums.ViewGasPumps,
	}

	//addOpts := middlewares.AuthMiddlewareOptions{
	//RequiredPermission: enums.AddGasPump,
	//}

	//editOpts := middlewares.AuthMiddlewareOptions{
	//RequiredPermission: enums.EditGasPump,
	//}

	adminOpts := middlewares.DefaultAuthmiddlewareOptions()

	router.GET("", gp.authMiddleware.Middleware(viewOpts), gp.controller.List)
	router.GET("/:id", gp.authMiddleware.Middleware(viewOpts), gp.controller.Get)
	router.PUT("/:id", gp.authMiddleware.Middleware(adminOpts), gp.controller.Update)
	router.POST("", gp.authMiddleware.Middleware(adminOpts), gp.controller.Create)
	router.GET("/:id/customer", gp.customerMiddleware.Middleware(), gp.controller.GetDetailForCustomer)

}
