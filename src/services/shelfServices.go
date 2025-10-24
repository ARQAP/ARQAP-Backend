package services

import (
	"github.com/ARQAP/ARQAP-Backend/src/models"
	"gorm.io/gorm"
)

type ShelfService struct {
	db *gorm.DB
}

// NewShelfService creates a new instance of ShelfService
func NewShelfService(db *gorm.DB) *ShelfService {
	return &ShelfService{db: db}
}

// GetAllShelfs retrieves all Shelf records from the database
func (s *ShelfService) GetAllShelfs() ([]models.ShelfModel, error) {
	var shelfs []models.ShelfModel
		result := s.db.Find(&shelfs)
	if result.Error != nil {
		return nil, result.Error
	}
	return shelfs, nil
}

// CreateShelf creates a new Shelf record in the database
func (s *ShelfService) CreateShelf(shelf *models.ShelfModel) (*models.ShelfModel, error) {
	result := s.db.Create(shelf)
	if result.Error != nil {
		return nil, result.Error
	}
	return shelf, nil
}

// UpdateShelf updates an existing Shelf record in the database
func (s *ShelfService) UpdateShelf(id int, updatedData *models.ShelfModel) (*models.ShelfModel, error) {
	var shelf models.ShelfModel
	if err := s.db.First(&shelf, "id = ?", id).Error; err != nil {
		return nil, err
	}
	if err := s.db.Model(&shelf).Updates(updatedData).Error; err != nil {
		return nil, err
	}
	return &shelf, nil
}

// DeleteShelf deletes an Shelf record from the database
func (s *ShelfService) DeleteShelf(id int) error {
	result := s.db.Delete(&models.ShelfModel{}, "id = ?", id)
	return result.Error
}

// GetShelfByID retrieves a Shelf record by ID
func (s *ShelfService) GetShelfByID (id int) (*models.ShelfModel, error) {
	var shelf models.ShelfModel
	result := s.db.First(&shelf, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &shelf, nil
}