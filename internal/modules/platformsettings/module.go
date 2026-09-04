// Package platformsettings es el punto de ensamblaje de la configuración global
// del SaaS: marca, contacto, URLs, textos legales y valores por defecto de
// negocio (días de prueba, plan inicial).
//
// Es un singleton: una sola fila. Se llama platformsettings y no platform para
// no confundirse con internal/platform/, que es infraestructura (config,
// servidor, almacenamiento), no un dominio de negocio.
package platformsettings

import (
	"dvra-api/internal/modules/platformsettings/repository"
	"dvra-api/internal/modules/platformsettings/service"
	"dvra-api/internal/modules/platformsettings/transport"
	"dvra-api/internal/platform/httpx"

	"gorm.io/gorm"
)

// Module agrupa las dependencias ya cableadas del módulo.
type Module struct {
	svc *service.PlatformSettingsService
}

// New construye el módulo con la conexión inyectada.
func New(db *gorm.DB) *Module {
	return &Module{
		svc: service.NewPlatformSettingsService(repository.NewPlatformSettingsRepository(db)),
	}
}

// Settings expone el service para los consumidores internos que necesitan los
// valores por defecto de negocio (días de prueba, plan inicial). Otros módulos
// deben declarar su propio puerto mínimo en vez de depender de este tipo.
func (m *Module) Settings() *service.PlatformSettingsService { return m.svc }

// Name identifica al módulo en logs y en la auditoría de rutas.
func (m *Module) Name() string { return "platformsettings" }

// RegisterRoutes monta la ruta pública y las de landlord.
func (m *Module) RegisterRoutes(r httpx.Router) {
	transport.RegisterRoutes(r.Public, r.Platform, m.svc)
}
