package routes

import (
	"smartgas-payment/api/v1/controllers"
	"smartgas-payment/internal/middlewares"

	"github.com/gin-gonic/gin"
)

type UserRoutes struct {
	controller     controllers.UserController
	authMiddleware *middlewares.AuthMiddleware
}

func ProvideUserRoutes(userController controllers.UserController, authMiddleware *middlewares.AuthMiddleware) *UserRoutes {
	return &UserRoutes{
		controller:     userController,
		authMiddleware: authMiddleware,
	}
}

func (ur UserRoutes) Setup(group *gin.RouterGroup) {
	router := group.Group("/users")

	defaultOpts := middlewares.DefaultAuthmiddlewareOptions()
	router.GET("/me", ur.authMiddleware.Middleware(middlewares.AuthMiddlewareOptions{}), ur.controller.Me)
	router.GET("", ur.authMiddleware.Middleware(defaultOpts), ur.controller.List)
	router.POST("", ur.authMiddleware.Middleware(defaultOpts), ur.controller.Create)
	router.PUT("/:id", ur.authMiddleware.Middleware(defaultOpts), ur.controller.Update)
	router.GET("/:id", ur.authMiddleware.Middleware(defaultOpts), ur.controller.GetUserDetail)
}
