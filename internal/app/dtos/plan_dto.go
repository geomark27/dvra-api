package dtos

// PlanDTO represents plan data transfer object
type PlanDTO struct {
	Name                string  `json:"name" binding:"required"`
	Slug                string  `json:"slug" binding:"required"`
	Description         string  `json:"description"`
	Price               float64 `json:"price" binding:"required,min=0"`
	CurrencyValueID     uint    `json:"currency_value_id" binding:"required"`
	BillingCycleValueID uint    `json:"billing_cycle_value_id" binding:"required"`
	IsActive            bool    `json:"is_active"`
	IsPublic            bool    `json:"is_public"`
	TrialDays           int     `json:"trial_days" binding:"min=0"`
	DisplayOrder        int     `json:"display_order"`
	MaxUsers            int     `json:"max_users"`
	MaxJobs             int     `json:"max_jobs"`
	MaxCandidates       int     `json:"max_candidates"`
	MaxApplications     int     `json:"max_applications"`
	MaxStorageGB        int     `json:"max_storage_gb"`
	CanExportData       bool    `json:"can_export_data"`
	CanUseCustomBrand   bool    `json:"can_use_custom_brand"`
	CanUseAPI           bool    `json:"can_use_api"`
	CanUseIntegrations  bool    `json:"can_use_integrations"`
	CanUseStaffing      bool    `json:"can_use_staffing"`
	SupportLevelValueID uint    `json:"support_level_value_id" binding:"required"`
}

// UpdatePlanDTO represents partial plan update
type UpdatePlanDTO struct {
	Name                *string  `json:"name"`
	Description         *string  `json:"description"`
	Price               *float64 `json:"price" binding:"omitempty,min=0"`
	CurrencyValueID     *uint    `json:"currency_value_id"`
	BillingCycleValueID *uint    `json:"billing_cycle_value_id"`
	IsActive            *bool    `json:"is_active"`
	IsPublic            *bool    `json:"is_public"`
	TrialDays           *int     `json:"trial_days" binding:"omitempty,min=0"`
	DisplayOrder        *int     `json:"display_order"`
	MaxUsers            *int     `json:"max_users"`
	MaxJobs             *int     `json:"max_jobs"`
	MaxCandidates       *int     `json:"max_candidates"`
	MaxApplications     *int     `json:"max_applications"`
	MaxStorageGB        *int     `json:"max_storage_gb"`
	CanExportData       *bool    `json:"can_export_data"`
	CanUseCustomBrand   *bool    `json:"can_use_custom_brand"`
	CanUseAPI           *bool    `json:"can_use_api"`
	CanUseIntegrations  *bool    `json:"can_use_integrations"`
	CanUseStaffing      *bool    `json:"can_use_staffing"`
	SupportLevelValueID *uint    `json:"support_level_value_id"`
}

// PlanResponse represents plan response
type PlanResponse struct {
	ID                  uint    `json:"id"`
	Name                string  `json:"name"`
	Slug                string  `json:"slug"`
	Description         string  `json:"description"`
	Price               float64 `json:"price"`
	Currency            string  `json:"currency"`
	CurrencyValueID     uint    `json:"currency_value_id"`
	BillingCycle        string  `json:"billing_cycle"`
	BillingCycleValueID uint    `json:"billing_cycle_value_id"`
	IsActive            bool    `json:"is_active"`
	IsPublic            bool    `json:"is_public"`
	TrialDays           int     `json:"trial_days"`
	DisplayOrder        int     `json:"display_order"`
	MaxUsers            int     `json:"max_users"`
	MaxJobs             int     `json:"max_jobs"`
	MaxCandidates       int     `json:"max_candidates"`
	MaxApplications     int     `json:"max_applications"`
	MaxStorageGB        int     `json:"max_storage_gb"`
	CanExportData       bool    `json:"can_export_data"`
	CanUseCustomBrand   bool    `json:"can_use_custom_brand"`
	CanUseAPI           bool    `json:"can_use_api"`
	CanUseIntegrations  bool    `json:"can_use_integrations"`
	CanUseStaffing      bool    `json:"can_use_staffing"`
	SupportLevel        string  `json:"support_level"`
	SupportLevelValueID uint    `json:"support_level_value_id"`
	CreatedAt           string  `json:"created_at"`
	UpdatedAt           string  `json:"updated_at"`
}

// AssignPlanToCompanyDTO represents plan assignment to company
type AssignPlanToCompanyDTO struct {
	CompanyID uint `json:"company_id" binding:"required"`
	PlanID    uint `json:"plan_id" binding:"required"`
}

// TogglePlanStatusDTO represents plan status toggle
type TogglePlanStatusDTO struct {
	IsActive bool `json:"is_active"`
}
