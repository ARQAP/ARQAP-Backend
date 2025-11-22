package services

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	urlpkg "net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ARQAP/ARQAP-Backend/src/dtos"
	"github.com/ARQAP/ARQAP-Backend/src/models"
	"github.com/ARQAP/ARQAP-Backend/src/utils"
	excelize "github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

// Cache entry
type CacheEntry struct {
	Data      interface{}
	ExpiresAt time.Time
}

type ImportResult struct {
	Imported int
	Errors   []string
}

type ArtefactService struct {
	db    *gorm.DB
	cache map[string]*CacheEntry
	mutex sync.RWMutex
}

type CreateArtefactWithMentionsDTO struct {
	Artefact models.ArtefactModel  `json:"artefact"`
	Mentions []models.MentionModel `json:"mentions"`
}

func NewArtefactService(db *gorm.DB) *ArtefactService {
	service := &ArtefactService{
		db:    db,
		cache: make(map[string]*CacheEntry),
	}

	// Clean up cache every 30 minutes
	go service.cleanupCache()

	return service
}

func (s *ArtefactService) cleanupCache() {
	ticker := time.NewTicker(30 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		s.mutex.Lock()
		now := time.Now()
		for key, entry := range s.cache {
			if now.After(entry.ExpiresAt) {
				delete(s.cache, key)
			}
		}
		s.mutex.Unlock()
	}
}

func (s *ArtefactService) setCache(key string, data interface{}, duration time.Duration) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.cache[key] = &CacheEntry{
		Data:      data,
		ExpiresAt: time.Now().Add(duration),
	}
}

func (s *ArtefactService) getCache(key string) (interface{}, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	entry, exists := s.cache[key]
	if !exists || time.Now().After(entry.ExpiresAt) {
		return nil, false
	}

	return entry.Data, true
}

func (s *ArtefactService) invalidateCache(pattern string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for key := range s.cache {
		if contains(key, pattern) {
			delete(s.cache, key)
		}
	}
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		(len(s) > len(substr) && s[:len(substr)] == substr))
}

// ======================= ARTEFACTOS COMPLETOS =======================

func (s *ArtefactService) GetAllArtefacts(shelfId *int) ([]models.ArtefactModel, error) {
	// Si hay filtro por shelf, no usar cache general
	var cacheKey string
	if shelfId != nil {
		cacheKey = fmt.Sprintf("artefacts_shelf_%d", *shelfId)
	} else {
		cacheKey = "all_artefacts"
	}

	// Try to get from cache
	if cached, found := s.getCache(cacheKey); found {
		return cached.([]models.ArtefactModel), nil
	}

	// If not in cache, query DB
	var artefacts []models.ArtefactModel
	query := s.db.Preload("Picture").
		Preload("HistoricalRecord").
		Preload("Archaeologist").
		Preload("ArchaeologicalSite").
		Preload("PhysicalLocation").
		Preload("Collection").
		Preload("InplClassifier").
		Preload("InternalClassifier").
		Preload("PhysicalLocation.Shelf")

	// Aplicar filtro por shelfId si está presente
	if shelfId != nil {
		query = query.Joins("JOIN physical_location_models ON physical_location_models.id = artefact_models.physical_location_id").
			Where("physical_location_models.shelf_id = ?", *shelfId)
	}

	err := query.Find(&artefacts).Error

	if err == nil {
		// Save to cache for 5 minutes
		s.setCache(cacheKey, artefacts, 5*time.Minute)
	}

	return artefacts, err
}

func (s *ArtefactService) GetArtefactByID(id int) (*models.ArtefactModel, error) {
	cacheKey := fmt.Sprintf("artefact_%d", id)

	// Try to get from cache
	if cached, found := s.getCache(cacheKey); found {
		artefact := cached.(models.ArtefactModel)
		return &artefact, nil
	}

	// If not in cache, query DB
	var artefact models.ArtefactModel

	err := s.db.Preload("Picture").
		Preload("HistoricalRecord").
		Preload("Archaeologist").
		Preload("ArchaeologicalSite").
		Preload("PhysicalLocation").
		Preload("Collection").
		Preload("InplClassifier").
		Preload("InternalClassifier").
		Preload("PhysicalLocation.Shelf").
		First(&artefact, id).Error
	if err != nil {
		return nil, err
	}

	// Save to cache for 10 minutes
	s.setCache(cacheKey, artefact, 10*time.Minute)

	return &artefact, nil
}

func (s *ArtefactService) CreateArtefact(artefact *models.ArtefactModel) error {
	if err := s.db.Create(artefact).Error; err != nil {
		return err
	}

	// Invalidate related cache
	s.invalidateCache("all_artefacts")
	s.invalidateCache("artefact_summaries")

	return nil
}

func (s *ArtefactService) UpdateArtefact(id int, artefact *models.ArtefactModel) error {
	if err := s.db.Where("id = ?", id).Updates(artefact).Error; err != nil {
		return err
	}

	// Invalidate specific and general cache
	s.invalidateCache(fmt.Sprintf("artefact_%d", id))
	s.invalidateCache("all_artefacts")
	s.invalidateCache("artefact_summaries")

	return nil
}

func (s *ArtefactService) DeleteArtefact(id int) error {
	// Obtain artefact with its associated Picture and HistoricalRecord paths
	var artefact models.ArtefactModel
	if err := s.db.Preload("Picture").Preload("HistoricalRecord").First(&artefact, id).Error; err != nil {
		return err
	}

	// Delete the files if they exist
	if len(artefact.Picture) > 0 && artefact.Picture[0].FilePath != "" {
		_ = os.Remove(artefact.Picture[0].FilePath)
	}
	if len(artefact.HistoricalRecord) > 0 && artefact.HistoricalRecord[0].FilePath != "" {
		_ = os.Remove(artefact.HistoricalRecord[0].FilePath)
	}

	// Delete from DB
	if err := s.db.Delete(&artefact, id).Error; err != nil {
		return err
	}

	// Invalidate cache
	s.invalidateCache(fmt.Sprintf("artefact_%d", id))
	s.invalidateCache("all_artefacts")
	s.invalidateCache("artefact_summaries")

	return nil
}

// ======================= FOTO Y DOCUMENTO =======================

func (s *ArtefactService) GetPictureByArtefactID(artefactID int) (*models.PictureModel, error) {
	cacheKey := fmt.Sprintf("picture_%d", artefactID)
	if cached, found := s.getCache(cacheKey); found {
		pic := cached.(models.PictureModel)
		return &pic, nil
	}

	var picture models.PictureModel
	err := s.db.Where("artefact_id = ?", artefactID).First(&picture).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// caso normal: aún no hay foto
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	s.setCache(cacheKey, picture, 30*time.Minute)
	return &picture, nil
}

