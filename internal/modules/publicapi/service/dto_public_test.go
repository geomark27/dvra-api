package service

import (
	"testing"

	"dvra-api/internal/shared/models"
)

// La career page pública solo debe exponer el nombre y la ruta del
// departamento: nunca su identificador interno ni el responsable asignado.
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
