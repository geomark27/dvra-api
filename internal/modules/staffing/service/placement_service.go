package service

import (
	"dvra-api/internal/app/dtos"
	"dvra-api/internal/app/models"
	"dvra-api/internal/modules/staffing/domain"
	"dvra-api/internal/shared/apperr"
)

// PlacementService orquesta la colocación de candidatos en clientes finales.
// Depende de tres puertos: el repo de colocaciones, el de clientes finales, y
// ApplicationFinder (puerto cross-módulo hacia recruitment).
type PlacementService struct {
	repo    domain.PlacementRepository
	clients domain.StaffingClientRepository
	apps    domain.ApplicationFinder
}

func NewPlacementService(
	repo domain.PlacementRepository,
	clients domain.StaffingClientRepository,
	apps domain.ApplicationFinder,
) *PlacementService {
	return &PlacementService{repo: repo, clients: clients, apps: apps}
}

func (s *PlacementService) GetByCompanyID(companyID uint, filters dtos.PlacementFilters) ([]models.Placement, error) {
	return s.repo.GetByCompanyID(companyID, filters)
}

func (s *PlacementService) GetByID(id uint) (*models.Placement, error) {
	placement, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if placement == nil {
		return nil, apperr.NotFound("placement not found")
	}
	return placement, nil
}

// Create valida la integridad antes de colocar:
//  1. el cliente final pertenece al mismo tenant (companyID),
//  2. la application existe y pertenece al mismo tenant,
//  3. la application está en etapa 'hired',
//  4. no existe ya un placement para esa application.
//
// CandidateID y JobID se copian de la application (no se confía en el body).
func (s *PlacementService) Create(companyID uint, dto dtos.CreatePlacementDTO) (*models.Placement, error) {
	client, err := s.clients.GetByID(dto.StaffingClientID)
	if err != nil {
		return nil, err
	}
	if client == nil {
		return nil, apperr.NotFound("staffing client not found")
	}
	if client.CompanyID != companyID {
		return nil, apperr.Forbidden("staffing client does not belong to your company")
	}

	app, err := s.apps.FindByID(dto.ApplicationID)
	if err != nil {
		return nil, err
	}
	if app == nil {
		return nil, apperr.NotFound("application not found")
	}
	if app.CompanyID != companyID {
		return nil, apperr.Forbidden("application does not belong to your company")
	}
	if app.Stage != "hired" {
		return nil, apperr.BadRequest("application must be in 'hired' stage to create a placement")
	}

	exists, err := s.repo.ExistsByApplicationID(dto.ApplicationID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, apperr.Conflict("a placement already exists for this application")
	}

	status := dto.Status
	if status == "" {
		status = "active"
	}

	placement := &models.Placement{
		CompanyID:        companyID,
		StaffingClientID: dto.StaffingClientID,
		CandidateID:      app.CandidateID,
		JobID:            app.JobID,
		ApplicationID:    app.ID,
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

// companyID = 0 omite la validación de tenant (SuperAdmin).
func (s *PlacementService) Update(id, companyID uint, dto dtos.UpdatePlacementDTO) (*models.Placement, error) {
	placement, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if placement == nil {
		return nil, apperr.NotFound("placement not found")
	}
	if companyID > 0 && placement.CompanyID != companyID {
		return nil, apperr.Forbidden("access denied")
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

func (s *PlacementService) Delete(id, companyID uint) error {
	placement, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if placement == nil {
		return apperr.NotFound("placement not found")
	}
	if companyID > 0 && placement.CompanyID != companyID {
		return apperr.Forbidden("access denied")
	}
	return s.repo.Delete(id)
}
