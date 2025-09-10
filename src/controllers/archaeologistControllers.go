package controllers

import (
	"net/http"
	"strconv"

	"github.com/ARQAP/ARQAP-Backend/src/models"
	"github.com/ARQAP/ARQAP-Backend/src/services"
	"github.com/gin-gonic/gin"
)

type ArchaeologistController struct {
    service *services.ArchaeologistService
}

func NewArchaeologistController(service *services.ArchaeologistService) *ArchaeologistController {
    return &ArchaeologistController{service: service}
}

// GetArchaeologists handles GET requests to retrieve all archaeologist records
func (c *ArchaeologistController) GetArchaeologists(ctx *gin.Context) {
	archaeologists, err := c.service.GetAllArchaeologists()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, archaeologists)
}

// CreateArchaeologist handles POST requests to create a new archaeologist record
func (c *ArchaeologistController) CreateArchaeologist(ctx *gin.Context) {
	var archaeologist models.ArchaeologistModel
	if err := ctx.ShouldBindJSON(&archaeologist); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	createdArchaeologist, err := c.service.CreateArchaeologist(&archaeologist)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, createdArchaeologist)
}

// UpdateArchaeologist handles PUT requests to update an existing archaeologist record
func (c *ArchaeologistController) UpdateArchaeologist(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var updatedData models.ArchaeologistModel
	if err := ctx.ShouldBindJSON(&updatedData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedArchaeologist, err := c.service.UpdateArchaeologist(id, &updatedData)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, updatedArchaeologist)
}

// DeleteArchaeologist handles DELETE requests to remove an archaeologist record
func (c *ArchaeologistController) DeleteArchaeologist(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := c.service.DeleteArchaeologist(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusNoContent, nil)
}