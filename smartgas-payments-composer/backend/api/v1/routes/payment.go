package routes

import (
	"smartgas-payment/api/v1/controllers"
	"smartgas-payment/internal/enums"
	"smartgas-payment/internal/middlewares"

	"github.com/gin-gonic/gin"
)

type PaymentRoutes struct {
	customerAuthMiddleware *middlewares.CustomerAuthMiddleware
	authMiddleware         *middlewares.AuthMiddleware
	securityMiddleware     *middlewares.SecurityMiddleware
	controller             controllers.PaymentController
}

func ProvidePaymenRoutes(
	customerAuthMiddleware *middlewares.CustomerAuthMiddleware,
	controller controllers.PaymentController,
	authMiddleware *middlewares.AuthMiddleware,
	securityMiddleware *middlewares.SecurityMiddleware,
) *PaymentRoutes {
	return &PaymentRoutes{
		controller:             controller,
		customerAuthMiddleware: customerAuthMiddleware,
		authMiddleware:         authMiddleware,
		securityMiddleware:     securityMiddleware,
	}
}

func (pr *PaymentRoutes) Setup(group *gin.RouterGroup) {
	router := group.Group("/payments")

	viewOpts := middlewares.AuthMiddlewareOptions{
		RequiredPermission: enums.ViewPayments,
	}

	canDoPaymentActionOpts := middlewares.AuthMiddlewareOptions{
		RequiredPermission: enums.CanDoActionsPayments,
	}

	router.POST(
		"/create-intent",
		pr.customerAuthMiddleware.Middleware(),
		pr.controller.CreateIntent,
	)
	router.POST("/stripe-webhook", pr.controller.StripeWebhook)
	router.GET("", pr.authMiddleware.Middleware(viewOpts), pr.controller.List)
	router.POST("/:id/events", pr.securityMiddleware.Middleware(), pr.controller.AddEvent)
	router.GET(
		"/:id/customer-detail",
		pr.customerAuthMiddleware.Middleware(),
		pr.controller.GetByIDForCustomer,
	)
	router.GET("/:id/customer-detail-ws", pr.controller.PaymentNotifierWS)
	router.GET(
		"/provider",
		pr.customerAuthMiddleware.Middleware(),
		pr.controller.GetPaymentProvider,
	)
	router.POST("/invoicing/:id", pr.customerAuthMiddleware.Middleware(), pr.controller.SignInvoice)
	router.POST(
		"/invoicing/:id/resend",
		pr.customerAuthMiddleware.Middleware(),
		pr.controller.ResendInvoice,
	)
	router.GET(
		"/invoicing/:id/pdf",
		pr.customerAuthMiddleware.Middleware(),
		pr.controller.GetInvoicePDF,
	)
	router.POST(
		"/actions/:id",
		pr.authMiddleware.Middleware(canDoPaymentActionOpts),
		pr.controller.DoPaymentAction,
	)
	router.POST(
		"/create-intent-operation",
		pr.securityMiddleware.SmartGasEmployeeMiddleware(),
		pr.controller.CreateIntentOperation,
	)
}
