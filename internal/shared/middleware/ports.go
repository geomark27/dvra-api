package middleware

import "dvra-api/internal/shared/identity"

// Puertos del paquete: las interfaces mínimas que los middlewares necesitan,
// declaradas AQUÍ (en el consumidor) y no en quien las implementa.
//
// Es lo que permite que un módulo monte sus rutas sin arrastrar el paquete de
// servicios completo: antes `staffing/transport` -> `shared/middleware` ->
// `app/services` metía TODOS los services en el grafo de cualquier módulo.
//
// En Go el acoplamiento es estructural: quien las satisface no necesita
// importar este paquete ni declarar que las implementa.

// TokenValidator es lo mínimo que AuthMiddleware necesita del emisor de tokens.
// Debe devolver identity.ErrExpiredToken cuando el token esté vencido.
type TokenValidator interface {
	Session(token string) (*identity.SessionClaims, error)
}

// RoleResolver resuelve el rol vigente de una membresía. Se consulta en cada
// petición a propósito: así revocar un rol surte efecto de inmediato, sin
// esperar a que caduque el token.
type RoleResolver interface {
	ResolveRole(membershipID, userID, companyID uint) (identity.TenantRole, error)
}

// Authenticator es lo que recibe AuthMiddleware: emisor de tokens + resolutor
// de rol.
type Authenticator interface {
	TokenValidator
	RoleResolver
}

// PermissionChecker responde si un rol de tenant tiene un permiso.
// La satisface tal cual el servicio de RBAC relacional.
type PermissionChecker interface {
	Can(tenantRoleID uint, permissionKey string) (bool, error)
}

// FeatureChecker responde si el plan de una empresa habilita un feature.
// La satisface tal cual el servicio de planes.
type FeatureChecker interface {
	CompanyHasFeature(companyID uint, feature string) (bool, error)
}
