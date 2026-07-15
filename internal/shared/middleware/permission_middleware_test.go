package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"dvra-api/internal/app/services"
	"dvra-api/internal/shared/identity"
	"dvra-api/internal/shared/permissions"

	"github.com/gin-gonic/gin"
)

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
	tenantToken, err := jwt.GenerateAccessToken(8, "admin@example.com", identity.Context{
		Type: identity.ContextTenant, CompanyID: &companyID, MembershipID: &tenantMembershipID, Role: permissions.RoleAdmin,
	})
	if err != nil {
		t.Fatal(err)
	}

	router := gin.New()
	router.GET("/tenant", AuthMiddleware(jwt), RequirePermission(permissions.JobsView), okHandler)
	router.GET("/platform", AuthMiddleware(jwt), RequirePlatformPermission(permissions.PlatformCompaniesView), okHandler)

	assertStatus(t, router, "/tenant", platformToken, http.StatusForbidden)
	assertStatus(t, router, "/platform", tenantToken, http.StatusForbidden)
	assertStatus(t, router, "/tenant", tenantToken, http.StatusOK)
	assertStatus(t, router, "/platform", platformToken, http.StatusOK)
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
