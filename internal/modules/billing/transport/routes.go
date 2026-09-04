package transport

import (
	"dvra-api/internal/modules/billing/service"
	"dvra-api/internal/shared/middleware"
	"dvra-api/internal/shared/permissions"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes monta billing en los tres contextos:
//
//   - público: el catálogo de planes que pinta la página de precios, sin sesión.
//   - tenant: la empresa lee y edita su configuración, y consulta su
//     suscripción y su consumo frente a los límites del plan.
//   - landlord: solo el SuperAdmin administra los planes y los asigna. Un tenant
//     no puede cambiarse de plan ni tocar sus fechas de prueba (ADR-003).
func RegisterRoutes(public, protected, platform *gin.RouterGroup, plans service.PlanService, settings *service.CompanySettingsService, rbac middleware.PermissionChecker) {
	planH := NewPlanHandler(plans)
	settingsH := NewCompanySettingsHandler(settings)

	publicPlans := public.Group("/plans")
	{
		publicPlans.GET("", planH.GetPublicPlans)
		publicPlans.GET("/:slug", planH.GetPlanBySlug)
	}

	protected.GET("/settings", middleware.RequirePermission(rbac, permissions.SettingsView), settingsH.Get)
	protected.PUT("/settings", middleware.RequirePermission(rbac, permissions.SettingsUpdate), settingsH.Update)
	protected.GET("/subscription", middleware.RequirePermission(rbac, permissions.BillingView), settingsH.Subscription)
	protected.GET("/subscription/usage", middleware.RequirePermission(rbac, permissions.BillingView), settingsH.Usage)

	platform.GET("/plans", middleware.RequirePlatformPermission(permissions.PlatformPlansManage), planH.GetPlans)
	platform.POST("/plans", middleware.RequirePlatformPermission(permissions.PlatformPlansManage), planH.CreatePlan)
	platform.GET("/plans/:id", middleware.RequirePlatformPermission(permissions.PlatformPlansManage), planH.GetPlanByID)
	platform.PUT("/plans/:id", middleware.RequirePlatformPermission(permissions.PlatformPlansManage), planH.UpdatePlan)
	platform.PATCH("/plans/:id/toggle", middleware.RequirePlatformPermission(permissions.PlatformPlansManage), planH.TogglePlanStatus)
	platform.DELETE("/plans/:id", middleware.RequirePlatformPermission(permissions.PlatformPlansManage), planH.DeletePlan)
	platform.POST("/plans/assign", middleware.RequirePlatformPermission(permissions.PlatformPlansManage), planH.AssignPlanToCompany)
}
