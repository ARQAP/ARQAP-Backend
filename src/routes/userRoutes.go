package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ARQAP/ARQAP-Backend/src/controllers"
	"github.com/ARQAP/ARQAP-Backend/src/services"
	"github.com/ARQAP/ARQAP-Backend/src/middleware"
)

func SetupUserRoutes(router *gin.Engine, service *services.UserService) {
    UserController := controllers.NewUserController(service)

    // Public routes
    router.POST("/login", UserController.AuthenticateUser)
	router.POST("/register", UserController.CreateUser)
	router.GET("/users", UserController.GetAllUsers)

	// Protected routes
    user := router.Group("/users")
    user.Use(middleware.AuthMiddleware())
    {
        user.DELETE("/:id", UserController.DeleteUser)
    }
}