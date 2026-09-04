package authctx

import (
	"net/http"

	"dvra-api/internal/shared/tenant"

	"github.com/gin-gonic/gin"
)

// Helpers de scoping multi-tenant (ESTANDARES_DESARROLLO.md §4).
//
// Sustituyen al patrón repetido a mano en los handlers:
//
//	if authctx.IsSuperAdmin(c) { ... } else { companyID, ok := c.Get("company_id") ... }
//
// Cada copia de ese bloque es una oportunidad de fuga cross-tenant. Aquí vive
// una sola vez, y los helpers ya emiten el 403 y abortan.

// RequireScope resuelve la frontera de datos de la petición:
//   - contexto platform con SuperAdmin -> scope global
//   - contexto tenant                  -> scope de la empresa del token
//
// Si no hay un scope válido responde 403 y aborta; el handler solo debe
// hacer `return`.
//
//	scope, ok := authctx.RequireScope(c)
//	if !ok {
//		return
//	}
func RequireScope(c *gin.Context) (tenant.Scope, bool) {
	if IsSuperAdmin(c) {
		return tenant.GlobalScope(), true
	}

	scope, ok := TenantScope(c)
	if !ok || !scope.Valid() {
		abortForbidden(c, "No company context")
		return tenant.Scope{}, false
	}
	return scope, true
}

// RequireCompanyScope exige contexto de empresa: rechaza platform y SuperAdmin.
// Para operaciones que SIEMPRE pertenecen a una empresa (crear una vacante,
// subir un CV), donde un scope global no significa nada.
func RequireCompanyScope(c *gin.Context) (tenant.Scope, bool) {
	scope, ok := TenantScope(c)
	if !ok || !scope.Valid() {
		abortForbidden(c, "No company context")
		return tenant.Scope{}, false
	}
	return scope, true
}

// RequireCompanyID devuelve la empresa del token para forzarla en los DTOs de
// creación (§4.2: el company_id del body nunca se respeta).
func RequireCompanyID(c *gin.Context) (uint, bool) {
	companyID, ok := CompanyID(c)
	if !ok || companyID == 0 {
		abortForbidden(c, "No company context")
		return 0, false
	}
	return companyID, true
}

func abortForbidden(c *gin.Context, message string) {
	c.JSON(http.StatusForbidden, gin.H{"error": message})
	c.Abort()
}
