package repositories

import (
	"dvra-api/internal/app/models"
	"dvra-api/internal/database"

	"gorm.io/gorm"
)

// ApplicationRepository define el contrato del repositorio de applications
type ApplicationRepository interface {
	GetAll() ([]models.Application, error)
	GetByID(id uint) (*models.Application, error)
	GetByJobID(jobID uint) ([]models.Application, error)
	GetByCandidateID(candidateID uint) ([]models.Application, error)
	GetByCompanyID(companyID uint) ([]models.Application, error)
	GetByStage(stage string, companyID uint) ([]models.Application, error)
	Create(application *models.Application) (*models.Application, error)
	Update(application *models.Application) (*models.Application, error)
	Delete(id uint) error
}

// applicationRepository es la implementación con GORM
type applicationRepository struct{}

// NewApplicationRepository crea una nueva instancia de ApplicationRepository
func NewApplicationRepository() ApplicationRepository {
	return &applicationRepository{}
}

func (r *applicationRepository) GetAll() ([]models.Application, error) {
	var applications []models.Application
	if err := database.DB.Preload("Job").Preload("Candidate").Preload("Company").Find(&applications).Error; err != nil {
		return nil, err
	}
	return applications, nil
}

func (r *applicationRepository) GetByID(id uint) (*models.Application, error) {
	var application models.Application
	if err := database.DB.Preload("Job").Preload("Candidate").Preload("Company").First(&application, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &application, nil
}

func (r *applicationRepository) GetByJobID(jobID uint) ([]models.Application, error) {
	var applications []models.Application
	if err := database.DB.Where("job_id = ?", jobID).Preload("Candidate").Find(&applications).Error; err != nil {
		return nil, err
	}
	return applications, nil
}

func (r *applicationRepository) GetByCandidateID(candidateID uint) ([]models.Application, error) {
	var applications []models.Application
	if err := database.DB.Where("candidate_id = ?", candidateID).Preload("Job").Find(&applications).Error; err != nil {
		return nil, err
	}
	return applications, nil
}

func (r *applicationRepository) GetByCompanyID(companyID uint) ([]models.Application, error) {
	var applications []models.Application
	if err := database.DB.Where("company_id = ?", companyID).Preload("Job").Preload("Candidate").Find(&applications).Error; err != nil {
		return nil, err
	}
	return applications, nil
}

func (r *applicationRepository) GetByStage(stage string, companyID uint) ([]models.Application, error) {
	var applications []models.Application
	if err := database.DB.Where("stage = ? AND company_id = ?", stage, companyID).Preload("Job").Preload("Candidate").Find(&applications).Error; err != nil {
		return nil, err
	}
	return applications, nil
}

func (r *applicationRepository) Create(application *models.Application) (*models.Application, error) {
	if err := database.DB.Create(application).Error; err != nil {
		return nil, err
	}
	return application, nil
}

func (r *applicationRepository) Update(application *models.Application) (*models.Application, error) {
	if err := database.DB.Save(application).Error; err != nil {
		return nil, err
	}
	return application, nil
}

func (r *applicationRepository) Delete(id uint) error {
	return database.DB.Delete(&models.Application{}, id).Error
}
