package dtos

import "time"

// CreateApplicationDTO represents the data needed to create an application
type CreateApplicationDTO struct {
	JobID       uint   `json:"job_id" validate:"required,min=1"`
	CandidateID uint   `json:"candidate_id" validate:"required,min=1"`
	CompanyID   uint   `json:"company_id" validate:"required,min=1"`
	Stage       string `json:"stage" validate:"required,oneof=applied screening technical offer hired rejected"`
	Rating      *int   `json:"rating,omitempty" validate:"omitempty,min=1,max=5"`
}

// UpdateApplicationDTO represents the data needed to update an application
type UpdateApplicationDTO struct {
	Stage      *string    `json:"stage,omitempty" validate:"omitempty,oneof=applied screening technical offer hired rejected"`
	Rating     *int       `json:"rating,omitempty" validate:"omitempty,min=1,max=5"`
	RejectedAt *time.Time `json:"rejected_at,omitempty"`
	HiredAt    *time.Time `json:"hired_at,omitempty"`
}
