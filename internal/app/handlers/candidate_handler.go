package handlers

import (
	"dvra-api/internal/shared/authctx"
	"mime"
	"net/http"
	"strconv"

	"dvra-api/internal/app/dtos"
	"dvra-api/internal/app/services"
	"dvra-api/internal/shared/apperr"

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

	// Tenant users only see candidates allowed by their active membership.
	scope, ok := authctx.TenantScope(c)
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "No company context"})
		return
	}

	candidates, err := h.candidateService.GetCandidatesWithScope(scope)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve candidates"})
		return
	}
	candidatesDTO := dtos.ToCandidateResponseList(candidates)
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"candidates": candidatesDTO, "count": len(candidatesDTO)}})
}

func (h *CandidateHandler) GetCandidate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid candidate ID"})
		return
	}
	scope, _ := authctx.TenantScope(c)
	candidate, err := h.candidateService.GetCandidateByID(uint(id), scope)
	if err != nil {
		c.JSON(apperr.StatusCode(err), gin.H{"error": err.Error()})
		return
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

	companyID, ok := authctx.CompanyID(c)
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "No company context"})
		return
	}
	dto.CompanyID = companyID

	candidate, err := h.candidateService.CreateCandidate(dto)
	if err != nil {
		c.JSON(apperr.StatusCode(err), gin.H{"error": err.Error()})
		return
	}
	candidateDTO := dtos.ToCandidateResponse(candidate)
	c.JSON(http.StatusCreated, gin.H{"status": "success", "data": candidateDTO})
}

func (h *CandidateHandler) UpdateCandidate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid candidate ID"})
		return
	}
	scope, _ := authctx.TenantScope(c)

	var dto dtos.UpdateCandidateDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	candidate, err := h.candidateService.UpdateCandidate(uint(id), scope, dto)
	if err != nil {
		c.JSON(apperr.StatusCode(err), gin.H{"error": err.Error()})
		return
	}
	candidateDTO := dtos.ToCandidateResponse(candidate)
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": candidateDTO})
}

func (h *CandidateHandler) DeleteCandidate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid candidate ID"})
		return
	}
	scope, _ := authctx.TenantScope(c)

	if err := h.candidateService.DeleteCandidate(uint(id), scope); err != nil {
		c.JSON(apperr.StatusCode(err), gin.H{"error": err.Error()})
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
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid candidate ID"})
		return
	}
	scope, _ := authctx.TenantScope(c)
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 12*1024*1024)
	header, err := c.FormFile("resume")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Resume file is required"})
		return
	}
	src, err := header.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid resume file"})
		return
	}
	defer src.Close()
	stored, err := h.candidateService.StoreResume(uint(id), scope, services.ResumeUpload{
		Reader: src, OriginalName: header.Filename, Size: header.Size,
	})
	if err != nil {
		c.JSON(apperr.StatusCode(err), gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   gin.H{"resume_url": "/api/v1/candidates/" + strconv.FormatUint(id, 10) + "/resume", "size_bytes": stored.SizeBytes},
	})
}

// DownloadResume godoc
// @Summary      Download candidate resume
// @Description  Descarga privada del CV de un candidato del tenant activo
// @Tags         Candidates
// @Produce      application/octet-stream
// @Param        id  path  int  true  "Candidate ID"
// @Success      200 {file} binary
// @Failure      403 {object} map[string]interface{}
// @Failure      404 {object} map[string]interface{}
// @Security     BearerAuth
// @Router       /candidates/{id}/resume [get]
func (h *CandidateHandler) DownloadResume(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid candidate ID"})
		return
	}
	scope, _ := authctx.TenantScope(c)
	stored, reader, err := h.candidateService.OpenResume(uint(id), scope)
	if err != nil {
		c.JSON(apperr.StatusCode(err), gin.H{"error": err.Error()})
		return
	}
	defer reader.Close()
	disposition := mime.FormatMediaType("attachment", map[string]string{"filename": stored.OriginalName})
	c.Header("Content-Disposition", disposition)
	c.DataFromReader(http.StatusOK, stored.SizeBytes, stored.ContentType, reader, nil)
}
