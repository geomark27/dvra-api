package dtos

import (
	"testing"

	"dvra-api/internal/app/models"
)

func TestDepartmentResponsesExposePathWithSafeFallback(t *testing.T) {
	department := &models.Department{
		BaseModel: models.BaseModel{ID: 2}, Name: "Backend", Code: "TECH-BE",
		Path: []models.DepartmentPathItem{
			{ID: 1, Name: "Technology", Code: "TECH"},
			{ID: 2, Name: "Backend", Code: "TECH-BE"},
		},
	}
	response := ToDepartmentResponse(department)
	if len(response.Path) != 2 || response.Path[0].Name != "Technology" || response.Path[1].Name != "Backend" {
		t.Fatalf("unexpected department path: %#v", response.Path)
	}

	department.Path = nil
	response = ToDepartmentResponse(department)
	if len(response.Path) != 1 || response.Path[0].ID != department.ID {
		t.Fatalf("missing fallback path: %#v", response.Path)
	}
}

func TestPublicJobResponseIncludesOnlyPublicDepartmentData(t *testing.T) {
	job := &models.Job{Department: &models.Department{
		BaseModel: models.BaseModel{ID: 2}, Name: "Backend", Code: "TECH-BE",
		Manager: &models.User{Email: "private@example.com"},
		Path:    []models.DepartmentPathItem{{ID: 1, Name: "Technology"}, {ID: 2, Name: "Backend"}},
	}}
	response := ToPublicJobResponse(job)
	if response.Department == nil || len(response.Department.Path) != 2 {
		t.Fatalf("public department path missing: %#v", response.Department)
	}
}
