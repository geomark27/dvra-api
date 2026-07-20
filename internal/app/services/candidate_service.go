package services

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"dvra-api/internal/app/dtos"
	"dvra-api/internal/app/models"
	"dvra-api/internal/app/repositories"
	"dvra-api/internal/shared/apperr"
	billingctx "dvra-api/internal/shared/billing"
	"dvra-api/internal/shared/tenant"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// CandidateService define el contrato del servicio de candidates
type CandidateService interface {
	GetAllCandidates() ([]models.Candidate, error)
	GetCandidateByID(id uint, scope tenant.Scope) (*models.Candidate, error)
	GetCandidatesByCompanyID(companyID uint) ([]models.Candidate, error)
	GetCandidatesWithScope(scope tenant.Scope) ([]models.Candidate, error)
	CreateCandidate(dto dtos.CreateCandidateDTO) (*models.Candidate, error)
	UpdateCandidate(id uint, scope tenant.Scope, dto dtos.UpdateCandidateDTO) (*models.Candidate, error)
	DeleteCandidate(id uint, scope tenant.Scope) error
	StoreResume(id uint, scope tenant.Scope, upload ResumeUpload) (*models.StoredFile, error)
	OpenResume(id uint, scope tenant.Scope) (*models.StoredFile, io.ReadCloser, error)
}

type ResumeUpload struct {
	Reader       io.Reader
	OriginalName string
	Size         int64
}

type FileStorage interface {
	Save(key string, src io.Reader) (int64, error)
	Open(key string) (io.ReadCloser, error)
	Delete(key string) error
}

type candidateService struct {
	candidateRepo repositories.CandidateRepository
	planRepo      repositories.PlanRepository
	db            *gorm.DB
	storage       FileStorage
}

func NewCandidateService(candidateRepo repositories.CandidateRepository, planRepo repositories.PlanRepository, db *gorm.DB, storage FileStorage) CandidateService {
	return &candidateService{candidateRepo: candidateRepo, planRepo: planRepo, db: db, storage: storage}
}

func (s *candidateService) GetAllCandidates() ([]models.Candidate, error) {
	return s.candidateRepo.GetAll()
}

func (s *candidateService) GetCandidateByID(id uint, scope tenant.Scope) (*models.Candidate, error) {
	candidate, err := s.candidateRepo.GetByID(id, scope)
	if err != nil {
		return nil, err
	}
	if candidate == nil {
		return nil, apperr.NotFound("candidate not found")
	}
	return candidate, nil
}

func (s *candidateService) GetCandidatesByCompanyID(companyID uint) ([]models.Candidate, error) {
	return s.candidateRepo.GetByCompanyID(companyID)
}

func (s *candidateService) GetCandidatesWithScope(scope tenant.Scope) ([]models.Candidate, error) {
	return s.candidateRepo.GetByScope(scope)
}

