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

// GetAllShelves retrieves all Shelf records from the database
func (s *ShelfService) GetAllShelves() ([]models.ShelfModel, error) {
	var shelves []models.ShelfModel
	result := s.db.Find(&shelves)
	if result.Error != nil {
		return nil, result.Error
	}
	return shelves, nil
}

// GetShelfByID retrieves a Shelf record by its ID
func (s *ShelfService) GetShelfByID(id int) (*models.ShelfModel, error) {
	var shelf models.ShelfModel
	result := s.db.First(&shelf, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &shelf, nil
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

// DeleteShelf removes a Shelf record from the database
func (s *ShelfService) DeleteShelf(id int) error {
	result := s.db.Delete(&models.ShelfModel{}, id)
	return result.Error
}