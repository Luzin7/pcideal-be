package routes

import (
	"github.com/Luzin7/pcideal-be/internal/http/controllers"
	"github.com/gin-gonic/gin"
)

func SetupRouter(partController *controllers.PartController) *gin.Engine {
	router := gin.Default()

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

	// auth := api.Group("/auth")
	// {
	// 	auth.POST("/login", authController.Login)
	// 	auth.POST("/register", authController.Register)
	// }

	return router
}
