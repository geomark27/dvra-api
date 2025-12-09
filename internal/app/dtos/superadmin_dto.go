package dtos

// CreateCompanyWithAdminDTO - Crear empresa con primer admin
type CreateCompanyWithAdminDTO struct {
	CompanyName    string  `json:"company_name" validate:"required"`
	CompanySlug    string  `json:"company_slug" validate:"required"`
	PlanSlug       *string `json:"plan_slug" validate:"omitempty"` // Opcional, default "free"
	AdminEmail     string  `json:"admin_email" validate:"required,email"`
	AdminPassword  string  `json:"admin_password" validate:"required,min=8"`
	AdminFirstName string  `json:"admin_first_name" validate:"required"`
	AdminLastName  string  `json:"admin_last_name" validate:"required"`
}

// ChangePlanDTO - Cambiar plan de empresa
type ChangePlanDTO struct {
	NewPlan string `json:"new_plan" validate:"required,oneof=free professional enterprise"`
}

// SuspendCompanyDTO - Suspender empresa
type SuspendCompanyDTO struct {
	Reason string `json:"reason" validate:"required"`
}

// CompanyWithStatsDTO - Empresa con estadísticas
type CompanyWithStatsDTO struct {
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Slug        string  `json:"slug"`
	PlanTier    string  `json:"plan_tier"`
	Status      string  `json:"status"`
	UserCount   int     `json:"user_count"`
	JobCount    int     `json:"job_count"`
	CreatedAt   string  `json:"created_at"`
	TrialEndsAt *string `json:"trial_ends_at"`
}

// GlobalAnalyticsDTO - Analytics globales del sistema
type GlobalAnalyticsDTO struct {
	TotalCompanies     int     `json:"total_companies"`
	ActiveCompanies    int     `json:"active_companies"`
	SuspendedCompanies int     `json:"suspended_companies"`
	TotalUsers         int     `json:"total_users"`
	TotalJobs          int     `json:"total_jobs"`
	TotalApplications  int     `json:"total_applications"`
	MonthlyRevenue     float64 `json:"monthly_revenue"`
	ChurnRate          float64 `json:"churn_rate"`
}
