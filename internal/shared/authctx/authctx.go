// Package authctx provee accesores tipados a los claims que AuthMiddleware
// inyecta en el contexto de Gin. Es el único lugar del proyecto donde se
// escriben las keys del contexto y el literal "superadmin": los handlers
// deben usar estos helpers en lugar de c.Get(...) directo.
package authctx

import (
	"dvra-api/internal/shared/permissions"

	"github.com/gin-gonic/gin"
)

// Keys del contexto (las setea AuthMiddleware)
const (
	KeyUserID    = "user_id"
	KeyEmail     = "email"
	KeyRole      = "role"
	KeyCompanyID = "company_id"
)

// Role devuelve el rol del token, o cadena vacía si no hay sesión.
func Role(c *gin.Context) string {
	role, ok := c.Get(KeyRole)
	if !ok {
		return ""
	}
	r, _ := role.(string)
	return r
}

// IsSuperAdmin reporta si la sesión es de SuperAdmin (acceso global).
func IsSuperAdmin(c *gin.Context) bool {
	return Role(c) == permissions.RoleSuperAdmin
}

// UserID devuelve el id del usuario autenticado.
func UserID(c *gin.Context) (uint, bool) {
	v, ok := c.Get(KeyUserID)
	if !ok {
		return 0, false
	}
	id, ok := v.(uint)
	return id, ok
}

// CompanyID devuelve la empresa activa del token.
// SuperAdmin no tiene empresa: devuelve (0, false).
func CompanyID(c *gin.Context) (uint, bool) {
	v, ok := c.Get(KeyCompanyID)
	if !ok {
		return 0, false
	}
	id, ok := v.(uint)
	return id, ok
}
