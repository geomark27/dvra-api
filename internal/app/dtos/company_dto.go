package dtos

import (
	"dvra-api/internal/app/models"
	"time"
)

// CompanyResponseDTO represents the company data in API responses
type CompanyResponseDTO struct {
	ID          uint       `json:"id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Name        string     `json:"name"`
	Slug        string     `json:"slug"`
	LogoURL     string     `json:"logo_url,omitempty"`
	PlanTier    string     `json:"plan_tier"`
	TrialEndsAt *time.Time `json:"trial_ends_at,omitempty"`
	Timezone    string     `json:"timezone"`
}

// CreateCompanyDTO represents the data needed to create a company
type CreateCompanyDTO struct {
	Name        string     `json:"name" validate:"required,min=2,max=255"`
	Slug        string     `json:"slug" validate:"required,min=2,max=100,alphanum"`
	LogoURL     string     `json:"logo_url,omitempty"`
	PlanTier    string     `json:"plan_tier" validate:"required,oneof=free premium enterprise"`
	TrialEndsAt *time.Time `json:"trial_ends_at,omitempty"`
	Timezone    string     `json:"timezone" validate:"required"`
}

// UpdateCompanyDTO represents the data needed to update a company
type UpdateCompanyDTO struct {
	Name        *string    `json:"name,omitempty" validate:"omitempty,min=2,max=255"`
	Slug        *string    `json:"slug,omitempty" validate:"omitempty,min=2,max=100,alphanum"`
	LogoURL     *string    `json:"logo_url,omitempty"`
	PlanTier    *string    `json:"plan_tier,omitempty" validate:"omitempty,oneof=free premium enterprise"`
	TrialEndsAt *time.Time `json:"trial_ends_at,omitempty"`
	Timezone    *string    `json:"timezone,omitempty"`
}

// ToCompanyResponse converts a Company model to CompanyResponseDTO
func ToCompanyResponse(company *models.Company) CompanyResponseDTO {
	return CompanyResponseDTO{
		ID:          company.ID,
		CreatedAt:   company.CreatedAt,
		UpdatedAt:   company.UpdatedAt,
		Name:        company.Name,
		Slug:        company.Slug,
		LogoURL:     company.LogoURL,
		PlanTier:    company.PlanTier,
		TrialEndsAt: company.TrialEndsAt,
		Timezone:    company.Timezone,
	}
}

// ToCompanyResponseList converts a slice of Company models to CompanyResponseDTO slice
func ToCompanyResponseList(companies []models.Company) []CompanyResponseDTO {
	result := make([]CompanyResponseDTO, len(companies))
	for i, company := range companies {
		result[i] = ToCompanyResponse(&company)
	}
	return result
}
