package service

import (
	"dvra-api/internal/shared/apperr"
	"dvra-api/internal/shared/models"
)

// Puertos de salida del módulo billing, declarados en el consumidor
// (ESTANDARES_DESARROLLO.md §3.3 regla 3).

// PlanRepository es la persistencia de los planes de suscripción.
type PlanRepository interface {
	Create(plan *models.Plan) (*models.Plan, error)
	FindByID(id uint) (*models.Plan, error)
	FindBySlug(slug string) (*models.Plan, error)
	FindActiveBySlug(slug string) (*models.Plan, error)
	FindAll() ([]models.Plan, error)
	FindActive() ([]models.Plan, error)
	FindPublic() ([]models.Plan, error)
	Update(plan *models.Plan) (*models.Plan, error)
	Delete(id uint) error
	ExistsBySlug(slug string) (bool, error)
}

// CompanyReader es lo ÚNICO que billing necesita del dominio iam: resolver una
// empresa por id al asignarle un plan. Es la interfaz mínima a propósito — el
// composition root inyecta el repositorio de iam, que la satisface sin saberlo.
type CompanyReader interface {
	GetByID(id uint) (*models.Company, error)
}

// SubscriptionReader resuelve la suscripción vigente de una empresa.
//
// Existe para que CompanyHasFeature —que alimenta el gating por plan de TODAS
// las rutas premium— no llame a GORM desde el service y pueda probarse sin base
// de datos. Lo que decide qué suscripción es "la vigente" (el estado debe ser
// trialing, active o past_due) vive en el repositorio.
//
// PublishJobWithinLimit no pasa por aquí a propósito: necesita una transacción
// con bloqueo de fila para que dos publicaciones simultáneas no se salten el
// límite del plan, y eso no cabe tras un puerto de solo lectura sin introducir
// una abstracción de unidad de trabajo que el ADR-001 descarta.
type SubscriptionReader interface {
	CurrentPlan(companyID uint) (*models.Plan, error)
}

// ErrCompanyNotFound es el error de billing cuando la empresa referenciada no
// existe. Es propio del módulo y no se toma prestado de iam: un módulo no
// depende de las variables de error de otro.
var ErrCompanyNotFound = apperr.NotFound("company not found")
