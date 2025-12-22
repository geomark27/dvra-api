package models

import "gorm.io/gorm"

// Country represents the country entity in the database
type Country struct {
	gorm.Model
	Name        string `gorm:"type:varchar(100);not null" json:"name"`
	NumericCode string `gorm:"type:varchar(3);not null;" json:"numeric_code"`
	Iso2        string `gorm:"type:varchar(2);not null;uniqueIndex" json:"iso2"`
	Iso3        string `gorm:"type:varchar(3);not null;uniqueIndex" json:"iso3"`
	PhoneCode   string `gorm:"type:varchar(10);not null;" json:"phone_code"`
	Timezones   string `gorm:"type:text;" json:"timezones"`
	SubregionID *uint  `gorm:"index" json:"subregion_id"`
	IsActive    bool   `gorm:"type:boolean;default:true" json:"is_active"`

	// Relaciones
	Subregion *Subregion `gorm:"foreignKey:SubregionID" json:"subregion,omitempty"`
	States    []State    `gorm:"foreignKey:CountryID" json:"states,omitempty"`
}

// TableName overrides the table name (optional)
func (Country) TableName() string {
	return "countries"
}
