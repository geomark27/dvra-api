package repositories

import (
	"dvra-api/internal/app/models"
	"dvra-api/internal/shared/permissions"
	"dvra-api/internal/shared/tenant"

	"gorm.io/gorm"
)

// ApplicationRepository define el contrato del repositorio de applications
type ApplicationRepository interface {
	GetAll() ([]models.Application, error)
	GetByID(id uint, scope tenant.Scope) (*models.Application, error)
	GetByJobID(jobID uint) ([]models.Application, error)
	GetByCandidateID(candidateID uint) ([]models.Application, error)
	GetByCompanyID(companyID uint) ([]models.Application, error)
	GetByScope(scope tenant.Scope) ([]models.Application, error)
	GetByStage(stage string, companyID uint) ([]models.Application, error)
	GetByCandidateAndJob(candidateID, jobID uint) (*models.Application, error)
	Create(application *models.Application) (*models.Application, error)
	Update(application *models.Application) (*models.Application, error)
	Delete(id uint) error
}

// applicationRepository es la implementación con GORM
type applicationRepository struct{ db *gorm.DB }

// NewApplicationRepository crea una nueva instancia de ApplicationRepository
func NewApplicationRepository(db *gorm.DB) ApplicationRepository {
	return &applicationRepository{db: db}
}

func (r *applicationRepository) GetAll() ([]models.Application, error) {
	var applications []models.Application
	if err := r.db.Preload("Job").Preload("Candidate").Preload("Company").Find(&applications).Error; err != nil {
		return nil, err
	}
	return applications, nil
}

func (r *applicationRepository) GetByID(id uint, scope tenant.Scope) (*models.Application, error) {
	var application models.Application
	query := tenant.Apply(r.db, scope, "applications.company_id")
	query = applyApplicationActorScope(query, scope)
	if err := query.Preload("Job").Preload("Candidate").Preload("Company").Where("applications.id = ?", id).First(&application).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &application, nil
}

func (r *applicationRepository) GetByScope(scope tenant.Scope) ([]models.Application, error) {
	var applications []models.Application
	query := tenant.Apply(r.db, scope, "applications.company_id")
	query = applyApplicationActorScope(query, scope)
	if err := query.Preload("Job").Preload("Candidate").Find(&applications).Error; err != nil {
		return nil, err
	}
	return applications, nil
}

func (r *applicationRepository) GetByJobID(jobID uint) ([]models.Application, error) {
	var applications []models.Application
	if err := r.db.Where("job_id = ?", jobID).Preload("Candidate").Find(&applications).Error; err != nil {
		return nil, err
	}
	return applications, nil
}

func (r *applicationRepository) GetByCandidateID(candidateID uint) ([]models.Application, error) {
	var applications []models.Application
	if err := r.db.Where("candidate_id = ?", candidateID).Preload("Job").Find(&applications).Error; err != nil {
		return nil, err
	}
	return applications, nil
}

func (r *applicationRepository) GetByCompanyID(companyID uint) ([]models.Application, error) {
	var applications []models.Application
	if err := r.db.Where("company_id = ?", companyID).Preload("Job").Preload("Candidate").Find(&applications).Error; err != nil {
		return nil, err
	}
	return applications, nil
}

func (r *applicationRepository) GetByStage(stage string, companyID uint) ([]models.Application, error) {
	var applications []models.Application
	if err := r.db.Where("stage = ? AND company_id = ?", stage, companyID).Preload("Job").Preload("Candidate").Find(&applications).Error; err != nil {
		return nil, err
	}
	return applications, nil
}

// GetByCandidateAndJob checks if a candidate has already applied to a specific job
func (r *applicationRepository) GetByCandidateAndJob(candidateID, jobID uint) (*models.Application, error) {
	var application models.Application
	if err := r.db.Where("candidate_id = ? AND job_id = ?", candidateID, jobID).First(&application).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &application, nil
}

func (r *applicationRepository) Create(application *models.Application) (*models.Application, error) {
	if err := r.db.Create(application).Error; err != nil {
		return nil, err
	}
	return application, nil
}

func (r *applicationRepository) Update(application *models.Application) (*models.Application, error) {
	if err := r.db.Save(application).Error; err != nil {
		return nil, err
	}
	return application, nil
}

func (r *applicationRepository) Delete(id uint) error {
	return r.db.Delete(&models.Application{}, id).Error
}

func applyApplicationActorScope(query *gorm.DB, scope tenant.Scope) *gorm.DB {
	if scope.Role != permissions.RoleHiringManager {
		return query
	}
	return query.
		Joins("JOIN jobs AS visibility_jobs ON visibility_jobs.id = applications.job_id AND visibility_jobs.deleted_at IS NULL").
		Where("visibility_jobs.hiring_manager = ?", scope.UserID)
}
