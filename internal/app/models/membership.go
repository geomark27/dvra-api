package models

import (
	"time"

	"gorm.io/gorm"
)

// Membership represents the companymembership entity in the database
type Membership struct {
	gorm.Model
	UserID    uint       `gorm:"not null;index:idx_user_company,priority:1" json:"user_id"`
	CompanyID *uint      `gorm:"index:idx_user_company,priority:2" json:"company_id,omitempty"` // Nullable para SuperAdmin
	Role      string     `gorm:"type:varchar(50);not null" json:"role"`
	Status    string     `gorm:"type:varchar(50);default:'active'" json:"status"`
	IsDefault bool       `gorm:"default:false" json:"is_default"`
	InvitedBy *uint      `gorm:"" json:"invited_by,omitempty"`
	InvitedAt *time.Time `gorm:"type:timestamp" json:"invited_at,omitempty"`
	JoinedAt  *time.Time `gorm:"type:timestamp" json:"joined_at,omitempty"`

	// Relaciones (para eager loading)
	User    *User    `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Company *Company `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
}

// TableName overrides the table name (optional)
func (Membership) TableName() string {
	return "memberships"
}
