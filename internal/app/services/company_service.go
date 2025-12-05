package services

import (
	"fmt"

	"dvra-api/internal/app/dtos"
	"dvra-api/internal/app/models"
	"dvra-api/internal/app/repositories"
)

// CompanyService define el contrato del servicio de companies
type CompanyService interface {
	GetAllCompanies() ([]models.Company, error)
	GetCompanyByID(id uint) (*models.Company, error)
	GetCompanyBySlug(slug string) (*models.Company, error)
	CreateCompany(dto dtos.CreateCompanyDTO) (*models.Company, error)
	UpdateCompany(id uint, dto dtos.UpdateCompanyDTO) (*models.Company, error)
	DeleteCompany(id uint) error
	GetCompanyWithMembers(companyID uint) (*models.Company, error)
}

// companyService es la implementación privada del servicio
type companyService struct {
	companyRepo repositories.CompanyRepository
}

// NewCompanyService crea una nueva instancia de CompanyService
func NewCompanyService(companyRepo repositories.CompanyRepository) CompanyService {
	return &companyService{
		companyRepo: companyRepo,
	}
}

func (s *companyService) GetAllCompanies() ([]models.Company, error) {
	return s.companyRepo.GetAll()
}

func (s *companyService) GetCompanyByID(id uint) (*models.Company, error) {
	company, err := s.companyRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if company == nil {
		return nil, fmt.Errorf("company not found")
	}
	return company, nil
}

func (s *companyService) GetCompanyBySlug(slug string) (*models.Company, error) {
	company, err := s.companyRepo.GetBySlug(slug)
	if err != nil {
		return nil, err
	}
	if company == nil {
		return nil, fmt.Errorf("company not found")
	}
	return company, nil
}

func (s *companyService) CreateCompany(dto dtos.CreateCompanyDTO) (*models.Company, error) {
	// Verificar que el slug no exista
	existing, err := s.companyRepo.GetBySlug(dto.Slug)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, fmt.Errorf("company with slug '%s' already exists", dto.Slug)
	}

	company := &models.Company{
		Name:        dto.Name,
		Slug:        dto.Slug,
		LogoURL:     dto.LogoURL,
		PlanTier:    dto.PlanTier,
		TrialEndsAt: dto.TrialEndsAt,
		Timezone:    dto.Timezone,
	}

	return s.companyRepo.Create(company)
}

func (s *companyService) UpdateCompany(id uint, dto dtos.UpdateCompanyDTO) (*models.Company, error) {
	company, err := s.companyRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if company == nil {
		return nil, fmt.Errorf("company not found")
	}

	// Actualizar solo campos proporcionados
	if dto.Name != nil {
		company.Name = *dto.Name
	}
	if dto.Slug != nil {
		company.Slug = *dto.Slug
	}
	if dto.LogoURL != nil {
		company.LogoURL = *dto.LogoURL
	}
	if dto.PlanTier != nil {
		company.PlanTier = *dto.PlanTier
	}
	if dto.TrialEndsAt != nil {
		company.TrialEndsAt = dto.TrialEndsAt
	}
	if dto.Timezone != nil {
		company.Timezone = *dto.Timezone
	}

	return s.companyRepo.Update(company)
}

func (s *companyService) DeleteCompany(id uint) error {
	company, err := s.companyRepo.GetByID(id)
	if err != nil {
		return err
	}
	if company == nil {
		return fmt.Errorf("company not found")
	}

	return s.companyRepo.Delete(id)
}

func (s *companyService) GetCompanyWithMembers(companyID uint) (*models.Company, error) {
	return s.companyRepo.GetCompaniesWithMembers(companyID)
}
