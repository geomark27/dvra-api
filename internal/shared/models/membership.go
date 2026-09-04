package models

import (
	"time"
)

// Membership represents the companymembership entity in the database
type Membership struct {
	BaseModel

	UserID       uint       `gorm:"not null;uniqueIndex:idx_user_company,priority:1,where:deleted_at IS NULL" json:"user_id"`
	CompanyID    uint       `gorm:"not null;uniqueIndex:idx_user_company,priority:2,where:deleted_at IS NULL" json:"company_id"`
	TenantRoleID uint       `gorm:"not null;index" json:"tenant_role_id"`
	Status       string     `gorm:"type:varchar(50);default:'active'" json:"status"`
	IsDefault    bool       `gorm:"default:false" json:"is_default"`
	InvitedBy    *uint      `gorm:"" json:"invited_by,omitempty"`
	InvitedAt    *time.Time `gorm:"type:timestamp" json:"invited_at,omitempty"`
	JoinedAt     *time.Time `gorm:"type:timestamp" json:"joined_at,omitempty"`

	// Relaciones (para eager loading)
	User       *User       `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Company    *Company    `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	TenantRole *TenantRole `gorm:"foreignKey:TenantRoleID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"tenant_role,omitempty"`
}

// TableName overrides the table name (optional)
func (Membership) TableName() string {
	return "memberships"
}
