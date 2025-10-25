package controllers

import (
	"io"
	"net/http"
	"strconv"

	"github.com/ARQAP/ARQAP-Backend/src/models"
	"github.com/ARQAP/ARQAP-Backend/src/services"
	"github.com/gin-gonic/gin"
)

type INPLClassifierController struct {
	service *services.INPLService
}

func NewINPLClassifierController(service *services.INPLService) *INPLClassifierController {
	return &INPLClassifierController{service: service}
}

// GetAllINPLClassifiers handles GET requests to retrieve all INPLClassifier records (optionally with photos via ?preload=true)
func (c *INPLClassifierController) GetAllINPLClassifiers(ctx *gin.Context) {
	preload := ctx.DefaultQuery("preload", "false") == "true"
	classifiers, err := c.service.GetAllClassifiers(preload)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, classifiers)
}

// CreateINPLClassifier handles POST requests (multipart) to create a new INPLClassifier with N photos (fichas[] files)
func (c *INPLClassifierController) CreateINPLClassifier(ctx *gin.Context) {
	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid multipart form-data"})
		return
	}
	files := form.File["fichas[]"]
	if len(files) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "upload at least one file in 'fichas[]'"})
		return
	}

	var uploads []services.FichaUpload
	var closers []io.ReadCloser

	for _, fh := range files {
		src, err := fh.Open()
		if err != nil {
			closeAll(closers)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to open one of the uploaded files"})
			return
		}
		closers = append(closers, src)
		uploads = append(uploads, services.FichaUpload{
			Reader:       src,
			OriginalName: fh.Filename,
			ContentType:  fh.Header.Get("Content-Type"),
			Size:         fh.Size,
		})
	}
	defer closeAll(closers)

	classifier, fichas, err := c.service.CreateClassifierWithFichas(uploads)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"classifier": classifier, "fichas": fichas})
}

// DeleteINPLClassifier handles DELETE requests to delete an INPLClassifier by ID (and its photos)
func (c *INPLClassifierController) DeleteINPLClassifier(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid INPLClassifier ID"})
		return
	}
	if err := c.service.DeleteClassifier(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "INPLClassifier deleted successfully"})
}

// UpdateINPLClassifier handles UPDATE requests to update an INPLClassifier by ID (placeholder if no fields exist)
func (c *INPLClassifierController) UpdateINPLClassifier(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid INPLClassifier ID"})
		return
	}
	// If someday the classifier has fields, bind them here.
	var body models.INPLClassifierModel
	_ = ctx.ShouldBindJSON(&body)

	updated, err := c.service.UpdateClassifier(id, &body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, updated)
}

// GetINPLClassifierByID handles GET requests to retrieve a single INPLClassifier by ID (with photos)
func (c *INPLClassifierController) GetINPLClassifierByID(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid INPLClassifier ID"})
		return
	}
	classifier, err := c.service.GetClassifierByID(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, classifier)
}

// AddFichasToINPLClassifier handles POST (multipart) to attach N photos to an existing classifier
func (c *INPLClassifierController) AddFichasToINPLClassifier(ctx *gin.Context) {
	idParam := ctx.Param("id")
	classifierID, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid INPLClassifier ID"})
		return
	}

	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid multipart form-data"})
		return
	}
	files := form.File["fichas[]"]
	if len(files) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "upload at least one file in 'fichas[]'"})
		return
	}

	var uploads []services.FichaUpload
	var closers []io.ReadCloser
	for _, fh := range files {
		src, err := fh.Open()
		if err != nil {
			closeAll(closers)
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to open one of the uploaded files"})
			return
		}
		closers = append(closers, src)
		uploads = append(uploads, services.FichaUpload{
			Reader:       src,
			OriginalName: fh.Filename,
			ContentType:  fh.Header.Get("Content-Type"),
			Size:         fh.Size,
		})
	}
	defer closeAll(closers)

	created, err := c.service.AddFichasToClassifier(classifierID, uploads)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"fichas": created})
}

// ListFichasByINPLClassifier handles GET requests to list all photos by INPLClassifier ID
func (c *INPLClassifierController) ListFichasByINPLClassifier(ctx *gin.Context) {
	idParam := ctx.Param("id")
	classifierID, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid INPLClassifier ID"})
		return
	}
	list, err := c.service.ListFichasByClassifier(classifierID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, list)
}

// ReplaceFicha handles PUT (multipart) to replace a single photo by ficha ID
func (c *INPLClassifierController) ReplaceFicha(ctx *gin.Context) {
	idParam := ctx.Param("id")
	fichaID, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ficha ID"})
		return
	}
	fh, err := ctx.FormFile("ficha")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "missing 'ficha' file"})
		return
	}
	src, err := fh.Open()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to open 'ficha' file"})
		return
	}
	defer src.Close()

	upload := services.FichaUpload{
		Reader:       src,
		OriginalName: fh.Filename,
		ContentType:  fh.Header.Get("Content-Type"),
		Size:         fh.Size,
	}
	updated, err := c.service.ReplaceFicha(fichaID, upload)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, updated)
}

// DeleteFicha handles DELETE requests to remove a single photo by ficha ID
func (c *INPLClassifierController) DeleteFicha(ctx *gin.Context) {
	idParam := ctx.Param("id")
	fichaID, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ficha ID"})
		return
	}
	if err := c.service.DeleteFicha(fichaID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Photo deleted successfully"})
}

// helpers

func closeAll(rr []io.ReadCloser) {
	for _, r := range rr {
		_ = r.Close()
	}
}
