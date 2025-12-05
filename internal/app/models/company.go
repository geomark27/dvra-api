package models

import (
	"time"

	"gorm.io/gorm"
)

// Company represents the company entity in the database
type Company struct {
	gorm.Model
	Name        string       `gorm:"type:varchar(255);not null" json:"name"`
	Slug        string       `gorm:"type:varchar(100);uniqueIndex;not null" json:"slug"`
	LogoURL     string       `gorm:"type:text" json:"logo_url,omitempty"`
	PlanTier    string       `gorm:"type:varchar(50);not null;default:'free'" json:"plan_tier"`
	TrialEndsAt *time.Time   `gorm:"type:timestamp" json:"trial_ends_at,omitempty"`
	Timezone    string       `gorm:"type:varchar(100);default:'America/Bogota'" json:"timezone"`
	Memberships []Membership `gorm:"foreignKey:CompanyID" json:"memberships,omitempty"`
}

func (Company) TableName() string {
	return "companies"
}

func (c *Company) IsTrialActive() bool {
	if c.TrialEndsAt == nil {
		return false
	}
	return time.Now().Before(*c.TrialEndsAt)
}
