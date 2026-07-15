package services

import (
	"time"

	"dvra-api/internal/app/dtos"
	"dvra-api/internal/app/models"
	"dvra-api/internal/app/repositories"
	"dvra-api/internal/shared/apperr"
)

// MembershipService define el contrato del servicio de memberships
type MembershipService interface {
	GetAllMemberships() ([]models.Membership, error)
	GetMembershipByID(id uint) (*models.Membership, error)
	GetMembershipsByUserID(userID uint) ([]models.Membership, error)
	GetMembershipsByCompanyID(companyID uint) ([]models.Membership, error)
	CreateMembership(dto dtos.CreateMembershipDTO) (*models.Membership, error)
	UpdateMembership(id, actorUserID, actorCompanyID uint, global bool, dto dtos.UpdateMembershipDTO) (*models.Membership, error)
	DeleteMembership(id uint) error
}

type membershipService struct {
	membershipRepo repositories.MembershipRepository
}

func NewMembershipService(membershipRepo repositories.MembershipRepository) MembershipService {
	return &membershipService{membershipRepo: membershipRepo}
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

func (s *membershipService) CreateMembership(dto dtos.CreateMembershipDTO) (*models.Membership, error) {
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
		UserID:    dto.UserID,
		CompanyID: *dto.CompanyID,
		Role:      dto.Role,
		Status:    status,
		IsDefault: dto.IsDefault,
		InvitedBy: dto.InvitedBy,
		InvitedAt: &now,
	}

	return s.membershipRepo.Create(membership)
}

func (s *membershipService) UpdateMembership(id, actorUserID, actorCompanyID uint, global bool, dto dtos.UpdateMembershipDTO) (*models.Membership, error) {
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
		if dto.Role != nil && membership.UserID == actorUserID {
			return nil, apperr.Forbidden("you cannot change your own role")
		}
	}

	if dto.Role != nil {
		if *dto.Role == "superadmin" {
			return nil, apperr.Forbidden("superadmin is a reserved platform role")
		}
		membership.Role = *dto.Role
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
