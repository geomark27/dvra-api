package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"dvra-api/internal/shared/identity"
	"dvra-api/internal/shared/permissions"

	"github.com/gin-gonic/gin"
)

// fakeAuth implementa Authenticator sin tocar el emisor real de tokens: el
// paquete de middlewares no debe depender de los servicios de la aplicación.
// El round-trip real de firma y expiración lo cubre jwt_service_test.go.
type fakeAuth struct {
	sessions map[string]*identity.SessionClaims

	// rol vigente que devuelve ResolveRole; si roleErr != nil, falla (rol
	// revocado o membresía suspendida).
	role    identity.TenantRole
	roleErr error

	// permiso que Can concede al rol vigente.
	granted string
}

func (f fakeAuth) Session(token string) (*identity.SessionClaims, error) {
	claims, ok := f.sessions[token]
	if !ok {
		return nil, identity.ErrInvalidToken
	}
	if token == "expired" {
		return nil, identity.ErrExpiredToken
	}
	// Copia: el middleware sobrescribe Role/TenantRoleID con el rol vigente.
	copied := *claims
	return &copied, nil
}

func (f fakeAuth) ResolveRole(_, _, _ uint) (identity.TenantRole, error) {
	if f.roleErr != nil {
		return identity.TenantRole{}, f.roleErr
	}
	return f.role, nil
}

func (f fakeAuth) Can(roleID uint, permission string) (bool, error) {
	return roleID == f.role.ID && permission == f.granted, nil
}

func platformClaims() *identity.SessionClaims {
	membershipID := uint(1)
	return &identity.SessionClaims{
		UserID:       7,
		Email:        "root@example.com",
		Role:         permissions.RoleSuperAdmin,
		ContextType:  identity.ContextPlatform,
		MembershipID: &membershipID,
	}
}

func tenantClaims() *identity.SessionClaims {
	companyID, membershipID, tokenRoleID := uint(9), uint(2), uint(1)
	return &identity.SessionClaims{
		UserID:       8,
		Email:        "admin@example.com",
		ContextType:  identity.ContextTenant,
		CompanyID:    &companyID,
		MembershipID: &membershipID,
		TenantRoleID: &tokenRoleID,
	}
}

func TestPlatformAndTenantPermissionsCannotCrossContexts(t *testing.T) {
	gin.SetMode(gin.TestMode)

	auth := fakeAuth{
		sessions: map[string]*identity.SessionClaims{
			"platform": platformClaims(),
			"tenant":   tenantClaims(),
		},
		role:    identity.TenantRole{ID: 1, Key: permissions.RoleAdmin},
		granted: permissions.JobsView,
	}

	router := gin.New()
	router.GET("/tenant", AuthMiddleware(auth), RequirePermission(auth, permissions.JobsView), okHandler)
	router.GET("/platform", AuthMiddleware(auth), RequirePlatformPermission(permissions.PlatformCompaniesView), okHandler)

	assertStatus(t, router, "/tenant", "platform", http.StatusForbidden)
	assertStatus(t, router, "/platform", "tenant", http.StatusForbidden)
	assertStatus(t, router, "/tenant", "tenant", http.StatusOK)
	assertStatus(t, router, "/platform", "platform", http.StatusOK)
}

// El rol efectivo se resuelve contra la base en cada petición, no se lee del
// token: el token dice rol 1 y el vigente es el 4.
func TestTenantMiddlewareUsesCurrentMembershipRole(t *testing.T) {
	gin.SetMode(gin.TestMode)

	currentRole := identity.TenantRole{ID: 4, Key: permissions.RoleInterviewer}
	auth := fakeAuth{
		sessions: map[string]*identity.SessionClaims{"tenant": tenantClaims()},
		role:     currentRole,
	}

	router := gin.New()
	router.GET("/current-role", AuthMiddleware(auth), func(c *gin.Context) {
		roleID, exists := c.Get("tenant_role_id")
		if !exists || roleID != currentRole.ID {
			c.Status(http.StatusInternalServerError)
			return
		}
		if c.GetString("role") != permissions.RoleInterviewer {
			c.Status(http.StatusInternalServerError)
			return
		}
		c.Status(http.StatusOK)
	})

	assertStatus(t, router, "/current-role", "tenant", http.StatusOK)
}

// Un rol revocado o una membresía suspendida cortan la petición aunque el token
// siga siendo válido.
func TestTenantSessionIsRejectedWhenRoleIsRevoked(t *testing.T) {
	gin.SetMode(gin.TestMode)

	auth := fakeAuth{
		sessions: map[string]*identity.SessionClaims{"tenant": tenantClaims()},
		roleErr:  errors.New("membership is not active"),
	}

	router := gin.New()
	router.GET("/tenant", AuthMiddleware(auth), okHandler)

	assertStatus(t, router, "/tenant", "tenant", http.StatusUnauthorized)
}

func TestAuthMiddlewareRejectsBadTokens(t *testing.T) {
	gin.SetMode(gin.TestMode)

	auth := fakeAuth{
		sessions: map[string]*identity.SessionClaims{
			"tenant":  tenantClaims(),
			"expired": tenantClaims(),
		},
		role: identity.TenantRole{ID: 1, Key: permissions.RoleAdmin},
	}

	router := gin.New()
	router.GET("/tenant", AuthMiddleware(auth), okHandler)

	assertStatus(t, router, "/tenant", "expired", http.StatusUnauthorized)
	assertStatus(t, router, "/tenant", "desconocido", http.StatusUnauthorized)

	// Sin header y con header mal formado.
	for _, header := range []string{"", "Basic abc", "Bearer"} {
		request := httptest.NewRequest(http.MethodGet, "/tenant", nil)
		if header != "" {
			request.Header.Set("Authorization", header)
		}
		response := httptest.NewRecorder()
		router.ServeHTTP(response, request)
		if response.Code != http.StatusUnauthorized {
			t.Fatalf("header %q devolvió %d, se esperaba 401", header, response.Code)
		}
	}
}

func okHandler(c *gin.Context) { c.Status(http.StatusOK) }

func assertStatus(t *testing.T, handler http.Handler, path, token string, expected int) {
	t.Helper()
	request := httptest.NewRequest(http.MethodGet, path, nil)
	request.Header.Set("Authorization", "Bearer "+token)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != expected {
		t.Fatalf("%s returned %d, expected %d", path, response.Code, expected)
	}
}
