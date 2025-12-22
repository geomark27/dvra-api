package handlers

import (
	"net/http"

	"dvra-api/internal/app/dtos"
	"dvra-api/internal/app/services"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service *services.AuthService
}

func NewAuthHandler(service *services.AuthService) *AuthHandler {
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
// @Param        user  body      dtos.RegisterDTO  true  "Datos del usuario"
// @Success      201   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]interface{}
// @Failure      409   {object}  map[string]interface{}
// @Router       /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var dto dtos.RegisterDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.service.Register(&dto)
	if err != nil {
		if err == services.ErrEmailExists {
			c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
// @Param        credentials  body      dtos.LoginDTO  true  "Credenciales de acceso"
// @Success      200          {object}  map[string]interface{}
// @Failure      400          {object}  map[string]interface{}
// @Failure      401          {object}  map[string]interface{}
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var dto dtos.LoginDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.service.LoginWithCompanies(&dto)
	if err != nil {
		if err == services.ErrInvalidCredentials {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
// @Param        token  body      dtos.RefreshTokenDTO  true  "Refresh token"
// @Success      200    {object}  map[string]interface{}
// @Failure      400    {object}  map[string]interface{}
// @Failure      401    {object}  map[string]interface{}
// @Router       /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var dto dtos.RefreshTokenDTO
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
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var dto dtos.ChangePasswordDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.service.ChangePassword(userID.(uint), &dto)
	if err != nil {
		if err == services.ErrInvalidPassword {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid old password"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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

func (h *AuthHandler) SuperAdminLogin(c *gin.Context) {
	var dto dtos.SuperAdminLoginDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.service.SuperAdminLogin(&dto)
	if err != nil {
		if err == services.ErrNotSuperAdmin {
			c.JSON(http.StatusForbidden, gin.H{"error": "User is not a superadmin"})
			return
		}
		if err == services.ErrInvalidCredentials {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// RegisterCompany godoc
// @Summary      Registrar empresa con admin
// @Description  Crea una nueva empresa y su usuario administrador
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        company  body      dtos.RegisterCompanyDTO  true  "Datos de empresa y admin"
// @Success      201      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]interface{}
// @Failure      409      {object}  map[string]interface{}
// @Router       /auth/register-company [post]
func (h *AuthHandler) RegisterCompany(c *gin.Context) {
	var dto dtos.RegisterCompanyDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.service.RegisterCompany(&dto)
	if err != nil {
		if err == services.ErrEmailExists {
			c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response)
}

func (h *AuthHandler) SwitchCompany(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var dto dtos.SwitchCompanyDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.service.SwitchCompany(userID.(uint), &dto)
	if err != nil {
		if err == services.ErrNoMembership {
			c.JSON(http.StatusForbidden, gin.H{"error": "You do not belong to this company"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"companies": companies})
}
