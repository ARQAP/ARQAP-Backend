package models

type CollectionModel struct {
	Id          int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string `json:"name" gorm:"column:name;type:varchar(255);not null"`
	Description string `json:"description" gorm:"column:description;type:text"`
	Year        int    `json:"year" gorm:"column:year;type:int"`
}
