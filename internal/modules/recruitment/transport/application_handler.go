package transport

import (
	"dvra-api/internal/shared/authctx"
	"net/http"
	"strconv"

	"dvra-api/internal/modules/recruitment/service"
	"dvra-api/internal/shared/apperr"

	"github.com/geomark27/loom-go/pkg/helpers"
	"github.com/gin-gonic/gin"
)

type ApplicationHandler struct {
	applicationService service.ApplicationService
	logger             helpers.Logger
}

func NewApplicationHandler(applicationService service.ApplicationService) *ApplicationHandler {
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

	// Tenant users only see applications allowed by their active membership.
	scope, ok := authctx.TenantScope(c)
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "No company context"})
		return
	}

	applications, err := h.applicationService.GetApplicationsWithScope(scope)
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
	scope, _ := authctx.TenantScope(c)
	application, err := h.applicationService.GetApplicationByID(uint(id), scope)
	if err != nil {
		c.JSON(apperr.StatusCode(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": application})
}

// CreateApplication godoc
// @Summary      Crear postulación
// @Description  Crea una nueva postulación a un empleo
// @Tags         Applications
// @Accept       json
// @Produce      json
// @Param        application  body      service.CreateApplicationDTO  true  "Datos de la postulación"
// @Success      201          {object}  map[string]interface{}
// @Failure      400          {object}  map[string]interface{}
// @Failure      403          {object}  map[string]interface{}
// @Failure      500          {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /applications [post]
func (h *ApplicationHandler) CreateApplication(c *gin.Context) {
	var dto service.CreateApplicationDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	companyID, ok := authctx.CompanyID(c)
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "No company context"})
		return
	}
	dto.CompanyID = companyID

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
	scope, _ := authctx.TenantScope(c)

	var dto service.UpdateApplicationDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	application, err := h.applicationService.UpdateApplication(uint(id), scope, dto)
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
	scope, _ := authctx.TenantScope(c)

	if err := h.applicationService.DeleteApplication(uint(id), scope); err != nil {
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

	scope, ok := authctx.TenantScope(c)
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "No company context"})
		return
	}

	applicationsByStage, err := h.applicationService.GetApplicationsGroupedByStage(scope)
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
// @Param        body  body      service.MoveApplicationDTO  true  "Stage data"
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
	scope, _ := authctx.TenantScope(c)

	var dto service.MoveApplicationDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	application, err := h.applicationService.MoveToStage(uint(id), scope, dto.Stage)
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
// @Param        body  body      service.RateApplicationDTO  true  "Rating data"
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
	scope, _ := authctx.TenantScope(c)

	var dto service.RateApplicationDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	application, err := h.applicationService.RateApplication(uint(id), scope, dto.Rating)
	if err != nil {
		c.JSON(apperr.StatusCode(err), gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": application})
}
