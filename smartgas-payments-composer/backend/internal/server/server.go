package server

import (
	"log"
	"smartgas-payment/config"
	"smartgas-payment/internal/injectors"
	"strconv"

	docs "smartgas-payment/docs"

	"github.com/stripe/stripe-go/v72"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
)

func Setup() *gin.Engine {
	a, err := injectors.InitializeServer()

	if err != nil {
		panic(err)
	}

	router := a.Router

	// Registering swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	return router

}

// @title           Smart Gas API Specification
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath  /

// @securitydefinitions.apiKey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and paste the access token

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func RunServer() {

	cfg := config.ConfigSettings

	stripe.Key = cfg.StripeSecretKey

	docs.SwaggerInfo.Host = cfg.Host

	app := Setup()

	if !cfg.Debug {

		log.Printf("Server running on %v\n", cfg.Host)
	}

	log.Fatalln(app.Run(":" + strconv.Itoa(cfg.Port)))

}
