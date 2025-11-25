package controllers

import (
	"net/http"
	"strconv"

	"github.com/ARQAP/ARQAP-Backend/src/models"
	"github.com/ARQAP/ARQAP-Backend/src/services"
	"github.com/gin-gonic/gin"
)

type InternalMovementController struct {
	service *services.InternalMovementService
}

func NewInternalMovementController(service *services.InternalMovementService) *InternalMovementController {
	return &InternalMovementController{service: service}
}

// GetAllInternalMovements handles GET requests to retrieve all internal movement records
func (c *InternalMovementController) GetAllInternalMovements(ctx *gin.Context) {
	movements, err := c.service.GetAllInternalMovements()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, movements)
}

// GetInternalMovementByID handles GET requests to retrieve an internal movement by its ID
func (c *InternalMovementController) GetInternalMovementByID(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movement ID"})
		return
	}

	movement, err := c.service.GetInternalMovementByID(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, movement)
}

// GetInternalMovementsByArtefactID handles GET requests to retrieve all movements for a specific artefact
func (c *InternalMovementController) GetInternalMovementsByArtefactID(ctx *gin.Context) {
	artefactIdParam := ctx.Param("artefactId")
	artefactId, err := strconv.Atoi(artefactIdParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid artefact ID"})
		return
	}

	movements, err := c.service.GetInternalMovementsByArtefactID(artefactId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, movements)
}

// GetActiveInternalMovementByArtefactID handles GET requests to retrieve the active movement for a specific artefact
func (c *InternalMovementController) GetActiveInternalMovementByArtefactID(ctx *gin.Context) {
	artefactIdParam := ctx.Param("artefactId")
	artefactId, err := strconv.Atoi(artefactIdParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid artefact ID"})
		return
	}

	movement, err := c.service.GetActiveInternalMovementByArtefactID(artefactId)
	if err != nil {
		// Si no hay movimiento activo, devolver 404
		ctx.JSON(http.StatusNotFound, gin.H{"error": "No active movement found"})
		return
	}
	ctx.JSON(http.StatusOK, movement)
}

// CreateInternalMovement handles POST requests to create a new internal movement record
func (c *InternalMovementController) CreateInternalMovement(ctx *gin.Context) {
	var movement models.InternalMovementModel
	if err := ctx.ShouldBindJSON(&movement); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	createdMovement, err := c.service.CreateInternalMovement(&movement)
	if err != nil {
		// Si el error indica que la pieza no está disponible, devolver 400 Bad Request
		if err.Error() == "la pieza arqueológica no está disponible para movimientos internos (ya está prestada)" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, createdMovement)
}

// CreateBatchInternalMovements handles POST requests to create multiple internal movements in a batch
func (c *InternalMovementController) CreateBatchInternalMovements(ctx *gin.Context) {
	var movements []models.InternalMovementModel
	if err := ctx.ShouldBindJSON(&movements); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert to pointers
	movementPointers := make([]*models.InternalMovementModel, len(movements))
	for i := range movements {
		movementPointers[i] = &movements[i]
	}

	createdMovements, err := c.service.CreateBatchInternalMovements(movementPointers)
	if err != nil {
		// Si el error indica que la pieza no está disponible, devolver 400 Bad Request
		if err.Error() == "la pieza arqueológica no está disponible para movimientos internos (ya está prestada)" {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, createdMovements)
}

// UpdateInternalMovement handles PUT requests to update an existing internal movement record
func (c *InternalMovementController) UpdateInternalMovement(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movement ID"})
		return
	}

	var movement models.InternalMovementModel
	if err := ctx.ShouldBindJSON(&movement); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedMovement, err := c.service.UpdateInternalMovement(id, &movement)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, updatedMovement)
}

// DeleteInternalMovement handles DELETE requests to remove an internal movement record by its ID
func (c *InternalMovementController) DeleteInternalMovement(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid movement ID"})
		return
	}

	if err := c.service.DeleteInternalMovement(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusNoContent, nil)
}

