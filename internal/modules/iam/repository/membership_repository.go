package repository

import (
	"dvra-api/internal/shared/models"

	"gorm.io/gorm"
)

// MembershipRepository define el contrato del repositorio de memberships
// membershipRepository es la implementación con GORM
type MembershipRepository struct{ db *gorm.DB }

// NewMembershipRepository crea una nueva instancia de MembershipRepository
func NewMembershipRepository(db *gorm.DB) *MembershipRepository {
	return &MembershipRepository{db: db}
}

func (r *MembershipRepository) GetAll() ([]models.Membership, error) {
	var memberships []models.Membership
	if err := r.db.Preload("User").Preload("Company").Preload("TenantRole").Find(&memberships).Error; err != nil {
		return nil, err
	}
	return memberships, nil
}

func (r *MembershipRepository) GetByID(id uint) (*models.Membership, error) {
	var membership models.Membership
	if err := r.db.Preload("User").Preload("Company").Preload("TenantRole").First(&membership, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &membership, nil
}

func (r *MembershipRepository) GetByUserID(userID uint) ([]models.Membership, error) {
	var memberships []models.Membership
	if err := r.db.Where("user_id = ?", userID).Preload("Company").Preload("TenantRole").Find(&memberships).Error; err != nil {
		return nil, err
	}
	return memberships, nil
}

func (r *MembershipRepository) GetByCompanyID(companyID uint) ([]models.Membership, error) {
	var memberships []models.Membership
	if err := r.db.Where("company_id = ?", companyID).Preload("User").Preload("TenantRole").Find(&memberships).Error; err != nil {
		return nil, err
	}
	return memberships, nil
}

func (r *MembershipRepository) GetByUserAndCompany(userID uint, companyID uint) (*models.Membership, error) {
	var membership models.Membership
	if err := r.db.Where("user_id = ? AND company_id = ?", userID, companyID).Preload("TenantRole").First(&membership).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &membership, nil
}

func (r *MembershipRepository) Create(membership *models.Membership) (*models.Membership, error) {
	if err := r.db.Create(membership).Error; err != nil {
		return nil, err
	}
	return membership, nil
}

func (r *MembershipRepository) Update(membership *models.Membership) (*models.Membership, error) {
	if err := r.db.Save(membership).Error; err != nil {
		return nil, err
	}
	return membership, nil
}

func (r *MembershipRepository) Delete(id uint) error {
	return r.db.Delete(&models.Membership{}, id).Error
}
