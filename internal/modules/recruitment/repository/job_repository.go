package repository

import (
	"dvra-api/internal/modules/recruitment/service"
	"dvra-api/internal/shared/models"
	"dvra-api/internal/shared/permissions"
	"dvra-api/internal/shared/tenant"
	"dvra-api/internal/shared/tenant/tenantgorm"

	"gorm.io/gorm"
)

// JobRepository define el contrato del repositorio de jobs
// jobRepository es la implementación con GORM
type JobRepository struct{ db *gorm.DB }

// NewJobRepository crea una nueva instancia de JobRepository
func NewJobRepository(db *gorm.DB) *JobRepository {
	return &JobRepository{db: db}
}

func (r *JobRepository) GetAll() ([]models.Job, error) {
	var jobs []models.Job
	if err := r.db.Preload("Company").Preload("Department.Manager").Find(&jobs).Error; err != nil {
		return nil, err
	}
	if err := attachJobDepartmentPaths(r.db, jobs); err != nil {
		return nil, err
	}
	return jobs, nil
}

func (r *JobRepository) GetAllWithFilters(filters service.JobFilters) ([]models.Job, error) {
	var jobs []models.Job
	query := r.db.Preload("Company").Preload("City.State.Country").Preload("Department.Manager")

	// Aplicar filtros
	if filters.Status != "" {
		query = query.Where("status = ?", filters.Status)
	}
	if filters.LocationType != "" {
		query = query.Where("location_type = ?", filters.LocationType)
	}
	if filters.CityID != nil {
		query = query.Where("city_id = ?", *filters.CityID)
	}
	if filters.StaffingClientID != nil {
		query = query.Where("staffing_client_id = ?", *filters.StaffingClientID)
	}
	if filters.DepartmentID != nil {
		query = query.Where("department_id = ?", *filters.DepartmentID)
	}

	if err := query.Find(&jobs).Error; err != nil {
		return nil, err
	}
	if err := attachJobDepartmentPaths(r.db, jobs); err != nil {
		return nil, err
	}
	return jobs, nil
}

func (r *JobRepository) GetByID(id uint, scope tenant.Scope) (*models.Job, error) {
	var job models.Job
	query := tenantgorm.Apply(r.db, scope, "jobs.company_id")
	query = applyJobActorScope(query, scope)
	if err := query.Preload("Company").Preload("City.State.Country").Preload("Department.Manager").Preload("Applications").Where("jobs.id = ?", id).First(&job).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	jobs := []models.Job{job}
	if err := attachJobDepartmentPaths(r.db, jobs); err != nil {
		return nil, err
	}
	return &jobs[0], nil
}

func (r *JobRepository) GetByScopeWithFilters(scope tenant.Scope, filters service.JobFilters) ([]models.Job, error) {
	var jobs []models.Job
	query := tenantgorm.Apply(r.db.Preload("City.State.Country").Preload("Department.Manager"), scope, "jobs.company_id")
	query = applyJobActorScope(query, scope)
	query = applyJobFilters(query, filters)
	if err := query.Find(&jobs).Error; err != nil {
		return nil, err
	}
	if err := attachJobDepartmentPaths(r.db, jobs); err != nil {
		return nil, err
	}
	return jobs, nil
}

func (r *JobRepository) GetByCompanyID(companyID uint) ([]models.Job, error) {
	var jobs []models.Job
	if err := r.db.Where("company_id = ?", companyID).Preload("Department.Manager").Find(&jobs).Error; err != nil {
		return nil, err
	}
	if err := attachJobDepartmentPaths(r.db, jobs); err != nil {
		return nil, err
	}
	return jobs, nil
}

func (r *JobRepository) GetByCompanyIDWithFilters(companyID uint, filters service.JobFilters) ([]models.Job, error) {
	var jobs []models.Job
	query := r.db.Preload("City.State.Country").Preload("Department.Manager").Where("company_id = ?", companyID)
	query = applyJobFilters(query, filters)

	if err := query.Find(&jobs).Error; err != nil {
		return nil, err
	}
	if err := attachJobDepartmentPaths(r.db, jobs); err != nil {
		return nil, err
	}
	return jobs, nil
}

func applyJobFilters(query *gorm.DB, filters service.JobFilters) *gorm.DB {
	if filters.Status != "" {
		query = query.Where("status = ?", filters.Status)
	}
	if filters.LocationType != "" {
		query = query.Where("location_type = ?", filters.LocationType)
	}
	if filters.CityID != nil {
		query = query.Where("city_id = ?", *filters.CityID)
	}
	if filters.StaffingClientID != nil {
		query = query.Where("staffing_client_id = ?", *filters.StaffingClientID)
	}
	if filters.DepartmentID != nil {
		query = query.Where("department_id = ?", *filters.DepartmentID)
	}
	return query
}

func applyJobActorScope(query *gorm.DB, scope tenant.Scope) *gorm.DB {
	if scope.Role == permissions.RoleHiringManager {
		return query.Where("jobs.hiring_manager = ?", scope.UserID)
	}
	return query
}

func (r *JobRepository) GetByStatus(status string, companyID uint) ([]models.Job, error) {
	var jobs []models.Job
	if err := r.db.Where("status = ? AND company_id = ?", status, companyID).Preload("Department.Manager").Find(&jobs).Error; err != nil {
		return nil, err
	}
	if err := attachJobDepartmentPaths(r.db, jobs); err != nil {
		return nil, err
	}
	return jobs, nil
}

// GetPublishedByCompanyID returns all published jobs for a company - for public career page
func (r *JobRepository) GetPublishedByCompanyID(companyID uint) ([]models.Job, error) {
	var jobs []models.Job
	if err := r.db.
		Where("company_id = ? AND status = ?", companyID, "published").
		Preload("Company").
		Preload("City.State.Country").
		Preload("Department").
		Order("created_at DESC").
		Find(&jobs).Error; err != nil {
		return nil, err
	}
	if err := attachJobDepartmentPaths(r.db, jobs); err != nil {
		return nil, err
	}
	return jobs, nil
}

func (r *JobRepository) Create(job *models.Job) (*models.Job, error) {
	if err := r.db.Create(job).Error; err != nil {
		return nil, err
	}
	return job, nil
}

func (r *JobRepository) Update(job *models.Job) (*models.Job, error) {
	if err := r.db.Save(job).Error; err != nil {
		return nil, err
	}
	return job, nil
}

func (r *JobRepository) Delete(id uint) error {
	return r.db.Delete(&models.Job{}, id).Error
}
