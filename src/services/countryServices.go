package services

import (
	"github.com/ARQAP/ARQAP-Backend/src/models"
	"gorm.io/gorm"
)

type CountryService struct {
	db *gorm.DB
}

// NewCountryService creates a new instance of CountryService
func NewCountryService(db *gorm.DB) *CountryService {
	return &CountryService{db: db}
}

// GetAllcountries retrieves all country records from the database
func (s *CountryService) GetAllCountries() ([]models.CountryModel, error) {
	var countries []models.CountryModel
	result := s.db.Find(&countries)
	if result.Error != nil {
		return nil, result.Error
	}
	return countries, nil
}

// GetCountryByID retrieves a Country record by ID
func (s *CountryService) GetCountryByID (id int) (*models.CountryModel, error) {
	var country models.CountryModel
	result := s.db.First(&country, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &country, nil
}

// CreateCountry creates a new Country record in the database
func (s *CountryService) CreateCountry(country *models.CountryModel) (*models.CountryModel, error) {
	result := s.db.Create(country)
	if result.Error != nil {
		return nil, result.Error
	}
	return country, nil
}

// DeleteCountry deletes a Country record by ID
func (s *CountryService) DeleteCountry(id int) error {
	result := s.db.Delete(&models.CountryModel{}, id)
	return result.Error
}

// UpdateCountry updates an existing Country record
func (s *CountryService) UpdateCountry(id int,updatedCountry *models.CountryModel) (*models.CountryModel, error) {
	var country models.CountryModel
	result := s.db.First(&country, id)
	if result.Error != nil {
		return nil, result.Error
	}
	country = *updatedCountry
	s.db.Save(&country)
	return &country, nil
}