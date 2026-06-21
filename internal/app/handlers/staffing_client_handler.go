package handlers

import (
	"net/http"
	"strconv"

	"dvra-api/internal/app/dtos"
	"dvra-api/internal/app/services"
	"dvra-api/internal/shared/apperr"
	"dvra-api/internal/shared/authctx"

	"github.com/gin-gonic/gin"
)

type StaffingClientHandler struct {
	service services.StaffingClientService
}

func NewStaffingClientHandler(service services.StaffingClientService) *StaffingClientHandler {
	return &StaffingClientHandler{service: service}
}

// GetStaffingClients godoc
// @Summary      Listar clientes finales (staffing)
// @Description  Lista los clientes finales de la empresa del token. Requiere un plan con el módulo staffing habilitado.
// @Tags         StaffingClients
// @Accept       json
// @Produce      json
// @Param        status  query     string  false  "Filtrar por estado (active, inactive, prospect)"
// @Success      200     {object}  map[string]interface{}
// @Failure      403     {object}  map[string]interface{}
// @Failure      500     {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /staffing-clients [get]
func (h *StaffingClientHandler) GetStaffingClients(c *gin.Context) {
	var filters dtos.StaffingClientFilters
	if err := c.ShouldBindQuery(&filters); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// SuperAdmin no tiene empresa propia: este listado es por tenant.
	companyID, ok := authctx.CompanyID(c)
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "No company context"})
		return
	}

	clients, err := h.service.GetByCompanyID(companyID, filters)
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
// @Description  Obtiene un cliente final por ID (validado contra la empresa del token)
// @Tags         StaffingClients
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID del cliente final"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      403  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /staffing-clients/{id} [get]
func (h *StaffingClientHandler) GetStaffingClient(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	client, err := h.service.GetByID(uint(id))
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
// @Description  Crea un cliente final dentro de la empresa del token (company_id se fuerza del token). Requiere plan con módulo staffing.
// @Tags         StaffingClients
// @Accept       json
// @Produce      json
// @Param        staffing_client  body      dtos.CreateStaffingClientDTO  true  "Datos del cliente final"
// @Success      201              {object}  map[string]interface{}
// @Failure      400              {object}  map[string]interface{}
// @Failure      403              {object}  map[string]interface{}
// @Failure      409              {object}  map[string]interface{}
// @Failure      500              {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /staffing-clients [post]
func (h *StaffingClientHandler) CreateStaffingClient(c *gin.Context) {
	var dto dtos.CreateStaffingClientDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Forzar company_id del token para usuarios normales (evita escalada cross-tenant)
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

	client, err := h.service.Create(dto)
	if err != nil {
		c.JSON(apperr.StatusCode(err), gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"status": "success", "data": dtos.ToStaffingClientResponse(client)})
}

// UpdateStaffingClient godoc
// @Summary      Actualizar cliente final
// @Description  Actualiza parcialmente un cliente final (validado contra la empresa del token)
// @Tags         StaffingClients
// @Accept       json
// @Produce      json
// @Param        id               path      int                           true  "ID del cliente final"
// @Param        staffing_client  body      dtos.UpdateStaffingClientDTO  true  "Campos a actualizar"
// @Success      200              {object}  map[string]interface{}
// @Failure      400              {object}  map[string]interface{}
// @Failure      403              {object}  map[string]interface{}
// @Failure      404              {object}  map[string]interface{}
// @Failure      409              {object}  map[string]interface{}
// @Failure      500              {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /staffing-clients/{id} [put]
func (h *StaffingClientHandler) UpdateStaffingClient(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	// companyID = 0 → SuperAdmin (el service omite la validación de tenant)
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

	updated, err := h.service.Update(uint(id), companyID, dto)
	if err != nil {
		c.JSON(apperr.StatusCode(err), gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": dtos.ToStaffingClientResponse(updated)})
}

// DeleteStaffingClient godoc
// @Summary      Eliminar cliente final
// @Description  Elimina (soft delete) un cliente final (validado contra la empresa del token)
// @Tags         StaffingClients
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID del cliente final"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      403  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
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

	if err := h.service.Delete(uint(id), companyID); err != nil {
		c.JSON(apperr.StatusCode(err), gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Staffing client deleted"})
}
