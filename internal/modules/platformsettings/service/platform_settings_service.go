package service

import (
	"dvra-api/internal/shared/apperr"
	"dvra-api/internal/shared/models"
)

// PlatformSettingsService maneja la lógica de negocio de platform settings
type PlatformSettingsService struct {
	repo PlatformSettingsRepository
}

// NewPlatformSettingsService crea una nueva instancia del servicio
func NewPlatformSettingsService(repo PlatformSettingsRepository) *PlatformSettingsService {
	return &PlatformSettingsService{repo: repo}
}

// GetPublic obtiene la configuración pública (para frontend)
func (s *PlatformSettingsService) GetPublic() (*PlatformSettingsPublicDTO, error) {
	settings, err := s.repo.GetCached()
	if err != nil {
		return nil, err
	}

	dto := ToPlatformSettingsPublicDTO(settings)
	return &dto, nil
}

// GetFull obtiene la configuración completa (solo SuperAdmin)
func (s *PlatformSettingsService) GetFull() (*PlatformSettingsFullDTO, error) {
	settings, err := s.repo.Get()
	if err != nil {
		return nil, err
	}

	dto := ToPlatformSettingsFullDTO(settings)
	return &dto, nil
}

// Update actualiza la configuración de la plataforma
func (s *PlatformSettingsService) Update(updateDTO *UpdatePlatformSettingsDTO, updatedByID uint) (*PlatformSettingsFullDTO, error) {
	// Obtener settings actuales
	settings, err := s.repo.Get()
	if err != nil {
		return nil, err
	}

	// Aplicar cambios
	if err := s.applyUpdates(settings, updateDTO); err != nil {
		return nil, err
	}

	// Registrar quién hizo el cambio
	settings.UpdatedByID = &updatedByID

	// Guardar
	if err := s.repo.Update(settings); err != nil {
		return nil, err
	}

	dto := ToPlatformSettingsFullDTO(settings)
	return &dto, nil
}

// applyUpdates aplica los cambios del DTO al modelo
func (s *PlatformSettingsService) applyUpdates(settings *models.PlatformSettings, dto *UpdatePlatformSettingsDTO) error {
	// Branding
	if dto.PlatformName != nil {
		if *dto.PlatformName == "" {
			return apperr.BadRequest("platform_name cannot be empty")
		}
		settings.PlatformName = *dto.PlatformName
	}
	if dto.PlatformShort != nil {
		if *dto.PlatformShort == "" {
			return apperr.BadRequest("platform_short cannot be empty")
		}
		settings.PlatformShort = *dto.PlatformShort
	}
	if dto.Tagline != nil {
		settings.Tagline = *dto.Tagline
	}
	if dto.LogoURL != nil {
		settings.LogoURL = *dto.LogoURL
	}
	if dto.LogoDarkURL != nil {
		settings.LogoDarkURL = *dto.LogoDarkURL
	}
	if dto.FaviconURL != nil {
		settings.FaviconURL = *dto.FaviconURL
	}
	if dto.PrimaryColor != nil {
		settings.PrimaryColor = *dto.PrimaryColor
	}

	// Contact
	if dto.SupportEmail != nil {
		if *dto.SupportEmail == "" {
			return apperr.BadRequest("support_email cannot be empty")
		}
		settings.SupportEmail = *dto.SupportEmail
	}
	if dto.SalesEmail != nil {
		settings.SalesEmail = *dto.SalesEmail
	}

	// URLs
	if dto.MarketingURL != nil {
		settings.MarketingURL = *dto.MarketingURL
	}
	if dto.DocsURL != nil {
		settings.DocsURL = *dto.DocsURL
	}
	if dto.TermsURL != nil {
		settings.TermsURL = *dto.TermsURL
	}
	if dto.PrivacyURL != nil {
		settings.PrivacyURL = *dto.PrivacyURL
	}

	// Business defaults
	if dto.DefaultTrialDays != nil {
		if *dto.DefaultTrialDays < 0 {
			return apperr.BadRequest("default_trial_days cannot be negative")
		}
		settings.DefaultTrialDays = *dto.DefaultTrialDays
	}
	if dto.DefaultPlanID != nil {
		settings.DefaultPlanID = dto.DefaultPlanID
	}

	// Legal
	if dto.LegalCompanyName != nil {
		settings.LegalCompanyName = *dto.LegalCompanyName
	}
	if dto.LegalAddress != nil {
		settings.LegalAddress = *dto.LegalAddress
	}
	if dto.LegalCountry != nil {
		settings.LegalCountry = *dto.LegalCountry
	}
	if dto.LegalTaxID != nil {
		settings.LegalTaxID = *dto.LegalTaxID
	}

	// Social
	if dto.TwitterURL != nil {
		settings.TwitterURL = *dto.TwitterURL
	}
	if dto.LinkedinURL != nil {
		settings.LinkedinURL = *dto.LinkedinURL
	}
	if dto.GithubURL != nil {
		settings.GithubURL = *dto.GithubURL
	}

	return nil
}

// GetDefaultTrialDays obtiene los días de prueba por defecto
func (s *PlatformSettingsService) GetDefaultTrialDays() (int, error) {
	settings, err := s.repo.GetCached()
	if err != nil {
		return 14, err // Default fallback
	}
	return settings.DefaultTrialDays, nil
}

// GetDefaultPlanTier obtiene el tier de plan por defecto
func (s *PlatformSettingsService) GetDefaultPlanTier() (string, error) {
	settings, err := s.repo.GetCached()
	if err != nil {
		return "basic", err // Default fallback
	}
	if settings.DefaultPlan == nil {
		return "", nil
	}
	return settings.DefaultPlan.Slug, nil
}
