package middleware

import (
	"net/http"

	"dvra-api/internal/app/services"
	"dvra-api/internal/shared/authctx"
	"dvra-api/internal/shared/permissions"

	"github.com/gin-gonic/gin"
)

// RequirePermission autoriza una acción tenant solo si el contexto y rol del
// token tienen el permiso indicado. Debe aplicarse después de AuthMiddleware:
//
//	jobs.POST("", middleware.RequirePermission(rbac, permissions.JobsCreate), h.CreateJob)
func RequirePermission(rbac services.TenantAuthorizer, permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !authctx.IsTenantContext(c) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		roleID, ok := authctx.TenantRoleID(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		allowed, err := rbac.Can(roleID, permission)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to evaluate permissions"})
			c.Abort()
			return
		}
		if !allowed {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func RequirePlatformPermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role := authctx.Role(c)
		if !authctx.IsPlatformContext(c) || !permissions.PlatformCan(role, permission) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient platform permissions"})
			c.Abort()
			return
		}
		c.Next()
	}
}
