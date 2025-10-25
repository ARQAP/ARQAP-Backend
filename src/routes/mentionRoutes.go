package routes

import (
	"github.com/ARQAP/ARQAP-Backend/src/controllers"
	"github.com/ARQAP/ARQAP-Backend/src/middleware"
	"github.com/ARQAP/ARQAP-Backend/src/services"
	"github.com/gin-gonic/gin"
)

func SetupMentionRoutes(router *gin.Engine, service *services.MentionService) {
	mentionController := controllers.NewMentionController(service)

	// Protected routes
	mention := router.Group("/mentions")
	mention.Use(middleware.AuthMiddleware())
	{
		mention.GET("/", mentionController.GetMentions)
		mention.GET("/:id", mentionController.GetMentionByID)
		mention.POST("/", mentionController.CreateMention)
		mention.PUT("/:id", mentionController.UpdateMention)
		mention.DELETE("/:id", mentionController.DeleteMention)
	}
}