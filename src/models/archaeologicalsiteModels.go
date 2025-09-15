package models

type ArchaeologicalSiteModel struct {
	Id          int         `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string      `json:"Name" gorm:"column:Name;type:varchar(50);not null"`
	Location    string      `json:"Location" gorm:"column:Location;type:varchar(50);not null"`
	Description string      `json:"Description" gorm:"column:Description;type:varchar(255);not null"`
	RegionID    int         `json:"regionId" gorm:"column:region_id;not null"`
	Region      RegionModel `json:"region" gorm:"foreignKey:RegionID;references:ID"`
}
