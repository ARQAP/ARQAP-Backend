package services

import (
	"github.com/ARQAP/ARQAP-Backend/src/models"
	"gorm.io/gorm"
)

type TestService struct {
	db *gorm.DB
}

// NewTestService creates a new instance of TestService
func NewTestService(db *gorm.DB) *TestService {
	return &TestService{db: db}
}

// GetAllTests retrieves all Test records from the database
func (s *TestService) GetAllTests() ([]models.TestModel, error) {
	var tests []models.TestModel
	result := s.db.Find(&tests)
	if result.Error != nil {
		return nil, result.Error
	}
	return tests, nil
}

// CreateTest creates a new Test record in the database
func (s *TestService) CreateTest(test *models.TestModel) (*models.TestModel, error) {
	result := s.db.Create(test)
	if result.Error != nil {
		return nil, result.Error
	}
	return test, nil
}