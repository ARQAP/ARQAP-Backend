package models

type RegionModel struct {
	ID        int          `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string       `json:"name" gorm:"column:name;type:varchar(255);not null"`
	CountryID int          `json:"countryId" gorm:"column:country_id;not null"`
	Country   CountryModel `json:"country" gorm:"foreignKey:CountryID;references:Id"`
}