func (s *ArtefactService) GetHistoricalRecordByArtefactID(artefactID int) (*models.HistoricalRecordModel, error) {
	cacheKey := fmt.Sprintf("historical_record_%d", artefactID)

	// Try to get from cache
	if cached, found := s.getCache(cacheKey); found {
		record := cached.(models.HistoricalRecordModel)
		return &record, nil
	}

	// If not in cache, query DB
	var record models.HistoricalRecordModel
	err := s.db.Where("artefact_id = ?", artefactID).First(&record).Error
	if err != nil {
		return nil, err
	}

	// Save to cache for 30 minutes
	s.setCache(cacheKey, record, 30*time.Minute)

	return &record, nil
}

func (s *ArtefactService) SavePicture(picture *models.PictureModel) error {
	var existing models.PictureModel
	err := s.db.Where("artefact_id = ?", picture.ArtefactID).First(&existing).Error

	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		// no había foto: crear
		if err := s.db.Create(picture).Error; err != nil {
			return err
		}
	case err != nil:
		// error real de DB
		return err
	default:
		// había foto: borrar archivo viejo y actualizar registro
		if existing.FilePath != "" {
			_ = os.Remove(existing.FilePath)
		}
		// aseguramos update sobre el registro existente
		picture.ID = existing.ID
		if err := s.db.Model(&existing).Updates(map[string]interface{}{
			"file_path":    picture.FilePath,
			"content_type": picture.ContentType,
			"size":         picture.Size,
			"updated_at":   time.Now(),
		}).Error; err != nil {
			return err
		}
	}

	s.invalidateCache(fmt.Sprintf("picture_%d", picture.ArtefactID))
	s.invalidateCache(fmt.Sprintf("artefact_%d", picture.ArtefactID))
	s.invalidateCache("all_artefacts")
	s.invalidateCache("artefact_summaries")
	return nil
}

func (s *ArtefactService) SaveHistoricalRecord(record *models.HistoricalRecordModel) error {
	// Check if a document already exists for this artefact
	var existing models.HistoricalRecordModel
	if err := s.db.Where("artefact_id = ?", record.ArtefactID).First(&existing).Error; err == nil {
		// Already exists, delete previous file
		_ = os.Remove(existing.FilePath)
		// Update existing record
		if err := s.db.Where("artefact_id = ?", record.ArtefactID).Updates(record).Error; err != nil {
			return err
		}
	} else {
		// Does not exist, create new
		if err := s.db.Create(record).Error; err != nil {
			return err
		}
	}

	// Invalidate related cache
	s.invalidateCache(fmt.Sprintf("historical_record_%d", record.ArtefactID))
	s.invalidateCache(fmt.Sprintf("artefact_%d", record.ArtefactID))
	s.invalidateCache("all_artefacts")
	s.invalidateCache("artefact_summaries")

	return nil
}

// ======================= ARTEFACTO + MENCIONES =======================

