package migrations

import (
	"fmt"
	"strings"
	"time"

	"dvra-api/internal/app/models"
	"dvra-api/internal/shared/permissions"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type migration struct {
	id  string
	run func(*gorm.DB) error
}

var dataMigrations = []migration{
	{
		id: "20260715_jobs_active_to_published",
		run: func(tx *gorm.DB) error {
			return tx.Model(&models.Job{}).Where("status = ?", "active").Update("status", "published").Error
		},
	},
	{
		id: "20260715_tenant_superadmin_to_admin",
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
		id: "20260717_seed_relational_tenant_rbac",
		run: func(tx *gorm.DB) error {
			return seedTenantRBAC(tx)
		},
	},
	{
		id: "20260717_admin_full_tenant_access",
		run: func(tx *gorm.DB) error {
			return seedTenantRBAC(tx)
		},
	},
	{
		id: "20260720_departments_module",
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
		id: "20260715_split_platform_memberships",
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
		id: "20260715_require_tenant_membership_company",
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
		id: "20260717_relational_tenant_rbac",
		run: func(tx *gorm.DB) error {
			if err := tx.AutoMigrate(&models.TenantRole{}, &models.Permission{}, &models.TenantRolePermission{}); err != nil {
				return err
			}
			if err := seedTenantRBAC(tx); err != nil {
				return err
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
