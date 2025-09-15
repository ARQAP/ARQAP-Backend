package services

import (
	"github.com/ARQAP/ARQAP-Backend/src/models"
	"gorm.io/gorm"
)

type ArchaeologicalSiteService struct {
	db *gorm.DB
}

// NewArchaeologicalSiteService creates a new instance of ArchaeologicalSiteService
func NewArchaeologicalSiteService(db *gorm.DB) *ArchaeologicalSiteService {
	return &ArchaeologicalSiteService{db: db}
}

// GetAllArchaeologicalSites retrieves all ArchaeologicalSite records from the database
func (s *ArchaeologicalSiteService) GetAllArchaeologicalSites() ([]models.ArchaeologicalSiteModel, error) {
	var archaeologicalSites []models.ArchaeologicalSiteModel
	result := s.db.Preload("Region.Country").Find(&archaeologicalSites)
	if result.Error != nil {
		return nil, result.Error
	}
	return archaeologicalSites, nil
}

// CreateArchaeologicalSite creates a new ArchaeologicalSite record in the database
func (s *ArchaeologicalSiteService) CreateArchaeologicalSite(archaeologicalSite *models.ArchaeologicalSiteModel) (*models.ArchaeologicalSiteModel, error) {
	result := s.db.Create(archaeologicalSite)
	if result.Error != nil {
		return nil, result.Error
	}
	return archaeologicalSite, nil
}

// UpdateArchaeologicalSite updates an existing ArchaeologicalSite record in the database
func (s *ArchaeologicalSiteService) UpdateArchaeologicalSite(id int, updatedData *models.ArchaeologicalSiteModel) (*models.ArchaeologicalSiteModel, error) {
	var archaeologicalSite models.ArchaeologicalSiteModel
	if err := s.db.First(&archaeologicalSite, "id = ?", id).Error; err != nil {
		return nil, err
	}
	if err := s.db.Model(&archaeologicalSite).Updates(updatedData).Error; err != nil {
		return nil, err
	}
	return &archaeologicalSite, nil
}

// DeleteArchaeologicalSite deletes an ArchaeologicalSite record from the database
func (s *ArchaeologicalSiteService) DeleteArchaeologicalSite(id int) error {
	result := s.db.Delete(&models.ArchaeologicalSiteModel{}, "id = ?", id)
	return result.Error
}