func (s *ArtefactService) CreateArtefactWithMentions(dto *CreateArtefactWithMentionsDTO) (*models.ArtefactModel, error) {
	artefact := dto.Artefact

	err := s.db.Transaction(func(tx *gorm.DB) error {
		// 1) Crear artefacto
		if err := tx.Create(&artefact).Error; err != nil {
			return err
		}

		// 2) Crear menciones (si hay)
		if len(dto.Mentions) > 0 {
			var mentionsToCreate []models.MentionModel

			for _, m := range dto.Mentions {
				// evitar menciones totalmente vacías
				if strings.TrimSpace(m.Title) == "" && strings.TrimSpace(m.Link) == "" {
					continue
				}
				m.ArtefactId = &artefact.ID
				mentionsToCreate = append(mentionsToCreate, m)
			}

			if len(mentionsToCreate) > 0 {
				if err := tx.Create(&mentionsToCreate).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Invalidate cache (igual que CreateArtefact, pero más completo)
	s.invalidateCache("all_artefacts")
	s.invalidateCache("artefact_summaries")
	s.invalidateCache(fmt.Sprintf("artefact_%d", artefact.ID))

	return &artefact, nil
}

// ======================= RESÚMENES LIVIANOS =======================

func (s *ArtefactService) GetArtefactSummaries(shelfId *int) ([]dtos.ArtefactSummaryDTO, error) {
	// Cache key dinámico según el filtro
	cacheKey := "artefact_summaries"
	if shelfId != nil {
		cacheKey = fmt.Sprintf("artefact_summaries_shelf_%d", *shelfId)
	}

	if cached, found := s.getCache(cacheKey); found {
		return cached.([]dtos.ArtefactSummaryDTO), nil
	}

	type summaryRow struct {
		ID                     int
		Name                   string
		Material               string
		CollectionName         *string `gorm:"column:collection_name"`
		ArchaeologistFirstName *string `gorm:"column:archaeologist_first_name"`
		ArchaeologistLastName  *string `gorm:"column:archaeologist_last_name"`
		ArchaeologicalSiteName *string `gorm:"column:archaeological_site_name"`
		ShelfCode              *int    `gorm:"column:shelf_code"`
		Level                  *int    `gorm:"column:level"`
		Column                 *string `gorm:"column:column"`
	}

	var rows []summaryRow

	query := s.db.Table("artefact_models AS a").
		Select(`a.id,
			a.name,
			a.material,
			c.name AS collection_name,
			ar.firstname AS archaeologist_first_name,
			ar.lastname AS archaeologist_last_name,
			site."Name" AS archaeological_site_name,
			sh.code AS shelf_code,
			pl.level AS level,
			pl.column AS column`).
		Joins("LEFT JOIN collection_models c ON c.id = a.collection_id").
		Joins("LEFT JOIN archaeologist_models ar ON ar.id = a.archaeologist_id").
		Joins("LEFT JOIN archaeological_site_models site ON site.id = a.archaeological_site_id").
		Joins("LEFT JOIN physical_location_models pl ON pl.id = a.physical_location_id").
		Joins("LEFT JOIN shelf_models sh ON sh.id = pl.shelf_id")

	if shelfId != nil {
		query = query.Where("pl.shelf_id = ?", *shelfId)
	}

	if err := query.Scan(&rows).Error; err != nil {
		return nil, err
	}

	summaries := make([]dtos.ArtefactSummaryDTO, 0, len(rows))

	for _, row := range rows {
		dto := dtos.ArtefactSummaryDTO{
			ID:                     row.ID,
			Name:                   row.Name,
			Material:               row.Material,
			CollectionName:         row.CollectionName,
			ArchaeologicalSiteName: row.ArchaeologicalSiteName,
			ShelfCode:              row.ShelfCode,
			Level:                  row.Level,
			Column:                 row.Column,
		}

		var nameParts []string
		if row.ArchaeologistFirstName != nil {
			first := strings.TrimSpace(*row.ArchaeologistFirstName)
			if first != "" {
				nameParts = append(nameParts, first)
			}
		}
		if row.ArchaeologistLastName != nil {
			last := strings.TrimSpace(*row.ArchaeologistLastName)
			if last != "" {
				nameParts = append(nameParts, last)
			}
		}
		if len(nameParts) > 0 {
			fullName := strings.Join(nameParts, " ")
			dto.ArchaeologistName = &fullName
		}

		summaries = append(summaries, dto)
	}

	s.setCache(cacheKey, summaries, 5*time.Minute)

	return summaries, nil
}

func (s *ArtefactService) ImportArtefactsFromExcel(r io.Reader) (*ImportResult, error) {
	log.Println("[IMPORT] Iniciando importación de artefactos desde Excel...")

	f, err := excelize.OpenReader(r)
	if err != nil {
		return nil, fmt.Errorf("archivo excel inválido: %w", err)
	}
	defer f.Close()

	rows, err := f.GetRows("BASE BRUCH DIEGO")
	if err != nil {
		return nil, fmt.Errorf("no se pudo leer la hoja BASE BRUCH DIEGO: %w", err)
	}

	log.Printf("[IMPORT] Total de filas en BASE BRUCH DIEGO: %d", len(rows))
	result := &ImportResult{Imported: 0, Errors: []string{}}

	// ===============================
	// 0) Leer hoja TOPOGRÁFICO y crear mapa código -> archivos
	// ===============================
	_, err = f.GetRows("TOPOGRÁFICO")
	if err != nil {
		// Si no existe la hoja, continuamos sin asociar archivos
		result.Errors = append(result.Errors, "Advertencia: no se encontró la hoja TOPOGRÁFICO, se importarán artefactos sin archivos asociados")
	}

	// Mapa: código numérico -> {foto, ficha histórica, ficha INPL}
	// Estructura real: Col A = ACLARACIONES (texto), Col B = FOTO PIEZA (código/hipervínculo), Col C = FICHA INPL (código/hipervínculo), Col D = FICHA HST. (código/hipervínculo)
	// Los hipervínculos en B, C y D apuntan a los archivos reales
	topoMap := make(map[string]struct {
		foto           string
		fichaHistorica string
		fichaINPL      string
	})

	if err == nil {
		log.Println("[IMPORT] Leyendo hoja TOPOGRÁFICO para asociar archivos...")
		// Leer hipervínculos de las celdas en lugar de solo el texto
		sheetName := "TOPOGRÁFICO"
		rows, err := f.GetRows(sheetName)
		if err == nil {
			topoCount := 0
			for rowIdx, row := range rows {
				// Saltar encabezados (primera fila) y filas vacías
				if rowIdx < 2 || len(row) < 2 {
					continue
				}

				// Obtener código de la columna B (índice 1, fila rowIdx+1 porque GetRows es 0-indexed pero las celdas son 1-indexed)
				cellB := fmt.Sprintf("B%d", rowIdx+1)
				codigo, _ := f.GetCellValue(sheetName, cellB)
				codigo = strings.TrimSpace(codigo)

				// Saltar si no hay código o es texto descriptivo
				if codigo == "" || !regexp.MustCompile(`^\d+$`).MatchString(codigo) {
					continue
				}

				// Extraer hipervínculos de las celdas B, C y D
				fotoPath := ""
				inplPath := ""
				historicaPath := ""

				// Columna B - Foto
				cellBHyperlink, target, _ := f.GetCellHyperLink(sheetName, cellB)
				if cellBHyperlink && target != "" {
					fotoPath = target
				}

				// Columna C - INPL
				cellC := fmt.Sprintf("C%d", rowIdx+1)
				cellCHyperlink, target, _ := f.GetCellHyperLink(sheetName, cellC)
				if cellCHyperlink && target != "" {
					inplPath = target
				}

				// Columna D - Histórica
				cellD := fmt.Sprintf("D%d", rowIdx+1)
				cellDHyperlink, target, _ := f.GetCellHyperLink(sheetName, cellD)
				if cellDHyperlink && target != "" {
					historicaPath = target
				}

				// Guardar en el mapa
				topoMap[codigo] = struct {
					foto           string
					fichaHistorica string
					fichaINPL      string
				}{
					foto:           fotoPath,
					fichaHistorica: historicaPath,
					fichaINPL:      inplPath,
				}
				topoCount++
			}
			log.Printf("[IMPORT] Procesados %d códigos en hoja TOPOGRÁFICO", topoCount)
		}
	}

	// ===============================
	// 1) Asegurar Colección "Bruch"
	// ===============================
	var bruchCollection models.CollectionModel
	if err := s.db.
		Where("name = ?", "Colección Bruch").
		First(&bruchCollection).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			bruchCollection = models.CollectionModel{
				Name: "Colección Bruch",
			}
			if err := s.db.Create(&bruchCollection).Error; err != nil {
				return nil, fmt.Errorf("no se pudo crear la colección Bruch: %w", err)
			}
		} else {
			return nil, fmt.Errorf("error buscando colección Bruch: %w", err)
		}
	}
	collectionID := bruchCollection.Id // references:Id en tu GORM

	// ==========================================
	// 2) Cache en memoria de arqueólogos por nombre
	//     key: nombre completo leído del Excel
	//     value: ID en la base
	// ==========================================
	archaeologistCache := make(map[string]int)

	// ===============================
	// 3) Recorrer filas del Excel
	// ===============================
	for i, row := range rows {
		// Fila vacía o sin código → la salto
		if len(row) == 0 || strings.TrimSpace(row[0]) == "" {
			continue
		}

		// ---------------------------------
		// 3.1. Nombre completo del arqueólogo desde Excel
		// ---------------------------------
		var archaeologistID *int

		if len(row) > 1 {
			fullName := strings.TrimSpace(row[1]) // ej: "Carlos Bruch"

			if fullName != "" {
				// Primero miro en cache para no pegarle siempre a la BD
				if id, ok := archaeologistCache[fullName]; ok {
					idCopy := id
					archaeologistID = &idCopy
				} else {
					// No está en cache, buscar/crear en la base
					// Partimos el nombre en FirstName + LastName (muy simple)
					firstName := fullName
					lastName := ""
					parts := strings.Fields(fullName)
					if len(parts) > 1 {
						firstName = parts[0]
						lastName = strings.Join(parts[1:], " ")
					}

					var arch models.ArchaeologistModel
					err := s.db.
						Where(&models.ArchaeologistModel{
							FirstName: firstName,
							LastName:  lastName,
						}).
						First(&arch).Error

					if errors.Is(err, gorm.ErrRecordNotFound) {
						// Crear nuevo arqueólogo
						arch = models.ArchaeologistModel{
							FirstName: firstName,
							LastName:  lastName,
						}
						if err := s.db.Create(&arch).Error; err != nil {
							result.Errors = append(result.Errors, fmt.Sprintf(
								"Fila %d: no se pudo crear arqueólogo %s: %v",
								i+1, fullName, err,
							))
							// sigo con la fila, pero sin archaeologist_id
						} else {
							archaeologistCache[fullName] = arch.Id
							idCopy := arch.Id
							archaeologistID = &idCopy
						}
					} else if err != nil {
						// error distinto a not found
						result.Errors = append(result.Errors, fmt.Sprintf(
							"Fila %d: error buscando arqueólogo %s: %v",
							i+1, fullName, err,
						))
					} else {
						// Encontrado correctamente
						archaeologistCache[fullName] = arch.Id
						idCopy := arch.Id
						archaeologistID = &idCopy
					}
				}
			}
		}

		// ---------------------------------
		// 3.2. Datos del artefacto
		// ---------------------------------

		// Col 0: código inventario → Name
		name := strings.TrimSpace(row[0])

		// Col 7: material
		material := ""
		if len(row) > 7 {
			material = strings.TrimSpace(row[7])
		}

		// Col 8–13: tipología + procedencia → Description
		typology := ""
		if len(row) > 8 {
			typology = strings.TrimSpace(row[8])
		}

		region, country, province, locality, site := "", "", "", "", ""
		if len(row) > 9 {
			region = strings.TrimSpace(row[9])
		}
		if len(row) > 10 {
			country = strings.TrimSpace(row[10])
		}
		if len(row) > 11 {
			province = strings.TrimSpace(row[11])
		}
		if len(row) > 12 {
			locality = strings.TrimSpace(row[12])
		}
		if len(row) > 13 {
			site = strings.TrimSpace(row[13])
		}

		// armamos una descripción básica
		descParts := []string{}
		if typology != "" {
			descParts = append(descParts, typology)
		}
		if region != "" {
			descParts = append(descParts, region)
		}

		loc := []string{}
		if site != "" {
			loc = append(loc, site)
		}
		if locality != "" {
			loc = append(loc, locality)
		}
		if province != "" {
			loc = append(loc, province)
		}
		if country != "" {
			loc = append(loc, country)
		}
		if len(loc) > 0 {
			descParts = append(descParts, strings.Join(loc, ", "))
		}

		var description *string
		if len(descParts) > 0 {
			d := strings.Join(descParts, " – ")
			description = &d
		}

		artefact := models.ArtefactModel{
			Name:         name,
			Material:     material,
			Available:    true,
			Description:  description,
			CollectionID: &collectionID,
			// puede ir nil si no pudimos crear/buscar al arqueólogo
			ArchaeologistID: archaeologistID,
		}

		if err := s.db.Create(&artefact).Error; err != nil {
			log.Printf("[IMPORT] ERROR en fila %d: %v", i+1, err)
			result.Errors = append(result.Errors, fmt.Sprintf("Fila %d: %v", i+1, err))
			continue
		}

		log.Printf("[IMPORT] Artefacto creado: %s (ID: %d)", name, artefact.ID)

		// ===============================
		// 3.3. Asociar archivos desde hoja TOPOGRÁFICO
		// ===============================
		// Extraer número del código (ej: 5788 de "MLP-Ar-CB-5788")
		re := regexp.MustCompile(`(\d+)$`)
		matches := re.FindStringSubmatch(name)
		if len(matches) > 1 && len(topoMap) > 0 {
			codigoNum := matches[1]
			if topoData, found := topoMap[codigoNum]; found {
				// Asociar foto - si hay ruta directa la usamos, sino buscamos por código
				if topoData.foto != "" {
					if strings.HasPrefix(topoData.foto, "http://") || strings.HasPrefix(topoData.foto, "https://") || strings.Contains(topoData.foto, string(filepath.Separator)) {
						// Es una ruta de archivo o URL
						log.Printf("[IMPORT] Descargando foto para %s desde: %s", name, topoData.foto)
						if err := s.associatePictureFromPath(&artefact, topoData.foto); err != nil {
							log.Printf("[IMPORT] ERROR asociando foto para %s: %v", name, err)
							result.Errors = append(result.Errors, fmt.Sprintf("Fila %d: error asociando foto: %v", i+1, err))
						} else {
							log.Printf("[IMPORT] Foto asociada exitosamente para %s", name)
						}
					} else {
						// Es solo un código, buscar archivo
						if err := s.associatePictureFromCode(&artefact, topoData.foto); err != nil {
							result.Errors = append(result.Errors, fmt.Sprintf("Fila %d: error asociando foto: %v", i+1, err))
						}
					}
				}

				// Asociar ficha histórica
				if topoData.fichaHistorica != "" {
					if strings.HasPrefix(topoData.fichaHistorica, "http://") || strings.HasPrefix(topoData.fichaHistorica, "https://") || strings.Contains(topoData.fichaHistorica, string(filepath.Separator)) {
						log.Printf("[IMPORT] Descargando ficha histórica para %s desde: %s", name, topoData.fichaHistorica)
						if err := s.associateHistoricalRecordFromPath(&artefact, topoData.fichaHistorica); err != nil {
							log.Printf("[IMPORT] ERROR asociando ficha histórica para %s: %v", name, err)
							result.Errors = append(result.Errors, fmt.Sprintf("Fila %d: error asociando ficha histórica: %v", i+1, err))
						} else {
							log.Printf("[IMPORT] Ficha histórica asociada exitosamente para %s", name)
						}
					} else {
						if err := s.associateHistoricalRecordFromCode(&artefact, topoData.fichaHistorica); err != nil {
							result.Errors = append(result.Errors, fmt.Sprintf("Fila %d: error asociando ficha histórica: %v", i+1, err))
						}
					}
				}

				// Asociar ficha INPL
				if topoData.fichaINPL != "" {
					if strings.HasPrefix(topoData.fichaINPL, "http://") || strings.HasPrefix(topoData.fichaINPL, "https://") || strings.Contains(topoData.fichaINPL, string(filepath.Separator)) {
						log.Printf("[IMPORT] Descargando ficha INPL para %s desde: %s", name, topoData.fichaINPL)
						if err := s.associateINPLFromPath(&artefact, topoData.fichaINPL); err != nil {
							log.Printf("[IMPORT] ERROR asociando ficha INPL para %s: %v", name, err)
							result.Errors = append(result.Errors, fmt.Sprintf("Fila %d: error asociando ficha INPL: %v", i+1, err))
						} else {
							log.Printf("[IMPORT] Ficha INPL asociada exitosamente para %s", name)
						}
					} else {
						if err := s.associateINPLFromCode(&artefact, topoData.fichaINPL); err != nil {
							result.Errors = append(result.Errors, fmt.Sprintf("Fila %d: error asociando ficha INPL: %v", i+1, err))
						}
					}
				}
			}
		}

		result.Imported++
	}

	s.invalidateCache("all_artefacts")

	log.Printf("[IMPORT] ========================================")
	log.Printf("[IMPORT] Importación completada")
	log.Printf("[IMPORT] Artefactos importados: %d", result.Imported)
	log.Printf("[IMPORT] Errores encontrados: %d", len(result.Errors))
	if len(result.Errors) > 0 {
		log.Printf("[IMPORT] Primeros 5 errores:")
		for i, err := range result.Errors {
			if i >= 5 {
				break
			}
			log.Printf("[IMPORT]   - %s", err)
		}
		if len(result.Errors) > 5 {
			log.Printf("[IMPORT]   ... y %d errores más", len(result.Errors)-5)
		}
	}
	log.Printf("[IMPORT] ========================================")

	if result.Imported == 0 && len(result.Errors) > 0 {
		return result, fmt.Errorf("no se pudo importar ninguna pieza")
	}

	return result, nil
}

// ===============================
// Funciones auxiliares para asociar archivos
// ===============================

// min retorna el mínimo de dos enteros
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// downloadFileFromURL descarga un archivo desde una URL y lo guarda temporalmente
// Retorna la ruta del archivo temporal y el nombre original del archivo (si se pudo determinar)
func downloadFileFromURL(url string) (string, string, error) {
	originalURL := url
	log.Printf("[DOWNLOAD] Iniciando descarga desde: %s", originalURL)

	// Si es una URL de Google Drive, usar la API
	if utils.IsGoogleDriveURL(url) {
		log.Printf("[DOWNLOAD] Detectada URL de Google Drive, usando API...")

		// Extraer ID del archivo
		fileID, err := utils.ExtractFileIDFromURL(url)
		if err != nil {
			return "", "", fmt.Errorf("error extrayendo ID de archivo: %w", err)
		}

		// Descargar usando la API
		fileBody, filename, err := utils.DownloadFileFromGoogleDrive(fileID)
		if err != nil {
			return "", "", fmt.Errorf("error descargando archivo desde Google Drive API: %w", err)
		}
		defer fileBody.Close()

		// Crear archivo temporal
		tmpFile, err := os.CreateTemp("", "download_*")
		if err != nil {
			return "", "", fmt.Errorf("error creando archivo temporal: %w", err)
		}
		defer tmpFile.Close()

		// Copiar contenido al archivo temporal
		_, err = io.Copy(tmpFile, fileBody)
		if err != nil {
			os.Remove(tmpFile.Name())
			return "", "", fmt.Errorf("error copiando contenido: %w", err)
		}

		log.Printf("[DOWNLOAD] Archivo descargado exitosamente: %s", filename)
		return tmpFile.Name(), filename, nil
	}

	// Crear cliente HTTP con timeout y cookies para mantener sesión
	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Timeout: 30 * time.Second,
		Jar:     jar,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// Seguir redirecciones automáticamente
			return nil
		},
	}

	// Descargar el archivo
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", "", fmt.Errorf("error creando request: %w", err)
	}

	// Agregar headers para evitar bloqueos
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	log.Printf("[DOWNLOAD] Realizando petición HTTP...")
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[DOWNLOAD] ERROR en petición HTTP: %v", err)
		return "", "", fmt.Errorf("error descargando archivo: %w", err)
	}
	defer resp.Body.Close()

	log.Printf("[DOWNLOAD] Respuesta recibida: código %d, Content-Type: %s", resp.StatusCode, resp.Header.Get("Content-Type"))

	// Google Drive puede devolver 200 pero con HTML de confirmación para archivos grandes
	if resp.StatusCode != http.StatusOK {
		log.Printf("[DOWNLOAD] ERROR: código de estado %d", resp.StatusCode)
		return "", "", fmt.Errorf("error descargando archivo: código de estado %d", resp.StatusCode)
	}

	// Verificar content-type para detectar si Google Drive devolvió HTML en lugar del archivo
	contentType := resp.Header.Get("Content-Type")
	if strings.Contains(contentType, "text/html") && strings.Contains(url, "drive.google.com") {
		log.Printf("[DOWNLOAD] Google Drive devolvió HTML (página de confirmación), intentando extraer token...")

		// Leer el body HTML
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", "", fmt.Errorf("error leyendo respuesta HTML: %w", err)
		}
		bodyStr := string(bodyBytes)

		// Buscar el token de confirmación en el HTML
		// Google Drive usa diferentes patrones, intentamos varios
		var confirmToken string

		// Patrón 1: name="confirm" value="TOKEN" (formulario)
		re1 := regexp.MustCompile(`name=["']confirm["']\s+value=["']([^"']+)["']`)
		matches1 := re1.FindStringSubmatch(bodyStr)
		if len(matches1) > 1 {
			confirmToken = matches1[1]
			log.Printf("[DOWNLOAD] Token encontrado (patrón 1): %s", confirmToken[:min(20, len(confirmToken))])
		}

		// Patrón 2: id="download-form" con input hidden
		if confirmToken == "" {
			re2 := regexp.MustCompile(`id=["']download-form["'][^>]*>.*?<input[^>]*name=["']confirm["'][^>]*value=["']([^"']+)["']`)
			matches2 := re2.FindStringSubmatch(bodyStr)
			if len(matches2) > 1 {
				confirmToken = matches2[1]
				log.Printf("[DOWNLOAD] Token encontrado (patrón 2): %s", confirmToken[:min(20, len(confirmToken))])
			}
		}

		// Patrón 3: confirm=TOKEN en la URL de descarga (href)
		if confirmToken == "" {
			re3 := regexp.MustCompile(`href=["']([^"']*uc[^"']*[?&]confirm=([^"']+)[^"']*)["']`)
			matches3 := re3.FindStringSubmatch(bodyStr)
			if len(matches3) > 2 {
				confirmToken = matches3[2]
				log.Printf("[DOWNLOAD] Token encontrado (patrón 3): %s", confirmToken[:min(20, len(confirmToken))])
			}
		}

		// Patrón 4: window.location con confirm
		if confirmToken == "" {
			re4 := regexp.MustCompile(`window\.location\s*=\s*["']([^"']*[?&]confirm=([^"']+)[^"']*)["']`)
			matches4 := re4.FindStringSubmatch(bodyStr)
			if len(matches4) > 2 {
				confirmToken = matches4[2]
				log.Printf("[DOWNLOAD] Token encontrado (patrón 4): %s", confirmToken[:min(20, len(confirmToken))])
			}
		}

		// Patrón 5: /uc?export=download&id=FILE_ID&confirm=TOKEN
		if confirmToken == "" {
			re5 := regexp.MustCompile(`/uc\?export=download[^"']*[&?]confirm=([a-zA-Z0-9_-]+)`)
			matches5 := re5.FindStringSubmatch(bodyStr)
			if len(matches5) > 1 {
				confirmToken = matches5[1]
				log.Printf("[DOWNLOAD] Token encontrado (patrón 5): %s", confirmToken[:min(20, len(confirmToken))])
			}
		}

		// Patrón 6: Buscar cualquier input con name="confirm"
		if confirmToken == "" {
			re6 := regexp.MustCompile(`<input[^>]*name=["']confirm["'][^>]*value=["']([^"']+)["']`)
			matches6 := re6.FindStringSubmatch(bodyStr)
			if len(matches6) > 1 {
				confirmToken = matches6[1]
				log.Printf("[DOWNLOAD] Token encontrado (patrón 6): %s", confirmToken[:min(20, len(confirmToken))])
			}
		}

		if confirmToken == "" {
			log.Printf("[DOWNLOAD] No se pudo extraer token de confirmación del HTML")
			log.Printf("[DOWNLOAD] Guardando muestra del HTML para debug (primeros 500 caracteres)...")
			htmlPreview := bodyStr
			if len(htmlPreview) > 500 {
				htmlPreview = htmlPreview[:500]
			}
			log.Printf("[DOWNLOAD] HTML preview: %s", htmlPreview)

			// Intentar con confirm=t como fallback
			log.Printf("[DOWNLOAD] Intentando con confirm=t...")
			parsedURL, err := urlpkg.Parse(url)
			if err == nil {
				q := parsedURL.Query()
				q.Set("confirm", "t")
				parsedURL.RawQuery = q.Encode()
				url = parsedURL.String()

				// Hacer nueva petición con confirm=t
				req2, err := http.NewRequest("GET", url, nil)
				if err != nil {
					return "", "", fmt.Errorf("error creando request con confirm: %w", err)
				}
				req2.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
				// Copiar cookies de la primera petición
				for _, cookie := range client.Jar.Cookies(parsedURL) {
					req2.AddCookie(cookie)
				}

				resp2, err := client.Do(req2)
				if err != nil {
					return "", "", fmt.Errorf("error en segunda petición: %w", err)
				}
				defer resp2.Body.Close()

				if resp2.StatusCode != http.StatusOK {
					return "", "", fmt.Errorf("error en segunda petición: código %d", resp2.StatusCode)
				}

				// Verificar si ahora es el archivo real
				contentType2 := resp2.Header.Get("Content-Type")
				if strings.Contains(contentType2, "text/html") {
					// Intentar una vez más leyendo el body para ver si hay un token
					bodyBytes2, _ := io.ReadAll(resp2.Body)
					bodyStr2 := string(bodyBytes2)

					// Buscar token nuevamente con todos los patrones
					reFinal := regexp.MustCompile(`name=["']confirm["']\s+value=["']([^"']+)["']`)
					matchesFinal := reFinal.FindStringSubmatch(bodyStr2)
					if len(matchesFinal) > 1 {
						confirmToken = matchesFinal[1]
						log.Printf("[DOWNLOAD] Token encontrado en segunda respuesta, continuando con POST...")
						// Continuar con el flujo POST más abajo (salir del if y usar el token)
					} else {
						return "", "", fmt.Errorf("google drive aún requiere confirmación después de intentar confirm=t")
					}
				} else {
					// Usar la nueva respuesta - archivo descargado exitosamente
					resp = resp2
					contentType = contentType2
					log.Printf("[DOWNLOAD] Segunda petición exitosa con confirm=t")
					// Salir del bloque de confirmación y continuar con el procesamiento normal
					confirmToken = "" // Marcar que ya no necesitamos POST
				}
			} else {
				return "", "", fmt.Errorf("no se pudo parsear URL para agregar confirm: %w", err)
			}
		}

		// Si tenemos un token (ya sea del HTML original o de la segunda respuesta), hacer POST
		if confirmToken != "" {
			log.Printf("[DOWNLOAD] Token de confirmación encontrado, haciendo petición POST...")

			// Hacer petición POST con el token
			parsedURL, err := urlpkg.Parse(url)
			if err != nil {
				return "", "", fmt.Errorf("error parseando URL: %w", err)
			}

			// Crear formulario con el token
			formData := urlpkg.Values{}
			formData.Set("confirm", confirmToken)

			req2, err := http.NewRequest("POST", parsedURL.String(), strings.NewReader(formData.Encode()))
			if err != nil {
				return "", "", fmt.Errorf("error creando request POST: %w", err)
			}
			req2.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
			req2.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			resp2, err := client.Do(req2)
			if err != nil {
				return "", "", fmt.Errorf("error en petición POST: %w", err)
			}
			defer resp2.Body.Close()

			if resp2.StatusCode != http.StatusOK {
				return "", "", fmt.Errorf("error en petición POST: código %d", resp2.StatusCode)
			}

			// Verificar si ahora es el archivo real
			contentType2 := resp2.Header.Get("Content-Type")
			if strings.Contains(contentType2, "text/html") {
				return "", "", fmt.Errorf("google drive aún requiere confirmación después de enviar token")
			}

			// Usar la nueva respuesta
			resp = resp2
			contentType = contentType2
			log.Printf("[DOWNLOAD] Petición POST exitosa, descargando archivo...")
		}
	}

	// Intentar extraer el nombre del archivo del header Content-Disposition
	var filename string
	contentDisposition := resp.Header.Get("Content-Disposition")
	if contentDisposition != "" {
		// Buscar filename= o filename*=
		if idx := strings.Index(contentDisposition, "filename="); idx != -1 {
			filename = contentDisposition[idx+9:]
			// Remover comillas si las hay
			filename = strings.Trim(filename, "\"'")
			// Si tiene encoding (filename*=), intentar decodificar
			filename = strings.TrimPrefix(filename, "UTF-8''")
		}
	}

	// Si no se pudo obtener del header, intentar de la URL
	if filename == "" {
		parsedURL, err := urlpkg.Parse(originalURL)
		if err == nil {
			// Intentar obtener de la query string
			if name := parsedURL.Query().Get("name"); name != "" {
				filename = name
			} else {
				// Obtener del path
				path := parsedURL.Path
				if path != "" && path != "/" {
					filename = filepath.Base(path)
				}
			}
		}
	}

	// Determinar extensión basada en Content-Type si no se tiene
	ext := filepath.Ext(filename)
	if ext == "" {
		switch {
		case strings.Contains(contentType, "image/jpeg") || strings.Contains(contentType, "image/jpg"):
			ext = ".jpg"
		case strings.Contains(contentType, "image/png"):
			ext = ".png"
		case strings.Contains(contentType, "image/gif"):
			ext = ".gif"
		case strings.Contains(contentType, "image/webp"):
			ext = ".webp"
		case strings.Contains(contentType, "application/pdf"):
			ext = ".pdf"
		case strings.Contains(contentType, "application/zip"):
			ext = ".zip"
		default:
			ext = ".bin"
		}
	}

	// Si no tenemos nombre, usar uno genérico con la extensión
	if filename == "" {
		filename = "download" + ext
	} else if filepath.Ext(filename) == "" {
		// Si el nombre no tiene extensión, agregarla
		filename = filename + ext
	}

	// Crear archivo temporal con la extensión correcta
	tmpFile, err := os.CreateTemp("", "download_*"+ext)
	if err != nil {
		return "", "", fmt.Errorf("error creando archivo temporal: %w", err)
	}

	// Copiar contenido
	log.Printf("[DOWNLOAD] Copiando contenido a archivo temporal...")
	if _, err := io.Copy(tmpFile, resp.Body); err != nil {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
		log.Printf("[DOWNLOAD] ERROR copiando contenido: %v", err)
		return "", "", fmt.Errorf("error copiando archivo: %w", err)
	}

	// Cerrar el archivo para asegurar que se escribió todo
	if err := tmpFile.Close(); err != nil {
		os.Remove(tmpFile.Name())
		log.Printf("[DOWNLOAD] ERROR cerrando archivo: %v", err)
		return "", "", fmt.Errorf("error cerrando archivo temporal: %w", err)
	}

	log.Printf("[DOWNLOAD] Descarga completada: %s (nombre: %s)", tmpFile.Name(), filename)
	return tmpFile.Name(), filename, nil
}

