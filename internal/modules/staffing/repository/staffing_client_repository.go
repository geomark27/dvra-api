// Package repository implementa los puertos de persistencia del módulo staffing
// con GORM. Recibe *gorm.DB inyectado (no usa el global database.DB), lo que lo
// hace testeable con SQLite en memoria.
package repository

import (
	"dvra-api/internal/app/models"
	"dvra-api/internal/modules/staffing/domain"

	"gorm.io/gorm"
)

type staffingClientRepository struct {
	db *gorm.DB
}

// NewStaffingClientRepository devuelve la implementación del puerto.
func NewStaffingClientRepository(db *gorm.DB) domain.StaffingClientRepository {
	return &staffingClientRepository{db: db}
}

func (r *staffingClientRepository) GetByCompanyID(companyID uint, status string) ([]models.StaffingClient, error) {
	var clients []models.StaffingClient
	query := r.db.Where("company_id = ?", companyID)
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
	if err := r.db.First(&client, id).Error; err != nil {
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
	query := r.db.Model(&models.StaffingClient{}).Where("company_id = ? AND slug = ?", companyID, slug)
	if excludeID > 0 {
		query = query.Where("id <> ?", excludeID)
	}
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *staffingClientRepository) Create(client *models.StaffingClient) (*models.StaffingClient, error) {
	if err := r.db.Create(client).Error; err != nil {
		return nil, err
	}
	return client, nil
}

func (r *staffingClientRepository) Update(client *models.StaffingClient) (*models.StaffingClient, error) {
	if err := r.db.Save(client).Error; err != nil {
		return nil, err
	}
	return client, nil
}

func (r *staffingClientRepository) Delete(id uint) error {
	return r.db.Delete(&models.StaffingClient{}, id).Error
}
