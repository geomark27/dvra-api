package services

import (
	"fmt"
	"time"

	"dvra-api/internal/app/dtos"
	"dvra-api/internal/app/models"
	"dvra-api/internal/app/repositories"
)

// MembershipService define el contrato del servicio de memberships
type MembershipService interface {
	GetAllMemberships() ([]models.Membership, error)
	GetMembershipByID(id uint) (*models.Membership, error)
	GetMembershipsByUserID(userID uint) ([]models.Membership, error)
	GetMembershipsByCompanyID(companyID uint) ([]models.Membership, error)
	CreateMembership(dto dtos.CreateMembershipDTO) (*models.Membership, error)
	UpdateMembership(id uint, dto dtos.UpdateMembershipDTO) (*models.Membership, error)
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
		return nil, fmt.Errorf("membership not found")
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
	// Verificar si ya existe membership para user+company
	if dto.CompanyID != nil {
		existing, err := s.membershipRepo.GetByUserAndCompany(dto.UserID, *dto.CompanyID)
		if err != nil {
			return nil, err
		}
		if existing != nil {
			return nil, fmt.Errorf("membership already exists for this user and company")
		}
	}

	now := time.Now()
	membership := &models.Membership{
		UserID:    dto.UserID,
		CompanyID: dto.CompanyID,
		Role:      dto.Role,
		Status:    dto.Status,
		IsDefault: dto.IsDefault,
		InvitedBy: dto.InvitedBy,
		InvitedAt: &now,
	}

	return s.membershipRepo.Create(membership)
}

func (s *membershipService) UpdateMembership(id uint, dto dtos.UpdateMembershipDTO) (*models.Membership, error) {
	membership, err := s.membershipRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if membership == nil {
		return nil, fmt.Errorf("membership not found")
	}

	if dto.Role != nil {
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
		return fmt.Errorf("membership not found")
	}
	return s.membershipRepo.Delete(id)
}
