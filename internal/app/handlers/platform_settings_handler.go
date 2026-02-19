package handlers

import (
	"dvra-api/internal/app/dtos"
	"dvra-api/internal/app/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// PlatformSettingsHandler maneja los endpoints de platform settings
type PlatformSettingsHandler struct {
	service *services.PlatformSettingsService
}

// NewPlatformSettingsHandler crea una nueva instancia del handler
func NewPlatformSettingsHandler(service *services.PlatformSettingsService) *PlatformSettingsHandler {
	return &PlatformSettingsHandler{service: service}
}

// ============================================================================
// PÚBLICO (Sin autenticación)
// ============================================================================

// GetPublicSettings obtiene la configuración pública de la plataforma
// @Summary      Get public platform settings
// @Description  Returns public platform configuration (branding, contact, URLs)
// @Tags         public
// @Accept       json
// @Produce      json
// @Success      200 {object} dtos.PlatformSettingsPublicDTO
// @Failure      500 {object} map[string]string
// @Router       /public/platform-settings [get]
func (h *PlatformSettingsHandler) GetPublicSettings(c *gin.Context) {
	settings, err := h.service.GetPublic()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get platform settings"})
		return
	}

	c.JSON(http.StatusOK, settings)
}

// ============================================================================
// SUPERADMIN (Requiere autenticación y rol superadmin)
// ============================================================================

// GetFullSettings obtiene la configuración completa de la plataforma
// @Summary      Get full platform settings (SuperAdmin only)
// @Description  Returns complete platform configuration including business defaults and legal info
// @Tags         superadmin
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} dtos.PlatformSettingsFullDTO
// @Failure      401 {object} map[string]string
// @Failure      403 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /superadmin/platform-settings [get]
func (h *PlatformSettingsHandler) GetFullSettings(c *gin.Context) {
	settings, err := h.service.GetFull()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get platform settings"})
		return
	}

	c.JSON(http.StatusOK, settings)
}

// UpdateSettings actualiza la configuración de la plataforma
// @Summary      Update platform settings (SuperAdmin only)
// @Description  Updates platform configuration (partial update supported)
// @Tags         superadmin
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        settings body dtos.UpdatePlatformSettingsDTO true "Settings to update"
// @Success      200 {object} dtos.PlatformSettingsFullDTO
// @Failure      400 {object} map[string]string
// @Failure      401 {object} map[string]string
// @Failure      403 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /superadmin/platform-settings [put]
func (h *PlatformSettingsHandler) UpdateSettings(c *gin.Context) {
	var updateDTO dtos.UpdatePlatformSettingsDTO
	if err := c.ShouldBindJSON(&updateDTO); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Obtener el ID del usuario del contexto (del middleware de autenticación)
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	settings, err := h.service.Update(&updateDTO, userID.(uint))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, settings)
}
