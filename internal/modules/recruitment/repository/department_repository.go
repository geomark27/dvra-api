package repository

import (
	"strings"

	"dvra-api/internal/shared/models"

	"gorm.io/gorm"
)

type DepartmentRepository struct{ db *gorm.DB }

func NewDepartmentRepository(db *gorm.DB) *DepartmentRepository {
	return &DepartmentRepository{db: db}
}

func (r *DepartmentRepository) List(companyID uint, includeInactive bool) ([]models.Department, error) {
	var departments []models.Department
	query := r.db.Where("company_id = ?", companyID).Preload("Manager").Order("name ASC")
	if !includeInactive {
		query = query.Where("is_active = ?", true)
	}
	if err := query.Find(&departments).Error; err != nil {
		return nil, err
	}
	ids := make([]uint, len(departments))
	for i := range departments {
		ids[i] = departments[i].ID
	}
	paths, err := findDepartmentPaths(r.db, ids)
	if err != nil {
		return nil, err
	}
	attachDepartmentPaths(departments, paths)
	return departments, nil
}

func (r *DepartmentRepository) GetByID(id, companyID uint) (*models.Department, error) {
	var department models.Department
	if err := r.db.Preload("Manager").Where("id = ? AND company_id = ?", id, companyID).First(&department).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	paths, err := findDepartmentPaths(r.db, []uint{department.ID})
	if err != nil {
		return nil, err
	}
	department.Path = paths[department.ID]
	return &department, nil
}

func (r *DepartmentRepository) NameExists(companyID uint, parentID *uint, name string, excludeID uint) (bool, error) {
	query := r.db.Model(&models.Department{}).
		Where("company_id = ? AND LOWER(name) = ?", companyID, strings.ToLower(strings.TrimSpace(name)))
	if parentID == nil || *parentID == 0 {
		query = query.Where("parent_department_id IS NULL")
	} else {
		query = query.Where("parent_department_id = ?", *parentID)
	}
	if excludeID > 0 {
		query = query.Where("id <> ?", excludeID)
	}
	var count int64
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *DepartmentRepository) CodeExists(companyID uint, code string, excludeID uint) (bool, error) {
	code = strings.TrimSpace(code)
	if code == "" {
		return false, nil
	}
	query := r.db.Model(&models.Department{}).Where("company_id = ? AND LOWER(code) = ?", companyID, strings.ToLower(code))
	if excludeID > 0 {
		query = query.Where("id <> ?", excludeID)
	}
	var count int64
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *DepartmentRepository) UserIsActiveMember(userID, companyID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.Membership{}).
		Where("user_id = ? AND company_id = ? AND status = ?", userID, companyID, "active").Count(&count).Error
	return count > 0, err
}

func (r *DepartmentRepository) HasActiveChildren(id, companyID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.Department{}).
		Where("parent_department_id = ? AND company_id = ? AND is_active = ?", id, companyID, true).Count(&count).Error
	return count > 0, err
}

func (r *DepartmentRepository) Create(department *models.Department) (*models.Department, error) {
	if err := r.db.Create(department).Error; err != nil {
		return nil, err
	}
	return r.GetByID(department.ID, department.CompanyID)
}

func (r *DepartmentRepository) Update(department *models.Department) (*models.Department, error) {
	if err := r.db.Save(department).Error; err != nil {
		return nil, err
	}
	return r.GetByID(department.ID, department.CompanyID)
}
