package repository

import (
	"dvra-api/internal/shared/models"

	"gorm.io/gorm"
)

// PlanRepository interface defines plan data access methods
type PlanRepository struct {
	db *gorm.DB
}

func withPlanCatalogs(db *gorm.DB) *gorm.DB {
	return db.Preload("Currency").Preload("BillingCycle").Preload("SupportLevel")
}

// NewPlanRepository creates a new plan repository
func NewPlanRepository(db *gorm.DB) *PlanRepository {
	return &PlanRepository{db: db}
}

// Create creates a new plan
func (r *PlanRepository) Create(plan *models.Plan) (*models.Plan, error) {
	err := r.db.Create(plan).Error
	if err != nil {
		return nil, err
	}
	return r.FindByID(plan.ID)
}

// FindByID finds a plan by ID
func (r *PlanRepository) FindByID(id uint) (*models.Plan, error) {
	var plan models.Plan
	err := withPlanCatalogs(r.db).First(&plan, id).Error
	if err != nil {
		return nil, err
	}
	return &plan, nil
}

// FindBySlug finds a plan by slug
func (r *PlanRepository) FindBySlug(slug string) (*models.Plan, error) {
	var plan models.Plan
	err := withPlanCatalogs(r.db).Where("slug = ?", slug).First(&plan).Error
	if err != nil {
		return nil, err
	}
	return &plan, nil
}

// FindActiveBySlug finds an active plan by slug
func (r *PlanRepository) FindActiveBySlug(slug string) (*models.Plan, error) {
	var plan models.Plan
	err := withPlanCatalogs(r.db).Where("slug = ? AND is_active = ?", slug, true).First(&plan).Error
	if err != nil {
		return nil, err
	}
	return &plan, nil
}

// FindAll retrieves all plans
func (r *PlanRepository) FindAll() ([]models.Plan, error) {
	var plans []models.Plan
	err := withPlanCatalogs(r.db).Order("display_order ASC, name ASC").Find(&plans).Error
	if err != nil {
		return nil, err
	}
	return plans, nil
}

// FindActive retrieves all active plans
func (r *PlanRepository) FindActive() ([]models.Plan, error) {
	var plans []models.Plan
	err := withPlanCatalogs(r.db).Where("is_active = ?", true).
		Order("display_order ASC, name ASC").
		Find(&plans).Error
	if err != nil {
		return nil, err
	}
	return plans, nil
}

// FindPublic retrieves all public plans (for pricing page)
func (r *PlanRepository) FindPublic() ([]models.Plan, error) {
	var plans []models.Plan
	err := withPlanCatalogs(r.db).Where("is_active = ? AND is_public = ?", true, true).
		Order("display_order ASC, name ASC").
		Find(&plans).Error
	if err != nil {
		return nil, err
	}
	return plans, nil
}

// Update updates an existing plan
func (r *PlanRepository) Update(plan *models.Plan) (*models.Plan, error) {
	err := r.db.Save(plan).Error
	if err != nil {
		return nil, err
	}
	return r.FindByID(plan.ID)
}

// Delete soft deletes a plan
func (r *PlanRepository) Delete(id uint) error {
	return r.db.Delete(&models.Plan{}, id).Error
}

// ExistsBySlug checks if a plan with the given slug exists
func (r *PlanRepository) ExistsBySlug(slug string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Plan{}).Where("slug = ?", slug).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
