package services

import (
	"time"

	"github.com/ARQAP/ARQAP-Backend/src/models"
	"gorm.io/gorm"
)

type InternalMovementService struct {
	db *gorm.DB
}

// NewInternalMovementService creates a new instance of InternalMovementService
func NewInternalMovementService(db *gorm.DB) *InternalMovementService {
	return &InternalMovementService{db: db}
}

// Helper function to create a string pointer
func stringPtr(s string) *string {
	return &s
}

// GetAllInternalMovements retrieves all InternalMovement records from the database
func (s *InternalMovementService) GetAllInternalMovements() ([]models.InternalMovementModel, error) {
	var movements []models.InternalMovementModel

	result := s.db.
		Preload("Artefact").
		Preload("Artefact.InternalClassifier").
		Preload("FromPhysicalLocation").
		Preload("FromPhysicalLocation.Shelf").
		Preload("ToPhysicalLocation").
		Preload("ToPhysicalLocation.Shelf").
		Preload("Requester").
		Order("movement_date DESC, movement_time DESC").
		Find(&movements)

	return movements, result.Error
}

// GetInternalMovementByID retrieves an InternalMovement record by its ID
func (s *InternalMovementService) GetInternalMovementByID(id int) (*models.InternalMovementModel, error) {
	var movement models.InternalMovementModel

	result := s.db.
		Preload("Artefact").
		Preload("Artefact.InternalClassifier").
		Preload("FromPhysicalLocation").
		Preload("FromPhysicalLocation.Shelf").
		Preload("ToPhysicalLocation").
		Preload("ToPhysicalLocation.Shelf").
		Preload("Requester").
		First(&movement, id)

	if result.Error != nil {
		return nil, result.Error
	}
	return &movement, nil
}

// GetInternalMovementsByArtefactID retrieves all movements for a specific artefact
func (s *InternalMovementService) GetInternalMovementsByArtefactID(artefactId int) ([]models.InternalMovementModel, error) {
	var movements []models.InternalMovementModel

	result := s.db.
		Where("artefact_id = ?", artefactId).
		Preload("Artefact").
		Preload("Artefact.InternalClassifier").
		Preload("FromPhysicalLocation").
		Preload("FromPhysicalLocation.Shelf").
		Preload("ToPhysicalLocation").
		Preload("ToPhysicalLocation.Shelf").
		Preload("Requester").
		Order("movement_date DESC, movement_time DESC").
		Find(&movements)

	return movements, result.Error
}

// GetActiveInternalMovementByArtefactID retrieves the active movement (not finished) for a specific artefact
func (s *InternalMovementService) GetActiveInternalMovementByArtefactID(artefactId int) (*models.InternalMovementModel, error) {
	var movement models.InternalMovementModel

	result := s.db.
		Where("artefact_id = ? AND return_date IS NULL AND return_time IS NULL", artefactId).
		Preload("Artefact").
		Preload("Artefact.InternalClassifier").
		Preload("FromPhysicalLocation").
		Preload("FromPhysicalLocation.Shelf").
		Preload("ToPhysicalLocation").
		Preload("ToPhysicalLocation.Shelf").
		Preload("Requester").
		Order("movement_date DESC, movement_time DESC").
		First(&movement)

	if result.Error != nil {
		return nil, result.Error
	}
	return &movement, nil
}

