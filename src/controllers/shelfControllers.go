package controllers

import (
	"net/http"
	"strconv"

	"github.com/ARQAP/ARQAP-Backend/src/models"
	"github.com/ARQAP/ARQAP-Backend/src/services"
	"github.com/gin-gonic/gin"
)

type shelfController struct {
	service *services.ShelfService
}

type ShelfController struct {
	service *services.ShelfService
}

func NewShelfController(service *services.ShelfService) *ShelfController {
	return &ShelfController{service: service}
}

// GetAllShelfs handles GET requests to retrieve all shelf records
func (c *ShelfController) GetAllShelfs(ctx *gin.Context) {
	shelfs, err := c.service.GetAllShelfs()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, shelfs)
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

// DeleteShelf handles DELETE requests to delete a shelf record by ID
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
	ctx.JSON(http.StatusOK, gin.H{"message": "Shelf deleted successfully"})
}

// UpdateShelf handles UPDATE requests to update a contry record by ID
func (c *ShelfController) UpdateShelf(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid shelf ID"})
		return
	}
	var shelf models.ShelfModel
	if err := ctx.ShouldBindJSON(&shelf); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updatedShelf, err := c.service.UpdateShelf(id, &shelf)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, updatedShelf)
}

// GetShelfByID handles GET request to retrive a shelf record by ID
func (c *ShelfController) GetShelfByID(ctx *gin.Context){
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid shelf ID"})
	}
	shelf, err := c.service.GetShelfByID(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	ctx.JSON(http.StatusOK, shelf)
}