package services

import (
	"errors"

	"dvra-api/internal/app/dtos"
	"dvra-api/internal/app/models"
	"dvra-api/internal/app/repositories"

	"gorm.io/gorm"
)

var (
	ErrPlanNotFound    = errors.New("plan not found")
	ErrPlanSlugExists  = errors.New("plan with this slug already exists")
	ErrPlanInUse       = errors.New("plan is currently in use by companies")
	ErrInvalidPlanData = errors.New("invalid plan data")
)

// PlanService interface defines plan business logic
type PlanService interface {
	CreatePlan(dto *dtos.PlanDTO) (*dtos.PlanResponse, error)
	GetPlanByID(id uint) (*dtos.PlanResponse, error)
	GetPlanBySlug(slug string) (*dtos.PlanResponse, error)
	GetAllPlans() ([]dtos.PlanResponse, error)
	GetActivePlans() ([]dtos.PlanResponse, error)
	GetPublicPlans() ([]dtos.PlanResponse, error)
	UpdatePlan(id uint, dto *dtos.UpdatePlanDTO) (*dtos.PlanResponse, error)
	TogglePlanStatus(id uint, isActive bool) (*dtos.PlanResponse, error)
	DeletePlan(id uint) error
	AssignPlanToCompany(companyID uint, planSlug string) error
}

type planService struct {
	planRepo    repositories.PlanRepository
	companyRepo repositories.CompanyRepository
	db          *gorm.DB
}

// NewPlanService creates a new plan service
func NewPlanService(
	planRepo repositories.PlanRepository,
	companyRepo repositories.CompanyRepository,
	db *gorm.DB,
) PlanService {
	return &planService{
		planRepo:    planRepo,
		companyRepo: companyRepo,
		db:          db,
	}
}

// CreatePlan creates a new plan
func (s *planService) CreatePlan(dto *dtos.PlanDTO) (*dtos.PlanResponse, error) {
	// Check if slug already exists
	exists, err := s.planRepo.ExistsBySlug(dto.Slug)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrPlanSlugExists
	}

	plan := &models.Plan{
		Name:               dto.Name,
		Slug:               dto.Slug,
		Description:        dto.Description,
		Price:              dto.Price,
		Currency:           dto.Currency,
		BillingCycle:       dto.BillingCycle,
		IsActive:           dto.IsActive,
		IsPublic:           dto.IsPublic,
		TrialDays:          dto.TrialDays,
		DisplayOrder:       dto.DisplayOrder,
		MaxUsers:           dto.MaxUsers,
		MaxJobs:            dto.MaxJobs,
		MaxCandidates:      dto.MaxCandidates,
		MaxApplications:    dto.MaxApplications,
		MaxStorageGB:       dto.MaxStorageGB,
		CanExportData:      dto.CanExportData,
		CanUseCustomBrand:  dto.CanUseCustomBrand,
		CanUseAPI:          dto.CanUseAPI,
		CanUseIntegrations: dto.CanUseIntegrations,
		SupportLevel:       dto.SupportLevel,
	}

	createdPlan, err := s.planRepo.Create(plan)
	if err != nil {
		return nil, err
	}

	return s.toPlanResponse(createdPlan), nil
}

// GetPlanByID retrieves a plan by ID
func (s *planService) GetPlanByID(id uint) (*dtos.PlanResponse, error) {
	plan, err := s.planRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPlanNotFound
		}
		return nil, err
	}
	return s.toPlanResponse(plan), nil
}

// GetPlanBySlug retrieves a plan by slug
func (s *planService) GetPlanBySlug(slug string) (*dtos.PlanResponse, error) {
	plan, err := s.planRepo.FindBySlug(slug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPlanNotFound
		}
		return nil, err
	}
	return s.toPlanResponse(plan), nil
}

// GetAllPlans retrieves all plans (SuperAdmin only)
func (s *planService) GetAllPlans() ([]dtos.PlanResponse, error) {
	plans, err := s.planRepo.FindAll()
	if err != nil {
		return nil, err
	}
	return s.toPlanResponseList(plans), nil
}

// GetActivePlans retrieves all active plans
func (s *planService) GetActivePlans() ([]dtos.PlanResponse, error) {
	plans, err := s.planRepo.FindActive()
	if err != nil {
		return nil, err
	}
	return s.toPlanResponseList(plans), nil
}

// GetPublicPlans retrieves all public plans (for pricing page)
func (s *planService) GetPublicPlans() ([]dtos.PlanResponse, error) {
	plans, err := s.planRepo.FindPublic()
	if err != nil {
		return nil, err
	}
	return s.toPlanResponseList(plans), nil
}

