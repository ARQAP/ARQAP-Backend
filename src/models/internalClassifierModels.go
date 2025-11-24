package models

type InternalClassifierModel struct {
	Id     int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Number *int   `json:"number" gorm:"column:number;type:int;uniqueIndex:idx_internal_classifier_name_number"`
	Name   string `json:"name" gorm:"column:name;type:varchar(255);not null;uniqueIndex:idx_internal_classifier_name_number"`
}
