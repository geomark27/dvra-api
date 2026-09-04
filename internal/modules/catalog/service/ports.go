package service

import "dvra-api/internal/shared/models"

// SystemValueRepository es el puerto de salida hacia la persistencia del
// catálogo (módulo → categoría → valor, con overrides por empresa).
//
// Se declara en el consumidor (ESTANDARES_DESARROLLO.md §3.3 regla 3).
type SystemValueRepository interface {
	GetCategories(module string) ([]models.SystemValueCategory, error)
	GetCategory(module, category string) (*models.SystemValueCategory, error)
	GetCategoryByKey(category string) (*models.SystemValueCategory, error)
	GetCategoryByID(id uint) (*models.SystemValueCategory, error)
	GetByCategoryID(categoryID uint, companyID *uint) ([]models.SystemValue, error)
	GetAll() ([]models.SystemValue, error)
	Create(value *models.SystemValue) error
	Update(value *models.SystemValue) error
	Delete(id uint) error
	GetByID(id uint) (*models.SystemValue, error)
	GetOverride(companyID, valueID uint) (*models.CompanySystemValueOverride, error)
	UpsertOverride(override *models.CompanySystemValueOverride) error
}
