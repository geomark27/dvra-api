package models

import (
	"gorm.io/gorm"
)

// Role represents the role entity in the database
type Role struct {
	gorm.Model
	Name        string `gorm:"type:varchar(50);uniqueIndex;not null" json:"name"`
	Slug        string `gorm:"type:varchar(50);uniqueIndex;not null" json:"slug"`
	Description string `gorm:"type:text" json:"description,omitempty"`
	Level       int    `gorm:"not null;default:0" json:"level"` // Para jerarquía: 100=superadmin, 50=admin, 10=user
	IsSystem    bool   `gorm:"default:false" json:"is_system"`  // Roles del sistema no se pueden eliminar
}

func (Role) TableName() string {
	return "roles"
}

// Constantes para roles del sistema
const (
	RoleSuperAdmin = "superadmin"
	RoleAdmin      = "admin"
	RoleRecruiter  = "recruiter"
	RoleUser       = "user"
)
