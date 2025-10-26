package services

import (
	"github.com/ARQAP/ARQAP-Backend/src/models"
	"gorm.io/gorm"
)

type RequesterService struct {
	db *gorm.DB
}

// NewRequesterService creates a new instance of RequesterService
func NewRequesterService(db *gorm.DB) *RequesterService {
	return &RequesterService{db: db}
}

// GetAllRequesters retrieves all Requester records from the database
func (s *RequesterService) GetAllRequesters() ([]models.RequesterModel, error) {
	var requesters []models.RequesterModel
	result := s.db.Find(&requesters)
	if result.Error != nil {
		return nil, result.Error
	}
	return requesters, nil
}

// GetRequesterByID retrieves a Requester record by its ID
func (s *RequesterService) GetRequesterByID(id int) (*models.RequesterModel, error) {
	var requester models.RequesterModel
	result := s.db.First(&requester, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &requester, nil
}

// CreateRequester creates a new Requester record in the database
func (s *RequesterService) CreateRequester(requester *models.RequesterModel) (*models.RequesterModel, error){
	result := s.db.Create(requester)
	if result.Error != nil {
		return nil, result.Error
	}
	return requester, nil
}

// DeleteRequester deletes a Requester record by its ID
func (s *RequesterService) DeleteRequester(id int) error {
	result := s.db.Delete(&models.RequesterModel{}, id)
	return result.Error
}

// UpdateRequester updates an existing Requester record
func (s *RequesterService) UpdateRequester(id int, updatedRequester *models.RequesterModel) (*models.RequesterModel, error) {
	var requester models.RequesterModel
	result := s.db.First(&requester, id)
	if result.Error != nil {
		return nil, result.Error
	}
    // Set the ID to ensure we update the correct record
    updatedRequester.Id = id
    
    // Use Updates instead of replacing the whole object
    result = s.db.Model(&requester).Updates(updatedRequester)
    if result.Error != nil {
        return nil, result.Error
    }
    
    // Fetch the updated record
    result = s.db.First(&requester, id)
    if result.Error != nil {
        return nil, result.Error
    }
    
    return &requester, nil
}