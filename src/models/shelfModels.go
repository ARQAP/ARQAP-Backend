package models

type ShelfModel struct {
	ID       int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Code     int    `json:"code" gorm:"column:code;type:int;not null"`
}
