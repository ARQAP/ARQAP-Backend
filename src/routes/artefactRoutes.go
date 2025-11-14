package routes

import (
	"github.com/ARQAP/ARQAP-Backend/src/controllers"
	"github.com/ARQAP/ARQAP-Backend/src/middleware"
	"github.com/ARQAP/ARQAP-Backend/src/services"
	"github.com/gin-gonic/gin"
)

func SetupArtefactRoutes(router *gin.Engine, service *services.ArtefactService) {
	controller := controllers.NewArtefactController(service)

	// Protected routes
	artefactGroup := router.Group("/artefacts")
	artefactGroup.Use(middleware.AuthMiddleware())
	{
		// CRUD
		artefactGroup.GET("", controller.GetAllArtefacts)
		artefactGroup.GET("/:id", controller.GetArtefactByID)
		artefactGroup.POST("/", controller.CreateArtefact)
		artefactGroup.PUT("/:id", controller.UpdateArtefact)
		artefactGroup.DELETE("/:id", controller.DeleteArtefact)

		// Upload
		artefactGroup.POST("/:id/picture", controller.UploadPicture)
		artefactGroup.POST("/:id/historical-record", controller.UploadHistoricalRecord)

		// Serve
		artefactGroup.GET("/:id/picture", controller.ServePicture)
		artefactGroup.GET("/:id/historical-record", controller.ServeHistoricalRecord)
	}
}
