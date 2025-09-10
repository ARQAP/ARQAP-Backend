package services

import (
	"github.com/ARQAP/ARQAP-Backend/src/models"
	"gorm.io/gorm"
)

type ArchaeologistService struct {
	db *gorm.DB
}

// NewArchaeologistService creates a new instance of ArchaeologistService
func NewArchaeologistService(db *gorm.DB) *ArchaeologistService {
	return &ArchaeologistService{db: db}
}

// GetAllArchaeologists retrieves all Archaeologist records from the database
func (s *ArchaeologistService) GetAllArchaeologists() ([]models.ArchaeologistModel, error) {
	var archaeologists []models.ArchaeologistModel
	result := s.db.Find(&archaeologists)
	if result.Error != nil {
		return nil, result.Error
	}
	return archaeologists, nil
}

// CreateArchaeologist creates a new Archaeologist record in the database
func (s *ArchaeologistService) CreateArchaeologist(archaeologist *models.ArchaeologistModel) (*models.ArchaeologistModel, error) {
	result := s.db.Create(archaeologist)
	if result.Error != nil {
		return nil, result.Error
	}
	return archaeologist, nil
}

// UpdateArchaeologist updates an existing Archaeologist record in the database
func (s *ArchaeologistService) UpdateArchaeologist(id int, updatedData *models.ArchaeologistModel) (*models.ArchaeologistModel, error) {
	var archaeologist models.ArchaeologistModel
	if err := s.db.First(&archaeologist, "id = ?", id).Error; err != nil {
		return nil, err
	}
	if err := s.db.Model(&archaeologist).Updates(updatedData).Error; err != nil {
		return nil, err
	}
	return &archaeologist, nil
}

// DeleteArchaeologist deletes an Archaeologist record from the database
func (s *ArchaeologistService) DeleteArchaeologist(id int) error {
	result := s.db.Delete(&models.ArchaeologistModel{}, "id = ?", id)
	return result.Error
}