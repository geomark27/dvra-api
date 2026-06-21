package services

import (
	"fmt"
	"os"
	"path/filepath"

	"dvra-api/internal/app/dtos"
	"dvra-api/internal/app/models"
	"dvra-api/internal/app/repositories"
	"dvra-api/internal/shared/apperr"
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
		return nil, apperr.NotFound("company not found")
	}
	return company, nil
}

func (s *companyService) GetCompanyBySlug(slug string) (*models.Company, error) {
	company, err := s.companyRepo.GetBySlug(slug)
	if err != nil {
		return nil, err
	}
	if company == nil {
		return nil, apperr.NotFound("company not found")
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
		return nil, apperr.Conflict(fmt.Sprintf("company with slug '%s' already exists", dto.Slug))
	}

	company := &models.Company{
		Name:        dto.Name,
		Slug:        dto.Slug,
		LogoURL:     dto.LogoURL,
		PlanTier:    dto.PlanTier,
		TrialEndsAt: dto.TrialEndsAt,
		Timezone:    dto.Timezone,
	}

	// Crear la empresa en la base de datos
	createdCompany, err := s.companyRepo.Create(company)
	if err != nil {
		return nil, err
	}

	// Crear estructura de directorios para la empresa
	if err := createCompanyDirectories(createdCompany.Slug); err != nil {
		// Log error pero no fallar - los directorios se crearán al subir archivos
		fmt.Printf("Warning: failed to create company directories for %s: %v\n", createdCompany.Slug, err)
	}

	return createdCompany, nil
}

// createCompanyDirectories crea la estructura de directorios para una empresa
func createCompanyDirectories(slug string) error {
	baseDir := filepath.Join("uploads", "companies", slug)

	// Crear subdirectorios: logo, resumes, documents
	dirs := []string{
		filepath.Join(baseDir, "logo"),
		filepath.Join(baseDir, "resumes"),
		filepath.Join(baseDir, "documents"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}

func (s *companyService) UpdateCompany(id uint, dto dtos.UpdateCompanyDTO) (*models.Company, error) {
	company, err := s.companyRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if company == nil {
		return nil, apperr.NotFound("company not found")
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
		return apperr.NotFound("company not found")
	}

	return s.companyRepo.Delete(id)
}

func (s *companyService) GetCompanyWithMembers(companyID uint) (*models.Company, error) {
	return s.companyRepo.GetCompaniesWithMembers(companyID)
}
