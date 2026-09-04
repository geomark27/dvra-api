package models

// SystemValueModule groups catalogs by bounded context (identity,
// recruitment, staffing, billing, etc.). The key is a stable technical
// identifier; its user-facing name can change without breaking references.
type SystemValueModule struct {
	BaseModel
	Key         string                `gorm:"type:varchar(80);not null;uniqueIndex" json:"key"`
	Name        string                `gorm:"type:varchar(120);not null" json:"name"`
	Description string                `gorm:"type:text" json:"description,omitempty"`
	IsSystem    bool                  `gorm:"not null;default:true" json:"is_system"`
	IsActive    bool                  `gorm:"not null;default:true" json:"is_active"`
	Categories  []SystemValueCategory `gorm:"foreignKey:ModuleID" json:"categories,omitempty"`
}

func (SystemValueModule) TableName() string { return "system_value_modules" }

// SystemValueCategory defines the ownership rules for a catalog. Closed
// categories reject tenant-created values; overridable categories allow a
// tenant to customize presentation without modifying the global row.
type SystemValueCategory struct {
	BaseModel
	ModuleID            uint               `gorm:"not null;uniqueIndex:idx_system_value_category_key,priority:1" json:"module_id"`
	Key                 string             `gorm:"type:varchar(80);not null;uniqueIndex:idx_system_value_category_key,priority:2" json:"key"`
	Name                string             `gorm:"type:varchar(120);not null" json:"name"`
	Description         string             `gorm:"type:text" json:"description,omitempty"`
	IsSystem            bool               `gorm:"not null;default:true" json:"is_system"`
	IsExtensible        bool               `gorm:"not null;default:false" json:"is_extensible"`
	AllowTenantValues   bool               `gorm:"not null;default:false" json:"allow_tenant_values"`
	AllowTenantOverride bool               `gorm:"not null;default:false" json:"allow_tenant_override"`
	IsActive            bool               `gorm:"not null;default:true" json:"is_active"`
	Module              *SystemValueModule `gorm:"foreignKey:ModuleID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"module,omitempty"`
	Values              []SystemValue      `gorm:"foreignKey:CategoryID" json:"values,omitempty"`
}

func (SystemValueCategory) TableName() string { return "system_value_categories" }

// SystemValue is either global (CompanyID=nil) or tenant-owned. Key is stable
// inside its category. For API compatibility it is serialized as "value".
type SystemValue struct {
	BaseModel
	CategoryID   uint                 `gorm:"not null;index" json:"category_id"`
	CompanyID    *uint                `gorm:"index" json:"company_id,omitempty"`
	Key          string               `gorm:"type:varchar(100);not null" json:"value"`
	Label        string               `gorm:"type:varchar(200);not null" json:"label"`
	Description  *string              `gorm:"type:text" json:"description,omitempty"`
	DisplayOrder int                  `gorm:"not null;default:0" json:"display_order"`
	IsSystem     bool                 `gorm:"not null;default:false" json:"is_system"`
	IsActive     bool                 `gorm:"not null;default:true" json:"is_active"`
	Category     *SystemValueCategory `gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"category,omitempty"`
	Company      *Company             `gorm:"foreignKey:CompanyID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"company,omitempty"`
}

func (SystemValue) TableName() string { return "system_values" }

// CompanySystemValueOverride stores presentation-only tenant customizations
// for a global value. It never duplicates or mutates the global catalog row.
type CompanySystemValueOverride struct {
	BaseModel
	CompanyID     uint         `gorm:"not null;uniqueIndex:idx_company_system_value_override,priority:1" json:"company_id"`
	SystemValueID uint         `gorm:"not null;uniqueIndex:idx_company_system_value_override,priority:2" json:"system_value_id"`
	CustomLabel   *string      `gorm:"type:varchar(200)" json:"custom_label,omitempty"`
	DisplayOrder  *int         `json:"display_order,omitempty"`
	IsEnabled     *bool        `json:"is_enabled,omitempty"`
	Company       *Company     `gorm:"foreignKey:CompanyID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"company,omitempty"`
	SystemValue   *SystemValue `gorm:"foreignKey:SystemValueID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"system_value,omitempty"`
}

func (CompanySystemValueOverride) TableName() string {
	return "company_system_value_overrides"
}
