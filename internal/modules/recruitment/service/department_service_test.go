package service

import (
	"errors"
	"testing"

	"dvra-api/internal/shared/apperr"
	"dvra-api/internal/shared/models"
)

type departmentRepoFake struct {
	departments    map[uint]*models.Department
	managerAllowed bool
	hasChildren    bool
	created        *models.Department
	updated        *models.Department
}

func (f *departmentRepoFake) List(uint, bool) ([]models.Department, error) { return nil, nil }
func (f *departmentRepoFake) GetByID(id, companyID uint) (*models.Department, error) {
	department := f.departments[id]
	if department == nil || department.CompanyID != companyID {
		return nil, nil
	}
	copy := *department
	return &copy, nil
}
func (f *departmentRepoFake) NameExists(uint, *uint, string, uint) (bool, error) {
	return false, nil
}
func (f *departmentRepoFake) CodeExists(uint, string, uint) (bool, error) { return false, nil }
func (f *departmentRepoFake) UserIsActiveMember(uint, uint) (bool, error) {
	return f.managerAllowed, nil
}
func (f *departmentRepoFake) HasActiveChildren(uint, uint) (bool, error) {
	return f.hasChildren, nil
}
func (f *departmentRepoFake) Create(department *models.Department) (*models.Department, error) {
	f.created = department
	return department, nil
}
func (f *departmentRepoFake) Update(department *models.Department) (*models.Department, error) {
	f.updated = department
	return department, nil
}

func TestCreateDepartmentRejectsParentFromAnotherTenant(t *testing.T) {
	parentID := uint(2)
	repo := &departmentRepoFake{departments: map[uint]*models.Department{
		parentID: {BaseModel: models.BaseModel{ID: parentID}, CompanyID: 99, Name: "Other", IsActive: true},
	}}
	svc := NewDepartmentService(repo)
	_, err := svc.CreateDepartment(7, CreateDepartmentDTO{Name: "Backend", ParentDepartmentID: &parentID})
	assertDepartmentErrorCode(t, err, 400)
}

func TestUpdateDepartmentRejectsHierarchyCycle(t *testing.T) {
	rootID, childID, grandchildID := uint(1), uint(2), uint(3)
	repo := &departmentRepoFake{departments: map[uint]*models.Department{
		rootID:       {BaseModel: models.BaseModel{ID: rootID}, CompanyID: 7, Name: "Tech", IsActive: true},
		childID:      {BaseModel: models.BaseModel{ID: childID}, CompanyID: 7, ParentDepartmentID: &rootID, Name: "Engineering", IsActive: true},
		grandchildID: {BaseModel: models.BaseModel{ID: grandchildID}, CompanyID: 7, ParentDepartmentID: &childID, Name: "Backend", IsActive: true},
	}}
	svc := NewDepartmentService(repo)
	_, err := svc.UpdateDepartment(rootID, 7, UpdateDepartmentDTO{ParentDepartmentID: &grandchildID})
	assertDepartmentErrorCode(t, err, 400)
}

func TestCreateDepartmentRequiresManagerFromTenant(t *testing.T) {
	managerID := uint(8)
	repo := &departmentRepoFake{departments: make(map[uint]*models.Department)}
	svc := NewDepartmentService(repo)
	_, err := svc.CreateDepartment(7, CreateDepartmentDTO{Name: "People", ManagerUserID: &managerID})
	assertDepartmentErrorCode(t, err, 400)
}

func TestDeactivateDepartmentRejectsActiveChildren(t *testing.T) {
	repo := &departmentRepoFake{
		departments: map[uint]*models.Department{1: {BaseModel: models.BaseModel{ID: 1}, CompanyID: 7, Name: "Tech", IsActive: true}},
		hasChildren: true,
	}
	svc := NewDepartmentService(repo)
	_, err := svc.DeactivateDepartment(1, 7)
	assertDepartmentErrorCode(t, err, 409)
}

func TestCreateRootDepartmentNormalizesOptionalIDs(t *testing.T) {
	zero := uint(0)
	repo := &departmentRepoFake{departments: make(map[uint]*models.Department)}
	svc := NewDepartmentService(repo)
	created, err := svc.CreateDepartment(7, CreateDepartmentDTO{
		Name: "  Technology  ", Code: " TECH ", ParentDepartmentID: &zero, ManagerUserID: &zero,
	})
	if err != nil {
		t.Fatal(err)
	}
	if created.ParentDepartmentID != nil || created.ManagerUserID != nil || created.Name != "Technology" || created.Code != "TECH" {
		t.Fatalf("unexpected normalized department: %#v", created)
	}
}

func TestUpdateAndDeactivateAppliesAllChanges(t *testing.T) {
	name := "People Operations"
	active := false
	repo := &departmentRepoFake{departments: map[uint]*models.Department{
		1: {BaseModel: models.BaseModel{ID: 1}, CompanyID: 7, Name: "People", IsActive: true},
	}}
	svc := NewDepartmentService(repo)
	updated, err := svc.UpdateDepartment(1, 7, UpdateDepartmentDTO{Name: &name, IsActive: &active})
	if err != nil {
		t.Fatal(err)
	}
	if updated.Name != name || updated.IsActive {
		t.Fatalf("update was not applied before deactivation: %#v", updated)
	}
}

func assertDepartmentErrorCode(t *testing.T, err error, code int) {
	t.Helper()
	var appErr *apperr.AppError
	if !errors.As(err, &appErr) || appErr.Code != code {
		t.Fatalf("expected app error %d, got %v", code, err)
	}
}
