package routes

import (
	"smartgas-payment/api/v1/controllers"
	"smartgas-payment/internal/enums"
	"smartgas-payment/internal/middlewares"

	"github.com/gin-gonic/gin"
)

type SynchronizationRoute struct {
	controller     controllers.SynchronizationController
	authMiddleware *middlewares.AuthMiddleware
}

func ProvideSynchronizationRoutes(
	controller controllers.SynchronizationController,
	authMiddleware *middlewares.AuthMiddleware,
) *SynchronizationRoute {
	return &SynchronizationRoute{
		authMiddleware: authMiddleware,
		controller:     controller,
	}
}

func (sr *SynchronizationRoute) Setup(group *gin.RouterGroup) {
	router := group.Group("/synchronizations")

	//router.Use(sr.authMiddleware.Middleware(middlewares.DefaultAuthmiddlewareOptions()))

	viewSyncsOpts := middlewares.AuthMiddlewareOptions{
		RequiredPermission: enums.ViewSynchronizations,
	}

	addSyncOpts := middlewares.AuthMiddlewareOptions{
		RequiredPermission: enums.AddSynchronization,
	}

	router.GET("", sr.authMiddleware.Middleware(viewSyncsOpts), sr.controller.List)
	router.GET("/:id/details", sr.authMiddleware.Middleware(viewSyncsOpts), sr.controller.ListDetails)
	router.GET("/last", sr.authMiddleware.Middleware(viewSyncsOpts), sr.controller.GetLastSync)
	router.POST("/now", sr.authMiddleware.Middleware(addSyncOpts), sr.controller.SyncNow)
}
