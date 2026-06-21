// Package transport es el adaptador de entrada HTTP del módulo staffing:
// handlers Gin + registro de rutas. Mapea HTTP <-> casos de uso (service).
package transport

import (
	"net/http"
	"strconv"

	"dvra-api/internal/app/dtos"
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
	var filters dtos.StaffingClientFilters
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
			"staffing_clients": dtos.ToStaffingClientResponseList(clients),
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

	client, err := h.svc.GetByID(uint(id))
	if err != nil {
		c.JSON(apperr.StatusCode(err), gin.H{"error": err.Error()})
		return
	}

	if !authctx.IsSuperAdmin(c) {
		companyID, ok := authctx.CompanyID(c)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "No company context"})
			return
		}
		if client.CompanyID != companyID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": dtos.ToStaffingClientResponse(client)})
}

// CreateStaffingClient godoc
// @Summary      Crear cliente final
// @Tags         StaffingClients
// @Accept       json
// @Produce      json
// @Param        staffing_client  body      dtos.CreateStaffingClientDTO  true  "Datos del cliente final"
// @Success      201              {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /staffing-clients [post]
func (h *StaffingClientHandler) CreateStaffingClient(c *gin.Context) {
	var dto dtos.CreateStaffingClientDTO
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
	c.JSON(http.StatusCreated, gin.H{"status": "success", "data": dtos.ToStaffingClientResponse(client)})
}

// UpdateStaffingClient godoc
// @Summary      Actualizar cliente final
// @Tags         StaffingClients
// @Accept       json
// @Produce      json
// @Param        id               path      int                           true  "ID del cliente final"
// @Param        staffing_client  body      dtos.UpdateStaffingClientDTO  true  "Campos a actualizar"
// @Success      200              {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /staffing-clients/{id} [put]
func (h *StaffingClientHandler) UpdateStaffingClient(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var companyID uint
	if !authctx.IsSuperAdmin(c) {
		cid, ok := authctx.CompanyID(c)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "No company context"})
			return
		}
		companyID = cid
	}

	var dto dtos.UpdateStaffingClientDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updated, err := h.svc.Update(uint(id), companyID, dto)
	if err != nil {
		c.JSON(apperr.StatusCode(err), gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": dtos.ToStaffingClientResponse(updated)})
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

	var companyID uint
	if !authctx.IsSuperAdmin(c) {
		cid, ok := authctx.CompanyID(c)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "No company context"})
			return
		}
		companyID = cid
	}

	if err := h.svc.Delete(uint(id), companyID); err != nil {
		c.JSON(apperr.StatusCode(err), gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Staffing client deleted"})
}
