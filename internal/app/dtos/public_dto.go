package dtos

import (
	"dvra-api/internal/app/models"
	"time"
)

// ===================================
// PUBLIC CAREER PAGE DTOs
// ===================================

// PublicCompanyResponseDTO represents public company info for career page
type PublicCompanyResponseDTO struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Slug     string `json:"slug"`
	LogoURL  string `json:"logo_url,omitempty"`
	Timezone string `json:"timezone"`
}

// PublicJobResponseDTO represents a job listing for public career page
type PublicJobResponseDTO struct {
	ID           uint                   `json:"id"`
	Title        string                 `json:"title"`
	Description  string                 `json:"description"`
	Requirements string                 `json:"requirements,omitempty"`
	Benefits     string                 `json:"benefits,omitempty"`
	SalaryMin    *float64               `json:"salary_min,omitempty"`
	SalaryMax    *float64               `json:"salary_max,omitempty"`
	LocationType string                 `json:"location_type"`
	City         *PublicCityDTO         `json:"city,omitempty"`
	CreatedAt    time.Time              `json:"created_at"`
	Company      *PublicCompanyShortDTO `json:"company,omitempty"`
}

// PublicCityDTO for public responses (simplified version)
type PublicCityDTO struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// PublicCompanyShortDTO for nested responses in jobs
type PublicCompanyShortDTO struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	Slug    string `json:"slug"`
	LogoURL string `json:"logo_url,omitempty"`
}

// PublicApplicationDTO represents the data needed to apply to a job publicly
type PublicApplicationDTO struct {
	FirstName   string `json:"first_name" binding:"required,min=2,max=100"`
	LastName    string `json:"last_name" binding:"required,min=2,max=100"`
	Email       string `json:"email" binding:"required,email"`
	Phone       string `json:"phone,omitempty"`
	LinkedinURL string `json:"linkedin_url,omitempty"`
	GithubURL   string `json:"github_url,omitempty"`
	CoverLetter string `json:"cover_letter,omitempty"`
	// ResumeURL will be set after file upload
	ResumeURL string `json:"-"`
}

// PublicApplicationResponseDTO represents the response after applying
type PublicApplicationResponseDTO struct {
	ID          uint      `json:"id"`
	JobID       uint      `json:"job_id"`
	JobTitle    string    `json:"job_title"`
	CompanyName string    `json:"company_name"`
	AppliedAt   time.Time `json:"applied_at"`
	Message     string    `json:"message"`
}

// ===================================
// CONVERTERS
// ===================================

// ToPublicCompanyResponse converts Company model to public DTO
func ToPublicCompanyResponse(company *models.Company) PublicCompanyResponseDTO {
	return PublicCompanyResponseDTO{
		ID:       company.ID,
		Name:     company.Name,
		Slug:     company.Slug,
		LogoURL:  company.LogoURL,
		Timezone: company.Timezone,
	}
}

// ToPublicJobResponse converts Job model to public DTO
func ToPublicJobResponse(job *models.Job) PublicJobResponseDTO {
	dto := PublicJobResponseDTO{
		ID:           job.ID,
		Title:        job.Title,
		Description:  job.Description,
		Requirements: job.Requirements,
		Benefits:     job.Benefits,
		SalaryMin:    job.SalaryMin,
		SalaryMax:    job.SalaryMax,
		LocationType: job.LocationType,
		CreatedAt:    job.CreatedAt,
	}

	// Include city if present
	if job.City != nil {
		dto.City = &PublicCityDTO{
			ID:   job.City.ID,
			Name: job.City.Name,
		}
	}

	// Include company if present
	if job.Company != nil {
		dto.Company = &PublicCompanyShortDTO{
			ID:      job.Company.ID,
			Name:    job.Company.Name,
			Slug:    job.Company.Slug,
			LogoURL: job.Company.LogoURL,
		}
	}

	return dto
}

// ToPublicJobResponseList converts a slice of Job models to public DTOs
func ToPublicJobResponseList(jobs []models.Job) []PublicJobResponseDTO {
	result := make([]PublicJobResponseDTO, len(jobs))
	for i, job := range jobs {
		result[i] = ToPublicJobResponse(&job)
	}
	return result
}
