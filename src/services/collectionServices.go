package services

import (
	"github.com/ARQAP/ARQAP-Backend/src/models"
	"gorm.io/gorm"
)

type CollectionService struct {
	db *gorm.DB
}

// NewCollectionService creates a new instance of CollectionService
func NewCollectionService(db *gorm.DB) *CollectionService {
	return &CollectionService{db: db}
}

// GetAllCollections retrieves all collection records from the database
func (s *CollectionService) GetAllCollections() ([]models.CollectionModel, error) {
	var collections []models.CollectionModel
	result := s.db.Find(&collections)
	if result.Error != nil {
		return nil, result.Error
	}
	return collections, nil
}

// CreateCollection creates a new collection record in the database
func (s *CollectionService) CreateCollection(collection *models.CollectionModel) (*models.CollectionModel, error) {
	result := s.db.Create(collection)
	if result.Error != nil {
		return nil, result.Error
	}
	return collection, nil
}

// UpdateCollection updates an existing collection record in the database
func (s *CollectionService) UpdateCollection(id int, updatedData *models.CollectionModel) (*models.CollectionModel, error) {
	var collection models.CollectionModel
	if err := s.db.First(&collection, "id = ?", id).Error; err != nil {
		return nil, err
	}
	if err := s.db.Model(&collection).Updates(updatedData).Error; err != nil {
		return nil, err
	}
	return &collection, nil
}

// DeleteCollection deletes a collection record from the database
func (s *CollectionService) DeleteCollection(id int) error {
	result := s.db.Delete(&models.CollectionModel{}, "id = ?", id)
	return result.Error
}
