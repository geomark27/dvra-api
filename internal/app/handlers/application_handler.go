package handlers

import (
	"net/http"
	"strconv"

	"dvra-api/internal/app/dtos"
	"dvra-api/internal/app/services"

	"github.com/geomark27/loom-go/pkg/helpers"
	"github.com/gin-gonic/gin"
)

type ApplicationHandler struct {
	applicationService services.ApplicationService
	logger             helpers.Logger
}

func NewApplicationHandler(applicationService services.ApplicationService) *ApplicationHandler {
	return &ApplicationHandler{applicationService: applicationService, logger: helpers.NewLogger()}
}

// GetApplications godoc
// @Summary      Listar postulaciones
// @Description  Obtiene postulaciones (todas si es SuperAdmin, de la empresa si es usuario normal)
// @Tags         Applications
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      403  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /applications [get]
func (h *ApplicationHandler) GetApplications(c *gin.Context) {
	role, _ := c.Get("role")

	// SuperAdmin puede ver todas las applications
	if role == "superadmin" {
		applications, err := h.applicationService.GetAllApplications()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve applications"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"applications": applications, "count": len(applications)}})
		return
	}

	// Usuarios normales solo ven applications de su empresa
	companyIDVal, exists := c.Get("company_id")
	if !exists {
		c.JSON(http.StatusForbidden, gin.H{"error": "No company context"})
		return
	}

	companyID := companyIDVal.(uint)
	applications, err := h.applicationService.GetApplicationsByCompanyID(companyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve applications"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"applications": applications, "count": len(applications)}})
}

func (h *ApplicationHandler) GetApplication(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	application, err := h.applicationService.GetApplicationByID(uint(id))
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
		if application.CompanyID != companyID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": application})
}

// CreateApplication godoc
// @Summary      Crear postulación
// @Description  Crea una nueva postulación a un empleo
// @Tags         Applications
// @Accept       json
// @Produce      json
// @Param        application  body      dtos.CreateApplicationDTO  true  "Datos de la postulación"
// @Success      201          {object}  map[string]interface{}
// @Failure      400          {object}  map[string]interface{}
// @Failure      403          {object}  map[string]interface{}
// @Failure      500          {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /applications [post]
func (h *ApplicationHandler) CreateApplication(c *gin.Context) {
	var dto dtos.CreateApplicationDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Forzar company_id del token para usuarios normales
	role, _ := c.Get("role")
	if role != "superadmin" {
		companyIDVal, exists := c.Get("company_id")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "No company context"})
			return
		}
		companyID := companyIDVal.(uint)
		dto.CompanyID = companyID
	}

	application, err := h.applicationService.CreateApplication(dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"status": "success", "data": application})
}

func (h *ApplicationHandler) UpdateApplication(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	// Validar que la application pertenece a la empresa del usuario
	role, _ := c.Get("role")
	if role != "superadmin" {
		application, err := h.applicationService.GetApplicationByID(uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Application not found"})
			return
		}
		companyIDVal, exists := c.Get("company_id")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "No company context"})
			return
		}
		companyID := companyIDVal.(uint)
		if application.CompanyID != companyID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}
	}

	var dto dtos.UpdateApplicationDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	application, err := h.applicationService.UpdateApplication(uint(id), dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": application})
}

func (h *ApplicationHandler) DeleteApplication(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	// Validar que la application pertenece a la empresa del usuario
	role, _ := c.Get("role")
	if role != "superadmin" {
		application, err := h.applicationService.GetApplicationByID(uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Application not found"})
			return
		}
		companyIDVal, exists := c.Get("company_id")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "No company context"})
			return
		}
		companyID := companyIDVal.(uint)
		if application.CompanyID != companyID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}
	}

	if err := h.applicationService.DeleteApplication(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Application deleted"})
}
