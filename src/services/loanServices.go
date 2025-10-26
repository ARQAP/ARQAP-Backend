package services

import (
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
	result := s.db.Preload("Requester").Find(&loans)
	return loans, result.Error
}

// GetLoanByID retrieves a Loan record by its ID
func (s *LoanService) GetLoanByID(id int) (*models.LoanModel, error) {
	var loan models.LoanModel
	result := s.db.Preload("Requester").First(&loan, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &loan, nil
}

// CreateLoan creates a new Loan record in the database
func (s *LoanService) CreateLoan(mention *models.LoanModel) (*models.LoanModel, error) {
	result := s.db.Create(mention)
	if result.Error != nil {
		return nil, result.Error
	}
	return mention, nil
}

// DeleteLoan deletes a Loan record by its ID
func (s *LoanService) DeleteLoan(id int) error {
	result := s.db.Delete(&models.LoanModel{}, id)
	return result.Error
}

// UpdateLoan updates an existing Loan record
func (s *LoanService) UpdateLoan(id int, updatedLoan *models.LoanModel) (*models.LoanModel, error) {
    var loan models.LoanModel
    result := s.db.First(&loan, id)
    if result.Error != nil {
        return nil, result.Error
    }
    
    updatedLoan.Id = id

    result = s.db.Model(&loan).Updates(updatedLoan)
    if result.Error != nil {
        return nil, result.Error
    }
    
    result = s.db.First(&loan, id)
    if result.Error != nil {
        return nil, result.Error
    }
    
    return &loan, nil
}