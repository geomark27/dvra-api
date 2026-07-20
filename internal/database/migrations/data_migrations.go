package migrations

import (
	"fmt"
	"strings"
	"time"

	"dvra-api/internal/app/models"
	"dvra-api/internal/database/seeders"
	"dvra-api/internal/shared/permissions"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type migration struct {
	id     string
	module string
	run    func(*gorm.DB) error
}

var dataMigrations = []migration{
	{
		id:     "20260720_seed_modular_catalogs",
		module: "core",
		run: func(tx *gorm.DB) error {
			return createCatalogIndexes(tx)
		},
	},
	{id: "20260720_seed_identity_catalogs", module: "identity", run: func(tx *gorm.DB) error { return seeders.SeedCatalogs(tx, "identity") }},
	{id: "20260720_seed_recruitment_catalogs", module: "recruitment", run: func(tx *gorm.DB) error { return seeders.SeedCatalogs(tx, "recruitment") }},
	{id: "20260720_seed_staffing_catalogs", module: "staffing", run: func(tx *gorm.DB) error { return seeders.SeedCatalogs(tx, "staffing") }},
	{id: "20260720_seed_billing_catalogs", module: "billing", run: func(tx *gorm.DB) error { return seeders.SeedCatalogs(tx, "billing") }},
	{id: "20260720_company_subscriptions", module: "billing", run: migrateCompanySubscriptions},
	{id: "20260720_platform_default_plan", module: "billing", run: migratePlatformDefaultPlan},
	{
		id: "20260715_jobs_active_to_published", module: "recruitment",
		run: func(tx *gorm.DB) error {
			return tx.Model(&models.Job{}).Where("status = ?", "active").Update("status", "published").Error
		},
	},
	{
		id: "20260715_tenant_superadmin_to_admin", module: "identity",
		run: func(tx *gorm.DB) error {
			if !tx.Migrator().HasColumn("memberships", "role") {
				return nil
			}
			return tx.Model(&models.Membership{}).
				Where("company_id IS NOT NULL AND role = ?", "superadmin").
				Update("role", "admin").Error
		},
	},
	{
		id: "20260717_seed_relational_tenant_rbac", module: "identity",
		run: func(tx *gorm.DB) error {
			return seedTenantRBAC(tx)
		},
	},
	{
		id: "20260717_admin_full_tenant_access", module: "identity",
		run: func(tx *gorm.DB) error {
			return seedTenantRBAC(tx)
		},
	},
	{
		id: "20260720_departments_module", module: "recruitment",
		run: func(tx *gorm.DB) error {
			if err := seedTenantRBAC(tx); err != nil {
				return err
			}
			statements := []string{
				"CREATE UNIQUE INDEX IF NOT EXISTS idx_departments_company_parent_name ON departments (company_id, COALESCE(parent_department_id, 0), LOWER(name)) WHERE deleted_at IS NULL",
				"CREATE UNIQUE INDEX IF NOT EXISTS idx_departments_company_code ON departments (company_id, LOWER(code)) WHERE deleted_at IS NULL AND code <> ''",
			}
			for _, statement := range statements {
				if err := tx.Exec(statement).Error; err != nil {
					return err
				}
			}
			return nil
		},
	},
}

var beforeSchemaMigrations = []migration{
	{
		id: "20260720_relational_modular_catalogs", module: "core",
		run: migrateLegacySystemValues,
	},
	{id: "20260720_relational_plan_catalogs", module: "billing", run: migrateLegacyPlanCatalogs},
	{
		id: "20260715_split_platform_memberships", module: "identity",
		run: func(tx *gorm.DB) error {
			if !tx.Migrator().HasTable("memberships") {
				return nil
			}
			type legacyMembership struct {
				ID     uint
				UserID uint
				Role   string
				Status string
			}
			var legacy []legacyMembership
			if err := tx.Table("memberships").
				Where("company_id IS NULL").
				Find(&legacy).Error; err != nil {
				return err
			}
			for _, item := range legacy {
				if item.Role != "superadmin" {
					return fmt.Errorf("membership %d has no company and unsupported platform role %q", item.ID, item.Role)
				}
				status := item.Status
				if status == "" {
					status = "active"
				}
				profile := models.PlatformMembership{UserID: item.UserID, Role: item.Role, Status: status}
				if err := tx.Where("user_id = ?", item.UserID).FirstOrCreate(&profile).Error; err != nil {
					return err
				}
				if err := tx.Exec("DELETE FROM memberships WHERE id = ?", item.ID).Error; err != nil {
					return err
				}
			}
			type duplicate struct {
				UserID    uint
				CompanyID uint
				Count     int64
			}
			var duplicates []duplicate
			if err := tx.Table("memberships").
				Select("user_id, company_id, COUNT(*) AS count").
				Where("deleted_at IS NULL").
				Group("user_id, company_id").
				Having("COUNT(*) > 1").
				Scan(&duplicates).Error; err != nil {
				return err
			}
			if len(duplicates) > 0 {
				return fmt.Errorf("duplicate active tenant memberships must be resolved before migration")
			}
			return nil
		},
	},
	{
		id: "20260715_require_tenant_membership_company", module: "identity",
		run: func(tx *gorm.DB) error {
			if !tx.Migrator().HasTable("memberships") {
				return nil
			}
			var nullCount int64
			if err := tx.Table("memberships").Where("company_id IS NULL").Count(&nullCount).Error; err != nil {
				return err
			}
			if nullCount > 0 {
				return fmt.Errorf("%d memberships still have no tenant company", nullCount)
			}
			type duplicate struct {
				UserID    uint
				CompanyID uint
				Count     int64
			}
			var duplicates []duplicate
			if err := tx.Table("memberships").
				Select("user_id, company_id, COUNT(*) AS count").
				Where("deleted_at IS NULL").
				Group("user_id, company_id").
				Having("COUNT(*) > 1").
				Scan(&duplicates).Error; err != nil {
				return err
			}
			if len(duplicates) > 0 {
				return fmt.Errorf("duplicate active tenant memberships must be resolved before migration")
			}
			if err := tx.Exec("ALTER TABLE memberships ALTER COLUMN company_id SET NOT NULL").Error; err != nil {
				return err
			}
			if err := tx.Exec("DROP INDEX IF EXISTS idx_user_company").Error; err != nil {
				return err
			}
			return tx.Exec("CREATE UNIQUE INDEX idx_user_company ON memberships (user_id, company_id) WHERE deleted_at IS NULL").Error
		},
	},
	{
		id: "20260717_relational_tenant_rbac", module: "identity",
		run: func(tx *gorm.DB) error {
			if err := tx.AutoMigrate(&models.TenantRole{}, &models.Permission{}, &models.TenantRolePermission{}); err != nil {
				return err
			}
			if err := seedTenantRBAC(tx); err != nil {
				return err
			}
			// En una instalación nueva memberships todavía no existe. AutoMigrate
			// la creará después con tenant_role_id y sus constraints finales, por
			// lo que aquí solo corresponde migrar filas de una tabla legacy real.
			if !tx.Migrator().HasTable("memberships") {
				return nil
			}
			if !tx.Migrator().HasColumn("memberships", "tenant_role_id") {
				if err := tx.Exec("ALTER TABLE memberships ADD COLUMN tenant_role_id BIGINT").Error; err != nil {
					return err
				}
			}
			if tx.Migrator().HasColumn("memberships", "role") {
				if err := tx.Exec("UPDATE memberships SET role = 'interviewer' WHERE role = 'user'").Error; err != nil {
					return err
				}
				if err := tx.Exec("UPDATE memberships SET role = 'admin' WHERE role = 'superadmin'").Error; err != nil {
					return err
				}
				var invalidRoles []string
				if err := tx.Table("memberships").Distinct("role").
					Where("role NOT IN ?", permissions.TenantRoles()).Pluck("role", &invalidRoles).Error; err != nil {
					return err
				}
				if len(invalidRoles) > 0 {
					return fmt.Errorf("memberships contain unsupported tenant roles: %v", invalidRoles)
				}
				if err := tx.Exec(`UPDATE memberships AS m
					SET tenant_role_id = tr.id
					FROM tenant_roles AS tr
					WHERE tr.key = m.role`).Error; err != nil {
					return err
				}
			}
			var missing int64
			if err := tx.Table("memberships").Where("tenant_role_id IS NULL OR tenant_role_id = 0").Count(&missing).Error; err != nil {
				return err
			}
			if missing > 0 {
				return fmt.Errorf("%d memberships could not be assigned to a tenant role", missing)
			}
			statements := []string{
				"ALTER TABLE memberships ALTER COLUMN tenant_role_id SET NOT NULL",
				"CREATE INDEX IF NOT EXISTS idx_memberships_tenant_role_id ON memberships (tenant_role_id)",
				"ALTER TABLE memberships DROP CONSTRAINT IF EXISTS fk_memberships_tenant_role",
				"ALTER TABLE memberships ADD CONSTRAINT fk_memberships_tenant_role FOREIGN KEY (tenant_role_id) REFERENCES tenant_roles(id) ON UPDATE CASCADE ON DELETE RESTRICT",
				"ALTER TABLE memberships DROP CONSTRAINT IF EXISTS memberships_role_check",
			}
			for _, statement := range statements {
				if err := tx.Exec(statement).Error; err != nil {
					return err
				}
			}
			if tx.Migrator().HasColumn("memberships", "role") {
				return tx.Exec("ALTER TABLE memberships DROP COLUMN role").Error
			}
			return nil
		},
	},
}

func migrateLegacySystemValues(tx *gorm.DB) error {
	if !tx.Migrator().HasTable("system_values") || !tx.Migrator().HasColumn("system_values", "category") {
		return nil
	}
	statements := []string{
		`CREATE TABLE IF NOT EXISTS system_value_modules (
			id BIGSERIAL PRIMARY KEY, created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), deleted_at TIMESTAMPTZ,
			key VARCHAR(80) NOT NULL, name VARCHAR(120) NOT NULL, description TEXT NOT NULL DEFAULT '', is_system BOOLEAN NOT NULL DEFAULT TRUE, is_active BOOLEAN NOT NULL DEFAULT TRUE)`,
		"CREATE UNIQUE INDEX IF NOT EXISTS idx_system_value_modules_key ON system_value_modules (key)",
		"CREATE INDEX IF NOT EXISTS idx_system_value_modules_deleted_at ON system_value_modules (deleted_at)",
		`CREATE TABLE IF NOT EXISTS system_value_categories (
			id BIGSERIAL PRIMARY KEY, created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), deleted_at TIMESTAMPTZ,
			module_id BIGINT NOT NULL, key VARCHAR(80) NOT NULL, name VARCHAR(120) NOT NULL, description TEXT NOT NULL DEFAULT '', is_system BOOLEAN NOT NULL DEFAULT TRUE,
			is_extensible BOOLEAN NOT NULL DEFAULT FALSE, allow_tenant_values BOOLEAN NOT NULL DEFAULT FALSE, allow_tenant_override BOOLEAN NOT NULL DEFAULT FALSE, is_active BOOLEAN NOT NULL DEFAULT TRUE)`,
		"CREATE UNIQUE INDEX IF NOT EXISTS idx_system_value_category_key ON system_value_categories (module_id, key)",
		"CREATE INDEX IF NOT EXISTS idx_system_value_categories_deleted_at ON system_value_categories (deleted_at)",
		"ALTER TABLE system_values ADD COLUMN IF NOT EXISTS category_id BIGINT",
		"ALTER TABLE system_values ADD COLUMN IF NOT EXISTS key VARCHAR(100)",
		"ALTER TABLE system_values ADD COLUMN IF NOT EXISTS is_system BOOLEAN NOT NULL DEFAULT FALSE",
	}
	for _, statement := range statements {
		if err := tx.Exec(statement).Error; err != nil {
			return err
		}
	}
	if err := seeders.SeedCatalogStructure(tx, ""); err != nil {
		return err
	}

	var legacyCategories []string
	if err := tx.Table("system_values").Distinct("category").Pluck("category", &legacyCategories).Error; err != nil {
		return err
	}
	for _, legacyKey := range legacyCategories {
		moduleKey, categoryKey, known := seeders.CatalogOwner(legacyKey)
		if !known {
			moduleKey, categoryKey = "legacy", legacyKey
			module := models.SystemValueModule{Key: moduleKey, Name: "Catálogos heredados", IsSystem: true, IsActive: true}
			if err := tx.Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "key"}}, DoNothing: true}).Create(&module).Error; err != nil {
				return err
			}
			if err := tx.Where("key = ?", moduleKey).First(&module).Error; err != nil {
				return err
			}
			category := models.SystemValueCategory{ModuleID: module.ID, Key: categoryKey, Name: legacyKey, IsSystem: true, IsExtensible: true, AllowTenantValues: true, AllowTenantOverride: true, IsActive: true}
			if err := tx.Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "module_id"}, {Name: "key"}}, DoNothing: true}).Create(&category).Error; err != nil {
				return err
			}
		}
		if err := tx.Exec(`UPDATE system_values AS sv SET category_id = c.id, key = sv.value, is_system = (sv.company_id IS NULL)
			FROM system_value_categories AS c JOIN system_value_modules AS m ON m.id = c.module_id
			WHERE sv.category = ? AND m.key = ? AND c.key = ?`, legacyKey, moduleKey, categoryKey).Error; err != nil {
			return err
		}
	}
	var missing int64
	if err := tx.Table("system_values").Where("category_id IS NULL OR key IS NULL OR key = ''").Count(&missing).Error; err != nil {
		return err
	}
	if missing > 0 {
		return fmt.Errorf("%d legacy system values could not be assigned to a modular category", missing)
	}
	finalStatements := []string{
		"DROP INDEX IF EXISTS idx_system_values_category_value",
		"ALTER TABLE system_values ALTER COLUMN category_id SET NOT NULL",
		"ALTER TABLE system_values ALTER COLUMN key SET NOT NULL",
		"ALTER TABLE system_values DROP COLUMN category",
		"ALTER TABLE system_values DROP COLUMN value",
	}
	for _, statement := range finalStatements {
		if err := tx.Exec(statement).Error; err != nil {
			return err
		}
	}
	return createCatalogIndexes(tx)
}

