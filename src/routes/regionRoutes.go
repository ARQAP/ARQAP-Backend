package routes

import (
	"github.com/ARQAP/ARQAP-Backend/src/controllers"
	"github.com/ARQAP/ARQAP-Backend/src/middleware"
	"github.com/ARQAP/ARQAP-Backend/src/services"
	"github.com/gin-gonic/gin"
)

func SetupRegionRoutes(router *gin.Engine, service *services.RegionService) {
	regionController := controllers.NewRegionController(service)

	// Protected routes
	region := router.Group("/regions")
	region.Use(middleware.AuthMiddleware())
	{
		region.GET("/", regionController.GetRegions)
		region.POST("/", regionController.CreateRegion)
		region.PUT("/:id", regionController.UpdateRegion)
		region.DELETE("/:id", regionController.DeleteRegion)
	}
}
