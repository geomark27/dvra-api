package models

// Role represents the role entity in the database
type Role struct {
	BaseModel
	Name        string `gorm:"type:varchar(50);uniqueIndex;not null" json:"name"`
	Slug        string `gorm:"type:varchar(50);uniqueIndex;not null" json:"slug"`
	Description string `gorm:"type:text" json:"description,omitempty"`
	Level       int    `gorm:"not null;default:0" json:"level"` // Para jerarquía: 50=admin, 30=recruiter, 10=user
	IsSystem    bool   `gorm:"default:false" json:"is_system"`  // Roles del sistema no se pueden eliminar
}

func (Role) TableName() string {
	return "roles"
}

// Constantes para roles del sistema
const (
	RoleAdmin     = "admin"
	RoleRecruiter = "recruiter"
	RoleUser      = "user"
)