func migrateLegacyPlanCatalogs(tx *gorm.DB) error {
	if !tx.Migrator().HasTable("plans") || !tx.Migrator().HasColumn("plans", "currency") {
		return nil
	}
	if err := seeders.SeedCatalogs(tx, "billing"); err != nil {
		return err
	}
	statements := []string{
		"ALTER TABLE plans ADD COLUMN IF NOT EXISTS currency_value_id BIGINT",
		"ALTER TABLE plans ADD COLUMN IF NOT EXISTS billing_cycle_value_id BIGINT",
		"ALTER TABLE plans ADD COLUMN IF NOT EXISTS support_level_value_id BIGINT",
		`UPDATE plans p SET currency_value_id = sv.id FROM system_values sv JOIN system_value_categories c ON c.id = sv.category_id JOIN system_value_modules m ON m.id = c.module_id WHERE m.key = 'billing' AND c.key = 'currency' AND sv.key = p.currency AND sv.company_id IS NULL`,
		`UPDATE plans p SET billing_cycle_value_id = sv.id FROM system_values sv JOIN system_value_categories c ON c.id = sv.category_id JOIN system_value_modules m ON m.id = c.module_id WHERE m.key = 'billing' AND c.key = 'billing_cycle' AND sv.key = p.billing_cycle AND sv.company_id IS NULL`,
		`UPDATE plans p SET support_level_value_id = sv.id FROM system_values sv JOIN system_value_categories c ON c.id = sv.category_id JOIN system_value_modules m ON m.id = c.module_id WHERE m.key = 'billing' AND c.key = 'support_level' AND sv.key = p.support_level AND sv.company_id IS NULL`,
	}
	for _, statement := range statements {
		if err := tx.Exec(statement).Error; err != nil {
			return err
		}
	}
	var missing int64
	if err := tx.Table("plans").Where("currency_value_id IS NULL OR billing_cycle_value_id IS NULL OR support_level_value_id IS NULL").Count(&missing).Error; err != nil {
		return err
	}
	if missing > 0 {
		return fmt.Errorf("%d plans contain unsupported billing catalog values", missing)
	}
	final := []string{
		"ALTER TABLE plans ALTER COLUMN currency_value_id SET NOT NULL", "ALTER TABLE plans ALTER COLUMN billing_cycle_value_id SET NOT NULL", "ALTER TABLE plans ALTER COLUMN support_level_value_id SET NOT NULL",
		"ALTER TABLE plans DROP COLUMN currency", "ALTER TABLE plans DROP COLUMN billing_cycle", "ALTER TABLE plans DROP COLUMN support_level",
	}
	for _, statement := range final {
		if err := tx.Exec(statement).Error; err != nil {
			return err
		}
	}
	return nil
}

