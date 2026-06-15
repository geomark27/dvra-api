package services

import (
	"fmt"

	"dvra-api/internal/app/dtos"
	"dvra-api/internal/app/models"
	"dvra-api/internal/app/repositories"
)

// StaffingClientService define el contrato del servicio de clientes finales
type StaffingClientService interface {
	GetByCompanyID(companyID uint, filters dtos.StaffingClientFilters) ([]models.StaffingClient, error)
	GetByID(id uint) (*models.StaffingClient, error)
	Create(dto dtos.CreateStaffingClientDTO) (*models.StaffingClient, error)
	Update(id uint, dto dtos.UpdateStaffingClientDTO) (*models.StaffingClient, error)
	Delete(id uint) error
}

type staffingClientService struct {
	repo repositories.StaffingClientRepository
}

// NewStaffingClientService crea una nueva instancia de StaffingClientService
func NewStaffingClientService(repo repositories.StaffingClientRepository) StaffingClientService {
	return &staffingClientService{repo: repo}
}

func (s *staffingClientService) GetByCompanyID(companyID uint, filters dtos.StaffingClientFilters) ([]models.StaffingClient, error) {
	return s.repo.GetByCompanyID(companyID, filters.Status)
}

func (s *staffingClientService) GetByID(id uint) (*models.StaffingClient, error) {
	client, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if client == nil {
		return nil, fmt.Errorf("staffing client not found")
	}
	return client, nil
}

func (s *staffingClientService) Create(dto dtos.CreateStaffingClientDTO) (*models.StaffingClient, error) {
	// El slug debe ser único dentro de la empresa
	exists, err := s.repo.ExistsBySlug(dto.CompanyID, dto.Slug, 0)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("a staffing client with this slug already exists")
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

func (s *staffingClientService) Update(id uint, dto dtos.UpdateStaffingClientDTO) (*models.StaffingClient, error) {
	client, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if client == nil {
		return nil, fmt.Errorf("staffing client not found")
	}

	if dto.Slug != nil && *dto.Slug != client.Slug {
		exists, err := s.repo.ExistsBySlug(client.CompanyID, *dto.Slug, client.ID)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, fmt.Errorf("a staffing client with this slug already exists")
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

func (s *staffingClientService) Delete(id uint) error {
	client, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if client == nil {
		return fmt.Errorf("staffing client not found")
	}
	return s.repo.Delete(id)
}
