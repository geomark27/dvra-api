package handlers

import (
	"net/http"

	"dvra-api/internal/app/dtos"
	"dvra-api/internal/app/services"
	"dvra-api/internal/shared/apperr"
	"dvra-api/internal/shared/authctx"

	"github.com/gin-gonic/gin"
)

type CompanySettingsHandler struct {
	service *services.CompanySettingsService
}

func NewCompanySettingsHandler(service *services.CompanySettingsService) *CompanySettingsHandler {
	return &CompanySettingsHandler{service: service}
}

func (h *CompanySettingsHandler) Get(c *gin.Context) {
	companyID, ok := authctx.CompanyID(c)
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "tenant context required"})
		return
	}
	settings, err := h.service.Get(companyID)
	respondCompanySettings(c, settings, err)
}

func (h *CompanySettingsHandler) Update(c *gin.Context) {
	companyID, ok := authctx.CompanyID(c)
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "tenant context required"})
		return
	}
	var dto dtos.UpdateCompanySettingsDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	settings, err := h.service.Update(companyID, dto)
	respondCompanySettings(c, settings, err)
}

func (h *CompanySettingsHandler) Subscription(c *gin.Context) {
	companyID, ok := authctx.CompanyID(c)
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "tenant context required"})
		return
	}
	subscription, err := h.service.Subscription(companyID)
	respondCompanySettings(c, subscription, err)
}

func (h *CompanySettingsHandler) Usage(c *gin.Context) {
	companyID, ok := authctx.CompanyID(c)
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "tenant context required"})
		return
	}
	usage, err := h.service.Usage(companyID)
	respondCompanySettings(c, usage, err)
}

func respondCompanySettings(c *gin.Context, data interface{}, err error) {
	if err != nil {
		c.JSON(apperr.StatusCode(err), gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": data})
}
