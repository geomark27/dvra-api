package service

import (
	"strings"

	"dvra-api/internal/shared/apperr"
	"dvra-api/internal/shared/models"
)

type DepartmentService interface {
	ListDepartments(companyID uint, includeInactive bool) ([]models.Department, error)
	GetDepartment(id, companyID uint) (*models.Department, error)
	CreateDepartment(companyID uint, dto CreateDepartmentDTO) (*models.Department, error)
	UpdateDepartment(id, companyID uint, dto UpdateDepartmentDTO) (*models.Department, error)
	DeactivateDepartment(id, companyID uint) (*models.Department, error)
}

type departmentService struct {
	repo DepartmentRepository
}

func NewDepartmentService(repo DepartmentRepository) DepartmentService {
	return &departmentService{repo: repo}
}

func (s *departmentService) ListDepartments(companyID uint, includeInactive bool) ([]models.Department, error) {
	return s.repo.List(companyID, includeInactive)
}

func (s *departmentService) GetDepartment(id, companyID uint) (*models.Department, error) {
	department, err := s.repo.GetByID(id, companyID)
	if err != nil {
		return nil, err
	}
	if department == nil {
		return nil, apperr.NotFound("department not found")
	}
	return department, nil
}

func (s *departmentService) CreateDepartment(companyID uint, dto CreateDepartmentDTO) (*models.Department, error) {
	parentID := normalizeOptionalID(dto.ParentDepartmentID)
	managerID := normalizeOptionalID(dto.ManagerUserID)
	if err := s.validateParent(companyID, 0, parentID); err != nil {
		return nil, err
	}
	if err := s.validateManager(companyID, managerID); err != nil {
		return nil, err
	}
	if err := s.validateUniqueness(companyID, parentID, dto.Name, dto.Code, 0); err != nil {
		return nil, err
	}
	department := &models.Department{
		CompanyID: companyID, ParentDepartmentID: parentID, Name: strings.TrimSpace(dto.Name),
		Code: strings.TrimSpace(dto.Code), Description: strings.TrimSpace(dto.Description),
		ManagerUserID: managerID, IsActive: true,
	}
	return s.repo.Create(department)
}

func (s *departmentService) UpdateDepartment(id, companyID uint, dto UpdateDepartmentDTO) (*models.Department, error) {
	department, err := s.GetDepartment(id, companyID)
	if err != nil {
		return nil, err
	}
	parentID := department.ParentDepartmentID
	if dto.ParentDepartmentID != nil {
		parentID = normalizeOptionalID(dto.ParentDepartmentID)
		if err := s.validateParent(companyID, id, parentID); err != nil {
			return nil, err
		}
	}
	managerID := department.ManagerUserID
	if dto.ManagerUserID != nil {
		managerID = normalizeOptionalID(dto.ManagerUserID)
		if err := s.validateManager(companyID, managerID); err != nil {
			return nil, err
		}
	}
	name, code := department.Name, department.Code
	if dto.Name != nil {
		name = strings.TrimSpace(*dto.Name)
	}
	if dto.Code != nil {
		code = strings.TrimSpace(*dto.Code)
	}
	if err := s.validateUniqueness(companyID, parentID, name, code, id); err != nil {
		return nil, err
	}
	department.ParentDepartmentID, department.ManagerUserID = parentID, managerID
	department.Name, department.Code = name, code
	if dto.Description != nil {
		department.Description = strings.TrimSpace(*dto.Description)
	}
	if dto.IsActive != nil {
		if !*dto.IsActive {
			hasChildren, childrenErr := s.repo.HasActiveChildren(id, companyID)
			if childrenErr != nil {
				return nil, childrenErr
			}
			if hasChildren {
				return nil, apperr.Conflict("deactivate child departments first")
			}
			department.IsActive = false
		} else if parentID != nil {
			parent, parentErr := s.repo.GetByID(*parentID, companyID)
			if parentErr != nil {
				return nil, parentErr
			}
			if parent == nil || !parent.IsActive {
				return nil, apperr.BadRequest("an active department requires an active parent")
			}
			department.IsActive = true
		} else {
			department.IsActive = true
		}
	}
	return s.repo.Update(department)
}

func (s *departmentService) DeactivateDepartment(id, companyID uint) (*models.Department, error) {
	department, err := s.GetDepartment(id, companyID)
	if err != nil {
		return nil, err
	}
	hasChildren, err := s.repo.HasActiveChildren(id, companyID)
	if err != nil {
		return nil, err
	}
	if hasChildren {
		return nil, apperr.Conflict("deactivate child departments first")
	}
	department.IsActive = false
	return s.repo.Update(department)
}

func (s *departmentService) validateParent(companyID, departmentID uint, parentID *uint) error {
	visited := make(map[uint]bool)
	for parentID != nil {
		if *parentID == departmentID && departmentID > 0 {
			return apperr.BadRequest("a department cannot be its own ancestor")
		}
		if visited[*parentID] {
			return apperr.BadRequest("department hierarchy contains a cycle")
		}
		visited[*parentID] = true
		parent, err := s.repo.GetByID(*parentID, companyID)
		if err != nil {
			return err
		}
		if parent == nil {
			return apperr.BadRequest("parent department does not belong to this company")
		}
		if !parent.IsActive {
			return apperr.BadRequest("parent department is inactive")
		}
		parentID = parent.ParentDepartmentID
	}
	return nil
}

func (s *departmentService) validateManager(companyID uint, managerID *uint) error {
	if managerID == nil {
		return nil
	}
	belongs, err := s.repo.UserIsActiveMember(*managerID, companyID)
	if err != nil {
		return err
	}
	if !belongs {
		return apperr.BadRequest("manager must be an active member of this company")
	}
	return nil
}

func (s *departmentService) validateUniqueness(companyID uint, parentID *uint, name, code string, excludeID uint) error {
	exists, err := s.repo.NameExists(companyID, parentID, name, excludeID)
	if err != nil {
		return err
	}
	if exists {
		return apperr.Conflict("a department with this name already exists at the same level")
	}
	exists, err = s.repo.CodeExists(companyID, code, excludeID)
	if err != nil {
		return err
	}
	if exists {
		return apperr.Conflict("department code already exists in this company")
	}
	return nil
}

func normalizeOptionalID(id *uint) *uint {
	if id == nil || *id == 0 {
		return nil
	}
	value := *id
	return &value
}
