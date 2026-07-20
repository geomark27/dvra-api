package seeders

import (
	"fmt"
	"log"

	"gorm.io/gorm"
)

// DatabaseSeeder orchestrates all seeders
type DatabaseSeeder struct{}

// Run executes all registered seeders
func (s *DatabaseSeeder) Run(db *gorm.DB) error {
	log.Println("🌱 Running seeders...")

	for _, registered := range AllSeeders {
		if err := registered.Seeder.Run(db); err != nil {
			log.Printf("❌ Seeder failed: %v", err)
			return err
		}
	}

	log.Println("✅ All seeders executed successfully")
	return nil
}

// RunModule executes only seeders owned by a bounded context. Demo data is a
// separate explicit module and is never pulled as a dependency.
func (s *DatabaseSeeder) RunModule(db *gorm.DB, module string) error {
	found := false
	for _, registered := range AllSeeders {
		if registered.Module != module {
			continue
		}
		found = true
		log.Printf("🌱 Running %s seeder", registered.Name)
		if err := registered.Seeder.Run(db); err != nil {
			return err
		}
	}
	if !found {
		return fmt.Errorf("no seeders registered for module %q", module)
	}
	return nil
}
