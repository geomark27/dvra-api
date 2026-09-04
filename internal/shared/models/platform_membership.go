package models

// PlatformMembership grants a user access to the SaaS landlord context.
// Tenant access is represented exclusively by Membership.
type PlatformMembership struct {
	BaseModel
	UserID uint   `gorm:"not null;uniqueIndex:idx_platform_user_active,where:deleted_at IS NULL" json:"user_id"`
	Role   string `gorm:"type:varchar(50);not null" json:"role"`
	Status string `gorm:"type:varchar(50);not null;default:'active'" json:"status"`
	User   *User  `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (PlatformMembership) TableName() string { return "platform_memberships" }
