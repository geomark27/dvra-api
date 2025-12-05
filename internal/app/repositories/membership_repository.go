package repositories

import (
	"dvra-api/internal/app/models"
	"dvra-api/internal/database"

	"gorm.io/gorm"
)

// MembershipRepository define el contrato del repositorio de memberships
type MembershipRepository interface {
	GetAll() ([]models.Membership, error)
	GetByID(id uint) (*models.Membership, error)
	GetByUserID(userID uint) ([]models.Membership, error)
	GetByCompanyID(companyID uint) ([]models.Membership, error)
	GetByUserAndCompany(userID uint, companyID uint) (*models.Membership, error)
	Create(membership *models.Membership) (*models.Membership, error)
	Update(membership *models.Membership) (*models.Membership, error)
	Delete(id uint) error
}

// membershipRepository es la implementación con GORM
type membershipRepository struct{}

// NewMembershipRepository crea una nueva instancia de MembershipRepository
func NewMembershipRepository() MembershipRepository {
	return &membershipRepository{}
}

func (r *membershipRepository) GetAll() ([]models.Membership, error) {
	var memberships []models.Membership
	if err := database.DB.Preload("User").Preload("Company").Find(&memberships).Error; err != nil {
		return nil, err
	}
	return memberships, nil
}

func (r *membershipRepository) GetByID(id uint) (*models.Membership, error) {
	var membership models.Membership
	if err := database.DB.Preload("User").Preload("Company").First(&membership, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &membership, nil
}

func (r *membershipRepository) GetByUserID(userID uint) ([]models.Membership, error) {
	var memberships []models.Membership
	if err := database.DB.Where("user_id = ?", userID).Preload("Company").Find(&memberships).Error; err != nil {
		return nil, err
	}
	return memberships, nil
}

func (r *membershipRepository) GetByCompanyID(companyID uint) ([]models.Membership, error) {
	var memberships []models.Membership
	if err := database.DB.Where("company_id = ?", companyID).Preload("User").Find(&memberships).Error; err != nil {
		return nil, err
	}
	return memberships, nil
}

func (r *membershipRepository) GetByUserAndCompany(userID uint, companyID uint) (*models.Membership, error) {
	var membership models.Membership
	if err := database.DB.Where("user_id = ? AND company_id = ?", userID, companyID).First(&membership).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &membership, nil
}

func (r *membershipRepository) Create(membership *models.Membership) (*models.Membership, error) {
	if err := database.DB.Create(membership).Error; err != nil {
		return nil, err
	}
	return membership, nil
}

func (r *membershipRepository) Update(membership *models.Membership) (*models.Membership, error) {
	if err := database.DB.Save(membership).Error; err != nil {
		return nil, err
	}
	return membership, nil
}

func (r *membershipRepository) Delete(id uint) error {
	return database.DB.Delete(&models.Membership{}, id).Error
}
