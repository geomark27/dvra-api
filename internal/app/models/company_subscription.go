package models

import "time"

type CompanySubscription struct {
	BaseModel
	CompanyID              uint         `gorm:"not null;index" json:"company_id"`
	PlanID                 uint         `gorm:"not null;index" json:"plan_id"`
	StatusValueID          uint         `gorm:"not null;index" json:"status_value_id"`
	CurrencyValueID        uint         `gorm:"not null;index" json:"currency_value_id"`
	BillingCycleValueID    uint         `gorm:"not null;index" json:"billing_cycle_value_id"`
	PriceSnapshot          float64      `gorm:"type:decimal(10,2);not null;default:0" json:"price_snapshot"`
	StartedAt              time.Time    `gorm:"not null" json:"started_at"`
	CurrentPeriodStart     time.Time    `gorm:"not null" json:"current_period_start"`
	CurrentPeriodEnd       *time.Time   `json:"current_period_end,omitempty"`
	TrialEndsAt            *time.Time   `json:"trial_ends_at,omitempty"`
	EndedAt                *time.Time   `json:"ended_at,omitempty"`
	CancelAtPeriodEnd      bool         `gorm:"not null;default:false" json:"cancel_at_period_end"`
	ExternalCustomerID     string       `gorm:"type:varchar(255)" json:"external_customer_id,omitempty"`
	ExternalSubscriptionID string       `gorm:"type:varchar(255)" json:"external_subscription_id,omitempty"`
	Company                *Company     `gorm:"foreignKey:CompanyID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"company,omitempty"`
	Plan                   *Plan        `gorm:"foreignKey:PlanID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"plan,omitempty"`
	Status                 *SystemValue `gorm:"foreignKey:StatusValueID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"status,omitempty"`
	Currency               *SystemValue `gorm:"foreignKey:CurrencyValueID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"currency,omitempty"`
	BillingCycle           *SystemValue `gorm:"foreignKey:BillingCycleValueID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT" json:"billing_cycle,omitempty"`
}

func (CompanySubscription) TableName() string { return "company_subscriptions" }