// UpdatePlan updates an existing plan
func (s *planService) UpdatePlan(id uint, dto *dtos.UpdatePlanDTO) (*dtos.PlanResponse, error) {
	plan, err := s.planRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPlanNotFound
		}
		return nil, err
	}

	// Update only provided fields
	if dto.Name != nil {
		plan.Name = *dto.Name
	}
	if dto.Description != nil {
		plan.Description = *dto.Description
	}
	if dto.Price != nil {
		plan.Price = *dto.Price
	}
	if dto.Currency != nil {
		plan.Currency = *dto.Currency
	}
	if dto.BillingCycle != nil {
		plan.BillingCycle = *dto.BillingCycle
	}
	if dto.IsActive != nil {
		plan.IsActive = *dto.IsActive
	}
	if dto.IsPublic != nil {
		plan.IsPublic = *dto.IsPublic
	}
	if dto.TrialDays != nil {
		plan.TrialDays = *dto.TrialDays
	}
	if dto.DisplayOrder != nil {
		plan.DisplayOrder = *dto.DisplayOrder
	}
	if dto.MaxUsers != nil {
		plan.MaxUsers = *dto.MaxUsers
	}
	if dto.MaxJobs != nil {
		plan.MaxJobs = *dto.MaxJobs
	}
	if dto.MaxCandidates != nil {
		plan.MaxCandidates = *dto.MaxCandidates
	}
	if dto.MaxApplications != nil {
		plan.MaxApplications = *dto.MaxApplications
	}
	if dto.MaxStorageGB != nil {
		plan.MaxStorageGB = *dto.MaxStorageGB
	}
	if dto.CanExportData != nil {
		plan.CanExportData = *dto.CanExportData
	}
	if dto.CanUseCustomBrand != nil {
		plan.CanUseCustomBrand = *dto.CanUseCustomBrand
	}
	if dto.CanUseAPI != nil {
		plan.CanUseAPI = *dto.CanUseAPI
	}
	if dto.CanUseIntegrations != nil {
		plan.CanUseIntegrations = *dto.CanUseIntegrations
	}
	if dto.SupportLevel != nil {
		plan.SupportLevel = *dto.SupportLevel
	}

	updatedPlan, err := s.planRepo.Update(plan)
	if err != nil {
		return nil, err
	}

	return s.toPlanResponse(updatedPlan), nil
}

// TogglePlanStatus enables/disables a plan
func (s *planService) TogglePlanStatus(id uint, isActive bool) (*dtos.PlanResponse, error) {
	plan, err := s.planRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPlanNotFound
		}
		return nil, err
	}

	plan.IsActive = isActive
	updatedPlan, err := s.planRepo.Update(plan)
	if err != nil {
		return nil, err
	}

	return s.toPlanResponse(updatedPlan), nil
}

// DeletePlan soft deletes a plan
func (s *planService) DeletePlan(id uint) error {
	// Check if plan exists
	_, err := s.planRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrPlanNotFound
		}
		return err
	}

	// Check if plan is in use by companies
	var count int64
	err = s.db.Model(&models.Company{}).Where("plan_tier = ?", id).Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return ErrPlanInUse
	}

	return s.planRepo.Delete(id)
}

// AssignPlanToCompany assigns a plan to a company
func (s *planService) AssignPlanToCompany(companyID uint, planSlug string) error {
	// Verify plan exists and is active
	plan, err := s.planRepo.FindBySlug(planSlug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrPlanNotFound
		}
		return err
	}

	if !plan.IsActive {
		return errors.New("plan is not active")
	}

	// Verify company exists
	company, err := s.companyRepo.GetByID(companyID)
	if err != nil {
		return ErrCompanyNotFound
	}

	// Update company plan
	company.PlanTier = planSlug
	_, err = s.companyRepo.Update(company)
	return err
}

// Helper functions
func (s *planService) toPlanResponse(plan *models.Plan) *dtos.PlanResponse {
	return &dtos.PlanResponse{
		ID:                 plan.ID,
		Name:               plan.Name,
		Slug:               plan.Slug,
		Description:        plan.Description,
		Price:              plan.Price,
		Currency:           plan.Currency,
		BillingCycle:       plan.BillingCycle,
		IsActive:           plan.IsActive,
		IsPublic:           plan.IsPublic,
		TrialDays:          plan.TrialDays,
		DisplayOrder:       plan.DisplayOrder,
		MaxUsers:           plan.MaxUsers,
		MaxJobs:            plan.MaxJobs,
		MaxCandidates:      plan.MaxCandidates,
		MaxApplications:    plan.MaxApplications,
		MaxStorageGB:       plan.MaxStorageGB,
		CanExportData:      plan.CanExportData,
		CanUseCustomBrand:  plan.CanUseCustomBrand,
		CanUseAPI:          plan.CanUseAPI,
		CanUseIntegrations: plan.CanUseIntegrations,
		SupportLevel:       plan.SupportLevel,
		CreatedAt:          plan.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:          plan.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func (s *planService) toPlanResponseList(plans []models.Plan) []dtos.PlanResponse {
	responses := make([]dtos.PlanResponse, len(plans))
	for i, plan := range plans {
		responses[i] = *s.toPlanResponse(&plan)
	}
	return responses
}
