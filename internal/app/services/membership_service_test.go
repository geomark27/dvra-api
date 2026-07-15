package services

import (
	"errors"
	"testing"

	"dvra-api/internal/app/dtos"
	"dvra-api/internal/app/models"
	"dvra-api/internal/shared/apperr"
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

func TestMembershipUpdateRejectsReservedRole(t *testing.T) {
	companyID := uint(3)
	repo := &membershipRepoFake{membership: &models.Membership{UserID: 8, CompanyID: companyID, Role: "user"}}
	svc := NewMembershipService(repo)
	role := "superadmin"
	_, err := svc.UpdateMembership(1, 9, companyID, false, dtos.UpdateMembershipDTO{Role: &role})
	var appErr *apperr.AppError
	if !errors.As(err, &appErr) || appErr.Code != 403 {
		t.Fatalf("expected forbidden, got %v", err)
	}
}

func TestMembershipUpdateRejectsOwnRoleAndOtherTenant(t *testing.T) {
	companyID := uint(3)
	repo := &membershipRepoFake{membership: &models.Membership{UserID: 8, CompanyID: companyID, Role: "user"}}
	svc := NewMembershipService(repo)
	role := "admin"
	if _, err := svc.UpdateMembership(1, 8, companyID, false, dtos.UpdateMembershipDTO{Role: &role}); err == nil {
		t.Fatal("own role change was accepted")
	}
	if _, err := svc.UpdateMembership(1, 9, 4, false, dtos.UpdateMembershipDTO{}); err == nil {
		t.Fatal("cross-tenant membership update was accepted")
	}
}

func TestTenantMembershipRequiresCompany(t *testing.T) {
	svc := NewMembershipService(&membershipRepoFake{})
	if _, err := svc.CreateMembership(dtos.CreateMembershipDTO{UserID: 8, Role: "admin"}); err == nil {
		t.Fatal("tenant membership without company was accepted")
	}
}
