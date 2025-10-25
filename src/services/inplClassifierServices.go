package services

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/ARQAP/ARQAP-Backend/src/models"
	"gorm.io/gorm"
)

type FichaUpload struct {
	Reader       io.Reader
	OriginalName string
	ContentType  string
	Size         int64
}

type INPLService struct {
	db         *gorm.DB
	uploadRoot string
}

func NewINPLService(db *gorm.DB, uploadRoot string) *INPLService {
	return &INPLService{db: db, uploadRoot: uploadRoot}
}

// CreateClassifierWithFichas creates the classifier (no fields) + N photos in a single transaction.
func (s *INPLService) CreateClassifierWithFichas(files []FichaUpload) (*models.INPLClassifierModel, []models.INPLFicha, error) {
	if len(files) == 0 {
		return nil, nil, errors.New("at least one photo is required")
	}
	for i, f := range files {
		if f.Reader == nil || f.Size <= 0 {
			return nil, nil, fmt.Errorf("photo %d is invalid", i)
		}
		if !isAllowedImageType(f.ContentType) {
			return nil, nil, fmt.Errorf("photo %d: content-type not allowed: %s", i, f.ContentType)
		}
	}

	tx := s.db.Begin()
	if err := tx.Error; err != nil {
		return nil, nil, err
	}

	// 1) create empty classifier
	cls := &models.INPLClassifierModel{}
	if err := tx.Create(cls).Error; err != nil {
		tx.Rollback()
		return nil, nil, err
	}

	// 2) create folder per classifier
	dir := filepath.Join(s.uploadRoot, strconv.Itoa(cls.ID))
	if err := os.MkdirAll(dir, 0o755); err != nil {
		tx.Rollback()
		return nil, nil, err
	}

	var saved []string
	var fichas []models.INPLFicha
	for i, fu := range files {
		name := buildSafeFilename(fu.OriginalName, i)
		path := filepath.Join(dir, name)
		if err := saveToFile(path, fu.Reader); err != nil {
			cleanupFiles(saved)
			tx.Rollback()
			return nil, nil, err
		}
		saved = append(saved, path)

		rec := models.INPLFicha{
			INPLClassifierID: cls.ID,
			Filename:         name,
			OriginalName:     fu.OriginalName,
			FilePath:         path,
			ContentType:      fu.ContentType,
			Size:             fu.Size,
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		}
		if err := tx.Create(&rec).Error; err != nil {
			cleanupFiles(saved)
			tx.Rollback()
			return nil, nil, err
		}
		fichas = append(fichas, rec)
	}

	if err := tx.Commit().Error; err != nil {
		cleanupFiles(saved)
		return nil, nil, err
	}
	return cls, fichas, nil
}

// isAllowedImageType checks whether the content-type is one of the allowed image types.
func isAllowedImageType(ct string) bool {
	ct = strings.ToLower(ct)
	return ct == "image/jpeg" || ct == "image/jpg" || ct == "image/png" || ct == "image/webp"
}

// sanitizeFilename removes dangerous characters from the filename.
func sanitizeFilename(name string) string {
	name = strings.TrimSpace(name)
	name = strings.ReplaceAll(name, "..", "")
	name = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') ||
			(r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') ||
			strings.ContainsRune("._- ", r) {
			return r
		}
		return '_'
	}, name)
	return name
}

// buildSafeFilename generates a safe filename, appending an index if necessary.
func buildSafeFilename(original string, idx int) string {
	base := sanitizeFilename(original)
	if base == "" {
		base = fmt.Sprintf("photo_%d", time.Now().UnixNano())
	}

	if idx >= 0 {
		ext := filepath.Ext(base)
		name := strings.TrimSuffix(base, ext)
		return fmt.Sprintf("%s_%d%s", name, idx, ext)
	}
	return base
}

// saveToFile writes the contents of r to the given path.
func saveToFile(path string, r io.Reader) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, r)
	return err
}

// cleanupFiles deletes the files at the provided paths (best-effort).
func cleanupFiles(paths []string) {
	for _, p := range paths {
		_ = os.Remove(p)
	}
}

