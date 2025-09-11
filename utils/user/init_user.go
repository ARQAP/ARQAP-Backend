package main

import (
	"log"
	"os"

	"github.com/ARQAP/ARQAP-Backend/src/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	// Migrate schema if not exists
	if err := db.AutoMigrate(&models.UserModel{}); err != nil {
		log.Fatalf("failed to migrate user model: %v", err)
	}

	var user models.UserModel
	result := db.Where("username = ?", "arqap").First(&user)
	if result.Error == nil {
		log.Println("User 'arqap' already exists")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("arqap"), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("failed to hash password: %v", err)
	}

	newUser := models.UserModel{
		Username: "arqap",
		Password: string(hashedPassword),
	}
	if err := db.Create(&newUser).Error; err != nil {
		log.Fatalf("failed to create user: %v", err)
	}
	log.Println("User 'arqap' created")
}
