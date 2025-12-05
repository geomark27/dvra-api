package dtos

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
