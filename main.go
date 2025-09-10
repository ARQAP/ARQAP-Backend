package main

import (
	"log"
	"os"

	"github.com/ARQAP/ARQAP-Backend/src/db"
	"github.com/ARQAP/ARQAP-Backend/src/models"
	"github.com/ARQAP/ARQAP-Backend/src/routes"
	"github.com/ARQAP/ARQAP-Backend/src/services"
	"github.com/gin-gonic/gin"
)

func main() {

	// Database connection
	db, err := db.Connect()
	if err != nil {
		log.Fatalf("Error connecting to database: %v\n", err)
	}

	// Auto-migrate models
	if err := db.AutoMigrate(&models.TestModel{}, &models.UserModel{}); err != nil {
		log.Fatalf("Error during auto-migration: %v\n", err)
	}

	// Port and host setup
	host := os.Getenv("SERVER_HOST")
	if host == "" {
		host = ":8080"
	}

	// Gin router setup
	router := gin.Default()

	// Services setup
	testService := services.NewTestService(db)
	userService := services.NewUserService(db)

	// Routes setup
	routes.SetupTestRoutes(router, testService)
	routes.SetupUserRoutes(router, userService)

	router.GET("/", func(c *gin.Context) {
		c.String(200, "Hello from Gin!")
	})

	// Server run
	if err := router.Run(host); err != nil {
		log.Fatalf("Error starting server on %s: %v\n", host, err)
	}

	log.Printf("Server is running on %s\n", host)

}
