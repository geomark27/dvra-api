package middleware

import (
	"net/http"

	"dvra-api/internal/shared/authctx"
	"dvra-api/internal/shared/permissions"

	"github.com/gin-gonic/gin"
)

// RequirePermission autoriza una acción tenant solo si el contexto y rol del
// token tienen el permiso indicado. Debe aplicarse después de AuthMiddleware:
//
//	jobs.POST("", middleware.RequirePermission(permissions.JobsCreate), h.CreateJob)
func RequirePermission(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role := authctx.Role(c)
		if role == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		if !authctx.IsTenantContext(c) || !permissions.Can(role, permission) {
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
