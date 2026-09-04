// Package service contiene los casos de uso del módulo staffing.
// Depende solo de domain (puertos), models y apperr — nunca de gin ni gorm.
package service

import (
	"dvra-api/internal/modules/staffing/domain"
	"dvra-api/internal/shared/apperr"
	"dvra-api/internal/shared/models"
	"dvra-api/internal/shared/tenant"
)

// StaffingClientService orquesta la lógica de clientes finales.
type StaffingClientService struct {
	repo domain.StaffingClientRepository
}

func NewStaffingClientService(repo domain.StaffingClientRepository) *StaffingClientService {
	return &StaffingClientService{repo: repo}
}

func (s *StaffingClientService) GetByCompanyID(companyID uint, filters domain.StaffingClientFilters) ([]models.StaffingClient, error) {
	return s.repo.GetByCompanyID(companyID, filters.Status)
}

func (s *StaffingClientService) GetByID(id uint, scope tenant.Scope) (*models.StaffingClient, error) {
	client, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if client == nil {
		return nil, apperr.NotFound("staffing client not found")
	}
	if !scope.Allows(client.CompanyID) {
		return nil, apperr.Forbidden("access denied")
	}
	return client, nil
}

func (s *StaffingClientService) Create(dto CreateStaffingClientDTO) (*models.StaffingClient, error) {
	exists, err := s.repo.ExistsBySlug(dto.CompanyID, dto.Slug, 0)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, apperr.Conflict("a staffing client with this slug already exists")
	}

	status := dto.Status
	if status == "" {
		status = "active"
	}

	client := &models.StaffingClient{
		CompanyID:    dto.CompanyID,
		Name:         dto.Name,
		Slug:         dto.Slug,
		Industry:     dto.Industry,
		Website:      dto.Website,
		LogoURL:      dto.LogoURL,
		ContactName:  dto.ContactName,
		ContactEmail: dto.ContactEmail,
		ContactPhone: dto.ContactPhone,
		Status:       status,
		Notes:        dto.Notes,
	}
	return s.repo.Create(client)
}

func (s *StaffingClientService) Update(id uint, scope tenant.Scope, dto UpdateStaffingClientDTO) (*models.StaffingClient, error) {
	client, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if client == nil {
		return nil, apperr.NotFound("staffing client not found")
	}
	if !scope.Allows(client.CompanyID) {
		return nil, apperr.Forbidden("access denied")
	}

	if dto.Slug != nil && *dto.Slug != client.Slug {
		exists, err := s.repo.ExistsBySlug(client.CompanyID, *dto.Slug, client.ID)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, apperr.Conflict("a staffing client with this slug already exists")
		}
		client.Slug = *dto.Slug
	}
	if dto.Name != nil {
		client.Name = *dto.Name
	}
	if dto.Industry != nil {
		client.Industry = *dto.Industry
	}
	if dto.Website != nil {
		client.Website = *dto.Website
	}
	if dto.LogoURL != nil {
		client.LogoURL = *dto.LogoURL
	}
	if dto.ContactName != nil {
		client.ContactName = *dto.ContactName
	}
	if dto.ContactEmail != nil {
		client.ContactEmail = *dto.ContactEmail
	}
	if dto.ContactPhone != nil {
		client.ContactPhone = *dto.ContactPhone
	}
	if dto.Status != nil {
		client.Status = *dto.Status
	}
	if dto.Notes != nil {
		client.Notes = *dto.Notes
	}

	return s.repo.Update(client)
}

func (s *StaffingClientService) Delete(id uint, scope tenant.Scope) error {
	client, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if client == nil {
		return apperr.NotFound("staffing client not found")
	}
	if !scope.Allows(client.CompanyID) {
		return apperr.Forbidden("access denied")
	}
	return s.repo.Delete(id)
}
