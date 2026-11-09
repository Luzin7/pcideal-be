package routes

import (
	"os"

	"github.com/Luzin7/pcideal-be/infra/http/controllers"
	partControllers "github.com/Luzin7/pcideal-be/infra/http/controllers/part"
	"github.com/Luzin7/pcideal-be/infra/http/middlewares"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(
	getAllPartsController *partControllers.GetAllPartsController,
	getPartByIDController *partControllers.GetPartByIDController,
	getBuildRecsController *partControllers.GetBuildRecommendationsController,
	healthController *controllers.HealthController,
) (*gin.Engine, error) {
	router := gin.Default()

	if err := middlewares.ConfigureTrustedProxies(router, os.Getenv("TRUSTED_PROXIES")); err != nil {
		return nil, err
	}

	router.Use(middlewares.IPLoggingMiddleware())
	router.Use(middlewares.SecurityHeadersMiddleware())

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000", "https://www.pcideal.online", "https://pcideal.online"}
	router.Use(cors.New(config))

	router.GET("/health", healthController.Handle)

	api := router.Group("/api")

	parts := api.Group("/parts")
	{
		parts.GET("/", getAllPartsController.Handle)
		parts.GET("/:id", getPartByIDController.Handle)
	}

	builds := api.Group("/builds")
	{
		builds.POST("/recommendations", getBuildRecsController.Handle)
	}

	return router, nil
}
