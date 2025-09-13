package services

import (
	"github.com/ARQAP/ARQAP-Backend/src/models"
	"gorm.io/gorm"
)

type RegionService struct {
	db *gorm.DB
}

// NewRegionService creates a new instance of RegionService
func NewRegionService(db *gorm.DB) *RegionService {
	return &RegionService{db: db}
}

// GetAllRegions retrieves all Region records from the database
func (s *RegionService) GetAllRegions() ([]models.RegionModel, error) {
	var regions []models.RegionModel
	result := s.db.Preload("Country").Find(&regions)
	if result.Error != nil {
		return nil, result.Error
	}
	return regions, nil
}

// CreateRegion creates a new Region record in the database
func (s *RegionService) CreateRegion(region *models.RegionModel) (*models.RegionModel, error) {
	result := s.db.Create(region)
	if result.Error != nil {
		return nil, result.Error
	}
	return region, nil
}

// UpdateRegion updates an existing Region record in the database
func (s *RegionService) UpdateRegion(id int, updatedData *models.RegionModel) (*models.RegionModel, error) {
	var region models.RegionModel
	if err := s.db.First(&region, "id = ?", id).Error; err != nil {
		return nil, err
	}
	if err := s.db.Model(&region).Updates(updatedData).Error; err != nil {
		return nil, err
	}
	return &region, nil
}

// DeleteRegion deletes an Region record from the database
func (s *RegionService) DeleteRegion(id int) error {
	result := s.db.Delete(&models.RegionModel{}, "id = ?", id)
	return result.Error
}
