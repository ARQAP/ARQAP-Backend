package services

import (
	"github.com/ARQAP/ARQAP-Backend/src/models"
	"gorm.io/gorm"
)

type MentionService struct {
	db *gorm.DB
}

// NewMentionService creates a new instance of MentionService
func NewMentionService(db *gorm.DB) *MentionService {
	return &MentionService{db: db}
}

// GetAllMentions retrieves all Mention records from the database
func (s *MentionService) GetAllMentions() ([]models.MentionModel, error) {
	var mentions []models.MentionModel
	result := s.db.Find(&mentions)
	if result.Error != nil {
		return nil, result.Error
	}
	return mentions, nil
}

// GetMentionsByArtefactID retrieves Mention records associated with a specific Artefact ID
func (s *MentionService) GetMentionsByArtefactID(artefactID int) ([]models.MentionModel, error) {
	var mentions []models.MentionModel
	result := s.db.Where("artefact_id = ?", artefactID).Find(&mentions)
	if result.Error != nil {
		return nil, result.Error
	}
	return mentions, nil
}

// GetMentionByID retrieves a Mention record by its ID
func (s *MentionService) GetMentionByID(id int) (*models.MentionModel, error) {
	var mention models.MentionModel
	result := s.db.First(&mention, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &mention, nil
}

// CreateMention creates a new Mention record in the database
func (s *MentionService) CreateMention(mention *models.MentionModel) (*models.MentionModel, error) {
	result := s.db.Create(mention)
	if result.Error != nil {
		return nil, result.Error
	}
	return mention, nil
}

// DeleteMention deletes a Mention record by its ID
func (s *MentionService) DeleteMention(id int) error {
	result := s.db.Delete(&models.MentionModel{}, id)
	return result.Error
}

// UpdateMention updates an existing Mention record
func (s *MentionService) UpdateMention(id int, updatedMention *models.MentionModel) (*models.MentionModel, error) {
	var mention models.MentionModel
	result := s.db.First(&mention, id)
	if result.Error != nil {
		return nil, result.Error
	}

	// Set the ID to ensure we update the correct record
	updatedMention.Id = id

	// Use Updates instead of replacing the whole object
	result = s.db.Model(&mention).Updates(updatedMention)
	if result.Error != nil {
		return nil, result.Error
	}

	// Fetch the updated record
	result = s.db.First(&mention, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &mention, nil
}
