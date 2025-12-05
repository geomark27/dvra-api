package services

import (
	"fmt"

	"dvra-api/internal/app/dtos"
	"dvra-api/internal/app/models"
	"dvra-api/internal/app/repositories"
)

// JobService define el contrato del servicio de jobs
type JobService interface {
	GetAllJobs() ([]models.Job, error)
	GetJobByID(id uint) (*models.Job, error)
	GetJobsByCompanyID(companyID uint) ([]models.Job, error)
	GetJobsByStatus(status string, companyID uint) ([]models.Job, error)
	CreateJob(dto dtos.CreateJobDTO) (*models.Job, error)
	UpdateJob(id uint, dto dtos.UpdateJobDTO) (*models.Job, error)
	DeleteJob(id uint) error
}

type jobService struct {
	jobRepo repositories.JobRepository
}

func NewJobService(jobRepo repositories.JobRepository) JobService {
	return &jobService{jobRepo: jobRepo}
}

func (s *jobService) GetAllJobs() ([]models.Job, error) {
	return s.jobRepo.GetAll()
}

func (s *jobService) GetJobByID(id uint) (*models.Job, error) {
	job, err := s.jobRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if job == nil {
		return nil, fmt.Errorf("job not found")
	}
	return job, nil
}

func (s *jobService) GetJobsByCompanyID(companyID uint) ([]models.Job, error) {
	return s.jobRepo.GetByCompanyID(companyID)
}

func (s *jobService) GetJobsByStatus(status string, companyID uint) ([]models.Job, error) {
	return s.jobRepo.GetByStatus(status, companyID)
}

func (s *jobService) CreateJob(dto dtos.CreateJobDTO) (*models.Job, error) {
	status := dto.Status
	if status == "" {
		status = "draft"
	}

	job := &models.Job{
		CompanyID:         dto.CompanyID,
		Title:             dto.Title,
		Description:       dto.Description,
		Location:          dto.Location,
		Status:            status,
		AssignedRecruiter: dto.AssignedRecruiter,
		HiringManager:     dto.HiringManager,
	}

	return s.jobRepo.Create(job)
}

func (s *jobService) UpdateJob(id uint, dto dtos.UpdateJobDTO) (*models.Job, error) {
	job, err := s.jobRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if job == nil {
		return nil, fmt.Errorf("job not found")
	}

	if dto.Title != nil {
		job.Title = *dto.Title
	}
	if dto.Description != nil {
		job.Description = *dto.Description
	}
	if dto.Location != nil {
		job.Location = *dto.Location
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

	return s.jobRepo.Update(job)
}

func (s *jobService) DeleteJob(id uint) error {
	job, err := s.jobRepo.GetByID(id)
	if err != nil {
		return err
	}
	if job == nil {
		return fmt.Errorf("job not found")
	}
	return s.jobRepo.Delete(id)
}
