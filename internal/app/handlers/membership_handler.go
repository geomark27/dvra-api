package handlers

import (
	"net/http"
	"strconv"

	"dvra-api/internal/app/dtos"
	"dvra-api/internal/app/services"

	"github.com/geomark27/loom-go/pkg/helpers"
	"github.com/gin-gonic/gin"
)

type MembershipHandler struct {
	membershipService services.MembershipService
	logger            helpers.Logger
}

func NewMembershipHandler(membershipService services.MembershipService) *MembershipHandler {
	return &MembershipHandler{
		membershipService: membershipService,
		logger:            helpers.NewLogger(),
	}
}

// GetMemberships godoc
// @Summary      Listar membresías
// @Description  Obtiene membresías (todas si es SuperAdmin, de la empresa si es usuario normal)
// @Tags         Memberships
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      403  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /memberships [get]
func (h *MembershipHandler) GetMemberships(c *gin.Context) {
	role, _ := c.Get("role")

	// SuperAdmin puede ver todas las memberships
	if role == "superadmin" {
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

func (h *MembershipHandler) GetMembership(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(400, gin.H{"error": "ID inválido"})
		return
	}

	membership, err := h.membershipService.GetMembershipByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	// Validar acceso: SuperAdmin o miembro de la misma empresa
	role, _ := c.Get("role")
	if role != "superadmin" {
		companyIDVal, exists := c.Get("company_id")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "No company context"})
			return
		}
		companyID := companyIDVal.(uint)
		if membership.CompanyID == nil || *membership.CompanyID != companyID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": membership})
}

// CreateMembership godoc
// @Summary      Crear membresía
// @Description  Asigna un usuario a una empresa (solo SuperAdmin)
// @Tags         Memberships
// @Accept       json
// @Produce      json
// @Param        membership  body      dtos.CreateMembershipDTO  true  "Datos de la membresía"
// @Success      201         {object}  map[string]interface{}
// @Failure      400         {object}  map[string]interface{}
// @Failure      403         {object}  map[string]interface{}
// @Failure      500         {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /admin/memberships [post]
func (h *MembershipHandler) CreateMembership(c *gin.Context) {
	// SOLO SuperAdmin puede crear memberships (MVP)
	// Los clientes crean usuarios que automáticamente se agregan a su empresa
	// Sistema de invitaciones → Fase 2
	role, _ := c.Get("role")
	if role != "superadmin" {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Only superadmin can assign users to companies. Regular users should create new users instead.",
		})
		return
	}

	var dto dtos.CreateMembershipDTO
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"status": "success", "data": membership})
}

func (h *MembershipHandler) UpdateMembership(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(400, gin.H{"error": "ID inválido"})
		return
	}

	// Validar que el membership pertenece a la empresa del usuario
	role, _ := c.Get("role")
	if role != "superadmin" {
		membership, err := h.membershipService.GetMembershipByID(uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Membership not found"})
			return
		}
		companyIDVal, exists := c.Get("company_id")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "No company context"})
			return
		}
		companyID := companyIDVal.(uint)
		if membership.CompanyID == nil || *membership.CompanyID != companyID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}
	}

	var dto dtos.UpdateMembershipDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	membership, err := h.membershipService.UpdateMembership(uint(id), dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": membership})
}

func (h *MembershipHandler) DeleteMembership(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(400, gin.H{"error": "ID inválido"})
		return
	}

	// Validar que el membership pertenece a la empresa del usuario
	role, _ := c.Get("role")
	if role != "superadmin" {
		membership, err := h.membershipService.GetMembershipByID(uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Membership not found"})
			return
		}
		companyIDVal, exists := c.Get("company_id")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "No company context"})
			return
		}
		companyID := companyIDVal.(uint)
		if membership.CompanyID == nil || *membership.CompanyID != companyID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}
	}

	if err := h.membershipService.DeleteMembership(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Membership deleted"})
}
