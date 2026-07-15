package services

import (
	"errors"

	"dvra-api/internal/app/dtos"
	"dvra-api/internal/app/models"
	"dvra-api/internal/app/repositories"
	"dvra-api/internal/shared/apperr"
	"dvra-api/internal/shared/identity"
	"dvra-api/internal/shared/permissions"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Errores de dominio de auth. Son valores *apperr.AppError, de modo que llevan
// su código HTTP y los handlers pueden usar apperr.StatusCode(err); las
// comparaciones errors.Is/== siguen funcionando por identidad. Los mensajes de
// credenciales son deliberadamente genéricos (no revelan si el email existe).
var (
	ErrInvalidCredentials = apperr.Unauthorized("invalid email or password")
	ErrEmailExists        = apperr.Conflict("email already exists")
	ErrUserNotFound       = apperr.NotFound("user not found")
	ErrInvalidPassword    = apperr.Unauthorized("invalid password")
	ErrCompanyNotFound    = apperr.NotFound("company not found")
	ErrNoMembership       = apperr.Forbidden("user does not belong to this company")
)

// AuthService handles authentication business logic
type AuthService struct {
	userRepo   repositories.UserRepository
	planRepo   repositories.PlanRepository
	jwtService JWTService
	db         *gorm.DB
}

// NewAuthService creates a new auth service
func NewAuthService(userRepo repositories.UserRepository, planRepo repositories.PlanRepository, jwtService JWTService, db *gorm.DB) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
		planRepo:   planRepo,
		jwtService: jwtService,
		db:         db,
	}
}

// Register creates a new user account
func (s *AuthService) Register(dto *dtos.RegisterDTO) (*dtos.LoginResponseDTO, error) {
	// Check if email already exists
	existingUser, _ := s.userRepo.FindByEmail(dto.Email)
	if existingUser != nil {
		return nil, ErrEmailExists
	}

	// Hash password
	hashedPassword, err := HashPassword(dto.Password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &models.User{
		Email:        dto.Email,
		PasswordHash: hashedPassword,
		FirstName:    dto.FirstName,
		LastName:     dto.LastName,
		IsActive:     true,
	}

	createdUser, err := s.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	// Generate tokens (no company yet, user just registered)
	context := identity.Context{Type: identity.ContextIdentity}
	accessToken, refreshToken, err := s.issueTokens(createdUser, context)
	if err != nil {
		return nil, err
	}

	return &dtos.LoginResponseDTO{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: dtos.UserResponse{
			ID:        createdUser.ID,
			Email:     createdUser.Email,
			FirstName: createdUser.FirstName,
			LastName:  createdUser.LastName,
			IsActive:  createdUser.IsActive,
		},
		ActiveContext: dtos.ToAuthContextResponse(context),
	}, nil
}

// Login authenticates a user and returns tokens
func (s *AuthService) Login(dto *dtos.LoginDTO) (*dtos.LoginResponseDTO, error) {
	// Find user by email
	user, err := s.userRepo.FindByEmail(dto.Email)
	if err != nil || user == nil {
		return nil, ErrInvalidCredentials
	}

	// Check if user is active
	if !user.IsActive {
		return nil, apperr.Forbidden("account is inactive")
	}

	// Verify password
	if err := ComparePassword(user.PasswordHash, dto.Password); err != nil {
		return nil, ErrInvalidCredentials
	}

	context, err := s.resolveInitialContext(user.ID)
	if err != nil {
		return nil, err
	}
	accessToken, refreshToken, err := s.issueTokens(user, context)
	if err != nil {
		return nil, err
	}

	// Update last login (best-effort: no debe bloquear el login)
	_ = s.userRepo.UpdateLastLogin(user.ID)

	return &dtos.LoginResponseDTO{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: dtos.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			IsActive:  user.IsActive,
		},
		ActiveContext: dtos.ToAuthContextResponse(context),
	}, nil
}

