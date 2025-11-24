package seed

import (
    "log"
    "golang.org/x/crypto/bcrypt"
    "gorm.io/gorm"
    "github.com/ARQAP/ARQAP-Backend/src/models"
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
}