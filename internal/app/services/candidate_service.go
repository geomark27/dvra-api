package services

import (
	"fmt"

	"dvra-api/internal/app/dtos"
	"dvra-api/internal/app/models"
	"dvra-api/internal/app/repositories"
)

// CandidateService define el contrato del servicio de candidates
type CandidateService interface {
	GetAllCandidates() ([]models.Candidate, error)
	GetCandidateByID(id uint) (*models.Candidate, error)
	GetCandidatesByCompanyID(companyID uint) ([]models.Candidate, error)
	CreateCandidate(dto dtos.CreateCandidateDTO) (*models.Candidate, error)
	UpdateCandidate(id uint, dto dtos.UpdateCandidateDTO) (*models.Candidate, error)
	DeleteCandidate(id uint) error
}

type candidateService struct {
	candidateRepo repositories.CandidateRepository
}

func NewCandidateService(candidateRepo repositories.CandidateRepository) CandidateService {
	return &candidateService{candidateRepo: candidateRepo}
}

func (s *candidateService) GetAllCandidates() ([]models.Candidate, error) {
	return s.candidateRepo.GetAll()
}

func (s *candidateService) GetCandidateByID(id uint) (*models.Candidate, error) {
	candidate, err := s.candidateRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if candidate == nil {
		return nil, fmt.Errorf("candidate not found")
	}
	return candidate, nil
}

func (s *candidateService) GetCandidatesByCompanyID(companyID uint) ([]models.Candidate, error) {
	return s.candidateRepo.GetByCompanyID(companyID)
}

func (s *candidateService) CreateCandidate(dto dtos.CreateCandidateDTO) (*models.Candidate, error) {
	// Verificar duplicado por email en la misma company
	existing, err := s.candidateRepo.GetByEmail(dto.Email, dto.CompanyID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, fmt.Errorf("candidate with email '%s' already exists in this company", dto.Email)
	}

	candidate := &models.Candidate{
		CompanyID:   dto.CompanyID,
		Email:       dto.Email,
		FirstName:   dto.FirstName,
		LastName:    dto.LastName,
		Phone:       dto.Phone,
		ResumeURL:   dto.ResumeURL,
		GithubURL:   dto.GithubURL,
		LinkedinURL: dto.LinkedinURL,
		Source:      dto.Source,
	}

	return s.candidateRepo.Create(candidate)
}

func (s *candidateService) UpdateCandidate(id uint, dto dtos.UpdateCandidateDTO) (*models.Candidate, error) {
	candidate, err := s.candidateRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if candidate == nil {
		return nil, fmt.Errorf("candidate not found")
	}

	if dto.Email != nil {
		candidate.Email = *dto.Email
	}
	if dto.FirstName != nil {
		candidate.FirstName = *dto.FirstName
	}
	if dto.LastName != nil {
		candidate.LastName = *dto.LastName
	}
	if dto.Phone != nil {
		candidate.Phone = *dto.Phone
	}
	if dto.ResumeURL != nil {
		candidate.ResumeURL = *dto.ResumeURL
	}
	if dto.GithubURL != nil {
		candidate.GithubURL = *dto.GithubURL
	}
	if dto.LinkedinURL != nil {
		candidate.LinkedinURL = *dto.LinkedinURL
	}
	if dto.Source != nil {
		candidate.Source = *dto.Source
	}

	return s.candidateRepo.Update(candidate)
}

func (s *candidateService) DeleteCandidate(id uint) error {
	candidate, err := s.candidateRepo.GetByID(id)
	if err != nil {
		return err
	}
	if candidate == nil {
		return fmt.Errorf("candidate not found")
	}
	return s.candidateRepo.Delete(id)
}
