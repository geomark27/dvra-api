// Package httpx define el contrato por el que un módulo monta sus rutas.
//
// Vive en platform/ y no en platform/server para que los módulos puedan
// importarlo sin importar el composition root: si viviera en server, habría
// ciclo (server -> modules -> server).
package httpx

import (
	"dvra-api/internal/shared/middleware"

	"github.com/gin-gonic/gin"
)

// Router agrupa los grupos de rutas y las políticas transversales que un módulo
// necesita para montarse. Sustituye a la lista de parámetros posicionales que
// crecía con cada dominio.
//
// Son tres grupos porque un mismo módulo puede servir a los tres contextos: iam
// expone /auth público, /users protegido y /platform/users de landlord con los
// mismos handlers.
type Router struct {
	// Public es /api/v1 sin autenticación.
	Public *gin.RouterGroup
	// Protected es /api/v1 con AuthMiddleware.
	Protected *gin.RouterGroup
	// Platform es /api/v1/platform con AuthMiddleware.
	Platform *gin.RouterGroup

	// RBAC resuelve permisos de rol para RequirePermission.
	RBAC middleware.PermissionChecker
	// Plans resuelve entitlements de plan para RequireFeature.
	Plans middleware.FeatureChecker
}

// Module es todo lo que el composition root necesita saber de un dominio.
//
// Coste de añadir un módulo nuevo: una línea en el registro.
type Module interface {
	// Name identifica al módulo en logs y en la auditoría de rutas.
	Name() string
	// RegisterRoutes monta las rutas del módulo en los grupos que le apliquen.
	RegisterRoutes(r Router)
}
