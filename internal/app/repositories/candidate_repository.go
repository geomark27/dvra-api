package repositories

import (
	"dvra-api/internal/app/models"
	"dvra-api/internal/shared/permissions"
	"dvra-api/internal/shared/tenant"

	"gorm.io/gorm"
)

// CandidateRepository define el contrato del repositorio de candidates
type CandidateRepository interface {
	GetAll() ([]models.Candidate, error)
	GetByID(id uint, scope tenant.Scope) (*models.Candidate, error)
	GetByCompanyID(companyID uint) ([]models.Candidate, error)
	GetByScope(scope tenant.Scope) ([]models.Candidate, error)
	GetByEmail(email string, companyID uint) (*models.Candidate, error)
	Create(candidate *models.Candidate) (*models.Candidate, error)
	Update(candidate *models.Candidate) (*models.Candidate, error)
	Delete(id uint) error
}

// candidateRepository es la implementación con GORM
type candidateRepository struct{ db *gorm.DB }

// NewCandidateRepository crea una nueva instancia de CandidateRepository
func NewCandidateRepository(db *gorm.DB) CandidateRepository {
	return &candidateRepository{db: db}
}

func (r *candidateRepository) GetAll() ([]models.Candidate, error) {
	var candidates []models.Candidate
	if err := r.db.Preload("Company").Find(&candidates).Error; err != nil {
		return nil, err
	}
	return candidates, nil
}

func (r *candidateRepository) GetByID(id uint, scope tenant.Scope) (*models.Candidate, error) {
	var candidate models.Candidate
	query := tenant.Apply(r.db, scope, "candidates.company_id")
	query = applyCandidateActorScope(query, scope)
	if err := query.Preload("Company").Preload("Applications").Where("candidates.id = ?", id).First(&candidate).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &candidate, nil
}

func (r *candidateRepository) GetByScope(scope tenant.Scope) ([]models.Candidate, error) {
	var candidates []models.Candidate
	query := tenant.Apply(r.db, scope, "candidates.company_id")
	query = applyCandidateActorScope(query, scope)
	if err := query.Find(&candidates).Error; err != nil {
		return nil, err
	}
	return candidates, nil
}

func (r *candidateRepository) GetByCompanyID(companyID uint) ([]models.Candidate, error) {
	var candidates []models.Candidate
	if err := r.db.Where("company_id = ?", companyID).Find(&candidates).Error; err != nil {
		return nil, err
	}
	return candidates, nil
}

func (r *candidateRepository) GetByEmail(email string, companyID uint) (*models.Candidate, error) {
	var candidate models.Candidate
	if err := r.db.Where("email = ? AND company_id = ?", email, companyID).First(&candidate).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &candidate, nil
}

func (r *candidateRepository) Create(candidate *models.Candidate) (*models.Candidate, error) {
	if err := r.db.Create(candidate).Error; err != nil {
		return nil, err
	}
	return candidate, nil
}

func (r *candidateRepository) Update(candidate *models.Candidate) (*models.Candidate, error) {
	if err := r.db.Save(candidate).Error; err != nil {
		return nil, err
	}
	return candidate, nil
}

func (r *candidateRepository) Delete(id uint) error {
	return r.db.Delete(&models.Candidate{}, id).Error
}

func applyCandidateActorScope(query *gorm.DB, scope tenant.Scope) *gorm.DB {
	if scope.Role != permissions.RoleHiringManager {
		return query
	}
	return query.
		Joins("JOIN applications AS visibility_applications ON visibility_applications.candidate_id = candidates.id AND visibility_applications.deleted_at IS NULL").
		Joins("JOIN jobs AS visibility_jobs ON visibility_jobs.id = visibility_applications.job_id AND visibility_jobs.deleted_at IS NULL").
		Where("visibility_jobs.hiring_manager = ?", scope.UserID).
		Distinct("candidates.*")
}
