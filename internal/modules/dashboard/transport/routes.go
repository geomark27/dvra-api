package transport

import (
	"dvra-api/internal/modules/dashboard/service"
	"dvra-api/internal/shared/middleware"
	"dvra-api/internal/shared/permissions"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes monta las métricas agregadas de la empresa activa.
func RegisterRoutes(protected *gin.RouterGroup, svc service.DashboardService, rbac middleware.PermissionChecker) {
	h := NewDashboardHandler(svc)

	dashboard := protected.Group("/dashboard")
	{
		dashboard.GET("/stats", middleware.RequirePermission(rbac, permissions.DashboardView), h.GetStats)
	}
}
