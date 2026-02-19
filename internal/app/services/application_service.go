package services

import (
	"fmt"
	"time"

	"dvra-api/internal/app/dtos"
	"dvra-api/internal/app/models"
	"dvra-api/internal/app/repositories"
)

// ApplicationService define el contrato del servicio de applications
type ApplicationService interface {
	GetAllApplications() ([]models.Application, error)
	GetApplicationByID(id uint) (*models.Application, error)
	GetApplicationsByJobID(jobID uint) ([]models.Application, error)
	GetApplicationsByCandidateID(candidateID uint) ([]models.Application, error)
	GetApplicationsByCompanyID(companyID uint) ([]models.Application, error)
	GetApplicationsByStage(stage string, companyID uint) ([]models.Application, error)
	GetApplicationsGroupedByStage(companyID uint) (map[string][]models.Application, error)
	CreateApplication(dto dtos.CreateApplicationDTO) (*models.Application, error)
	UpdateApplication(id uint, dto dtos.UpdateApplicationDTO) (*models.Application, error)
	MoveToStage(id uint, stage string) (*models.Application, error)
	RateApplication(id uint, rating int) (*models.Application, error)
	DeleteApplication(id uint) error
}

type applicationService struct {
	applicationRepo repositories.ApplicationRepository
}

func NewApplicationService(applicationRepo repositories.ApplicationRepository) ApplicationService {
	return &applicationService{applicationRepo: applicationRepo}
}

func (s *applicationService) GetAllApplications() ([]models.Application, error) {
	return s.applicationRepo.GetAll()
}

func (s *applicationService) GetApplicationByID(id uint) (*models.Application, error) {
	application, err := s.applicationRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if application == nil {
		return nil, fmt.Errorf("application not found")
	}
	return application, nil
}

func (s *applicationService) GetApplicationsByJobID(jobID uint) ([]models.Application, error) {
	return s.applicationRepo.GetByJobID(jobID)
}

func (s *applicationService) GetApplicationsByCandidateID(candidateID uint) ([]models.Application, error) {
	return s.applicationRepo.GetByCandidateID(candidateID)
}

func (s *applicationService) GetApplicationsByCompanyID(companyID uint) ([]models.Application, error) {
	return s.applicationRepo.GetByCompanyID(companyID)
}

func (s *applicationService) GetApplicationsByStage(stage string, companyID uint) ([]models.Application, error) {
	return s.applicationRepo.GetByStage(stage, companyID)
}

func (s *applicationService) CreateApplication(dto dtos.CreateApplicationDTO) (*models.Application, error) {
	now := time.Now()
	application := &models.Application{
		JobID:       dto.JobID,
		CandidateID: dto.CandidateID,
		CompanyID:   dto.CompanyID,
		Stage:       dto.Stage,
		Rating:      dto.Rating,
		Notes:       dto.Notes,
		AppliedAt:   now,
	}

	return s.applicationRepo.Create(application)
}

func (s *applicationService) UpdateApplication(id uint, dto dtos.UpdateApplicationDTO) (*models.Application, error) {
	application, err := s.applicationRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if application == nil {
		return nil, fmt.Errorf("application not found")
	}

	if dto.Stage != nil {
		application.Stage = *dto.Stage

		// Auto-actualizar timestamps según el stage
		now := time.Now()
		if *dto.Stage == "rejected" && application.RejectedAt == nil {
			application.RejectedAt = &now
		}
		if *dto.Stage == "hired" && application.HiredAt == nil {
			application.HiredAt = &now
		}
	}
	if dto.Rating != nil {
		application.Rating = dto.Rating
	}
	if dto.Notes != nil {
		application.Notes = *dto.Notes
	}
	if dto.RejectedAt != nil {
		application.RejectedAt = dto.RejectedAt
	}
	if dto.HiredAt != nil {
		application.HiredAt = dto.HiredAt
	}

	return s.applicationRepo.Update(application)
}

func (s *applicationService) DeleteApplication(id uint) error {
	application, err := s.applicationRepo.GetByID(id)
	if err != nil {
		return err
	}
	if application == nil {
		return fmt.Errorf("application not found")
	}
	return s.applicationRepo.Delete(id)
}

func (s *applicationService) GetApplicationsGroupedByStage(companyID uint) (map[string][]models.Application, error) {
	applications, err := s.applicationRepo.GetByCompanyID(companyID)
	if err != nil {
		return nil, err
	}

	// Initialize all stages
	stages := []string{"applied", "screening", "technical", "offer", "hired"}
	result := make(map[string][]models.Application)
	for _, stage := range stages {
		result[stage] = []models.Application{}
	}

	// Group applications by stage
	for _, app := range applications {
		result[app.Stage] = append(result[app.Stage], app)
	}

	return result, nil
}

func (s *applicationService) MoveToStage(id uint, stage string) (*models.Application, error) {
	application, err := s.applicationRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if application == nil {
		return nil, fmt.Errorf("application not found")
	}

	application.Stage = stage

	// Auto-update timestamps based on stage
	now := time.Now()
	if stage == "rejected" && application.RejectedAt == nil {
		application.RejectedAt = &now
	}
	if stage == "hired" && application.HiredAt == nil {
		application.HiredAt = &now
	}

	return s.applicationRepo.Update(application)
}

func (s *applicationService) RateApplication(id uint, rating int) (*models.Application, error) {
	application, err := s.applicationRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if application == nil {
		return nil, fmt.Errorf("application not found")
	}

	application.Rating = &rating
	return s.applicationRepo.Update(application)
}
