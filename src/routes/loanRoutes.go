package routes

import (
	"github.com/ARQAP/ARQAP-Backend/src/controllers"
	"github.com/ARQAP/ARQAP-Backend/src/middleware"
	"github.com/ARQAP/ARQAP-Backend/src/services"
	"github.com/gin-gonic/gin"
)

func SetupLoanRoutes(router *gin.Engine, service *services.LoanService) {

	loanController := controllers.NewLoanController(service)

	// Protected routes
	mention := router.Group("/loans")
	mention.Use(middleware.AuthMiddleware())
	{
		mention.GET("/", loanController.GetAllLoans)
		mention.GET("/:id", loanController.GetLoanByID)
		mention.POST("/", loanController.CreateLoan)
		mention.PUT("/:id", loanController.UpdateLoan)
		mention.DELETE("/:id", loanController.DeleteLoan)
	}
}