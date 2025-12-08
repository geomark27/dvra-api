package handlers

import (
	"net/http"
	"strconv"

	"dvra-api/internal/app/dtos"
	"dvra-api/internal/app/services"

	"github.com/geomark27/loom-go/pkg/helpers"
	"github.com/gin-gonic/gin"
)

// CompanyHandler handles company-related routes
type CompanyHandler struct {
	companyService services.CompanyService
	logger         helpers.Logger
}

// NewCompanyHandler creates a new CompanyHandler instance
func NewCompanyHandler(companyService services.CompanyService) *CompanyHandler {
	return &CompanyHandler{
		companyService: companyService,
		logger:         helpers.NewLogger(),
	}
}

// GetCompanies retrieves all companies
func (h *CompanyHandler) GetCompanies(c *gin.Context) {
	role, _ := c.Get("role")

	// SuperAdmin puede ver todas las empresas
	if role == "superadmin" {
		companies, err := h.companyService.GetAllCompanies()
		if err != nil {
			h.logger.Error("Failed to get companies", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve companies"})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Companies retrieved successfully",
			"data": gin.H{
				"companies": companies,
				"count":     len(companies),
			},
		})
		return
	}

	// Usuarios normales solo ven su empresa
	companyIDVal, exists := c.Get("company_id")
	if !exists {
		c.JSON(http.StatusForbidden, gin.H{"error": "No company context"})
		return
	}

	companyID := companyIDVal.(uint)
	company, err := h.companyService.GetCompanyByID(companyID)
	if err != nil {
		h.logger.Error("Failed to get company", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve company"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Company retrieved successfully",
		"data": gin.H{
			"companies": []interface{}{company},
			"count":     1,
		},
	})
}

// GetCompany retrieves a company by ID
func (h *CompanyHandler) GetCompany(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid company ID"})
		return
	}

	role, _ := c.Get("role")

	// Validar acceso: SuperAdmin o miembro de la empresa
	if role != "superadmin" {
		companyIDVal, exists := c.Get("company_id")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "No company context"})
			return
		}
		companyID := companyIDVal.(uint)
		if uint(id) != companyID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}
	}

	company, err := h.companyService.GetCompanyByID(uint(id))
	if err != nil {
		if err.Error() == "company not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Company not found"})
			return
		}
		h.logger.Error("Failed to get company", "error", err, "company_id", id)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve company"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Company retrieved successfully",
		"data":    company,
	})
}

// CreateCompany creates a new company
func (h *CompanyHandler) CreateCompany(c *gin.Context) {
	var dto dtos.CreateCompanyDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	if errors := helpers.ValidateStruct(&dto); len(errors) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": errors})
		return
	}

	company, err := h.companyService.CreateCompany(dto)
	if err != nil {
		h.logger.Error("Failed to create company", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Company created successfully",
		"data":    company,
	})
}

// UpdateCompany updates an existing company
func (h *CompanyHandler) UpdateCompany(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid company ID"})
		return
	}

	// Validar acceso: SuperAdmin o miembro de esa empresa
	role, _ := c.Get("role")
	if role != "superadmin" {
		companyIDVal, exists := c.Get("company_id")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "No company context"})
			return
		}
		companyID := companyIDVal.(uint)
		if uint(id) != companyID {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}
	}

	var dto dtos.UpdateCompanyDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	if errors := helpers.ValidateStruct(&dto); len(errors) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": errors})
		return
	}

	company, err := h.companyService.UpdateCompany(uint(id), dto)
	if err != nil {
		if err.Error() == "company not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Company not found"})
			return
		}
		h.logger.Error("Failed to update company", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update company"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Company updated successfully",
		"data":    company,
	})
}

// DeleteCompany deletes a company
func (h *CompanyHandler) DeleteCompany(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid company ID"})
		return
	}

	// Validar acceso: Solo SuperAdmin puede eliminar empresas
	role, _ := c.Get("role")
	if role != "superadmin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Only superadmin can delete companies"})
		return
	}

	if err := h.companyService.DeleteCompany(uint(id)); err != nil {
		if err.Error() == "company not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Company not found"})
			return
		}
		h.logger.Error("Failed to delete company", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete company"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Company deleted successfully",
	})
}
