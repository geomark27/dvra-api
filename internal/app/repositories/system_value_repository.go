package repositories

import (
	"errors"
	"sort"

	"dvra-api/internal/app/models"

	"gorm.io/gorm"
)

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

type systemValueRepository struct{ db *gorm.DB }

func NewSystemValueRepository(db *gorm.DB) SystemValueRepository {
	return &systemValueRepository{db: db}
}

func (r *systemValueRepository) GetCategories(module string) ([]models.SystemValueCategory, error) {
	var categories []models.SystemValueCategory
	query := r.db.Preload("Module").Where("system_value_categories.is_active = ?", true)
	if module != "" {
		query = query.Joins("JOIN system_value_modules m ON m.id = system_value_categories.module_id").Where("m.key = ?", module)
	}
	err := query.Order("system_value_categories.module_id, system_value_categories.name").Find(&categories).Error
	return categories, err
}

func (r *systemValueRepository) GetCategory(module, category string) (*models.SystemValueCategory, error) {
	var result models.SystemValueCategory
	err := r.db.Preload("Module").Joins("JOIN system_value_modules m ON m.id = system_value_categories.module_id").
		Where("m.key = ? AND system_value_categories.key = ? AND m.is_active = ? AND system_value_categories.is_active = ?", module, category, true, true).
		First(&result).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &result, err
}

func (r *systemValueRepository) GetCategoryByKey(category string) (*models.SystemValueCategory, error) {
	var results []models.SystemValueCategory
	if err := r.db.Preload("Module").Where("key = ? AND is_active = ?", category, true).Limit(2).Find(&results).Error; err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, nil
	}
	if len(results) > 1 {
		return nil, errors.New("ambiguous category; use module and category")
	}
	return &results[0], nil
}

func (r *systemValueRepository) GetCategoryByID(id uint) (*models.SystemValueCategory, error) {
	var category models.SystemValueCategory
	err := r.db.Preload("Module").First(&category, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &category, err
}

func (r *systemValueRepository) GetByCategoryID(categoryID uint, companyID *uint) ([]models.SystemValue, error) {
	var result []models.SystemValue
	query := r.db.Preload("Category.Module").Where("category_id = ? AND is_active = ?", categoryID, true)
	if companyID == nil {
		query = query.Where("company_id IS NULL")
	} else {
		query = query.Where("company_id IS NULL OR company_id = ?", *companyID)
	}
	if err := query.Order("display_order, label").Find(&result).Error; err != nil {
		return nil, err
	}
	if companyID == nil || len(result) == 0 {
		return result, nil
	}
	var overrides []models.CompanySystemValueOverride
	if err := r.db.Where("company_id = ? AND system_value_id IN ?", *companyID, valueIDs(result)).Find(&overrides).Error; err != nil {
		return nil, err
	}
	byValue := make(map[uint]models.CompanySystemValueOverride, len(overrides))
	for _, override := range overrides {
		byValue[override.SystemValueID] = override
	}
	effective := make([]models.SystemValue, 0, len(result))
	for _, value := range result {
		if override, ok := byValue[value.ID]; ok {
			if override.IsEnabled != nil && !*override.IsEnabled {
				continue
			}
			if override.CustomLabel != nil {
				value.Label = *override.CustomLabel
			}
			if override.DisplayOrder != nil {
				value.DisplayOrder = *override.DisplayOrder
			}
		}
		effective = append(effective, value)
	}
	sort.SliceStable(effective, func(i, j int) bool {
		if effective[i].DisplayOrder == effective[j].DisplayOrder {
			return effective[i].Label < effective[j].Label
		}
		return effective[i].DisplayOrder < effective[j].DisplayOrder
	})
	return effective, nil
}

func valueIDs(values []models.SystemValue) []uint {
	ids := make([]uint, len(values))
	for i := range values {
		ids[i] = values[i].ID
	}
	return ids
}

func (r *systemValueRepository) GetAll() ([]models.SystemValue, error) {
	var values []models.SystemValue
	err := r.db.Preload("Category.Module").Order("category_id, display_order").Find(&values).Error
	return values, err
}

func (r *systemValueRepository) Create(value *models.SystemValue) error {
	return r.db.Create(value).Error
}
func (r *systemValueRepository) Update(value *models.SystemValue) error {
	return r.db.Save(value).Error
}
func (r *systemValueRepository) Delete(id uint) error {
	return r.db.Delete(&models.SystemValue{}, id).Error
}

func (r *systemValueRepository) GetByID(id uint) (*models.SystemValue, error) {
	var value models.SystemValue
	err := r.db.Preload("Category.Module").First(&value, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &value, err
}

func (r *systemValueRepository) GetOverride(companyID, valueID uint) (*models.CompanySystemValueOverride, error) {
	var override models.CompanySystemValueOverride
	err := r.db.Where("company_id = ? AND system_value_id = ?", companyID, valueID).First(&override).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &override, err
}

func (r *systemValueRepository) UpsertOverride(override *models.CompanySystemValueOverride) error {
	var existing models.CompanySystemValueOverride
	err := r.db.Where("company_id = ? AND system_value_id = ?", override.CompanyID, override.SystemValueID).First(&existing).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return r.db.Create(override).Error
	}
	if err != nil {
		return err
	}
	existing.CustomLabel, existing.DisplayOrder, existing.IsEnabled = override.CustomLabel, override.DisplayOrder, override.IsEnabled
	return r.db.Save(&existing).Error
}
