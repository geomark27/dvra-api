package service

import (
	"time"

	"dvra-api/internal/shared/apperr"
	"dvra-api/internal/shared/models"
)

// MembershipService define el contrato del servicio de memberships
type MembershipService interface {
	GetAllMemberships() ([]models.Membership, error)
	GetMembershipByID(id uint) (*models.Membership, error)
	GetMembershipsByUserID(userID uint) ([]models.Membership, error)
	GetMembershipsByCompanyID(companyID uint) ([]models.Membership, error)
	CreateMembership(dto CreateMembershipDTO) (*models.Membership, error)
	UpdateMembership(id, actorUserID, actorCompanyID uint, global bool, dto UpdateMembershipDTO) (*models.Membership, error)
	DeleteMembership(id uint) error
}

type membershipService struct {
	membershipRepo MembershipRepository
	rbac           TenantAuthorizer
}

func NewMembershipService(membershipRepo MembershipRepository, rbac TenantAuthorizer) MembershipService {
	return &membershipService{membershipRepo: membershipRepo, rbac: rbac}
}

func (s *membershipService) GetAllMemberships() ([]models.Membership, error) {
	return s.membershipRepo.GetAll()
}

func (s *membershipService) GetMembershipByID(id uint) (*models.Membership, error) {
	membership, err := s.membershipRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if membership == nil {
		return nil, apperr.NotFound("membership not found")
	}
	return membership, nil
}

func (s *membershipService) GetMembershipsByUserID(userID uint) ([]models.Membership, error) {
	return s.membershipRepo.GetByUserID(userID)
}

func (s *membershipService) GetMembershipsByCompanyID(companyID uint) ([]models.Membership, error) {
	return s.membershipRepo.GetByCompanyID(companyID)
}

func (s *membershipService) CreateMembership(dto CreateMembershipDTO) (*models.Membership, error) {
	role, err := s.rbac.RoleByID(dto.TenantRoleID)
	if err != nil {
		return nil, apperr.BadRequest("invalid tenant_role_id")
	}
	if dto.CompanyID == nil || *dto.CompanyID == 0 {
		return nil, apperr.BadRequest("company_id is required for tenant memberships")
	}
	// Verificar si ya existe membership para user+company
	existing, err := s.membershipRepo.GetByUserAndCompany(dto.UserID, *dto.CompanyID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, apperr.Conflict("membership already exists for this user and company")
	}

	now := time.Now()
	status := dto.Status
	if status == "" {
		status = "active"
	}
	membership := &models.Membership{
		UserID:       dto.UserID,
		CompanyID:    *dto.CompanyID,
		TenantRoleID: dto.TenantRoleID,
		TenantRole:   role,
		Status:       status,
		IsDefault:    dto.IsDefault,
		InvitedBy:    dto.InvitedBy,
		InvitedAt:    &now,
	}

	return s.membershipRepo.Create(membership)
}

func (s *membershipService) UpdateMembership(id, actorUserID, actorCompanyID uint, global bool, dto UpdateMembershipDTO) (*models.Membership, error) {
	membership, err := s.membershipRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if membership == nil {
		return nil, apperr.NotFound("membership not found")
	}
	if !global {
		if membership.CompanyID != actorCompanyID {
			return nil, apperr.Forbidden("access denied")
		}
		if dto.TenantRoleID != nil && membership.UserID == actorUserID {
			return nil, apperr.Forbidden("you cannot change your own role")
		}
	}

	if dto.TenantRoleID != nil {
		role, err := s.rbac.RoleByID(*dto.TenantRoleID)
		if err != nil {
			return nil, apperr.BadRequest("invalid tenant_role_id")
		}
		membership.TenantRoleID = *dto.TenantRoleID
		membership.TenantRole = role
	}
	if dto.Status != nil {
		membership.Status = *dto.Status
	}
	if dto.IsDefault != nil {
		membership.IsDefault = *dto.IsDefault
	}
	if dto.JoinedAt != nil {
		membership.JoinedAt = dto.JoinedAt
	}

	return s.membershipRepo.Update(membership)
}

func (s *membershipService) DeleteMembership(id uint) error {
	membership, err := s.membershipRepo.GetByID(id)
	if err != nil {
		return err
	}
	if membership == nil {
		return apperr.NotFound("membership not found")
	}
	return s.membershipRepo.Delete(id)
}
