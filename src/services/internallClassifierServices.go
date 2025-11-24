package services

import (
	"errors"
	"fmt"

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
func (s *InternalClassifierService) GetInternalClassifierByID(id int) (*models.InternalClassifierModel, error) {
	var internalClassifier models.InternalClassifierModel
	result := s.db.First(&internalClassifier, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &internalClassifier, nil
}

// GetInternalClassifiersByName retrieves internal classifiers that match the provided name
func (s *InternalClassifierService) GetInternalClassifiersByName(name string) ([]models.InternalClassifierModel, error) {
	var internalClassifiers []models.InternalClassifierModel
	result := s.db.Where("name = ?", name).Find(&internalClassifiers)
	if result.Error != nil {
		return nil, result.Error
	}
	return internalClassifiers, nil
}

// GetAllInternalClassifierNames returns the distinct names of internal classifiers
func (s *InternalClassifierService) GetAllInternalClassifierNames() ([]string, error) {
	var names []string
	result := s.db.Model(&models.InternalClassifierModel{}).Distinct("name").Pluck("name", &names)
	if result.Error != nil {
		return nil, result.Error
	}
	return names, nil
}

// CreateInternalClassifier creates a new InternalClassifier record in the database
func (s *InternalClassifierService) CreateInternalClassifier(internalClassifier *models.InternalClassifierModel) (*models.InternalClassifierModel, error) {
	// Prevent duplicate where both name and number match an existing record
	var existing models.InternalClassifierModel
	var res *gorm.DB
	if internalClassifier.Number == nil {
		res = s.db.Where("name = ? AND number IS NULL", internalClassifier.Name).First(&existing)
	} else {
		res = s.db.Where("name = ? AND number = ?", internalClassifier.Name, *internalClassifier.Number).First(&existing)
	}
	if res.Error == nil {
		numStr := "null"
		if internalClassifier.Number != nil {
			numStr = fmt.Sprintf("%d", *internalClassifier.Number)
		}
		return nil, fmt.Errorf("internal classifier with name '%s' and number %s already exists", internalClassifier.Name, numStr)
	}
	if res.Error != nil && !errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, res.Error
	}

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
func (s *InternalClassifierService) UpdateInternalClassifier(id int, updatedInternalClassifier *models.InternalClassifierModel) (*models.InternalClassifierModel, error) {
	var internalClassifier models.InternalClassifierModel
	result := s.db.First(&internalClassifier, id)
	if result.Error != nil {
		return nil, result.Error
	}
	// Prevent duplicate where both name and number match another existing record
	var existing models.InternalClassifierModel
	var res *gorm.DB
	if updatedInternalClassifier.Number == nil {
		res = s.db.Where("name = ? AND number IS NULL AND id <> ?", updatedInternalClassifier.Name, id).First(&existing)
	} else {
		res = s.db.Where("name = ? AND number = ? AND id <> ?", updatedInternalClassifier.Name, *updatedInternalClassifier.Number, id).First(&existing)
	}
	if res.Error == nil {
		numStr := "null"
		if updatedInternalClassifier.Number != nil {
			numStr = fmt.Sprintf("%d", *updatedInternalClassifier.Number)
		}
		return nil, fmt.Errorf("internal classifier with name '%s' and number %s already exists", updatedInternalClassifier.Name, numStr)
	}
	if res.Error != nil && !errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return nil, res.Error
	}

	internalClassifier = *updatedInternalClassifier
	if err := s.db.Save(&internalClassifier).Error; err != nil {
		return nil, err
	}
	return &internalClassifier, nil
}
