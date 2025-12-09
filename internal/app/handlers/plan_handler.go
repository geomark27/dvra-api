package handlers

import (
	"net/http"
	"strconv"

	"dvra-api/internal/app/dtos"
	"dvra-api/internal/app/services"

	"github.com/geomark27/loom-go/pkg/helpers"
	"github.com/gin-gonic/gin"
)

// PlanHandler handles plan-related HTTP requests
type PlanHandler struct {
	planService services.PlanService
	logger      helpers.Logger
}

// NewPlanHandler creates a new plan handler
func NewPlanHandler(planService services.PlanService) *PlanHandler {
	return &PlanHandler{
		planService: planService,
		logger:      helpers.NewLogger(),
	}
}

// CreatePlan creates a new plan (SuperAdmin only)
// @Summary Create plan
// @Description Create a new subscription plan (SuperAdmin only)
// @Tags plans
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param plan body dtos.PlanDTO true "Plan data"
// @Success 201 {object} dtos.PlanResponse
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Router /admin/plans [post]
func (h *PlanHandler) CreatePlan(c *gin.Context) {
	var dto dtos.PlanDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	plan, err := h.planService.CreatePlan(&dto)
	if err != nil {
		if err == services.ErrPlanSlugExists {
			c.JSON(http.StatusConflict, gin.H{"error": "Plan with this slug already exists"})
			return
		}
		h.logger.Error("Failed to create plan", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create plan"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data":   gin.H{"plan": plan},
	})
}

// GetPlans retrieves all plans
// @Summary Get all plans
// @Description Get all plans (SuperAdmin sees all, others see only public active plans)
// @Tags plans
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /admin/plans [get]
func (h *PlanHandler) GetPlans(c *gin.Context) {
	role, _ := c.Get("role")

	var plans []dtos.PlanResponse
	var err error

	// SuperAdmin can see all plans
	if role == "superadmin" {
		plans, err = h.planService.GetAllPlans()
	} else {
		// Regular users only see public active plans
		plans, err = h.planService.GetPublicPlans()
	}

	if err != nil {
		h.logger.Error("Failed to retrieve plans", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve plans"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"plans": plans,
			"count": len(plans),
		},
	})
}

// GetPlanByID retrieves a plan by ID
// @Summary Get plan by ID
// @Description Get a specific plan by ID
// @Tags plans
// @Produce json
// @Security BearerAuth
// @Param id path int true "Plan ID"
// @Success 200 {object} dtos.PlanResponse
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /admin/plans/{id} [get]
func (h *PlanHandler) GetPlanByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid plan ID"})
		return
	}

	plan, err := h.planService.GetPlanByID(uint(id))
	if err != nil {
		if err == services.ErrPlanNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Plan not found"})
			return
		}
		h.logger.Error("Failed to retrieve plan", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve plan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   gin.H{"plan": plan},
	})
}

// GetPlanBySlug retrieves a plan by slug
// @Summary Get plan by slug
// @Description Get a specific plan by slug
// @Tags plans
// @Produce json
// @Param slug path string true "Plan slug"
// @Success 200 {object} dtos.PlanResponse
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /plans/{slug} [get]
func (h *PlanHandler) GetPlanBySlug(c *gin.Context) {
	slug := c.Param("slug")

	plan, err := h.planService.GetPlanBySlug(slug)
	if err != nil {
		if err == services.ErrPlanNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Plan not found"})
			return
		}
		h.logger.Error("Failed to retrieve plan", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve plan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   gin.H{"plan": plan},
	})
}

// GetPublicPlans retrieves all public active plans (no auth required)
// @Summary Get public plans
// @Description Get all public active plans for pricing page
// @Tags plans
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]string
// @Router /plans [get]
func (h *PlanHandler) GetPublicPlans(c *gin.Context) {
	plans, err := h.planService.GetPublicPlans()
	if err != nil {
		h.logger.Error("Failed to retrieve public plans", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve plans"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"plans": plans,
			"count": len(plans),
		},
	})
}

// UpdatePlan updates an existing plan (SuperAdmin only)
// @Summary Update plan
// @Description Update an existing plan (SuperAdmin only)
// @Tags plans
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Plan ID"
// @Param plan body dtos.UpdatePlanDTO true "Plan data"
// @Success 200 {object} dtos.PlanResponse
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /admin/plans/{id} [put]
func (h *PlanHandler) UpdatePlan(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid plan ID"})
		return
	}

	var dto dtos.UpdatePlanDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	plan, err := h.planService.UpdatePlan(uint(id), &dto)
	if err != nil {
		if err == services.ErrPlanNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Plan not found"})
			return
		}
		h.logger.Error("Failed to update plan", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update plan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Plan updated successfully",
		"data":    gin.H{"plan": plan},
	})
}

// TogglePlanStatus enables/disables a plan (SuperAdmin only)
// @Summary Toggle plan status
// @Description Enable or disable a plan (SuperAdmin only)
// @Tags plans
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Plan ID"
// @Param status body dtos.TogglePlanStatusDTO true "Plan status"
// @Success 200 {object} dtos.PlanResponse
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /admin/plans/{id}/toggle [patch]
func (h *PlanHandler) TogglePlanStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid plan ID"})
		return
	}

	var dto dtos.TogglePlanStatusDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	plan, err := h.planService.TogglePlanStatus(uint(id), dto.IsActive)
	if err != nil {
		if err == services.ErrPlanNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Plan not found"})
			return
		}
		h.logger.Error("Failed to toggle plan status", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to toggle plan status"})
		return
	}

	statusText := "disabled"
	if dto.IsActive {
		statusText = "enabled"
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Plan " + statusText + " successfully",
		"data":    gin.H{"plan": plan},
	})
}

// DeletePlan soft deletes a plan (SuperAdmin only)
// @Summary Delete plan
// @Description Soft delete a plan (SuperAdmin only)
// @Tags plans
// @Produce json
// @Security BearerAuth
// @Param id path int true "Plan ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Router /admin/plans/{id} [delete]
func (h *PlanHandler) DeletePlan(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid plan ID"})
		return
	}

	err = h.planService.DeletePlan(uint(id))
	if err != nil {
		if err == services.ErrPlanNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Plan not found"})
			return
		}
		if err == services.ErrPlanInUse {
			c.JSON(http.StatusConflict, gin.H{"error": "Cannot delete plan that is in use by companies"})
			return
		}
		h.logger.Error("Failed to delete plan", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete plan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Plan deleted successfully",
	})
}

// AssignPlanToCompany assigns a plan to a company (SuperAdmin only)
// @Summary Assign plan to company
// @Description Assign a subscription plan to a company (SuperAdmin only)
// @Tags plans
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param assignment body dtos.AssignPlanToCompanyDTO true "Plan assignment"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /admin/plans/assign [post]
func (h *PlanHandler) AssignPlanToCompany(c *gin.Context) {
	var dto dtos.AssignPlanToCompanyDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get plan by ID to get slug
	plan, err := h.planService.GetPlanByID(dto.PlanID)
	if err != nil {
		if err == services.ErrPlanNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Plan not found"})
			return
		}
		h.logger.Error("Failed to retrieve plan", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign plan"})
		return
	}

	err = h.planService.AssignPlanToCompany(dto.CompanyID, plan.Slug)
	if err != nil {
		if err == services.ErrCompanyNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Company not found"})
			return
		}
		h.logger.Error("Failed to assign plan to company", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assign plan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Plan assigned to company successfully",
	})
}
