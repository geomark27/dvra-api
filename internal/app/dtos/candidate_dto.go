package dtos

import (
	"dvra-api/internal/app/models"
	"time"
)

// CandidateResponseDTO represents the candidate data in API responses
type CandidateResponseDTO struct {
	ID          uint      `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	CompanyID   uint      `json:"company_id"`
	Email       string    `json:"email"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Phone       string    `json:"phone,omitempty"`
	ResumeURL   string    `json:"resume_url,omitempty"`
	GithubURL   string    `json:"github_url,omitempty"`
	LinkedinURL string    `json:"linkedin_url,omitempty"`
	Source      string    `json:"source,omitempty"`
}

// CreateCandidateDTO represents the data needed to create a candidate
type CreateCandidateDTO struct {
	CompanyID   uint   `json:"company_id" validate:"required,min=1"`
	Email       string `json:"email" validate:"required,email"`
	FirstName   string `json:"first_name" validate:"required,min=2,max=100"`
	LastName    string `json:"last_name" validate:"required,min=2,max=100"`
	Phone       string `json:"phone,omitempty" validate:"omitempty,max=50"`
	ResumeURL   string `json:"resume_url,omitempty" validate:"omitempty,url"`
	GithubURL   string `json:"github_url,omitempty" validate:"omitempty,url"`
	LinkedinURL string `json:"linkedin_url,omitempty" validate:"omitempty,url"`
	Source      string `json:"source,omitempty" validate:"omitempty,oneof=linkedin referral direct_apply agency"`
}

// UpdateCandidateDTO represents the data needed to update a candidate
type UpdateCandidateDTO struct {
	Email       *string `json:"email,omitempty" validate:"omitempty,email"`
	FirstName   *string `json:"first_name,omitempty" validate:"omitempty,min=2,max=100"`
	LastName    *string `json:"last_name,omitempty" validate:"omitempty,min=2,max=100"`
	Phone       *string `json:"phone,omitempty" validate:"omitempty,max=50"`
	ResumeURL   *string `json:"resume_url,omitempty" validate:"omitempty,url"`
	GithubURL   *string `json:"github_url,omitempty" validate:"omitempty,url"`
	LinkedinURL *string `json:"linkedin_url,omitempty" validate:"omitempty,url"`
	Source      *string `json:"source,omitempty" validate:"omitempty,oneof=linkedin referral direct_apply agency"`
}

// ToCandidateResponse converts a Candidate model to CandidateResponseDTO
func ToCandidateResponse(candidate *models.Candidate) CandidateResponseDTO {
	return CandidateResponseDTO{
		ID:          candidate.ID,
		CreatedAt:   candidate.CreatedAt,
		UpdatedAt:   candidate.UpdatedAt,
		CompanyID:   candidate.CompanyID,
		Email:       candidate.Email,
		FirstName:   candidate.FirstName,
		LastName:    candidate.LastName,
		Phone:       candidate.Phone,
		ResumeURL:   candidate.ResumeURL,
		GithubURL:   candidate.GithubURL,
		LinkedinURL: candidate.LinkedinURL,
		Source:      candidate.Source,
	}
}

// ToCandidateResponseList converts a slice of Candidate models to CandidateResponseDTO slice
func ToCandidateResponseList(candidates []models.Candidate) []CandidateResponseDTO {
	result := make([]CandidateResponseDTO, len(candidates))
	for i, candidate := range candidates {
		result[i] = ToCandidateResponse(&candidate)
	}
	return result
}
