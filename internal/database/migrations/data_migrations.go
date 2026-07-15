package migrations

import (
	"fmt"
	"time"

	"dvra-api/internal/app/models"

	"gorm.io/gorm"
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
			return tx.Model(&models.Membership{}).
				Where("company_id IS NOT NULL AND role = ?", "superadmin").
				Update("role", "admin").Error
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
