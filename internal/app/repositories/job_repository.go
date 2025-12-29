package repositories

import (
	"dvra-api/internal/app/dtos"
	"dvra-api/internal/app/models"
	"dvra-api/internal/database"

	"gorm.io/gorm"
)

// JobRepository define el contrato del repositorio de jobs
type JobRepository interface {
	GetAll() ([]models.Job, error)
	GetAllWithFilters(filters dtos.JobFilters) ([]models.Job, error)
	GetByID(id uint) (*models.Job, error)
	GetByCompanyID(companyID uint) ([]models.Job, error)
	GetByCompanyIDWithFilters(companyID uint, filters dtos.JobFilters) ([]models.Job, error)
	GetByStatus(status string, companyID uint) ([]models.Job, error)
	Create(job *models.Job) (*models.Job, error)
	Update(job *models.Job) (*models.Job, error)
	Delete(id uint) error
}

// jobRepository es la implementación con GORM
type jobRepository struct{}

// NewJobRepository crea una nueva instancia de JobRepository
func NewJobRepository() JobRepository {
	return &jobRepository{}
}

func (r *jobRepository) GetAll() ([]models.Job, error) {
	var jobs []models.Job
	if err := database.DB.Preload("Company").Find(&jobs).Error; err != nil {
		return nil, err
	}
	return jobs, nil
}

func (r *jobRepository) GetAllWithFilters(filters dtos.JobFilters) ([]models.Job, error) {
	var jobs []models.Job
	query := database.DB.Preload("Company").Preload("City.State.Country")

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

	if err := query.Find(&jobs).Error; err != nil {
		return nil, err
	}
	return jobs, nil
}

func (r *jobRepository) GetByID(id uint) (*models.Job, error) {
	var job models.Job
	if err := database.DB.Preload("Company").Preload("City.State.Country").Preload("Applications").First(&job, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &job, nil
}

func (r *jobRepository) GetByCompanyID(companyID uint) ([]models.Job, error) {
	var jobs []models.Job
	if err := database.DB.Where("company_id = ?", companyID).Find(&jobs).Error; err != nil {
		return nil, err
	}
	return jobs, nil
}

func (r *jobRepository) GetByCompanyIDWithFilters(companyID uint, filters dtos.JobFilters) ([]models.Job, error) {
	var jobs []models.Job
	query := database.DB.Preload("City.State.Country").Where("company_id = ?", companyID)

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

	if err := query.Find(&jobs).Error; err != nil {
		return nil, err
	}
	return jobs, nil
}

func (r *jobRepository) GetByStatus(status string, companyID uint) ([]models.Job, error) {
	var jobs []models.Job
	if err := database.DB.Where("status = ? AND company_id = ?", status, companyID).Find(&jobs).Error; err != nil {
		return nil, err
	}
	return jobs, nil
}

func (r *jobRepository) Create(job *models.Job) (*models.Job, error) {
	if err := database.DB.Create(job).Error; err != nil {
		return nil, err
	}
	return job, nil
}

func (r *jobRepository) Update(job *models.Job) (*models.Job, error) {
	if err := database.DB.Save(job).Error; err != nil {
		return nil, err
	}
	return job, nil
}

func (r *jobRepository) Delete(id uint) error {
	return database.DB.Delete(&models.Job{}, id).Error
}
