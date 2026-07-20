package services

import (
	"fmt"
	"strings"

	"dvra-api/internal/app/dtos"
	"dvra-api/internal/app/models"
	"dvra-api/internal/app/repositories"
	"dvra-api/internal/shared/apperr"
)

type SystemValueService interface {
	GetCategories(module string) ([]models.SystemValueCategory, error)
	GetByCategory(module, category string, companyID *uint) ([]models.SystemValue, error)
	GetByLegacyCategory(category string, companyID *uint) ([]models.SystemValue, error)
	GetAll() ([]models.SystemValue, error)
	CreateGlobal(dto dtos.CreateSystemValueDTO) (*models.SystemValue, error)
	CreateTenant(companyID uint, dto dtos.CreateSystemValueDTO) (*models.SystemValue, error)
	UpdateGlobal(id uint, dto dtos.UpdateSystemValueDTO) (*models.SystemValue, error)
	UpdateTenant(companyID, id uint, dto dtos.UpdateSystemValueDTO) (*models.SystemValue, error)
	DeleteGlobal(id uint) error
	DeleteTenant(companyID, id uint) error
	UpsertTenantOverride(companyID, valueID uint, dto dtos.UpsertSystemValueOverrideDTO) error
}

type systemValueService struct {
	repo repositories.SystemValueRepository
}

func NewSystemValueService(repo repositories.SystemValueRepository) SystemValueService {
	return &systemValueService{repo: repo}
}

func (s *systemValueService) GetCategories(module string) ([]models.SystemValueCategory, error) {
	return s.repo.GetCategories(module)
}

func (s *systemValueService) GetByCategory(module, category string, companyID *uint) ([]models.SystemValue, error) {
	definition, err := s.repo.GetCategory(module, category)
	if err != nil {
		return nil, err
	}
	if definition == nil {
		return nil, apperr.NotFound("catalog category not found")
	}
	return s.repo.GetByCategoryID(definition.ID, companyID)
}

func (s *systemValueService) GetByLegacyCategory(category string, companyID *uint) ([]models.SystemValue, error) {
	definition, err := s.repo.GetCategoryByKey(category)
	if err != nil {
		return nil, apperr.BadRequest(err.Error())
	}
	if definition == nil {
		return nil, apperr.NotFound("catalog category not found")
	}
	return s.repo.GetByCategoryID(definition.ID, companyID)
}

func (s *systemValueService) GetAll() ([]models.SystemValue, error) { return s.repo.GetAll() }

func (s *systemValueService) CreateGlobal(dto dtos.CreateSystemValueDTO) (*models.SystemValue, error) {
	return s.create(nil, true, dto)
}

func (s *systemValueService) CreateTenant(companyID uint, dto dtos.CreateSystemValueDTO) (*models.SystemValue, error) {
	category, err := s.repo.GetCategoryByID(dto.CategoryID)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, apperr.NotFound("catalog category not found")
	}
	if !category.IsExtensible || !category.AllowTenantValues {
		return nil, apperr.Forbidden("this catalog does not allow tenant values")
	}
	return s.create(&companyID, false, dto)
}

func (s *systemValueService) create(companyID *uint, system bool, dto dtos.CreateSystemValueDTO) (*models.SystemValue, error) {
	key := strings.TrimSpace(dto.Value)
	value := &models.SystemValue{CategoryID: dto.CategoryID, CompanyID: companyID, Key: key, Label: strings.TrimSpace(dto.Label), Description: dto.Description, DisplayOrder: dto.DisplayOrder, IsSystem: system, IsActive: true}
	if err := s.repo.Create(value); err != nil {
		return nil, fmt.Errorf("create catalog value: %w", err)
	}
	return value, nil
}

func applySystemValueUpdate(value *models.SystemValue, dto dtos.UpdateSystemValueDTO) {
	if dto.Label != nil {
		value.Label = strings.TrimSpace(*dto.Label)
	}
	if dto.Description != nil {
		value.Description = dto.Description
	}
	if dto.DisplayOrder != nil {
		value.DisplayOrder = *dto.DisplayOrder
	}
	if dto.IsActive != nil {
		value.IsActive = *dto.IsActive
	}
}

func (s *systemValueService) UpdateGlobal(id uint, dto dtos.UpdateSystemValueDTO) (*models.SystemValue, error) {
	value, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if value == nil || value.CompanyID != nil {
		return nil, apperr.NotFound("global catalog value not found")
	}
	applySystemValueUpdate(value, dto)
	return value, s.repo.Update(value)
}

func (s *systemValueService) UpdateTenant(companyID, id uint, dto dtos.UpdateSystemValueDTO) (*models.SystemValue, error) {
	value, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if value == nil || value.CompanyID == nil || *value.CompanyID != companyID {
		return nil, apperr.NotFound("tenant catalog value not found")
	}
	applySystemValueUpdate(value, dto)
	return value, s.repo.Update(value)
}

func (s *systemValueService) DeleteGlobal(id uint) error {
	value, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if value == nil || value.CompanyID != nil {
		return apperr.NotFound("global catalog value not found")
	}
	return s.repo.Delete(id)
}

func (s *systemValueService) DeleteTenant(companyID, id uint) error {
	value, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if value == nil || value.CompanyID == nil || *value.CompanyID != companyID {
		return apperr.NotFound("tenant catalog value not found")
	}
	return s.repo.Delete(id)
}

func (s *systemValueService) UpsertTenantOverride(companyID, valueID uint, dto dtos.UpsertSystemValueOverrideDTO) error {
	value, err := s.repo.GetByID(valueID)
	if err != nil {
		return err
	}
	if value == nil || value.CompanyID != nil {
		return apperr.NotFound("global catalog value not found")
	}
	if value.Category == nil || !value.Category.AllowTenantOverride {
		return apperr.Forbidden("this catalog does not allow tenant overrides")
	}
	return s.repo.UpsertOverride(&models.CompanySystemValueOverride{CompanyID: companyID, SystemValueID: valueID, CustomLabel: dto.CustomLabel, DisplayOrder: dto.DisplayOrder, IsEnabled: dto.IsEnabled})
}