// findFileByCode busca un archivo por código numérico en un directorio base
// Intenta diferentes patrones de nombres: codigo.ext, foto_codigo.ext, codigo_foto.ext, etc.
func findFileByCode(baseDir string, codigo string, extensions []string) (string, error) {
	if baseDir == "" {
		// Si no hay directorio base configurado, intentar directorios comunes
		baseDir = "archivos_bruch"
	}

	// Patrones de nombres a intentar
	patterns := []string{
		codigo,                          // 5788.jpg
		fmt.Sprintf("foto_%s", codigo),  // foto_5788.jpg
		fmt.Sprintf("%s_foto", codigo),  // 5788_foto.jpg
		fmt.Sprintf("pieza_%s", codigo), // pieza_5788.jpg
	}

	for _, pattern := range patterns {
		for _, ext := range extensions {
			fullPath := filepath.Join(baseDir, pattern+ext)
			if _, err := os.Stat(fullPath); err == nil {
				return fullPath, nil
			}
		}
	}

	return "", fmt.Errorf("archivo no encontrado para código %s en %s", codigo, baseDir)
}

// associatePictureFromCode busca y asocia una foto usando el código numérico
func (s *ArtefactService) associatePictureFromCode(artefact *models.ArtefactModel, codigo string) error {
	// Buscar archivo por código
	extensions := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}
	baseDir := os.Getenv("BRUCH_FILES_BASE_DIR") // Directorio base configurable
	if baseDir == "" {
		baseDir = "archivos_bruch"
	}

	sourcePath, err := findFileByCode(baseDir, codigo, extensions)
	if err != nil {
		return err // Archivo no encontrado, pero no es crítico
	}

	return s.associatePictureFromPath(artefact, sourcePath)
}

