package services

import (
	"dvra-api/internal/app/dtos"
	"dvra-api/internal/app/models"
	"dvra-api/internal/app/repositories"
	"dvra-api/internal/shared/apperr"
)

// JobService define el contrato del servicio de jobs
type JobService interface {
	GetAllJobs() ([]models.Job, error)
	GetAllJobsWithFilters(filters dtos.JobFilters) ([]models.Job, error)
	GetJobByID(id uint) (*models.Job, error)
	GetJobsByCompanyID(companyID uint) ([]models.Job, error)
	GetJobsByCompanyIDWithFilters(companyID uint, filters dtos.JobFilters) ([]models.Job, error)
	GetJobsByStatus(status string, companyID uint) ([]models.Job, error)
	CreateJob(dto dtos.CreateJobDTO) (*models.Job, error)
	UpdateJob(id uint, dto dtos.UpdateJobDTO) (*models.Job, error)
	PublishJob(id uint) (*models.Job, error)
	CloseJob(id uint) (*models.Job, error)
	DeleteJob(id uint) error
}

// staffingClientReader es lo que job necesita de staffing: validar que un cliente
// final exista. Puerto definido por el consumidor — el composition root inyecta la
// implementación del módulo staffing, así recruitment no importa staffing (evita ciclo).
type staffingClientReader interface {
	GetByID(id uint) (*models.StaffingClient, error)
}

type jobService struct {
	jobRepo      repositories.JobRepository
	staffingRepo staffingClientReader
}

func NewJobService(jobRepo repositories.JobRepository, staffingRepo staffingClientReader) JobService {
	return &jobService{jobRepo: jobRepo, staffingRepo: staffingRepo}
}

// validateStaffingClient asegura que el cliente final exista y pertenezca a la
// misma empresa que el job (integridad cross-tenant). Sin esta validación un
// tenant podría enganchar jobs a clientes de otro tenant.
func (s *jobService) validateStaffingClient(companyID, staffingClientID uint) error {
	client, err := s.staffingRepo.GetByID(staffingClientID)
	if err != nil {
		return err
	}
	if client == nil {
		return apperr.NotFound("staffing client not found")
	}
	if client.CompanyID != companyID {
		return apperr.Forbidden("staffing client does not belong to your company")
	}
	return nil
}

func (s *jobService) GetAllJobs() ([]models.Job, error) {
	return s.jobRepo.GetAll()
}

func (s *jobService) GetAllJobsWithFilters(filters dtos.JobFilters) ([]models.Job, error) {
	return s.jobRepo.GetAllWithFilters(filters)
}

func (s *jobService) GetJobByID(id uint) (*models.Job, error) {
	job, err := s.jobRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if job == nil {
		return nil, apperr.NotFound("job not found")
	}
	return job, nil
}

func (s *jobService) GetJobsByCompanyID(companyID uint) ([]models.Job, error) {
	return s.jobRepo.GetByCompanyID(companyID)
}

func (s *jobService) GetJobsByCompanyIDWithFilters(companyID uint, filters dtos.JobFilters) ([]models.Job, error) {
	return s.jobRepo.GetByCompanyIDWithFilters(companyID, filters)
}

func (s *jobService) GetJobsByStatus(status string, companyID uint) ([]models.Job, error) {
	return s.jobRepo.GetByStatus(status, companyID)
}

func (s *jobService) CreateJob(dto dtos.CreateJobDTO) (*models.Job, error) {
	status := dto.Status
	if status == "" {
		status = "draft"
	}

	// Si el job se asigna a un cliente final, debe ser del mismo tenant
	if dto.StaffingClientID != nil {
		if err := s.validateStaffingClient(dto.CompanyID, *dto.StaffingClientID); err != nil {
			return nil, err
		}
	}

	job := &models.Job{
		CompanyID:         dto.CompanyID,
		Title:             dto.Title,
		Description:       dto.Description,
		LocationType:      dto.LocationType,
		CityID:            dto.CityID,
		SalaryMin:         dto.SalaryMin,
		SalaryMax:         dto.SalaryMax,
		Requirements:      dto.Requirements,
		Benefits:          dto.Benefits,
		Status:            status,
		AssignedRecruiter: dto.AssignedRecruiter,
		HiringManager:     dto.HiringManager,
		StaffingClientID:  dto.StaffingClientID,
	}

	return s.jobRepo.Create(job)
}

func (s *jobService) UpdateJob(id uint, dto dtos.UpdateJobDTO) (*models.Job, error) {
	job, err := s.jobRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if job == nil {
		return nil, apperr.NotFound("job not found")
	}

	if dto.Title != nil {
		job.Title = *dto.Title
	}
	if dto.Description != nil {
		job.Description = *dto.Description
	}
	if dto.LocationType != nil {
		job.LocationType = *dto.LocationType
	}
	if dto.CityID != nil {
		job.CityID = dto.CityID
	}
	if dto.SalaryMin != nil {
		job.SalaryMin = dto.SalaryMin
	}
	if dto.SalaryMax != nil {
		job.SalaryMax = dto.SalaryMax
	}
	if dto.Requirements != nil {
		job.Requirements = *dto.Requirements
	}
	if dto.Benefits != nil {
		job.Benefits = *dto.Benefits
	}
	if dto.Status != nil {
		job.Status = *dto.Status
	}
	if dto.AssignedRecruiter != nil {
		job.AssignedRecruiter = dto.AssignedRecruiter
	}
	if dto.HiringManager != nil {
		job.HiringManager = dto.HiringManager
	}
	if dto.StaffingClientID != nil {
		if err := s.validateStaffingClient(job.CompanyID, *dto.StaffingClientID); err != nil {
			return nil, err
		}
		job.StaffingClientID = dto.StaffingClientID
	}

	return s.jobRepo.Update(job)
}

func (s *jobService) PublishJob(id uint) (*models.Job, error) {
	job, err := s.jobRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if job == nil {
		return nil, apperr.NotFound("job not found")
	}

	// Validar que el job tenga los campos mínimos requeridos
	if job.Title == "" || job.Description == "" {
		return nil, apperr.BadRequest("job must have title and description to be published")
	}

	job.Status = "active"
	return s.jobRepo.Update(job)
}

func (s *jobService) CloseJob(id uint) (*models.Job, error) {
	job, err := s.jobRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if job == nil {
		return nil, apperr.NotFound("job not found")
	}

	job.Status = "closed"
	return s.jobRepo.Update(job)
}

func (s *jobService) DeleteJob(id uint) error {
	job, err := s.jobRepo.GetByID(id)
	if err != nil {
		return err
	}
	if job == nil {
		return apperr.NotFound("job not found")
	}
	return s.jobRepo.Delete(id)
}
