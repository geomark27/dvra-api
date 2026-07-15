package services

import (
	"errors"
	"testing"

	"dvra-api/internal/app/dtos"
	"dvra-api/internal/app/models"
	"dvra-api/internal/app/repositories"
	"dvra-api/internal/shared/apperr"
	"dvra-api/internal/shared/tenant"
)

type applicationRepoStub struct {
	repositories.ApplicationRepository
	created bool
}

func (f *applicationRepoStub) GetByCandidateAndJob(uint, uint) (*models.Application, error) {
	return nil, nil
}

func (f *applicationRepoStub) Create(application *models.Application) (*models.Application, error) {
	f.created = true
	return application, nil
}

type applicationJobRepoStub struct {
	repositories.JobRepository
	job *models.Job
}

func (f *applicationJobRepoStub) GetByID(uint, tenant.Scope) (*models.Job, error) {
	return f.job, nil
}

type applicationCandidateRepoStub struct {
	repositories.CandidateRepository
	candidate *models.Candidate
}

func (f *applicationCandidateRepoStub) GetByID(uint, tenant.Scope) (*models.Candidate, error) {
	return f.candidate, nil
}

func TestCreateApplicationRejectsCrossTenantEntities(t *testing.T) {
	tests := []struct {
		name      string
		job       *models.Job
		candidate *models.Candidate
	}{
		{name: "job", job: nil, candidate: &models.Candidate{CompanyID: 2}},
		{name: "candidate", job: &models.Job{CompanyID: 2}, candidate: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			applications := &applicationRepoStub{}
			svc := NewApplicationService(
				applications,
				&applicationJobRepoStub{job: tt.job},
				&applicationCandidateRepoStub{candidate: tt.candidate},
			)
			_, err := svc.CreateApplication(dtos.CreateApplicationDTO{
				CompanyID: 2, JobID: 10, CandidateID: 20, Stage: "applied",
			})
			var appErr *apperr.AppError
			if !errors.As(err, &appErr) || appErr.Code != 404 {
				t.Fatalf("expected not found for cross-tenant %s, got %v", tt.name, err)
			}
			if applications.created {
				t.Fatal("application was persisted after tenant validation failed")
			}
		})
	}
}

func TestCreateApplicationAcceptsSameTenantEntities(t *testing.T) {
	applications := &applicationRepoStub{}
	svc := NewApplicationService(
		applications,
		&applicationJobRepoStub{job: &models.Job{CompanyID: 2}},
		&applicationCandidateRepoStub{candidate: &models.Candidate{CompanyID: 2}},
	)
	created, err := svc.CreateApplication(dtos.CreateApplicationDTO{
		CompanyID: 2, JobID: 10, CandidateID: 20, Stage: "applied",
	})
	if err != nil {
		t.Fatal(err)
	}
	if !applications.created || created.CompanyID != 2 {
		t.Fatal("same-tenant application was not persisted")
	}
}
