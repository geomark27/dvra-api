package service

import (
	"dvra-api/internal/shared/models"
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
	CompanyID   uint   `json:"company_id"`
	Email       string `json:"email" binding:"required,email"`
	FirstName   string `json:"first_name" binding:"required,min=2,max=100"`
	LastName    string `json:"last_name" binding:"required,min=2,max=100"`
	Phone       string `json:"phone,omitempty" binding:"omitempty,max=50"`
	ResumeURL   string `json:"resume_url,omitempty" binding:"omitempty,url"`
	GithubURL   string `json:"github_url,omitempty" binding:"omitempty,url"`
	LinkedinURL string `json:"linkedin_url,omitempty" binding:"omitempty,url"`
	Source      string `json:"source,omitempty" binding:"omitempty,oneof=linkedin referral direct_apply agency career_page"`
}

// UpdateCandidateDTO represents the data needed to update a candidate
type UpdateCandidateDTO struct {
	Email       *string `json:"email,omitempty" binding:"omitempty,email"`
	FirstName   *string `json:"first_name,omitempty" binding:"omitempty,min=2,max=100"`
	LastName    *string `json:"last_name,omitempty" binding:"omitempty,min=2,max=100"`
	Phone       *string `json:"phone,omitempty" binding:"omitempty,max=50"`
	ResumeURL   *string `json:"resume_url,omitempty" binding:"omitempty,url"`
	GithubURL   *string `json:"github_url,omitempty" binding:"omitempty,url"`
	LinkedinURL *string `json:"linkedin_url,omitempty" binding:"omitempty,url"`
	Source      *string `json:"source,omitempty" binding:"omitempty,oneof=linkedin referral direct_apply agency career_page"`
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
