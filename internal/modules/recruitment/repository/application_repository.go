package repository

import (
	"dvra-api/internal/shared/models"
	"dvra-api/internal/shared/permissions"
	"dvra-api/internal/shared/tenant"
	"dvra-api/internal/shared/tenant/tenantgorm"

	"gorm.io/gorm"
)

// ApplicationRepository define el contrato del repositorio de applications
// applicationRepository es la implementación con GORM
type ApplicationRepository struct{ db *gorm.DB }

// NewApplicationRepository crea una nueva instancia de ApplicationRepository
func NewApplicationRepository(db *gorm.DB) *ApplicationRepository {
	return &ApplicationRepository{db: db}
}

func (r *ApplicationRepository) GetAll() ([]models.Application, error) {
	var applications []models.Application
	if err := r.db.Preload("Job").Preload("Candidate").Preload("Company").Find(&applications).Error; err != nil {
		return nil, err
	}
	return applications, nil
}

func (r *ApplicationRepository) GetByID(id uint, scope tenant.Scope) (*models.Application, error) {
	var application models.Application
	query := tenantgorm.Apply(r.db, scope, "applications.company_id")
	query = applyApplicationActorScope(query, scope)
	if err := query.Preload("Job").Preload("Candidate").Preload("Company").Where("applications.id = ?", id).First(&application).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &application, nil
}

func (r *ApplicationRepository) GetByScope(scope tenant.Scope) ([]models.Application, error) {
	var applications []models.Application
	query := tenantgorm.Apply(r.db, scope, "applications.company_id")
	query = applyApplicationActorScope(query, scope)
	if err := query.Preload("Job").Preload("Candidate").Find(&applications).Error; err != nil {
		return nil, err
	}
	return applications, nil
}

func (r *ApplicationRepository) GetByJobID(jobID uint) ([]models.Application, error) {
	var applications []models.Application
	if err := r.db.Where("job_id = ?", jobID).Preload("Candidate").Find(&applications).Error; err != nil {
		return nil, err
	}
	return applications, nil
}

func (r *ApplicationRepository) GetByCandidateID(candidateID uint) ([]models.Application, error) {
	var applications []models.Application
	if err := r.db.Where("candidate_id = ?", candidateID).Preload("Job").Find(&applications).Error; err != nil {
		return nil, err
	}
	return applications, nil
}

func (r *ApplicationRepository) GetByCompanyID(companyID uint) ([]models.Application, error) {
	var applications []models.Application
	if err := r.db.Where("company_id = ?", companyID).Preload("Job").Preload("Candidate").Find(&applications).Error; err != nil {
		return nil, err
	}
	return applications, nil
}

func (r *ApplicationRepository) GetByStage(stage string, companyID uint) ([]models.Application, error) {
	var applications []models.Application
	if err := r.db.Where("stage = ? AND company_id = ?", stage, companyID).Preload("Job").Preload("Candidate").Find(&applications).Error; err != nil {
		return nil, err
	}
	return applications, nil
}

// GetByCandidateAndJob checks if a candidate has already applied to a specific job
func (r *ApplicationRepository) GetByCandidateAndJob(candidateID, jobID uint) (*models.Application, error) {
	var application models.Application
	if err := r.db.Where("candidate_id = ? AND job_id = ?", candidateID, jobID).First(&application).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &application, nil
}

func (r *ApplicationRepository) Create(application *models.Application) (*models.Application, error) {
	if err := r.db.Create(application).Error; err != nil {
		return nil, err
	}
	return application, nil
}

func (r *ApplicationRepository) Update(application *models.Application) (*models.Application, error) {
	if err := r.db.Save(application).Error; err != nil {
		return nil, err
	}
	return application, nil
}

func (r *ApplicationRepository) Delete(id uint) error {
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