func createCatalogIndexes(tx *gorm.DB) error {
	statements := []string{
		"CREATE UNIQUE INDEX IF NOT EXISTS idx_system_values_global_key ON system_values (category_id, key) WHERE company_id IS NULL AND deleted_at IS NULL",
		"CREATE UNIQUE INDEX IF NOT EXISTS idx_system_values_tenant_key ON system_values (category_id, company_id, key) WHERE company_id IS NOT NULL AND deleted_at IS NULL",
	}
	for _, statement := range statements {
		if err := tx.Exec(statement).Error; err != nil {
			return err
		}
	}
	return nil
}

func migrateCompanySubscriptions(tx *gorm.DB) error {
	if !tx.Migrator().HasColumn("companies", "plan_tier") {
		return nil
	}
	type legacyCompany struct {
		ID          uint
		PlanTier    string
		TrialEndsAt *time.Time
	}
	var companies []legacyCompany
	if err := tx.Table("companies").Select("id, plan_tier, trial_ends_at").Find(&companies).Error; err != nil {
		return err
	}
	resolve := func(category, key string) (*models.SystemValue, error) {
		var value models.SystemValue
		err := tx.Joins("JOIN system_value_categories c ON c.id = system_values.category_id").Joins("JOIN system_value_modules m ON m.id = c.module_id").
			Where("m.key = 'billing' AND c.key = ? AND system_values.key = ? AND system_values.company_id IS NULL", category, key).First(&value).Error
		return &value, err
	}
	status, err := resolve("subscription_status", "active")
	if err != nil {
		return err
	}
	for _, company := range companies {
		var count int64
		if err := tx.Model(&models.CompanySubscription{}).Where("company_id = ?", company.ID).Count(&count).Error; err != nil {
			return err
		}
		if count == 0 {
			var plan models.Plan
			if err := tx.Where("slug = ?", company.PlanTier).First(&plan).Error; err != nil {
				return fmt.Errorf("company %d references missing plan %q: %w", company.ID, company.PlanTier, err)
			}
			now := time.Now().UTC()
			subscription := models.CompanySubscription{CompanyID: company.ID, PlanID: plan.ID, StatusValueID: status.ID, CurrencyValueID: plan.CurrencyValueID, BillingCycleValueID: plan.BillingCycleValueID, PriceSnapshot: plan.Price, StartedAt: now, CurrentPeriodStart: now, TrialEndsAt: company.TrialEndsAt}
			if err := tx.Create(&subscription).Error; err != nil {
				return err
			}
		}
		settings := models.CompanySettings{CompanyID: company.ID, PrimaryColor: "#2563eb", SecondaryColor: "#0f172a", Locale: "es", CareerPageEnabled: true}
		if err := tx.Where("company_id = ?", company.ID).FirstOrCreate(&settings).Error; err != nil {
			return err
		}
	}
	return tx.Exec("ALTER TABLE companies DROP COLUMN plan_tier").Error
}

