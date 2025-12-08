package handlers

import (
	"net/http"
	"strconv"

	"dvra-api/internal/app/dtos"
	"dvra-api/internal/app/services"

	"github.com/geomark27/loom-go/pkg/helpers"
	"github.com/gin-gonic/gin"
)

// UserHandler handles user-related routes
type UserHandler struct {
	userService services.UserService
	logger      helpers.Logger
}

// NewUserHandler creates a new UserHandler instance
func NewUserHandler(userService services.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
		logger:      helpers.NewLogger(),
	}
}

// GetUsers retrieves all users
func (h *UserHandler) GetUsers(c *gin.Context) {
	role, _ := c.Get("role")

	// SuperAdmin puede ver todos los usuarios
	if role == "superadmin" {
		users, err := h.userService.GetAllUsers()
		if err != nil {
			h.logger.Error("Failed to get users", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to retrieve users",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Users retrieved successfully",
			"data": gin.H{
				"users": users,
				"count": len(users),
			},
		})
		return
	}

	// Usuarios normales solo ven usuarios de su empresa
	companyIDVal, exists := c.Get("company_id")
	if !exists {
		c.JSON(http.StatusForbidden, gin.H{"error": "No company context"})
		return
	}

	companyID := companyIDVal.(uint)
	users, err := h.userService.GetUsersByCompanyID(companyID)
	if err != nil {
		h.logger.Error("Failed to get users", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve users",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Users retrieved successfully",
		"data": gin.H{
			"users": users,
			"count": len(users),
		},
	})
}

// GetUser retrieves a user by ID
func (h *UserHandler) GetUser(c *gin.Context) {
	idParam, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}
	id := uint(idParam)

	user, err := h.userService.GetUserByID(id)
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "User not found",
			})
			return
		}
		h.logger.Error("Failed to get user", "error", err, "user_id", id)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve user",
		})
		return
	}

	// Validar acceso: SuperAdmin puede ver cualquier user
	// Usuarios normales solo pueden ver users de su empresa
	role, _ := c.Get("role")
	if role != "superadmin" {
		companyIDVal, exists := c.Get("company_id")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "No company context"})
			return
		}
		companyID := companyIDVal.(uint)

		// Verificar si el user pertenece a la empresa
		userBelongsToCompany := false
		for _, membership := range user.Memberships {
			if membership.CompanyID != nil && *membership.CompanyID == companyID {
				userBelongsToCompany = true
				break
			}
		}

		if !userBelongsToCompany {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "User retrieved successfully",
		"data":    user,
	})
}

// CreateUser creates a new user
func (h *UserHandler) CreateUser(c *gin.Context) {
	var dto dtos.CreateUserDTO

	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}
	// Validate DTO
	if errors := helpers.ValidateStruct(&dto); len(errors) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation failed",
			"details": errors,
		})
		return
	}

	user, err := h.userService.CreateUser(dto)
	if err != nil {
		h.logger.Error("Failed to create user", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create user",
		})
		return
	}
	h.logger.Info("User created successfully", "user_id", user.ID)

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "User created successfully",
		"data":    user,
	})
}

// UpdateUser updates an existing user
func (h *UserHandler) UpdateUser(c *gin.Context) {
	idParam, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}
	id := uint(idParam)

	var dto dtos.UpdateUserDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}
	// Validate DTO
	if errors := helpers.ValidateStruct(&dto); len(errors) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation failed",
			"details": errors,
		})
		return
	}

	user, err := h.userService.UpdateUser(id, dto)
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "User not found",
			})
			return
		}
		h.logger.Error("Failed to update user", "error", err, "user_id", id)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update user",
		})
		return
	}
	h.logger.Info("User updated successfully", "user_id", id)

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "User updated successfully",
		"data":    user,
	})
}

// DeleteUser deletes a user
func (h *UserHandler) DeleteUser(c *gin.Context) {
	idParam, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}
	id := uint(idParam)

	if err := h.userService.DeleteUser(id); err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "User not found",
			})
			return
		}
		h.logger.Error("Failed to delete user", "error", err, "user_id", id)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete user",
		})
		return
	}
	h.logger.Info("User deleted successfully", "user_id", id)

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "User deleted successfully",
	})
}
