package handlers

import (
	"net/http"
	"strconv"

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
	candidates, err := h.candidateService.GetAllCandidates()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve candidates"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"candidates": candidates, "count": len(candidates)}})
}

func (h *CandidateHandler) GetCandidate(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	candidate, err := h.candidateService.GetCandidateByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": candidate})
}

func (h *CandidateHandler) CreateCandidate(c *gin.Context) {
	var dto dtos.CreateCandidateDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	candidate, err := h.candidateService.CreateCandidate(dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"status": "success", "data": candidate})
}

func (h *CandidateHandler) UpdateCandidate(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
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
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": candidate})
}

func (h *CandidateHandler) DeleteCandidate(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := h.candidateService.DeleteCandidate(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Candidate deleted"})
}
