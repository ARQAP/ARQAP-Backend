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

// FichaDTO is the data transfer object for INPLFicha
type FichaDTO struct {
	ID               int    `json:"id"`
	INPLClassifierID int    `json:"inplClassifierId"`
	Filename         string `json:"filename"`
	ContentType      string `json:"contentType"`
	Size             int64  `json:"size"`
	URL              string `json:"url"`
}

// INPLClassifierDTO is the data transfer object for INPLClassifierModel
type INPLClassifierDTO struct {
	ID         int        `json:"id"`
	INPLFichas []FichaDTO `json:"inplFichas"`
}

// toDTO converts an INPLClassifierModel to INPLClassifierDTO
func (c *INPLClassifierController) toDTO(ctx *gin.Context, m *models.INPLClassifierModel) INPLClassifierDTO {
	out := INPLClassifierDTO{ID: m.ID}
	out.INPLFichas = make([]FichaDTO, 0, len(m.INPLFichas))
	base := baseURL(ctx)
	for _, f := range m.INPLFichas {
		out.INPLFichas = append(out.INPLFichas, FichaDTO{
			ID:               f.ID,
			INPLClassifierID: f.INPLClassifierID,
			Filename:         f.Filename,
			ContentType:      f.ContentType,
			Size:             f.Size,
			URL:              base + "/inplFichas/" + strconv.Itoa(f.ID) + "/download",
		})
	}
	return out
}

// baseURL constructs the base URL from the request context
func baseURL(ctx *gin.Context) string {
	scheme := "http"
	if ctx.Request.TLS != nil || ctx.Request.Header.Get("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}
	return scheme + "://" + ctx.Request.Host
}

// GetAllINPLClassifiers retrieves all INPLClassifiers
func (c *INPLClassifierController) GetAllINPLClassifiers(ctx *gin.Context) {
	preload := ctx.DefaultQuery("preload", "false") == "true"
	classifiers, err := c.service.GetAllClassifiers(preload)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !preload {
		ctx.JSON(http.StatusOK, classifiers)
		return
	}
	dtos := make([]INPLClassifierDTO, 0, len(classifiers))
	for i := range classifiers {
		dtos = append(dtos, c.toDTO(ctx, &classifiers[i]))
	}
	ctx.JSON(http.StatusOK, dtos)
}

// CreateINPLClassifier creates a new INPLClassifier
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

	classifier, _, err := c.service.CreateClassifierWithFichas(uploads)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	reloaded, err := c.service.GetClassifierByID(classifier.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, c.toDTO(ctx, reloaded))
}

// UpdateINPLClassifier updates an existing INPLClassifier
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

// UpdateINPLClassifier updates an existing INPLClassifier
func (c *INPLClassifierController) UpdateINPLClassifier(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid INPLClassifier ID"})
		return
	}
	var body models.INPLClassifierModel
	_ = ctx.ShouldBindJSON(&body)

	updated, err := c.service.UpdateClassifier(id, &body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, updated)
}

// GetINPLClassifierByID retrieves an INPLClassifier by its ID
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
	ctx.JSON(http.StatusOK, c.toDTO(ctx, classifier))
}

// ListFichasByINPLClassifier lists all fichas associated with a given INPLClassifier
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

	if _, err := c.service.AddFichasToClassifier(classifierID, uploads); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	reloaded, err := c.service.GetClassifierByID(classifierID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, c.toDTO(ctx, reloaded))
}

// ListFichasByINPLClassifier lists all fichas associated with a given INPLClassifier
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

// ReplaceFicha replaces an existing ficha with a new uploaded file
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

// DeleteFicha deletes an INPLFicha by its ID
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

// DownloadFicha serves the file associated with the given ficha ID
func (c *INPLClassifierController) DownloadFicha(ctx *gin.Context) {
	idParam := ctx.Param("id")
	fichaID, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	f, err := c.service.GetFichaByID(fichaID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	ctx.Header("Content-Disposition", `inline; filename="`+f.Filename+`"`)
	ctx.Header("Content-Type", f.ContentType)
	ctx.File(f.FilePath)
}

// closeAll closes all provided ReadClosers
func closeAll(rr []io.ReadCloser) {
	for _, r := range rr {
		_ = r.Close()
	}
}
