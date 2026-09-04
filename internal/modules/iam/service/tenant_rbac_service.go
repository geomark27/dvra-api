package service

import (
	"dvra-api/internal/shared/identity"
	"dvra-api/internal/shared/models"

	"gorm.io/gorm"
)

type TenantAuthorizer interface {
	Can(tenantRoleID uint, permissionKey string) (bool, error)
	Permissions(tenantRoleID uint) ([]string, error)
	RoleByID(id uint) (*models.TenantRole, error)
	RoleByKey(key string) (*models.TenantRole, error)
	RoleForMembership(membershipID, userID, companyID uint) (*models.TenantRole, error)
}

type TenantRBACService struct{ db *gorm.DB }

func NewTenantRBACService(db *gorm.DB) *TenantRBACService {
	return &TenantRBACService{db: db}
}

func (s *TenantRBACService) Can(tenantRoleID uint, permissionKey string) (bool, error) {
	var count int64
	err := s.db.Table("tenant_role_permissions trp").
		Joins("JOIN tenant_roles tr ON tr.id = trp.tenant_role_id AND tr.deleted_at IS NULL AND tr.is_active = TRUE").
		Joins("JOIN permissions p ON p.id = trp.permission_id AND p.deleted_at IS NULL AND p.is_active = TRUE").
		Where("trp.tenant_role_id = ? AND p.key = ?", tenantRoleID, permissionKey).
		Count(&count).Error
	return count > 0, err
}

func (s *TenantRBACService) Permissions(tenantRoleID uint) ([]string, error) {
	keys := make([]string, 0)
	err := s.db.Table("tenant_role_permissions trp").
		Joins("JOIN tenant_roles tr ON tr.id = trp.tenant_role_id AND tr.deleted_at IS NULL AND tr.is_active = TRUE").
		Joins("JOIN permissions p ON p.id = trp.permission_id AND p.deleted_at IS NULL AND p.is_active = TRUE").
		Where("trp.tenant_role_id = ?", tenantRoleID).
		Order("p.key").
		Pluck("p.key", &keys).Error
	return keys, err
}

func (s *TenantRBACService) RoleByID(id uint) (*models.TenantRole, error) {
	var role models.TenantRole
	if err := s.db.Where("id = ? AND is_active = ?", id, true).First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (s *TenantRBACService) RoleByKey(key string) (*models.TenantRole, error) {
	var role models.TenantRole
	if err := s.db.Where("key = ? AND is_active = ?", key, true).First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (s *TenantRBACService) RoleForMembership(membershipID, userID, companyID uint) (*models.TenantRole, error) {
	var membership models.Membership
	err := s.db.Preload("TenantRole", "is_active = ?", true).
		Where("id = ? AND user_id = ? AND company_id = ? AND status = ?", membershipID, userID, companyID, "active").
		First(&membership).Error
	if err != nil {
		return nil, err
	}
	if membership.TenantRole == nil {
		return nil, gorm.ErrRecordNotFound
	}
	return membership.TenantRole, nil
}

func (s *TenantRBACService) ListRoles() ([]models.TenantRole, error) {
	var roles []models.TenantRole
	err := s.db.Where("is_active = ?", true).Order("id").Find(&roles).Error
	return roles, err
}

// ResolveRole devuelve el rol vigente de una membresía en la vista neutral que
// consumen los middlewares (satisface middleware.RoleResolver).
//
// Se resuelve contra la base en cada petición a propósito: así revocar un rol o
// suspender una membresía surte efecto en la llamada siguiente, sin esperar a
// que caduque el token.
func (s *TenantRBACService) ResolveRole(membershipID, userID, companyID uint) (identity.TenantRole, error) {
	role, err := s.RoleForMembership(membershipID, userID, companyID)
	if err != nil {
		return identity.TenantRole{}, err
	}
	return identity.TenantRole{ID: role.ID, Key: role.Key}, nil
}