// RefreshToken generates new access token from refresh token
func (s *AuthService) RefreshToken(dto *dtos.RefreshTokenDTO) (*dtos.RefreshTokenResponseDTO, error) {
	// Validate refresh token
	claims, err := s.jwtService.ValidateRefreshToken(dto.RefreshToken)
	if err != nil {
		return nil, apperr.Unauthorized("invalid refresh token")
	}

	// Get user
	user, err := s.userRepo.FindByID(int(claims.UserID))
	if err != nil || user == nil {
		return nil, ErrUserNotFound
	}

	if !user.IsActive {
		return nil, apperr.Forbidden("account is inactive")
	}

	context, err := s.resolveRefreshContext(user.ID, claims)
	if err != nil {
		return nil, err
	}
	accessToken, newRefreshToken, err := s.issueTokens(user, context)
	if err != nil {
		return nil, err
	}

	return &dtos.RefreshTokenResponseDTO{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

// ChangePassword changes user password
func (s *AuthService) ChangePassword(userID uint, dto *dtos.ChangePasswordDTO) error {
	// Get user
	user, err := s.userRepo.FindByID(int(userID))
	if err != nil || user == nil {
		return ErrUserNotFound
	}

	// Verify old password
	if err := ComparePassword(user.PasswordHash, dto.OldPassword); err != nil {
		return ErrInvalidPassword
	}

	// Hash new password
	hashedPassword, err := HashPassword(dto.NewPassword)
	if err != nil {
		return err
	}

	// Update password
	return s.userRepo.UpdatePassword(userID, hashedPassword)
}

// GetMe returns current user info
func (s *AuthService) GetMe(userID uint) (*dtos.UserResponse, error) {
	user, err := s.userRepo.FindByID(int(userID))
	if err != nil || user == nil {
		return nil, ErrUserNotFound
	}

	return &dtos.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		IsActive:  user.IsActive,
	}, nil
}

// RegisterCompany creates a new company with its first admin user
func (s *AuthService) RegisterCompany(dto *dtos.RegisterCompanyDTO) (*dtos.RegisterCompanyResponseDTO, error) {
	// Check if admin email already exists
	existingUser, _ := s.userRepo.FindByEmail(dto.AdminEmail)
	if existingUser != nil {
		return nil, ErrEmailExists
	}

	// Verify that the free plan exists and is active
	freePlan, err := s.planRepo.FindActiveBySlug("free")
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("free plan is not available, please contact support")
		}
		return nil, err
	}

	// Start transaction
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. Create company with validated free plan
	timezone := dto.Timezone
	if timezone == "" {
		timezone = "America/Bogota"
	}

	company := models.Company{
		Name:     dto.CompanyName,
		Slug:     dto.CompanySlug,
		PlanTier: freePlan.Slug, // Use validated free plan slug
		Timezone: timezone,
	}

	if err := tx.Create(&company).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 2. Create admin user
	hashedPassword, err := HashPassword(dto.AdminPassword)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	admin := models.User{
		Email:        dto.AdminEmail,
		PasswordHash: hashedPassword,
		FirstName:    dto.AdminFirstName,
		LastName:     dto.AdminLastName,
		IsActive:     true,
	}

	if err := tx.Create(&admin).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 3. Create membership (admin → company)
	membership := models.Membership{
		UserID:    admin.ID,
		CompanyID: company.ID,
		Role:      permissions.RoleAdmin,
		Status:    "active",
		IsDefault: true, // First company is always default
	}

	if err := tx.Create(&membership).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	membershipID := membership.ID
	companyID := company.ID
	context := identity.Context{Type: identity.ContextTenant, CompanyID: &companyID, MembershipID: &membershipID, Role: membership.Role}
	accessToken, refreshToken, err := s.issueTokens(&admin, context)
	if err != nil {
		return nil, err
	}

	return &dtos.RegisterCompanyResponseDTO{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Company: dtos.CompanyResponse{
			ID:       company.ID,
			Name:     company.Name,
			Slug:     company.Slug,
			PlanTier: company.PlanTier,
		},
		Admin: dtos.UserResponse{
			ID:        admin.ID,
			Email:     admin.Email,
			FirstName: admin.FirstName,
			LastName:  admin.LastName,
			IsActive:  admin.IsActive,
		},
		ActiveContext: dtos.ToAuthContextResponse(context),
	}, nil
}

