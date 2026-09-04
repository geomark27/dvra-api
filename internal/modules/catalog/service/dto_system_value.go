package service

type CreateSystemValueDTO struct {
	CategoryID   uint    `json:"category_id" binding:"required"`
	Value        string  `json:"value" binding:"required,min=1,max=100"`
	Label        string  `json:"label" binding:"required,min=1,max=200"`
	Description  *string `json:"description,omitempty"`
	DisplayOrder int     `json:"display_order"`
}

type UpdateSystemValueDTO struct {
	Label        *string `json:"label,omitempty" binding:"omitempty,min=1,max=200"`
	Description  *string `json:"description,omitempty"`
	DisplayOrder *int    `json:"display_order,omitempty"`
	IsActive     *bool   `json:"is_active,omitempty"`
}

type UpsertSystemValueOverrideDTO struct {
	CustomLabel  *string `json:"custom_label,omitempty" binding:"omitempty,min=1,max=200"`
	DisplayOrder *int    `json:"display_order,omitempty"`
	IsEnabled    *bool   `json:"is_enabled,omitempty"`
}