// CreateInternalMovement creates a new InternalMovement record in the database
// y actualiza la ubicación física de la pieza
// Si la pieza ya tiene un movimiento activo, lo finaliza primero y usa su ubicación destino como origen del nuevo
func (s *InternalMovementService) CreateInternalMovement(movement *models.InternalMovementModel) (*models.InternalMovementModel, error) {
	err := s.db.Transaction(func(tx *gorm.DB) error {
		// 1) Buscar movimientos activos previos de la misma pieza
		var activeMovements []models.InternalMovementModel
		if err := tx.Where("artefact_id = ? AND return_date IS NULL AND return_time IS NULL", movement.ArtefactId).
			Order("movement_date DESC, movement_time DESC").
			Find(&activeMovements).Error; err != nil {
			return err
		}

		// Si hay un movimiento activo previo, usar su ubicación destino como origen del nuevo movimiento
		if len(activeMovements) > 0 {
			// Tomar el movimiento más reciente (el primero después de ordenar DESC)
			mostRecentActive := activeMovements[0]

			// Si el nuevo movimiento no tiene FromPhysicalLocation especificado,
			// usar la ubicación destino del movimiento activo previo
			if movement.FromPhysicalLocationId == nil && mostRecentActive.ToPhysicalLocationId != nil {
				movement.FromPhysicalLocationId = mostRecentActive.ToPhysicalLocationId
			}

			// Finalizar todos los movimientos activos previos
			now := time.Now()
			for _, activeMovement := range activeMovements {
				activeMovement.ReturnDate = &now
				activeMovement.ReturnTime = &now

				if err := tx.Model(&activeMovement).Updates(map[string]interface{}{
					"return_date": activeMovement.ReturnDate,
					"return_time": activeMovement.ReturnTime,
				}).Error; err != nil {
					return err
				}
			}
			// NO devolvemos la pieza a la ubicación origen aquí porque vamos a crear un nuevo movimiento
		} else {
			// Si no hay movimiento activo previo, obtener la ubicación actual de la pieza
			var artefact models.ArtefactModel
			if err := tx.First(&artefact, movement.ArtefactId).Error; err != nil {
				return err
			}

			// Si no se especificó FromPhysicalLocation, usar la ubicación actual de la pieza
			if movement.FromPhysicalLocationId == nil && artefact.PhysicalLocationID != nil {
				fromLocationId := *artefact.PhysicalLocationID
				movement.FromPhysicalLocationId = &fromLocationId
			}
		}

		// 2) Crear el nuevo movimiento
		if err := tx.Create(movement).Error; err != nil {
			return err
		}

		// 3) Mover la pieza a la ubicación destino del nuevo movimiento
		if movement.ArtefactId != 0 {
			updateData := map[string]interface{}{}
			if movement.ToPhysicalLocationId != nil {
				updateData["physical_location_id"] = *movement.ToPhysicalLocationId
			} else {
				updateData["physical_location_id"] = nil
			}

			if err := tx.Model(&models.ArtefactModel{}).
				Where("id = ?", movement.ArtefactId).
				Updates(updateData).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Devolver el movimiento con las relaciones cargadas
	if err := s.db.
		Preload("Artefact").
		Preload("Artefact.InternalClassifier").
		Preload("FromPhysicalLocation").
		Preload("FromPhysicalLocation.Shelf").
		Preload("ToPhysicalLocation").
		Preload("ToPhysicalLocation.Shelf").
		Preload("Requester").
		First(movement, movement.Id).Error; err != nil {
		return nil, err
	}

	return movement, nil
}

// DeleteInternalMovement deletes an InternalMovement record by its ID
func (s *InternalMovementService) DeleteInternalMovement(id int) error {
	result := s.db.Delete(&models.InternalMovementModel{}, id)
	return result.Error
}

// UpdateInternalMovement updates an existing InternalMovement record
// Si se está finalizando el movimiento (returnDate/returnTime se establecen), devuelve la pieza a la ubicación origen
func (s *InternalMovementService) UpdateInternalMovement(id int, updatedMovement *models.InternalMovementModel) (*models.InternalMovementModel, error) {
	var movement models.InternalMovementModel

	err := s.db.Transaction(func(tx *gorm.DB) error {
		// 1) Obtener el movimiento actual
		if err := tx.First(&movement, id).Error; err != nil {
			return err
		}

		// 2) Verificar si se está finalizando el movimiento
		isFinishing := (movement.ReturnDate == nil && movement.ReturnTime == nil) &&
			(updatedMovement.ReturnDate != nil && updatedMovement.ReturnTime != nil)

		// 3) Asegurar que el ID se mantenga
		updatedMovement.Id = id

		// 4) Actualizar campos del movimiento
		if err := tx.Model(&movement).Updates(updatedMovement).Error; err != nil {
			return err
		}

		// 5) Si se está finalizando, crear movimiento de devolución y devolver la pieza a la ubicación origen original
		if isFinishing {
			// Buscar el primer movimiento de la cadena (el más antiguo) para obtener la ubicación origen original
			var firstMovement models.InternalMovementModel
			if err := tx.Where("artefact_id = ?", movement.ArtefactId).
				Order("movement_date ASC, movement_time ASC").
				First(&firstMovement).Error; err != nil {
				// Si no se encuentra el primer movimiento, usar la ubicación origen del movimiento actual
				if movement.FromPhysicalLocationId != nil {
					if err := tx.Model(&models.ArtefactModel{}).
						Where("id = ?", movement.ArtefactId).
						Update("physical_location_id", *movement.FromPhysicalLocationId).Error; err != nil {
						return err
					}
				}
			} else {
				// Obtener la ubicación origen original del primer movimiento
				originalLocationId := firstMovement.FromPhysicalLocationId

				// Crear un nuevo movimiento de devolución desde la ubicación actual hacia la ubicación origen original
				if originalLocationId != nil && movement.ToPhysicalLocationId != nil {
					// Solo crear el movimiento de devolución si la ubicación actual es diferente a la original
					if *movement.ToPhysicalLocationId != *originalLocationId {
						returnMovement := &models.InternalMovementModel{
							MovementDate:           *updatedMovement.ReturnDate,
							MovementTime:           *updatedMovement.ReturnTime,
							ReturnDate:             updatedMovement.ReturnDate, // Ya finalizado
							ReturnTime:             updatedMovement.ReturnTime, // Ya finalizado
							ArtefactId:             movement.ArtefactId,
							FromPhysicalLocationId: movement.ToPhysicalLocationId, // Desde donde está ahora
							ToPhysicalLocationId:   originalLocationId,            // Hacia la ubicación original
							Reason:                 stringPtr("Devolución a ubicación original"),
							Observations:           stringPtr("Movimiento automático de devolución"),
							RequesterId:            movement.RequesterId, // Mantener el mismo requester si existe
						}

						if err := tx.Create(returnMovement).Error; err != nil {
							return err
						}
					}
				}

				// Actualizar la ubicación de la pieza a la ubicación origen original
				if originalLocationId != nil {
					if err := tx.Model(&models.ArtefactModel{}).
						Where("id = ?", movement.ArtefactId).
						Update("physical_location_id", *originalLocationId).Error; err != nil {
						return err
					}
				} else {
					// Si no hay ubicación origen original, dejar la pieza sin ubicación
					if err := tx.Model(&models.ArtefactModel{}).
						Where("id = ?", movement.ArtefactId).
						Update("physical_location_id", nil).Error; err != nil {
						return err
					}
				}
			}
		} else if !isFinishing {
			// Si no se está finalizando pero cambió la ubicación destino, actualizar la pieza
			if updatedMovement.ToPhysicalLocationId != nil &&
				(movement.ToPhysicalLocationId == nil || *movement.ToPhysicalLocationId != *updatedMovement.ToPhysicalLocationId) {
				updateData := map[string]interface{}{}
				if updatedMovement.ToPhysicalLocationId != nil {
					updateData["physical_location_id"] = *updatedMovement.ToPhysicalLocationId
				} else {
					updateData["physical_location_id"] = nil
				}

				if err := tx.Model(&models.ArtefactModel{}).
					Where("id = ?", movement.ArtefactId).
					Updates(updateData).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Devolver el movimiento actualizado con preload
	if err := s.db.
		Preload("Artefact").
		Preload("Artefact.InternalClassifier").
		Preload("FromPhysicalLocation").
		Preload("FromPhysicalLocation.Shelf").
		Preload("ToPhysicalLocation").
		Preload("ToPhysicalLocation.Shelf").
		Preload("Requester").
		First(&movement, id).Error; err != nil {
		return nil, err
	}

	return &movement, nil
}
