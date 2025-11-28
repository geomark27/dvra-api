package models

import "gorm.io/gorm"

// Company represents the company entity in the database
type Company struct {
	gorm.Model
	Name        string `gorm:"size:100;not null" json:"name"`
	Description string `gorm:"type:text" json:"description"`
	// Add your fields here
	// Example:
	// Price    float64 `gorm:"not null;default:0" json:"price"`
	// Stock    int     `gorm:"default:0" json:"stock"`
	// IsActive bool    `gorm:"default:true" json:"is_active"`
}

// TableName overrides the table name (optional)
// func (Company) TableName() string {
//     return "companys"
// }
