package models

import "time"

type ArtefactModel struct {
	ID                   int                      `json:"id" gorm:"primaryKey;autoIncrement"`
	Name                 string                   `json:"name" gorm:"type:varchar(100);not null"`
	Material             *string                  `json:"material" gorm:"type:varchar(100)"`
	Observation          *string                  `json:"observation" gorm:"type:text"`
	Available            bool                     `json:"available" gorm:"type:boolean;default:true;not null"`
	Picture              []PictureModel           `json:"picture,omitempty" gorm:"foreignKey:ArtefactID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	HistoricalRecord     []HistoricalRecordModel  `json:"historicalRecord,omitempty" gorm:"foreignKey:ArtefactID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Description          *string                  `json:"description" gorm:"type:text"`
	CollectionID         *int                     `json:"collectionId" gorm:"column:collection_id"`
	Collection           *CollectionModel         `json:"collection" gorm:"foreignKey:CollectionID;references:Id"`
	ArchaeologistID      *int                     `json:"archaeologistId" gorm:"column:archaeologist_id"`
	Archaeologist        *ArchaeologistModel      `json:"archaeologist" gorm:"foreignKey:ArchaeologistID;references:Id"`
	ArchaeologicalSiteId *int                     `json:"archaeologicalSiteId" gorm:"column:archaeological_site_id"`
	ArchaeologicalSite   *ArchaeologicalSiteModel `json:"archaeologicalSite" gorm:"foreignKey:ArchaeologicalSiteId;references:Id"`
	InplClassifierID     *int                     `json:"inplClassifierId" gorm:"column:inpl_classifier_id"`
	InplClassifier       *INPLClassifierModel     `json:"inplClassifier" gorm:"foreignKey:InplClassifierID;references:ID"`
	InternalClassifierID *int                     `json:"internalClassifierId" gorm:"column:internal_classifier_id"`
	InternalClassifier   *InternalClassifierModel `json:"internalClassifier" gorm:"foreignKey:InternalClassifierID;references:Id"`
	PhysicalLocationID   *int                     `json:"physicalLocationId" gorm:"column:physical_location_id"`
	PhysicalLocation     *PhysicalLocationModel   `json:"physicalLocation" gorm:"foreignKey:PhysicalLocationID;references:ID"`
}

type PictureModel struct {
	ID           int       `json:"id" gorm:"primaryKey;autoIncrement"`
	ArtefactID   int       `json:"artefactId" gorm:"not null;uniqueIndex"`
	Filename     string    `json:"filename" gorm:"type:varchar(255);not null"`
	OriginalName string    `json:"originalName" gorm:"type:varchar(255)"`
	FilePath     string    `json:"filePath" gorm:"type:varchar(500);not null"`
	ContentType  string    `json:"contentType" gorm:"type:varchar(50)"`
	Size         int64     `json:"size"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type HistoricalRecordModel struct {
	ID           int       `json:"id" gorm:"primaryKey;autoIncrement"`
	ArtefactID   int       `json:"artefactId" gorm:"not null;uniqueIndex"`
	Filename     string    `json:"filename" gorm:"type:varchar(255);not null"`
	OriginalName string    `json:"originalName" gorm:"type:varchar(255)"`
	FilePath     string    `json:"filePath" gorm:"type:varchar(500);not null"`
	ContentType  string    `json:"contentType" gorm:"type:varchar(50)"`
	Size         int64     `json:"size"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}
