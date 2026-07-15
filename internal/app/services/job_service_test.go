package services

import (
	"errors"
	"testing"

	"dvra-api/internal/app/models"
	"dvra-api/internal/app/repositories"
	"dvra-api/internal/shared/apperr"
	"dvra-api/internal/shared/tenant"
)

type jobRepoStub struct {
	repositories.JobRepository
	job *models.Job
}

func (f *jobRepoStub) GetByID(uint, tenant.Scope) (*models.Job, error) {
	return f.job, nil
}

type planServiceStub struct {
	PlanService
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
	svc := NewJobService(&jobRepoStub{job: job}, nil, plans)

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
	svc := NewJobService(&jobRepoStub{job: job}, nil, plans)

	_, err := svc.PublishJob(1, tenant.ForCompany(2))
	var appErr *apperr.AppError
	if !errors.As(err, &appErr) || appErr.Code != 409 {
		t.Fatalf("expected conflict, got %v", err)
	}
	if job.Status != "draft" {
		t.Fatal("in-memory status changed after plan limit rejection")
	}
}
