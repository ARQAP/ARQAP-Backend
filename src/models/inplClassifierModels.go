package models

import "time"

type INPLClassifierModel struct {
	ID         int         `json:"id" gorm:"primaryKey;autoIncrement"`
	INPLFichas []INPLFicha `json:"inplFichas,omitempty" gorm:"foreignKey:INPLClassifierID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type INPLFicha struct {
	ID               int       `json:"id" gorm:"primaryKey;autoIncrement"`
	INPLClassifierID int       `json:"inplClassifierId" gorm:"column:inpl_classifier_id;index;not null"`
	Filename         string    `json:"filename" gorm:"type:varchar(255);not null"`
	OriginalName     string    `json:"originalName" gorm:"column:original_name;type:varchar(255)"`
	FilePath         string    `json:"filePath" gorm:"column:file_path;type:varchar(500);not null"`
	ContentType      string    `json:"contentType" gorm:"column:content_type;type:varchar(50)"`
	Size             int64     `json:"size"`
	CreatedAt        time.Time `json:"createdAt" gorm:"column:created_at"`
	UpdatedAt        time.Time `json:"updatedAt" gorm:"column:updated_at"`
}
