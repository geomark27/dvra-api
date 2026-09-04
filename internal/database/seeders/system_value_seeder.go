package seeders

import (
	"fmt"
	"log"

	"dvra-api/internal/shared/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type CatalogValueDefinition struct {
	Key, Label string
	Order      int
}

type CatalogCategoryDefinition struct {
	Key, Name           string
	IsExtensible        bool
	AllowTenantValues   bool
	AllowTenantOverride bool
	Values              []CatalogValueDefinition
}

type CatalogModuleDefinition struct {
	Key, Name  string
	Categories []CatalogCategoryDefinition
}

var CatalogModules = []CatalogModuleDefinition{
	{Key: "identity", Name: "Identidad y acceso", Categories: []CatalogCategoryDefinition{
		{Key: "membership_status", Name: "Estado de membresía", Values: values("active", "Activo", "inactive", "Inactivo", "pending", "Pendiente", "suspended", "Suspendido", "removed", "Eliminado")},
		{Key: "platform_membership_status", Name: "Estado de acceso landlord", Values: values("active", "Activo", "inactive", "Inactivo", "suspended", "Suspendido")},
	}},
	{Key: "recruitment", Name: "Reclutamiento", Categories: []CatalogCategoryDefinition{
		{Key: "job_status", Name: "Estado de vacante", AllowTenantOverride: true, Values: values("draft", "Borrador", "published", "Publicada", "closed", "Cerrada")},
		{Key: "application_stage", Name: "Etapa del pipeline", IsExtensible: true, AllowTenantValues: true, AllowTenantOverride: true, Values: values("applied", "Aplicado", "screening", "En revisión", "technical", "Prueba técnica", "interview", "Entrevista", "offer", "Oferta", "hired", "Contratado", "rejected", "Rechazado")},
		{Key: "work_mode", Name: "Modalidad de trabajo", AllowTenantOverride: true, Values: values("remote", "Remoto", "onsite", "Presencial", "hybrid", "Híbrido")},
		{Key: "contract_type", Name: "Tipo de contrato", IsExtensible: true, AllowTenantValues: true, AllowTenantOverride: true, Values: values("full_time", "Tiempo completo", "part_time", "Medio tiempo", "contractor", "Contratista", "internship", "Pasantía", "temporary", "Temporal")},
		{Key: "experience_level", Name: "Nivel de experiencia", IsExtensible: true, AllowTenantValues: true, AllowTenantOverride: true, Values: values("entry", "Junior / Sin experiencia", "mid", "Semi-Senior", "senior", "Senior", "lead", "Lead / Principal")},
		{Key: "priority", Name: "Prioridad", AllowTenantOverride: true, Values: values("low", "Baja", "medium", "Media", "high", "Alta", "urgent", "Urgente")},
		{Key: "candidate_source", Name: "Fuente del candidato", IsExtensible: true, AllowTenantValues: true, AllowTenantOverride: true, Values: values("linkedin", "LinkedIn", "website", "Sitio web", "referral", "Referido", "job_board", "Bolsa de trabajo", "direct", "Aplicación directa", "agency", "Agencia", "career_page", "Career page", "other", "Otro")},
	}},
	{Key: "staffing", Name: "Staffing", Categories: []CatalogCategoryDefinition{
		{Key: "client_status", Name: "Estado de cliente", Values: values("active", "Activo", "inactive", "Inactivo", "prospect", "Prospecto")},
		{Key: "placement_status", Name: "Estado de colocación", Values: values("active", "Activo", "ended", "Finalizada", "suspended", "Suspendida")},
		{Key: "placement_contract_type", Name: "Tipo de contrato de colocación", IsExtensible: true, AllowTenantValues: true, AllowTenantOverride: true, Values: values("outsourcing", "Outsourcing", "staffing", "Staffing", "project", "Proyecto")},
		{Key: "billing_rate_type", Name: "Periodicidad de tarifa", Values: values("hourly", "Por hora", "monthly", "Mensual")},
	}},
	{Key: "billing", Name: "Planes y facturación", Categories: []CatalogCategoryDefinition{
		{Key: "subscription_status", Name: "Estado de suscripción", Values: values("trialing", "En prueba", "active", "Activa", "past_due", "Pago vencido", "paused", "Pausada", "canceled", "Cancelada", "expired", "Expirada")},
		{Key: "billing_cycle", Name: "Ciclo de facturación", Values: values("monthly", "Mensual", "yearly", "Anual")},
		{Key: "currency", Name: "Moneda", Values: values("USD", "Dólar estadounidense", "EUR", "Euro")},
		{Key: "support_level", Name: "Nivel de soporte", Values: values("email", "Correo", "priority", "Prioritario", "dedicated", "Dedicado")},
	}},
}

