package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"dvra-api/internal/app/models"
	"dvra-api/internal/app/services"
	"dvra-api/internal/shared/identity"
	"dvra-api/internal/shared/permissions"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type middlewareTenantAuthorizer struct {
	roleID uint
	role   *models.TenantRole
}

func (f middlewareTenantAuthorizer) Can(roleID uint, permission string) (bool, error) {
	return roleID == f.roleID && permission == permissions.JobsView, nil
}
func (f middlewareTenantAuthorizer) Permissions(uint) ([]string, error) { return nil, nil }
func (f middlewareTenantAuthorizer) RoleByID(id uint) (*models.TenantRole, error) {
	if id != f.roleID {
		return nil, gorm.ErrRecordNotFound
	}
	return f.role, nil
}
func (f middlewareTenantAuthorizer) RoleByKey(string) (*models.TenantRole, error) {
	return nil, gorm.ErrRecordNotFound
}
func (f middlewareTenantAuthorizer) RoleForMembership(_, _, _ uint) (*models.TenantRole, error) {
	return f.role, nil
}

func TestPlatformAndTenantPermissionsCannotCrossContexts(t *testing.T) {
	gin.SetMode(gin.TestMode)
	jwt := services.NewJWTService("access-secret-with-more-than-32-chars", "refresh-secret-with-more-than-32-chars")
	platformMembershipID := uint(1)
	platformToken, err := jwt.GenerateAccessToken(7, "root@example.com", identity.Context{
		Type: identity.ContextPlatform, MembershipID: &platformMembershipID, Role: permissions.RoleSuperAdmin,
	})
	if err != nil {
		t.Fatal(err)
	}
	companyID, tenantMembershipID := uint(9), uint(2)
	tenantRoleID := uint(1)
	tenantToken, err := jwt.GenerateAccessToken(8, "admin@example.com", identity.Context{
		Type: identity.ContextTenant, CompanyID: &companyID, MembershipID: &tenantMembershipID, TenantRoleID: &tenantRoleID,
	})
	if err != nil {
		t.Fatal(err)
	}
	rbac := middlewareTenantAuthorizer{
		roleID: tenantRoleID,
		role:   &models.TenantRole{BaseModel: models.BaseModel{ID: tenantRoleID}, Key: permissions.RoleAdmin, IsActive: true},
	}

	router := gin.New()
	router.GET("/tenant", AuthMiddleware(jwt, rbac), RequirePermission(rbac, permissions.JobsView), okHandler)
	router.GET("/platform", AuthMiddleware(jwt, rbac), RequirePlatformPermission(permissions.PlatformCompaniesView), okHandler)

	assertStatus(t, router, "/tenant", platformToken, http.StatusForbidden)
	assertStatus(t, router, "/platform", tenantToken, http.StatusForbidden)
	assertStatus(t, router, "/tenant", tenantToken, http.StatusOK)
	assertStatus(t, router, "/platform", platformToken, http.StatusOK)
}

func TestTenantMiddlewareUsesCurrentMembershipRole(t *testing.T) {
	gin.SetMode(gin.TestMode)
	jwt := services.NewJWTService("access-secret-with-more-than-32-chars", "refresh-secret-with-more-than-32-chars")
	companyID, membershipID, tokenRoleID := uint(9), uint(2), uint(1)
	token, err := jwt.GenerateAccessToken(8, "user@example.com", identity.Context{
		Type: identity.ContextTenant, CompanyID: &companyID, MembershipID: &membershipID, TenantRoleID: &tokenRoleID,
	})
	if err != nil {
		t.Fatal(err)
	}
	currentRoleID := uint(4)
	rbac := middlewareTenantAuthorizer{
		roleID: currentRoleID,
		role:   &models.TenantRole{BaseModel: models.BaseModel{ID: currentRoleID}, Key: permissions.RoleInterviewer, IsActive: true},
	}
	router := gin.New()
	router.GET("/current-role", AuthMiddleware(jwt, rbac), func(c *gin.Context) {
		roleID, exists := c.Get("tenant_role_id")
		if !exists || roleID != currentRoleID {
			c.Status(http.StatusInternalServerError)
			return
		}
		c.Status(http.StatusOK)
	})
	assertStatus(t, router, "/current-role", token, http.StatusOK)
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
