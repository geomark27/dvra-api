package seeders

import (
	"dvra-api/internal/app/models"
	"log"

	"gorm.io/gorm"
)

// PlatformSettingsSeeder seeds default platform settings
type PlatformSettingsSeeder struct{}

// Run executes the platform settings seeder
func (s *PlatformSettingsSeeder) Run(db *gorm.DB) error {
	return SeedPlatformSettings(db)
}

// SeedPlatformSettings creates the default platform settings (singleton)
func SeedPlatformSettings(db *gorm.DB) error {
	// Check if settings already exist
	var count int64
	if err := db.Model(&models.PlatformSettings{}).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		log.Println("⏭️  Platform settings already exist, skipping...")
		return nil
	}

	settings := models.PlatformSettings{
		// Branding - valores por defecto configurables
		PlatformName:  "DVRA ATS",
		PlatformShort: "DVRA",
		Tagline:       "Modern Applicant Tracking System for Growing Teams",
		LogoURL:       "",
		LogoDarkURL:   "",
		FaviconURL:    "",
		PrimaryColor:  "#0ea5e9", // sky-500

		// Contact
		SupportEmail: "support@dvra.io",
		SalesEmail:   "sales@dvra.io",

		// URLs
		MarketingURL: "https://dvra.io",
		DocsURL:      "https://docs.dvra.io",
		TermsURL:     "/legal/terms",
		PrivacyURL:   "/legal/privacy",

		// Business defaults
		DefaultTrialDays: 14,
		DefaultPlanTier:  "starter",

		// Legal (vacío por defecto, configurable)
		LegalCompanyName: "",
		LegalAddress:     "",
		LegalCountry:     "",
		LegalTaxID:       "",

		// Social media
		TwitterURL:  "",
		LinkedinURL: "",
		GithubURL:   "",
	}

	if err := db.Create(&settings).Error; err != nil {
		return err
	}

	log.Println("✅ Platform settings seeded successfully")
	return nil
}
