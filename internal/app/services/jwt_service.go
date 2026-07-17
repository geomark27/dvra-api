package services

import (
	"errors"
	"strconv"
	"time"

	"dvra-api/internal/shared/identity"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")
)

// JWTClaims represents the claims stored in JWT
type JWTClaims struct {
	UserID       uint   `json:"user_id"`
	CompanyID    *uint  `json:"company_id"` // Present only in tenant context.
	Email        string `json:"email"`
	Role         string `json:"role"`
	ContextType  string `json:"context_type"`
	MembershipID *uint  `json:"membership_id,omitempty"`
	TenantRoleID *uint  `json:"tenant_role_id,omitempty"`
	TokenType    string `json:"token_type,omitempty"`
	jwt.RegisteredClaims
}

type RefreshClaims struct {
	UserID       uint   `json:"user_id"`
	TokenType    string `json:"token_type"`
	ContextType  string `json:"context_type"`
	CompanyID    *uint  `json:"company_id,omitempty"`
	MembershipID *uint  `json:"membership_id,omitempty"`
	jwt.RegisteredClaims
}

// JWTService handles JWT token operations
type JWTService interface {
	GenerateAccessToken(userID uint, email string, context identity.Context) (string, error)
	GenerateRefreshToken(userID uint, context identity.Context) (string, error)
	ValidateAccessToken(tokenString string) (*JWTClaims, error)
	ValidateRefreshToken(tokenString string) (*RefreshClaims, error)
}

type jwtService struct {
	accessSecret  string
	refreshSecret string
	accessTTL     time.Duration
	refreshTTL    time.Duration
}

// NewJWTService creates a new JWT service
func NewJWTService(accessSecret, refreshSecret string) JWTService {
	return &jwtService{
		accessSecret:  accessSecret,
		refreshSecret: refreshSecret,
		accessTTL:     1 * time.Hour,       // Access token: 1 hour
		refreshTTL:    30 * 24 * time.Hour, // Refresh token: 30 days
	}
}

// GenerateAccessToken generates a new access token
func (s *jwtService) GenerateAccessToken(userID uint, email string, context identity.Context) (string, error) {
	role := context.Role
	if context.Type == identity.ContextTenant {
		role = "" // Tenant authorization is referenced by tenant_role_id.
	}
	claims := JWTClaims{
		UserID:       userID,
		CompanyID:    context.CompanyID,
		Email:        email,
		Role:         role,
		ContextType:  context.Type,
		MembershipID: context.MembershipID,
		TenantRoleID: context.TenantRoleID,
		TokenType:    "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.accessTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.accessSecret))
}

// GenerateRefreshToken generates a new refresh token
func (s *jwtService) GenerateRefreshToken(userID uint, context identity.Context) (string, error) {
	claims := RefreshClaims{
		UserID:       userID,
		TokenType:    "refresh",
		ContextType:  context.Type,
		CompanyID:    context.CompanyID,
		MembershipID: context.MembershipID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.FormatUint(uint64(userID), 10),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.refreshTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.refreshSecret))
}

// ValidateAccessToken validates a context-aware access token. Legacy tokens
// without token_type/context_type are intentionally rejected.
func (s *jwtService) ValidateAccessToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, ErrInvalidToken
		}
		return []byte(s.accessSecret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid && claims.TokenType == "access" && claims.UserID > 0 {
		if !validTokenContext(claims.ContextType, claims.CompanyID, claims.MembershipID, claims.TenantRoleID, claims.Role) {
			return nil, ErrInvalidToken
		}
		return claims, nil
	}

	return nil, ErrInvalidToken
}

// ValidateRefreshToken only accepts refresh tokens signed with the refresh
// secret. Access tokens can never be exchanged at the refresh endpoint.
func (s *jwtService) ValidateRefreshToken(tokenString string) (*RefreshClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &RefreshClaims{}, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, ErrInvalidToken
		}
		return []byte(s.refreshSecret), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}
	claims, ok := token.Claims.(*RefreshClaims)
	if !ok || !token.Valid || claims.TokenType != "refresh" || claims.UserID == 0 ||
		!validRefreshContext(claims.ContextType, claims.CompanyID, claims.MembershipID) {
		return nil, ErrInvalidToken
	}
	return claims, nil
}

func validTokenContext(contextType string, companyID, membershipID, tenantRoleID *uint, role string) bool {
	switch contextType {
	case identity.ContextPlatform:
		return companyID == nil && membershipID != nil && *membershipID > 0 && tenantRoleID == nil && role != ""
	case identity.ContextTenant:
		return companyID != nil && *companyID > 0 && membershipID != nil && *membershipID > 0 && tenantRoleID != nil && *tenantRoleID > 0 && role == ""
	case identity.ContextIdentity:
		return companyID == nil && membershipID == nil && tenantRoleID == nil && role == ""
	default:
		return false
	}
}

func validRefreshContext(contextType string, companyID, membershipID *uint) bool {
	switch contextType {
	case identity.ContextPlatform:
		return companyID == nil && membershipID != nil && *membershipID > 0
	case identity.ContextTenant:
		return companyID != nil && *companyID > 0 && membershipID != nil && *membershipID > 0
	case identity.ContextIdentity:
		return companyID == nil && membershipID == nil
	default:
		return false
	}
}
