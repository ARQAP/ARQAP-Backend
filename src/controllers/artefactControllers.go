package controllers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/ARQAP/ARQAP-Backend/src/models"
	"github.com/ARQAP/ARQAP-Backend/src/services"
	"github.com/gin-gonic/gin"
)

type ArtefactController struct {
	service *services.ArtefactService
}

func NewArtefactController(service *services.ArtefactService) *ArtefactController {
	return &ArtefactController{service: service}
}

// ======================= ARTEFACTOS COMPLETOS =======================

// TODO: Considerar usar un DTO para optimizar memoria y performance
func (ac *ArtefactController) GetAllArtefacts(c *gin.Context) {
	// Obtener query parameter shelfId si existe
	shelfIdStr := c.Query("shelfId")
	var shelfId *int
	if shelfIdStr != "" {
		parsedId, err := strconv.Atoi(shelfIdStr)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid shelfId parameter"})
			return
		}
		shelfId = &parsedId
	}

	artefacts, err := ac.service.GetAllArtefacts(shelfId)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, artefacts)
}

// TODO: Considerar usar un DTO para optimizar memoria y performance
func (ac *ArtefactController) GetArtefactByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID format"})
		return
	}

	artefact, err := ac.service.GetArtefactByID(id)
	if err != nil {
		c.JSON(404, gin.H{"error": "Artefact not found"})
		return
	}
	c.JSON(200, artefact)
}

func (ac *ArtefactController) CreateArtefact(c *gin.Context) {
	var artefact models.ArtefactModel
	if err := c.ShouldBindJSON(&artefact); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Validaciones obligatorias
	if strings.TrimSpace(artefact.Name) == "" {
		c.JSON(400, gin.H{"error": "El nombre es obligatorio"})
		return
	}

	if strings.TrimSpace(artefact.Material) == "" {
		c.JSON(400, gin.H{"error": "El material es obligatorio"})
		return
	}

	if err := ac.service.CreateArtefact(&artefact); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, artefact)
}

func (ac *ArtefactController) UpdateArtefact(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID format"})
		return
	}

	var artefact models.ArtefactModel
	if err := c.ShouldBindJSON(&artefact); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// Validaciones obligatorias
	if strings.TrimSpace(artefact.Name) == "" {
		c.JSON(400, gin.H{"error": "El nombre es obligatorio"})
		return
	}

	if strings.TrimSpace(artefact.Material) == "" {
		c.JSON(400, gin.H{"error": "El material es obligatorio"})
		return
	}

	if err := ac.service.UpdateArtefact(id, &artefact); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Artefact updated successfully"})
}

func (ac *ArtefactController) DeleteArtefact(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID format"})
		return
	}

	if err := ac.service.DeleteArtefact(id); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Artefact deleted successfully"})
}

// ======================= FOTO =======================

func (ac *ArtefactController) UploadPicture(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID format"})
		return
	}

	// Verify that the artefact exists
	_, err = ac.service.GetArtefactByID(id)
	if err != nil {
		c.JSON(404, gin.H{"error": "Artefact not found"})
		return
	}

	file, header, err := c.Request.FormFile("picture")
	if err != nil {
		c.JSON(400, gin.H{"error": "No file uploaded"})
		return
	}
	defer file.Close()

	// Validate file type
	if !strings.HasPrefix(header.Header.Get("Content-Type"), "image/") {
		c.JSON(400, gin.H{"error": "File must be an image"})
		return
	}

	// Create directories if they don't exist
	uploadDir := "uploads/pictures"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.JSON(500, gin.H{"error": "Could not create upload directory"})
		return
	}

	// Generate unique filename
	filename := fmt.Sprintf("artefact_%d_%d_%s", id, time.Now().Unix(), header.Filename)
	filePath := filepath.Join(uploadDir, filename)

	// Save file
	dst, err := os.Create(filePath)
	if err != nil {
		c.JSON(500, gin.H{"error": "Could not save file"})
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		c.JSON(500, gin.H{"error": "Could not save file"})
		return
	}

	// Save metadata to DB
	picture := models.PictureModel{
		ArtefactID:   id,
		Filename:     filename,
		OriginalName: header.Filename,
		FilePath:     filePath,
		ContentType:  header.Header.Get("Content-Type"),
		Size:         header.Size,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := ac.service.SavePicture(&picture); err != nil {
		// Clean up file if DB save fails
		_ = os.Remove(filePath)
		c.JSON(500, gin.H{"error": "Could not save picture metadata"})
		return
	}

	c.JSON(200, picture)
}

