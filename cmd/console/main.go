package main

import (
	"log"
	"strings"

	"dvra-api/internal/database"
	"dvra-api/internal/database/migrations"
	"dvra-api/internal/database/seeders"
	"dvra-api/internal/platform/config"
	"dvra-api/internal/shared/models"
	"dvra-api/internal/shared/permissions"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

func main() {
	// Cargar variables de entorno desde .env
	_ = godotenv.Load()

	rootCmd := &cobra.Command{
		Use:   "console",
		Short: "Database console tools for <no value>",
		Long:  `Command-line tools for managing database migrations and seeders`,
	}

	// migrate command
	migrateCmd := &cobra.Command{
		Use:   "migrate",
		Short: "Run database migrations",
		Long:  `Automatically migrate database schema based on your models`,
		Run:   runMigrate,
	}
	migrateCmd.Flags().Bool("seed", false, "Run seeders after migration")
	migrateCmd.Flags().Bool("fresh", false, "Drop all tables and migrate from scratch")
	migrateCmd.Flags().String("module", "", "Migrate one module (core, identity, recruitment, staffing, billing)")

	// seed command
	seedCmd := &cobra.Command{
		Use:   "seed",
		Short: "Run database seeders",
		Long:  `Populate database with initial or test data`,
		Run:   runSeed,
	}
	seedCmd.Flags().String("module", "", "Seed one module (core, identity, recruitment, staffing, billing, demo)")
	platformAdminCmd := &cobra.Command{
		Use:   "platform-admin",
		Short: "Grant landlord superadmin access to an existing user",
		Run:   runPlatformAdmin,
	}
	platformAdminCmd.Flags().String("email", "", "Email of the existing user")
	_ = platformAdminCmd.MarkFlagRequired("email")

	rootCmd.AddCommand(migrateCmd, seedCmd, platformAdminCmd)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func runPlatformAdmin(cmd *cobra.Command, args []string) {
	cfg := config.Load()
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}
	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer func() { _ = database.CloseDB() }()

	email, _ := cmd.Flags().GetString("email")
	email = strings.ToLower(strings.TrimSpace(email))
	var user models.User
	if err := db.Where("LOWER(email) = ?", email).First(&user).Error; err != nil {
		log.Fatalf("User %q was not found; create the identity first", email)
	}
	profile := models.PlatformMembership{UserID: user.ID, Role: permissions.RoleSuperAdmin, Status: "active"}
	if err := db.Where("user_id = ?", user.ID).
		Assign(map[string]interface{}{"role": permissions.RoleSuperAdmin, "status": "active"}).
		FirstOrCreate(&profile).Error; err != nil {
		log.Fatalf("Error granting platform access: %v", err)
	}
	log.Printf("Platform superadmin access granted to %s", email)
}

func runMigrate(cmd *cobra.Command, args []string) {
	// Load configuration
	cfg := config.Load()
	if err := cfg.Validate(); err != nil {
		log.Fatalf("❌ Invalid configuration: %v", err)
	}

	// Initialize database
	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatalf("❌ Error connecting to database: %v", err)
	}
	defer func() { _ = database.CloseDB() }()

	// Check for fresh flag
	fresh, _ := cmd.Flags().GetBool("fresh")
	module, _ := cmd.Flags().GetString("module")
	module = strings.TrimSpace(module)
	if fresh && module != "" {
		log.Fatal("❌ --fresh cannot be combined with --module because it would invalidate migration history for other modules")
	}
	if fresh {
		log.Println("🗑️  Dropping all tables...")
		modelsToDrop := make([]interface{}, 0, len(database.AllModels)+1)
		modelsToDrop = append(modelsToDrop, &models.SchemaMigration{})
		modelsToDrop = append(modelsToDrop, database.AllModels...)
		if err := db.Migrator().DropTable(modelsToDrop...); err != nil {
			log.Printf("⚠️  Warning dropping tables: %v", err)
		}
	}

	// Bootstrap migration metadata and landlord profiles before tightening the
	// tenant membership schema.
	log.Println("🔄 Running migrations...")
	if !db.Migrator().HasTable(&models.SchemaMigration{}) {
		if err := db.AutoMigrate(&models.SchemaMigration{}); err != nil {
			log.Fatalf("❌ Migration bootstrap error: %v", err)
		}
	}
	if db.Migrator().HasTable("memberships") {
		if err := db.AutoMigrate(&models.PlatformMembership{}); err != nil {
			log.Fatalf("❌ Platform membership migration error: %v", err)
		}
	}
	if module == "" {
		err = migrations.RunBeforeSchema(db)
	} else {
		err = migrations.RunBeforeSchemaModule(db, module)
	}
	if err != nil {
		log.Fatalf("❌ Pre-schema data migration error: %v", err)
	}
	modelsToMigrate := database.AllModels
	if module != "" {
		modelsToMigrate, err = database.ModelsForModule(module)
		if err != nil {
			log.Fatalf("❌ %v", err)
		}
	}
	if err := db.AutoMigrate(modelsToMigrate...); err != nil {
		log.Fatalf("❌ Migration error: %v", err)
	}
	if module == "" {
		err = migrations.Run(db)
	} else {
		err = migrations.RunModule(db, module)
	}
	if err != nil {
		log.Fatalf("❌ Data migration error: %v", err)
	}
	log.Println("✅ Migrations completed successfully")

	// Check for seed flag
	withSeed, _ := cmd.Flags().GetBool("seed")
	if withSeed {
		runSeedLogic(db, module)
	}
}

func runSeed(cmd *cobra.Command, args []string) {
	// Load configuration
	cfg := config.Load()
	if err := cfg.Validate(); err != nil {
		log.Fatalf("❌ Invalid configuration: %v", err)
	}

	// Initialize database
	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatalf("❌ Error connecting to database: %v", err)
	}
	defer func() { _ = database.CloseDB() }()

	module, _ := cmd.Flags().GetString("module")
	runSeedLogic(db, strings.TrimSpace(module))
}

func runSeedLogic(db *gorm.DB, module string) {
	seeder := &seeders.DatabaseSeeder{}
	var err error
	if module == "" {
		err = seeder.Run(db)
	} else {
		err = seeder.RunModule(db, module)
	}
	if err != nil {
		log.Fatalf("❌ Seeder error: %v", err)
	}
}
