package controllers

import (
	"net/http"
	"strconv"

	"github.com/ARQAP/ARQAP-Backend/src/models"
	"github.com/ARQAP/ARQAP-Backend/src/services"
	"github.com/gin-gonic/gin"
)

type MentionController struct {
	service *services.MentionService
}

func NewMentionController(service *services.MentionService) *MentionController {
	return &MentionController{service: service}
}

// GetMentions handles GET requests to retrieve all mention records
func (c *MentionController) GetMentions(ctx *gin.Context) {
	mentions, err := c.service.GetAllMentions()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, mentions)
}

// GetMentionByID handles GET requests to retrieve a mention by its ID
func (c *MentionController) GetMentionByID(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	mention, err := c.service.GetMentionByID(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, mention)
}

// CreateMention handles POST requests to create a new mention record
func (c *MentionController) CreateMention(ctx *gin.Context) {
	var mention models.MentionModel
	if err := ctx.ShouldBindJSON(&mention); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	createdMention, err := c.service.CreateMention(&mention)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, createdMention)
}

// UpdateMention handles PUT requests to update an existing mention record
func (c *MentionController) UpdateMention(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	var mention models.MentionModel
	if err := ctx.ShouldBindJSON(&mention); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updatedMention, err := c.service.UpdateMention(id, &mention)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, updatedMention)
}

// DeleteMention handles DELETE requests to remove a mention record
func (c *MentionController) DeleteMention(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	if err := c.service.DeleteMention(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusNoContent, nil)
}