func (ac *ArtefactController) ServePicture(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID format"})
		return
	}

	picture, err := ac.service.GetPictureByArtefactID(id)
	if err != nil || picture == nil {
		c.JSON(404, gin.H{"error": "Picture not found"})
		return
	}

	// Verify that the file exists
	fileInfo, err := os.Stat(picture.FilePath)
	if os.IsNotExist(err) {
		c.JSON(404, gin.H{"error": "Picture file not found"})
		return
	}

	// Cache headers
	lastModified := fileInfo.ModTime().UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT")
	etag := fmt.Sprintf(`"%d-%d"`, picture.ID, picture.UpdatedAt.Unix())

	// Cache for 1 year (images rarely change)
	c.Header("Cache-Control", "public, max-age=31536000") // 1 year
	c.Header("ETag", etag)
	c.Header("Last-Modified", lastModified)

	// Verify If-None-Match (ETag)
	if match := c.GetHeader("If-None-Match"); match == etag {
		c.Status(304) // Not Modified
		return
	}

	// Verify If-Modified-Since
	if modSince := c.GetHeader("If-Modified-Since"); modSince != "" {
		if t, err := time.Parse("Mon, 02 Jan 2006 15:04:05 GMT", modSince); err == nil {
			if !fileInfo.ModTime().After(t) {
				c.Status(304) // Not Modified
				return
			}
		}
	}

	// Serve file with correct content type
	c.Header("Content-Type", picture.ContentType)
	c.File(picture.FilePath)
}

// ======================= DOCUMENTO HISTÓRICO =======================

func (ac *ArtefactController) UploadHistoricalRecord(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID format"})
		return
	}

	// Verify that the artefact exists
	_, err = ac.service.GetArtefactByID(id)
	if err != nil {
		c.JSON(404, gin.H{"error": "Artefact not found"})
		return
	}

	file, header, err := c.Request.FormFile("document")
	if err != nil {
		c.JSON(400, gin.H{"error": "No file uploaded"})
		return
	}
	defer file.Close()

	// Validate file type
	contentType := header.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") && contentType != "application/pdf" {
		c.JSON(400, gin.H{"error": "File must be an image or PDF"})
		return
	}

	// Create directories if they don't exist
	uploadDir := "uploads/historical_records"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.JSON(500, gin.H{"error": "Could not create upload directory"})
		return
	}

	// Generate unique filename
	filename := fmt.Sprintf("record_%d_%d_%s", id, time.Now().Unix(), header.Filename)
	filePath := filepath.Join(uploadDir, filename)

	// Save file
	dst, err := os.Create(filePath)
	if err != nil {
		c.JSON(500, gin.H{"error": "Could not save file"})
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		c.JSON(500, gin.H{"error": "Could not save file"})
		return
	}

	// Save metadata to DB
	record := models.HistoricalRecordModel{
		ArtefactID:   id,
		Filename:     filename,
		OriginalName: header.Filename,
		FilePath:     filePath,
		ContentType:  contentType,
		Size:         header.Size,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := ac.service.SaveHistoricalRecord(&record); err != nil {
		// Clean up file if DB save fails
		_ = os.Remove(filePath)
		c.JSON(500, gin.H{"error": "Could not save document metadata"})
		return
	}

	c.JSON(200, record)
}

