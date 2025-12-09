package admin

import (
	"dvra-api/internal/app/dtos"
	"dvra-api/internal/app/models"
	"dvra-api/internal/app/repositories"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// SuperAdminCompaniesService handles global company management
type SuperAdminCompaniesService struct {
	db             *gorm.DB
	companyRepo    repositories.CompanyRepository
	userRepo       repositories.UserRepository
	membershipRepo repositories.MembershipRepository
	planRepo       repositories.PlanRepository
}

// NewSuperAdminCompaniesService creates a new SuperAdminCompaniesService
func NewSuperAdminCompaniesService(
	db *gorm.DB,
	companyRepo repositories.CompanyRepository,
	userRepo repositories.UserRepository,
	membershipRepo repositories.MembershipRepository,
	planRepo repositories.PlanRepository,
) *SuperAdminCompaniesService {
	return &SuperAdminCompaniesService{
		db:             db,
		companyRepo:    companyRepo,
		userRepo:       userRepo,
		membershipRepo: membershipRepo,
		planRepo:       planRepo,
	}
}

// GetAllCompanies retrieves all companies with pagination and filters (no company scoping)
func (s *SuperAdminCompaniesService) GetAllCompanies(page, limit int, search, planTier string) ([]dtos.CompanyWithStatsDTO, int64, error) {
	var companies []models.Company
	var total int64

	query := s.db.Model(&models.Company{})

	// Apply filters
	if search != "" {
		query = query.Where("name ILIKE ? OR slug ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if planTier != "" {
		query = query.Where("plan_tier = ?", planTier)
	}

	// Count total
	query.Count(&total)

	// Paginate
	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&companies).Error; err != nil {
		return nil, 0, err
	}

	// Build response with stats
	result := make([]dtos.CompanyWithStatsDTO, len(companies))
	for i, company := range companies {
		// Count users
		var userCount int64
		s.db.Model(&models.Membership{}).Where("company_id = ?", company.ID).Count(&userCount)

		// Count jobs
		var jobCount int64
		s.db.Model(&models.Job{}).Where("company_id = ?", company.ID).Count(&jobCount)

		var trialEndsAt *string
		if company.TrialEndsAt != nil {
			formatted := company.TrialEndsAt.Format(time.RFC3339)
			trialEndsAt = &formatted
		}

		result[i] = dtos.CompanyWithStatsDTO{
			ID:          company.ID,
			Name:        company.Name,
			Slug:        company.Slug,
			PlanTier:    company.PlanTier,
			Status:      company.PlanTier, // or add a separate Status field
			UserCount:   int(userCount),
			JobCount:    int(jobCount),
			CreatedAt:   company.CreatedAt.Format(time.RFC3339),
			TrialEndsAt: trialEndsAt,
		}
	}

	return result, total, nil
}

// CreateCompanyWithAdmin creates a new company and its first admin user
func (s *SuperAdminCompaniesService) CreateCompanyWithAdmin(dto dtos.CreateCompanyWithAdminDTO) (*models.Company, *models.User, error) {
	// Check if admin email already exists
	existingUser, _ := s.userRepo.FindByEmail(dto.AdminEmail)
	if existingUser != nil {
		return nil, nil, errors.New("admin email already exists")
	}

	// Determine which plan to use (default to "free" if not specified)
	planSlug := "free"
	if dto.PlanSlug != nil && *dto.PlanSlug != "" {
		planSlug = *dto.PlanSlug
	}

	// Verify that the plan exists and is active
	plan, err := s.planRepo.FindActiveBySlug(planSlug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, errors.New("plan '" + planSlug + "' is not available or inactive")
		}
		return nil, nil, err
	}

	// Start transaction
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. Create company with validated plan
	trialEnds := time.Now().AddDate(0, 1, 0) // 1 month trial
	company := models.Company{
		Name:        dto.CompanyName,
		Slug:        dto.CompanySlug,
		PlanTier:    plan.Slug, // Use validated plan slug
		TrialEndsAt: &trialEnds,
		Timezone:    "America/Bogota",
	}

	if err := tx.Create(&company).Error; err != nil {
		tx.Rollback()
		return nil, nil, err
	}

	// 2. Create admin user
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(dto.AdminPassword), bcrypt.DefaultCost)
	if err != nil {
		tx.Rollback()
		return nil, nil, err
	}

	admin := models.User{
		Email:        dto.AdminEmail,
		PasswordHash: string(hashedPassword),
		FirstName:    dto.AdminFirstName,
		LastName:     dto.AdminLastName,
		IsActive:     true,
	}

	if err := tx.Create(&admin).Error; err != nil {
		tx.Rollback()
		return nil, nil, err
	}

	// 3. Create membership
	joinedAt := time.Now()
	membership := models.Membership{
		UserID:    admin.ID,
		CompanyID: &company.ID,
		Role:      models.RoleAdmin,
		Status:    "active",
		IsDefault: true,
		JoinedAt:  &joinedAt,
	}

	if err := tx.Create(&membership).Error; err != nil {
		tx.Rollback()
		return nil, nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, nil, err
	}

	return &company, &admin, nil
}

