package handlers

import (
	"dvra-api/internal/shared/authctx"
	"net/http"

	"dvra-api/internal/app/services"

	"github.com/geomark27/loom-go/pkg/helpers"
	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	dashboardService services.DashboardService
	logger           helpers.Logger
}

func NewDashboardHandler(dashboardService services.DashboardService) *DashboardHandler {
	return &DashboardHandler{dashboardService: dashboardService, logger: helpers.NewLogger()}
}

// GetStats godoc
// @Summary      Get dashboard statistics
// @Description  Obtiene estadísticas del dashboard para la empresa actual
// @Tags         Dashboard
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      403  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /dashboard/stats [get]
func (h *DashboardHandler) GetStats(c *gin.Context) {

	// SuperAdmin gets stats for all companies (or can pass company_id as query param)
	if authctx.IsSuperAdmin(c) {
		// For now, superadmin gets empty stats or could aggregate
		// You could add a ?company_id query param for superadmin to see specific company
		queryCompanyID := c.Query("company_id")
		if queryCompanyID != "" {
			// Parse and use the query company_id
			var companyID uint
			if _, err := parseUint(queryCompanyID); err == nil {
				companyID = uint(mustParseUint(queryCompanyID))
				stats, err := h.dashboardService.GetStats(companyID)
				if err != nil {
					h.logger.Error("Failed to get dashboard stats", map[string]interface{}{"error": err.Error()})
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve dashboard stats"})
					return
				}
				c.JSON(http.StatusOK, gin.H{"status": "success", "data": stats})
				return
			}
		}
		// Return message that company_id is needed for superadmin
		c.JSON(http.StatusBadRequest, gin.H{"error": "SuperAdmin must provide company_id query parameter"})
		return
	}

	// Normal users get stats for their company
	companyIDVal, exists := c.Get("company_id")
	if !exists {
		c.JSON(http.StatusForbidden, gin.H{"error": "No company context"})
		return
	}

	companyID := companyIDVal.(uint)
	stats, err := h.dashboardService.GetStats(companyID)
	if err != nil {
		h.logger.Error("Failed to get dashboard stats", map[string]interface{}{"error": err.Error()})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve dashboard stats"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": stats})
}

// Helper functions
func parseUint(s string) (uint64, error) {
	var n uint64
	_, err := parseUintHelper(s, &n)
	return n, err
}

func parseUintHelper(s string, n *uint64) (bool, error) {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false, nil
		}
		*n = *n*10 + uint64(c-'0')
	}
	return true, nil
}

func mustParseUint(s string) uint64 {
	var n uint64
	for _, c := range s {
		n = n*10 + uint64(c-'0')
	}
	return n
}
