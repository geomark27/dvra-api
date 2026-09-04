package transport

import (
	"dvra-api/internal/modules/catalog/service"
	"dvra-api/internal/shared/middleware"
	"dvra-api/internal/shared/permissions"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes monta el catálogo en sus dos contextos (ver ADR-003):
//
//   - tenant: la empresa lee los valores que le aplican y, en las categorías
//     extensibles, añade los suyos o sobrescribe la etiqueta de uno global.
//   - landlord: el SuperAdmin administra el catálogo global, el que heredan
//     todas las empresas.
//
// El aislamiento es el de siempre: un tenant nunca ve ni toca los valores de
// otro; sus altas y overrides llevan su company_id, forzado desde el token.
func RegisterRoutes(protected, platform *gin.RouterGroup, svc service.SystemValueService, rbac middleware.PermissionChecker) {
	h := NewSystemValueHandler(svc)

	systemValues := protected.Group("/system-values")
	{
		systemValues.GET("", middleware.RequirePermission(rbac, permissions.CatalogsView), h.GetCategories)
		// Deprecada: no distingue módulo, así que dos módulos con la misma
		// categoría colisionan. Se conserva por compatibilidad del frontend.
		systemValues.GET("/:category", middleware.RequirePermission(rbac, permissions.CatalogsView), h.GetByCategory)
		systemValues.GET("/modules/:module/:category", middleware.RequirePermission(rbac, permissions.CatalogsView), h.GetByModuleCategory)
		systemValues.POST("", middleware.RequirePermission(rbac, permissions.CatalogsManage), h.CreateTenant)
		systemValues.PUT("/:id", middleware.RequirePermission(rbac, permissions.CatalogsManage), h.UpdateTenant)
		systemValues.DELETE("/:id", middleware.RequirePermission(rbac, permissions.CatalogsManage), h.DeleteTenant)
		systemValues.PUT("/:id/override", middleware.RequirePermission(rbac, permissions.CatalogsManage), h.UpsertTenantOverride)
	}

	platform.GET("/system-values", middleware.RequirePlatformPermission(permissions.PlatformCatalogsManage), h.GetAll)
	platform.POST("/system-values", middleware.RequirePlatformPermission(permissions.PlatformCatalogsManage), h.Create)
	platform.PUT("/system-values/:id", middleware.RequirePlatformPermission(permissions.PlatformCatalogsManage), h.Update)
	platform.DELETE("/system-values/:id", middleware.RequirePlatformPermission(permissions.PlatformCatalogsManage), h.Delete)
}
