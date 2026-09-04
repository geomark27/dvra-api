package middleware

import (
	"errors"
	"net/http"
	"strings"

	"dvra-api/internal/shared/authctx"
	"dvra-api/internal/shared/identity"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware valida el token y deja la sesión en el contexto.
//
// En contexto tenant resuelve el rol contra la base en CADA petición, no desde
// el token: así revocar un rol o suspender una membresía surte efecto en la
// llamada siguiente, sin esperar a que el token caduque.
func AuthMiddleware(auth Authenticator) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		token, ok := bearerToken(header)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		claims, err := auth.Session(token)
		if err != nil {
			if errors.Is(err, identity.ErrExpiredToken) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Token has expired"})
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			}
			c.Abort()
			return
		}

		if claims.ContextType == identity.ContextTenant {
			if claims.MembershipID == nil || claims.CompanyID == nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
				c.Abort()
				return
			}

			role, roleErr := auth.ResolveRole(*claims.MembershipID, claims.UserID, *claims.CompanyID)
			if roleErr != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Tenant role is no longer active"})
				c.Abort()
				return
			}

			roleID := role.ID
			claims.TenantRoleID = &roleID
			claims.Role = role.Key
		}

		authctx.SetSession(c, claims)
		c.Next()
	}
}

// bearerToken extrae el token de un header Authorization con formato Bearer.
func bearerToken(header string) (string, bool) {
	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", false
	}
	return parts[1], true
}
