package repositories

import (
	"fmt"
	"os"
	"testing"
	"time"

	"dvra-api/internal/app/dtos"
	"dvra-api/internal/app/models"
	"dvra-api/internal/shared/permissions"
	"dvra-api/internal/shared/tenant"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestHiringManagerVisibilityWithPostgres(t *testing.T) {
	dsn := os.Getenv("DVRA_TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("DVRA_TEST_DATABASE_URL is not set")
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatalf("open postgres: %v", err)
	}
	tx := db.Begin()
	if tx.Error != nil {
		t.Fatalf("begin transaction: %v", tx.Error)
	}
	t.Cleanup(func() { _ = tx.Rollback().Error })

	suffix := time.Now().UnixNano()
	company := models.Company{Name: "Visibility A", Slug: fmt.Sprintf("visibility-a-%d", suffix), PlanTier: "free"}
	otherCompany := models.Company{Name: "Visibility B", Slug: fmt.Sprintf("visibility-b-%d", suffix), PlanTier: "free"}
	mustCreate(t, tx, &company, &otherCompany)

	manager := models.User{Email: fmt.Sprintf("manager-%d@example.com", suffix), PasswordHash: "not-used", IsActive: true}
	mustCreate(t, tx, &manager)
	managerID := manager.ID
	assignedJob := models.Job{CompanyID: company.ID, Title: "Assigned", Status: "draft", HiringManager: &managerID}
	unassignedJob := models.Job{CompanyID: company.ID, Title: "Not assigned", Status: "draft"}
	otherJob := models.Job{CompanyID: otherCompany.ID, Title: "Other tenant", Status: "draft", HiringManager: &managerID}
	mustCreate(t, tx, &assignedJob, &unassignedJob, &otherJob)

	visibleCandidate := models.Candidate{CompanyID: company.ID, Email: fmt.Sprintf("visible-%d@example.com", suffix)}
	hiddenCandidate := models.Candidate{CompanyID: company.ID, Email: fmt.Sprintf("hidden-%d@example.com", suffix)}
	otherCandidate := models.Candidate{CompanyID: otherCompany.ID, Email: fmt.Sprintf("other-%d@example.com", suffix)}
	mustCreate(t, tx, &visibleCandidate, &hiddenCandidate, &otherCandidate)
	visibleApplication := models.Application{CompanyID: company.ID, JobID: assignedJob.ID, CandidateID: visibleCandidate.ID, Stage: "applied", AppliedAt: time.Now()}
	hiddenApplication := models.Application{CompanyID: company.ID, JobID: unassignedJob.ID, CandidateID: hiddenCandidate.ID, Stage: "applied", AppliedAt: time.Now()}
	otherApplication := models.Application{CompanyID: otherCompany.ID, JobID: otherJob.ID, CandidateID: otherCandidate.ID, Stage: "applied", AppliedAt: time.Now()}
	mustCreate(t, tx, &visibleApplication, &hiddenApplication, &otherApplication)

	scope := tenant.ForActor(company.ID, manager.ID, permissions.RoleHiringManager)
	jobs, err := NewJobRepository(tx).GetByScopeWithFilters(scope, dtos.JobFilters{})
	assertOnlyID(t, "jobs", jobIDs(jobs), assignedJob.ID, err)
	candidates, err := NewCandidateRepository(tx).GetByScope(scope)
	assertOnlyID(t, "candidates", candidateIDs(candidates), visibleCandidate.ID, err)
	applications, err := NewApplicationRepository(tx).GetByScope(scope)
	assertOnlyID(t, "applications", applicationIDs(applications), visibleApplication.ID, err)

	job, err := NewJobRepository(tx).GetByID(unassignedJob.ID, scope)
	if err != nil || job != nil {
		t.Fatalf("unassigned job must be hidden, got job=%v err=%v", job, err)
	}
	candidate, err := NewCandidateRepository(tx).GetByID(hiddenCandidate.ID, scope)
	if err != nil || candidate != nil {
		t.Fatalf("candidate from unassigned job must be hidden, got candidate=%v err=%v", candidate, err)
	}
	application, err := NewApplicationRepository(tx).GetByID(hiddenApplication.ID, scope)
	if err != nil || application != nil {
		t.Fatalf("application from unassigned job must be hidden, got application=%v err=%v", application, err)
	}
}

func mustCreate(t *testing.T, db *gorm.DB, values ...interface{}) {
	t.Helper()
	for _, value := range values {
		if err := db.Create(value).Error; err != nil {
			t.Fatalf("create fixture: %v", err)
		}
	}
}

func assertOnlyID(t *testing.T, resource string, got []uint, want uint, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("list %s: %v", resource, err)
	}
	if len(got) != 1 || got[0] != want {
		t.Fatalf("%s visible ids = %v, want [%d]", resource, got, want)
	}
}

func jobIDs(items []models.Job) []uint {
	ids := make([]uint, len(items))
	for i := range items {
		ids[i] = items[i].ID
	}
	return ids
}

func candidateIDs(items []models.Candidate) []uint {
	ids := make([]uint, len(items))
	for i := range items {
		ids[i] = items[i].ID
	}
	return ids
}

func applicationIDs(items []models.Application) []uint {
	ids := make([]uint, len(items))
	for i := range items {
		ids[i] = items[i].ID
	}
	return ids
}
