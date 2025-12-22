package admin

import (
	"dvra-api/internal/app/dtos"
	"dvra-api/internal/app/services/admin"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// SuperAdminCompaniesHandler handles SuperAdmin company management endpoints
type SuperAdminCompaniesHandler struct {
	service *admin.SuperAdminCompaniesService
}

// NewSuperAdminCompaniesHandler creates a new SuperAdminCompaniesHandler
func NewSuperAdminCompaniesHandler(service *admin.SuperAdminCompaniesService) *SuperAdminCompaniesHandler {
	return &SuperAdminCompaniesHandler{
		service: service,
	}
}

// GetAllCompanies godoc
// @Summary      Listar todas las empresas (SuperAdmin)
// @Description  Obtiene todas las empresas con paginación y filtros
// @Tags         SuperAdmin
// @Accept       json
// @Produce      json
// @Param        page       query  int     false  "Página"       default(1)
// @Param        limit      query  int     false  "Límite"      default(20)
// @Param        search     query  string  false  "Búsqueda"
// @Param        plan_tier  query  string  false  "Filtrar por plan"
// @Success      200        {object}  map[string]interface{}
// @Failure      500        {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /admin/companies [get]
func (h *SuperAdminCompaniesHandler) GetAllCompanies(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	search := c.Query("search")
	planTier := c.Query("plan_tier")

	companies, total, err := h.service.GetAllCompanies(page, limit, search, planTier)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"companies": companies,
		"pagination": gin.H{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// CreateCompany godoc
// @Summary      Crear empresa con admin (SuperAdmin)
// @Description  Crea una nueva empresa y su usuario administrador
// @Tags         SuperAdmin
// @Accept       json
// @Produce      json
// @Param        company  body      dtos.CreateCompanyWithAdminDTO  true  "Datos de empresa y admin"
// @Success      201      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]interface{}
// @Failure      409      {object}  map[string]interface{}
// @Failure      500      {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /admin/companies [post]
func (h *SuperAdminCompaniesHandler) CreateCompany(c *gin.Context) {
	var dto dtos.CreateCompanyWithAdminDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	company, admin, err := h.service.CreateCompanyWithAdmin(dto)
	if err != nil {
		if err.Error() == "admin email already exists" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"company": company,
		"admin":   admin,
		"message": "Company and admin created successfully",
	})
}

func (h *SuperAdminCompaniesHandler) ChangeCompanyPlan(c *gin.Context) {
	companyID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid company ID"})
		return
	}

	var dto dtos.ChangePlanDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.ChangeCompanyPlan(uint(companyID), dto.NewPlan); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Plan updated successfully",
		"company_id": companyID,
		"new_plan":   dto.NewPlan,
	})
}

func (h *SuperAdminCompaniesHandler) SuspendCompany(c *gin.Context) {
	companyID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid company ID"})
		return
	}

	var dto dtos.SuspendCompanyDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.SuspendCompany(uint(companyID), dto.Reason); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Company suspended successfully",
		"company_id": companyID,
		"reason":     dto.Reason,
	})
}

func (h *SuperAdminCompaniesHandler) GetCompanyUsers(c *gin.Context) {
	companyID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid company ID"})
		return
	}

	users, err := h.service.GetCompanyUsers(uint(companyID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"company_id": companyID,
		"users":      users,
		"count":      len(users),
	})
}

// GetGlobalAnalytics godoc
// @Summary      Analíticas globales (SuperAdmin)
// @Description  Obtiene estadísticas globales de la plataforma
// @Tags         SuperAdmin
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /admin/analytics [get]
func (h *SuperAdminCompaniesHandler) GetGlobalAnalytics(c *gin.Context) {
	analytics, err := h.service.GetGlobalAnalytics()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, analytics)
}
