package seeders

import (
	"dvra-api/internal/app/models"
	"dvra-api/internal/shared/catalog"
	"dvra-api/internal/shared/permissions"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/gorm"
)

// CompanySeeder seeds the company table
type CompanySeeder struct{}

// Run implements the Seeder interface
func (s *CompanySeeder) Run(db *gorm.DB) error {
	if os.Getenv("SEED_DEMO_DATA") != "true" {
		log.Println("⏭️  Demo company disabled (set SEED_DEMO_DATA=true to enable)")
		return nil
	}
	if os.Getenv("ENVIRONMENT") == "production" {
		return fmt.Errorf("demo company cannot be seeded in production")
	}
	var existing models.Company
	result := db.Where("slug = ?", "azentic-sys").First(&existing)

	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		log.Printf("❌ Error checking existing company: %v", result.Error)
		return result.Error
	}

	if result.Error == gorm.ErrRecordNotFound {
		// Crear empresa inicial
		company := models.Company{
			Name:        "Azentic Sys",
			Slug:        "azentic-sys",
			LogoURL:     "",
			TrialEndsAt: nil,
			Timezone:    "America/Bogota",
		}

		if err := db.Create(&company).Error; err != nil {
			log.Printf("❌ Error creating company: %v", err)
			return err
		}
		var plan models.Plan
		if err := db.Where("slug = ?", "free").First(&plan).Error; err != nil {
			return err
		}
		status, err := catalog.ResolveGlobal(db, "billing", "subscription_status", "active")
		if err != nil || status == nil {
			return fmt.Errorf("resolve subscription status: %w", err)
		}
		if plan.CurrencyValueID == 0 || plan.BillingCycleValueID == 0 {
			return fmt.Errorf("plan billing catalogs are missing")
		}
		now := time.Now().UTC()
		subscription := models.CompanySubscription{CompanyID: company.ID, PlanID: plan.ID, StatusValueID: status.ID, CurrencyValueID: plan.CurrencyValueID, BillingCycleValueID: plan.BillingCycleValueID, PriceSnapshot: plan.Price, StartedAt: now, CurrentPeriodStart: now}
		if err := db.Create(&subscription).Error; err != nil {
			return err
		}
		settings := models.CompanySettings{CompanyID: company.ID, DisplayName: company.Name, PrimaryColor: "#2563eb", SecondaryColor: "#0f172a", Locale: "es", CareerPageEnabled: true}
		if err := db.Create(&settings).Error; err != nil {
			return err
		}

		log.Printf("✅ Company '%s' created successfully (ID: %d)", company.Name, company.ID)

		// Crear membership entre admin user y la empresa
		var adminUser models.User
		if err := db.Where("email = ?", "admin@azentic.com").First(&adminUser).Error; err != nil {
			log.Printf("❌ Admin user not found: %v", err)
			return err
		}
		var adminRole models.TenantRole
		if err := db.Where("key = ? AND is_active = ?", permissions.RoleAdmin, true).First(&adminRole).Error; err != nil {
			log.Printf("❌ Tenant admin role not found: %v", err)
			return err
		}

		joinedAt := time.Now()
		membership := models.Membership{
			UserID:       adminUser.ID,
			CompanyID:    company.ID,
			TenantRoleID: adminRole.ID,
			Status:       "active",
			IsDefault:    true,
			JoinedAt:     &joinedAt,
		}

		if err := db.Create(&membership).Error; err != nil {
			log.Printf("❌ Error creating membership: %v", err)
			return err
		}

		name := adminUser.FirstName + " " + adminUser.LastName
		log.Printf("✅ Membership created: User '%s' is admin of '%s'", name, company.Name)
	} else {
		log.Println("⏭️  Company 'Azentic Sys' already exists, skipping")
	}

	return nil
}
