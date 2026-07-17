package services

import (
	"dvra-api/internal/app/models"

	"gorm.io/gorm"
)

type tenantAuthorizerFake struct {
	roles   map[uint]*models.TenantRole
	allowed map[uint]map[string]bool
}

func (f *tenantAuthorizerFake) Can(roleID uint, permission string) (bool, error) {
	return f.allowed[roleID][permission], nil
}

func (f *tenantAuthorizerFake) Permissions(roleID uint) ([]string, error) {
	result := make([]string, 0)
	for permission, allowed := range f.allowed[roleID] {
		if allowed {
			result = append(result, permission)
		}
	}
	return result, nil
}

func (f *tenantAuthorizerFake) RoleByID(id uint) (*models.TenantRole, error) {
	role, ok := f.roles[id]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	return role, nil
}

func (f *tenantAuthorizerFake) RoleByKey(key string) (*models.TenantRole, error) {
	for _, role := range f.roles {
		if role.Key == key {
			return role, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func (f *tenantAuthorizerFake) RoleForMembership(_, _, _ uint) (*models.TenantRole, error) {
	for _, role := range f.roles {
		return role, nil
	}
	return nil, gorm.ErrRecordNotFound
}

func testTenantAuthorizer() *tenantAuthorizerFake {
	return &tenantAuthorizerFake{roles: map[uint]*models.TenantRole{
		1: {BaseModel: models.BaseModel{ID: 1}, Key: "admin", Name: "Administrador", IsActive: true},
		4: {BaseModel: models.BaseModel{ID: 4}, Key: "interviewer", Name: "Entrevistador", IsActive: true},
	}, allowed: make(map[uint]map[string]bool)}
}
