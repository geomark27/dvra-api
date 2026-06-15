package services

import (
	"fmt"

	"dvra-api/internal/app/dtos"
	"dvra-api/internal/app/models"
	"dvra-api/internal/app/repositories"
)

// PlacementService define el contrato del servicio de colocaciones
type PlacementService interface {
	GetByCompanyID(companyID uint, filters dtos.PlacementFilters) ([]models.Placement, error)
	GetByID(id uint) (*models.Placement, error)
	Create(companyID uint, dto dtos.CreatePlacementDTO) (*models.Placement, error)
	Update(id uint, dto dtos.UpdatePlacementDTO) (*models.Placement, error)
	Delete(id uint) error
}

type placementService struct {
	repo            repositories.PlacementRepository
	applicationRepo repositories.ApplicationRepository
	staffingRepo    repositories.StaffingClientRepository
}

// NewPlacementService crea una nueva instancia de PlacementService.
// Recibe los repos de application y staffing client para validar la integridad
// cruzada (mismo tenant) al crear una colocación.
func NewPlacementService(
	repo repositories.PlacementRepository,
	applicationRepo repositories.ApplicationRepository,
	staffingRepo repositories.StaffingClientRepository,
) PlacementService {
	return &placementService{
		repo:            repo,
		applicationRepo: applicationRepo,
		staffingRepo:    staffingRepo,
	}
}

func (s *placementService) GetByCompanyID(companyID uint, filters dtos.PlacementFilters) ([]models.Placement, error) {
	return s.repo.GetByCompanyID(companyID, filters)
}

func (s *placementService) GetByID(id uint) (*models.Placement, error) {
	placement, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if placement == nil {
		return nil, fmt.Errorf("placement not found")
	}
	return placement, nil
}

// Create valida la integridad antes de colocar:
//  1. el cliente final pertenece al mismo tenant (companyID),
//  2. la application existe y pertenece al mismo tenant,
//  3. la application está en etapa 'hired'.
//
// CandidateID y JobID se copian de la application (no se confía en el body).
func (s *placementService) Create(companyID uint, dto dtos.CreatePlacementDTO) (*models.Placement, error) {
	// 1. Cliente final del mismo tenant
	client, err := s.staffingRepo.GetByID(dto.StaffingClientID)
	if err != nil {
		return nil, err
	}
	if client == nil {
		return nil, fmt.Errorf("staffing client not found")
	}
	if client.CompanyID != companyID {
		return nil, fmt.Errorf("staffing client does not belong to your company")
	}

	// 2. Application del mismo tenant
	application, err := s.applicationRepo.GetByID(dto.ApplicationID)
	if err != nil {
		return nil, err
	}
	if application == nil {
		return nil, fmt.Errorf("application not found")
	}
	if application.CompanyID != companyID {
		return nil, fmt.Errorf("application does not belong to your company")
	}

	// 3. Solo se coloca a candidatos contratados
	if application.Stage != "hired" {
		return nil, fmt.Errorf("application must be in 'hired' stage to create a placement")
	}

	status := dto.Status
	if status == "" {
		status = "active"
	}

	jobID := application.JobID
	appID := application.ID

	placement := &models.Placement{
		CompanyID:        companyID,
		StaffingClientID: dto.StaffingClientID,
		CandidateID:      application.CandidateID,
		JobID:            &jobID,
		ApplicationID:    &appID,
		StartDate:        dto.StartDate,
		EndDate:          dto.EndDate,
		ContractType:     dto.ContractType,
		Position:         dto.Position,
		BillRateAmount:   dto.BillRateAmount,
		BillRateCurrency: dto.BillRateCurrency,
		BillRateType:     dto.BillRateType,
		PayRateAmount:    dto.PayRateAmount,
		Status:           status,
		Notes:            dto.Notes,
	}
	return s.repo.Create(placement)
}

func (s *placementService) Update(id uint, dto dtos.UpdatePlacementDTO) (*models.Placement, error) {
	placement, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if placement == nil {
		return nil, fmt.Errorf("placement not found")
	}

	if dto.StartDate != nil {
		placement.StartDate = dto.StartDate
	}
	if dto.EndDate != nil {
		placement.EndDate = dto.EndDate
	}
	if dto.ContractType != nil {
		placement.ContractType = *dto.ContractType
	}
	if dto.Position != nil {
		placement.Position = *dto.Position
	}
	if dto.BillRateAmount != nil {
		placement.BillRateAmount = dto.BillRateAmount
	}
	if dto.BillRateCurrency != nil {
		placement.BillRateCurrency = *dto.BillRateCurrency
	}
	if dto.BillRateType != nil {
		placement.BillRateType = *dto.BillRateType
	}
	if dto.PayRateAmount != nil {
		placement.PayRateAmount = dto.PayRateAmount
	}
	if dto.Status != nil {
		placement.Status = *dto.Status
	}
	if dto.Notes != nil {
		placement.Notes = *dto.Notes
	}

	return s.repo.Update(placement)
}

func (s *placementService) Delete(id uint) error {
	placement, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if placement == nil {
		return fmt.Errorf("placement not found")
	}
	return s.repo.Delete(id)
}
