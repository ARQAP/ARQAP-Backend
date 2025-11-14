package services

import (
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/ARQAP/ARQAP-Backend/src/models"
	"gorm.io/gorm"
)

// Cache entry
type CacheEntry struct {
	Data      interface{}
	ExpiresAt time.Time
}

type ArtefactService struct {
	db    *gorm.DB
	cache map[string]*CacheEntry
	mutex sync.RWMutex
}

func NewArtefactService(db *gorm.DB) *ArtefactService {
	service := &ArtefactService{
		db:    db,
		cache: make(map[string]*CacheEntry),
	}

	// Clean up cache every 30 minutes
	go service.cleanupCache()

	return service
}

func (s *ArtefactService) cleanupCache() {
	ticker := time.NewTicker(30 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		s.mutex.Lock()
		now := time.Now()
		for key, entry := range s.cache {
			if now.After(entry.ExpiresAt) {
				delete(s.cache, key)
			}
		}
		s.mutex.Unlock()
	}
}

func (s *ArtefactService) setCache(key string, data interface{}, duration time.Duration) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.cache[key] = &CacheEntry{
		Data:      data,
		ExpiresAt: time.Now().Add(duration),
	}
}

func (s *ArtefactService) getCache(key string) (interface{}, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	entry, exists := s.cache[key]
	if !exists || time.Now().After(entry.ExpiresAt) {
		return nil, false
	}

	return entry.Data, true
}

func (s *ArtefactService) invalidateCache(pattern string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for key := range s.cache {
		if contains(key, pattern) {
			delete(s.cache, key)
		}
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		(len(s) > len(substr) && s[:len(substr)] == substr))
}

func (s *ArtefactService) GetAllArtefacts(shelfId *int) ([]models.ArtefactModel, error) {
	// Si hay filtro por shelf, no usar cache general
	var cacheKey string
	if shelfId != nil {
		cacheKey = fmt.Sprintf("artefacts_shelf_%d", *shelfId)
	} else {
		cacheKey = "all_artefacts"
	}

	// Try to get from cache
	if cached, found := s.getCache(cacheKey); found {
		return cached.([]models.ArtefactModel), nil
	}

	// If not in cache, query DB
	var artefacts []models.ArtefactModel
	query := s.db.Preload("Picture").Preload("HistoricalRecord").Preload("Archaeologist").Preload("ArchaeologicalSite").Preload("PhysicalLocation").Preload("Collection").Preload("InplClassifier").Preload("InternalClassifier").Preload("PhysicalLocation.Shelf")

	// Aplicar filtro por shelfId si está presente
	if shelfId != nil {
		query = query.Joins("JOIN physical_location_models ON physical_location_models.id = artefact_models.physical_location_id").
			Where("physical_location_models.shelf_id = ?", *shelfId)
	}

	err := query.Find(&artefacts).Error

	if err == nil {
		// Save to cache for 5 minutes
		s.setCache(cacheKey, artefacts, 5*time.Minute)
	}

	return artefacts, err
}

func (s *ArtefactService) GetArtefactByID(id int) (*models.ArtefactModel, error) {
	cacheKey := fmt.Sprintf("artefact_%d", id)

	// Try to get from cache
	if cached, found := s.getCache(cacheKey); found {
		artefact := cached.(models.ArtefactModel)
		return &artefact, nil
	}

	// If not in cache, query DB
	var artefact models.ArtefactModel

	err := s.db.Preload("Picture").Preload("HistoricalRecord").Preload("Archaeologist").Preload("ArchaeologicalSite").Preload("PhysicalLocation").Preload("Collection").Preload("InplClassifier").Preload("InternalClassifier").Preload("PhysicalLocation.Shelf").First(&artefact, id).Error
	if err != nil {
		return nil, err
	}

	// Save to cache for 10 minutes
	s.setCache(cacheKey, artefact, 10*time.Minute)

	return &artefact, nil
}

func (s *ArtefactService) CreateArtefact(artefact *models.ArtefactModel) error {
	if err := s.db.Create(artefact).Error; err != nil {
		return err
	}

	// Invalidate related cache
	s.invalidateCache("all_artefacts")

	return nil
}

func (s *ArtefactService) UpdateArtefact(id int, artefact *models.ArtefactModel) error {
	if err := s.db.Where("id = ?", id).Updates(artefact).Error; err != nil {
		return err
	}

	// Invalidate specific and general cache
	s.invalidateCache(fmt.Sprintf("artefact_%d", id))
	s.invalidateCache("all_artefacts")

	return nil
}

