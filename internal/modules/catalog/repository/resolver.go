package repository

// Resolución de valores globales de catálogo por (módulo, categoría, clave).
//
// Vivía en internal/shared/catalog, pero no era una utilidad transversal: es
// lógica de este dominio con consultas GORM propias, así que su sitio es el
// módulo (hallazgo P10).
//
// Siguen siendo funciones sin estado que reciben la conexión, porque sus
// consumidores actuales (los seeders y los services de billing e iam) todavía
// no están migrados y pasan el *gorm.DB que ya tienen. Cuando esos dominios
// pasen a módulo, declararán su propio puerto mínimo y recibirán el resolver
// inyectado en lugar de llamar aquí.

import (
	"errors"

	"dvra-api/internal/shared/models"

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
