package models

type ArtefactModel struct {
	ID                      int                     `json:"id" gorm:"primaryKey;autoIncrement"`
	Name                    string                  `json:"name" gorm:"type:varchar(100);not null"`
	Material                *string                 `json:"material" gorm:"type:varchar(100)"`
	Observation             *string                 `json:"observation" gorm:"type:text"`
	Available               bool                    `json:"available" gorm:"type:boolean;default:true;not null"`
	// Picture              []PictureModel          `json:"picture,omitempty" gorm:"foreignKey:ArtefactID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`          // REVISAR PARA IMPLEMENTACIÓN DE IMÁGENES
	// HistoricalRecord     []HistoricalRecordModel `json:"historicalRecord,omitempty" gorm:"foreignKey:ArtefactID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"` // REVISAR PARA IMPLEMENTACIÓN DE IMÁGENES
	Description             *string                 `json:"description" gorm:"type:text"`
	// CollectionID         *int                    `json:"collectionId" gorm:"column:collection_id"`
	// Collection           CollectionModel         `json:"collection" gorm:"foreignKey:CollectionID;references:ID"`
	ArchaeologistID         *int                    `json:"archaeologistId" gorm:"column:archaeologist_id"`
	Archaeologist           ArchaeologistModel      `json:"archaeologist" gorm:"foreignKey:ArchaeologistID;references:ID"`
	// InplClassifierID     *int                    `json:"inplClassifierId" gorm:"column:inpl_classifier_id"`
	// InplClassifier       InplClassifierModel     `json:"inplClassifier" gorm:"foreignKey:InplClassifierID;references:ID"`
	// InternalClassifierID *int                    `json:"internalClassifierId" gorm:"column:internal_classifier_id"`
	// InternalClassifier   InternalClassifierModel `json:"internalClassifier" gorm:"foreignKey:InternalClassifierID;references:ID"`
	PhysicalLocationID      *int                    `json:"physicalLocationId" gorm:"column:physical_location_id"`
	PhysicalLocation        PhysicalLocationModel   `json:"physicalLocation" gorm:"foreignKey:PhysicalLocationID;references:ID"`
}
