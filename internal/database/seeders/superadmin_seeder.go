package seeders

import (
	"dvra-api/internal/app/models"
	"log"
	"time"

	"gorm.io/gorm"
)

// SuperAdminSeeder asigna el rol de SuperAdmin al usuario modo dios
type SuperAdminSeeder struct{}

// Run executes the superadmin role assignment
func (s *SuperAdminSeeder) Run(db *gorm.DB) error {
	// Buscar el usuario SuperAdmin
	var superAdmin models.User
	if err := db.Where("email = ?", "superadmin@dvra.com").First(&superAdmin).Error; err != nil {
		log.Printf("❌ SuperAdmin user not found: %v", err)
		return err
	}

	// Verificar si ya tiene un membership con rol superadmin
	// NOTA: SuperAdmin NO pertenece a ninguna empresa específica
	// Pero podríamos crear un membership especial con CompanyID = 0 o NULL
	// O simplemente guardar su rol en una tabla separada
	// Por ahora, vamos a crear un membership "global" con CompanyID = 0

	var existingMembership models.Membership
	result := db.Where("user_id = ? AND role = ?", superAdmin.ID, models.RoleSuperAdmin).First(&existingMembership)

	if result.Error == gorm.ErrRecordNotFound {
		joinedAt := time.Now()
		membership := models.Membership{
			UserID:    superAdmin.ID,
			CompanyID: nil, // nil = Sin empresa (SuperAdmin global)
			Role:      models.RoleSuperAdmin,
			Status:    "active",
			IsDefault: true,
			JoinedAt:  &joinedAt,
		}

		if err := db.Create(&membership).Error; err != nil {
			log.Printf("❌ Error creating superadmin membership: %v", err)
			return err
		}

		log.Printf("✅ SuperAdmin role assigned to '%s' (Global access)", superAdmin.Name)
	} else if result.Error != nil {
		log.Printf("❌ Error checking superadmin membership: %v", result.Error)
		return result.Error
	} else {
		log.Println("⏭️  SuperAdmin membership already exists, skipping")
	}

	return nil
}