// GetUserCompanies returns all companies that a user belongs to
func (s *AuthService) GetUserCompanies(userID uint) ([]dtos.CompanyResponse, error) {
	var memberships []models.Membership
	err := s.db.Preload("Company").Where("user_id = ? AND status = ?", userID, "active").Find(&memberships).Error
	if err != nil {
		return nil, err
	}

	companies := make([]dtos.CompanyResponse, 0, len(memberships))
	for _, membership := range memberships {
		if membership.Company != nil {
			companies = append(companies, dtos.CompanyResponse{
				ID:           membership.Company.ID,
				Name:         membership.Company.Name,
				Slug:         membership.Company.Slug,
				PlanTier:     membership.Company.PlanTier,
				MembershipID: membership.ID,
				Role:         membership.Role,
			})
		}
	}

	return companies, nil
}

// SwitchCompany generates a new token for a different company context
func (s *AuthService) SwitchCompany(userID uint, dto *dtos.SwitchCompanyDTO) (*dtos.SwitchCompanyResponseDTO, error) {
	// Verify user has membership in this company
	var membership models.Membership
	err := s.db.Preload("Company").
		Where("user_id = ? AND company_id = ? AND status = ?", userID, dto.CompanyID, "active").
		First(&membership).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrNoMembership
		}
		return nil, err
	}

	if membership.Company == nil {
		return nil, ErrCompanyNotFound
	}

	// Get user email
	user, err := s.userRepo.FindByID(int(userID))
	if err != nil {
		return nil, ErrUserNotFound
	}

	membershipID := membership.ID
	context := identity.Context{Type: identity.ContextTenant, CompanyID: &dto.CompanyID, MembershipID: &membershipID, Role: membership.Role}
	accessToken, refreshToken, err := s.issueTokens(user, context)
	if err != nil {
		return nil, err
	}

	return &dtos.SwitchCompanyResponseDTO{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Company: dtos.CompanyResponse{
			ID:           membership.Company.ID,
			Name:         membership.Company.Name,
			Slug:         membership.Company.Slug,
			PlanTier:     membership.Company.PlanTier,
			MembershipID: membership.ID,
			Role:         membership.Role,
		},
	}, nil
}

func (s *AuthService) SwitchToPlatform(userID uint) (*dtos.SwitchPlatformResponseDTO, error) {
	user, err := s.userRepo.FindByID(int(userID))
	if err != nil || user == nil {
		return nil, ErrUserNotFound
	}
	var membership models.PlatformMembership
	if err := s.db.Where("user_id = ? AND status = ?", userID, "active").First(&membership).Error; err != nil {
		return nil, apperr.Forbidden("user has no active platform access")
	}
	membershipID := membership.ID
	context := identity.Context{Type: identity.ContextPlatform, MembershipID: &membershipID, Role: membership.Role}
	accessToken, refreshToken, err := s.issueTokens(user, context)
	if err != nil {
		return nil, err
	}
	return &dtos.SwitchPlatformResponseDTO{
		AccessToken: accessToken, RefreshToken: refreshToken,
		ActiveContext: dtos.ToAuthContextResponse(context),
	}, nil
}

// LoginWithCompanies authenticates a user and returns tokens with companies list
func (s *AuthService) LoginWithCompanies(dto *dtos.LoginDTO) (*dtos.LoginResponseWithCompaniesDTO, error) {
	// Find user by email
	user, err := s.userRepo.FindByEmail(dto.Email)
	if err != nil || user == nil {
		return nil, ErrInvalidCredentials
	}

	// Check if user is active
	if !user.IsActive {
		return nil, apperr.Forbidden("account is inactive")
	}

	// Verify password
	if err := ComparePassword(user.PasswordHash, dto.Password); err != nil {
		return nil, ErrInvalidCredentials
	}

	context, err := s.resolveInitialContext(user.ID)
	if err != nil {
		return nil, err
	}

	// Get all user's companies
	companies, err := s.GetUserCompanies(user.ID)
	if err != nil {
		return nil, err
	}
	platformAccess, err := s.getPlatformAccess(user.ID)
	if err != nil {
		return nil, err
	}

	accessToken, refreshToken, err := s.issueTokens(user, context)
	if err != nil {
		return nil, err
	}

	// Update last login (best-effort: no debe bloquear el login)
	_ = s.userRepo.UpdateLastLogin(user.ID)

	return &dtos.LoginResponseWithCompaniesDTO{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: dtos.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			IsActive:  user.IsActive,
		},
		Companies:      companies,
		ActiveContext:  dtos.ToAuthContextResponse(context),
		PlatformAccess: platformAccess,
	}, nil
}

