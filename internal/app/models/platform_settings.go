package models

// PlatformSettings representa la configuración global del SaaS
// Solo debe existir UNA fila en esta tabla (singleton)
type PlatformSettings struct {
	BaseModel
	// =========================================================================
	// BRANDING
	// =========================================================================
	PlatformName  string `gorm:"type:varchar(100);not null;default:'DVRA ATS'" json:"platform_name"`
	PlatformShort string `gorm:"type:varchar(50);not null;default:'DVRA'" json:"platform_short"`
	Tagline       string `gorm:"type:varchar(255);default:'Applicant Tracking System'" json:"tagline"`
	LogoURL       string `gorm:"type:text" json:"logo_url,omitempty"`
	LogoDarkURL   string `gorm:"type:text" json:"logo_dark_url,omitempty"`
	FaviconURL    string `gorm:"type:text" json:"favicon_url,omitempty"`
	PrimaryColor  string `gorm:"type:varchar(7);default:'#2563eb'" json:"primary_color"`

	// =========================================================================
	// CONTACT & SUPPORT
	// =========================================================================
	SupportEmail string `gorm:"type:varchar(255);default:'support@dvra.app'" json:"support_email"`
	SalesEmail   string `gorm:"type:varchar(255)" json:"sales_email,omitempty"`

	// =========================================================================
	// URLs
	// =========================================================================
	MarketingURL string `gorm:"type:varchar(255);default:'https://dvra.app'" json:"marketing_url"`
	DocsURL      string `gorm:"type:varchar(255)" json:"docs_url,omitempty"`
	TermsURL     string `gorm:"type:varchar(255)" json:"terms_url,omitempty"`
	PrivacyURL   string `gorm:"type:varchar(255)" json:"privacy_url,omitempty"`

	// =========================================================================
	// BUSINESS DEFAULTS
	// =========================================================================
	DefaultTrialDays int    `gorm:"type:int;default:14" json:"default_trial_days"`
	DefaultPlanTier  string `gorm:"type:varchar(50);default:'free'" json:"default_plan_tier"`

	// =========================================================================
	// LEGAL / COMPANY INFO
	// =========================================================================
	LegalCompanyName string `gorm:"type:varchar(255)" json:"legal_company_name,omitempty"`
	LegalAddress     string `gorm:"type:text" json:"legal_address,omitempty"`
	LegalCountry     string `gorm:"type:varchar(100)" json:"legal_country,omitempty"`
	LegalTaxID       string `gorm:"type:varchar(100)" json:"legal_tax_id,omitempty"`

	// =========================================================================
	// SOCIAL MEDIA
	// =========================================================================
	TwitterURL  string `gorm:"type:varchar(255)" json:"twitter_url,omitempty"`
	LinkedinURL string `gorm:"type:varchar(255)" json:"linkedin_url,omitempty"`
	GithubURL   string `gorm:"type:varchar(255)" json:"github_url,omitempty"`

	// =========================================================================
	// METADATA
	// =========================================================================
	UpdatedByID *uint `gorm:"type:int" json:"updated_by_id,omitempty"`
}

func (PlatformSettings) TableName() string {
	return "platform_settings"
}
