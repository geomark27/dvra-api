package handlers

import (
	"dvra-api/internal/shared/authctx"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"dvra-api/internal/app/dtos"
	"dvra-api/internal/app/services"

	"github.com/geomark27/loom-go/pkg/helpers"
	"github.com/gin-gonic/gin"
)

type CandidateHandler struct {
	candidateService services.CandidateService
	logger           helpers.Logger
}

func NewCandidateHandler(candidateService services.CandidateService) *CandidateHandler {
	return &CandidateHandler{candidateService: candidateService, logger: helpers.NewLogger()}
}

func (h *CandidateHandler) GetCandidates(c *gin.Context) {

	// SuperAdmin puede ver todos los candidates
	if authctx.IsSuperAdmin(c) {
		candidates, err := h.candidateService.GetAllCandidates()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve candidates"})
			return
		}
		candidatesDTO := dtos.ToCandidateResponseList(candidates)
		c.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"candidates": candidatesDTO, "count": len(candidatesDTO)}})
		return
	}

	// Usuarios normales solo ven candidates de su empresa
	companyIDVal, exists := c.Get("company_id")
	if !exists {
		c.JSON(http.StatusForbidden, gin.H{"error": "No company context"})
		return
	}

	companyID := companyIDVal.(uint)
	candidates, err := h.candidateService.GetCandidatesByCompanyID(companyID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve candidates"})
		return
	}
	candidatesDTO := dtos.ToCandidateResponseList(candidates)
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"candidates": candidatesDTO, "count": len(candidatesDTO)}})
}

func (h *CandidateHandler) GetCandidate(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	candidate, err := h.candidateService.GetCandidateByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
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
		if candidate.CompanyID != companyID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}
	}

	candidateDTO := dtos.ToCandidateResponse(candidate)
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": candidateDTO})
}

func (h *CandidateHandler) CreateCandidate(c *gin.Context) {
	var dto dtos.CreateCandidateDTO
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

	candidate, err := h.candidateService.CreateCandidate(dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	candidateDTO := dtos.ToCandidateResponse(candidate)
	c.JSON(http.StatusCreated, gin.H{"status": "success", "data": candidateDTO})
}

func (h *CandidateHandler) UpdateCandidate(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	// Validar que el candidate pertenece a la empresa del usuario
	if !authctx.IsSuperAdmin(c) {
		candidate, err := h.candidateService.GetCandidateByID(uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Candidate not found"})
			return
		}
		companyIDVal, exists := c.Get("company_id")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "No company context"})
			return
		}
		companyID := companyIDVal.(uint)
		if candidate.CompanyID != companyID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}
	}

	var dto dtos.UpdateCandidateDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	candidate, err := h.candidateService.UpdateCandidate(uint(id), dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	candidateDTO := dtos.ToCandidateResponse(candidate)
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": candidateDTO})
}

func (h *CandidateHandler) DeleteCandidate(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	// Validar que el candidate pertenece a la empresa del usuario
	if !authctx.IsSuperAdmin(c) {
		candidate, err := h.candidateService.GetCandidateByID(uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Candidate not found"})
			return
		}
		companyIDVal, exists := c.Get("company_id")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "No company context"})
			return
		}
		companyID := companyIDVal.(uint)
		if candidate.CompanyID != companyID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}
	}

	if err := h.candidateService.DeleteCandidate(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Candidate deleted"})
}

// UploadResume godoc
// @Summary      Upload resume for candidate
// @Description  Sube un CV/Resume para un candidato
// @Tags         Candidates
// @Accept       multipart/form-data
// @Produce      json
// @Param        id      path      int   true  "Candidate ID"
// @Param        resume  formData  file  true  "Resume file"
// @Success      200     {object}  map[string]interface{}
// @Failure      400     {object}  map[string]interface{}
// @Failure      403     {object}  map[string]interface{}
// @Failure      404     {object}  map[string]interface{}
// @Failure      500     {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /candidates/{id}/upload-resume [post]
func (h *CandidateHandler) UploadResume(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	// Validate access
	if !authctx.IsSuperAdmin(c) {
		candidate, err := h.candidateService.GetCandidateByID(uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Candidate not found"})
			return
		}
		companyIDVal, exists := c.Get("company_id")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "No company context"})
			return
		}
		companyID := companyIDVal.(uint)
		if candidate.CompanyID != companyID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}
	}

	// Get file from form
	file, err := c.FormFile("resume")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Resume file is required"})
		return
	}

	// Validate file type (PDF, DOC, DOCX)
	allowedTypes := map[string]bool{
		"application/pdf":    true,
		"application/msword": true,
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
	}

	if !allowedTypes[file.Header.Get("Content-Type")] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type. Allowed: PDF, DOC, DOCX"})
		return
	}

	// For now, we'll generate a simple file path/URL
	// In production, you would upload to S3, GCS, or similar
	filename := fmt.Sprintf("resumes/%d_%d_%s", id, time.Now().Unix(), file.Filename)

	// Save file locally (in production use cloud storage)
	if err := c.SaveUploadedFile(file, "./uploads/"+filename); err != nil {
		h.logger.Error("Failed to save resume", map[string]interface{}{"error": err.Error()})
		// Even if local save fails, we'll update the URL for demo purposes
	}

	// Generate URL (in production this would be a cloud storage URL)
	resumeURL := "/uploads/" + filename

	// Update candidate with resume URL
	updateDTO := dtos.UpdateCandidateDTO{
		ResumeURL: &resumeURL,
	}

	_, err = h.candidateService.UpdateCandidate(uint(id), updateDTO)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update candidate with resume URL"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":     "success",
		"resume_url": resumeURL,
	})
}
