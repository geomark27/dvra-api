package models

import "gorm.io/gorm"

// State represents the state entity in the database
type State struct {
	gorm.Model
	Name        string `gorm:"type:varchar(255);not null" json:"name"`
	CountryID   uint   `gorm:"not null;index" json:"country_id"`
	CountryCode string `gorm:"type:varchar(3);not null;" json:"country_code"`
	IsActive    bool   `gorm:"type:boolean;default:true" json:"is_active"`

	// Relaciones
	Country Country `gorm:"foreignKey:CountryID" json:"country,omitempty"`
	Cities  []City  `gorm:"foreignKey:StateID" json:"cities,omitempty"`
}

// TableName overrides the table name (optional)
func (State) TableName() string {
	return "states"
}