// HashPassword hashes a plain text password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// ComparePassword compares a hashed password with plain text password
func ComparePassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (s *AuthService) issueTokens(user *models.User, context identity.Context) (string, string, error) {
	accessToken, err := s.jwtService.GenerateAccessToken(user.ID, user.Email, context)
	if err != nil {
		return "", "", err
	}
	refreshToken, err := s.jwtService.GenerateRefreshToken(user.ID, context)
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

func (s *AuthService) resolveInitialContext(userID uint) (identity.Context, error) {
	var platform models.PlatformMembership
	err := s.db.Where("user_id = ? AND status = ?", userID, "active").First(&platform).Error
	if err == nil {
		membershipID := platform.ID
		return identity.Context{Type: identity.ContextPlatform, MembershipID: &membershipID, Role: platform.Role}, nil
	}
	if err != gorm.ErrRecordNotFound {
		return identity.Context{}, err
	}

	var membership models.Membership
	err = s.db.Where("user_id = ? AND status = ?", userID, "active").
		Order("is_default DESC, id ASC").First(&membership).Error
	if err == nil {
		companyID := membership.CompanyID
		membershipID := membership.ID
		return identity.Context{Type: identity.ContextTenant, CompanyID: &companyID, MembershipID: &membershipID, Role: membership.Role}, nil
	}
	if err != gorm.ErrRecordNotFound {
		return identity.Context{}, err
	}
	var count int64
	if err := s.db.Model(&models.Membership{}).Where("user_id = ?", userID).Count(&count).Error; err != nil {
		return identity.Context{}, err
	}
	var platformCount int64
	if err := s.db.Model(&models.PlatformMembership{}).Where("user_id = ?", userID).Count(&platformCount).Error; err != nil {
		return identity.Context{}, err
	}
	if count+platformCount > 0 {
		return identity.Context{}, apperr.Forbidden("user has no active authorization context")
	}
	return identity.Context{Type: identity.ContextIdentity}, nil
}

func (s *AuthService) resolveRefreshContext(userID uint, claims *RefreshClaims) (identity.Context, error) {
	switch claims.ContextType {
	case identity.ContextPlatform:
		var platform models.PlatformMembership
		if err := s.db.Where("id = ? AND user_id = ? AND status = ?", *claims.MembershipID, userID, "active").First(&platform).Error; err != nil {
			return identity.Context{}, apperr.Forbidden("platform access is no longer active")
		}
		membershipID := platform.ID
		return identity.Context{Type: identity.ContextPlatform, MembershipID: &membershipID, Role: platform.Role}, nil
	case identity.ContextTenant:
		var membership models.Membership
		if err := s.db.Where("id = ? AND user_id = ? AND company_id = ? AND status = ?", *claims.MembershipID, userID, *claims.CompanyID, "active").First(&membership).Error; err != nil {
			return identity.Context{}, apperr.Forbidden("tenant access is no longer active")
		}
		companyID := membership.CompanyID
		membershipID := membership.ID
		return identity.Context{Type: identity.ContextTenant, CompanyID: &companyID, MembershipID: &membershipID, Role: membership.Role}, nil
	case identity.ContextIdentity:
		return identity.Context{Type: identity.ContextIdentity}, nil
	default:
		return identity.Context{}, apperr.Unauthorized("invalid authorization context")
	}
}

func (s *AuthService) getPlatformAccess(userID uint) (*dtos.PlatformAccessResponse, error) {
	var membership models.PlatformMembership
	if err := s.db.Where("user_id = ? AND status = ?", userID, "active").First(&membership).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &dtos.PlatformAccessResponse{MembershipID: membership.ID, Role: membership.Role}, nil
}
