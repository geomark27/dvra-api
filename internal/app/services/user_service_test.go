package services

import (
	"testing"

	"dvra-api/internal/app/dtos"
	"dvra-api/internal/app/models"
)

type userRepoFake struct {
	user    *models.User
	updated bool
	deleted bool
	removed bool
}

func (f *userRepoFake) GetAll() ([]models.User, error)              { return nil, nil }
func (f *userRepoFake) GetByID(uint) (*models.User, error)          { return f.user, nil }
func (f *userRepoFake) GetByCompanyID(uint) ([]models.User, error)  { return nil, nil }
func (f *userRepoFake) Create(u *models.User) (*models.User, error) { return u, nil }
func (f *userRepoFake) Update(u *models.User) (*models.User, error) {
	f.updated = true
	return u, nil
}
func (f *userRepoFake) Delete(uint) error                                 { f.deleted = true; return nil }
func (f *userRepoFake) RemoveFromCompany(uint, uint) error                { f.removed = true; return nil }
func (f *userRepoFake) GetByEmail(string) (*models.User, error)           { return nil, nil }
func (f *userRepoFake) GetUserWithMemberships(uint) (*models.User, error) { return f.user, nil }
func (f *userRepoFake) FindByEmail(string) (*models.User, error)          { return nil, nil }
func (f *userRepoFake) FindByID(int) (*models.User, error)                { return f.user, nil }
func (f *userRepoFake) UpdateLastLogin(uint) error                        { return nil }
func (f *userRepoFake) UpdatePassword(uint, string) error                 { return nil }

func TestUserMutationsRejectOtherTenant(t *testing.T) {
	companyID := uint(2)
	repo := &userRepoFake{user: &models.User{Memberships: []models.Membership{{CompanyID: companyID}}}}
	svc := NewUserService(repo, nil, testTenantAuthorizer())
	name := "Changed"
	if _, err := svc.UpdateUser(1, 3, false, dtos.UpdateUserDTO{FirstName: &name}); err == nil {
		t.Fatal("cross-tenant update was accepted")
	}
	if err := svc.DeleteUser(1, 3, false); err == nil {
		t.Fatal("cross-tenant delete was accepted")
	}
	if repo.updated || repo.deleted {
		t.Fatal("repository mutation occurred after access denial")
	}
}

func TestUserMutationsAllowSameTenant(t *testing.T) {
	companyID := uint(2)
	repo := &userRepoFake{user: &models.User{Memberships: []models.Membership{{CompanyID: companyID}}}}
	svc := NewUserService(repo, nil, testTenantAuthorizer())
	name := "Changed"
	if _, err := svc.UpdateUser(1, companyID, false, dtos.UpdateUserDTO{FirstName: &name}); err != nil {
		t.Fatal(err)
	}
	if err := svc.DeleteUser(1, companyID, false); err != nil {
		t.Fatal(err)
	}
	if !repo.updated || !repo.removed || repo.deleted {
		t.Fatal("same-tenant mutations did not reach repository")
	}
}

func TestCreateUserRejectsUnknownTenantRole(t *testing.T) {
	svc := NewUserService(&userRepoFake{}, nil, testTenantAuthorizer())
	if _, err := svc.CreateUser(dtos.CreateUserDTO{TenantRoleID: 999}); err == nil {
		t.Fatal("unknown tenant role ID was accepted")
	}
}
