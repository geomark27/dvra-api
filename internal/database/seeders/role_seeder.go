package seeders

import (
	"dvra-api/internal/app/models"
	"log"

	"gorm.io/gorm"
)

// RoleSeeder seeds initial roles
type RoleSeeder struct{}

// Run executes the role seeder
func (s *RoleSeeder) Run(db *gorm.DB) error {
	roles := []models.Role{
		{
			Name:        "Super Admin",
			Slug:        models.RoleSuperAdmin,
			Description: "Acceso total al sistema, gestiona todas las empresas",
			Level:       100,
			IsSystem:    true,
		},
		{
			Name:        "Admin",
			Slug:        models.RoleAdmin,
			Description: "Administrador de empresa, acceso total a su empresa",
			Level:       50,
			IsSystem:    true,
		},
		{
			Name:        "Recruiter",
			Slug:        models.RoleRecruiter,
			Description: "Reclutador, gestiona jobs y candidatos",
			Level:       30,
			IsSystem:    true,
		},
		{
			Name:        "User",
			Slug:        models.RoleUser,
			Description: "Usuario genérico con permisos limitados",
			Level:       10,
			IsSystem:    true,
		},
	}

	for _, role := range roles {
		var existing models.Role
		result := db.Where("slug = ?", role.Slug).First(&existing)

		if result.Error == gorm.ErrRecordNotFound {
			if err := db.Create(&role).Error; err != nil {
				log.Printf("❌ Error creating role %s: %v", role.Slug, err)
				return err
			}
			log.Printf("✅ Role '%s' created successfully", role.Name)
		} else if result.Error != nil {
			log.Printf("❌ Error checking role %s: %v", role.Slug, result.Error)
			return result.Error
		} else {
			log.Printf("⏭️  Role '%s' already exists, skipping", role.Name)
		}
	}

	return nil
}
