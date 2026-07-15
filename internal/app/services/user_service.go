package services

import (
	"time"

	"dvra-api/internal/app/dtos"
	"dvra-api/internal/app/models"
	"dvra-api/internal/app/repositories"
	"dvra-api/internal/shared/apperr"
	"gorm.io/gorm"
)

// UserService define el contrato del servicio de usuarios
type UserService interface {
	GetAllUsers() ([]models.User, error)
	GetUserByID(id uint) (*models.User, error)
	GetUsersByCompanyID(companyID uint) ([]models.User, error)
	CreateUser(dto dtos.CreateUserDTO) (*models.User, error)
	UpdateUser(id, companyID uint, global bool, dto dtos.UpdateUserDTO) (*models.User, error)
	DeleteUser(id, companyID uint, global bool) error
}

// userService es la implementación privada del servicio
type userService struct {
	userRepo repositories.UserRepository
	db       *gorm.DB
}

// NewUserService crea una nueva instancia de UserService
func NewUserService(userRepo repositories.UserRepository, db *gorm.DB) UserService {
	return &userService{
		userRepo: userRepo,
		db:       db,
	}
}

// GetAllUsers obtiene todos los usuarios
func (s *userService) GetAllUsers() ([]models.User, error) {
	return s.userRepo.GetAll()
}

// GetUsersByCompanyID obtiene usuarios de una empresa específica
func (s *userService) GetUsersByCompanyID(companyID uint) ([]models.User, error) {
	return s.userRepo.GetByCompanyID(companyID)
}

// GetUserByID obtiene un usuario por su ID
func (s *userService) GetUserByID(id uint) (*models.User, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, apperr.NotFound("user not found")
	}
	return user, nil
}

// CreateUser crea un nuevo usuario
func (s *userService) CreateUser(dto dtos.CreateUserDTO) (*models.User, error) {
	// Aquí puedes agregar lógica de negocio adicional:
	// - Validar que el email no exista
	// - Hashear contraseña
	// - Enviar email de bienvenida
	// - etc.

	// Verificar si el email ya existe
	existingUser, err := s.userRepo.GetByEmail(dto.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, apperr.Conflict("email already exists")
	}

	hashedPassword, err := HashPassword(dto.Password)
	if err != nil {
		return nil, err
	}
	user := &models.User{
		FirstName:    dto.FirstName,
		LastName:     dto.LastName,
		Email:        dto.Email,
		PasswordHash: hashedPassword,
		IsActive:     true,
	}
	invitedAt := time.Now()
	err = s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}
		membership := &models.Membership{
			UserID: user.ID, CompanyID: dto.CompanyID, Role: dto.Role,
			Status: "active", IsDefault: true, InvitedBy: &dto.InvitedBy,
			InvitedAt: &invitedAt, JoinedAt: &invitedAt,
		}
		return tx.Create(membership).Error
	})
	if err != nil {
		return nil, err
	}
	return user, nil
}

// UpdateUser actualiza un usuario existente
func (s *userService) UpdateUser(id, companyID uint, global bool, dto dtos.UpdateUserDTO) (*models.User, error) {
	existingUser, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if existingUser == nil {
		return nil, apperr.NotFound("user not found")
	}
	if !global && !userHasCompany(existingUser, companyID) {
		return nil, apperr.Forbidden("access denied")
	}

	// Actualizar campos si se proporcionan
	if dto.FirstName != nil {
		existingUser.FirstName = *dto.FirstName
	}
	if dto.LastName != nil {
		existingUser.LastName = *dto.LastName
	}
	if dto.Email != nil {
		existingUser.Email = *dto.Email
	}

	updatedUser, err := s.userRepo.Update(existingUser)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}

// DeleteUser elimina un usuario
func (s *userService) DeleteUser(id, companyID uint, global bool) error {
	existingUser, err := s.userRepo.GetByID(id)
	if err != nil {
		return err
	}
	if existingUser == nil {
		return apperr.NotFound("user not found")
	}
	if !global && !userHasCompany(existingUser, companyID) {
		return apperr.Forbidden("access denied")
	}

	if global {
		return s.userRepo.Delete(id)
	}
	if err := s.userRepo.RemoveFromCompany(id, companyID); err != nil {
		if err == gorm.ErrRecordNotFound {
			return apperr.NotFound("tenant membership not found")
		}
		return err
	}
	return nil
}

func userHasCompany(user *models.User, companyID uint) bool {
	for _, membership := range user.Memberships {
		if membership.CompanyID == companyID {
			return true
		}
	}
	return false
}
