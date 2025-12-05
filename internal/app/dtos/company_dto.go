package dtos

import "time"

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
