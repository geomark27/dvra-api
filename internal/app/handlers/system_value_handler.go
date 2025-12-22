package handlers

import (
	"dvra-api/internal/app/dtos"
	"dvra-api/internal/app/services"
	"net/http"
	"strconv"

	"github.com/geomark27/loom-go/pkg/helpers"
	"github.com/gin-gonic/gin"
)

type SystemValueHandler struct {
	service services.SystemValueService
	logger  helpers.Logger
}

func NewSystemValueHandler(service services.SystemValueService) *SystemValueHandler {
	return &SystemValueHandler{
		service: service,
		logger:  helpers.NewLogger(),
	}
}

// GetByCategory godoc
// @Summary      Obtener valores por categoría
// @Description  Retorna todos los valores de sistema de una categoría específica
// @Tags         System Values
// @Accept       json
// @Produce      json
// @Param        category  path  string  true  "Categoría (ej: employment_type, contract_type)"
// @Success      200       {object}  map[string]interface{}
// @Failure      500       {object}  map[string]interface{}
// @Router       /system-values/{category} [get]
func (h *SystemValueHandler) GetByCategory(c *gin.Context) {
	category := c.Param("category")

	// Get company_id from header if present
	var companyID *uint
	if companyIDStr := c.GetHeader("X-Company-ID"); companyIDStr != "" {
		if id, err := strconv.ParseUint(companyIDStr, 10, 32); err == nil {
			companyIDUint := uint(id)
			companyID = &companyIDUint
		}
	}

	values, err := h.service.GetByCategory(category, companyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve system values"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   values,
	})
}

// GetAll godoc
// @Summary      Listar todos los valores de sistema
// @Description  Obtiene todos los valores de sistema (solo SuperAdmin)
// @Tags         System Values
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /admin/system-values [get]
func (h *SystemValueHandler) GetAll(c *gin.Context) {
	values, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve system values"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   values,
	})
}

// Create godoc
// @Summary      Crear valor de sistema
// @Description  Crea un nuevo valor de sistema (solo SuperAdmin)
// @Tags         System Values
// @Accept       json
// @Produce      json
// @Param        value  body      dtos.CreateSystemValueDTO  true  "Datos del valor"
// @Success      201    {object}  map[string]interface{}
// @Failure      400    {object}  map[string]interface{}
// @Failure      500    {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /admin/system-values [post]
func (h *SystemValueHandler) Create(c *gin.Context) {
	var dto dtos.CreateSystemValueDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	value, err := h.service.Create(dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "System value created successfully",
		"data":    value,
	})
}

// Update updates a system value (admin only)
func (h *SystemValueHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var dto dtos.UpdateSystemValueDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	value, err := h.service.Update(uint(id), dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "System value updated successfully",
		"data":    value,
	})
}

// Delete deletes a system value (admin only)
func (h *SystemValueHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.service.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "System value deleted successfully",
	})
}
