// Package billing es el punto de ensamblaje de planes, suscripciones y
// configuración de empresa.
//
// Es el mayor proveedor del sistema: de él dependen el gating por plan
// (RequireFeature), el límite de vacantes publicadas de recruitment, y el alta
// de empresa de iam.
package billing

import (
	"dvra-api/internal/modules/billing/repository"
	"dvra-api/internal/modules/billing/service"
	"dvra-api/internal/modules/billing/transport"
	"dvra-api/internal/platform/httpx"

	"gorm.io/gorm"
)

// Module agrupa las dependencias ya cableadas del módulo.
type Module struct {
	plans    service.PlanService
	settings *service.CompanySettingsService
}

// New construye el módulo. companies es el puerto hacia iam, que el composition
// root satisface con el repositorio de empresas.
func New(db *gorm.DB, companies service.CompanyReader) *Module {
	return &Module{
		plans: service.NewPlanService(
			repository.NewPlanRepository(db),
			companies,
			repository.NewSubscriptionRepository(db),
			db,
		),
		settings: service.NewCompanySettingsService(db),
	}
}

// Plans expone el service de planes.
//
// Satisface directamente middleware.FeatureChecker (CompanyHasFeature), que es
// como el gating por plan llega a las rutas sin que shared/middleware conozca
// este módulo.
func (m *Module) Plans() service.PlanService { return m.plans }

// Name identifica al módulo en logs y en la auditoría de rutas.
func (m *Module) Name() string { return "billing" }

// RegisterRoutes monta las rutas públicas, de tenant y de landlord.
func (m *Module) RegisterRoutes(r httpx.Router) {
	transport.RegisterRoutes(r.Public, r.Protected, r.Platform, m.plans, m.settings, r.RBAC)
}
