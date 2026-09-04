package transport

import (
	"net/http"

	"dvra-api/internal/modules/iam/service"
	"dvra-api/internal/shared/apperr"
	"dvra-api/internal/shared/authctx"
	"dvra-api/internal/shared/permissions"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service *service.AuthService
}

func NewAuthHandler(service *service.AuthService) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

// Register godoc
// @Summary      Registrar nuevo usuario (deprecated)
// @Description  Registra un nuevo usuario. Usar RegisterCompany en su lugar
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        user  body      service.RegisterDTO  true  "Datos del usuario"
// @Success      201   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]interface{}
// @Failure      409   {object}  map[string]interface{}
// @Router       /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var dto service.RegisterDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.service.Register(&dto)
	if err != nil {
		c.JSON(apperr.StatusCode(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response)
}

// Login godoc
// @Summary      Iniciar sesión
// @Description  Autentica un usuario y retorna tokens JWT
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        credentials  body      service.LoginDTO  true  "Credenciales de acceso"
// @Success      200          {object}  map[string]interface{}
// @Failure      400          {object}  map[string]interface{}
// @Failure      401          {object}  map[string]interface{}
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var dto service.LoginDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.service.LoginWithCompanies(&dto)
	if err != nil {
		c.JSON(apperr.StatusCode(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// RefreshToken godoc
// @Summary      Refrescar token de acceso
// @Description  Genera un nuevo access token usando el refresh token
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        token  body      service.RefreshTokenDTO  true  "Refresh token"
// @Success      200    {object}  map[string]interface{}
// @Failure      400    {object}  map[string]interface{}
// @Failure      401    {object}  map[string]interface{}
// @Router       /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var dto service.RefreshTokenDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.service.RefreshToken(&dto)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetMe godoc
// @Summary      Obtener usuario actual
// @Description  Retorna la información del usuario autenticado
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /auth/me [get]
func (h *AuthHandler) GetMe(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	response, err := h.service.GetMe(userID.(uint))
	if err != nil {
		c.JSON(apperr.StatusCode(err), gin.H{"error": err.Error()})
		return
	}

	role := authctx.Role(c)
	response.Role = role
	context := &service.AuthContextResponse{Type: authctx.ContextType(c), Role: role}
	if companyID, ok := authctx.CompanyID(c); ok {
		context.CompanyID = &companyID
	}
	if membershipID, ok := authctx.MembershipID(c); ok {
		context.MembershipID = &membershipID
	}
	response.ActiveContext = context
	if authctx.IsPlatformContext(c) {
		response.Permissions = permissions.ForPlatform(role)
	} else if authctx.IsTenantContext(c) {
		roleID, ok := authctx.TenantRoleID(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant role missing"})
			return
		}
		tenantRole, rolePermissions, authErr := h.service.GetTenantAuthorization(roleID)
		if authErr != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant role is no longer active"})
			return
		}
		response.TenantRole = service.ToTenantRoleResponse(tenantRole)
		response.Permissions = rolePermissions
		context.TenantRoleID = &roleID
		context.TenantRole = service.ToTenantRoleResponse(tenantRole)
	}

	c.JSON(http.StatusOK, response)
}

func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var dto service.ChangePasswordDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.ChangePassword(userID.(uint), &dto)
	if err != nil {
		c.JSON(apperr.StatusCode(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	// En un sistema con JWT stateless, el logout es manejado por el cliente
	// descartando los tokens. En el futuro, se puede implementar una blacklist
	// de tokens en Redis para invalidar tokens antes de su expiración.
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

// RegisterCompany godoc
// @Summary      Registrar empresa con admin
// @Description  Crea una nueva empresa y su usuario administrador
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        company  body      service.RegisterCompanyDTO  true  "Datos de empresa y admin"
// @Success      201      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]interface{}
// @Failure      409      {object}  map[string]interface{}
// @Router       /auth/register-company [post]
func (h *AuthHandler) RegisterCompany(c *gin.Context) {
	var dto service.RegisterCompanyDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.service.RegisterCompany(&dto)
	if err != nil {
		c.JSON(apperr.StatusCode(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response)
}

// SwitchCompany godoc
// @Summary      Switch active tenant
// @Description  Issues new access and refresh tokens for an active tenant membership
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        body body service.SwitchCompanyDTO true "Tenant context"
// @Success      200 {object} service.SwitchCompanyResponseDTO
// @Security     BearerAuth
// @Router       /auth/switch-company [post]
func (h *AuthHandler) SwitchCompany(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var dto service.SwitchCompanyDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.service.SwitchCompany(userID.(uint), &dto)
	if err != nil {
		c.JSON(apperr.StatusCode(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// SwitchToPlatform godoc
// @Summary      Switch to landlord context
// @Description  Issues new access and refresh tokens for an active platform membership
// @Tags         Authentication
// @Produce      json
// @Success      200 {object} service.SwitchPlatformResponseDTO
// @Security     BearerAuth
// @Router       /auth/switch-platform [post]
func (h *AuthHandler) SwitchToPlatform(c *gin.Context) {
	userID, ok := authctx.UserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	response, err := h.service.SwitchToPlatform(userID)
	if err != nil {
		c.JSON(apperr.StatusCode(err), gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, response)
}

func (h *AuthHandler) GetMyCompanies(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	companies, err := h.service.GetUserCompanies(userID.(uint))
	if err != nil {
		c.JSON(apperr.StatusCode(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"companies": companies})
}
