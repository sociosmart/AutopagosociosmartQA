package routes

import (
	"smartgas-payment/api/v1/controllers"
	"smartgas-payment/internal/enums"
	"smartgas-payment/internal/middlewares"

	"github.com/gin-gonic/gin"
)

type CustomerRoutes struct {
	customerAuthMiddleware *middlewares.CustomerAuthMiddleware
	controller             controllers.CustomerController
	authMiddleware         *middlewares.AuthMiddleware
}

func ProvideCustomerRoutes(
	customerAuthMiddleware *middlewares.CustomerAuthMiddleware,
	controller controllers.CustomerController,
	authMiddleware *middlewares.AuthMiddleware,
) *CustomerRoutes {
	return &CustomerRoutes{
		customerAuthMiddleware: customerAuthMiddleware,
		controller:             controller,
		authMiddleware:         authMiddleware,
	}
}

func (cr *CustomerRoutes) Setup(group *gin.RouterGroup) {
	router := group.Group("/customers")

	viewAllCustomersPerms := middlewares.AuthMiddlewareOptions{
		RequiredPermission: enums.ViewAllCustomers,
	}
	router.GET(
		"/payment-methods",
		cr.customerAuthMiddleware.Middleware(),
		cr.controller.ListPaymenthMethods,
	)
	router.DELETE(
		"/payment-methods/:card_id",
		cr.customerAuthMiddleware.Middleware(),
		cr.controller.DeleteCard,
	)
	router.GET("/level", cr.customerAuthMiddleware.Middleware(), cr.controller.GetElegibilityLevel)
	router.GET(
		"/payment-methods-swit",
		cr.customerAuthMiddleware.Middleware(),
		cr.controller.ListPaymenthMethodsSwit,
	)
	router.GET("/all",
		cr.authMiddleware.Middleware(viewAllCustomersPerms),
		cr.controller.ListAll,
	)
}
