package service

import (
	"testing"

	"dvra-api/internal/shared/models"
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
