package routes

import (
	"smartgas-payment/api/v1/controllers"
	"smartgas-payment/internal/enums"
	"smartgas-payment/internal/middlewares"

	"github.com/gin-gonic/gin"
)

type ElebilityRoutes struct {
	authMiddleware *middlewares.AuthMiddleware
	controller     controllers.ElegibilityController
}

func ProvideElebilityRoutes(
	authMiddleware *middlewares.AuthMiddleware,
	controller controllers.ElegibilityController,
) *ElebilityRoutes {
	return &ElebilityRoutes{
		authMiddleware: authMiddleware,
		controller:     controller,
	}
}

func (er *ElebilityRoutes) Setup(group *gin.RouterGroup) {
	router := group.Group("/elegibility")

	viewElegibilityLevels := middlewares.AuthMiddlewareOptions{
		RequiredPermission: enums.ViewElegibilityLevels,
	}

	editElegibilityLevel := middlewares.AuthMiddlewareOptions{
		RequiredPermission: enums.EditElegibilityLevel,
	}

	addElegibilityLevel := middlewares.AuthMiddlewareOptions{
		RequiredPermission: enums.AddElegibilityLevel,
	}

	viewCustomerLevels := middlewares.AuthMiddlewareOptions{
		RequiredPermission: enums.ViewCustomerLevels,
	}

	viewAllLevelsPerms := middlewares.AuthMiddlewareOptions{
		RequiredPermission: enums.ViewAllElegebilityLevels,
	}

	updateCustomerLevelPerms := middlewares.AuthMiddlewareOptions{
		RequiredPermission: enums.EditCustomerLevel,
	}

	createCustomerLevelPerms := middlewares.AuthMiddlewareOptions{
		RequiredPermission: enums.AddCustomerLevel,
	}

	router.GET(
		"/levels",
		er.authMiddleware.Middleware(viewElegibilityLevels),
		er.controller.LevelList,
	)
	router.POST(
		"/levels",
		er.authMiddleware.Middleware(addElegibilityLevel),
		er.controller.CreateLevel,
	)
	router.PUT(
		"/levels/:id",
		er.authMiddleware.Middleware(editElegibilityLevel),
		er.controller.UpdateLevel,
	)
	router.GET(
		"/customers/levels",
		er.authMiddleware.Middleware(viewCustomerLevels),
		er.controller.CustomerLevelList,
	)
	router.GET(
		"/levels/all",
		er.authMiddleware.Middleware(viewAllLevelsPerms),
		er.controller.LevelListAll,
	)

	router.PUT(
		"/customers/levels/:id",
		er.authMiddleware.Middleware(updateCustomerLevelPerms),
		er.controller.UpdateCustomerLevel,
	)

	router.POST(
		"/customers/levels",
		er.authMiddleware.Middleware(createCustomerLevelPerms),
		er.controller.CreateCustomerLevel,
	)
}