func (ac *ArtefactController) ServeHistoricalRecord(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid ID format"})
		return
	}

	record, err := ac.service.GetHistoricalRecordByArtefactID(id)
	if err != nil || record == nil {
		c.JSON(404, gin.H{"error": "Historical record not found"})
		return
	}

	// Verify that the file exists
	fileInfo, err := os.Stat(record.FilePath)
	if os.IsNotExist(err) {
		c.JSON(404, gin.H{"error": "Document file not found"})
		return
	}

	// Cache headers
	lastModified := fileInfo.ModTime().UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT")
	etag := fmt.Sprintf(`"%d-%d"`, record.ID, record.UpdatedAt.Unix())

	// Cache for 1 month (historical documents change less)
	c.Header("Cache-Control", "public, max-age=2592000") // 30 days
	c.Header("ETag", etag)
	c.Header("Last-Modified", lastModified)

	// Verify cache
	if match := c.GetHeader("If-None-Match"); match == etag {
		c.Status(304) // Not Modified
		return
	}

	if modSince := c.GetHeader("If-Modified-Since"); modSince != "" {
		if t, err := time.Parse("Mon, 02 Jan 2006 15:04:05 GMT", modSince); err == nil {
			if !fileInfo.ModTime().After(t) {
				c.Status(304) // Not Modified
				return
			}
		}
	}

	// Serve historical document with correct content type
	c.Header("Content-Type", record.ContentType)
	c.File(record.FilePath)
}

// ======================= ARTEFACTO + MENCIONES =======================

func (ac *ArtefactController) CreateArtefactWithMentions(c *gin.Context) {
	var dto services.CreateArtefactWithMentionsDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	artefact := dto.Artefact

	// Validaciones obligatorias, igual que en CreateArtefact
	if strings.TrimSpace(artefact.Name) == "" {
		c.JSON(400, gin.H{"error": "El nombre es obligatorio"})
		return
	}

	if strings.TrimSpace(artefact.Material) == "" {
		c.JSON(400, gin.H{"error": "El material es obligatorio"})
		return
	}

	created, err := ac.service.CreateArtefactWithMentions(&dto)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, created)
}

// ======================= RESUMENES (ENDPOINT NUEVO) =======================

func (ac *ArtefactController) GetArtefactSummaries(c *gin.Context) {
<<<<<<< HEAD
	// Obtener query parameter shelfId si existe
	shelfIdStr := c.Query("shelfId")
	var shelfId *int
	if shelfIdStr != "" {
=======
	var shelfId *int
	if shelfIdStr := c.Query("shelfId"); shelfIdStr != "" {
>>>>>>> 1d4a55074754b5ea8e380ac4e25e79101d1cb3f3
		parsedId, err := strconv.Atoi(shelfIdStr)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid shelfId parameter"})
			return
		}
		shelfId = &parsedId
	}

	summaries, err := ac.service.GetArtefactSummaries(shelfId)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, summaries)
}

func (c *ArtefactController) ImportArtefactsFromExcel(ctx *gin.Context) {
    file, err := ctx.FormFile("file")
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "error": "no se recibió el archivo",
            "detail": err.Error(),
        })
        return
    }

    f, err := file.Open()
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "error": "no se pudo abrir el archivo",
            "detail": err.Error(),
        })
        return
    }
    defer f.Close()

    result, err := c.service.ImportArtefactsFromExcel(f)
    if err != nil {
        // 👇 manejar el caso en que result sea nil
        if result != nil {
            ctx.JSON(http.StatusBadRequest, gin.H{
                "error":  err.Error(),
                "errors": result.Errors,
            })
        } else {
            ctx.JSON(http.StatusBadRequest, gin.H{
                "error": err.Error(),
            })
        }
        return
    }

    ctx.JSON(http.StatusOK, gin.H{
        "message":  "Importación completada",
        "imported": result.Imported,
        "errors":   result.Errors, // pueden ser warnings de filas puntuales
    })
}