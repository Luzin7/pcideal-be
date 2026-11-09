package routes

import (
	"net/http"
	"os"
	"time"

	"github.com/Luzin7/pcideal-be/infra/http/controllers"
	partControllers "github.com/Luzin7/pcideal-be/infra/http/controllers/part"
	"github.com/Luzin7/pcideal-be/infra/http/middlewares"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Dependencies struct {
	GetAllPartsController  *partControllers.GetAllPartsController
	GetPartByIDController  *partControllers.GetPartByIDController
	GetBuildRecsController *partControllers.GetBuildRecommendationsController
	HealthController       *controllers.HealthController
	ClientIPResolver       middlewares.ClientIPResolver
}

func SetupRouter(deps Dependencies) (*gin.Engine, error) {
	allowedOrigins, err := middlewares.ParseAllowedOrigins(os.Getenv("CORS_ALLOWED_ORIGINS"))
	if err != nil {
		return nil, err
	}

	router := gin.Default()

	if err := middlewares.ConfigureTrustedProxies(router, os.Getenv("TRUSTED_PROXIES")); err != nil {
		return nil, err
	}

	router.Use(middlewares.IPLoggingMiddleware(deps.ClientIPResolver))
	router.Use(middlewares.SecurityHeadersMiddleware())

	config := cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodHead, http.MethodOptions},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "x-api-key"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}
	router.Use(cors.New(config))

	router.GET("/health", deps.HealthController.Handle)

	api := router.Group("/api")

	parts := api.Group("/parts")
	{
		parts.GET("/", deps.GetAllPartsController.Handle)
		parts.GET("/:id", deps.GetPartByIDController.Handle)
	}

	builds := api.Group("/builds")
	{
		builds.POST("/recommendations", deps.GetBuildRecsController.Handle)
	}

	return router, nil
}
