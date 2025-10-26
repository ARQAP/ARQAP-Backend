package routes

import (
	"github.com/ARQAP/ARQAP-Backend/src/controllers"
	"github.com/ARQAP/ARQAP-Backend/src/middleware"
	"github.com/ARQAP/ARQAP-Backend/src/services"
	"github.com/gin-gonic/gin"
)

func SetupRequesterRoutes(router *gin.Engine, service *services.RequesterService) {

	requesterController := controllers.NewRequesterController(service)

	// Protected routes
	requester := router.Group("/requesters")
	requester.Use(middleware.AuthMiddleware())
	{
		requester.GET("/", requesterController.GetAllRequesters)
		requester.GET("/:id", requesterController.GetRequesterByID)
		requester.POST("/", requesterController.CreateRequester)
		requester.PUT("/:id", requesterController.UpdateRequester)
		requester.DELETE("/:id", requesterController.DeleteRequester)
	}
}