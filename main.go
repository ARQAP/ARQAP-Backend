package main

import (
	"log"
	"os"
	"path/filepath"

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
	if err := db.AutoMigrate(
		&models.UserModel{},
		&models.ArchaeologistModel{},
		&models.CountryModel{},
		&models.RegionModel{},
		&models.ArchaeologicalSiteModel{},
		&models.PhysicalLocationModel{},
		&models.CollectionModel{},
		&models.ShelfModel{},
		&models.InternalClassifierModel{},
		&models.INPLClassifierModel{},
		&models.INPLFicha{},
		&models.ArtefactModel{},
		&models.PictureModel{},
		&models.HistoricalRecordModel{}
	); err != nil {
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
	archaeologicalsiteService := services.NewArchaeologicalSiteService(db)
	countryService := services.NewCountryService(db)
	regionService := services.NewRegionService(db)
	archaeologistService := services.NewArchaeologistService(db)
	userService := services.NewUserService(db)
	collectionService := services.NewCollectionService(db)
	shelfService := services.NewShelfService(db)
	physicalLocationService := services.NewPhysicalLocationService(db)
	artefactService := services.NewArtefactService(db)
	internalLocationService := services.NewInternalClassifierService(db)

	// INPL uploads root (from env or default)
	inplUploadRoot := os.Getenv("INPL_UPLOAD_ROOT")
	if inplUploadRoot == "" {
		inplUploadRoot = filepath.Join("uploads", "inpl")
	}
	_ = os.MkdirAll(inplUploadRoot, 0o755)

	inplClassifierService := services.NewINPLService(db, inplUploadRoot)

	// Routes setup
	routes.SetupArchaeologicalSiteRoutes(router, archaeologicalsiteService)
	routes.SetupCountriesRoutes(router, countryService)
	routes.SetupRegionRoutes(router, regionService)
	routes.SetupArchaeologistRoutes(router, archaeologistService)
	routes.SetupUserRoutes(router, userService)
	routes.SetupPhysicalLocationRoutes(router, physicalLocationService)
	routes.SetupArtefactRoutes(router, artefactService)
	routes.SetupCollectionRoutes(router, collectionService)
	routes.SetupShelfsRoutes(router, shelfService)
	routes.SetupInternalClassifiersRoutes(router, internalLocationService)
	routes.SetupINPLClassifiersRoutes(router, inplClassifierService)

	// Test route
	router.GET("/", func(c *gin.Context) {
		c.String(200, "Hello from Gin! Server is up and running.")
	})

	// Server run
	if err := router.Run(host); err != nil {
		log.Fatalf("Error starting server on %s: %v\n", host, err)
	}

	log.Printf("Server is running on %s\n", host)
}
