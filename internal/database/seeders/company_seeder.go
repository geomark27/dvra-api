package seeders

import (
	"dvra-api/internal/app/models"
	"log"
	"time"

	"gorm.io/gorm"
)

// CompanySeeder seeds the company table
type CompanySeeder struct{}

// Run implements the Seeder interface
func (s *CompanySeeder) Run(db *gorm.DB) error {
	var existing models.Company
	result := db.Where("slug = ?", "azentic-sys").First(&existing)

	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		log.Printf("❌ Error checking existing company: %v", result.Error)
		return result.Error
	}

	if result.Error == gorm.ErrRecordNotFound {
		// Crear empresa inicial
		trialEnds := time.Now().AddDate(0, 1, 0) // 1 mes de trial
		company := models.Company{
			Name:        "Azentic Sys",
			Slug:        "azentic-sys",
			LogoURL:     "",
			PlanTier:    "free",
			TrialEndsAt: &trialEnds,
			Timezone:    "America/Bogota",
		}

		if err := db.Create(&company).Error; err != nil {
			log.Printf("❌ Error creating company: %v", err)
			return err
		}

		log.Printf("✅ Company '%s' created successfully (ID: %d)", company.Name, company.ID)

		// Crear membership entre admin user y la empresa
		var adminUser models.User
		if err := db.Where("email = ?", "admin@azentic.com").First(&adminUser).Error; err != nil {
			log.Printf("❌ Admin user not found: %v", err)
			return err
		}

		joinedAt := time.Now()
		membership := models.Membership{
			UserID:    adminUser.ID,
			CompanyID: &company.ID,
			Role:      models.RoleAdmin, // Admin de la empresa
			Status:    "active",
			IsDefault: true,
			JoinedAt:  &joinedAt,
		}

		if err := db.Create(&membership).Error; err != nil {
			log.Printf("❌ Error creating membership: %v", err)
			return err
		}

		log.Printf("✅ Membership created: User '%s' is admin of '%s'", adminUser.Name, company.Name)
	} else {
		log.Println("⏭️  Company 'Azentic Sys' already exists, skipping")
	}

	return nil
}
