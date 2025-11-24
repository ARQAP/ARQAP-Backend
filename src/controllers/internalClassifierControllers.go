package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"sort"
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

// GetInternalClassifiersByName handles GET requests to retrieve internal classifiers by name
func (c *InternalClassifierController) GetInternalClassifiersByName(ctx *gin.Context) {
	name := ctx.Param("name")
	if name == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "name parameter is required"})
		return
	}
	internalClassifiers, err := c.service.GetInternalClassifiersByName(name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, internalClassifiers)
}

// GetAllInternalClassifierNames handles GET requests to retrieve distinct classifier names
func (c *InternalClassifierController) GetAllInternalClassifierNames(ctx *gin.Context) {
	names, err := c.service.GetAllInternalClassifierNames()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Fallback deduplication in controller to guarantee no repeats
	unique := make(map[string]struct{}, len(names))
	var out []string
	for _, n := range names {
		if _, ok := unique[n]; !ok {
			unique[n] = struct{}{}
			out = append(out, n)
		}
	}
	sort.Strings(out)
	ctx.JSON(http.StatusOK, out)
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
		var dupErr *services.DuplicateNameNumberError
		if errors.As(err, &dupErr) {
			// Diferenciar mensaje según si Number es nil o no
			if dupErr.Number == nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("El clasificador con el nombre '%s' ya se encuentra creado.", dupErr.Name)})
			} else {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("El clasificador '%s' con número %d ya se encuentra creado.", dupErr.Name, *dupErr.Number)})
			}
			return
		}
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
		var dupErr *services.DuplicateNameNumberError
		if errors.As(err, &dupErr) {
			if dupErr.Number == nil {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("El clasificador con el nombre '%s' ya se encuentra creado.", dupErr.Name)})
			} else {
				ctx.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("El clasificador '%s' con número %d ya se encuentra creado.", dupErr.Name, *dupErr.Number)})
			}
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, updatedInternalClassifier)
}

// (GetById removed) GetInternalClassifierByID was removed because lookup by id is not required.
