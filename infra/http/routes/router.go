package routes

import (
	"github.com/Luzin7/pcideal-be/infra/http/controllers"
	"github.com/Luzin7/pcideal-be/infra/http/middlewares"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(partController *controllers.PartController) *gin.Engine {
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
		parts.GET("/", partController.GetAllParts)
		parts.GET("/:id", partController.GetPartByID)
	}

	builds := api.Group("/builds")
	{
		builds.POST("/recommendations", partController.GetBuildRecomendations)
	}

	return router
}