// associateHistoricalRecordFromCode busca y asocia una ficha histórica usando el código numérico
func (s *ArtefactService) associateHistoricalRecordFromCode(artefact *models.ArtefactModel, codigo string) error {
	// Buscar archivo por código
	extensions := []string{".pdf", ".jpg", ".jpeg", ".png"}
	baseDir := os.Getenv("BRUCH_FILES_BASE_DIR")
	if baseDir == "" {
		baseDir = "archivos_bruch"
	}

	// Patrones específicos para fichas históricas
	patterns := []string{
		codigo,
		fmt.Sprintf("historica_%s", codigo),
		fmt.Sprintf("%s_historica", codigo),
		fmt.Sprintf("ficha_%s", codigo),
		fmt.Sprintf("hst_%s", codigo),
	}

	for _, pattern := range patterns {
		for _, ext := range extensions {
			fullPath := filepath.Join(baseDir, pattern+ext)
			if _, err := os.Stat(fullPath); err == nil {
				return s.associateHistoricalRecordFromPath(artefact, fullPath)
			}
		}
	}

	return fmt.Errorf("ficha histórica no encontrada para código %s", codigo)
}

// associateINPLFromCode busca y asocia una ficha INPL usando el código numérico
func (s *ArtefactService) associateINPLFromCode(artefact *models.ArtefactModel, codigo string) error {
	// Buscar archivo por código
	extensions := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}
	baseDir := os.Getenv("BRUCH_FILES_BASE_DIR")
	if baseDir == "" {
		baseDir = "archivos_bruch"
	}

	// Patrones específicos para fichas INPL
	patterns := []string{
		codigo,
		fmt.Sprintf("inpl_%s", codigo),
		fmt.Sprintf("%s_inpl", codigo),
		fmt.Sprintf("ficha_inpl_%s", codigo),
	}

	for _, pattern := range patterns {
		for _, ext := range extensions {
			fullPath := filepath.Join(baseDir, pattern+ext)
			if _, err := os.Stat(fullPath); err == nil {
				return s.associateINPLFromPath(artefact, fullPath)
			}
		}
	}

	return fmt.Errorf("ficha INPL no encontrada para código %s", codigo)
}

