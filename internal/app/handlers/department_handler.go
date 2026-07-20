package handlers

import (
	"net/http"
	"strconv"

	"dvra-api/internal/app/dtos"
	"dvra-api/internal/app/services"
	"dvra-api/internal/shared/apperr"
	"dvra-api/internal/shared/authctx"

	"github.com/gin-gonic/gin"
)

type DepartmentHandler struct{ service services.DepartmentService }

func NewDepartmentHandler(service services.DepartmentService) *DepartmentHandler {
	return &DepartmentHandler{service: service}
}

// ListDepartments godoc
// @Summary Listar departamentos
// @Description Devuelve la estructura organizacional del tenant en formato plano o árbol
// @Tags Departments
// @Produce json
// @Param tree query bool false "Construir jerarquía"
// @Param include_inactive query bool false "Incluir departamentos inactivos"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /departments [get]
func (h *DepartmentHandler) ListDepartments(c *gin.Context) {
	h.listDepartments(c, c.Query("tree") == "true")
}

// GetDepartmentTree godoc
// @Summary Obtener árbol de departamentos
// @Tags Departments
// @Produce json
// @Param include_inactive query bool false "Incluir departamentos inactivos"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /departments/tree [get]
func (h *DepartmentHandler) GetDepartmentTree(c *gin.Context) {
	h.listDepartments(c, true)
}

func (h *DepartmentHandler) listDepartments(c *gin.Context, tree bool) {
	companyID, ok := authctx.CompanyID(c)
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "No company context"})
		return
	}
	departments, err := h.service.ListDepartments(companyID, c.Query("include_inactive") == "true")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve departments"})
		return
	}
	var data interface{} = dtos.ToDepartmentResponseList(departments)
	if tree {
		data = dtos.ToDepartmentTree(departments)
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": data})
}

// GetDepartment godoc
// @Summary Obtener departamento
// @Tags Departments
// @Produce json
// @Param id path int true "Department ID"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /departments/{id} [get]
func (h *DepartmentHandler) GetDepartment(c *gin.Context) {
	id, companyID, ok := departmentRequestScope(c)
	if !ok {
		return
	}
	department, err := h.service.GetDepartment(id, companyID)
	if err != nil {
		c.JSON(apperr.StatusCode(err), gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": dtos.ToDepartmentResponse(department)})
}

// CreateDepartment godoc
// @Summary Crear departamento o subdepartamento
// @Tags Departments
// @Accept json
// @Produce json
// @Param department body dtos.CreateDepartmentDTO true "Departamento"
// @Success 201 {object} map[string]interface{}
// @Security BearerAuth
// @Router /departments [post]
func (h *DepartmentHandler) CreateDepartment(c *gin.Context) {
	companyID, ok := authctx.CompanyID(c)
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "No company context"})
		return
	}
	var dto dtos.CreateDepartmentDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	department, err := h.service.CreateDepartment(companyID, dto)
	if err != nil {
		c.JSON(apperr.StatusCode(err), gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"status": "success", "data": dtos.ToDepartmentResponse(department)})
}

// UpdateDepartment godoc
// @Summary Actualizar departamento
// @Description Use parent_department_id=0 para moverlo a la raíz y manager_user_id=0 para quitar responsable
// @Tags Departments
// @Accept json
// @Produce json
// @Param id path int true "Department ID"
// @Param department body dtos.UpdateDepartmentDTO true "Cambios"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /departments/{id} [put]
func (h *DepartmentHandler) UpdateDepartment(c *gin.Context) {
	id, companyID, ok := departmentRequestScope(c)
	if !ok {
		return
	}
	var dto dtos.UpdateDepartmentDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	department, err := h.service.UpdateDepartment(id, companyID, dto)
	if err != nil {
		c.JSON(apperr.StatusCode(err), gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": dtos.ToDepartmentResponse(department)})
}

// DeactivateDepartment godoc
// @Summary Desactivar departamento
// @Tags Departments
// @Produce json
// @Param id path int true "Department ID"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /departments/{id} [delete]
func (h *DepartmentHandler) DeactivateDepartment(c *gin.Context) {
	id, companyID, ok := departmentRequestScope(c)
	if !ok {
		return
	}
	department, err := h.service.DeactivateDepartment(id, companyID)
	if err != nil {
		c.JSON(apperr.StatusCode(err), gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": dtos.ToDepartmentResponse(department)})
}

func departmentRequestScope(c *gin.Context) (uint, uint, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid department ID"})
		return 0, 0, false
	}
	companyID, ok := authctx.CompanyID(c)
	if !ok {
		c.JSON(http.StatusForbidden, gin.H{"error": "No company context"})
		return 0, 0, false
	}
	return uint(id), companyID, true
}
