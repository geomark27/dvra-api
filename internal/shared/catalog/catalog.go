package catalog

import (
	"errors"

	"dvra-api/internal/app/models"

	"gorm.io/gorm"
)

func ResolveGlobal(db *gorm.DB, module, category, key string) (*models.SystemValue, error) {
	var value models.SystemValue
	err := db.Joins("JOIN system_value_categories c ON c.id = system_values.category_id").
		Joins("JOIN system_value_modules m ON m.id = c.module_id").
		Where("m.key = ? AND c.key = ? AND system_values.key = ? AND system_values.company_id IS NULL AND system_values.is_active = ?", module, category, key, true).
		First(&value).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &value, err
}

func IsGlobalValue(db *gorm.DB, id uint, module, category string) (bool, error) {
	var count int64
	err := db.Model(&models.SystemValue{}).
		Joins("JOIN system_value_categories c ON c.id = system_values.category_id").
		Joins("JOIN system_value_modules m ON m.id = c.module_id").
		Where("system_values.id = ? AND system_values.company_id IS NULL AND system_values.is_active = ? AND m.key = ? AND c.key = ?", id, true, module, category).Count(&count).Error
	return count == 1, err
}