// ChangeCompanyPlan updates the plan tier of a company
func (s *SuperAdminCompaniesService) ChangeCompanyPlan(companyID uint, newPlan string) error {
	return s.db.Model(&models.Company{}).
		Where("id = ?", companyID).
		Update("plan_tier", newPlan).Error
}

// SuspendCompany suspends a company (sets plan to suspended)
func (s *SuperAdminCompaniesService) SuspendCompany(companyID uint, reason string) error {
	// Update plan tier to indicate suspension
	// In the future, you could add a "suspend_reason" field to Company model
	return s.db.Model(&models.Company{}).
		Where("id = ?", companyID).
		Update("plan_tier", "suspended").Error
}

// GetCompanyUsers retrieves all users from a specific company
func (s *SuperAdminCompaniesService) GetCompanyUsers(companyID uint) ([]models.User, error) {
	var users []models.User
	err := s.db.
		Joins("JOIN memberships ON memberships.user_id = users.id").
		Where("memberships.company_id = ?", companyID).
		Find(&users).Error
	return users, err
}

// GetGlobalAnalytics returns system-wide analytics
func (s *SuperAdminCompaniesService) GetGlobalAnalytics() (*dtos.GlobalAnalyticsDTO, error) {
	analytics := &dtos.GlobalAnalyticsDTO{}

	// Count companies by status
	var totalCompanies, activeCompanies, suspendedCompanies int64
	s.db.Model(&models.Company{}).Count(&totalCompanies)
	s.db.Model(&models.Company{}).Where("plan_tier != ?", "suspended").Count(&activeCompanies)
	s.db.Model(&models.Company{}).Where("plan_tier = ?", "suspended").Count(&suspendedCompanies)

	analytics.TotalCompanies = int(totalCompanies)
	analytics.ActiveCompanies = int(activeCompanies)
	analytics.SuspendedCompanies = int(suspendedCompanies)

	// Count total users
	var totalUsers int64
	s.db.Model(&models.User{}).Count(&totalUsers)
	analytics.TotalUsers = int(totalUsers)

	// Count total jobs
	var totalJobs int64
	s.db.Model(&models.Job{}).Count(&totalJobs)
	analytics.TotalJobs = int(totalJobs)

	// Count total applications
	var totalApplications int64
	s.db.Model(&models.Application{}).Count(&totalApplications)
	analytics.TotalApplications = int(totalApplications)

	// Calculate monthly revenue (simplified)
	// TODO: Implement actual billing tracking
	analytics.MonthlyRevenue = calculateMonthlyRevenue(s.db)

	// Calculate churn rate (simplified)
	// TODO: Implement actual churn tracking
	analytics.ChurnRate = 0.0

	return analytics, nil
}

// calculateMonthlyRevenue is a helper to estimate MRR
func calculateMonthlyRevenue(db *gorm.DB) float64 {
	var companies []models.Company
	db.Where("plan_tier != ?", "suspended").Find(&companies)

	planPrices := map[string]float64{
		"free":         0,
		"professional": 149,
		"enterprise":   399,
	}

	var mrr float64
	for _, company := range companies {
		if price, exists := planPrices[company.PlanTier]; exists {
			mrr += price
		}
	}

	return mrr
}
