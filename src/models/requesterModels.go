package models

type RequesterType string

const (
	Investigator RequesterType = "Investigador"
	Department   RequesterType = "Departamento"
	Exhibition   RequesterType = "Exhibición"
)

type RequesterModel struct {
	Id          int           `json:"id" gorm:"primaryKey;autoIncrement"`
	Type        RequesterType `json:"type" gorm:"column:type;type:varchar(50);not null"`
	FirstName   *string        `json:"firstname" gorm:"column:firstname;type:varchar(50)"`
	LastName    *string        `json:"lastname" gorm:"column:lastname;type:varchar(50)"`
	Dni         *string        `json:"dni" gorm:"column:dni;type:varchar(20);unique"`
	Email       *string        `json:"email" gorm:"column:email;type:varchar(100)"`
	PhoneNumber *string        `json:"phoneNumber" gorm:"column:phone_number;type:varchar(20)"`
}
