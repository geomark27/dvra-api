package dtos

import (
	"dvra-api/internal/app/models"
	"time"
)

// JobFilters represents filters for listing jobs
type JobFilters struct {
	Status       string `form:"status"`
	LocationType string `form:"location_type"`
	CityID       *uint  `form:"city_id"`
}

// JobResponseDTO represents the job data in API responses
type JobResponseDTO struct {
	ID                uint            `json:"id"`
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
	CompanyID         uint            `json:"company_id"`
	Company           *models.Company `json:"company,omitempty"`
	Title             string          `json:"title"`
	Description       string          `json:"description"`
	SalaryMin         *float64        `json:"salary_min,omitempty"`
	SalaryMax         *float64        `json:"salary_max,omitempty"`
	Requirements      string          `json:"requirements,omitempty"`
	Benefits          string          `json:"benefits,omitempty"`
	Status            string          `json:"status"`
	LocationType      string          `json:"location_type"`
	CityID            *uint           `json:"city_id,omitempty"`
	City              *models.City    `json:"city,omitempty"`
	AssignedRecruiter *uint           `json:"assigned_recruiter,omitempty"`
	HiringManager     *uint           `json:"hiring_manager,omitempty"`
}

// CreateJobDTO represents the data needed to create a job
type CreateJobDTO struct {
	CompanyID         uint     `json:"company_id" validate:"required,min=1"`
	Title             string   `json:"title" validate:"required,min=3,max=255"`
	Description       string   `json:"description" validate:"required"`
	SalaryMin         *float64 `json:"salary_min,omitempty"`
	SalaryMax         *float64 `json:"salary_max,omitempty"`
	Requirements      string   `json:"requirements,omitempty"`
	Benefits          string   `json:"benefits,omitempty"`
	Status            string   `json:"status" validate:"omitempty,oneof=draft active on_hold closed"`
	LocationType      string   `json:"location_type" binding:"required"` // Validar contra system_values.work_mode
	CityID            *uint    `json:"city_id,omitempty"`
	AssignedRecruiter *uint    `json:"assigned_recruiter,omitempty"`
	HiringManager     *uint    `json:"hiring_manager,omitempty"`
}

// UpdateJobDTO represents the data needed to update a job
type UpdateJobDTO struct {
	Title             *string  `json:"title,omitempty" validate:"omitempty,min=3,max=255"`
	Description       *string  `json:"description,omitempty"`
	SalaryMin         *float64 `json:"salary_min,omitempty"`
	SalaryMax         *float64 `json:"salary_max,omitempty"`
	Requirements      *string  `json:"requirements,omitempty"`
	Benefits          *string  `json:"benefits,omitempty"`
	Status            *string  `json:"status,omitempty" validate:"omitempty,oneof=draft active on_hold closed"`
	LocationType      *string  `json:"location_type,omitempty"` // Validar contra system_values.work_mode
	CityID            *uint    `json:"city_id,omitempty"`
	AssignedRecruiter *uint    `json:"assigned_recruiter,omitempty"`
	HiringManager     *uint    `json:"hiring_manager,omitempty"`
}

// ToJobResponse converts a Job model to JobResponseDTO
func ToJobResponse(job *models.Job) JobResponseDTO {
	return JobResponseDTO{
		ID:                job.ID,
		CreatedAt:         job.CreatedAt,
		UpdatedAt:         job.UpdatedAt,
		CompanyID:         job.CompanyID,
		Company:           job.Company,
		Title:             job.Title,
		Description:       job.Description,
		SalaryMin:         job.SalaryMin,
		SalaryMax:         job.SalaryMax,
		Requirements:      job.Requirements,
		Benefits:          job.Benefits,
		Status:            job.Status,
		LocationType:      job.LocationType,
		CityID:            job.CityID,
		City:              job.City,
		AssignedRecruiter: job.AssignedRecruiter,
		HiringManager:     job.HiringManager,
	}
}

// ToJobResponseList converts a slice of Job models to JobResponseDTO slice
func ToJobResponseList(jobs []models.Job) []JobResponseDTO {
	result := make([]JobResponseDTO, len(jobs))
	for i, job := range jobs {
		result[i] = ToJobResponse(&job)
	}
	return result
}
