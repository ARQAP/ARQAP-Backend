package controllers

import (
	"net/http"
	"strconv"

	"github.com/ARQAP/ARQAP-Backend/src/models"
	"github.com/ARQAP/ARQAP-Backend/src/services"
	"github.com/gin-gonic/gin"
)

type RequesterController struct {
	service *services.RequesterService
}

func NewRequesterController(service *services.RequesterService) *RequesterController {
	return &RequesterController{service: service}
}

// GetAllRequesters handles GET requests to retrieve all requester records
func (c *RequesterController) GetAllRequesters(ctx *gin.Context) {
	requesters, err := c.service.GetAllRequesters()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, requesters)
}

// GetRequesterByID handles GET requests to retrieve a requester by its ID
func (c *RequesterController) GetRequesterByID(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid requester ID"})
		return
	}

	requester, err := c.service.GetRequesterByID(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, requester)
}

// CreateRequester handles POST requests to create a new requester record
func (c *RequesterController) CreateRequester(ctx *gin.Context) {
	var requester models.RequesterModel
	if err := ctx.ShouldBindJSON(&requester); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate requester type
	if requester.Type == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "type is required"})
		return
	}

	// Validate that type is one of the allowed values
	if requester.Type != models.Investigator &&
		requester.Type != models.Department &&
		requester.Type != models.Exhibition {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid type value"})
		return
	}

	createdRequester, err := c.service.CreateRequester(&requester)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, createdRequester)
}

// UpdateRequester handles PUT requests to update an existing requester record
func (c *RequesterController) UpdateRequester(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid requester ID"})
		return
	}

	var requester models.RequesterModel
	if err := ctx.ShouldBindJSON(&requester); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedRequester, err := c.service.UpdateRequester(id, &requester)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, updatedRequester)
}

// DeleteRequester handles DELETE requests to remove a requester record by its ID
func (c *RequesterController) DeleteRequester(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid requester ID"})
		return
	}
	if err := c.service.DeleteRequester(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusNoContent, nil)
}
