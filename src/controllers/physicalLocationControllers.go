package controllers

import (
	"net/http"
	"strconv"

	"github.com/ARQAP/ARQAP-Backend/src/models"
	"github.com/ARQAP/ARQAP-Backend/src/services"
	"github.com/gin-gonic/gin"
)

type PhysicalLocationController struct {
	service *services.PhysicalLocationService
}

func NewPhysicalLocationController(service *services.PhysicalLocationService) *PhysicalLocationController {
	return &PhysicalLocationController{service: service}
}

// Validation functions

func isValidLevelNumber(level models.LevelNumber) bool {
    return level == models.Level1 || level == models.Level2 || level == models.Level3 || level == models.Level4
}

func isValidColumnLetter(column models.ColumnLetter) bool {
    return column == models.ColumnA || column == models.ColumnB || column == models.ColumnC || column == models.ColumnD
}

// GetAllPhysicalLocations handles GET requests to retrieve all physical locations
func (pc *PhysicalLocationController) GetAllPhysicalLocations(c *gin.Context) {
	locations, err := pc.service.GetAllPhysicalLocations()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, locations)
}

// GetPhysicalLocationByID handles GET
func (pc *PhysicalLocationController) GetPhysicalLocationByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	location, err := pc.service.GetPhysicalLocationByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, location)
}

// CreatePhysicalLocation handles POST requests to create a new physical location
func (pc *PhysicalLocationController) CreatePhysicalLocation(c *gin.Context) {
	var location models.PhysicalLocationModel
	if err := c.ShouldBindJSON(&location); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Validate Level
	if !isValidLevelNumber(location.Level) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid level number"})
		return
	}

	// Validate Column
	if !isValidColumnLetter(location.Column) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid column letter"})
		return
	}

	if err := pc.service.CreatePhysicalLocation(&location); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, location)
}

// UpdatePhysicalLocation handles PUT requests to update an existing physical location
func (pc *PhysicalLocationController) UpdatePhysicalLocation(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var location models.PhysicalLocationModel
	if err := c.ShouldBindJSON(&location); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate Level
	if !isValidLevelNumber(location.Level) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid level number"})
		return
	}

	// Validate Column
	if !isValidColumnLetter(location.Column) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid column letter"})
		return
	}


	location.ID = id
	if err := pc.service.UpdatePhysicalLocation(&location); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, location)
}

// DeletePhysicalLocation handles DELETE requests to remove a physical location
func (pc *PhysicalLocationController) DeletePhysicalLocation(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := pc.service.DeletePhysicalLocation(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}