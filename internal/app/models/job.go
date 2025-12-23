package models

import (
	"gorm.io/gorm"
)

// Job represents the job entity in the database
type Job struct {
	gorm.Model
	CompanyID    uint     `gorm:"not null;index:idx_jobs_company_status,priority:1" json:"company_id"`
	Title        string   `gorm:"type:varchar(255);not null" json:"title"`
	Description  string   `gorm:"type:text" json:"description"`
	SalaryMin    *float64 `gorm:"type:decimal(12,2)" json:"salary_min,omitempty"`
	SalaryMax    *float64 `gorm:"type:decimal(12,2)" json:"salary_max,omitempty"`
	Requirements string   `gorm:"type:text" json:"requirements,omitempty"`
	Benefits     string   `gorm:"type:text" json:"benefits,omitempty"`
	Status       string   `gorm:"type:varchar(50);default:'draft';index:idx_jobs_company_status,priority:2" json:"status"`

	// LocationType uses system_values with category='work_mode' (remote, onsite, hybrid)
	LocationType string `gorm:"type:varchar(50);default:'onsite'" json:"location_type"`
	CityID       *uint  `gorm:"index" json:"city_id,omitempty"` // If onsite/hybrid, reference to cities table

	AssignedRecruiter *uint `gorm:"" json:"assigned_recruiter,omitempty"`
	HiringManager     *uint `gorm:"" json:"hiring_manager,omitempty"`

	// Relations
	Company      *Company      `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	City         *City         `gorm:"foreignKey:CityID" json:"city,omitempty"`
	Applications []Application `gorm:"foreignKey:JobID" json:"applications,omitempty"`
}

// TableName overrides the table name (optional)
func (Job) TableName() string {
	return "jobs"
}
