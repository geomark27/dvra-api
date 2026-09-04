package models

// Plan represents a subscription plan
type Plan struct {
	BaseModel
	Name                string       `gorm:"type:varchar(100);uniqueIndex;not null" json:"name"`
	Slug                string       `gorm:"type:varchar(100);uniqueIndex;not null" json:"slug"`
	Description         string       `gorm:"type:text" json:"description"`
	Price               float64      `gorm:"type:decimal(10,2);not null;default:0" json:"price"`
	CurrencyValueID     uint         `gorm:"not null;index" json:"currency_value_id"`
	BillingCycleValueID uint         `gorm:"not null;index" json:"billing_cycle_value_id"`
	IsActive            bool         `gorm:"default:true" json:"is_active"`
	IsPublic            bool         `gorm:"default:true" json:"is_public"` // Si aparece en listado público
	TrialDays           int          `gorm:"default:0" json:"trial_days"`
	DisplayOrder        int          `gorm:"default:0" json:"display_order"`
	MaxUsers            int          `gorm:"default:-1" json:"max_users"`        // -1 = ilimitado
	MaxJobs             int          `gorm:"default:-1" json:"max_jobs"`         // -1 = ilimitado
	MaxCandidates       int          `gorm:"default:-1" json:"max_candidates"`   // -1 = ilimitado
	MaxApplications     int          `gorm:"default:-1" json:"max_applications"` // -1 = ilimitado
	MaxStorageGB        int          `gorm:"default:-1" json:"max_storage_gb"`   // -1 = ilimitado
	CanExportData       bool         `gorm:"default:false" json:"can_export_data"`
	CanUseCustomBrand   bool         `gorm:"default:false" json:"can_use_custom_brand"`
	CanUseAPI           bool         `gorm:"default:false" json:"can_use_api"`
	CanUseIntegrations  bool         `gorm:"default:false" json:"can_use_integrations"`
	CanUseStaffing      bool         `gorm:"default:false" json:"can_use_staffing"` // Habilita módulo de staffing/outsourcing
	SupportLevelValueID uint         `gorm:"not null;index" json:"support_level_value_id"`
	Currency            *SystemValue `gorm:"foreignKey:CurrencyValueID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"currency,omitempty"`
	BillingCycle        *SystemValue `gorm:"foreignKey:BillingCycleValueID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"billing_cycle,omitempty"`
	SupportLevel        *SystemValue `gorm:"foreignKey:SupportLevelValueID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"support_level,omitempty"`
}

func (Plan) TableName() string {
	return "plans"
}

// IsUnlimited checks if a specific limit is unlimited
func (p *Plan) IsUnlimited(limitType string) bool {
	switch limitType {
	case "users":
		return p.MaxUsers == -1
	case "jobs":
		return p.MaxJobs == -1
	case "candidates":
		return p.MaxCandidates == -1
	case "applications":
		return p.MaxApplications == -1
	case "storage":
		return p.MaxStorageGB == -1
	default:
		return false
	}
}

// HasFeature checks if plan includes a specific feature
func (p *Plan) HasFeature(feature string) bool {
	switch feature {
	case "export_data":
		return p.CanExportData
	case "custom_brand":
		return p.CanUseCustomBrand
	case "api":
		return p.CanUseAPI
	case "integrations":
		return p.CanUseIntegrations
	case "staffing":
		return p.CanUseStaffing
	default:
		return false
	}
}
