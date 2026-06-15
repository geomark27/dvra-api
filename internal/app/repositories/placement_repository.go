package repositories

import (
	"dvra-api/internal/app/dtos"
	"dvra-api/internal/app/models"
	"dvra-api/internal/database"

	"gorm.io/gorm"
)

// PlacementRepository define el contrato del repositorio de colocaciones
type PlacementRepository interface {
	GetByCompanyID(companyID uint, filters dtos.PlacementFilters) ([]models.Placement, error)
	GetByID(id uint) (*models.Placement, error)
	Create(placement *models.Placement) (*models.Placement, error)
	Update(placement *models.Placement) (*models.Placement, error)
	Delete(id uint) error
}

// placementRepository es la implementación con GORM
type placementRepository struct{}

// NewPlacementRepository crea una nueva instancia de PlacementRepository
func NewPlacementRepository() PlacementRepository {
	return &placementRepository{}
}

func (r *placementRepository) GetByCompanyID(companyID uint, filters dtos.PlacementFilters) ([]models.Placement, error) {
	var placements []models.Placement
	query := database.DB.
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
	if err := database.DB.
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

func (r *placementRepository) Create(placement *models.Placement) (*models.Placement, error) {
	if err := database.DB.Create(placement).Error; err != nil {
		return nil, err
	}
	return placement, nil
}

func (r *placementRepository) Update(placement *models.Placement) (*models.Placement, error) {
	if err := database.DB.Save(placement).Error; err != nil {
		return nil, err
	}
	return placement, nil
}

func (r *placementRepository) Delete(id uint) error {
	return database.DB.Delete(&models.Placement{}, id).Error
}
