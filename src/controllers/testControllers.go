package controllers

import (
	"net/http"

	"github.com/ARQAP/ARQAP-Backend/src/models"
	"github.com/ARQAP/ARQAP-Backend/src/services"
	"github.com/gin-gonic/gin"
)

type TestController struct {
    service *services.TestService
}

func NewTestController(service *services.TestService) *TestController {
    return &TestController{service: service}
}

// GetTests handles GET requests to retrieve all test records
func (c *TestController) GetTests(ctx *gin.Context) {
	tests, err := c.service.GetAllTests()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, tests)
}

// CreateTests handles POST requests to create a new test record
func (c *TestController) CreateTests(ctx *gin.Context) {
	var test models.TestModel
	if err := ctx.ShouldBindJSON(&test); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	createdTest, err := c.service.CreateTest(&test)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, createdTest)
}