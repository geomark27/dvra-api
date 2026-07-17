package middleware

import (
	"net/http"
	"strings"

	"dvra-api/internal/app/services"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates JWT token and injects user info into context
func AuthMiddleware(jwtService services.JWTService, rbac services.TenantAuthorizer) gin.HandlerFunc {
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
		claims, err := jwtService.ValidateAccessToken(tokenString)
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
		if claims.ContextType == "tenant" {
			role, roleErr := rbac.RoleForMembership(*claims.MembershipID, claims.UserID, *claims.CompanyID)
			if roleErr != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant role is no longer active"})
				c.Abort()
				return
			}
			c.Set("tenant_role_id", role.ID)
			c.Set("role", role.Key)
		} else {
			c.Set("role", claims.Role)
		}
		c.Set("context_type", claims.ContextType)
		if claims.MembershipID != nil {
			c.Set("membership_id", *claims.MembershipID)
		}

		if claims.CompanyID != nil {
			c.Set("company_id", *claims.CompanyID)
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

// OptionalAuth validates token if present, but doesn't require it
func OptionalAuth(jwtService services.JWTService, rbac services.TenantAuthorizer) gin.HandlerFunc {
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
		claims, err := jwtService.ValidateAccessToken(tokenString)
		if err == nil {
			c.Set("user_id", claims.UserID)
			c.Set("email", claims.Email)
			if claims.ContextType == "tenant" {
				role, roleErr := rbac.RoleForMembership(*claims.MembershipID, claims.UserID, *claims.CompanyID)
				if roleErr != nil {
					c.Next()
					return
				}
				c.Set("tenant_role_id", role.ID)
				c.Set("role", role.Key)
			} else {
				c.Set("role", claims.Role)
			}
			c.Set("context_type", claims.ContextType)
			if claims.MembershipID != nil {
				c.Set("membership_id", *claims.MembershipID)
			}

			if claims.CompanyID != nil {
				c.Set("company_id", *claims.CompanyID)
			}
		}

		c.Next()
	}
}
