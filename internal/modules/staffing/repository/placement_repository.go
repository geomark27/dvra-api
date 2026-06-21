package repository

import (
	"dvra-api/internal/app/dtos"
	"dvra-api/internal/app/models"
	"dvra-api/internal/modules/staffing/domain"

	"gorm.io/gorm"
)

type placementRepository struct {
	db *gorm.DB
}

// NewPlacementRepository devuelve la implementación del puerto.
func NewPlacementRepository(db *gorm.DB) domain.PlacementRepository {
	return &placementRepository{db: db}
}

func (r *placementRepository) GetByCompanyID(companyID uint, filters dtos.PlacementFilters) ([]models.Placement, error) {
	var placements []models.Placement
	query := r.db.
		Preload("StaffingClient").
		Preload("Candidate").
		Where("company_id = ?", companyID)

	if filters.Status != "" {
		query = query.Where("status = ?", filters.Status)
	}
	if filters.StaffingClientID != nil {
		query = query.Where("staffing_client_id = ?", *filters.StaffingClientID)
	}
	if filters.CandidateID != nil {
		query = query.Where("candidate_id = ?", *filters.CandidateID)
	}

	if err := query.Order("created_at DESC").Find(&placements).Error; err != nil {
		return nil, err
	}
	return placements, nil
}

func (r *placementRepository) GetByID(id uint) (*models.Placement, error) {
	var placement models.Placement
	if err := r.db.
		Preload("StaffingClient").
		Preload("Candidate").
		Preload("Job").
		Preload("Application").
		First(&placement, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &placement, nil
}

func (r *placementRepository) ExistsByApplicationID(applicationID uint) (bool, error) {
	var count int64
	if err := r.db.Model(&models.Placement{}).
		Where("application_id = ?", applicationID).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *placementRepository) Create(placement *models.Placement) (*models.Placement, error) {
	if err := r.db.Create(placement).Error; err != nil {
		return nil, err
	}
	return placement, nil
}

func (r *placementRepository) Update(placement *models.Placement) (*models.Placement, error) {
	if err := r.db.Save(placement).Error; err != nil {
		return nil, err
	}
	return placement, nil
}

func (r *placementRepository) Delete(id uint) error {
	return r.db.Delete(&models.Placement{}, id).Error
}
