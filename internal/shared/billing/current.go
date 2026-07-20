package billing

import (
	"errors"

	"dvra-api/internal/app/models"

	"gorm.io/gorm"
)

var entitlementStatuses = []string{"trialing", "active", "past_due"}

func CurrentSubscription(db *gorm.DB, companyID uint) (*models.CompanySubscription, error) {
	var subscription models.CompanySubscription
	err := db.Preload("Plan").Preload("Status").Preload("Currency").Preload("BillingCycle").
		Joins("JOIN system_values subscription_status ON subscription_status.id = company_subscriptions.status_value_id").
		Where("company_subscriptions.company_id = ? AND subscription_status.key IN ?", companyID, entitlementStatuses).
		Order("company_subscriptions.started_at DESC, company_subscriptions.id DESC").First(&subscription).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &subscription, err
}

func CurrentPlan(db *gorm.DB, companyID uint) (*models.Plan, error) {
	subscription, err := CurrentSubscription(db, companyID)
	if err != nil || subscription == nil {
		return nil, err
	}
	return subscription.Plan, nil
}
