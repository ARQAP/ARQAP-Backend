package routes

import (
	"github.com/ARQAP/ARQAP-Backend/src/controllers"
	"github.com/ARQAP/ARQAP-Backend/src/middleware"
	"github.com/ARQAP/ARQAP-Backend/src/services"
	"github.com/gin-gonic/gin"
)

func SetupCountriesRoutes(router *gin.Engine, service *services.CountryService) {
	countryController := controllers.NewCountryController(service)

	// Protected routes
	country := router.Group("/countries")
	country.Use(middleware.AuthMiddleware())
	{
		country.GET("/", countryController.GetAllCountries)
		country.GET("/:id", countryController.GetCountryByID)
		country.POST("/", countryController.CreateCountry)
		country.PUT("/:id", countryController.UpdateCountry)
		country.DELETE("/:id", countryController.DeleteCountry)
	}
}
