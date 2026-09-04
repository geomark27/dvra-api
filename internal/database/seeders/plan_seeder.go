package seeders

import (
	catalog "dvra-api/internal/modules/catalog/repository"
	"dvra-api/internal/shared/models"
	"fmt"

	"gorm.io/gorm"
)

// PlanSeeder seeds default subscription plans
type PlanSeeder struct{}

// Run executes the plan seeder
func (s *PlanSeeder) Run(db *gorm.DB) error {
	return SeedPlans(db)
}

// SeedPlans creates the default subscription plans
func SeedPlans(db *gorm.DB) error {
	usd, err := catalog.ResolveGlobal(db, "billing", "currency", "USD")
	if err != nil || usd == nil {
		return fmt.Errorf("billing currency USD is not seeded")
	}
	monthly, err := catalog.ResolveGlobal(db, "billing", "billing_cycle", "monthly")
	if err != nil || monthly == nil {
		return fmt.Errorf("billing cycle monthly is not seeded")
	}
	email, err := catalog.ResolveGlobal(db, "billing", "support_level", "email")
	if err != nil || email == nil {
		return fmt.Errorf("support level email is not seeded")
	}
	priority, err := catalog.ResolveGlobal(db, "billing", "support_level", "priority")
	if err != nil || priority == nil {
		return fmt.Errorf("support level priority is not seeded")
	}
	dedicated, err := catalog.ResolveGlobal(db, "billing", "support_level", "dedicated")
	if err != nil || dedicated == nil {
		return fmt.Errorf("support level dedicated is not seeded")
	}
	plans := []models.Plan{
		{
			Name:                "Free",
			Slug:                "free",
			Description:         "Perfect for trying out Dvra ATS. Limited features to get you started.",
			Price:               0.00,
			CurrencyValueID:     usd.ID,
			BillingCycleValueID: monthly.ID,
			IsActive:            true,
			IsPublic:            true,
			TrialDays:           0,
			DisplayOrder:        1,
			MaxUsers:            2,
			MaxJobs:             3,
			MaxCandidates:       50,
			MaxApplications:     100,
			MaxStorageGB:        1,
			CanExportData:       false,
			CanUseCustomBrand:   false,
			CanUseAPI:           false,
			CanUseIntegrations:  false,
			CanUseStaffing:      false,
			SupportLevelValueID: email.ID,
		},
		{
			Name:                "Starter",
			Slug:                "starter",
			Description:         "Ideal for small teams starting their recruitment journey.",
			Price:               39.99,
			CurrencyValueID:     usd.ID,
			BillingCycleValueID: monthly.ID,
			IsActive:            true,
			IsPublic:            true,
			TrialDays:           14,
			DisplayOrder:        2,
			MaxUsers:            5,
			MaxJobs:             10,
			MaxCandidates:       200,
			MaxApplications:     500,
			MaxStorageGB:        5,
			CanExportData:       true,
			CanUseCustomBrand:   false,
			CanUseAPI:           false,
			CanUseIntegrations:  false,
			CanUseStaffing:      false,
			SupportLevelValueID: email.ID,
		},
		{
			Name:                "Professional",
			Slug:                "professional",
			Description:         "For growing companies with advanced recruitment needs.",
			Price:               79.99,
			CurrencyValueID:     usd.ID,
			BillingCycleValueID: monthly.ID,
			IsActive:            true,
			IsPublic:            true,
			TrialDays:           14,
			DisplayOrder:        3,
			MaxUsers:            15,
			MaxJobs:             50,
			MaxCandidates:       1000,
			MaxApplications:     5000,
			MaxStorageGB:        20,
			CanExportData:       true,
			CanUseCustomBrand:   true,
			CanUseAPI:           true,
			CanUseIntegrations:  true,
			CanUseStaffing:      true,
			SupportLevelValueID: priority.ID,
		},
		{
			Name:                "Enterprise",
			Slug:                "enterprise",
			Description:         "Unlimited power for large organizations with complex hiring processes.",
			Price:               159.99,
			CurrencyValueID:     usd.ID,
			BillingCycleValueID: monthly.ID,
			IsActive:            true,
			IsPublic:            true,
			TrialDays:           30,
			DisplayOrder:        4,
			MaxUsers:            -1, // Unlimited
			MaxJobs:             -1, // Unlimited
			MaxCandidates:       -1, // Unlimited
			MaxApplications:     -1, // Unlimited
			MaxStorageGB:        -1, // Unlimited
			CanExportData:       true,
			CanUseCustomBrand:   true,
			CanUseAPI:           true,
			CanUseIntegrations:  true,
			CanUseStaffing:      true,
			SupportLevelValueID: dedicated.ID,
		},
	}

	for _, plan := range plans {
		// Check if plan already exists
		var existingPlan models.Plan
		err := db.Where("slug = ?", plan.Slug).First(&existingPlan).Error
		if err == gorm.ErrRecordNotFound {
			// Plan doesn't exist, create it
			if err := db.Create(&plan).Error; err != nil {
				return err
			}
		}
		// If plan exists, skip (don't update existing plans)
	}

	return nil
}
