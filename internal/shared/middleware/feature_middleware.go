package middleware

import (
	"net/http"

	"dvra-api/internal/app/services"
	"dvra-api/internal/shared/authctx"

	"github.com/gin-gonic/gin"
)

// RequireFeature autoriza la acción solo si el plan de la empresa habilita el
// feature indicado (entitlement). Es ortogonal a RequirePermission, que valida
// el rol: una acción de staffing requiere AMBOS (el plan lo incluye Y el rol lo
// permite). SuperAdmin no tiene empresa, por lo que pasa siempre.
// Debe aplicarse después de AuthMiddleware:
//
//	staffing.Use(middleware.RequireFeature(planService, "staffing"))
func RequireFeature(planService services.PlanService, feature string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if authctx.IsSuperAdmin(c) {
			c.Next()
			return
		}

		companyID, ok := authctx.CompanyID(c)
		if !ok {
			c.JSON(http.StatusForbidden, gin.H{"error": "No company context"})
			c.Abort()
			return
		}

		enabled, err := planService.CompanyHasFeature(companyID, feature)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to verify plan entitlement"})
			c.Abort()
			return
		}
		if !enabled {
			c.JSON(http.StatusForbidden, gin.H{"error": "Your plan does not include this feature"})
			c.Abort()
			return
		}

		c.Next()
	}
}
