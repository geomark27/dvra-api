package transport

import (
	"dvra-api/internal/shared/authctx"
	"net/http"
	"strconv"

	"dvra-api/internal/modules/iam/service"
	"dvra-api/internal/shared/apperr"

	"github.com/geomark27/loom-go/pkg/helpers"
	"github.com/gin-gonic/gin"
)

type MembershipHandler struct {
	membershipService service.MembershipService
	logger            helpers.Logger
}

func NewMembershipHandler(membershipService service.MembershipService) *MembershipHandler {
	return &MembershipHandler{
		membershipService: membershipService,
		logger:            helpers.NewLogger(),
	}
}

// GetMemberships godoc
// @Summary      Listar membresías
// @Description  Obtiene las membresías autorizadas por el contexto activo
// @Tags         Memberships
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      403  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /memberships [get]
// @Router       /platform/memberships [get]
func (h *MembershipHandler) GetMemberships(c *gin.Context) {

	// SuperAdmin puede ver todas las memberships
	if authctx.IsSuperAdmin(c) {
		memberships, err := h.membershipService.GetAllMemberships()
		if err != nil {
			h.logger.Error("Failed to get memberships", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve memberships"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"memberships": memberships, "count": len(memberships)}})
		return
	}

	// Usuarios normales solo ven memberships de su empresa
	companyIDVal, exists := c.Get("company_id")
	if !exists {
		c.JSON(http.StatusForbidden, gin.H{"error": "No company context"})
		return
	}

	companyID := companyIDVal.(uint)
	memberships, err := h.membershipService.GetMembershipsByCompanyID(companyID)
	if err != nil {
		h.logger.Error("Failed to get memberships", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve memberships"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"memberships": memberships, "count": len(memberships)}})
}

// GetMembership godoc
// @Summary      Obtener membresía
// @Tags         Memberships
// @Produce      json
// @Param        id path int true "ID de la membresía"
// @Success      200 {object} map[string]interface{}
// @Security     BearerAuth
// @Router       /memberships/{id} [get]
// @Router       /platform/memberships/{id} [get]
func (h *MembershipHandler) GetMembership(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(400, gin.H{"error": "ID inválido"})
		return
	}

	membership, err := h.membershipService.GetMembershipByID(uint(id))
	if err != nil {
		c.JSON(apperr.StatusCode(err), gin.H{"error": err.Error()})
		return
	}

	// Validar acceso: SuperAdmin o miembro de la misma empresa
	if !authctx.IsSuperAdmin(c) {
		companyIDVal, exists := c.Get("company_id")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "No company context"})
			return
		}
		companyID := companyIDVal.(uint)
		if membership.CompanyID != companyID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": membership})
}

// CreateMembership godoc
// @Summary      Crear membresía
// @Description  Asigna un usuario a una empresa desde el contexto landlord
// @Tags         Memberships
// @Accept       json
// @Produce      json
// @Param        membership  body      service.CreateMembershipDTO  true  "Datos de la membresía"
// @Success      201         {object}  map[string]interface{}
// @Failure      400         {object}  map[string]interface{}
// @Failure      403         {object}  map[string]interface{}
// @Failure      500         {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /platform/memberships [post]
func (h *MembershipHandler) CreateMembership(c *gin.Context) {
	// SOLO SuperAdmin puede crear memberships (MVP)
	// Los clientes crean usuarios que automáticamente se agregan a su empresa
	// Sistema de invitaciones → Fase 2
	if !authctx.IsSuperAdmin(c) {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Only superadmin can assign users to companies. Regular users should create new users instead.",
		})
		return
	}

	var dto service.CreateMembershipDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// SuperAdmin debe especificar company_id explícitamente
	if dto.CompanyID == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "company_id is required"})
		return
	}

	membership, err := h.membershipService.CreateMembership(dto)
	if err != nil {
		c.JSON(apperr.StatusCode(err), gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"status": "success", "data": membership})
}

// UpdateMembership godoc
// @Summary      Actualizar rol o estado de una membresía
// @Tags         Memberships
// @Accept       json
// @Produce      json
// @Param        id path int true "ID de la membresía"
// @Param        membership body service.UpdateMembershipDTO true "Cambios"
// @Success      200 {object} map[string]interface{}
// @Security     BearerAuth
// @Router       /memberships/{id} [put]
// @Router       /platform/memberships/{id} [put]
func (h *MembershipHandler) UpdateMembership(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(400, gin.H{"error": "ID inválido"})
		return
	}

	global := authctx.IsSuperAdmin(c)
	var actorCompanyID uint
	if !global {
		var ok bool
		actorCompanyID, ok = authctx.CompanyID(c)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "No company context"})
			return
		}
	}
	actorUserID, _ := authctx.UserID(c)

	var dto service.UpdateMembershipDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	membership, err := h.membershipService.UpdateMembership(uint(id), actorUserID, actorCompanyID, global, dto)
	if err != nil {
		c.JSON(apperr.StatusCode(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": membership})
}

// DeleteMembership godoc
// @Summary      Eliminar una membresía
// @Tags         Memberships
// @Produce      json
// @Param        id path int true "ID de la membresía"
// @Success      200 {object} map[string]interface{}
// @Security     BearerAuth
// @Router       /memberships/{id} [delete]
// @Router       /platform/memberships/{id} [delete]
func (h *MembershipHandler) DeleteMembership(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(400, gin.H{"error": "ID inválido"})
		return
	}

	// Validar que el membership pertenece a la empresa del usuario
	if !authctx.IsSuperAdmin(c) {
		membership, err := h.membershipService.GetMembershipByID(uint(id))
		if err != nil {
			c.JSON(apperr.StatusCode(err), gin.H{"error": err.Error()})
			return
		}
		companyIDVal, exists := c.Get("company_id")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "No company context"})
			return
		}
		companyID := companyIDVal.(uint)
		if membership.CompanyID != companyID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}
	}

	if err := h.membershipService.DeleteMembership(uint(id)); err != nil {
		c.JSON(apperr.StatusCode(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Membership deleted"})
}
