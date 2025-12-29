package handlers

import (
	"net/http"
	"strconv"

	"dvra-api/internal/app/dtos"
	"dvra-api/internal/app/services"

	"github.com/geomark27/loom-go/pkg/helpers"
	"github.com/gin-gonic/gin"
)

type JobHandler struct {
	jobService services.JobService
	logger     helpers.Logger
}

func NewJobHandler(jobService services.JobService) *JobHandler {
	return &JobHandler{jobService: jobService, logger: helpers.NewLogger()}
}

func (h *JobHandler) GetJobs(c *gin.Context) {
	role, _ := c.Get("role")

	// Leer filtros de query params
	var filters dtos.JobFilters
	if err := c.ShouldBindQuery(&filters); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Debug: log filtros recibidos
	h.logger.Info("Filters received: Status=%s, LocationType=%s, CityID=%v", filters.Status, filters.LocationType, filters.CityID)

	// SuperAdmin puede ver todos los jobs
	if role == "superadmin" {
		jobs, err := h.jobService.GetAllJobsWithFilters(filters)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve jobs"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"jobs": dtos.ToJobResponseList(jobs), "count": len(jobs)}})
		return
	}

	// Usuarios normales solo ven jobs de su empresa
	companyIDVal, exists := c.Get("company_id")
	if !exists {
		c.JSON(http.StatusForbidden, gin.H{"error": "No company context"})
		return
	}

	companyID := companyIDVal.(uint)
	jobs, err := h.jobService.GetJobsByCompanyIDWithFilters(companyID, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve jobs"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Jobs retrieved successfully",
		"data": gin.H{
			"jobs":  dtos.ToJobResponseList(jobs),
			"count": len(jobs),
		},
	})
}

func (h *JobHandler) GetJob(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	job, err := h.jobService.GetJobByID(uint(id))
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
		if job.CompanyID != companyID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": dtos.ToJobResponse(job)})
}

func (h *JobHandler) CreateJob(c *gin.Context) {
	var dto dtos.CreateJobDTO
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
		dto.CompanyID = companyID // Forzar company del token, ignorar el enviado
	}

	job, err := h.jobService.CreateJob(dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"status": "success", "data": dtos.ToJobResponse(job)})
}

func (h *JobHandler) UpdateJob(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	// Validar que el job pertenece a la empresa del usuario
	role, _ := c.Get("role")
	if role != "superadmin" {
		job, err := h.jobService.GetJobByID(uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
			return
		}
		companyIDVal, exists := c.Get("company_id")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "No company context"})
			return
		}
		companyID := companyIDVal.(uint)
		if job.CompanyID != companyID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}
	}

	var dto dtos.UpdateJobDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updatedJob, err := h.jobService.UpdateJob(uint(id), dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": dtos.ToJobResponse(updatedJob)})
}

func (h *JobHandler) DeleteJob(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	// Validar que el job pertenece a la empresa del usuario
	role, _ := c.Get("role")
	if role != "superadmin" {
		job, err := h.jobService.GetJobByID(uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
			return
		}
		companyIDVal, exists := c.Get("company_id")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "No company context"})
			return
		}
		companyID := companyIDVal.(uint)
		if job.CompanyID != companyID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}
	}

	if err := h.jobService.DeleteJob(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Job deleted"})
}

// PublishJob godoc
// @Summary      Publicar empleo
// @Description  Cambia el estado del job a 'active'
// @Tags         Jobs
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Job ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      403  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /jobs/{id}/publish [patch]
func (h *JobHandler) PublishJob(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	// Validar que el job pertenece a la empresa del usuario
	role, _ := c.Get("role")
	if role != "superadmin" {
		job, err := h.jobService.GetJobByID(uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
			return
		}
		companyIDVal, exists := c.Get("company_id")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "No company context"})
			return
		}
		companyID := companyIDVal.(uint)
		if job.CompanyID != companyID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}
	}

	job, err := h.jobService.PublishJob(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Job published", "data": dtos.ToJobResponse(job)})
}

// CloseJob godoc
// @Summary      Cerrar empleo
// @Description  Cambia el estado del job a 'closed'
// @Tags         Jobs
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Job ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      403  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /jobs/{id}/close [patch]
func (h *JobHandler) CloseJob(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	// Validar que el job pertenece a la empresa del usuario
	role, _ := c.Get("role")
	if role != "superadmin" {
		job, err := h.jobService.GetJobByID(uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
			return
		}
		companyIDVal, exists := c.Get("company_id")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "No company context"})
			return
		}
		companyID := companyIDVal.(uint)
		if job.CompanyID != companyID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}
	}

	job, err := h.jobService.CloseJob(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Job closed", "data": dtos.ToJobResponse(job)})
}
