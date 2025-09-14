package services

import (
	"github.com/ARQAP/ARQAP-Backend/src/models"
	"gorm.io/gorm"
)

type PhysicalLocationService struct {
	db *gorm.DB
}

// NewPhysicalLocationService creates a new instance of PhysicalLocationService
func NewPhysicalLocationService(db *gorm.DB) *ShelfService {
	return &ShelfService{db: db}
}