package models

type RegionModel struct {
	ID        int          `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string       `json:"name" gorm:"type:varchar(255);not null"`
	CountryID int          `json:"countryId" gorm:"column:countryid;type:uuid;not null"`
	Country   CountryModel `json:"country" gorm:"foreignKey:CountryId;references:Id"`
}
