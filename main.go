package main

import (
	"log"
	"os"

	"github.com/ARQAP/ARQAP-Backend/src/db"
	"github.com/ARQAP/ARQAP-Backend/src/models"
	"github.com/ARQAP/ARQAP-Backend/src/routes"
	"github.com/ARQAP/ARQAP-Backend/src/seed"
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
	if err := db.AutoMigrate(&models.UserModel{}, &models.ArchaeologistModel{}, &models.CountryModel{}, &models.RegionModel{}); err != nil {
		log.Fatalf("Error during auto-migration: %v\n", err)
	}

	// Db seeding
	seed.Seed(db)

	// Port and host setup
	host := os.Getenv("SERVER_HOST")
	if host == "" {
		host = ":8080"
	}

	// Logs
	const Reset, Cyan = "\033[0m", "\033[36m"

	log.Printf("%s-----------------------------------------------: %s\n", Cyan, Reset)
	log.Printf("%sPGADMIN4 DASHBOARD: %s\n", Cyan, "http://localhost:5050")
	log.Printf("%s-----------------------------------------------: %s\n", Cyan, Reset)

	// Gin router setup
	router := gin.Default()

	// Services setup
	countryService := services.NewCountryService(db)
	regionService := services.NewRegionService(db)
	archaeologistService := services.NewArchaeologistService(db)
	userService := services.NewUserService(db)

	// Routes setup
	routes.SetupCountriesRoutes(router, countryService)
	routes.SetupRegionRoutes(router, regionService)
	routes.SetupArchaeologistRoutes(router, archaeologistService)
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
