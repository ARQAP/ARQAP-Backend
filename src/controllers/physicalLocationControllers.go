package controllers

import (
	"net/http"
	"strconv"

	"github.com/ARQAP/ARQAP-Backend/src/models"
	"github.com/ARQAP/ARQAP-Backend/src/services"
	"github.com/gin-gonic/gin"
)

type PhysicalLocationController struct {
	service *services.PhysicalLocationService
}

func NewPhysicalLocationController(service *services.PhysicalLocationService) *PhysicalLocationController {
	return &PhysicalLocationController{service: service}
}