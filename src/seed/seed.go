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
        return
    }

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

	// Other seeding operations can be added here
}