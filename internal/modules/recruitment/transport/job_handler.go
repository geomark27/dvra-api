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

type JobHandler struct {
	jobService service.JobService
	logger     helpers.Logger
}

func NewJobHandler(jobService service.JobService) *JobHandler {
	return &JobHandler{jobService: jobService, logger: helpers.NewLogger()}
}

func (h *JobHandler) GetJobs(c *gin.Context) {

	// Leer filtros de query params
	var filters service.JobFilters
	if err := c.ShouldBindQuery(&filters); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Debug: log filtros recibidos
	h.logger.Info("Filters received: Status=%s, LocationType=%s, CityID=%v", filters.Status, filters.LocationType, filters.CityID)

	// Tenant users only see jobs allowed by their active membership.
	scope, ok := authctx.TenantScope(c)
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "No company context"})
		return
	}

	jobs, err := h.jobService.GetJobsWithScope(scope, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve jobs"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Jobs retrieved successfully",
		"data": gin.H{
			"jobs":  service.ToJobResponseList(jobs),
			"count": len(jobs),
		},
	})
}

func (h *JobHandler) GetJob(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid job ID"})
		return
	}
	scope, _ := authctx.TenantScope(c)
	job, err := h.jobService.GetJobByID(uint(id), scope)
	if err != nil {
		c.JSON(apperr.StatusCode(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": service.ToJobResponse(job)})
}

func (h *JobHandler) CreateJob(c *gin.Context) {
	var dto service.CreateJobDTO
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

	job, err := h.jobService.CreateJob(dto)
	if err != nil {
		c.JSON(apperr.StatusCode(err), gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"status": "success", "data": service.ToJobResponse(job)})
}

func (h *JobHandler) UpdateJob(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid job ID"})
		return
	}
	scope, _ := authctx.TenantScope(c)

	var dto service.UpdateJobDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updatedJob, err := h.jobService.UpdateJob(uint(id), scope, dto)
	if err != nil {
		c.JSON(apperr.StatusCode(err), gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": service.ToJobResponse(updatedJob)})
}

func (h *JobHandler) DeleteJob(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid job ID"})
		return
	}
	scope, _ := authctx.TenantScope(c)

	if err := h.jobService.DeleteJob(uint(id), scope); err != nil {
		c.JSON(apperr.StatusCode(err), gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Job deleted"})
}

// PublishJob godoc
// @Summary      Publicar empleo
// @Description  Cambia el estado del job a 'published'
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
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid job ID"})
		return
	}
	scope, _ := authctx.TenantScope(c)

	job, err := h.jobService.PublishJob(uint(id), scope)
	if err != nil {
		c.JSON(apperr.StatusCode(err), gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Job published", "data": service.ToJobResponse(job)})
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
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid job ID"})
		return
	}
	scope, _ := authctx.TenantScope(c)

	job, err := h.jobService.CloseJob(uint(id), scope)
	if err != nil {
		c.JSON(apperr.StatusCode(err), gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Job closed", "data": service.ToJobResponse(job)})
}
