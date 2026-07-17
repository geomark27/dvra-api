package transport

import (
	"dvra-api/internal/app/services"
	"dvra-api/internal/modules/staffing/service"
	"dvra-api/internal/shared/middleware"
	"dvra-api/internal/shared/permissions"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes monta las rutas del módulo bajo el grupo protegido (rg).
// Todo el módulo está gateado por plan vía RequireFeature("staffing").
func RegisterRoutes(
	rg *gin.RouterGroup,
	clientSvc *service.StaffingClientService,
	placementSvc *service.PlacementService,
	planService services.PlanService,
	rbac services.TenantAuthorizer,
) {
	clientH := NewStaffingClientHandler(clientSvc)
	placementH := NewPlacementHandler(placementSvc)

	clients := rg.Group("/staffing-clients")
	clients.Use(middleware.RequireFeature(planService, "staffing"))
	{
		clients.GET("", middleware.RequirePermission(rbac, permissions.StaffingClientsView), clientH.GetStaffingClients)
		clients.POST("", middleware.RequirePermission(rbac, permissions.StaffingClientsCreate), clientH.CreateStaffingClient)
		clients.GET("/:id", middleware.RequirePermission(rbac, permissions.StaffingClientsView), clientH.GetStaffingClient)
		clients.PUT("/:id", middleware.RequirePermission(rbac, permissions.StaffingClientsUpdate), clientH.UpdateStaffingClient)
		clients.DELETE("/:id", middleware.RequirePermission(rbac, permissions.StaffingClientsDelete), clientH.DeleteStaffingClient)
	}

	placements := rg.Group("/placements")
	placements.Use(middleware.RequireFeature(planService, "staffing"))
	{
		placements.GET("", middleware.RequirePermission(rbac, permissions.PlacementsView), placementH.GetPlacements)
		placements.POST("", middleware.RequirePermission(rbac, permissions.PlacementsCreate), placementH.CreatePlacement)
		placements.GET("/:id", middleware.RequirePermission(rbac, permissions.PlacementsView), placementH.GetPlacement)
		placements.PUT("/:id", middleware.RequirePermission(rbac, permissions.PlacementsUpdate), placementH.UpdatePlacement)
		placements.DELETE("/:id", middleware.RequirePermission(rbac, permissions.PlacementsDelete), placementH.DeletePlacement)
	}
}
