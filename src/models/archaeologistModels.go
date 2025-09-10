package models

type ArchaeologistModel struct {
	Id        int    `json:"id" gorm:"primaryKey;autoIncrement"`
	FirstName string `json:"firstname" gorm:"column:firstname;type:varchar(50);not null"`
	LastName  string `json:"lastname" gorm:"column:lastname;type:varchar(50);not null"`
}
