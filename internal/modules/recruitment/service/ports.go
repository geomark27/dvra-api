package service

import (
	"dvra-api/internal/shared/apperr"
	"dvra-api/internal/shared/models"
	"dvra-api/internal/shared/tenant"
)

// Puertos de salida del módulo recruitment, declarados en el consumidor
// (ESTANDARES_DESARROLLO.md §3.3 regla 3).

// PlanReader es lo que recruitment necesita de billing para comprobar la cuota
// de almacenamiento antes de aceptar un CV.
type PlanReader interface {
	FindActiveBySlug(slug string) (*models.Plan, error)
}

// PlanGuard comprueba que la empresa puede publicar una vacante más según su
// plan. Va aparte de PlanReader porque la publicación abre una transacción con
// bloqueo de fila en billing, no una simple lectura.
type PlanGuard interface {
	PublishJobWithinLimit(job *models.Job) error
}

// StaffingClientReader resuelve el cliente final al que se atribuye una
// vacante en modo agencia. Puerto hacia staffing.
type StaffingClientReader interface {
	GetByID(id uint) (*models.StaffingClient, error)
}

// ErrPlanNotFound es el error propio del módulo. Billing tiene el suyo con el
// mismo mensaje: un módulo no toma prestadas las variables de error de otro.
var ErrPlanNotFound = apperr.NotFound("plan not found")

type ApplicationRepository interface {
	GetAll() ([]models.Application, error)
	GetByID(id uint, scope tenant.Scope) (*models.Application, error)
	GetByJobID(jobID uint) ([]models.Application, error)
	GetByCandidateID(candidateID uint) ([]models.Application, error)
	GetByCompanyID(companyID uint) ([]models.Application, error)
	GetByScope(scope tenant.Scope) ([]models.Application, error)
	GetByStage(stage string, companyID uint) ([]models.Application, error)
	GetByCandidateAndJob(candidateID, jobID uint) (*models.Application, error)
	Create(application *models.Application) (*models.Application, error)
	Update(application *models.Application) (*models.Application, error)
	Delete(id uint) error
}

type CandidateRepository interface {
	GetAll() ([]models.Candidate, error)
	GetByID(id uint, scope tenant.Scope) (*models.Candidate, error)
	GetByCompanyID(companyID uint) ([]models.Candidate, error)
	GetByScope(scope tenant.Scope) ([]models.Candidate, error)
	GetByEmail(email string, companyID uint) (*models.Candidate, error)
	Create(candidate *models.Candidate) (*models.Candidate, error)
	Update(candidate *models.Candidate) (*models.Candidate, error)
	Delete(id uint) error
}

type DepartmentRepository interface {
	List(companyID uint, includeInactive bool) ([]models.Department, error)
	GetByID(id, companyID uint) (*models.Department, error)
	NameExists(companyID uint, parentID *uint, name string, excludeID uint) (bool, error)
	CodeExists(companyID uint, code string, excludeID uint) (bool, error)
	UserIsActiveMember(userID, companyID uint) (bool, error)
	HasActiveChildren(id, companyID uint) (bool, error)
	Create(department *models.Department) (*models.Department, error)
	Update(department *models.Department) (*models.Department, error)
}

type JobRepository interface {
	GetAll() ([]models.Job, error)
	GetAllWithFilters(filters JobFilters) ([]models.Job, error)
	GetByID(id uint, scope tenant.Scope) (*models.Job, error)
	GetByCompanyID(companyID uint) ([]models.Job, error)
	GetByCompanyIDWithFilters(companyID uint, filters JobFilters) ([]models.Job, error)
	GetByScopeWithFilters(scope tenant.Scope, filters JobFilters) ([]models.Job, error)
	GetByStatus(status string, companyID uint) ([]models.Job, error)
	GetPublishedByCompanyID(companyID uint) ([]models.Job, error)
	Create(job *models.Job) (*models.Job, error)
	Update(job *models.Job) (*models.Job, error)
	Delete(id uint) error
}
