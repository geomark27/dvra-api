package repositories

import (
	"dvra-api/internal/app/models"

	"gorm.io/gorm"
)

// CompanyRepository define el contrato del repositorio de companies
type CompanyRepository interface {
	GetAll() ([]models.Company, error)
	GetByID(id uint) (*models.Company, error)
	GetBySlug(slug string) (*models.Company, error)
	Create(company *models.Company) (*models.Company, error)
	Update(company *models.Company) (*models.Company, error)
	Delete(id uint) error
	GetCompaniesWithMembers(companyID uint) (*models.Company, error)
}

// companyRepository es la implementación con GORM
type companyRepository struct{ db *gorm.DB }

// NewCompanyRepository crea una nueva instancia de CompanyRepository
func NewCompanyRepository(db *gorm.DB) CompanyRepository {
	return &companyRepository{db: db}
}

func (r *companyRepository) GetAll() ([]models.Company, error) {
	var companies []models.Company
	if err := r.db.Find(&companies).Error; err != nil {
		return nil, err
	}
	return companies, nil
}

func (r *companyRepository) GetByID(id uint) (*models.Company, error) {
	var company models.Company
	if err := r.db.First(&company, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &company, nil
}

func (r *companyRepository) GetBySlug(slug string) (*models.Company, error) {
	var company models.Company
	if err := r.db.Where("slug = ?", slug).First(&company).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &company, nil
}

func (r *companyRepository) Create(company *models.Company) (*models.Company, error) {
	if err := r.db.Create(company).Error; err != nil {
		return nil, err
	}
	return company, nil
}

func (r *companyRepository) Update(company *models.Company) (*models.Company, error) {
	if err := r.db.Save(company).Error; err != nil {
		return nil, err
	}
	return company, nil
}

func (r *companyRepository) Delete(id uint) error {
	return r.db.Delete(&models.Company{}, id).Error
}

func (r *companyRepository) GetCompaniesWithMembers(companyID uint) (*models.Company, error) {
	var company models.Company
	if err := r.db.Preload("Memberships").Preload("Memberships.User").First(&company, companyID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &company, nil
}
