package controllers

import (
	"net/http"
	"strconv"

	"github.com/ARQAP/ARQAP-Backend/src/models"
	"github.com/ARQAP/ARQAP-Backend/src/services"
	"github.com/gin-gonic/gin"
)

type InternalClassifierController struct {
	service *services.InternalClassifierService
}

func NewInternalClassifierController(service *services.InternalClassifierService) *InternalClassifierController {
	return &InternalClassifierController{service: service}
}

// GetAllInternalClassifiers handles GET requests to retrieve all internalClassifier records
func (c *InternalClassifierController) GetAllInternalClassifiers(ctx *gin.Context) {
	internalClassifiers, err := c.service.GetAllInternalClassifiers()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, internalClassifiers)
}

// CreateInternalClassifier handles POST requests to create a new internalClassifier record
func (c *InternalClassifierController) CreateInternalClassifier(ctx *gin.Context) {
	var internalClassifier models.InternalClassifierModel
	if err := ctx.ShouldBindJSON(&internalClassifier); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	createdInternalClassifier, err := c.service.CreateInternalClassifier(&internalClassifier)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, createdInternalClassifier)
}

// DeleteInternalClassifier handles DELETE requests to delete a internalClassifier record by ID
func (c *InternalClassifierController) DeleteInternalClassifier(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid internalClassifier ID"})
		return
	}
	if err := c.service.DeleteInternalClassifier(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "InternalClassifier deleted successfully"})
}

// UpdateInternalClassifier handles UPDATE requests to update a contry record by ID
func (c *InternalClassifierController) UpdateInternalClassifier(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid InternalClassifier ID"})
		return
	}
	var internalClassifier models.InternalClassifierModel
	if err := ctx.ShouldBindJSON(&internalClassifier); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updatedInternalClassifier, err := c.service.UpdateInternalClassifier(id, &internalClassifier)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, updatedInternalClassifier)
}

// GetInternalClassifierByID handles GET request to retrive a internalClassifier record by ID
func (c *InternalClassifierController) GetInternalClassifierByID(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid InternalClassifier ID"})
	}
	internalClassifier, err := c.service.GetInternalClassifierByID(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	ctx.JSON(http.StatusOK, internalClassifier)
}
