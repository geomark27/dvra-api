package handlers

import (
	"dvra-api/internal/shared/authctx"
	"net/http"
	"strconv"

	"dvra-api/internal/app/dtos"
	"dvra-api/internal/app/services"
	"dvra-api/internal/shared/apperr"

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

	// SuperAdmin puede ver todas las applications
	if authctx.IsSuperAdmin(c) {
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
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid application ID"})
		return
	}
	application, err := h.applicationService.GetApplicationByID(uint(id))
	if err != nil {
		c.JSON(apperr.StatusCode(err), gin.H{"error": err.Error()})
		return
	}

	// Validar acceso: SuperAdmin o miembro de la misma empresa
	if !authctx.IsSuperAdmin(c) {
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
	if !authctx.IsSuperAdmin(c) {
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
		c.JSON(apperr.StatusCode(err), gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"status": "success", "data": application})
}

func (h *ApplicationHandler) UpdateApplication(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid application ID"})
		return
	}

	// Validar que la application pertenece a la empresa del usuario
	if !authctx.IsSuperAdmin(c) {
		application, err := h.applicationService.GetApplicationByID(uint(id))
		if err != nil {
			c.JSON(apperr.StatusCode(err), gin.H{"error": err.Error()})
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
		c.JSON(apperr.StatusCode(err), gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": application})
}

func (h *ApplicationHandler) DeleteApplication(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid application ID"})
		return
	}

	// Validar que la application pertenece a la empresa del usuario
	if !authctx.IsSuperAdmin(c) {
		application, err := h.applicationService.GetApplicationByID(uint(id))
		if err != nil {
			c.JSON(apperr.StatusCode(err), gin.H{"error": err.Error()})
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
		c.JSON(apperr.StatusCode(err), gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Application deleted"})
}

// GetApplicationsByStage godoc
// @Summary      Get applications grouped by stage
// @Description  Obtiene postulaciones agrupadas por etapa (para Kanban board)
// @Tags         Applications
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      403  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /applications/by-stage [get]
func (h *ApplicationHandler) GetApplicationsByStage(c *gin.Context) {

	var companyID uint
	if authctx.IsSuperAdmin(c) {
		// SuperAdmin needs to specify company_id via query param
		queryCompanyID := c.Query("company_id")
		if queryCompanyID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "SuperAdmin must provide company_id query parameter"})
			return
		}
		id, err := strconv.ParseUint(queryCompanyID, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid company_id"})
			return
		}
		companyID = uint(id)
	} else {
		companyIDVal, exists := c.Get("company_id")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "No company context"})
			return
		}
		companyID = companyIDVal.(uint)
	}

	applicationsByStage, err := h.applicationService.GetApplicationsGroupedByStage(companyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve applications by stage"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": applicationsByStage})
}

// MoveApplication godoc
// @Summary      Move application to a stage
// @Description  Mueve una postulación a una etapa del pipeline
// @Tags         Applications
// @Accept       json
// @Produce      json
// @Param        id    path      int                      true  "Application ID"
// @Param        body  body      dtos.MoveApplicationDTO  true  "Stage data"
// @Success      200   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]interface{}
// @Failure      403   {object}  map[string]interface{}
// @Failure      404   {object}  map[string]interface{}
// @Failure      500   {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /applications/{id}/move [patch]
func (h *ApplicationHandler) MoveApplication(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid application ID"})
		return
	}

	// Validate access
	if !authctx.IsSuperAdmin(c) {
		application, err := h.applicationService.GetApplicationByID(uint(id))
		if err != nil {
			c.JSON(apperr.StatusCode(err), gin.H{"error": err.Error()})
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

	var dto dtos.MoveApplicationDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	application, err := h.applicationService.MoveToStage(uint(id), dto.Stage)
	if err != nil {
		c.JSON(apperr.StatusCode(err), gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": application})
}

// RateApplication godoc
// @Summary      Rate an application
// @Description  Asigna una calificación a una postulación
// @Tags         Applications
// @Accept       json
// @Produce      json
// @Param        id    path      int                      true  "Application ID"
// @Param        body  body      dtos.RateApplicationDTO  true  "Rating data"
// @Success      200   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]interface{}
// @Failure      403   {object}  map[string]interface{}
// @Failure      404   {object}  map[string]interface{}
// @Failure      500   {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /applications/{id}/rate [patch]
func (h *ApplicationHandler) RateApplication(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid application ID"})
		return
	}

	// Validate access
	if !authctx.IsSuperAdmin(c) {
		application, err := h.applicationService.GetApplicationByID(uint(id))
		if err != nil {
			c.JSON(apperr.StatusCode(err), gin.H{"error": err.Error()})
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

	var dto dtos.RateApplicationDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	application, err := h.applicationService.RateApplication(uint(id), dto.Rating)
	if err != nil {
		c.JSON(apperr.StatusCode(err), gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": application})
}