func values(items ...string) []CatalogValueDefinition {
	result := make([]CatalogValueDefinition, 0, len(items)/2)
	for i := 0; i+1 < len(items); i += 2 {
		result = append(result, CatalogValueDefinition{Key: items[i], Label: items[i+1], Order: len(result) + 1})
	}
	return result
}

type SystemValueSeeder struct{ Module string }

func (s *SystemValueSeeder) Run(db *gorm.DB) error { return SeedCatalogs(db, s.Module) }

func SystemValueSeederFunc(db *gorm.DB) error { return SeedCatalogs(db, "") }

// SeedCatalogs upserts release-controlled definitions. module="" seeds all.
func SeedCatalogs(db *gorm.DB, moduleFilter string) error {
	log.Printf("🌱 Seeding catalogs (module=%q)...", moduleFilter)
	if err := SeedCatalogStructure(db, moduleFilter); err != nil {
		return err
	}
	for _, definition := range CatalogModules {
		if moduleFilter != "" && definition.Key != moduleFilter {
			continue
		}
		var module models.SystemValueModule
		if err := db.Where("key = ?", definition.Key).First(&module).Error; err != nil {
			return err
		}
		for _, categoryDefinition := range definition.Categories {
			var category models.SystemValueCategory
			if err := db.Where("module_id = ? AND key = ?", module.ID, categoryDefinition.Key).First(&category).Error; err != nil {
				return err
			}
			for _, valueDefinition := range categoryDefinition.Values {
				value := models.SystemValue{CategoryID: category.ID, Key: valueDefinition.Key, Label: valueDefinition.Label, DisplayOrder: valueDefinition.Order, IsSystem: true, IsActive: true}
				if err := db.Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "category_id"}, {Name: "key"}}, TargetWhere: clause.Where{Exprs: []clause.Expression{clause.Expr{SQL: "company_id IS NULL AND deleted_at IS NULL"}}}, DoUpdates: clause.Assignments(map[string]interface{}{"label": value.Label, "display_order": value.DisplayOrder, "is_system": true, "is_active": true, "deleted_at": nil})}).Create(&value).Error; err != nil {
					return fmt.Errorf("seed %s.%s.%s: %w", module.Key, category.Key, value.Key, err)
				}
			}
		}
	}
	return nil
}

// SeedCatalogStructure can run while migrating the legacy system_values table,
// whose old NOT NULL columns still prevent inserting final value rows.
func SeedCatalogStructure(db *gorm.DB, moduleFilter string) error {
	for _, definition := range CatalogModules {
		if moduleFilter != "" && definition.Key != moduleFilter {
			continue
		}
		module := models.SystemValueModule{Key: definition.Key, Name: definition.Name, IsSystem: true, IsActive: true}
		if err := db.Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "key"}}, DoUpdates: clause.Assignments(map[string]interface{}{"name": definition.Name, "is_system": true, "is_active": true, "deleted_at": nil})}).Create(&module).Error; err != nil {
			return err
		}
		if err := db.Where("key = ?", definition.Key).First(&module).Error; err != nil {
			return err
		}
		for _, item := range definition.Categories {
			category := models.SystemValueCategory{ModuleID: module.ID, Key: item.Key, Name: item.Name, IsSystem: true, IsExtensible: item.IsExtensible, AllowTenantValues: item.AllowTenantValues, AllowTenantOverride: item.AllowTenantOverride, IsActive: true}
			updates := map[string]interface{}{"name": category.Name, "is_system": true, "is_extensible": category.IsExtensible, "allow_tenant_values": category.AllowTenantValues, "allow_tenant_override": category.AllowTenantOverride, "is_active": true, "deleted_at": nil}
			if err := db.Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "module_id"}, {Name: "key"}}, DoUpdates: clause.Assignments(updates)}).Create(&category).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

// CatalogOwner resolves old flat category names during an in-place migration.
func CatalogOwner(category string) (string, string, bool) {
	if category == "application_status" {
		category = "application_stage"
	}
	for _, module := range CatalogModules {
		for _, item := range module.Categories {
			if item.Key == category {
				return module.Key, item.Key, true
			}
		}
	}
	return "", "", false
}
