package seeders

import (
	"dvra-api/internal/app/models"
	"log"

	"gorm.io/gorm"
)

// SystemValueSeeder implements the Seeder interface
type SystemValueSeeder struct{}

// Run executes the seeder
func (s *SystemValueSeeder) Run(db *gorm.DB) error {
	return SystemValueSeederFunc(db)
}

// SystemValueSeederFunc seeds initial system values
func SystemValueSeederFunc(db *gorm.DB) error {
	log.Println("🌱 Seeding system values...")

	systemValues := []models.SystemValue{
		// Job statuses
		{Category: "job_status", Value: "draft", Label: "Borrador (no visible para candidatos)", DisplayOrder: 1, IsActive: true},
		{Category: "job_status", Value: "published", Label: "Publicada (visible para candidatos)", DisplayOrder: 2, IsActive: true},
		{Category: "job_status", Value: "closed", Label: "Cerrada", DisplayOrder: 3, IsActive: true},

		// Application statuses
		{Category: "application_status", Value: "applied", Label: "Aplicado", DisplayOrder: 1, IsActive: true},
		{Category: "application_status", Value: "screening", Label: "En Revisión", DisplayOrder: 2, IsActive: true},
		{Category: "application_status", Value: "technical", Label: "Prueba Técnica", DisplayOrder: 3, IsActive: true},
		{Category: "application_status", Value: "interview", Label: "Entrevista", DisplayOrder: 4, IsActive: true},
		{Category: "application_status", Value: "offer", Label: "Oferta", DisplayOrder: 5, IsActive: true},
		{Category: "application_status", Value: "hired", Label: "Contratado", DisplayOrder: 6, IsActive: true},
		{Category: "application_status", Value: "rejected", Label: "Rechazado", DisplayOrder: 7, IsActive: true},

		// Contract types
		{Category: "contract_type", Value: "full_time", Label: "Tiempo Completo", DisplayOrder: 1, IsActive: true},
		{Category: "contract_type", Value: "part_time", Label: "Medio Tiempo", DisplayOrder: 2, IsActive: true},
		{Category: "contract_type", Value: "contractor", Label: "Contratista", DisplayOrder: 3, IsActive: true},
		{Category: "contract_type", Value: "internship", Label: "Pasantía", DisplayOrder: 4, IsActive: true},
		{Category: "contract_type", Value: "temporary", Label: "Temporal", DisplayOrder: 5, IsActive: true},

		// Work modes
		{Category: "work_mode", Value: "remote", Label: "Remoto", DisplayOrder: 1, IsActive: true},
		{Category: "work_mode", Value: "onsite", Label: "Presencial", DisplayOrder: 2, IsActive: true},
		{Category: "work_mode", Value: "hybrid", Label: "Híbrido", DisplayOrder: 3, IsActive: true},

		// Experience levels
		{Category: "experience_level", Value: "entry", Label: "Junior / Sin experiencia", DisplayOrder: 1, IsActive: true},
		{Category: "experience_level", Value: "mid", Label: "Semi-Senior (2-5 años)", DisplayOrder: 2, IsActive: true},
		{Category: "experience_level", Value: "senior", Label: "Senior (5+ años)", DisplayOrder: 3, IsActive: true},
		{Category: "experience_level", Value: "lead", Label: "Lead / Principal", DisplayOrder: 4, IsActive: true},

		// Priority levels
		{Category: "priority", Value: "low", Label: "Baja", DisplayOrder: 1, IsActive: true},
		{Category: "priority", Value: "medium", Label: "Media", DisplayOrder: 2, IsActive: true},
		{Category: "priority", Value: "high", Label: "Alta", DisplayOrder: 3, IsActive: true},
		{Category: "priority", Value: "urgent", Label: "Urgente", DisplayOrder: 4, IsActive: true},

		// Source (where candidate came from)
		{Category: "candidate_source", Value: "linkedin", Label: "LinkedIn", DisplayOrder: 1, IsActive: true},
		{Category: "candidate_source", Value: "website", Label: "Sitio Web", DisplayOrder: 2, IsActive: true},
		{Category: "candidate_source", Value: "referral", Label: "Referido", DisplayOrder: 3, IsActive: true},
		{Category: "candidate_source", Value: "job_board", Label: "Bolsa de Trabajo", DisplayOrder: 4, IsActive: true},
		{Category: "candidate_source", Value: "direct", Label: "Aplicación Directa", DisplayOrder: 5, IsActive: true},
		{Category: "candidate_source", Value: "other", Label: "Otro", DisplayOrder: 6, IsActive: true},
	}

	for _, value := range systemValues {
		// Check if exists (category + value combination)
		var existing models.SystemValue
		result := db.Where("category = ? AND value = ? AND company_id IS NULL", value.Category, value.Value).First(&existing)

		if result.Error == gorm.ErrRecordNotFound {
			if err := db.Create(&value).Error; err != nil {
				log.Printf("⚠️  Error creating system value %s.%s: %v", value.Category, value.Value, err)
				continue
			}
			log.Printf("  ✅ Created: %s.%s (%s)", value.Category, value.Value, value.Label)
		} else {
			log.Printf("  ⏭️  Skipped: %s.%s (already exists)", value.Category, value.Value)
		}
	}

	log.Println("✅ System values seeded successfully")
	return nil
}
