package routes

import (
	partControllers "github.com/Luzin7/pcideal-be/infra/http/controllers/part"
	"github.com/Luzin7/pcideal-be/infra/http/middlewares"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(
	getAllPartsController *partControllers.GetAllPartsController,
	getPartByIDController *partControllers.GetPartByIDController,
	getBuildRecsController *partControllers.GetBuildRecommendationsController,
) *gin.Engine {
	router := gin.Default()
	router.SetTrustedProxies(nil)

	router.Use(middlewares.IPLoggingMiddleware())
	router.Use(middlewares.SecurityHeadersMiddleware())

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000", "https://www.pcideal.online", "https://pcideal.online"}
	router.Use(cors.New(config))

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

	return router
}
