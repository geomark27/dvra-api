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

func (h *MembershipHandler) GetMemberships(c *gin.Context) {
	memberships, err := h.membershipService.GetAllMemberships()
	if err != nil {
		h.logger.Error("Failed to get memberships", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve memberships"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"memberships": memberships, "count": len(memberships)}})
}

func (h *MembershipHandler) GetMembership(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	membership, err := h.membershipService.GetMembershipByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": membership})
}

func (h *MembershipHandler) CreateMembership(c *gin.Context) {
	var dto dtos.CreateMembershipDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
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
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := h.membershipService.DeleteMembership(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Membership deleted"})
}
