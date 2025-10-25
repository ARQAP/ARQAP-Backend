package routes

import (
	"github.com/ARQAP/ARQAP-Backend/src/controllers"
	"github.com/ARQAP/ARQAP-Backend/src/middleware"
	"github.com/ARQAP/ARQAP-Backend/src/services"
	"github.com/gin-gonic/gin"
)

func SetupINPLClassifiersRoutes(router *gin.Engine, service *services.INPLService) {
	inplController := controllers.NewINPLClassifierController(service)

	inplClassifiers := router.Group("/inplClassifiers")
	inplClassifiers.Use(middleware.AuthMiddleware())
	{
		inplClassifiers.GET("/", inplController.GetAllINPLClassifiers)
		inplClassifiers.GET("/:id", inplController.GetINPLClassifierByID)
		inplClassifiers.POST("/", inplController.CreateINPLClassifier)
		inplClassifiers.PUT("/:id", inplController.UpdateINPLClassifier)
		inplClassifiers.DELETE("/:id", inplController.DeleteINPLClassifier)

		inplClassifiers.POST("/:id/fichas", inplController.AddFichasToINPLClassifier)
		inplClassifiers.GET("/:id/fichas", inplController.ListFichasByINPLClassifier)
	}

	inplFichas := router.Group("/inplFichas")
	inplFichas.Use(middleware.AuthMiddleware())
	{
		inplFichas.PUT("/:id", inplController.ReplaceFicha)
		inplFichas.DELETE("/:id", inplController.DeleteFicha)
		inplFichas.GET("/:id/download", inplController.DownloadFicha)
	}
}
