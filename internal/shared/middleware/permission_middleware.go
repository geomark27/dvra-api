package middleware

import (
	"net/http"

	"dvra-api/internal/shared/authctx"
	"dvra-api/internal/shared/permissions"

	"github.com/gin-gonic/gin"
)

// RequirePermission autoriza la acción solo si el rol del token tiene el
// permiso indicado (ver internal/shared/permissions). SuperAdmin pasa siempre.
// Debe aplicarse después de AuthMiddleware:
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

		if !permissions.Can(role, permission) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}
