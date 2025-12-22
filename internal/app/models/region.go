package models

import "gorm.io/gorm"

// Region represents the region entity in the database
type Region struct {
	gorm.Model
	Name     string `gorm:"type:varchar(255);not null" json:"name"`
	IsActive bool   `gorm:"type:boolean;default:true" json:"is_active"`

	// Relaciones
	Subregions []Subregion `gorm:"foreignKey:RegionID" json:"subregions,omitempty"`
}

// TableName overrides the table name (optional)
func (Region) TableName() string {
	return "regions"
}