// associatePictureFromPath copia un archivo de imagen y lo asocia al artefacto
func (s *ArtefactService) associatePictureFromPath(artefact *models.ArtefactModel, sourcePath string) error {
	// Si es una URL, descargarla primero
	var actualPath string
	var originalFilename string
	var shouldDeleteTemp bool

	if strings.HasPrefix(sourcePath, "http://") || strings.HasPrefix(sourcePath, "https://") {
		downloadedPath, filename, err := downloadFileFromURL(sourcePath)
		if err != nil {
			return fmt.Errorf("error descargando archivo desde URL: %w", err)
		}
		actualPath = downloadedPath
		originalFilename = filename
		shouldDeleteTemp = true
		defer func() {
			if shouldDeleteTemp {
				os.Remove(actualPath)
			}
		}()
	} else {
		actualPath = sourcePath
		originalFilename = filepath.Base(sourcePath)
		// Verificar si el archivo existe
		if _, err := os.Stat(actualPath); os.IsNotExist(err) {
			return fmt.Errorf("archivo no encontrado: %s", actualPath)
		}
	}

	// Leer el archivo
	sourceFile, err := os.Open(actualPath)
	if err != nil {
		return fmt.Errorf("no se pudo abrir el archivo: %w", err)
	}
	defer sourceFile.Close()

	// Obtener información del archivo
	fileInfo, err := sourceFile.Stat()
	if err != nil {
		return fmt.Errorf("no se pudo obtener información del archivo: %w", err)
	}

	// Determinar content type basado en la extensión del archivo original
	ext := strings.ToLower(filepath.Ext(originalFilename))
	if ext == "" {
		ext = strings.ToLower(filepath.Ext(actualPath))
	}
	contentType := "image/jpeg"
	switch ext {
	case ".png":
		contentType = "image/png"
	case ".gif":
		contentType = "image/gif"
	case ".webp":
		contentType = "image/webp"
	}

	// Crear directorio de destino
	uploadDir := "uploads/pictures"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return fmt.Errorf("no se pudo crear directorio: %w", err)
	}

	// Generar nombre único usando el nombre original del archivo
	filename := fmt.Sprintf("artefact_%d_%d_%s", artefact.ID, time.Now().Unix(), originalFilename)
	destPath := filepath.Join(uploadDir, filename)

	// Copiar archivo
	destFile, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("no se pudo crear archivo destino: %w", err)
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, sourceFile); err != nil {
		os.Remove(destPath)
		return fmt.Errorf("no se pudo copiar archivo: %w", err)
	}

	// Crear registro en BD
	picture := models.PictureModel{
		ArtefactID:   artefact.ID,
		Filename:     filename,
		OriginalName: originalFilename,
		FilePath:     destPath,
		ContentType:  contentType,
		Size:         fileInfo.Size(),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.SavePicture(&picture); err != nil {
		os.Remove(destPath)
		return fmt.Errorf("no se pudo guardar metadata: %w", err)
	}

	// Si descargamos un archivo temporal, marcarlo para no eliminarlo (ya se copió)
	shouldDeleteTemp = false
	return nil
}

