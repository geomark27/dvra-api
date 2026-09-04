package models

import (
	"time"

	"gorm.io/gorm"
)

// BaseModel es el equivalente a gorm.Model pero con tags JSON en snake_case
// (id, created_at, updated_at, deleted_at). Todos los modelos del proyecto
// deben embeber BaseModel en lugar de gorm.Model: genera las mismas columnas
// en BD y mantiene el contrato JSON cuando un modelo se serializa directo.
type BaseModel struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