func (s *candidateService) CreateCandidate(dto dtos.CreateCandidateDTO) (*models.Candidate, error) {
	// Verificar duplicado por email en la misma company
	existing, err := s.candidateRepo.GetByEmail(dto.Email, dto.CompanyID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, apperr.Conflict(fmt.Sprintf("candidate with email '%s' already exists in this company", dto.Email))
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

func (s *candidateService) UpdateCandidate(id uint, scope tenant.Scope, dto dtos.UpdateCandidateDTO) (*models.Candidate, error) {
	candidate, err := s.candidateRepo.GetByID(id, scope)
	if err != nil {
		return nil, err
	}
	if candidate == nil {
		return nil, apperr.NotFound("candidate not found")
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

func (s *candidateService) DeleteCandidate(id uint, scope tenant.Scope) error {
	candidate, err := s.candidateRepo.GetByID(id, scope)
	if err != nil {
		return err
	}
	if candidate == nil {
		return apperr.NotFound("candidate not found")
	}
	var stored models.StoredFile
	hasStoredFile := false
	err = s.db.Transaction(func(tx *gorm.DB) error {
		queryErr := tx.Where("company_id = ? AND candidate_id = ?", candidate.CompanyID, candidate.ID).First(&stored).Error
		if queryErr != nil && queryErr != gorm.ErrRecordNotFound {
			return queryErr
		}
		hasStoredFile = queryErr == nil
		if err := tx.Delete(&models.Candidate{}, candidate.ID).Error; err != nil {
			return err
		}
		if hasStoredFile {
			return tx.Delete(&stored).Error
		}
		return nil
	})
	if err != nil {
		return err
	}
	if hasStoredFile {
		if err := s.storage.Delete(stored.Key); err != nil {
			return fmt.Errorf("candidate deleted but resume cleanup failed: %w", err)
		}
	}
	return nil
}

func (s *candidateService) StoreResume(id uint, scope tenant.Scope, upload ResumeUpload) (*models.StoredFile, error) {
	candidate, err := s.candidateRepo.GetByID(id, scope)
	if err != nil {
		return nil, err
	}
	if candidate == nil {
		return nil, apperr.NotFound("candidate not found")
	}
	const maxResumeSize = int64(10 * 1024 * 1024)
	if upload.Size <= 0 || upload.Size > maxResumeSize {
		return nil, apperr.BadRequest("resume file must be between 1 byte and 10MB")
	}
	ext := strings.ToLower(filepath.Ext(upload.OriginalName))
	if ext != ".pdf" && ext != ".doc" && ext != ".docx" {
		return nil, apperr.BadRequest("resume must be PDF, DOC, or DOCX format")
	}
	header := make([]byte, 8)
	n, readErr := io.ReadFull(upload.Reader, header)
	if readErr != nil && readErr != io.ErrUnexpectedEOF {
		return nil, apperr.BadRequest("invalid resume file")
	}
	header = header[:n]
	if !validResumeSignature(ext, header) {
		return nil, apperr.BadRequest("resume content does not match its extension")
	}
	random := make([]byte, 16)
	if _, err := rand.Read(random); err != nil {
		return nil, err
	}
	key := fmt.Sprintf("companies/%d/candidates/%d/%s%s", candidate.CompanyID, candidate.ID, hex.EncodeToString(random), ext)
	reader := io.MultiReader(strings.NewReader(string(header)), upload.Reader)
	contentType := resumeContentType(ext)

	var stored models.StoredFile
	var oldKey string
	saved := false
	err = s.db.Transaction(func(tx *gorm.DB) error {
		var company models.Company
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&company, candidate.CompanyID).Error; err != nil {
			return err
		}
		plan, err := billingctx.CurrentPlan(tx, company.ID)
		if err != nil {
			return err
		}
		if plan == nil {
			return ErrPlanNotFound
		}
		var used int64
		if err := tx.Model(&models.StoredFile{}).Where("company_id = ?", company.ID).Select("COALESCE(SUM(size_bytes), 0)").Scan(&used).Error; err != nil {
			return err
		}
		var existing models.StoredFile
		existingErr := tx.Where("company_id = ? AND candidate_id = ?", company.ID, candidate.ID).First(&existing).Error
		if existingErr != nil && existingErr != gorm.ErrRecordNotFound {
			return existingErr
		}
		if existingErr == nil {
			used -= existing.SizeBytes
			oldKey = existing.Key
		}
		if plan.MaxStorageGB >= 0 && used+upload.Size > int64(plan.MaxStorageGB)*1024*1024*1024 {
			return apperr.Conflict("storage limit reached: upgrade your plan or remove files")
		}
		written, err := s.storage.Save(key, io.LimitReader(reader, maxResumeSize+1))
		if err != nil {
			return err
		}
		saved = true
		if written > maxResumeSize || (plan.MaxStorageGB >= 0 && used+written > int64(plan.MaxStorageGB)*1024*1024*1024) {
			return apperr.Conflict("storage limit reached: upgrade your plan or reduce the file size")
		}
		stored = existing
		stored.CompanyID = company.ID
		stored.CandidateID = candidate.ID
		stored.Key = key
		stored.OriginalName = filepath.Base(upload.OriginalName)
		stored.ContentType = contentType
		stored.SizeBytes = written
		if existingErr == gorm.ErrRecordNotFound {
			if err := tx.Create(&stored).Error; err != nil {
				return err
			}
		} else if err := tx.Save(&stored).Error; err != nil {
			return err
		}
		resumeURL := fmt.Sprintf("/api/v1/candidates/%d/resume", candidate.ID)
		return tx.Model(&models.Candidate{}).Where("id = ? AND company_id = ?", candidate.ID, company.ID).Update("resume_url", resumeURL).Error
	})
	if err != nil {
		if saved {
			_ = s.storage.Delete(key)
		}
		return nil, err
	}
	if oldKey != "" && oldKey != key {
		_ = s.storage.Delete(oldKey)
	}
	return &stored, nil
}

func (s *candidateService) OpenResume(id uint, scope tenant.Scope) (*models.StoredFile, io.ReadCloser, error) {
	candidate, err := s.candidateRepo.GetByID(id, scope)
	if err != nil || candidate == nil {
		if err == nil {
			err = apperr.NotFound("candidate not found")
		}
		return nil, nil, err
	}
	var stored models.StoredFile
	if err := s.db.Where("company_id = ? AND candidate_id = ?", candidate.CompanyID, candidate.ID).First(&stored).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil, apperr.NotFound("resume not found")
		}
		return nil, nil, err
	}
	reader, err := s.storage.Open(stored.Key)
	if err != nil {
		return nil, nil, err
	}
	return &stored, reader, nil
}

func validResumeSignature(ext string, header []byte) bool {
	switch ext {
	case ".pdf":
		return len(header) >= 4 && string(header[:4]) == "%PDF"
	case ".doc":
		return len(header) >= 8 && string(header[:8]) == "\xd0\xcf\x11\xe0\xa1\xb1\x1a\xe1"
	case ".docx":
		return len(header) >= 4 && string(header[:2]) == "PK"
	default:
		return false
	}
}

func resumeContentType(ext string) string {
	if ext == ".pdf" {
		return "application/pdf"
	}
	if ext == ".doc" {
		return "application/msword"
	}
	return "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
}
