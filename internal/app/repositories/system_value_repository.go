package repositories

import (
	"dvra-api/internal/app/models"
	"errors"

	"gorm.io/gorm"
)

type SystemValueRepository interface {
	GetByCategory(category string, companyID *uint) ([]models.SystemValue, error)
	GetAll() ([]models.SystemValue, error)
	Create(value *models.SystemValue) error
	Update(value *models.SystemValue) error
	Delete(id uint) error
	GetByID(id uint) (*models.SystemValue, error)
}

type systemValueRepository struct {
	db *gorm.DB
}

func NewSystemValueRepository(db *gorm.DB) SystemValueRepository {
	return &systemValueRepository{db: db}
}

// GetByCategory retrieves all active values for a category (global + company-specific)
func (r *systemValueRepository) GetByCategory(category string, companyID *uint) ([]models.SystemValue, error) {
	var values []models.SystemValue

	query := r.db.Where("category = ? AND is_active = ?", category, true)

	// Include global values (company_id IS NULL) and company-specific values
	if companyID != nil {
		query = query.Where("company_id IS NULL OR company_id = ?", *companyID)
	} else {
		query = query.Where("company_id IS NULL")
	}

	err := query.Order("display_order ASC, label ASC").Find(&values).Error
	return values, err
}

// GetAll retrieves all system values
func (r *systemValueRepository) GetAll() ([]models.SystemValue, error) {
	var values []models.SystemValue
	err := r.db.Order("category ASC, display_order ASC").Find(&values).Error
	return values, err
}

// Create creates a new system value
func (r *systemValueRepository) Create(value *models.SystemValue) error {
	// Check for duplicate category + value + company_id
	var existing models.SystemValue
	query := r.db.Where("category = ? AND value = ?", value.Category, value.Value)

	if value.CompanyID != nil {
		query = query.Where("company_id = ?", *value.CompanyID)
	} else {
		query = query.Where("company_id IS NULL")
	}

	if err := query.First(&existing).Error; err == nil {
		return errors.New("system value already exists for this category and company")
	}

	return r.db.Create(value).Error
}

// Update updates a system value
func (r *systemValueRepository) Update(value *models.SystemValue) error {
	return r.db.Save(value).Error
}

// Delete soft deletes a system value
func (r *systemValueRepository) Delete(id uint) error {
	return r.db.Delete(&models.SystemValue{}, id).Error
}

// GetByID retrieves a system value by ID
func (r *systemValueRepository) GetByID(id uint) (*models.SystemValue, error) {
	var value models.SystemValue
	err := r.db.First(&value, id).Error
	if err != nil {
		return nil, err
	}
	return &value, nil
}
