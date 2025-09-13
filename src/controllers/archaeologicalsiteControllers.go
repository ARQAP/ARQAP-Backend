package controllers

import (
	"net/http"
	"strconv"

	"github.com/ARQAP/ARQAP-Backend/src/models"
	"github.com/ARQAP/ARQAP-Backend/src/services"
	"github.com/gin-gonic/gin"
)

type ArchaeologicalSiteController struct {
	service *services.ArchaeologicalSiteService
}

func NewArchaeologicalSiteController(service *services.ArchaeologicalSiteService) *ArchaeologicalSiteController {
	return &ArchaeologicalSiteController{service: service}
}

// GetArchaeologicalSites handles GET requests to retrieve all archaeologicalSite records
func (c *ArchaeologicalSiteController) GetArchaeologicalSites(ctx *gin.Context) {
	archaeologicalSite, err := c.service.GetAllArchaeologicalSites()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, archaeologicalSite)
}

// CreateArchaeologicalSite handles POST requests to create a new archaeologicalSite record
func (c *ArchaeologicalSiteController) CreateArchaeologicalSite(ctx *gin.Context) {
	var archaeologicalSite models.ArchaeologicalSiteModel
	if err := ctx.ShouldBindJSON(&archaeologicalSite); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	createdArchaeologicalSite, err := c.service.CreateArchaeologicalSite(&archaeologicalSite)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, createdArchaeologicalSite)
}

// UpdateArchaeologicalSite handles PUT requests to update an existing archaeologicalSite record
func (c *ArchaeologicalSiteController) UpdateArchaeologicalSite(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var updatedData models.ArchaeologicalSiteModel
	if err := ctx.ShouldBindJSON(&updatedData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedArchaeologicalSite, err := c.service.UpdateArchaeologicalSite(id, &updatedData)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, updatedArchaeologicalSite)
}

// DeleteArchaeologicalSite handles DELETE requests to remove an archaeologicalSite record
func (c *ArchaeologicalSiteController) DeleteArchaeologicalSite(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := c.service.DeleteArchaeologicalSite(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusNoContent, nil)
}
