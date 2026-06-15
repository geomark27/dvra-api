package repositories

import (
	"dvra-api/internal/app/models"
	"dvra-api/internal/database"

	"gorm.io/gorm"
)

// StaffingClientRepository define el contrato del repositorio de clientes finales
type StaffingClientRepository interface {
	GetByCompanyID(companyID uint, status string) ([]models.StaffingClient, error)
	GetByID(id uint) (*models.StaffingClient, error)
	ExistsBySlug(companyID uint, slug string, excludeID uint) (bool, error)
	Create(client *models.StaffingClient) (*models.StaffingClient, error)
	Update(client *models.StaffingClient) (*models.StaffingClient, error)
	Delete(id uint) error
}

// staffingClientRepository es la implementación con GORM
type staffingClientRepository struct{}

// NewStaffingClientRepository crea una nueva instancia de StaffingClientRepository
func NewStaffingClientRepository() StaffingClientRepository {
	return &staffingClientRepository{}
}

func (r *staffingClientRepository) GetByCompanyID(companyID uint, status string) ([]models.StaffingClient, error) {
	var clients []models.StaffingClient
	query := database.DB.Where("company_id = ?", companyID)
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if err := query.Order("name ASC").Find(&clients).Error; err != nil {
		return nil, err
	}
	return clients, nil
}

func (r *staffingClientRepository) GetByID(id uint) (*models.StaffingClient, error) {
	var client models.StaffingClient
	if err := database.DB.First(&client, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &client, nil
}

// ExistsBySlug verifica si ya existe un cliente con ese slug en la empresa.
// excludeID (>0) permite ignorar el propio registro al actualizar.
func (r *staffingClientRepository) ExistsBySlug(companyID uint, slug string, excludeID uint) (bool, error) {
	var count int64
	query := database.DB.Model(&models.StaffingClient{}).Where("company_id = ? AND slug = ?", companyID, slug)
	if excludeID > 0 {
		query = query.Where("id <> ?", excludeID)
	}
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *staffingClientRepository) Create(client *models.StaffingClient) (*models.StaffingClient, error) {
	if err := database.DB.Create(client).Error; err != nil {
		return nil, err
	}
	return client, nil
}

func (r *staffingClientRepository) Update(client *models.StaffingClient) (*models.StaffingClient, error) {
	if err := database.DB.Save(client).Error; err != nil {
		return nil, err
	}
	return client, nil
}

func (r *staffingClientRepository) Delete(id uint) error {
	return database.DB.Delete(&models.StaffingClient{}, id).Error
}
