// Package locations es el punto de ensamblaje del catálogo geográfico
// (región → subregión → país → estado → ciudad).
//
// Es un catálogo de referencia global: no pertenece a ninguna empresa, así que
// no lleva scoping por company_id, y se expone públicamente en modo lectura.
// No consume nada de otros módulos: no tiene ports.go.
package locations

import (
	"dvra-api/internal/modules/locations/repository"
	"dvra-api/internal/modules/locations/service"
	"dvra-api/internal/modules/locations/transport"
	"dvra-api/internal/platform/httpx"

	"gorm.io/gorm"
)

// Module agrupa las dependencias ya cableadas del módulo locations.
type Module struct {
	svc service.LocationService
}

// New construye el módulo con la conexión inyectada.
func New(db *gorm.DB) *Module {
	return &Module{
		svc: service.NewLocationService(repository.NewLocationRepository(db)),
	}
}

// Name identifica al módulo en logs y en la auditoría de rutas.
func (m *Module) Name() string { return "locations" }

// RegisterRoutes monta el catálogo bajo el grupo público.
func (m *Module) RegisterRoutes(r httpx.Router) {
	transport.RegisterRoutes(r.Public, m.svc)
}
