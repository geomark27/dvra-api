package repositories

import (
	"dvra-api/internal/app/models"

	"gorm.io/gorm"
)

// PlanRepository interface defines plan data access methods
type PlanRepository interface {
	Create(plan *models.Plan) (*models.Plan, error)
	FindByID(id uint) (*models.Plan, error)
	FindBySlug(slug string) (*models.Plan, error)
	FindActiveBySlug(slug string) (*models.Plan, error)
	FindAll() ([]models.Plan, error)
	FindActive() ([]models.Plan, error)
	FindPublic() ([]models.Plan, error)
	Update(plan *models.Plan) (*models.Plan, error)
	Delete(id uint) error
	ExistsBySlug(slug string) (bool, error)
}

type planRepository struct {
	db *gorm.DB
}

// NewPlanRepository creates a new plan repository
func NewPlanRepository(db *gorm.DB) PlanRepository {
	return &planRepository{db: db}
}

// Create creates a new plan
func (r *planRepository) Create(plan *models.Plan) (*models.Plan, error) {
	err := r.db.Create(plan).Error
	if err != nil {
		return nil, err
	}
	return plan, nil
}

// FindByID finds a plan by ID
func (r *planRepository) FindByID(id uint) (*models.Plan, error) {
	var plan models.Plan
	err := r.db.First(&plan, id).Error
	if err != nil {
		return nil, err
	}
	return &plan, nil
}

// FindBySlug finds a plan by slug
func (r *planRepository) FindBySlug(slug string) (*models.Plan, error) {
	var plan models.Plan
	err := r.db.Where("slug = ?", slug).First(&plan).Error
	if err != nil {
		return nil, err
	}
	return &plan, nil
}

// FindActiveBySlug finds an active plan by slug
func (r *planRepository) FindActiveBySlug(slug string) (*models.Plan, error) {
	var plan models.Plan
	err := r.db.Where("slug = ? AND is_active = ?", slug, true).First(&plan).Error
	if err != nil {
		return nil, err
	}
	return &plan, nil
}

// FindAll retrieves all plans
func (r *planRepository) FindAll() ([]models.Plan, error) {
	var plans []models.Plan
	err := r.db.Order("display_order ASC, name ASC").Find(&plans).Error
	if err != nil {
		return nil, err
	}
	return plans, nil
}

// FindActive retrieves all active plans
func (r *planRepository) FindActive() ([]models.Plan, error) {
	var plans []models.Plan
	err := r.db.Where("is_active = ?", true).
		Order("display_order ASC, name ASC").
		Find(&plans).Error
	if err != nil {
		return nil, err
	}
	return plans, nil
}

// FindPublic retrieves all public plans (for pricing page)
func (r *planRepository) FindPublic() ([]models.Plan, error) {
	var plans []models.Plan
	err := r.db.Where("is_active = ? AND is_public = ?", true, true).
		Order("display_order ASC, name ASC").
		Find(&plans).Error
	if err != nil {
		return nil, err
	}
	return plans, nil
}

// Update updates an existing plan
func (r *planRepository) Update(plan *models.Plan) (*models.Plan, error) {
	err := r.db.Save(plan).Error
	if err != nil {
		return nil, err
	}
	return plan, nil
}

// Delete soft deletes a plan
func (r *planRepository) Delete(id uint) error {
	return r.db.Delete(&models.Plan{}, id).Error
}

// ExistsBySlug checks if a plan with the given slug exists
func (r *planRepository) ExistsBySlug(slug string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Plan{}).Where("slug = ?", slug).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
