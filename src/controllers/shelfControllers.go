package controllers

import (
	"net/http"
	"strconv"

	"github.com/ARQAP/ARQAP-Backend/src/models"
	"github.com/ARQAP/ARQAP-Backend/src/services"
	"github.com/gin-gonic/gin"
)

type ShelfController struct {
	service *services.ShelfService
}

func NewShelfController(service *services.ShelfService) *ShelfController {
	return &ShelfController{service: service}
}

// GetAllShelves handles GET requests to retrieve all shelf records
func (c *ShelfController) GetAllShelves(ctx *gin.Context) {
	shelves, err := c.service.GetAllShelves()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, shelves)
}

// GetShelfByID handles GET requests to retrieve a shelf record by ID
func (c *ShelfController) GetShelfByID(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid shelf ID"})
		return
	}

	shelf, err := c.service.GetShelfByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, shelf)
}

// CreateShelf handles POST requests to create a new shelf record
func (c *ShelfController) CreateShelf(ctx *gin.Context) {
	var shelf models.ShelfModel
	if err := ctx.ShouldBindJSON(&shelf); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	createdShelf, err := c.service.CreateShelf(&shelf)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, createdShelf)
}

// UpdateShelf handles PUT requests to update an existing shelf record
func (c *ShelfController) UpdateShelf(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid shelf ID"})
		return
	}

	var updatedData models.ShelfModel
	if err := ctx.ShouldBindJSON(&updatedData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedShelf, err := c.service.UpdateShelf(id, &updatedData)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, updatedShelf)
}

// DeleteShelf handles DELETE requests to remove a shelf record
func (c *ShelfController) DeleteShelf(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid shelf ID"})
		return
	}

	if err := c.service.DeleteShelf(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusNoContent, nil)
}