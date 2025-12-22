package services

import (
	"dvra-api/internal/app/dtos"
	"dvra-api/internal/app/models"
	"dvra-api/internal/app/repositories"
	"errors"
)

type SystemValueService interface {
	GetByCategory(category string, companyID *uint) ([]models.SystemValue, error)
	GetAll() ([]models.SystemValue, error)
	Create(dto dtos.CreateSystemValueDTO) (*models.SystemValue, error)
	Update(id uint, dto dtos.UpdateSystemValueDTO) (*models.SystemValue, error)
	Delete(id uint) error
}

type systemValueService struct {
	repo repositories.SystemValueRepository
}

func NewSystemValueService(repo repositories.SystemValueRepository) SystemValueService {
	return &systemValueService{repo: repo}
}

func (s *systemValueService) GetByCategory(category string, companyID *uint) ([]models.SystemValue, error) {
	return s.repo.GetByCategory(category, companyID)
}

func (s *systemValueService) GetAll() ([]models.SystemValue, error) {
	return s.repo.GetAll()
}

func (s *systemValueService) Create(dto dtos.CreateSystemValueDTO) (*models.SystemValue, error) {
	value := &models.SystemValue{
		Category:     dto.Category,
		Value:        dto.Value,
		Label:        dto.Label,
		Description:  dto.Description,
		DisplayOrder: dto.DisplayOrder,
		CompanyID:    dto.CompanyID,
		IsActive:     true,
	}

	if err := s.repo.Create(value); err != nil {
		return nil, err
	}

	return value, nil
}

func (s *systemValueService) Update(id uint, dto dtos.UpdateSystemValueDTO) (*models.SystemValue, error) {
	value, err := s.repo.GetByID(id)
	if err != nil {
		return nil, errors.New("system value not found")
	}

	value.Label = dto.Label
	value.Description = dto.Description
	value.DisplayOrder = dto.DisplayOrder
	value.IsActive = dto.IsActive

	if err := s.repo.Update(value); err != nil {
		return nil, err
	}

	return value, nil
}

func (s *systemValueService) Delete(id uint) error {
	return s.repo.Delete(id)
}
