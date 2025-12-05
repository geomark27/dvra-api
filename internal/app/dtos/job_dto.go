package dtos

// CreateJobDTO represents the data needed to create a job
type CreateJobDTO struct {
	CompanyID         uint   `json:"company_id" validate:"required,min=1"`
	Title             string `json:"title" validate:"required,min=3,max=255"`
	Description       string `json:"description" validate:"required"`
	Location          string `json:"location,omitempty" validate:"omitempty,max=255"`
	Status            string `json:"status" validate:"omitempty,oneof=draft active on_hold closed"`
	AssignedRecruiter *uint  `json:"assigned_recruiter,omitempty"`
	HiringManager     *uint  `json:"hiring_manager,omitempty"`
}

// UpdateJobDTO represents the data needed to update a job
type UpdateJobDTO struct {
	Title             *string `json:"title,omitempty" validate:"omitempty,min=3,max=255"`
	Description       *string `json:"description,omitempty"`
	Location          *string `json:"location,omitempty" validate:"omitempty,max=255"`
	Status            *string `json:"status,omitempty" validate:"omitempty,oneof=draft active on_hold closed"`
	AssignedRecruiter *uint   `json:"assigned_recruiter,omitempty"`
	HiringManager     *uint   `json:"hiring_manager,omitempty"`
}
