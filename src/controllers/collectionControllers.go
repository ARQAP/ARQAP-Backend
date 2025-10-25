package controllers

import (
	"net/http"
	"strconv"

	"github.com/ARQAP/ARQAP-Backend/src/models"
	"github.com/ARQAP/ARQAP-Backend/src/services"
	"github.com/gin-gonic/gin"
)

type CollectionController struct {
	service *services.CollectionService
}

func NewCollectionController(service *services.CollectionService) *CollectionController {
	return &CollectionController{service: service}
}

// GetCollections handles GET requests to retrieve all collection records
func (c *CollectionController) GetCollections(ctx *gin.Context) {
	collections, err := c.service.GetAllCollections()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, collections)
}

// CreateCollection handles POST requests to create a new collection record
func (c *CollectionController) CreateCollection(ctx *gin.Context) {
	var collection models.CollectionModel
	if err := ctx.ShouldBindJSON(&collection); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	createdCollection, err := c.service.CreateCollection(&collection)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, createdCollection)
}

// UpdateCollection handles PUT requests to update an existing collection record
func (c *CollectionController) UpdateCollection(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	var updatedData models.CollectionModel
	if err := ctx.ShouldBindJSON(&updatedData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updatedCollection, err := c.service.UpdateCollection(id, &updatedData)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, updatedCollection)
}

// DeleteCollection handles DELETE requests to remove a collection record
func (c *CollectionController) DeleteCollection(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	if err := c.service.DeleteCollection(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusNoContent, nil)
}
