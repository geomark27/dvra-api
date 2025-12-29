package dtos

import (
	"dvra-api/internal/app/models"
	"time"
)

// ApplicationResponseDTO represents the application data in API responses
type ApplicationResponseDTO struct {
	ID          uint       `json:"id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	JobID       uint       `json:"job_id"`
	CandidateID uint       `json:"candidate_id"`
	CompanyID   uint       `json:"company_id"`
	Stage       string     `json:"stage"`
	Rating      *int       `json:"rating,omitempty"`
	RejectedAt  *time.Time `json:"rejected_at,omitempty"`
	HiredAt     *time.Time `json:"hired_at,omitempty"`
}

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

// ToApplicationResponse converts an Application model to ApplicationResponseDTO
func ToApplicationResponse(app *models.Application) ApplicationResponseDTO {
	return ApplicationResponseDTO{
		ID:          app.ID,
		CreatedAt:   app.CreatedAt,
		UpdatedAt:   app.UpdatedAt,
		JobID:       app.JobID,
		CandidateID: app.CandidateID,
		CompanyID:   app.CompanyID,
		Stage:       app.Stage,
		Rating:      app.Rating,
		RejectedAt:  app.RejectedAt,
		HiredAt:     app.HiredAt,
	}
}

// ToApplicationResponseList converts a slice of Application models to ApplicationResponseDTO slice
func ToApplicationResponseList(apps []models.Application) []ApplicationResponseDTO {
	result := make([]ApplicationResponseDTO, len(apps))
	for i, app := range apps {
		result[i] = ToApplicationResponse(&app)
	}
	return result
}