// AddFichasToClassifier adds N photos to an existing classifier.
func (s *INPLService) AddFichasToClassifier(classifierID int, files []FichaUpload) ([]models.INPLFicha, error) {
	if len(files) == 0 {
		return nil, errors.New("no files provided")
	}

	for i, f := range files {
		if f.Reader == nil || f.Size <= 0 {
			return nil, fmt.Errorf("photo %d is invalid", i)
		}
		if !isAllowedImageType(f.ContentType) {
			return nil, fmt.Errorf("photo %d: content-type not allowed: %s", i, f.ContentType)
		}
	}

	tx := s.db.Begin()
	if err := tx.Error; err != nil {
		return nil, err
	}

	var exists int64
	if err := tx.Model(&models.INPLClassifierModel{}).Where("id = ?", classifierID).Count(&exists).Error; err != nil || exists == 0 {
		tx.Rollback()
		return nil, errors.New("classifier does not exist")
	}

	dir := filepath.Join(s.uploadRoot, strconv.Itoa(classifierID))
	if err := os.MkdirAll(dir, 0o755); err != nil {
		tx.Rollback()
		return nil, err
	}

	var saved []string
	var out []models.INPLFicha

	for idx, fu := range files {
		filename := buildSafeFilename(fu.OriginalName, idx)
		final := filepath.Join(dir, filename)
		if err := saveToFile(final, fu.Reader); err != nil {
			cleanupFiles(saved)
			tx.Rollback()
			return nil, err
		}
		saved = append(saved, final)

		rec := models.INPLFicha{
			INPLClassifierID: classifierID,
			Filename:         filename,
			OriginalName:     fu.OriginalName,
			FilePath:         final,
			ContentType:      fu.ContentType,
			Size:             fu.Size,
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		}
		if err := tx.Create(&rec).Error; err != nil {
			cleanupFiles(saved)
			tx.Rollback()
			return nil, err
		}
		out = append(out, rec)
	}

	if err := tx.Commit().Error; err != nil {
		cleanupFiles(saved)
		return nil, err
	}
	return out, nil
}

// DeleteFicha removes a single photo both from the DB and from storage.
func (s *INPLService) DeleteFicha(fichaID int) error {
	var f models.INPLFicha
	if err := s.db.First(&f, fichaID).Error; err != nil {
		return err
	}
	if err := s.db.Delete(&models.INPLFicha{}, fichaID).Error; err != nil {
		return err
	}
	_ = os.Remove(f.FilePath)
	return nil
}

// GetClassifierByID fetches a classifier by ID, including its photos.
func (s *INPLService) GetClassifierByID(id int) (*models.INPLClassifierModel, error) {
	var classifier models.INPLClassifierModel
	result := s.db.Preload("INPLFichas").First(&classifier, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &classifier, nil
}

// GetAllClassifiers fetches all classifiers, optionally including their photos.
func (s *INPLService) GetAllClassifiers(preload bool) ([]models.INPLClassifierModel, error) {
	var list []models.INPLClassifierModel
	q := s.db
	if preload {
		q = q.Preload("INPLFichas")
	}
	if err := q.Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

// UpdateClassifier updates an existing INPLClassifier record.
func (s *INPLService) UpdateClassifier(id int, updatedClassifier *models.INPLClassifierModel) (*models.INPLClassifierModel, error) {
	var classifier models.INPLClassifierModel
	result := s.db.First(&classifier, id)
	if result.Error != nil {
		return nil, result.Error
	}
	classifier = *updatedClassifier
	s.db.Save(&classifier)
	return &classifier, nil
}

// ListFichasByClassifier lists the photos associated with a given classifier.
func (s *INPLService) ListFichasByClassifier(classifierID int) ([]models.INPLFicha, error) {
	var fichas []models.INPLFicha
	if err := s.db.Where("inpl_classifier_id = ?", classifierID).Find(&fichas).Error; err != nil {
		return nil, err
	}
	return fichas, nil
}

// DeleteClassifier removes a classifier and its associated photos, from both DB and storage.
func (s *INPLService) DeleteClassifier(id int) error {
	var fichas []models.INPLFicha
	if err := s.db.Where("inpl_classifier_id = ?", id).Find(&fichas).Error; err != nil {
		return err
	}

	if err := s.db.Delete(&models.INPLClassifierModel{}, id).Error; err != nil {
		return err
	}

	for _, f := range fichas {
		_ = os.Remove(f.FilePath)
	}

	dir := filepath.Join(s.uploadRoot, strconv.Itoa(id))
	_ = os.RemoveAll(dir)

	return nil
}

// ReplaceFicha replaces the file of an existing photo record.
func (s *INPLService) ReplaceFicha(fichaID int, file FichaUpload) (*models.INPLFicha, error) {
	if file.Reader == nil || file.Size <= 0 {
		return nil, errors.New("invalid file")
	}
	if !isAllowedImageType(file.ContentType) {
		return nil, fmt.Errorf("content-type not allowed: %s", file.ContentType)
	}

	var f models.INPLFicha
	if err := s.db.First(&f, fichaID).Error; err != nil {
		return nil, err
	}

	dir := filepath.Dir(f.FilePath)
	name := buildSafeFilename(file.OriginalName, 0)
	newPath := filepath.Join(dir, name)

	if err := saveToFile(newPath, file.Reader); err != nil {
		return nil, err
	}

	upd := map[string]any{
		"filename":      name,
		"original_name": file.OriginalName,
		"file_path":     newPath,
		"content_type":  file.ContentType,
		"size":          file.Size,
		"updated_at":    time.Now(),
	}
	if err := s.db.Model(&models.INPLFicha{}).Where("id = ?", fichaID).Updates(upd).Error; err != nil {
		_ = os.Remove(newPath)
		return nil, err
	}

	_ = os.Remove(f.FilePath)

	if err := s.db.First(&f, fichaID).Error; err != nil {
		return nil, err
	}
	return &f, nil
}
