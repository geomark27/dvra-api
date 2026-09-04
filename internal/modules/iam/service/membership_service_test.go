package service

import (
	"errors"
	"testing"

	"dvra-api/internal/shared/apperr"
	"dvra-api/internal/shared/models"
)

type membershipRepoFake struct{ membership *models.Membership }

func (f *membershipRepoFake) GetAll() ([]models.Membership, error)             { return nil, nil }
func (f *membershipRepoFake) GetByID(uint) (*models.Membership, error)         { return f.membership, nil }
func (f *membershipRepoFake) GetByUserID(uint) ([]models.Membership, error)    { return nil, nil }
func (f *membershipRepoFake) GetByCompanyID(uint) ([]models.Membership, error) { return nil, nil }
func (f *membershipRepoFake) GetByUserAndCompany(uint, uint) (*models.Membership, error) {
	return nil, nil
}
func (f *membershipRepoFake) Create(m *models.Membership) (*models.Membership, error) { return m, nil }
func (f *membershipRepoFake) Update(m *models.Membership) (*models.Membership, error) { return m, nil }
func (f *membershipRepoFake) Delete(uint) error                                       { return nil }

func TestMembershipUpdateRejectsUnknownRoleID(t *testing.T) {
	companyID := uint(3)
	repo := &membershipRepoFake{membership: &models.Membership{UserID: 8, CompanyID: companyID, TenantRoleID: 4}}
	svc := NewMembershipService(repo, testTenantAuthorizer())
	unknownRoleID := uint(999)
	_, err := svc.UpdateMembership(1, 9, companyID, false, UpdateMembershipDTO{TenantRoleID: &unknownRoleID})
	var appErr *apperr.AppError
	if !errors.As(err, &appErr) || appErr.Code != 400 {
		t.Fatalf("expected bad request, got %v", err)
	}
}

func TestMembershipUpdateRejectsOwnRoleAndOtherTenant(t *testing.T) {
	companyID := uint(3)
	repo := &membershipRepoFake{membership: &models.Membership{UserID: 8, CompanyID: companyID, TenantRoleID: 4}}
	svc := NewMembershipService(repo, testTenantAuthorizer())
	adminRoleID := uint(1)
	if _, err := svc.UpdateMembership(1, 8, companyID, false, UpdateMembershipDTO{TenantRoleID: &adminRoleID}); err == nil {
		t.Fatal("own role change was accepted")
	}
	if _, err := svc.UpdateMembership(1, 9, 4, false, UpdateMembershipDTO{}); err == nil {
		t.Fatal("cross-tenant membership update was accepted")
	}
}

func TestTenantMembershipRequiresCompany(t *testing.T) {
	svc := NewMembershipService(&membershipRepoFake{}, testTenantAuthorizer())
	if _, err := svc.CreateMembership(CreateMembershipDTO{UserID: 8, TenantRoleID: 1}); err == nil {
		t.Fatal("tenant membership without company was accepted")
	}
}

func TestTenantMembershipRejectsUnknownRole(t *testing.T) {
	companyID := uint(3)
	svc := NewMembershipService(&membershipRepoFake{}, testTenantAuthorizer())
	if _, err := svc.CreateMembership(CreateMembershipDTO{UserID: 8, CompanyID: &companyID, TenantRoleID: 999}); err == nil {
		t.Fatal("unknown tenant role ID was accepted")
	}
}
