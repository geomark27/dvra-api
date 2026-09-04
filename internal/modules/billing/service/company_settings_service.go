package service

import (
	"errors"

	"dvra-api/internal/modules/billing/repository"
	"dvra-api/internal/shared/apperr"
	"dvra-api/internal/shared/models"

	"gorm.io/gorm"
)

type CompanySettingsService struct{ db *gorm.DB }

func NewCompanySettingsService(db *gorm.DB) *CompanySettingsService {
	return &CompanySettingsService{db: db}
}

func (s *CompanySettingsService) Get(companyID uint) (*models.CompanySettings, error) {
	var settings models.CompanySettings
	err := s.db.Where("company_id = ?", companyID).First(&settings).Error
	if err == nil {
		return &settings, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	var company models.Company
	if err := s.db.First(&company, companyID).Error; err != nil {
		return nil, err
	}
	settings = models.CompanySettings{CompanyID: companyID, DisplayName: company.Name, LogoURL: company.LogoURL, PrimaryColor: "#2563eb", SecondaryColor: "#0f172a", Locale: "es", CareerPageEnabled: true}
	if err := s.db.Create(&settings).Error; err != nil {
		return nil, err
	}
	return &settings, nil
}

func (s *CompanySettingsService) Update(companyID uint, dto UpdateCompanySettingsDTO) (*models.CompanySettings, error) {
	settings, err := s.Get(companyID)
	if err != nil {
		return nil, err
	}
	if dto.DisplayName != nil {
		settings.DisplayName = *dto.DisplayName
	}
	if dto.LogoURL != nil {
		settings.LogoURL = *dto.LogoURL
	}
	if dto.PrimaryColor != nil {
		settings.PrimaryColor = *dto.PrimaryColor
	}
	if dto.SecondaryColor != nil {
		settings.SecondaryColor = *dto.SecondaryColor
	}
	if dto.Locale != nil {
		settings.Locale = *dto.Locale
	}
	if dto.CareerPageTitle != nil {
		settings.CareerPageTitle = *dto.CareerPageTitle
	}
	if dto.CareerPageDescription != nil {
		settings.CareerPageDescription = *dto.CareerPageDescription
	}
	if dto.CareerPageEnabled != nil {
		settings.CareerPageEnabled = *dto.CareerPageEnabled
	}
	if dto.CustomDomain != nil {
		settings.CustomDomain = *dto.CustomDomain
	}
	if dto.ReplyToEmail != nil {
		settings.ReplyToEmail = *dto.ReplyToEmail
	}
	if err := s.db.Save(settings).Error; err != nil {
		return nil, err
	}
	return settings, nil
}

func (s *CompanySettingsService) Subscription(companyID uint) (*models.CompanySubscription, error) {
	subscription, err := repository.CurrentSubscription(s.db, companyID)
	if err != nil {
		return nil, err
	}
	if subscription == nil {
		return nil, apperr.NotFound("company has no active subscription")
	}
	return subscription, nil
}

func (s *CompanySettingsService) Usage(companyID uint) (*SubscriptionUsageDTO, error) {
	usage := &SubscriptionUsageDTO{}
	queries := []struct {
		model  interface{}
		target *int64
	}{
		{&models.Membership{}, &usage.Users}, {&models.Job{}, &usage.Jobs}, {&models.Candidate{}, &usage.Candidates}, {&models.Application{}, &usage.Applications},
	}
	for _, query := range queries {
		if err := s.db.Model(query.model).Where("company_id = ?", companyID).Count(query.target).Error; err != nil {
			return nil, err
		}
	}
	if err := s.db.Model(&models.StoredFile{}).Where("company_id = ?", companyID).Select("COALESCE(SUM(size_bytes), 0)").Scan(&usage.StorageBytes).Error; err != nil {
		return nil, err
	}
	return usage, nil
}
