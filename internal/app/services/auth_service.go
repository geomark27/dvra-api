package services

import (
	"errors"

	"dvra-api/internal/app/dtos"
	"dvra-api/internal/app/models"
	"dvra-api/internal/app/repositories"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrEmailExists        = errors.New("email already exists")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidPassword    = errors.New("invalid password")
	ErrCompanyNotFound    = errors.New("company not found")
	ErrNoMembership       = errors.New("user does not belong to this company")
	ErrNotSuperAdmin      = errors.New("user is not a superadmin")
)

// AuthService handles authentication business logic
type AuthService struct {
	userRepo   repositories.UserRepository
	jwtService JWTService
	db         *gorm.DB
}

// NewAuthService creates a new auth service
func NewAuthService(userRepo repositories.UserRepository, jwtService JWTService, db *gorm.DB) *AuthService {
	return &AuthService{
		userRepo:   userRepo,
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
	accessToken, err := s.jwtService.GenerateAccessToken(createdUser.ID, nil, createdUser.Email, "user")
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.jwtService.GenerateRefreshToken(createdUser.ID)
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
	}, nil
}

// Login authenticates a user and returns tokens
func (s *AuthService) Login(dto *dtos.LoginDTO) (*dtos.LoginResponseDTO, error) {
	// Find user by email
	user, err := s.userRepo.FindByEmail(dto.Email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// Check if user is active
	if !user.IsActive {
		return nil, errors.New("account is inactive")
	}

	// Verify password
	if err := ComparePassword(user.PasswordHash, dto.Password); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Get user's default membership to include company context
	var membership models.Membership
	var companyID *uint
	var role string = "user"

	err = s.db.Where("user_id = ? AND is_default = ?", user.ID, true).
		First(&membership).Error

	if err == nil {
		// User has a default company
		companyID = membership.CompanyID
		role = membership.Role
	} else if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// Generate tokens
	accessToken, err := s.jwtService.GenerateAccessToken(user.ID, companyID, user.Email, role)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.jwtService.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	// Update last login
	s.userRepo.UpdateLastLogin(user.ID)

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
	}, nil
}

// RefreshToken generates new access token from refresh token
func (s *AuthService) RefreshToken(dto *dtos.RefreshTokenDTO) (*dtos.RefreshTokenResponseDTO, error) {
	// Validate refresh token
	claims, err := s.jwtService.ValidateToken(dto.RefreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	// Get user
	user, err := s.userRepo.FindByID(int(claims.UserID))
	if err != nil {
		return nil, ErrUserNotFound
	}

	if !user.IsActive {
		return nil, errors.New("account is inactive")
	}

	// Get user's default membership
	var membership models.Membership
	var companyID *uint
	var role string = "user"

	err = s.db.Where("user_id = ? AND is_default = ?", user.ID, true).
		First(&membership).Error

	if err == nil {
		companyID = membership.CompanyID
		role = membership.Role
	}

	// Generate new tokens
	accessToken, err := s.jwtService.GenerateAccessToken(user.ID, companyID, user.Email, role)
	if err != nil {
		return nil, err
	}

	newRefreshToken, err := s.jwtService.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &dtos.RefreshTokenResponseDTO{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

// SuperAdminLogin authenticates a superadmin user
func (s *AuthService) SuperAdminLogin(dto *dtos.SuperAdminLoginDTO) (*dtos.SuperAdminLoginResponseDTO, error) {
	// Find user by email
	user, err := s.userRepo.FindByEmail(dto.Email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// Check if user is active
	if !user.IsActive {
		return nil, errors.New("account is inactive")
	}

	// Verify user is superadmin
	if !user.IsSuperAdmin {
		return nil, ErrNotSuperAdmin
	}

	// Verify password
	if err := ComparePassword(user.PasswordHash, dto.Password); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Generate tokens WITHOUT company_id (superadmin is global)
	accessToken, err := s.jwtService.GenerateAccessToken(user.ID, nil, user.Email, models.RoleSuperAdmin)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.jwtService.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	// Update last login
	s.userRepo.UpdateLastLogin(user.ID)

	return &dtos.SuperAdminLoginResponseDTO{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: dtos.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			IsActive:  user.IsActive,
		},
		IsSuperAdmin: true,
	}, nil
}

// ChangePassword changes user password
func (s *AuthService) ChangePassword(userID uint, dto *dtos.ChangePasswordDTO) error {
	// Get user
	user, err := s.userRepo.FindByID(int(userID))
	if err != nil {
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
	if err != nil {
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

	// Start transaction
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. Create company
	timezone := dto.Timezone
	if timezone == "" {
		timezone = "America/Bogota"
	}

	company := models.Company{
		Name:     dto.CompanyName,
		Slug:     dto.CompanySlug,
		PlanTier: "free", // Default plan
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
		CompanyID: &company.ID,
		Role:      models.RoleAdmin,
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

	// Generate tokens with company context
	accessToken, err := s.jwtService.GenerateAccessToken(admin.ID, &company.ID, admin.Email, models.RoleAdmin)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.jwtService.GenerateRefreshToken(admin.ID)
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
				ID:       membership.Company.ID,
				Name:     membership.Company.Name,
				Slug:     membership.Company.Slug,
				PlanTier: membership.Company.PlanTier,
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

	// Generate new token with new company context
	accessToken, err := s.jwtService.GenerateAccessToken(userID, &dto.CompanyID, user.Email, membership.Role)
	if err != nil {
		return nil, err
	}

	return &dtos.SwitchCompanyResponseDTO{
		AccessToken: accessToken,
		Company: dtos.CompanyResponse{
			ID:       membership.Company.ID,
			Name:     membership.Company.Name,
			Slug:     membership.Company.Slug,
			PlanTier: membership.Company.PlanTier,
		},
	}, nil
}

// LoginWithCompanies authenticates a user and returns tokens with companies list
func (s *AuthService) LoginWithCompanies(dto *dtos.LoginDTO) (*dtos.LoginResponseWithCompaniesDTO, error) {
	// Find user by email
	user, err := s.userRepo.FindByEmail(dto.Email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// Check if user is active
	if !user.IsActive {
		return nil, errors.New("account is inactive")
	}

	// Verify password
	if err := ComparePassword(user.PasswordHash, dto.Password); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Get user's default membership
	var membership models.Membership
	var companyID *uint
	var role string = "user"

	err = s.db.Where("user_id = ? AND is_default = ?", user.ID, true).
		First(&membership).Error

	if err == nil {
		companyID = membership.CompanyID
		role = membership.Role
	} else if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	// Get all user's companies
	companies, _ := s.GetUserCompanies(user.ID)

	// Generate tokens
	accessToken, err := s.jwtService.GenerateAccessToken(user.ID, companyID, user.Email, role)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.jwtService.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	// Update last login
	s.userRepo.UpdateLastLogin(user.ID)

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
		Companies: companies,
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
