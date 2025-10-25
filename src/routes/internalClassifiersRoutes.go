package routes

import (
	"github.com/ARQAP/ARQAP-Backend/src/controllers"
	"github.com/ARQAP/ARQAP-Backend/src/middleware"
	"github.com/ARQAP/ARQAP-Backend/src/services"
	"github.com/gin-gonic/gin"
)

func SetupInternalClassifiersRoutes(router *gin.Engine, service *services.InternalClassifierService) {
	internalClassifierController := controllers.NewInternalClassifierController(service)

	// Protected routes
	internalClassifier := router.Group("/internalClassifiers")
	internalClassifier.Use(middleware.AuthMiddleware())
	{
		internalClassifier.GET("/", internalClassifierController.GetAllInternalClassifiers)
		internalClassifier.GET("/:id", internalClassifierController.GetInternalClassifierByID)
		internalClassifier.POST("/", internalClassifierController.CreateInternalClassifier)
		internalClassifier.PUT("/:id", internalClassifierController.UpdateInternalClassifier)
		internalClassifier.DELETE("/:id", internalClassifierController.DeleteInternalClassifier)
	}
}
