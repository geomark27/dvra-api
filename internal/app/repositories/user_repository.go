package repositories

import (
	"dvra-api/internal/app/models"
	"dvra-api/internal/database"

	"gorm.io/gorm"
)

// UserRepository define el contrato del repositorio de usuarios
type UserRepository interface {
	GetAll() ([]models.User, error)
	GetByID(id uint) (*models.User, error)
	Create(user *models.User) (*models.User, error)
	Update(user *models.User) (*models.User, error)
	Delete(id uint) error
	GetByEmail(email string) (*models.User, error)
	GetUserWithMemberships(userID uint) (*models.User, error)
}

// userRepository es la implementación con GORM
type userRepository struct{}

// NewUserRepository crea una nueva instancia de UserRepository
func NewUserRepository() UserRepository {
	return &userRepository{}
}

// GetAll obtiene todos los usuarios
func (r *userRepository) GetAll() ([]models.User, error) {
	var users []models.User
	if err := database.DB.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

// GetByID obtiene un usuario por su ID
func (r *userRepository) GetByID(id uint) (*models.User, error) {
	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// GetByEmail obtiene un usuario por su email
func (r *userRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	if err := database.DB.Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

// Create crea un nuevo usuario
func (r *userRepository) Create(user *models.User) (*models.User, error) {
	if err := database.DB.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// Update actualiza un usuario existente
func (r *userRepository) Update(user *models.User) (*models.User, error) {
	if err := database.DB.Save(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// Delete elimina un usuario (soft delete)
func (r *userRepository) Delete(id uint) error {
	return database.DB.Delete(&models.User{}, id).Error
}

// GetUserWithMemberships obtiene un usuario con sus memberships
func (r *userRepository) GetUserWithMemberships(userID uint) (*models.User, error) {
	var user models.User
	if err := database.DB.Preload("Memberships").Preload("Memberships.Company").First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