func (s *ArtefactService) DeleteArtefact(id int) error {
	// Obtain artefact with its associated Picture and HistoricalRecord paths
	var artefact models.ArtefactModel
	if err := s.db.Preload("Picture").Preload("HistoricalRecord").First(&artefact, id).Error; err != nil {
		return err
	}

	// Delete the files if they exist
	if len(artefact.Picture) > 0 && artefact.Picture[0].FilePath != "" {
		os.Remove(artefact.Picture[0].FilePath)
	}
	if len(artefact.HistoricalRecord) > 0 && artefact.HistoricalRecord[0].FilePath != "" {
		os.Remove(artefact.HistoricalRecord[0].FilePath)
	}

	// Delete from DB
	if err := s.db.Delete(&artefact, id).Error; err != nil {
		return err
	}

	// Invalidate cache
	s.invalidateCache(fmt.Sprintf("artefact_%d", id))
	s.invalidateCache("all_artefacts")

	return nil
}

func (s *ArtefactService) GetPictureByArtefactID(artefactID int) (*models.PictureModel, error) {
	cacheKey := fmt.Sprintf("picture_%d", artefactID)
	if cached, found := s.getCache(cacheKey); found {
		pic := cached.(models.PictureModel)
		return &pic, nil
	}

	var picture models.PictureModel
	err := s.db.Where("artefact_id = ?", artefactID).First(&picture).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// caso normal: aún no hay foto
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	s.setCache(cacheKey, picture, 30*time.Minute)
	return &picture, nil
}

func (s *ArtefactService) GetHistoricalRecordByArtefactID(artefactID int) (*models.HistoricalRecordModel, error) {
	cacheKey := fmt.Sprintf("historical_record_%d", artefactID)

	// Try to get from cache
	if cached, found := s.getCache(cacheKey); found {
		record := cached.(models.HistoricalRecordModel)
		return &record, nil
	}

	// If not in cache, query DB
	var record models.HistoricalRecordModel
	err := s.db.Where("artefact_id = ?", artefactID).First(&record).Error
	if err != nil {
		return nil, err
	}

	// Save to cache for 30 minutes
	s.setCache(cacheKey, record, 30*time.Minute)

	return &record, nil
}

func (s *ArtefactService) SavePicture(picture *models.PictureModel) error {
	var existing models.PictureModel
	err := s.db.Where("artefact_id = ?", picture.ArtefactID).First(&existing).Error

	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		// no había foto: crear
		if err := s.db.Create(picture).Error; err != nil {
			return err
		}
	case err != nil:
		// error real de DB
		return err
	default:
		// había foto: borrar archivo viejo y actualizar registro
		if existing.FilePath != "" {
			_ = os.Remove(existing.FilePath)
		}
		// aseguramos update sobre el registro existente
		picture.ID = existing.ID
		if err := s.db.Model(&existing).Updates(map[string]interface{}{
			"file_path":    picture.FilePath,
			"content_type": picture.ContentType,
			"size":         picture.Size,
			"updated_at":   time.Now(),
		}).Error; err != nil {
			return err
		}
	}

	s.invalidateCache(fmt.Sprintf("picture_%d", picture.ArtefactID))
	s.invalidateCache(fmt.Sprintf("artefact_%d", picture.ArtefactID))
	s.invalidateCache("all_artefacts")
	return nil
}

func (s *ArtefactService) SaveHistoricalRecord(record *models.HistoricalRecordModel) error {
	// Check if a document already exists for this artefact
	var existing models.HistoricalRecordModel
	if err := s.db.Where("artefact_id = ?", record.ArtefactID).First(&existing).Error; err == nil {
		// Already exists, delete previous file
		os.Remove(existing.FilePath)
		// Update existing record
		if err := s.db.Where("artefact_id = ?", record.ArtefactID).Updates(record).Error; err != nil {
			return err
		}
	} else {
		// Does not exist, create new
		if err := s.db.Create(record).Error; err != nil {
			return err
		}
	}

	// Invalidate related cache
	s.invalidateCache(fmt.Sprintf("historical_record_%d", record.ArtefactID))
	s.invalidateCache(fmt.Sprintf("artefact_%d", record.ArtefactID))
	s.invalidateCache("all_artefacts")

	return nil
}
