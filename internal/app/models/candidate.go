package models

import (
	"gorm.io/gorm"
)

// Candidate represents the candidate entity in the database
type Candidate struct {
	gorm.Model
	CompanyID uint `gorm:"not null;index:idx_candidates_company_email,priority:1" json:"company_id"`

	// Datos del candidato
	Email     string `gorm:"type:varchar(255);not null;index:idx_candidates_company_email,priority:2" json:"email"`
	FirstName string `gorm:"type:varchar(100)" json:"first_name"`
	LastName  string `gorm:"type:varchar(100)" json:"last_name"`
	Phone     string `gorm:"type:varchar(50)" json:"phone,omitempty"`

	// Archivos
	ResumeURL   string `gorm:"type:text" json:"resume_url,omitempty"`
	GithubURL   string `gorm:"type:text" json:"github_url,omitempty"`
	LinkedinURL string `gorm:"type:text" json:"linkedin_url,omitempty"`

	// Source tracking
	Source string `gorm:"type:varchar(100)" json:"source,omitempty"`
	// Valores: "linkedin", "referral", "direct_apply", "agency"

	// Relaciones
	Company      *Company      `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Applications []Application `gorm:"foreignKey:CandidateID" json:"applications,omitempty"`
}

// TableName overrides the table name (optional)
func (Candidate) TableName() string {
	return "candidates"
}
