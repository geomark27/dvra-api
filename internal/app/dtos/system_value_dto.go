package dtos

// CreateSystemValueDTO for creating a new system value
type CreateSystemValueDTO struct {
	Category     string  `json:"category" binding:"required"`
	Value        string  `json:"value" binding:"required"`
	Label        string  `json:"label" binding:"required"`
	Description  *string `json:"description"`
	DisplayOrder int     `json:"display_order"`
	CompanyID    *uint   `json:"company_id"`
}

// UpdateSystemValueDTO for updating a system value
type UpdateSystemValueDTO struct {
	Label        string  `json:"label"`
	Description  *string `json:"description"`
	DisplayOrder int     `json:"display_order"`
	IsActive     bool    `json:"is_active"`
}
