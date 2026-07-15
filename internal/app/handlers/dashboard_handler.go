package handlers

import (
	"net/http"

	"dvra-api/internal/app/services"
	"dvra-api/internal/shared/authctx"

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
	companyID, ok := authctx.CompanyID(c)
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "No company context"})
		return
	}

	stats, err := h.dashboardService.GetStats(companyID)
	if err != nil {
		h.logger.Error("Failed to get dashboard stats", map[string]interface{}{"error": err.Error()})
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve dashboard stats"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": stats})
}
