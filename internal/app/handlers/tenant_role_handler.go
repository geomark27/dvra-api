package handlers

import (
	"net/http"

	"dvra-api/internal/app/dtos"
	"dvra-api/internal/app/services"

	"github.com/gin-gonic/gin"
)

type TenantRoleHandler struct{ service *services.TenantRBACService }

func NewTenantRoleHandler(service *services.TenantRBACService) *TenantRoleHandler {
	return &TenantRoleHandler{service: service}
}

// GetTenantRoles godoc
// @Summary      Listar roles tenant
// @Description  Devuelve los perfiles activos que pueden asignarse a memberships
// @Tags         Tenant Roles
// @Produce      json
// @Success      200 {object} map[string]interface{}
// @Security     BearerAuth
// @Router       /tenant-roles [get]
func (h *TenantRoleHandler) GetTenantRoles(c *gin.Context) {
	roles, err := h.service.ListRoles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve tenant roles"})
		return
	}
	response := make([]*dtos.TenantRoleResponse, 0, len(roles))
	for i := range roles {
		response = append(response, dtos.ToTenantRoleResponse(&roles[i]))
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": response})
}
