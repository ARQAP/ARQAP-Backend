package seed

import (
	"log"

	"github.com/ARQAP/ARQAP-Backend/src/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB) {
	// Users
	var user models.UserModel
	result := db.Where("username = ?", "arqap").First(&user)
	if result.Error == nil {
		log.Println("User 'arqap' already exists")
	} else {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("arqap"), bcrypt.DefaultCost)

		newUser := models.UserModel{
			Username: "arqap",
			Password: string(hashedPassword),
		}
		if err := db.Create(&newUser).Error; err != nil {
			log.Printf("Failed to create user: %v\n", err)
		} else {
			log.Println("User 'arqap' created")
		}
	}

	// Shelves seeding - Create shelves from code 1 to 30
	log.Println("Checking and creating shelves from code 1 to 30...")
	createdCount := 0
	for i := 1; i <= 30; i++ {
		var existingShelf models.ShelfModel
		checkResult := db.Where("code = ?", i).First(&existingShelf)
		if checkResult.Error == nil {
			log.Printf("Shelf with code %d already exists, skipping\n", i)
		} else {
			shelf := models.ShelfModel{
				Code: i,
			}
			if err := db.Create(&shelf).Error; err != nil {
				log.Printf("Failed to create shelf with code %d: %v\n", i, err)
			} else {
				log.Printf("Shelf with code %d created\n", i)
				createdCount++
			}
		}
	}
	if createdCount > 0 {
		log.Printf("Finished creating %d new shelves\n", createdCount)
	} else {
		log.Println("All shelves already exist")
	}

	// Physical Locations seeding - Create physical locations for all shelves
	log.Println("Checking and creating physical locations for all shelves...")
	var allShelves []models.ShelfModel
	if err := db.Find(&allShelves).Error; err != nil {
		log.Printf("Failed to fetch shelves: %v\n", err)
		return
	}

	locationsCreatedCount := 0
	levels := []models.LevelNumber{models.Level1, models.Level2, models.Level3, models.Level4}
	columns := []models.ColumnLetter{models.ColumnA, models.ColumnB, models.ColumnC, models.ColumnD}

	for _, shelf := range allShelves {
		// Mesas de trabajo (códigos 28, 29, 30) solo tienen nivel 1, columna A
		if shelf.Code == 28 || shelf.Code == 29 || shelf.Code == 30 {
			var existingLocation models.PhysicalLocationModel
			checkResult := db.Where(`shelf_id = ? AND level = ? AND "column" = ?`, shelf.ID, models.Level1, models.ColumnA).First(&existingLocation)
			if checkResult.Error == nil {
				log.Printf("Physical location for shelf %d (Level 1, Column A) already exists, skipping\n", shelf.Code)
			} else {
				location := models.PhysicalLocationModel{
					ShelfId: shelf.ID,
					Level:   models.Level1,
					Column:  models.ColumnA,
				}
				if err := db.Create(&location).Error; err != nil {
					log.Printf("Failed to create physical location for shelf %d (Level 1, Column A): %v\n", shelf.Code, err)
				} else {
					log.Printf("Physical location created for shelf %d (Level 1, Column A)\n", shelf.Code)
					locationsCreatedCount++
				}
			}
		} else {
			// Estantes normales: crear todas las combinaciones de niveles 1-4 y columnas A-D
			for _, level := range levels {
				for _, column := range columns {
					var existingLocation models.PhysicalLocationModel
					checkResult := db.Where(`shelf_id = ? AND level = ? AND "column" = ?`, shelf.ID, level, column).First(&existingLocation)
					if checkResult.Error == nil {
						log.Printf("Physical location for shelf %d (Level %d, Column %s) already exists, skipping\n", shelf.Code, level, column)
					} else {
						location := models.PhysicalLocationModel{
							ShelfId: shelf.ID,
							Level:   level,
							Column:  column,
						}
						if err := db.Create(&location).Error; err != nil {
							log.Printf("Failed to create physical location for shelf %d (Level %d, Column %s): %v\n", shelf.Code, level, column, err)
						} else {
							log.Printf("Physical location created for shelf %d (Level %d, Column %s)\n", shelf.Code, level, column)
							locationsCreatedCount++
						}
					}
				}
			}
		}
	}

	if locationsCreatedCount > 0 {
		log.Printf("Finished creating %d new physical locations\n", locationsCreatedCount)
	} else {
		log.Println("All physical locations already exist")
	}
}
