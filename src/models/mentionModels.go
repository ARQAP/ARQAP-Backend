package models

type MentionModel struct {
	Id          int            `json:"id" gorm:"primaryKey;autoIncrement"`
	Title       string         `json:"title" gorm:"column:title;type:varchar(100);not null"`
	Link        string         `json:"link" gorm:"column:link;type:varchar(255);not null"`
	Description *string        `json:"description,omitempty" gorm:"column:description;type:text"`
	ArtefactId  *int           `json:"artefactId" gorm:"column:artefact_id"`
	Artefact    *ArtefactModel `json:"artefact" gorm:"foreignKey:ArtefactId;references:ID"`
}
