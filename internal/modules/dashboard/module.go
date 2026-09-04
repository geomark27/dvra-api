// Package dashboard es el punto de ensamblaje de las métricas agregadas que ve
// una empresa en su pantalla de inicio.
package dashboard

import (
	"dvra-api/internal/modules/dashboard/repository"
	"dvra-api/internal/modules/dashboard/service"
	"dvra-api/internal/modules/dashboard/transport"
	"dvra-api/internal/platform/httpx"

	"gorm.io/gorm"
)

// Module agrupa las dependencias ya cableadas del módulo.
type Module struct {
	svc service.DashboardService
}

// New construye el módulo con la conexión inyectada.
func New(db *gorm.DB) *Module {
	return &Module{svc: service.NewDashboardService(repository.NewDashboardRepository(db))}
}

// Name identifica al módulo en logs y en la auditoría de rutas.
func (m *Module) Name() string { return "dashboard" }

// RegisterRoutes monta las rutas de tenant.
func (m *Module) RegisterRoutes(r httpx.Router) {
	transport.RegisterRoutes(r.Protected, m.svc, r.RBAC)
}
