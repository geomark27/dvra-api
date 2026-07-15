// Package authctx provee accesores tipados a los claims que AuthMiddleware
// inyecta en el contexto de Gin. Es el único lugar del proyecto donde se
// escriben las keys del contexto y el literal "superadmin": los handlers
// deben usar estos helpers en lugar de c.Get(...) directo.
package authctx

import (
	"dvra-api/internal/shared/identity"
	"dvra-api/internal/shared/permissions"
	"dvra-api/internal/shared/tenant"

	"github.com/gin-gonic/gin"
)

// Keys del contexto (las setea AuthMiddleware)
const (
	KeyUserID       = "user_id"
	KeyEmail        = "email"
	KeyRole         = "role"
	KeyCompanyID    = "company_id"
	KeyMembershipID = "membership_id"
	KeyContextType  = "context_type"
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

func ContextType(c *gin.Context) string {
	value, _ := c.Get(KeyContextType)
	contextType, _ := value.(string)
	return contextType
}

func IsTenantContext(c *gin.Context) bool { return ContextType(c) == identity.ContextTenant }

func IsPlatformContext(c *gin.Context) bool { return ContextType(c) == identity.ContextPlatform }

// TenantScope construye el scope tipado de la sesión actual.
func TenantScope(c *gin.Context) (tenant.Scope, bool) {
	if !IsTenantContext(c) {
		return tenant.Scope{}, false
	}
	companyID, ok := CompanyID(c)
	if !ok {
		return tenant.Scope{}, false
	}
	userID, ok := UserID(c)
	if !ok {
		return tenant.Scope{}, false
	}
	return tenant.ForActor(companyID, userID, Role(c)), true
}

// IsSuperAdmin reporta si la sesión pertenece al landlord SuperAdmin.
func IsSuperAdmin(c *gin.Context) bool {
	return IsPlatformContext(c) && Role(c) == permissions.RoleSuperAdmin
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
// Los contextos platform e identity no tienen empresa: devuelve (0, false).
func CompanyID(c *gin.Context) (uint, bool) {
	v, ok := c.Get(KeyCompanyID)
	if !ok {
		return 0, false
	}
	id, ok := v.(uint)
	return id, ok
}

func MembershipID(c *gin.Context) (uint, bool) {
	v, ok := c.Get(KeyMembershipID)
	if !ok {
		return 0, false
	}
	id, ok := v.(uint)
	return id, ok
}
