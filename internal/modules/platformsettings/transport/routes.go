package transport

import (
	"dvra-api/internal/modules/platformsettings/service"
	"dvra-api/internal/shared/middleware"
	"dvra-api/internal/shared/permissions"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes monta la configuración global del SaaS en sus dos contextos:
//
//   - pública: el frontend la lee sin sesión para pintar marca, enlaces legales
//     y datos de contacto antes del login.
//   - landlord: solo el SuperAdmin la edita, con permiso de plataforma.
//
// Un tenant no tiene nada que hacer aquí: la configuración es del SaaS, no de
// una empresa (esa vive en company_settings, del módulo de billing).
func RegisterRoutes(public, platform *gin.RouterGroup, svc *service.PlatformSettingsService) {
	h := NewPlatformSettingsHandler(svc)

	public.GET("/public/platform-settings", h.GetPublicSettings)

	platform.GET("/settings",
		middleware.RequirePlatformPermission(permissions.PlatformSettingsManage),
		h.GetFullSettings)
	platform.PUT("/settings",
		middleware.RequirePlatformPermission(permissions.PlatformSettingsManage),
		h.UpdateSettings)
}
