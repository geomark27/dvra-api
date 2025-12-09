package seeders

import (
	"dvra-api/internal/app/models"

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
	plans := []models.Plan{
		{
			Name:               "Free",
			Slug:               "free",
			Description:        "Perfect for trying out Dvra ATS. Limited features to get you started.",
			Price:              0.00,
			Currency:           "USD",
			BillingCycle:       "monthly",
			IsActive:           true,
			IsPublic:           true,
			TrialDays:          0,
			DisplayOrder:       1,
			MaxUsers:           2,
			MaxJobs:            3,
			MaxCandidates:      50,
			MaxApplications:    100,
			MaxStorageGB:       1,
			CanExportData:      false,
			CanUseCustomBrand:  false,
			CanUseAPI:          false,
			CanUseIntegrations: false,
			SupportLevel:       "email",
		},
		{
			Name:               "Starter",
			Slug:               "starter",
			Description:        "Ideal for small teams starting their recruitment journey.",
			Price:              39.99,
			Currency:           "USD",
			BillingCycle:       "monthly",
			IsActive:           true,
			IsPublic:           true,
			TrialDays:          14,
			DisplayOrder:       2,
			MaxUsers:           5,
			MaxJobs:            10,
			MaxCandidates:      200,
			MaxApplications:    500,
			MaxStorageGB:       5,
			CanExportData:      true,
			CanUseCustomBrand:  false,
			CanUseAPI:          false,
			CanUseIntegrations: false,
			SupportLevel:       "email",
		},
		{
			Name:               "Professional",
			Slug:               "professional",
			Description:        "For growing companies with advanced recruitment needs.",
			Price:              79.99,
			Currency:           "USD",
			BillingCycle:       "monthly",
			IsActive:           true,
			IsPublic:           true,
			TrialDays:          14,
			DisplayOrder:       3,
			MaxUsers:           15,
			MaxJobs:            50,
			MaxCandidates:      1000,
			MaxApplications:    5000,
			MaxStorageGB:       20,
			CanExportData:      true,
			CanUseCustomBrand:  true,
			CanUseAPI:          true,
			CanUseIntegrations: true,
			SupportLevel:       "priority",
		},
		{
			Name:               "Enterprise",
			Slug:               "enterprise",
			Description:        "Unlimited power for large organizations with complex hiring processes.",
			Price:              159.99,
			Currency:           "USD",
			BillingCycle:       "monthly",
			IsActive:           true,
			IsPublic:           true,
			TrialDays:          30,
			DisplayOrder:       4,
			MaxUsers:           -1, // Unlimited
			MaxJobs:            -1, // Unlimited
			MaxCandidates:      -1, // Unlimited
			MaxApplications:    -1, // Unlimited
			MaxStorageGB:       -1, // Unlimited
			CanExportData:      true,
			CanUseCustomBrand:  true,
			CanUseAPI:          true,
			CanUseIntegrations: true,
			SupportLevel:       "dedicated",
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
