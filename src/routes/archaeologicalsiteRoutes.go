package routes

import (
	"github.com/ARQAP/ARQAP-Backend/src/controllers"
	"github.com/ARQAP/ARQAP-Backend/src/middleware"
	"github.com/ARQAP/ARQAP-Backend/src/services"
	"github.com/gin-gonic/gin"
)

func SetupArchaeologicalSiteRoutes(router *gin.Engine, service *services.ArchaeologicalSiteService) {
	archaeologicalSiteController := controllers.NewArchaeologicalSiteController(service)

	// Protected routes
	archaeologicalSite := router.Group("/archaeologicalSites")
	archaeologicalSite.Use(middleware.AuthMiddleware())
	{
		archaeologicalSite.GET("/", archaeologicalSiteController.GetArchaeologicalSites)
		archaeologicalSite.POST("/", archaeologicalSiteController.CreateArchaeologicalSite)
		archaeologicalSite.PUT("/:id", archaeologicalSiteController.UpdateArchaeologicalSite)
		archaeologicalSite.DELETE("/:id", archaeologicalSiteController.DeleteArchaeologicalSite)
	}
}
