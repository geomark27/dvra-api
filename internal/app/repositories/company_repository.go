package repositories

import (
	"dvra-api/internal/app/models"
	billingctx "dvra-api/internal/shared/billing"

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
	for i := range companies {
		if err := r.attachPlan(&companies[i]); err != nil {
			return nil, err
		}
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
	if err := r.attachPlan(&company); err != nil {
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
	if err := r.attachPlan(&company); err != nil {
		return nil, err
	}
	return &company, nil
}

func (r *companyRepository) attachPlan(company *models.Company) error {
	plan, err := billingctx.CurrentPlan(r.db, company.ID)
	if err != nil {
		return err
	}
	if plan != nil {
		company.PlanTier = plan.Slug
	}
	return nil
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
