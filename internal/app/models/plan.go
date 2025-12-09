package models

import (
	"gorm.io/gorm"
)

// Plan represents a subscription plan
type Plan struct {
	gorm.Model
	Name               string  `gorm:"type:varchar(100);uniqueIndex;not null" json:"name"`
	Slug               string  `gorm:"type:varchar(100);uniqueIndex;not null" json:"slug"`
	Description        string  `gorm:"type:text" json:"description"`
	Price              float64 `gorm:"type:decimal(10,2);not null;default:0" json:"price"`
	Currency           string  `gorm:"type:varchar(3);not null;default:'USD'" json:"currency"`
	BillingCycle       string  `gorm:"type:varchar(20);not null;default:'monthly'" json:"billing_cycle"` // monthly, yearly
	IsActive           bool    `gorm:"default:true" json:"is_active"`
	IsPublic           bool    `gorm:"default:true" json:"is_public"` // Si aparece en listado público
	TrialDays          int     `gorm:"default:0" json:"trial_days"`
	DisplayOrder       int     `gorm:"default:0" json:"display_order"`
	MaxUsers           int     `gorm:"default:-1" json:"max_users"`        // -1 = ilimitado
	MaxJobs            int     `gorm:"default:-1" json:"max_jobs"`         // -1 = ilimitado
	MaxCandidates      int     `gorm:"default:-1" json:"max_candidates"`   // -1 = ilimitado
	MaxApplications    int     `gorm:"default:-1" json:"max_applications"` // -1 = ilimitado
	MaxStorageGB       int     `gorm:"default:-1" json:"max_storage_gb"`   // -1 = ilimitado
	CanExportData      bool    `gorm:"default:false" json:"can_export_data"`
	CanUseCustomBrand  bool    `gorm:"default:false" json:"can_use_custom_brand"`
	CanUseAPI          bool    `gorm:"default:false" json:"can_use_api"`
	CanUseIntegrations bool    `gorm:"default:false" json:"can_use_integrations"`
	SupportLevel       string  `gorm:"type:varchar(50);default:'email'" json:"support_level"` // email, priority, dedicated
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
	default:
		return false
	}
}
