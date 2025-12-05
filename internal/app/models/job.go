package models

import (
	"time"

	"gorm.io/gorm"
)

// Job represents the job entity in the database
type Job struct {
	gorm.Model
	CompanyID string `gorm:"type:uuid;not null;index:idx_jobs_company_status,priority:1" json:"company_id"`

	// Datos del job
	Title       string `gorm:"type:varchar(255);not null" json:"title"`
	Description string `gorm:"type:text" json:"description"`
	Location    string `gorm:"type:varchar(255)" json:"location"`

	// Estado
	Status string `gorm:"type:varchar(50);default:'draft';index:idx_jobs_company_status,priority:2" json:"status"`
	// Valores: "draft", "active", "on_hold", "closed"

	// Asignación
	AssignedRecruiter *string `gorm:"type:uuid" json:"assigned_recruiter,omitempty"`
	HiringManager     *string `gorm:"type:uuid" json:"hiring_manager,omitempty"`

	// Timestamps
	CreatedAt time.Time      `gorm:"type:timestamp;default:now()" json:"created_at"`
	UpdatedAt time.Time      `gorm:"type:timestamp;default:now()" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relaciones
	Company      *Company      `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Applications []Application `gorm:"foreignKey:JobID" json:"applications,omitempty"`
}

// TableName overrides the table name (optional)
func (Job) TableName() string {
	return "jobs"
}
