package models

type TestModel struct {
	Id  int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Msg string `json:"msg" gorm:"type:text;not null"`
}
