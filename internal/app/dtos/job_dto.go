package dtos

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
