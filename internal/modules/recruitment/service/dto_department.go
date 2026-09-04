package service

import (
	"time"

	"dvra-api/internal/shared/models"
)

type CreateDepartmentDTO struct {
	Name               string `json:"name" binding:"required,min=2,max=150"`
	Code               string `json:"code,omitempty" binding:"omitempty,max=50"`
	Description        string `json:"description,omitempty" binding:"omitempty,max=2000"`
	ParentDepartmentID *uint  `json:"parent_department_id,omitempty"`
	ManagerUserID      *uint  `json:"manager_user_id,omitempty"`
}

type UpdateDepartmentDTO struct {
	Name               *string `json:"name,omitempty" binding:"omitempty,min=2,max=150"`
	Code               *string `json:"code,omitempty" binding:"omitempty,max=50"`
	Description        *string `json:"description,omitempty" binding:"omitempty,max=2000"`
	ParentDepartmentID *uint   `json:"parent_department_id,omitempty"`
	ManagerUserID      *uint   `json:"manager_user_id,omitempty"`
	IsActive           *bool   `json:"is_active,omitempty"`
}

type DepartmentManagerResponse struct {
	ID        uint   `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

type DepartmentResponse struct {
	ID                 uint                        `json:"id"`
	CreatedAt          time.Time                   `json:"created_at"`
	UpdatedAt          time.Time                   `json:"updated_at"`
	CompanyID          uint                        `json:"company_id"`
	ParentDepartmentID *uint                       `json:"parent_department_id,omitempty"`
	Name               string                      `json:"name"`
	Code               string                      `json:"code,omitempty"`
	Description        string                      `json:"description,omitempty"`
	ManagerUserID      *uint                       `json:"manager_user_id,omitempty"`
	Manager            *DepartmentManagerResponse  `json:"manager,omitempty"`
	IsActive           bool                        `json:"is_active"`
	Path               []models.DepartmentPathItem `json:"path"`
	Children           []*DepartmentResponse       `json:"children,omitempty"`
}

func ToDepartmentResponse(department *models.Department) *DepartmentResponse {
	response := &DepartmentResponse{
		ID: department.ID, CreatedAt: department.CreatedAt, UpdatedAt: department.UpdatedAt,
		CompanyID: department.CompanyID, ParentDepartmentID: department.ParentDepartmentID,
		Name: department.Name, Code: department.Code, Description: department.Description,
		ManagerUserID: department.ManagerUserID, IsActive: department.IsActive,
		Path: department.Path,
	}
	if len(response.Path) == 0 {
		response.Path = []models.DepartmentPathItem{{ID: department.ID, Name: department.Name, Code: department.Code}}
	}
	if department.Manager != nil {
		response.Manager = &DepartmentManagerResponse{
			ID: department.Manager.ID, FirstName: department.Manager.FirstName,
			LastName: department.Manager.LastName, Email: department.Manager.Email,
		}
	}
	return response
}

func ToDepartmentResponseList(departments []models.Department) []*DepartmentResponse {
	result := make([]*DepartmentResponse, 0, len(departments))
	for i := range departments {
		result = append(result, ToDepartmentResponse(&departments[i]))
	}
	return result
}

func ToDepartmentTree(departments []models.Department) []*DepartmentResponse {
	nodes := make(map[uint]*DepartmentResponse, len(departments))
	for i := range departments {
		nodes[departments[i].ID] = ToDepartmentResponse(&departments[i])
	}
	roots := make([]*DepartmentResponse, 0)
	for i := range departments {
		node := nodes[departments[i].ID]
		if departments[i].ParentDepartmentID == nil {
			roots = append(roots, node)
			continue
		}
		parent, ok := nodes[*departments[i].ParentDepartmentID]
		if !ok {
			roots = append(roots, node)
			continue
		}
		parent.Children = append(parent.Children, node)
	}
	return roots
}
