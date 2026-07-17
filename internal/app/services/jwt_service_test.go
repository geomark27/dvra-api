package services

import (
	"errors"
	"testing"

	"dvra-api/internal/shared/identity"
)

func TestJWTSeparatesAccessAndRefreshTokens(t *testing.T) {
	svc := NewJWTService("access-secret-with-more-than-32-chars", "refresh-secret-with-more-than-32-chars")
	companyID := uint(9)
	membershipID := uint(12)
	tenantRoleID := uint(1)
	context := identity.Context{Type: identity.ContextTenant, CompanyID: &companyID, MembershipID: &membershipID, TenantRoleID: &tenantRoleID, Role: "admin"}
	access, err := svc.GenerateAccessToken(7, "user@example.com", context)
	if err != nil {
		t.Fatal(err)
	}
	refresh, err := svc.GenerateRefreshToken(7, context)
	if err != nil {
		t.Fatal(err)
	}

	accessClaims, err := svc.ValidateAccessToken(access)
	if err != nil || accessClaims.UserID != 7 || accessClaims.TokenType != "access" || accessClaims.ContextType != identity.ContextTenant {
		t.Fatalf("unexpected access claims: %#v, %v", accessClaims, err)
	}
	if accessClaims.TenantRoleID == nil || *accessClaims.TenantRoleID != tenantRoleID || accessClaims.Role != "" {
		t.Fatalf("tenant access token did not reference its role ID: %#v", accessClaims)
	}
	refreshClaims, err := svc.ValidateRefreshToken(refresh)
	if err != nil || refreshClaims.UserID != 7 || refreshClaims.Subject != "7" || refreshClaims.CompanyID == nil || *refreshClaims.CompanyID != companyID {
		t.Fatalf("unexpected refresh claims: %#v, %v", refreshClaims, err)
	}
	if _, err := svc.ValidateRefreshToken(access); !errors.Is(err, ErrInvalidToken) {
		t.Fatalf("access token accepted as refresh: %v", err)
	}
	if _, err := svc.ValidateAccessToken(refresh); !errors.Is(err, ErrInvalidToken) {
		t.Fatalf("refresh token accepted as access: %v", err)
	}
}

func TestJWTRejectsTokenFromAnotherSecret(t *testing.T) {
	issuer := NewJWTService("issuer-access-secret-with-32-chars", "issuer-refresh-secret-with-32-chars")
	validator := NewJWTService("other-access-secret-with-32-chars", "other-refresh-secret-with-32-chars")
	token, err := issuer.GenerateAccessToken(1, "a@example.com", identity.Context{Type: identity.ContextIdentity})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := validator.ValidateAccessToken(token); !errors.Is(err, ErrInvalidToken) {
		t.Fatalf("token signed with another secret was accepted: %v", err)
	}
}

func TestJWTRejectsMixedAuthorizationContext(t *testing.T) {
	svc := NewJWTService("access-secret-with-more-than-32-chars", "refresh-secret-with-more-than-32-chars")
	companyID, membershipID := uint(9), uint(3)
	token, err := svc.GenerateAccessToken(1, "root@example.com", identity.Context{
		Type: identity.ContextPlatform, CompanyID: &companyID, MembershipID: &membershipID, Role: "superadmin",
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := svc.ValidateAccessToken(token); !errors.Is(err, ErrInvalidToken) {
		t.Fatalf("mixed platform/tenant context was accepted: %v", err)
	}
}
