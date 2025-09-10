package routes

import (
	"github.com/ARQAP/ARQAP-Backend/src/controllers"
	"github.com/ARQAP/ARQAP-Backend/src/middleware"
	"github.com/ARQAP/ARQAP-Backend/src/services"
	"github.com/gin-gonic/gin"
)

func SetupArchaeologistRoutes(router *gin.Engine, service *services.ArchaeologistService) {
	archaeologistController := controllers.NewArchaeologistController(service)

	// Protected routes
	archaeologist := router.Group("/archaeologists")
	archaeologist.Use(middleware.AuthMiddleware())
	{
		archaeologist.GET("/", archaeologistController.GetArchaeologists)
		archaeologist.POST("/", archaeologistController.CreateArchaeologist)
		archaeologist.PUT("/:id", archaeologistController.UpdateArchaeologist)
		archaeologist.DELETE("/:id", archaeologistController.DeleteArchaeologist)
	}
}
