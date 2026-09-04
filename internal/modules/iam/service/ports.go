package service

import (
	"dvra-api/internal/shared/models"
)

// Puertos de salida del módulo iam, declarados en el consumidor
// (ESTANDARES_DESARROLLO.md §3.3 regla 3).

type UserRepository interface {
	GetAll() ([]models.User, error)
	GetByID(id uint) (*models.User, error)
	GetByCompanyID(companyID uint) ([]models.User, error)
	Create(user *models.User) (*models.User, error)
	Update(user *models.User) (*models.User, error)
	Delete(id uint) error
	RemoveFromCompany(userID, companyID uint) error
	GetByEmail(email string) (*models.User, error)
	GetUserWithMemberships(userID uint) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	FindByID(id int) (*models.User, error)
	UpdateLastLogin(userID uint) error
	UpdatePassword(userID uint, hashedPassword string) error
}

type CompanyRepository interface {
	GetAll() ([]models.Company, error)
	GetByID(id uint) (*models.Company, error)
	GetBySlug(slug string) (*models.Company, error)
	Create(company *models.Company) (*models.Company, error)
	Update(company *models.Company) (*models.Company, error)
	Delete(id uint) error
	GetCompaniesWithMembers(companyID uint) (*models.Company, error)
}

type MembershipRepository interface {
	GetAll() ([]models.Membership, error)
	GetByID(id uint) (*models.Membership, error)
	GetByUserID(userID uint) ([]models.Membership, error)
	GetByCompanyID(companyID uint) ([]models.Membership, error)
	GetByUserAndCompany(userID uint, companyID uint) (*models.Membership, error)
	Create(membership *models.Membership) (*models.Membership, error)
	Update(membership *models.Membership) (*models.Membership, error)
	Delete(id uint) error
}

// PlanReader es lo ÚNICO que iam necesita de billing: resolver el plan inicial
// al dar de alta una empresa. El composition root inyecta el repositorio de
// planes, que lo satisface sin conocer este paquete.
type PlanReader interface {
	FindActiveBySlug(slug string) (*models.Plan, error)
}
