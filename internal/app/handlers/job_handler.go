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

// GetJobs godoc
// @Summary      Listar empleos
// @Description  Obtiene empleos (todos si es SuperAdmin, de la empresa si es usuario normal)
// @Tags         Jobs
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      403  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /jobs [get]
func (h *JobHandler) GetJobs(c *gin.Context) {
	role, _ := c.Get("role")

	// SuperAdmin puede ver todos los jobs
	if role == "superadmin" {
		jobs, err := h.jobService.GetAllJobs()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve jobs"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"jobs": jobs, "count": len(jobs)}})
		return
	}

	// Usuarios normales solo ven jobs de su empresa
	companyIDVal, exists := c.Get("company_id")
	if !exists {
		c.JSON(http.StatusForbidden, gin.H{"error": "No company context"})
		return
	}

	companyID := companyIDVal.(uint)
	jobs, err := h.jobService.GetJobsByCompanyID(companyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve jobs"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"jobs": jobs, "count": len(jobs)}})
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

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": job})
}

// CreateJob godoc
// @Summary      Crear empleo
// @Description  Crea un nuevo empleo en la empresa actual
// @Tags         Jobs
// @Accept       json
// @Produce      json
// @Param        job  body      dtos.CreateJobDTO  true  "Datos del empleo"
// @Success      201  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      403  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /jobs [post]
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
	c.JSON(http.StatusCreated, gin.H{"status": "success", "data": job})
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
	job, err := h.jobService.UpdateJob(uint(id), dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": job})
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
