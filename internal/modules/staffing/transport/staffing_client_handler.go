// Package transport es el adaptador de entrada HTTP del módulo staffing:
// handlers Gin + registro de rutas. Mapea HTTP <-> casos de uso (service).
package transport

import (
	"dvra-api/internal/modules/staffing/domain"
	"net/http"
	"strconv"

	"dvra-api/internal/modules/staffing/service"
	"dvra-api/internal/shared/apperr"
	"dvra-api/internal/shared/authctx"

	"github.com/gin-gonic/gin"
)

type StaffingClientHandler struct {
	svc *service.StaffingClientService
}

func NewStaffingClientHandler(svc *service.StaffingClientService) *StaffingClientHandler {
	return &StaffingClientHandler{svc: svc}
}

// GetStaffingClients godoc
// @Summary      Listar clientes finales (staffing)
// @Description  Lista los clientes finales de la empresa del token. Requiere un plan con el módulo staffing habilitado.
// @Tags         StaffingClients
// @Produce      json
// @Param        status  query     string  false  "Filtrar por estado (active, inactive, prospect)"
// @Success      200     {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /staffing-clients [get]
func (h *StaffingClientHandler) GetStaffingClients(c *gin.Context) {
	var filters domain.StaffingClientFilters
	if err := c.ShouldBindQuery(&filters); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	companyID, ok := authctx.CompanyID(c)
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "No company context"})
		return
	}

	clients, err := h.svc.GetByCompanyID(companyID, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve staffing clients"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"staffing_clients": service.ToStaffingClientResponseList(clients),
			"count":            len(clients),
		},
	})
}

// GetStaffingClient godoc
// @Summary      Obtener cliente final
// @Tags         StaffingClients
// @Produce      json
// @Param        id   path      int  true  "ID del cliente final"
// @Success      200  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /staffing-clients/{id} [get]
func (h *StaffingClientHandler) GetStaffingClient(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	scope, ok := authctx.TenantScope(c)
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "No company context"})
		return
	}
	client, err := h.svc.GetByID(uint(id), scope)
	if err != nil {
		c.JSON(apperr.StatusCode(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": service.ToStaffingClientResponse(client)})
}

// CreateStaffingClient godoc
// @Summary      Crear cliente final
// @Tags         StaffingClients
// @Accept       json
// @Produce      json
// @Param        staffing_client  body      service.CreateStaffingClientDTO  true  "Datos del cliente final"
// @Success      201              {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /staffing-clients [post]
func (h *StaffingClientHandler) CreateStaffingClient(c *gin.Context) {
	var dto service.CreateStaffingClientDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !authctx.IsSuperAdmin(c) {
		companyID, ok := authctx.CompanyID(c)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "No company context"})
			return
		}
		dto.CompanyID = companyID
	}
	if dto.CompanyID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_id is required"})
		return
	}

	client, err := h.svc.Create(dto)
	if err != nil {
		c.JSON(apperr.StatusCode(err), gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"status": "success", "data": service.ToStaffingClientResponse(client)})
}

// UpdateStaffingClient godoc
// @Summary      Actualizar cliente final
// @Tags         StaffingClients
// @Accept       json
// @Produce      json
// @Param        id               path      int                           true  "ID del cliente final"
// @Param        staffing_client  body      service.UpdateStaffingClientDTO  true  "Campos a actualizar"
// @Success      200              {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /staffing-clients/{id} [put]
func (h *StaffingClientHandler) UpdateStaffingClient(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	scope, ok := authctx.TenantScope(c)
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "No company context"})
		return
	}

	var dto service.UpdateStaffingClientDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updated, err := h.svc.Update(uint(id), scope, dto)
	if err != nil {
		c.JSON(apperr.StatusCode(err), gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": service.ToStaffingClientResponse(updated)})
}

// DeleteStaffingClient godoc
// @Summary      Eliminar cliente final
// @Tags         StaffingClients
// @Produce      json
// @Param        id   path      int  true  "ID del cliente final"
// @Success      200  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /staffing-clients/{id} [delete]
func (h *StaffingClientHandler) DeleteStaffingClient(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	scope, ok := authctx.TenantScope(c)
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "No company context"})
		return
	}

	if err := h.svc.Delete(uint(id), scope); err != nil {
		c.JSON(apperr.StatusCode(err), gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Staffing client deleted"})
}
