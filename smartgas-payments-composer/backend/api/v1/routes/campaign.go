package routes

import (
	"smartgas-payment/api/v1/controllers"
	"smartgas-payment/internal/enums"
	"smartgas-payment/internal/middlewares"

	"github.com/gin-gonic/gin"
)

type CampaignRoutes struct {
	controller     controllers.CampaignController
	authMiddleware *middlewares.AuthMiddleware
}

func ProvideCampaingRoutes(
	controller controllers.CampaignController,
	authMiddleware *middlewares.AuthMiddleware,
) *CampaignRoutes {
	return &CampaignRoutes{
		authMiddleware: authMiddleware,
		controller:     controller,
	}
}

func (cr *CampaignRoutes) Setup(group *gin.RouterGroup) {
	router := group.Group("/campaigns")

	// defaultOpts := middlewares.DefaultAuthmiddlewareOptions()

	viewCampaignPerm := middlewares.AuthMiddlewareOptions{
		RequiredPermission: enums.ViewCampaigns,
	}

	editCampaignPerm := middlewares.AuthMiddlewareOptions{
		RequiredPermission: enums.EditCampaign,
	}

	addCampaignPerm := middlewares.AuthMiddlewareOptions{
		RequiredPermission: enums.AddCampaign,
	}

	router.GET("", cr.authMiddleware.Middleware(viewCampaignPerm), cr.controller.List)
	router.POST("", cr.authMiddleware.Middleware(addCampaignPerm), cr.controller.Create)
	router.PUT("/:id", cr.authMiddleware.Middleware(editCampaignPerm), cr.controller.Update)
	router.GET(
		"/:id",
		cr.authMiddleware.Middleware(viewCampaignPerm),
		cr.controller.GetCampaignDetail,
	)
}
