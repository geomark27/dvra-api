package dtos

// RegisterDTO represents the registration request
type RegisterDTO struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

// LoginDTO represents the login request
type LoginDTO struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponseDTO represents the login response
type LoginResponseDTO struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	User         UserResponse `json:"user"`
}

// UserResponse represents user data in auth response
type UserResponse struct {
	ID        uint   `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	IsActive  bool   `json:"is_active"`
}

// RefreshTokenDTO represents the refresh token request
type RefreshTokenDTO struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// RefreshTokenResponseDTO represents the refresh token response
type RefreshTokenResponseDTO struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// ChangePasswordDTO represents the change password request
type ChangePasswordDTO struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

// RegisterCompanyDTO represents company registration with admin user
type RegisterCompanyDTO struct {
	CompanyName    string `json:"company_name" binding:"required"`
	CompanySlug    string `json:"company_slug" binding:"required"`
	AdminEmail     string `json:"admin_email" binding:"required,email"`
	AdminPassword  string `json:"admin_password" binding:"required,min=8"`
	AdminFirstName string `json:"admin_first_name" binding:"required"`
	AdminLastName  string `json:"admin_last_name" binding:"required"`
	Timezone       string `json:"timezone"`
}

// RegisterCompanyResponseDTO represents company registration response
type RegisterCompanyResponseDTO struct {
	AccessToken  string          `json:"access_token"`
	RefreshToken string          `json:"refresh_token"`
	Company      CompanyResponse `json:"company"`
	Admin        UserResponse    `json:"admin"`
}

// CompanyResponse represents company data in response
type CompanyResponse struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Slug     string `json:"slug"`
	PlanTier string `json:"plan_tier"`
}

// LoginResponseWithCompaniesDTO represents login response with user's companies
type LoginResponseWithCompaniesDTO struct {
	AccessToken  string            `json:"access_token"`
	RefreshToken string            `json:"refresh_token"`
	User         UserResponse      `json:"user"`
	Companies    []CompanyResponse `json:"companies"`
}

// SwitchCompanyDTO represents switch company request
type SwitchCompanyDTO struct {
	CompanyID uint `json:"company_id" binding:"required"`
}

// SwitchCompanyResponseDTO represents switch company response
type SwitchCompanyResponseDTO struct {
	AccessToken string          `json:"access_token"`
	Company     CompanyResponse	 `json:"company"`
}

// SuperAdminLoginDTO represents superadmin login request
type SuperAdminLoginDTO struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// SuperAdminLoginResponseDTO represents superadmin login response
type SuperAdminLoginResponseDTO struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	User         UserResponse `json:"user"`
	IsSuperAdmin bool         `json:"is_superadmin"`
}
