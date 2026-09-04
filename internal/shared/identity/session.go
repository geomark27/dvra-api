package identity

import "errors"

// Errores de validación de token. Viven aquí y no en el emisor porque los
// consumen los middlewares, que no deben conocer al emisor concreto.
var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")
)

// SessionClaims es la vista neutral del token que consumen los middlewares y
// los accesores de contexto.
//
// No lleva jwt.RegisteredClaims: eso es un detalle del emisor. Mantenerlo fuera
// es lo que permite que shared/middleware no dependa del paquete de servicios.
type SessionClaims struct {
	UserID       uint
	Email        string
	Role         string
	ContextType  string
	CompanyID    *uint
	MembershipID *uint
	TenantRoleID *uint
}

// TenantRole es la vista mínima del rol relacional resuelto por RBAC: lo justo
// para autorizar y para poblar el contexto.
type TenantRole struct {
	ID  uint
	Key string
}
