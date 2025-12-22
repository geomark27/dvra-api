package models

import "gorm.io/gorm"

// City represents the city entity in the database
type City struct {
	gorm.Model
	Name        string       `gorm:"type:varchar(255);not null" json:"name"`
	StateID     uint         `gorm:"not null;index" json:"state_id"`
	IsActive    bool         `gorm:"type:boolean;default:true" json:"is_active"`
	
	// Relaciones
	State State `gorm:"foreignKey:StateID" json:"state,omitempty"`
}

// TableName overrides the table name (optional)
func (City) TableName() string {
    return "cities"
}
