// Package catalog es el punto de ensamblaje del catálogo de valores del sistema
// (ADR-003): jerarquía módulo → categoría → valor, con valores globales, valores
// propios del tenant y overrides de etiqueta por empresa.
//
// Es el catálogo del que dependen los desplegables de toda la aplicación: tipos
// de contrato, modalidades de trabajo, etapas del pipeline, monedas, estados de
// suscripción.
package catalog

import (
	"dvra-api/internal/modules/catalog/repository"
	"dvra-api/internal/modules/catalog/service"
	"dvra-api/internal/modules/catalog/transport"
	"dvra-api/internal/platform/httpx"

	"gorm.io/gorm"
)

// Module agrupa las dependencias ya cableadas del módulo.
type Module struct {
	svc service.SystemValueService
}

// New construye el módulo con la conexión inyectada.
func New(db *gorm.DB) *Module {
	return &Module{
		svc: service.NewSystemValueService(repository.NewSystemValueRepository(db)),
	}
}

// Name identifica al módulo en logs y en la auditoría de rutas.
func (m *Module) Name() string { return "catalog" }

// RegisterRoutes monta las rutas de tenant y las de landlord.
func (m *Module) RegisterRoutes(r httpx.Router) {
	transport.RegisterRoutes(r.Protected, r.Platform, m.svc, r.RBAC)
}
