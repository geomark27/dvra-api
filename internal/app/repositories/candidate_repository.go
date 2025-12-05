package repositories

import (
	"dvra-api/internal/app/models"
	"dvra-api/internal/database"

	"gorm.io/gorm"
)

// CandidateRepository define el contrato del repositorio de candidates
type CandidateRepository interface {
	GetAll() ([]models.Candidate, error)
	GetByID(id uint) (*models.Candidate, error)
	GetByCompanyID(companyID uint) ([]models.Candidate, error)
	GetByEmail(email string, companyID uint) (*models.Candidate, error)
	Create(candidate *models.Candidate) (*models.Candidate, error)
	Update(candidate *models.Candidate) (*models.Candidate, error)
	Delete(id uint) error
}

// candidateRepository es la implementación con GORM
type candidateRepository struct{}

// NewCandidateRepository crea una nueva instancia de CandidateRepository
func NewCandidateRepository() CandidateRepository {
	return &candidateRepository{}
}

func (r *candidateRepository) GetAll() ([]models.Candidate, error) {
	var candidates []models.Candidate
	if err := database.DB.Preload("Company").Find(&candidates).Error; err != nil {
		return nil, err
	}
	return candidates, nil
}

func (r *candidateRepository) GetByID(id uint) (*models.Candidate, error) {
	var candidate models.Candidate
	if err := database.DB.Preload("Company").Preload("Applications").First(&candidate, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &candidate, nil
}

func (r *candidateRepository) GetByCompanyID(companyID uint) ([]models.Candidate, error) {
	var candidates []models.Candidate
	if err := database.DB.Where("company_id = ?", companyID).Find(&candidates).Error; err != nil {
		return nil, err
	}
	return candidates, nil
}

func (r *candidateRepository) GetByEmail(email string, companyID uint) (*models.Candidate, error) {
	var candidate models.Candidate
	if err := database.DB.Where("email = ? AND company_id = ?", email, companyID).First(&candidate).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &candidate, nil
}

func (r *candidateRepository) Create(candidate *models.Candidate) (*models.Candidate, error) {
	if err := database.DB.Create(candidate).Error; err != nil {
		return nil, err
	}
	return candidate, nil
}

func (r *candidateRepository) Update(candidate *models.Candidate) (*models.Candidate, error) {
	if err := database.DB.Save(candidate).Error; err != nil {
		return nil, err
	}
	return candidate, nil
}

func (r *candidateRepository) Delete(id uint) error {
	return database.DB.Delete(&models.Candidate{}, id).Error
}
