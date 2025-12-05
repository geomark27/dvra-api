package repositories

import (
	"dvra-api/internal/app/models"
	"dvra-api/internal/database"

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
type companyRepository struct{}

// NewCompanyRepository crea una nueva instancia de CompanyRepository
func NewCompanyRepository() CompanyRepository {
	return &companyRepository{}
}

func (r *companyRepository) GetAll() ([]models.Company, error) {
	var companies []models.Company
	if err := database.DB.Find(&companies).Error; err != nil {
		return nil, err
	}
	return companies, nil
}

func (r *companyRepository) GetByID(id uint) (*models.Company, error) {
	var company models.Company
	if err := database.DB.First(&company, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &company, nil
}

func (r *companyRepository) GetBySlug(slug string) (*models.Company, error) {
	var company models.Company
	if err := database.DB.Where("slug = ?", slug).First(&company).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &company, nil
}

func (r *companyRepository) Create(company *models.Company) (*models.Company, error) {
	if err := database.DB.Create(company).Error; err != nil {
		return nil, err
	}
	return company, nil
}

func (r *companyRepository) Update(company *models.Company) (*models.Company, error) {
	if err := database.DB.Save(company).Error; err != nil {
		return nil, err
	}
	return company, nil
}

func (r *companyRepository) Delete(id uint) error {
	return database.DB.Delete(&models.Company{}, id).Error
}

func (r *companyRepository) GetCompaniesWithMembers(companyID uint) (*models.Company, error) {
	var company models.Company
	if err := database.DB.Preload("Memberships").Preload("Memberships.User").First(&company, companyID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &company, nil
}
