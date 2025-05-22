package routes

import (
	"github.com/Luzin7/pcideal-be/internal/http/controllers"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.GET("/ping", controllers.Ping)

	return router
}
