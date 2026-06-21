package transport

import (
	"net/http"
	"strconv"

	"dvra-api/internal/app/dtos"
	"dvra-api/internal/modules/staffing/service"
	"dvra-api/internal/shared/apperr"
	"dvra-api/internal/shared/authctx"

	"github.com/gin-gonic/gin"
)

type PlacementHandler struct {
	svc *service.PlacementService
}

func NewPlacementHandler(svc *service.PlacementService) *PlacementHandler {
	return &PlacementHandler{svc: svc}
}

// GetPlacements godoc
// @Summary      Listar colocaciones (placements)
// @Tags         Placements
// @Produce      json
// @Param        status              query     string  false  "Filtrar por estado (active, ended, suspended)"
// @Param        staffing_client_id  query     int     false  "Filtrar por cliente final"
// @Param        candidate_id        query     int     false  "Filtrar por candidato"
// @Success      200                 {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /placements [get]
func (h *PlacementHandler) GetPlacements(c *gin.Context) {
	var filters dtos.PlacementFilters
	if err := c.ShouldBindQuery(&filters); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	companyID, ok := authctx.CompanyID(c)
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "No company context"})
		return
	}

	placements, err := h.svc.GetByCompanyID(companyID, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve placements"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"placements": dtos.ToPlacementResponseList(placements),
			"count":      len(placements),
		},
	})
}

// GetPlacement godoc
// @Summary      Obtener colocación
// @Tags         Placements
// @Produce      json
// @Param        id   path      int  true  "ID de la colocación"
// @Success      200  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /placements/{id} [get]
func (h *PlacementHandler) GetPlacement(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	placement, err := h.svc.GetByID(uint(id))
	if err != nil {
		c.JSON(apperr.StatusCode(err), gin.H{"error": err.Error()})
		return
	}

	if !authctx.IsSuperAdmin(c) {
		companyID, ok := authctx.CompanyID(c)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "No company context"})
			return
		}
		if placement.CompanyID != companyID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": dtos.ToPlacementResponse(placement)})
}

// CreatePlacement godoc
// @Summary      Crear colocación
// @Description  Coloca a un candidato en un cliente final a partir de una Application en etapa 'hired'.
// @Tags         Placements
// @Accept       json
// @Produce      json
// @Param        placement  body      dtos.CreatePlacementDTO  true  "Datos de la colocación"
// @Success      201        {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /placements [post]
func (h *PlacementHandler) CreatePlacement(c *gin.Context) {
	var dto dtos.CreatePlacementDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	companyID, ok := authctx.CompanyID(c)
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "No company context"})
		return
	}

	placement, err := h.svc.Create(companyID, dto)
	if err != nil {
		c.JSON(apperr.StatusCode(err), gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"status": "success", "data": dtos.ToPlacementResponse(placement)})
}

// UpdatePlacement godoc
// @Summary      Actualizar colocación
// @Tags         Placements
// @Accept       json
// @Produce      json
// @Param        id         path      int                      true  "ID de la colocación"
// @Param        placement  body      dtos.UpdatePlacementDTO  true  "Campos a actualizar"
// @Success      200        {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /placements/{id} [put]
func (h *PlacementHandler) UpdatePlacement(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var companyID uint
	if !authctx.IsSuperAdmin(c) {
		cid, ok := authctx.CompanyID(c)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "No company context"})
			return
		}
		companyID = cid
	}

	var dto dtos.UpdatePlacementDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updated, err := h.svc.Update(uint(id), companyID, dto)
	if err != nil {
		c.JSON(apperr.StatusCode(err), gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": dtos.ToPlacementResponse(updated)})
}

// DeletePlacement godoc
// @Summary      Eliminar colocación
// @Tags         Placements
// @Produce      json
// @Param        id   path      int  true  "ID de la colocación"
// @Success      200  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /placements/{id} [delete]
func (h *PlacementHandler) DeletePlacement(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	var companyID uint
	if !authctx.IsSuperAdmin(c) {
		cid, ok := authctx.CompanyID(c)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "No company context"})
			return
		}
		companyID = cid
	}

	if err := h.svc.Delete(uint(id), companyID); err != nil {
		c.JSON(apperr.StatusCode(err), gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Placement deleted"})
}
