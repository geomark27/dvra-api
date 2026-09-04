package service

import (
	"errors"
	"testing"

	"dvra-api/internal/shared/apperr"
	"dvra-api/internal/shared/models"
	"dvra-api/internal/shared/tenant"
)

type jobRepoStub struct {
	JobRepository
	job     *models.Job
	created *models.Job
}

func (f *jobRepoStub) GetByID(uint, tenant.Scope) (*models.Job, error) {
	return f.job, nil
}

func (f *jobRepoStub) Create(job *models.Job) (*models.Job, error) {
	f.created = job
	return job, nil
}

type departmentReaderStub struct{ department *models.Department }

func (f departmentReaderStub) GetByID(_, companyID uint) (*models.Department, error) {
	if f.department == nil || f.department.CompanyID != companyID {
		return nil, nil
	}
	return f.department, nil
}

// Satisface el puerto planGuard, que es lo único que jobService necesita de
// billing.
type planServiceStub struct {
	called bool
	err    error
}

func (f *planServiceStub) PublishJobWithinLimit(*models.Job) error {
	f.called = true
	return f.err
}

func TestPublishJobUsesPlanLimitAndCanonicalStatus(t *testing.T) {
	plans := &planServiceStub{}
	job := &models.Job{CompanyID: 2, Title: "Backend", Description: "Build APIs", Status: "draft"}
	svc := NewJobService(&jobRepoStub{job: job}, nil, plans, nil)

	published, err := svc.PublishJob(1, tenant.ForCompany(2))
	if err != nil {
		t.Fatal(err)
	}
	if !plans.called {
		t.Fatal("plan limit was not checked")
	}
	if published.Status != "published" {
		t.Fatalf("expected published status, got %q", published.Status)
	}
}

func TestPublishJobPropagatesPlanLimitConflict(t *testing.T) {
	plans := &planServiceStub{err: apperr.Conflict("job limit reached")}
	job := &models.Job{CompanyID: 2, Title: "Backend", Description: "Build APIs", Status: "draft"}
	svc := NewJobService(&jobRepoStub{job: job}, nil, plans, nil)

	_, err := svc.PublishJob(1, tenant.ForCompany(2))
	var appErr *apperr.AppError
	if !errors.As(err, &appErr) || appErr.Code != 409 {
		t.Fatalf("expected conflict, got %v", err)
	}
	if job.Status != "draft" {
		t.Fatal("in-memory status changed after plan limit rejection")
	}
}

func TestCreateJobRequiresActiveDepartmentFromTenant(t *testing.T) {
	departmentID := uint(9)
	repo := &jobRepoStub{}
	svc := NewJobService(repo, nil, &planServiceStub{}, departmentReaderStub{
		department: &models.Department{BaseModel: models.BaseModel{ID: departmentID}, CompanyID: 99, IsActive: true},
	})
	_, err := svc.CreateJob(CreateJobDTO{
		CompanyID: 7, Title: "Backend Engineer", Description: "Build APIs",
		LocationType: "remote", DepartmentID: &departmentID,
	})
	var appErr *apperr.AppError
	if !errors.As(err, &appErr) || appErr.Code != 400 {
		t.Fatalf("expected invalid department, got %v", err)
	}
	if repo.created != nil {
		t.Fatal("job was persisted with a cross-tenant department")
	}
}
