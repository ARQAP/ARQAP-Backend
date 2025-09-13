package controllers

import (
	"net/http"
	"strconv"

	"github.com/ARQAP/ARQAP-Backend/src/models"
	"github.com/ARQAP/ARQAP-Backend/src/services"
	"github.com/gin-gonic/gin"
)

type RegionController struct {
	service *services.RegionService
}

func NewRegionController(service *services.RegionService) *RegionController {
	return &RegionController{service: service}
}

// GetRegions handles GET requests to retrieve all region records
func (c *RegionController) GetRegions(ctx *gin.Context) {
	regions, err := c.service.GetAllRegions()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, regions)
}

// CreateRegion handles POST requests to create a new region record
func (c *RegionController) CreateRegion(ctx *gin.Context) {
	var region models.RegionModel
	if err := ctx.ShouldBindJSON(&region); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	createdRegion, err := c.service.CreateRegion(&region)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, createdRegion)
}

// UpdateRegion handles PUT requests to update an existing region record
func (c *RegionController) UpdateRegion(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var updatedData models.RegionModel
	if err := ctx.ShouldBindJSON(&updatedData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedRegion, err := c.service.UpdateRegion(id, &updatedData)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, updatedRegion)
}

// DeleteRegion handles DELETE requests to remove an region record
func (c *RegionController) DeleteRegion(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := c.service.DeleteRegion(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusNoContent, gin.H{"message": "Region deleted successfully"})
}
