package models

type InternalClassifierModel struct {
	Id     int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Number int    `json:"number" gorm:"column:number;type:int;not null"`
	Color  string `json:"color" gorm:"column:color;type:varchar(255);not null"`
}
