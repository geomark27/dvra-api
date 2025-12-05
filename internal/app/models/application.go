package models

import (
	"time"

	"gorm.io/gorm"
)

// Application represents the application entity in the database
type Application struct {
	gorm.Model
	JobID       uint `gorm:"not null;index" json:"job_id"`
	CandidateID uint `gorm:"not null;index" json:"candidate_id"`

	// ⚠️ CRÍTICO: company_id (redundante pero útil para queries)
	CompanyID uint `gorm:"not null;index:idx_applications_company_stage,priority:1" json:"company_id"`

	// Pipeline
	Stage string `gorm:"type:varchar(100);not null;index:idx_applications_company_stage,priority:2" json:"stage"`
	// Valores: "applied", "screening", "technical", "offer", "hired", "rejected"

	// Rating
	Rating *int `gorm:"type:integer" json:"rating,omitempty"` // 1-5 estrellas

	// Timestamps de estado
	AppliedAt  time.Time  `gorm:"type:timestamp;default:now()" json:"applied_at"`
	RejectedAt *time.Time `gorm:"type:timestamp" json:"rejected_at,omitempty"`
	HiredAt    *time.Time `gorm:"type:timestamp" json:"hired_at,omitempty"`

	// Relaciones
	Job       *Job       `gorm:"foreignKey:JobID" json:"job,omitempty"`
	Candidate *Candidate `gorm:"foreignKey:CandidateID" json:"candidate,omitempty"`
	Company   *Company   `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
}

// TableName overrides the table name (optional)
func (Application) TableName() string {
	return "applications"
}
