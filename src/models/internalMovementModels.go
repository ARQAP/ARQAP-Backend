package models

import "time"

type InternalMovementModel struct {
	Id                     int                    `json:"id" gorm:"primaryKey;autoIncrement"`
	MovementDate           time.Time              `json:"movementDate" gorm:"type:date;not null"`
	MovementTime           time.Time              `json:"movementTime" gorm:"type:time;not null"`
	ReturnDate             *time.Time             `json:"returnDate" gorm:"type:date"`
	ReturnTime             *time.Time             `json:"returnTime" gorm:"type:time"`
	ArtefactId             int                    `json:"artefactId" gorm:"column:artefact_id;not null"`
	Artefact               *ArtefactModel         `json:"artefact" gorm:"foreignKey:ArtefactId;references:ID"`
	FromPhysicalLocationId *int                   `json:"fromPhysicalLocationId" gorm:"column:from_physical_location_id"`
	FromPhysicalLocation   *PhysicalLocationModel `json:"fromPhysicalLocation" gorm:"foreignKey:FromPhysicalLocationId;references:ID"`
	ToPhysicalLocationId   *int                   `json:"toPhysicalLocationId" gorm:"column:to_physical_location_id"`
	ToPhysicalLocation     *PhysicalLocationModel `json:"toPhysicalLocation" gorm:"foreignKey:ToPhysicalLocationId;references:ID"`
	Reason                 *string                `json:"reason" gorm:"type:text"`
	Observations           *string                `json:"observations" gorm:"type:text"`
	RequesterId            *int                   `json:"requesterId" gorm:"column:requester_id"`
	Requester              *RequesterModel        `json:"requester" gorm:"foreignKey:RequesterId;references:Id"`
}
