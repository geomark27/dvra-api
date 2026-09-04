// Package iam es el punto de ensamblaje de identidad y acceso: autenticación,
// usuarios, empresas, membresías y el RBAC relacional por tenant.
//
// Es el módulo del que cuelga la seguridad del sistema: emite los tokens,
// resuelve el rol vigente en cada petición y define la frontera entre el
// contexto de empresa y el de plataforma (ADR-002).
package iam

import (
	"dvra-api/internal/modules/iam/repository"
	"dvra-api/internal/modules/iam/service"
	"dvra-api/internal/modules/iam/transport"
	"dvra-api/internal/platform/httpx"
	"dvra-api/internal/shared/identity"
	"dvra-api/internal/shared/middleware"

	"gorm.io/gorm"
)

// Deps son las dependencias cross-módulo que iam necesita.
type Deps struct {
	// Plans resuelve el plan inicial al dar de alta una empresa (puerto hacia
	// billing).
	Plans service.PlanReader
}

// Module agrupa las dependencias ya cableadas del módulo.
type Module struct {
	auth        *service.AuthService
	jwt         service.JWTService
	rbac        *service.TenantRBACService
	users       service.UserService
	companies   service.CompanyService
	memberships service.MembershipService

	// Companies se expone para que otros módulos (vía un puerto que ellos
	// declaren) resuelvan una empresa sin importar este paquete. Lo consume
	// billing para asignar planes.
	Companies *repository.CompanyRepository
}

// New construye el módulo con la conexión inyectada.
func New(db *gorm.DB, accessSecret, refreshSecret string, deps Deps) *Module {
	userRepo := repository.NewUserRepository(db)
	companyRepo := repository.NewCompanyRepository(db)
	membershipRepo := repository.NewMembershipRepository(db)

	jwtService := service.NewJWTService(accessSecret, refreshSecret)
	rbac := service.NewTenantRBACService(db)

	return &Module{
		auth:        service.NewAuthService(userRepo, deps.Plans, jwtService, rbac, db),
		jwt:         jwtService,
		rbac:        rbac,
		users:       service.NewUserService(userRepo, db, rbac),
		companies:   service.NewCompanyService(companyRepo, deps.Plans, db),
		memberships: service.NewMembershipService(membershipRepo, rbac),
		Companies:   companyRepo,
	}
}

// Auth expone el autenticador que consumen los middlewares: el emisor de tokens
// resuelve la sesión y el RBAC el rol vigente.
//
// Sustituye al adaptador temporal que vivía en el composition root desde la
// Fase A: ahora ambas piezas satisfacen los puertos de forma nativa.
func (m *Module) Auth() middleware.Authenticator {
	return authenticator{tokens: m.jwt, roles: m.rbac}
}

// RBAC expone el verificador de permisos de tenant.
func (m *Module) RBAC() middleware.PermissionChecker { return m.rbac }

// authenticator compone las dos mitades del puerto middleware.Authenticator.
// Ambas son nativas del módulo; no hay traducción de tipos.
type authenticator struct {
	tokens service.JWTService
	roles  *service.TenantRBACService
}

func (a authenticator) Session(token string) (*identity.SessionClaims, error) {
	return a.tokens.Session(token)
}

func (a authenticator) ResolveRole(membershipID, userID, companyID uint) (identity.TenantRole, error) {
	return a.roles.ResolveRole(membershipID, userID, companyID)
}

// Name identifica al módulo en logs y en la auditoría de rutas.
func (m *Module) Name() string { return "iam" }

// RegisterRoutes monta sesión, tenant y landlord.
func (m *Module) RegisterRoutes(r httpx.Router) {
	transport.RegisterRoutes(r.Public, r.Protected, r.Platform,
		m.auth, m.users, m.companies, m.memberships, m.rbac, r.RBAC)
}
