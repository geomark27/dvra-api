package models

// TenantRole is a built-in or future custom authorization profile for tenant
// memberships. Landlord roles intentionally live outside this model.
type TenantRole struct {
	BaseModel
	Key         string `gorm:"type:varchar(50);uniqueIndex;not null" json:"key"`
	Name        string `gorm:"type:varchar(100);not null" json:"name"`
	Description string `gorm:"type:text" json:"description,omitempty"`
	IsSystem    bool   `gorm:"not null;default:true" json:"is_system"`
	IsActive    bool   `gorm:"not null;default:true;index" json:"is_active"`
}

func (TenantRole) TableName() string { return "tenant_roles" }

// Permission represents an atomic action that can be granted to a tenant role.
type Permission struct {
	BaseModel
	Key         string `gorm:"type:varchar(100);uniqueIndex;not null" json:"key"`
	Module      string `gorm:"type:varchar(50);not null;index" json:"module"`
	Action      string `gorm:"type:varchar(50);not null" json:"action"`
	Description string `gorm:"type:text" json:"description,omitempty"`
	IsActive    bool   `gorm:"not null;default:true;index" json:"is_active"`
}

func (Permission) TableName() string { return "permissions" }

// TenantRolePermission is the explicit many-to-many authorization mapping.
type TenantRolePermission struct {
	TenantRoleID uint `gorm:"primaryKey;not null" json:"tenant_role_id"`
	PermissionID uint `gorm:"primaryKey;not null" json:"permission_id"`

	TenantRole *TenantRole `gorm:"foreignKey:TenantRoleID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	Permission *Permission `gorm:"foreignKey:PermissionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"permission,omitempty"`
}

func (TenantRolePermission) TableName() string { return "tenant_role_permissions" }
