package services

import (
	"fmt"
	"time"

	"dvra-api/internal/app/dtos"
	"dvra-api/internal/app/models"
	"dvra-api/internal/app/repositories"
)

// PublicService define el contrato del servicio público para Career Page
type PublicService interface {
	// Company
	GetCompanyBySlug(slug string) (*models.Company, error)

	// Jobs
	GetPublishedJobsByCompanySlug(slug string) ([]models.Job, error)
	GetPublishedJobByID(jobID uint) (*models.Job, error)

	// Applications
	ApplyToJob(jobID uint, dto dtos.PublicApplicationDTO) (*models.Application, error)
}

type publicService struct {
	companyRepo     repositories.CompanyRepository
	jobRepo         repositories.JobRepository
	candidateRepo   repositories.CandidateRepository
	applicationRepo repositories.ApplicationRepository
}

// NewPublicService crea una nueva instancia de PublicService
func NewPublicService(
	companyRepo repositories.CompanyRepository,
	jobRepo repositories.JobRepository,
	candidateRepo repositories.CandidateRepository,
	applicationRepo repositories.ApplicationRepository,
) PublicService {
	return &publicService{
		companyRepo:     companyRepo,
		jobRepo:         jobRepo,
		candidateRepo:   candidateRepo,
		applicationRepo: applicationRepo,
	}
}

// GetCompanyBySlug obtiene una empresa por su slug
func (s *publicService) GetCompanyBySlug(slug string) (*models.Company, error) {
	company, err := s.companyRepo.GetBySlug(slug)
	if err != nil {
		return nil, err
	}
	if company == nil {
		return nil, fmt.Errorf("company not found")
	}
	return company, nil
}

// GetPublishedJobsByCompanySlug obtiene todos los jobs publicados de una empresa
func (s *publicService) GetPublishedJobsByCompanySlug(slug string) ([]models.Job, error) {
	// Primero obtener la empresa
	company, err := s.companyRepo.GetBySlug(slug)
	if err != nil {
		return nil, err
	}
	if company == nil {
		return nil, fmt.Errorf("company not found")
	}

	// Obtener jobs activos de la empresa
	jobs, err := s.jobRepo.GetPublishedByCompanyID(company.ID)
	if err != nil {
		return nil, err
	}

	return jobs, nil
}

// GetPublishedJobByID obtiene un job publicado por ID
func (s *publicService) GetPublishedJobByID(jobID uint) (*models.Job, error) {
	job, err := s.jobRepo.GetByID(jobID)
	if err != nil {
		return nil, err
	}
	if job == nil {
		return nil, fmt.Errorf("job not found")
	}

	// Verificar que el job esté publicado
	if job.Status != "published" {
		return nil, fmt.Errorf("job is not available")
	}

	return job, nil
}

// ApplyToJob permite a un candidato aplicar a un job
func (s *publicService) ApplyToJob(jobID uint, dto dtos.PublicApplicationDTO) (*models.Application, error) {
	// 1. Obtener el job y verificar que esté activo
	job, err := s.jobRepo.GetByID(jobID)
	if err != nil {
		return nil, err
	}
	if job == nil {
		return nil, fmt.Errorf("job not found")
	}
	if job.Status != "published" {
		return nil, fmt.Errorf("job is not accepting applications")
	}

	// 2. Buscar o crear candidato
	candidate, err := s.candidateRepo.GetByEmail(dto.Email, job.CompanyID)
	if err != nil {
		return nil, err
	}

	if candidate == nil {
		// Crear nuevo candidato
		candidate = &models.Candidate{
			CompanyID:   job.CompanyID,
			Email:       dto.Email,
			FirstName:   dto.FirstName,
			LastName:    dto.LastName,
			Phone:       dto.Phone,
			LinkedinURL: dto.LinkedinURL,
			GithubURL:   dto.GithubURL,
			ResumeURL:   dto.ResumeURL,
			Source:      "career_page",
		}
		candidate, err = s.candidateRepo.Create(candidate)
		if err != nil {
			return nil, fmt.Errorf("failed to create candidate: %w", err)
		}
	} else {
		// Actualizar datos del candidato si están vacíos
		updated := false
		if candidate.Phone == "" && dto.Phone != "" {
			candidate.Phone = dto.Phone
			updated = true
		}
		if candidate.LinkedinURL == "" && dto.LinkedinURL != "" {
			candidate.LinkedinURL = dto.LinkedinURL
			updated = true
		}
		if candidate.GithubURL == "" && dto.GithubURL != "" {
			candidate.GithubURL = dto.GithubURL
			updated = true
		}
		if dto.ResumeURL != "" {
			candidate.ResumeURL = dto.ResumeURL
			updated = true
		}
		if updated {
			candidate, err = s.candidateRepo.Update(candidate)
			if err != nil {
				return nil, fmt.Errorf("failed to update candidate: %w", err)
			}
		}
	}

	// 3. Verificar si ya existe una aplicación para este job
	existingApp, err := s.applicationRepo.GetByCandidateAndJob(candidate.ID, jobID)
	if err != nil {
		return nil, err
	}
	if existingApp != nil {
		return nil, fmt.Errorf("you have already applied to this job")
	}

	// 4. Crear la aplicación
	application := &models.Application{
		JobID:       jobID,
		CandidateID: candidate.ID,
		CompanyID:   job.CompanyID,
		Stage:       "applied",
		Notes:       dto.CoverLetter,
		AppliedAt:   time.Now(),
	}

	application, err = s.applicationRepo.Create(application)
	if err != nil {
		return nil, fmt.Errorf("failed to create application: %w", err)
	}

	// Cargar relaciones para la respuesta
	application.Job = job
	application.Candidate = candidate

	return application, nil
}
