package service

import (
	"time"

	"dvra-api/internal/shared/apperr"
	"dvra-api/internal/shared/models"
	"dvra-api/internal/shared/tenant"
)

// ApplicationService define el contrato del servicio de applications
type ApplicationService interface {
	GetAllApplications() ([]models.Application, error)
	GetApplicationByID(id uint, scope tenant.Scope) (*models.Application, error)
	GetApplicationsByJobID(jobID uint) ([]models.Application, error)
	GetApplicationsByCandidateID(candidateID uint) ([]models.Application, error)
	GetApplicationsByCompanyID(companyID uint) ([]models.Application, error)
	GetApplicationsWithScope(scope tenant.Scope) ([]models.Application, error)
	GetApplicationsByStage(stage string, companyID uint) ([]models.Application, error)
	GetApplicationsGroupedByStage(scope tenant.Scope) (map[string][]models.Application, error)
	CreateApplication(dto CreateApplicationDTO) (*models.Application, error)
	UpdateApplication(id uint, scope tenant.Scope, dto UpdateApplicationDTO) (*models.Application, error)
	MoveToStage(id uint, scope tenant.Scope, stage string) (*models.Application, error)
	RateApplication(id uint, scope tenant.Scope, rating int) (*models.Application, error)
	DeleteApplication(id uint, scope tenant.Scope) error
}

type applicationService struct {
	applicationRepo ApplicationRepository
	jobRepo         JobRepository
	candidateRepo   CandidateRepository
}

func NewApplicationService(applicationRepo ApplicationRepository, jobRepo JobRepository, candidateRepo CandidateRepository) ApplicationService {
	return &applicationService{applicationRepo: applicationRepo, jobRepo: jobRepo, candidateRepo: candidateRepo}
}

func (s *applicationService) GetAllApplications() ([]models.Application, error) {
	return s.applicationRepo.GetAll()
}

func (s *applicationService) GetApplicationByID(id uint, scope tenant.Scope) (*models.Application, error) {
	application, err := s.applicationRepo.GetByID(id, scope)
	if err != nil {
		return nil, err
	}
	if application == nil {
		return nil, apperr.NotFound("application not found")
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

func (s *applicationService) GetApplicationsWithScope(scope tenant.Scope) ([]models.Application, error) {
	return s.applicationRepo.GetByScope(scope)
}

func (s *applicationService) GetApplicationsByStage(stage string, companyID uint) ([]models.Application, error) {
	return s.applicationRepo.GetByStage(stage, companyID)
}

func (s *applicationService) CreateApplication(dto CreateApplicationDTO) (*models.Application, error) {
	job, err := s.jobRepo.GetByID(dto.JobID, tenant.ForCompany(dto.CompanyID))
	if err != nil {
		return nil, err
	}
	if job == nil {
		return nil, apperr.NotFound("job not found")
	}
	if job.CompanyID != dto.CompanyID {
		return nil, apperr.Forbidden("job does not belong to your company")
	}
	candidate, err := s.candidateRepo.GetByID(dto.CandidateID, tenant.ForCompany(dto.CompanyID))
	if err != nil {
		return nil, err
	}
	if candidate == nil {
		return nil, apperr.NotFound("candidate not found")
	}
	if candidate.CompanyID != dto.CompanyID {
		return nil, apperr.Forbidden("candidate does not belong to your company")
	}
	existing, err := s.applicationRepo.GetByCandidateAndJob(dto.CandidateID, dto.JobID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, apperr.Conflict("candidate already has an application for this job")
	}
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

func (s *applicationService) UpdateApplication(id uint, scope tenant.Scope, dto UpdateApplicationDTO) (*models.Application, error) {
	application, err := s.applicationRepo.GetByID(id, scope)
	if err != nil {
		return nil, err
	}
	if application == nil {
		return nil, apperr.NotFound("application not found")
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

func (s *applicationService) DeleteApplication(id uint, scope tenant.Scope) error {
	application, err := s.applicationRepo.GetByID(id, scope)
	if err != nil {
		return err
	}
	if application == nil {
		return apperr.NotFound("application not found")
	}
	return s.applicationRepo.Delete(id)
}

func (s *applicationService) GetApplicationsGroupedByStage(scope tenant.Scope) (map[string][]models.Application, error) {
	applications, err := s.applicationRepo.GetByScope(scope)
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

func (s *applicationService) MoveToStage(id uint, scope tenant.Scope, stage string) (*models.Application, error) {
	application, err := s.applicationRepo.GetByID(id, scope)
	if err != nil {
		return nil, err
	}
	if application == nil {
		return nil, apperr.NotFound("application not found")
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

func (s *applicationService) RateApplication(id uint, scope tenant.Scope, rating int) (*models.Application, error) {
	application, err := s.applicationRepo.GetByID(id, scope)
	if err != nil {
		return nil, err
	}
	if application == nil {
		return nil, apperr.NotFound("application not found")
	}

	application.Rating = &rating
	return s.applicationRepo.Update(application)
}
