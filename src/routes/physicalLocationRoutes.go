package routes

import (
	"github.com/ARQAP/ARQAP-Backend/src/controllers"
	"github.com/ARQAP/ARQAP-Backend/src/middleware"
	"github.com/ARQAP/ARQAP-Backend/src/services"
	"github.com/gin-gonic/gin"
)

func SetupPhysicalLocationRoutes(router *gin.Engine, service *services.PhysicalLocationService) {
	physicalLocationController := controllers.NewPhysicalLocationController(service)

	// Protected routes
	physicalLocation := router.Group("/physical-locations")
	physicalLocation.Use(middleware.AuthMiddleware())
	{
		physicalLocation.GET("/", physicalLocationController.GetAllPhysicalLocations)
		physicalLocation.GET("/:id", physicalLocationController.GetPhysicalLocationByID)
		physicalLocation.POST("/", physicalLocationController.CreatePhysicalLocation)
		physicalLocation.PUT("/:id", physicalLocationController.UpdatePhysicalLocation)
		physicalLocation.DELETE("/:id", physicalLocationController.DeletePhysicalLocation)
	}
}
