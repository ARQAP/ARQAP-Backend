package services

import (
	"github.com/ARQAP/ARQAP-Backend/src/models"
	"gorm.io/gorm"
)

type InternalClassifierService struct {
	db *gorm.DB
}

// NewInternalClassifierService creates a new instance of InternalClassifierService
func NewInternalClassifierService(db *gorm.DB) *InternalClassifierService {
	return &InternalClassifierService{db: db}
}

// GetAllinternalClassifiers retrieves all internalClassifier records from the database
func (s *InternalClassifierService) GetAllInternalClassifiers() ([]models.InternalClassifierModel, error) {
	var internalClassifiers []models.InternalClassifierModel
	result := s.db.Find(&internalClassifiers)
	if result.Error != nil {
		return nil, result.Error
	}
	return internalClassifiers, nil
}

// GetInternalClassifierByID retrieves a InternalClassifier record by ID
func (s *InternalClassifierService) GetInternalClassifierByID (id int) (*models.InternalClassifierModel, error) {
	var internalClassifier models.InternalClassifierModel
	result := s.db.First(&internalClassifier, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &internalClassifier, nil
}

// CreateInternalClassifier creates a new InternalClassifier record in the database
func (s *InternalClassifierService) CreateInternalClassifier(internalClassifier *models.InternalClassifierModel) (*models.InternalClassifierModel, error) {
	result := s.db.Create(internalClassifier)
	if result.Error != nil {
		return nil, result.Error
	}
	return internalClassifier, nil
}

// DeleteInternalClassifier deletes a InternalClassifier record by ID
func (s *InternalClassifierService) DeleteInternalClassifier(id int) error {
	result := s.db.Delete(&models.InternalClassifierModel{}, id)
	return result.Error
}

// UpdateInternalClassifier updates an existing InternalClassifier record
func (s *InternalClassifierService) UpdateInternalClassifier(id int,updatedInternalClassifier *models.InternalClassifierModel) (*models.InternalClassifierModel, error) {
	var internalClassifier models.InternalClassifierModel
	result := s.db.First(&internalClassifier, id)
	if result.Error != nil {
		return nil, result.Error
	}
	internalClassifier = *updatedInternalClassifier
	s.db.Save(&internalClassifier)
	return &internalClassifier, nil
}