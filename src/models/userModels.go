package models

type UserModel struct {
	Id       int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Username string `json:"username" gorm:"column:username;type:varchar(255);not null"`
	Password string `json:"password" gorm:"type:varchar(100);not null"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}
