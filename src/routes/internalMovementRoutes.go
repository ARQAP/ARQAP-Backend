package routes

import (
	"github.com/ARQAP/ARQAP-Backend/src/controllers"
	"github.com/ARQAP/ARQAP-Backend/src/middleware"
	"github.com/ARQAP/ARQAP-Backend/src/services"
	"github.com/gin-gonic/gin"
)

func SetupInternalMovementRoutes(router *gin.Engine, service *services.InternalMovementService) {
	internalMovementController := controllers.NewInternalMovementController(service)

	// Protected routes
	internalMovementGroup := router.Group("/internal-movements")
	internalMovementGroup.Use(middleware.AuthMiddleware())
	{
		internalMovementGroup.GET("/", internalMovementController.GetAllInternalMovements)
		internalMovementGroup.GET("/:id", internalMovementController.GetInternalMovementByID)
		internalMovementGroup.GET("/artefact/:artefactId", internalMovementController.GetInternalMovementsByArtefactID)
		internalMovementGroup.GET("/artefact/:artefactId/active", internalMovementController.GetActiveInternalMovementByArtefactID)
		internalMovementGroup.POST("/", internalMovementController.CreateInternalMovement)
		internalMovementGroup.PUT("/:id", internalMovementController.UpdateInternalMovement)
		internalMovementGroup.DELETE("/:id", internalMovementController.DeleteInternalMovement)
	}
}

