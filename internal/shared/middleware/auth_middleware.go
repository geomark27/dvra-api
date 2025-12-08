package middleware

import (
	"net/http"
	"strings"

	"dvra-api/internal/app/services"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates JWT token and injects user info into context
func AuthMiddleware(jwtService services.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Check if it's a Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Validate token
		claims, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			if err == services.ErrExpiredToken {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Token has expired"})
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			}
			c.Abort()
			return
		}

		// Inject claims into context
		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)

		if claims.CompanyID != nil {
			c.Set("company_id", *claims.CompanyID)
		}

		c.Next()
	}
}

// RequireRole checks if user has required role level
func RequireRole(minLevel int) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		userLevel := getRoleLevel(role.(string))
		if userLevel < minLevel {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireCompany checks if user has a company context
func RequireCompany() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, exists := c.Get("company_id")
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "Company context required"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireSuperAdmin validates that user is SuperAdmin (no company_id, role = superadmin)
func RequireSuperAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Must be authenticated
		_, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			c.Abort()
			return
		}

		// Get role and company_id
		role, roleExists := c.Get("role")
		companyID, hasCompany := c.Get("company_id")

		// SuperAdmin must have role = "superadmin"
		if !roleExists || role != "superadmin" {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "SuperAdmin access required",
			})
			c.Abort()
			return
		}

		// SuperAdmin must NOT have company context (global access)
		if hasCompany && companyID != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "SuperAdmin routes require no company context",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// OptionalAuth validates token if present, but doesn't require it
func OptionalAuth(jwtService services.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Next()
			return
		}

		tokenString := parts[1]
		claims, err := jwtService.ValidateToken(tokenString)
		if err == nil {
			c.Set("user_id", claims.UserID)
			c.Set("email", claims.Email)
			c.Set("role", claims.Role)

			if claims.CompanyID != nil {
				c.Set("company_id", *claims.CompanyID)
			}
		}

		c.Next()
	}
}

// getRoleLevel returns the numeric level for a role
func getRoleLevel(role string) int {
	levels := map[string]int{
		"superadmin":     100,
		"admin":          50,
		"recruiter":      30,
		"hiring_manager": 20,
		"user":           10,
	}

	if level, exists := levels[role]; exists {
		return level
	}
	return 0
}
