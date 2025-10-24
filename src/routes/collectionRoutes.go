package routes

import (
	"github.com/ARQAP/ARQAP-Backend/src/controllers"
	"github.com/ARQAP/ARQAP-Backend/src/middleware"
	"github.com/ARQAP/ARQAP-Backend/src/services"
	"github.com/gin-gonic/gin"
)

func SetupCollectionRoutes(router *gin.Engine, service *services.CollectionService) {
	collectionController := controllers.NewCollectionController(service)

	// Protected routes
	collection := router.Group("/collections")
	collection.Use(middleware.AuthMiddleware())
	{
		collection.GET("/", collectionController.GetCollections)
		collection.POST("/", collectionController.CreateCollection)
		collection.PUT("/:id", collectionController.UpdateCollection)
		collection.DELETE("/:id", collectionController.DeleteCollection)
	}
}
