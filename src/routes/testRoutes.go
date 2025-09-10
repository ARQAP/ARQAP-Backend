package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ARQAP/ARQAP-Backend/src/controllers"
	"github.com/ARQAP/ARQAP-Backend/src/services"
)

func SetupTestRoutes(router *gin.Engine, service *services.TestService) {
	testController := controllers.NewTestController(service)

	router.GET("/test", testController.GetTests)
	router.POST("/test", testController.CreateTests)

}