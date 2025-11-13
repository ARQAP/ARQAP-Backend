package models

import "time"

type LoanModel struct {
	Id          int             `json:"id" gorm:"primaryKey;autoIncrement"`
	LoanDate    time.Time       `json:"loanDate" gorm:"type:date;not null"`
	LoanTime    time.Time       `json:"loanTime" gorm:"type:time;not null"`
	ReturnDate  *time.Time      `json:"returnDate" gorm:"type:date"`
	ReturnTime  *time.Time      `json:"returnTime" gorm:"type:time"`
	ArtefactId  *int            `json:"artefactId" gorm:"column:artefact_id"`
	Artefact    *ArtefactModel  `json:"artefact" gorm:"foreignKey:ArtefactId;references:ID"`
	RequesterId *int            `json:"requesterId" gorm:"column:requester_id"`
	Requester   *RequesterModel `json:"requester" gorm:"foreignKey:RequesterId;references:Id"`
}
