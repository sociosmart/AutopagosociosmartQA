package app

import (
	"fmt"
	"net/http"
	"smartgas-payment/api/v1/routes"
	"smartgas-payment/config"
	"smartgas-payment/internal/dto"
	"smartgas-payment/internal/lang"
	"strings"

	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// @Summary Show health check
// @Description Authorize users
// @Tags HealthCheck
// @Produce json
// @Router /healthcheck [GET]
// @Success 200 {object} dto.GeneralMessage "Shows healthy if server is running"
func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, dto.GeneralMessage{Detail: lang.Healthy})
}

func configureCors(cfg config.Config, engine *gin.Engine) {
	allowdHosts := strings.Split(cfg.AllowedHosts, " ")
	engine.Use(cors.New(
		cors.Config{
			AllowOrigins: allowdHosts,
			AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
			AllowHeaders: []string{
				"Content-Type",
				"Access-Control-Allow-Origin",
				"Access-Control-Allow-Headers",
				"Authorization",
				"X-EMPLOYEE-ID",
				"X-EMPLOYEE-NIP",
				"X-GAS-STATION-ID",
			},
			ExposeHeaders: []string{"Content-Length"},
		},
	))
}

func ProvideGinApp(c config.Config, v1Routes routes.Routes) *gin.Engine {
	if !c.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	var engine *gin.Engine

	if gin.Mode() == "test" {
		engine = gin.New()
		engine.Use(gin.Recovery())
	} else {
		engine = gin.Default()
	}

	if gin.Mode() != "test" {
		rate := 0.0

		if !c.Debug {
			rate = 1.0
		}
		// Setup sentry
		if err := sentry.Init(sentry.ClientOptions{
			Dsn:           c.SentryDsn,
			EnableTracing: !c.Debug,
			// Set TracesSampleRate to 1.0 to capture 100%
			// of transactions for performance monitoring.
			// We recommend adjusting this value in production,
			TracesSampleRate:   rate,
			ProfilesSampleRate: rate,
			Environment:        c.Environment,
			Debug:              c.Debug,
			AttachStacktrace:   true,
		}); err != nil {
			fmt.Printf("Sentry initialization failed: %v", err)
		}

		engine.Use(sentrygin.New(sentrygin.Options{Repanic: true}))
	}

	configureCors(c, engine)

	engine.SetTrustedProxies(strings.Split(c.TrustedProxies, " "))

	// Register healthcheck
	engine.GET("/healthcheck", healthCheck)

	// Registering api v1
	v1Group := engine.Group("/api/v1")
	for _, route := range v1Routes {
		route.Setup(v1Group)
	}

	return engine
}