func migratePlatformDefaultPlan(tx *gorm.DB) error {
	if !tx.Migrator().HasColumn("platform_settings", "default_plan_tier") {
		return nil
	}
	if err := tx.Exec(`UPDATE platform_settings ps SET default_plan_id = p.id FROM plans p WHERE p.slug = ps.default_plan_tier`).Error; err != nil {
		return err
	}
	var missing int64
	if err := tx.Table("platform_settings").Where("default_plan_tier <> '' AND default_plan_id IS NULL").Count(&missing).Error; err != nil {
		return err
	}
	if missing > 0 {
		return fmt.Errorf("platform settings reference a missing default plan")
	}
	return tx.Exec("ALTER TABLE platform_settings DROP COLUMN default_plan_tier").Error
}

func seedTenantRBAC(tx *gorm.DB) error {
	for _, definition := range permissions.TenantRoleDefinitions() {
		role := models.TenantRole{
			Key: definition.Key, Name: definition.Name, Description: definition.Description,
			IsSystem: true, IsActive: true,
		}
		if err := tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "key"}},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"name": definition.Name, "description": definition.Description,
				"is_system": true, "is_active": true, "deleted_at": nil,
			}),
		}).Create(&role).Error; err != nil {
			return err
		}
	}

	permissionIDs := make(map[string]uint)
	for _, key := range permissions.AllTenantPermissions() {
		module, action, _ := strings.Cut(key, ".")
		permission := models.Permission{Key: key, Module: module, Action: action, IsActive: true}
		if err := tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "key"}},
			DoUpdates: clause.Assignments(map[string]interface{}{
				"module": module, "action": action, "is_active": true, "deleted_at": nil,
			}),
		}).Create(&permission).Error; err != nil {
			return err
		}
		var persisted models.Permission
		if err := tx.Where("key = ?", key).First(&persisted).Error; err != nil {
			return err
		}
		permissionIDs[key] = persisted.ID
	}

	for _, definition := range permissions.TenantRoleDefinitions() {
		var role models.TenantRole
		if err := tx.Where("key = ?", definition.Key).First(&role).Error; err != nil {
			return err
		}
		if err := tx.Where("tenant_role_id = ?", role.ID).Delete(&models.TenantRolePermission{}).Error; err != nil {
			return err
		}
		for _, key := range definition.Permissions {
			mapping := models.TenantRolePermission{TenantRoleID: role.ID, PermissionID: permissionIDs[key]}
			if err := tx.Create(&mapping).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

// Run ejecuta migraciones de datos versionadas después de AutoMigrate.
func Run(db *gorm.DB) error {
	return run(db, dataMigrations)
}

// RunBeforeSchema migrates legacy landlord memberships before Membership is
// changed to require a tenant company_id.
func RunBeforeSchema(db *gorm.DB) error {
	return run(db, beforeSchemaMigrations)
}

func RunModule(db *gorm.DB, module string) error {
	return run(db, migrationsForModule(dataMigrations, module))
}

func RunBeforeSchemaModule(db *gorm.DB, module string) error {
	return run(db, migrationsForModule(beforeSchemaMigrations, module))
}

func migrationsForModule(items []migration, module string) []migration {
	result := make([]migration, 0)
	for _, item := range items {
		if item.module == "core" || item.module == module {
			result = append(result, item)
		}
	}
	return result
}

func run(db *gorm.DB, items []migration) error {
	for _, item := range items {
		if err := db.Transaction(func(tx *gorm.DB) error {
			var count int64
			if err := tx.Model(&models.SchemaMigration{}).Where("id = ?", item.id).Count(&count).Error; err != nil {
				return err
			}
			if count > 0 {
				return nil
			}
			if err := item.run(tx); err != nil {
				return err
			}
			return tx.Create(&models.SchemaMigration{ID: item.id, AppliedAt: time.Now().UTC()}).Error
		}); err != nil {
			return fmt.Errorf("data migration %s failed: %w", item.id, err)
		}
	}
	return nil
}
