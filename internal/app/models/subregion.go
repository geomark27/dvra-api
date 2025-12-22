package models

import "gorm.io/gorm"

// Subregion represents the subregion entity in the database
type Subregion struct {
	gorm.Model
	Name     string `gorm:"type:varchar(255);not null" json:"name"`
	RegionID uint   `gorm:"not null;index" json:"region_id"`
	IsActive bool   `gorm:"type:boolean;default:true" json:"is_active"`

	// Relaciones
	Region    Region    `gorm:"foreignKey:RegionID" json:"region,omitempty"`
	Countries []Country `gorm:"foreignKey:SubregionID" json:"countries,omitempty"`
}

// TableName overrides the table name (optional)
func (Subregion) TableName() string {
	return "subregions"
}
