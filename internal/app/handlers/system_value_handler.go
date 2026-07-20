package handlers

import (
	"net/http"
	"strconv"

	"dvra-api/internal/app/dtos"
	"dvra-api/internal/app/services"
	"dvra-api/internal/shared/apperr"
	"dvra-api/internal/shared/authctx"

	"github.com/geomark27/loom-go/pkg/helpers"
	"github.com/gin-gonic/gin"
)

type SystemValueHandler struct {
	service services.SystemValueService
	logger  helpers.Logger
}

func NewSystemValueHandler(service services.SystemValueService) *SystemValueHandler {
	return &SystemValueHandler{service: service, logger: helpers.NewLogger()}
}

func tenantCompanyID(c *gin.Context) (uint, bool) {
	companyID, ok := authctx.CompanyID(c)
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "tenant context required"})
	}
	return companyID, ok
}

// GetCategories godoc
// @Summary Lista categorías de catálogos
// @Tags System Values
// @Security BearerAuth
// @Param module query string false "Módulo"
// @Success 200 {object} map[string]interface{}
// @Router /system-values [get]
func (h *SystemValueHandler) GetCategories(c *gin.Context) {
	categories, err := h.service.GetCategories(c.Query("module"))
	respondSystemValue(c, http.StatusOK, categories, err)
}

// GetByCategory preserves the original route while clients migrate to the
// unambiguous /system-values/modules/:module/:category endpoint.
func (h *SystemValueHandler) GetByCategory(c *gin.Context) {
	companyID, ok := tenantCompanyID(c)
	if !ok {
		return
	}
	values, err := h.service.GetByLegacyCategory(c.Param("category"), &companyID)
	respondSystemValue(c, http.StatusOK, values, err)
}

// GetByModuleCategory godoc
// @Summary Obtiene valores efectivos de un catálogo para el tenant activo
// @Tags System Values
// @Security BearerAuth
// @Param module path string true "Módulo"
// @Param category path string true "Categoría"
// @Success 200 {object} map[string]interface{}
// @Router /system-values/modules/{module}/{category} [get]
func (h *SystemValueHandler) GetByModuleCategory(c *gin.Context) {
	companyID, ok := tenantCompanyID(c)
	if !ok {
		return
	}
	values, err := h.service.GetByCategory(c.Param("module"), c.Param("category"), &companyID)
	respondSystemValue(c, http.StatusOK, values, err)
}

func (h *SystemValueHandler) GetAll(c *gin.Context) {
	values, err := h.service.GetAll()
	respondSystemValue(c, http.StatusOK, values, err)
}

// CreateTenant creates a value owned by the company in the JWT. company_id is
// deliberately absent from the request DTO.
func (h *SystemValueHandler) CreateTenant(c *gin.Context) {
	companyID, ok := tenantCompanyID(c)
	if !ok {
		return
	}
	var dto dtos.CreateSystemValueDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	value, err := h.service.CreateTenant(companyID, dto)
	respondSystemValue(c, http.StatusCreated, value, err)
}

func (h *SystemValueHandler) UpdateTenant(c *gin.Context) {
	companyID, ok := tenantCompanyID(c)
	if !ok {
		return
	}
	id, ok := systemValueID(c)
	if !ok {
		return
	}
	var dto dtos.UpdateSystemValueDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	value, err := h.service.UpdateTenant(companyID, id, dto)
	respondSystemValue(c, http.StatusOK, value, err)
}

func (h *SystemValueHandler) DeleteTenant(c *gin.Context) {
	companyID, ok := tenantCompanyID(c)
	if !ok {
		return
	}
	id, ok := systemValueID(c)
	if !ok {
		return
	}
	err := h.service.DeleteTenant(companyID, id)
	respondSystemValue(c, http.StatusOK, gin.H{"message": "catalog value deleted"}, err)
}

func (h *SystemValueHandler) UpsertTenantOverride(c *gin.Context) {
	companyID, ok := tenantCompanyID(c)
	if !ok {
		return
	}
	id, ok := systemValueID(c)
	if !ok {
		return
	}
	var dto dtos.UpsertSystemValueOverrideDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := h.service.UpsertTenantOverride(companyID, id, dto)
	respondSystemValue(c, http.StatusOK, gin.H{"message": "catalog override saved"}, err)
}

// Create creates a landlord-owned global value.
func (h *SystemValueHandler) Create(c *gin.Context) {
	var dto dtos.CreateSystemValueDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	value, err := h.service.CreateGlobal(dto)
	respondSystemValue(c, http.StatusCreated, value, err)
}

func (h *SystemValueHandler) Update(c *gin.Context) {
	id, ok := systemValueID(c)
	if !ok {
		return
	}
	var dto dtos.UpdateSystemValueDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	value, err := h.service.UpdateGlobal(id, dto)
	respondSystemValue(c, http.StatusOK, value, err)
}

func (h *SystemValueHandler) Delete(c *gin.Context) {
	id, ok := systemValueID(c)
	if !ok {
		return
	}
	err := h.service.DeleteGlobal(id)
	respondSystemValue(c, http.StatusOK, gin.H{"message": "global catalog value deleted"}, err)
}

func systemValueID(c *gin.Context) (uint, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return 0, false
	}
	return uint(id), true
}

func respondSystemValue(c *gin.Context, status int, data interface{}, err error) {
	if err != nil {
		c.JSON(apperr.StatusCode(err), gin.H{"error": err.Error()})
		return
	}
	c.JSON(status, gin.H{"status": "success", "data": data})
}
