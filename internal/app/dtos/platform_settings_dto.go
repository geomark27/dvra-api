package dtos

import "dvra-api/internal/app/models"

// ============================================================================
// PLATFORM SETTINGS DTOs
// ============================================================================

// PlatformSettingsPublicDTO - Datos públicos (para frontend sin auth)
type PlatformSettingsPublicDTO struct {
	PlatformName  string `json:"platform_name"`
	PlatformShort string `json:"platform_short"`
	Tagline       string `json:"tagline"`
	LogoURL       string `json:"logo_url,omitempty"`
	LogoDarkURL   string `json:"logo_dark_url,omitempty"`
	FaviconURL    string `json:"favicon_url,omitempty"`
	PrimaryColor  string `json:"primary_color"`
	SupportEmail  string `json:"support_email"`
	MarketingURL  string `json:"marketing_url"`
	DocsURL       string `json:"docs_url,omitempty"`
	TermsURL      string `json:"terms_url,omitempty"`
	PrivacyURL    string `json:"privacy_url,omitempty"`

	// Social media
	TwitterURL  string `json:"twitter_url,omitempty"`
	LinkedinURL string `json:"linkedin_url,omitempty"`
	GithubURL   string `json:"github_url,omitempty"`
}

// PlatformSettingsFullDTO - Todos los datos (solo SuperAdmin)
type PlatformSettingsFullDTO struct {
	ID uint `json:"id"`

	// Branding
	PlatformName  string `json:"platform_name"`
	PlatformShort string `json:"platform_short"`
	Tagline       string `json:"tagline"`
	LogoURL       string `json:"logo_url,omitempty"`
	LogoDarkURL   string `json:"logo_dark_url,omitempty"`
	FaviconURL    string `json:"favicon_url,omitempty"`
	PrimaryColor  string `json:"primary_color"`

	// Contact
	SupportEmail string `json:"support_email"`
	SalesEmail   string `json:"sales_email,omitempty"`

	// URLs
	MarketingURL string `json:"marketing_url"`
	DocsURL      string `json:"docs_url,omitempty"`
	TermsURL     string `json:"terms_url,omitempty"`
	PrivacyURL   string `json:"privacy_url,omitempty"`

	// Business defaults
	DefaultTrialDays int    `json:"default_trial_days"`
	DefaultPlanTier  string `json:"default_plan_tier"`

	// Legal
	LegalCompanyName string `json:"legal_company_name,omitempty"`
	LegalAddress     string `json:"legal_address,omitempty"`
	LegalCountry     string `json:"legal_country,omitempty"`
	LegalTaxID       string `json:"legal_tax_id,omitempty"`

	// Social
	TwitterURL  string `json:"twitter_url,omitempty"`
	LinkedinURL string `json:"linkedin_url,omitempty"`
	GithubURL   string `json:"github_url,omitempty"`

	// Metadata
	UpdatedAt   string `json:"updated_at"`
	UpdatedByID *uint  `json:"updated_by_id,omitempty"`
}

// UpdatePlatformSettingsDTO - Para actualizar settings
type UpdatePlatformSettingsDTO struct {
	// Branding
	PlatformName  *string `json:"platform_name,omitempty"`
	PlatformShort *string `json:"platform_short,omitempty"`
	Tagline       *string `json:"tagline,omitempty"`
	LogoURL       *string `json:"logo_url,omitempty"`
	LogoDarkURL   *string `json:"logo_dark_url,omitempty"`
	FaviconURL    *string `json:"favicon_url,omitempty"`
	PrimaryColor  *string `json:"primary_color,omitempty"`

	// Contact
	SupportEmail *string `json:"support_email,omitempty"`
	SalesEmail   *string `json:"sales_email,omitempty"`

	// URLs
	MarketingURL *string `json:"marketing_url,omitempty"`
	DocsURL      *string `json:"docs_url,omitempty"`
	TermsURL     *string `json:"terms_url,omitempty"`
	PrivacyURL   *string `json:"privacy_url,omitempty"`

	// Business defaults
	DefaultTrialDays *int    `json:"default_trial_days,omitempty"`
	DefaultPlanTier  *string `json:"default_plan_tier,omitempty"`

	// Legal
	LegalCompanyName *string `json:"legal_company_name,omitempty"`
	LegalAddress     *string `json:"legal_address,omitempty"`
	LegalCountry     *string `json:"legal_country,omitempty"`
	LegalTaxID       *string `json:"legal_tax_id,omitempty"`

	// Social
	TwitterURL  *string `json:"twitter_url,omitempty"`
	LinkedinURL *string `json:"linkedin_url,omitempty"`
	GithubURL   *string `json:"github_url,omitempty"`
}

// ============================================================================
// CONVERTERS
// ============================================================================

// ToPublicDTO convierte el modelo a DTO público
func ToPlatformSettingsPublicDTO(s *models.PlatformSettings) PlatformSettingsPublicDTO {
	return PlatformSettingsPublicDTO{
		PlatformName:  s.PlatformName,
		PlatformShort: s.PlatformShort,
		Tagline:       s.Tagline,
		LogoURL:       s.LogoURL,
		LogoDarkURL:   s.LogoDarkURL,
		FaviconURL:    s.FaviconURL,
		PrimaryColor:  s.PrimaryColor,
		SupportEmail:  s.SupportEmail,
		MarketingURL:  s.MarketingURL,
		DocsURL:       s.DocsURL,
		TermsURL:      s.TermsURL,
		PrivacyURL:    s.PrivacyURL,
		TwitterURL:    s.TwitterURL,
		LinkedinURL:   s.LinkedinURL,
		GithubURL:     s.GithubURL,
	}
}

// ToFullDTO convierte el modelo a DTO completo (para SuperAdmin)
func ToPlatformSettingsFullDTO(s *models.PlatformSettings) PlatformSettingsFullDTO {
	return PlatformSettingsFullDTO{
		ID:               s.ID,
		PlatformName:     s.PlatformName,
		PlatformShort:    s.PlatformShort,
		Tagline:          s.Tagline,
		LogoURL:          s.LogoURL,
		LogoDarkURL:      s.LogoDarkURL,
		FaviconURL:       s.FaviconURL,
		PrimaryColor:     s.PrimaryColor,
		SupportEmail:     s.SupportEmail,
		SalesEmail:       s.SalesEmail,
		MarketingURL:     s.MarketingURL,
		DocsURL:          s.DocsURL,
		TermsURL:         s.TermsURL,
		PrivacyURL:       s.PrivacyURL,
		DefaultTrialDays: s.DefaultTrialDays,
		DefaultPlanTier:  s.DefaultPlanTier,
		LegalCompanyName: s.LegalCompanyName,
		LegalAddress:     s.LegalAddress,
		LegalCountry:     s.LegalCountry,
		LegalTaxID:       s.LegalTaxID,
		TwitterURL:       s.TwitterURL,
		LinkedinURL:      s.LinkedinURL,
		GithubURL:        s.GithubURL,
		UpdatedAt:        s.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedByID:      s.UpdatedByID,
	}
}
