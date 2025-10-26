package controllers

import (
	"net/http"
	"strconv"

	"github.com/ARQAP/ARQAP-Backend/src/models"
	"github.com/ARQAP/ARQAP-Backend/src/services"
	"github.com/gin-gonic/gin"
)

type LoanController struct {
	service *services.LoanService
}

func NewLoanController(service *services.LoanService) *LoanController {
	return &LoanController{service: service}
}

// GetAllLoans handles GET requests to retrieve all loan records
func (c *LoanController) GetAllLoans(ctx *gin.Context) {
	loans, err := c.service.GetAllLoans()									
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, loans)
}

// GetLoanByID handles GET requests to retrieve a loan by its ID
func (c *LoanController) GetLoanByID(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid loan ID"})
		return
	}

	loan, err := c.service.GetLoanByID(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, loan)
}

// CreateLoan handles POST requests to create a new loan record
func (c *LoanController) CreateLoan(ctx *gin.Context) {
	var loan models.LoanModel
	if err := ctx.ShouldBindJSON(&loan); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	createdLoan, err := c.service.CreateLoan(&loan)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, createdLoan)
}

// UpdateLoan handles PUT requests to update an existing loan record
func (c *LoanController) UpdateLoan(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid loan ID"})
		return
	}

	var loan models.LoanModel
	if err := ctx.ShouldBindJSON(&loan); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedLoan, err := c.service.UpdateLoan(id, &loan)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, updatedLoan)
}

// DeleteLoan handles DELETE requests to remove a loan record by its ID
func (c *LoanController) DeleteLoan(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid loan ID"})
		return
	}

	if err := c.service.DeleteLoan(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusNoContent, nil)
}