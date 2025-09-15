package services

import (
	"github.com/ARQAP/ARQAP-Backend/src/models"
	"gorm.io/gorm"
)

type PhysicalLocationService struct {
	db *gorm.DB
}

// NewPhysicalLocationService creates a new instance of PhysicalLocationService
func NewPhysicalLocationService(db *gorm.DB) *PhysicalLocationService {
	return &PhysicalLocationService{db: db}
}

// GetAllPhysicalLocations retrieves all physical locations from the database
func (s *PhysicalLocationService) GetAllPhysicalLocations() ([]models.PhysicalLocationModel, error) {
	var locations []models.PhysicalLocationModel
	if err := s.db.Preload("Shelf").Find(&locations).Error; err != nil {
		return nil, err
	}
	return locations, nil
}

// GetPhysicalLocationByID retrieves a physical location by its ID
func (s *PhysicalLocationService) GetPhysicalLocationByID(id int) (*models.PhysicalLocationModel, error) {
	var location models.PhysicalLocationModel
	if err := s.db.First(&location, id).Error; err != nil {
		return nil, err
	}
	return &location, nil
}

// CreatePhysicalLocation creates a new physical location in the database
func (s *PhysicalLocationService) CreatePhysicalLocation(location *models.PhysicalLocationModel) error {
	if err := s.db.Create(location).Error; err != nil {
		return err
	}
	return nil
}

// UpdatePhysicalLocation updates an existing physical location in the database
func (s *PhysicalLocationService) UpdatePhysicalLocation(location *models.PhysicalLocationModel) error {
	if err := s.db.Save(location).Error; err != nil {
		return err
	}
	return nil
}

// DeletePhysicalLocation removes a physical location from the database
func (s *PhysicalLocationService) DeletePhysicalLocation(id int) error {
	if err := s.db.Delete(&models.PhysicalLocationModel{}, id).Error; err != nil {
		return err
	}
	return nil
}