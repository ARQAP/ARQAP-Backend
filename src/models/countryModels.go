package models

type CountryModel struct {
	Id   int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Name string `json:"name" gorm:"column:name;type:varchar(255);not null"`
}