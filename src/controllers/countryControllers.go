package controllers

import (
	"net/http"
	"strconv"

	"github.com/ARQAP/ARQAP-Backend/src/models"
	"github.com/ARQAP/ARQAP-Backend/src/services"
	"github.com/gin-gonic/gin"
)

type CountryController struct {
	service *services.CountryService
}

func NewCountryController(service *services.CountryService) *CountryController {
	return &CountryController{service: service}
}

// GetAllCountries handles GET requests to retrieve all country records
func (c *CountryController) GetAllCountries(ctx *gin.Context) {
	countries, err := c.service.GetAllCountries()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, countries)
}

// CreateCountry handles POST requests to create a new country record
func (c *CountryController) CreateCountry(ctx *gin.Context) {
    var country models.CountryModel
    if err := ctx.ShouldBindJSON(&country); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    createdCountry, err := c.service.CreateCountry(&country)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    ctx.JSON(http.StatusCreated, createdCountry)
}

// DeleteCountry handles DELETE requests to delete a country record by ID
func (c *CountryController) DeleteCountry(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid country ID"})
		return
	}
	if err := c.service.DeleteCountry(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Country deleted successfully"})
}

// UpdateCountry handles UPDATE requests to update a contry record by ID
func (c *CountryController) UpdateCountry(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid country ID"})
		return
	}
	var country models.CountryModel
	if err := ctx.ShouldBindJSON(&country); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updatedCountry, err := c.service.UpdateCountry(id, &country)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, updatedCountry)
}

// GetCountryByID handles GET request to retrive a country record by ID
func (c *CountryController) GetCountryByID(ctx *gin.Context){
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid country ID"})
	}
	country, err := c.service.GetCountryByID(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	ctx.JSON(http.StatusOK, country)
}