// associateHistoricalRecordFromPath copia un archivo de ficha histórica y lo asocia al artefacto
func (s *ArtefactService) associateHistoricalRecordFromPath(artefact *models.ArtefactModel, sourcePath string) error {
	// Si es una URL, descargarla primero
	var actualPath string
	var originalFilename string
	var shouldDeleteTemp bool

	if strings.HasPrefix(sourcePath, "http://") || strings.HasPrefix(sourcePath, "https://") {
		downloadedPath, filename, err := downloadFileFromURL(sourcePath)
		if err != nil {
			return fmt.Errorf("error descargando archivo desde URL: %w", err)
		}
		actualPath = downloadedPath
		originalFilename = filename
		shouldDeleteTemp = true
		defer func() {
			if shouldDeleteTemp {
				os.Remove(actualPath)
			}
		}()
	} else {
		actualPath = sourcePath
		originalFilename = filepath.Base(sourcePath)
		// Verificar si el archivo existe
		if _, err := os.Stat(actualPath); os.IsNotExist(err) {
			return fmt.Errorf("archivo no encontrado: %s", actualPath)
		}
	}

	// Leer el archivo
	sourceFile, err := os.Open(actualPath)
	if err != nil {
		return fmt.Errorf("no se pudo abrir el archivo: %w", err)
	}
	defer sourceFile.Close()

	// Obtener información del archivo
	fileInfo, err := sourceFile.Stat()
	if err != nil {
		return fmt.Errorf("no se pudo obtener información del archivo: %w", err)
	}

	// Determinar content type basado en la extensión del archivo original
	ext := strings.ToLower(filepath.Ext(originalFilename))
	if ext == "" {
		ext = strings.ToLower(filepath.Ext(actualPath))
	}
	contentType := "application/pdf"
	if strings.HasPrefix(ext, ".jpg") || strings.HasPrefix(ext, ".jpeg") {
		contentType = "image/jpeg"
	} else if ext == ".png" {
		contentType = "image/png"
	}

	// Crear directorio de destino
	uploadDir := "uploads/historical_records"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return fmt.Errorf("no se pudo crear directorio: %w", err)
	}

	// Generar nombre único usando el nombre original del archivo
	filename := fmt.Sprintf("record_%d_%d_%s", artefact.ID, time.Now().Unix(), originalFilename)
	destPath := filepath.Join(uploadDir, filename)

	// Copiar archivo
	destFile, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("no se pudo crear archivo destino: %w", err)
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, sourceFile); err != nil {
		os.Remove(destPath)
		return fmt.Errorf("no se pudo copiar archivo: %w", err)
	}

	// Crear registro en BD
	record := models.HistoricalRecordModel{
		ArtefactID:   artefact.ID,
		Filename:     filename,
		OriginalName: originalFilename,
		FilePath:     destPath,
		ContentType:  contentType,
		Size:         fileInfo.Size(),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.SaveHistoricalRecord(&record); err != nil {
		os.Remove(destPath)
		return fmt.Errorf("no se pudo guardar metadata: %w", err)
	}

	// Si descargamos un archivo temporal, marcarlo para no eliminarlo (ya se copió)
	shouldDeleteTemp = false
	return nil
}

