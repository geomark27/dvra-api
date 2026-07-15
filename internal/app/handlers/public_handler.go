package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"dvra-api/internal/app/dtos"
	"dvra-api/internal/app/services"
	"dvra-api/internal/shared/apperr"

	"github.com/geomark27/loom-go/pkg/helpers"
	"github.com/gin-gonic/gin"
)

// PublicHandler handles public career page endpoints
type PublicHandler struct {
	publicService services.PublicService
	logger        helpers.Logger
}

// NewPublicHandler creates a new PublicHandler instance
func NewPublicHandler(publicService services.PublicService) *PublicHandler {
	return &PublicHandler{
		publicService: publicService,
		logger:        helpers.NewLogger(),
	}
}

// GetCompanyBySlug godoc
// @Summary      Get company info by slug
// @Description  Returns public company information for career page
// @Tags         Public
// @Accept       json
// @Produce      json
// @Param        slug   path      string  true  "Company slug"
// @Success      200    {object}  map[string]interface{}
// @Failure      404    {object}  map[string]interface{}
// @Router       /public/companies/{slug} [get]
func (h *PublicHandler) GetCompanyBySlug(c *gin.Context) {
	slug := c.Param("slug")

	company, err := h.publicService.GetCompanyBySlug(slug)
	if err != nil {
		h.logger.Error("Failed to get company by slug: %v", err)
		c.JSON(apperr.StatusCode(err), gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   dtos.ToPublicCompanyResponse(company),
	})
}

// GetPublishedJobsByCompany godoc
// @Summary      Get published jobs by company slug
// @Description  Returns all published jobs for a company's career page
// @Tags         Public
// @Accept       json
// @Produce      json
// @Param        slug   path      string  true  "Company slug"
// @Success      200    {object}  map[string]interface{}
// @Failure      404    {object}  map[string]interface{}
// @Router       /public/companies/{slug}/jobs [get]
func (h *PublicHandler) GetPublishedJobsByCompany(c *gin.Context) {
	slug := c.Param("slug")

	jobs, err := h.publicService.GetPublishedJobsByCompanySlug(slug)
	if err != nil {
		h.logger.Error("Failed to get jobs for company %s: %v", slug, err)
		c.JSON(apperr.StatusCode(err), gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"jobs":  dtos.ToPublicJobResponseList(jobs),
			"count": len(jobs),
		},
	})
}

// GetPublishedJobByID godoc
// @Summary      Get published job by ID
// @Description  Returns a published job's details for career page
// @Tags         Public
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Job ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /public/jobs/{id} [get]
func (h *PublicHandler) GetPublishedJobByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid job ID",
		})
		return
	}

	job, err := h.publicService.GetPublishedJobByID(uint(id))
	if err != nil {
		h.logger.Error("Failed to get published job %d: %v", id, err)
		c.JSON(apperr.StatusCode(err), gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   dtos.ToPublicJobResponse(job),
	})
}

// ApplyToJob godoc
// @Summary      Apply to a job
// @Description  Submit a job application (public, no auth required)
// @Tags         Public
// @Accept       multipart/form-data
// @Produce      json
// @Param        id           path      int     true   "Job ID"
// @Param        first_name   formData  string  true   "First name"
// @Param        last_name    formData  string  true   "Last name"
// @Param        email        formData  string  true   "Email"
// @Param        phone        formData  string  false  "Phone"
// @Param        linkedin_url formData  string  false  "LinkedIn URL"
// @Param        github_url   formData  string  false  "GitHub URL"
// @Param        cover_letter formData  string  false  "Cover letter"
// @Param        resume       formData  file    false  "Resume file (PDF, DOC, DOCX)"
// @Success      201  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      409  {object}  map[string]interface{}
// @Router       /public/jobs/{id}/apply [post]
func (h *PublicHandler) ApplyToJob(c *gin.Context) {
	// Parse job ID
	jobID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Invalid job ID",
		})
		return
	}

	// Get job first to validate and get company slug for file organization
	_, err = h.publicService.GetPublishedJobByID(uint(jobID))
	if err != nil {
		h.logger.Error("Failed to get job %d: %v", jobID, err)
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Job not found or not accepting applications",
		})
		return
	}
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 12*1024*1024)

	// Parse form data
	var dto dtos.PublicApplicationDTO
	dto.FirstName = strings.TrimSpace(c.PostForm("first_name"))
	dto.LastName = strings.TrimSpace(c.PostForm("last_name"))
	dto.Email = strings.TrimSpace(strings.ToLower(c.PostForm("email")))
	dto.Phone = strings.TrimSpace(c.PostForm("phone"))
	dto.LinkedinURL = strings.TrimSpace(c.PostForm("linkedin_url"))
	dto.GithubURL = strings.TrimSpace(c.PostForm("github_url"))
	dto.CoverLetter = strings.TrimSpace(c.PostForm("cover_letter"))

	// Basic validation
	if dto.FirstName == "" || dto.LastName == "" || dto.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "First name, last name, and email are required",
		})
		return
	}

	var upload *services.ResumeUpload
	file, header, err := c.Request.FormFile("resume")
	if err == nil {
		defer file.Close()
		upload = &services.ResumeUpload{Reader: file, OriginalName: header.Filename, Size: header.Size}
	} else if err != http.ErrMissingFile {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "Invalid multipart request"})
		return
	}

	// Process application
	application, err := h.publicService.ApplyToJob(uint(jobID), dto, upload)
	if err != nil {
		h.logger.Error("Failed to apply to job %d: %v", jobID, err)

		// apperr clasifica el código; los errores internos (no-apperr) caen en 500
		// con mensaje genérico para no filtrar detalles al público.
		code := apperr.StatusCode(err)
		message := err.Error()
		if code == http.StatusInternalServerError {
			message = "Failed to submit application"
		}
		c.JSON(code, gin.H{
			"status":  "error",
			"message": message,
		})
		return
	}

	// Build response
	response := dtos.PublicApplicationResponseDTO{
		ID:        application.ID,
		JobID:     application.JobID,
		AppliedAt: application.AppliedAt,
		Message:   "Your application has been submitted successfully!",
	}

	if application.Job != nil {
		response.JobTitle = application.Job.Title
		if application.Job.Company != nil {
			response.CompanyName = application.Job.Company.Name
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data":   response,
	})
}
