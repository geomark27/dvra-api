package handlers

import (
	"net/http"
	"strconv"

	"dvra-api/internal/app/dtos"
	"dvra-api/internal/app/services"
	"dvra-api/internal/shared/authctx"

	"github.com/gin-gonic/gin"
)

type PlacementHandler struct {
	service services.PlacementService
}

func NewPlacementHandler(service services.PlacementService) *PlacementHandler {
	return &PlacementHandler{service: service}
}

// GetPlacements godoc
// @Summary      Listar colocaciones (placements)
// @Description  Lista las colocaciones de la empresa del token. Requiere un plan con el módulo staffing habilitado.
// @Tags         Placements
// @Accept       json
// @Produce      json
// @Param        status              query     string  false  "Filtrar por estado (active, ended, suspended)"
// @Param        staffing_client_id  query     int     false  "Filtrar por cliente final"
// @Param        candidate_id        query     int     false  "Filtrar por candidato"
// @Success      200                 {object}  map[string]interface{}
// @Failure      403                 {object}  map[string]interface{}
// @Failure      500                 {object}  map[string]interface{}
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

	placements, err := h.service.GetByCompanyID(companyID, filters)
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
// @Description  Obtiene una colocación por ID (validada contra la empresa del token)
// @Tags         Placements
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID de la colocación"
// @Success      200  {object}  map[string]interface{}
// @Failure      403  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /placements/{id} [get]
func (h *PlacementHandler) GetPlacement(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	placement, err := h.service.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
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
// @Description  Coloca a un candidato en un cliente final a partir de una Application en etapa 'hired'. Valida que la application y el cliente pertenezcan a la empresa del token. candidate_id/job_id se copian de la application.
// @Tags         Placements
// @Accept       json
// @Produce      json
// @Param        placement  body      dtos.CreatePlacementDTO  true  "Datos de la colocación (staffing_client_id + application_id requeridos)"
// @Success      201        {object}  map[string]interface{}
// @Failure      400        {object}  map[string]interface{}
// @Failure      403        {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /placements [post]
func (h *PlacementHandler) CreatePlacement(c *gin.Context) {
	var dto dtos.CreatePlacementDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// La colocación siempre se crea en el contexto de una empresa: el service
	// valida que application y cliente final pertenezcan a este company_id.
	companyID, ok := authctx.CompanyID(c)
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "No company context"})
		return
	}

	placement, err := h.service.Create(companyID, dto)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"status": "success", "data": dtos.ToPlacementResponse(placement)})
}

// UpdatePlacement godoc
// @Summary      Actualizar colocación
// @Description  Actualiza datos de contrato, billing y estado de una colocación (el origen application/candidate/cliente es inmutable). Validada contra la empresa del token.
// @Tags         Placements
// @Accept       json
// @Produce      json
// @Param        id         path      int                      true  "ID de la colocación"
// @Param        placement  body      dtos.UpdatePlacementDTO  true  "Campos a actualizar"
// @Success      200        {object}  map[string]interface{}
// @Failure      400        {object}  map[string]interface{}
// @Failure      403        {object}  map[string]interface{}
// @Failure      404        {object}  map[string]interface{}
// @Failure      500        {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /placements/{id} [put]
func (h *PlacementHandler) UpdatePlacement(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	if !authctx.IsSuperAdmin(c) {
		placement, err := h.service.GetByID(uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Placement not found"})
			return
		}
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

	var dto dtos.UpdatePlacementDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updated, err := h.service.Update(uint(id), dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": dtos.ToPlacementResponse(updated)})
}

// DeletePlacement godoc
// @Summary      Eliminar colocación
// @Description  Elimina (soft delete) una colocación (validada contra la empresa del token)
// @Tags         Placements
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ID de la colocación"
// @Success      200  {object}  map[string]interface{}
// @Failure      403  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /placements/{id} [delete]
func (h *PlacementHandler) DeletePlacement(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)

	if !authctx.IsSuperAdmin(c) {
		placement, err := h.service.GetByID(uint(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Placement not found"})
			return
		}
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

	if err := h.service.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "message": "Placement deleted"})
}