// associateINPLFromPath copia un archivo de ficha INPL y lo asocia al artefacto
func (s *ArtefactService) associateINPLFromPath(artefact *models.ArtefactModel, sourcePath string) error {
	// Si es una URL, descargarla primero
	var actualPath string
	var originalFilename string
	var shouldDeleteTemp bool

	if strings.HasPrefix(sourcePath, "http://") || strings.HasPrefix(sourcePath, "https://") {
		downloadedPath, filename, err := downloadFileFromURL(sourcePath)
		if err != nil {
			return fmt.Errorf("error descargando archivo desde URL: %w", err)
		}
		actualPath = downloadedPath
		originalFilename = filename
		shouldDeleteTemp = true
		defer func() {
			if shouldDeleteTemp {
				os.Remove(actualPath)
			}
		}()
	} else {
		actualPath = sourcePath
		originalFilename = filepath.Base(sourcePath)
		// Verificar si el archivo existe
		if _, err := os.Stat(actualPath); os.IsNotExist(err) {
			return fmt.Errorf("archivo no encontrado: %s", actualPath)
		}
	}

	// Leer el archivo
	sourceFile, err := os.Open(actualPath)
	if err != nil {
		return fmt.Errorf("no se pudo abrir el archivo: %w", err)
	}
	defer sourceFile.Close()

	// Obtener información del archivo
	fileInfo, err := sourceFile.Stat()
	if err != nil {
		return fmt.Errorf("no se pudo obtener información del archivo: %w", err)
	}

	// Determinar content type basado en la extensión del archivo original (INPL son imágenes)
	ext := strings.ToLower(filepath.Ext(originalFilename))
	if ext == "" {
		ext = strings.ToLower(filepath.Ext(actualPath))
	}
	contentType := "image/jpeg"
	switch ext {
	case ".png":
		contentType = "image/png"
	case ".gif":
		contentType = "image/gif"
	case ".webp":
		contentType = "image/webp"
	}

	// Crear o obtener INPLClassifier para este artefacto
	var inplClassifier models.INPLClassifierModel
	if artefact.InplClassifierID != nil {
		// Ya tiene un INPLClassifier, usarlo
		if err := s.db.First(&inplClassifier, *artefact.InplClassifierID).Error; err != nil {
			return fmt.Errorf("no se pudo encontrar INPLClassifier existente: %w", err)
		}
	} else {
		// Crear nuevo INPLClassifier
		if err := s.db.Create(&inplClassifier).Error; err != nil {
			return fmt.Errorf("no se pudo crear INPLClassifier: %w", err)
		}
		// Asociar al artefacto
		artefact.InplClassifierID = &inplClassifier.ID
		if err := s.db.Model(artefact).Update("inpl_classifier_id", inplClassifier.ID).Error; err != nil {
			return fmt.Errorf("no se pudo asociar INPLClassifier al artefacto: %w", err)
		}
	}

	// Crear directorio de destino (usando estructura similar a INPLService)
	uploadRoot := "uploads/inpl"
	if envRoot := os.Getenv("INPL_UPLOAD_ROOT"); envRoot != "" {
		uploadRoot = envRoot
	}
	dir := filepath.Join(uploadRoot, strconv.Itoa(inplClassifier.ID))
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("no se pudo crear directorio: %w", err)
	}

	// Generar nombre único usando el nombre original del archivo
	filename := fmt.Sprintf("ficha_%d_%d_%s", inplClassifier.ID, time.Now().Unix(), originalFilename)
	destPath := filepath.Join(dir, filename)

	// Copiar archivo
	destFile, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("no se pudo crear archivo destino: %w", err)
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, sourceFile); err != nil {
		os.Remove(destPath)
		return fmt.Errorf("no se pudo copiar archivo: %w", err)
	}

	// Crear registro en BD
	ficha := models.INPLFicha{
		INPLClassifierID: inplClassifier.ID,
		Filename:         filename,
		OriginalName:     originalFilename,
		FilePath:         destPath,
		ContentType:      contentType,
		Size:             fileInfo.Size(),
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	if err := s.db.Create(&ficha).Error; err != nil {
		os.Remove(destPath)
		return fmt.Errorf("no se pudo guardar ficha INPL: %w", err)
	}

	// Si descargamos un archivo temporal, marcarlo para no eliminarlo (ya se copió)
	shouldDeleteTemp = false
	return nil
}
