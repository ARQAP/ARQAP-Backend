package services

import (
	"errors"

	"github.com/ARQAP/ARQAP-Backend/src/models"
	"gorm.io/gorm"
)

type LoanService struct {
	db *gorm.DB
}

// NewLoanService creates a new instance of LoanService
func NewLoanService(db *gorm.DB) *LoanService {
	return &LoanService{db: db}
}

// GetAllLoans retrieves all Loan records from the database
func (s *LoanService) GetAllLoans() ([]models.LoanModel, error) {
	var loans []models.LoanModel

	result := s.db.
		Preload("Requester").
		Preload("Artefact").
		Preload("Artefact.InternalClassifier").
		Find(&loans)

	return loans, result.Error
}

// GetLoanByID retrieves a Loan record by its ID
func (s *LoanService) GetLoanByID(id int) (*models.LoanModel, error) {
	var loan models.LoanModel

	result := s.db.
		Preload("Requester").
		Preload("Artefact").
		Preload("Artefact.InternalClassifier").
		First(&loan, id)

	if result.Error != nil {
		return nil, result.Error
	}
	return &loan, nil
}

// CreateLoan creates a new Loan record in the database
// y marca la pieza asociada como no disponible (available = false)
func (s *LoanService) CreateLoan(loan *models.LoanModel) (*models.LoanModel, error) {
	err := s.db.Transaction(func(tx *gorm.DB) error {
		// 1) Verificar que la pieza esté disponible antes de crear el préstamo
		if loan.ArtefactId != nil && *loan.ArtefactId != 0 {
			var artefact models.ArtefactModel
			if err := tx.First(&artefact, *loan.ArtefactId).Error; err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return errors.New("la pieza arqueológica no existe")
				}
				return err
			}

			if !artefact.Available {
				return errors.New("la pieza arqueológica no está disponible para préstamo (ya está prestada)")
			}
		}

		// 2) Crear el préstamo
		if err := tx.Create(loan).Error; err != nil {
			return err
		}

		// 3) Marcar la pieza como NO disponible
		if loan.ArtefactId != nil && *loan.ArtefactId != 0 {
			if err := tx.Model(&models.ArtefactModel{}).
				Where("id = ?", *loan.ArtefactId).
				Update("available", false).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Opcional: devolver el préstamo con las relaciones cargadas
	if err := s.db.
		Preload("Requester").
		Preload("Artefact").
		Preload("Artefact.InternalClassifier").
		First(loan, loan.Id).Error; err != nil {
		return nil, err
	}

	return loan, nil
}

// DeleteLoan deletes a Loan record by its ID
func (s *LoanService) DeleteLoan(id int) error {
	result := s.db.Delete(&models.LoanModel{}, id)
	return result.Error
}

// UpdateLoan updates an existing Loan record
// y (asumiendo que se usa para finalizar el préstamo) vuelve a marcar la pieza como disponible
func (s *LoanService) UpdateLoan(id int, updatedLoan *models.LoanModel) (*models.LoanModel, error) {
	var loan models.LoanModel

	err := s.db.Transaction(func(tx *gorm.DB) error {
		// 1) Obtener el préstamo actual
		if err := tx.First(&loan, id).Error; err != nil {
			return err
		}

		// 2) Asegurar que el ID se mantenga
		updatedLoan.Id = id

		// 3) Actualizar campos del préstamo
		if err := tx.Model(&loan).Updates(updatedLoan).Error; err != nil {
			return err
		}

		// 4) Volver a marcar la pieza como disponible
		if loan.ArtefactId != nil && *loan.ArtefactId != 0 {
			if err := tx.Model(&models.ArtefactModel{}).
				Where("id = ?", *loan.ArtefactId).
				Update("available", true).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Devolver el préstamo actualizado con preload
	if err := s.db.
		Preload("Requester").
		Preload("Artefact").
		Preload("Artefact.InternalClassifier").
		First(&loan, id).Error; err != nil {
		return nil, err
	}

	return &loan, nil
}
