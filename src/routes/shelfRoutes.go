package routes

import (
	"github.com/ARQAP/ARQAP-Backend/src/controllers"
	"github.com/ARQAP/ARQAP-Backend/src/middleware"
	"github.com/ARQAP/ARQAP-Backend/src/services"
	"github.com/gin-gonic/gin"
)

func SetupShelfRoutes(router *gin.Engine, service *services.ShelfService) {
	shelfController := controllers.NewShelfController(service)

	// Protected routes
	shelf := router.Group("/shelves")
	shelf.Use(middleware.AuthMiddleware())
	{
		shelf.GET("/", shelfController.GetAllShelves)
		shelf.GET("/:id", shelfController.GetShelfByID)
		shelf.POST("/", shelfController.CreateShelf)
		shelf.PUT("/:id", shelfController.UpdateShelf)
		shelf.DELETE("/:id", shelfController.DeleteShelf)
	}
}
