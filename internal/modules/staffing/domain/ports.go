// Package domain define el centro del módulo staffing: los puertos (interfaces)
// que el módulo necesita. No importa gin ni gorm — la infraestructura los implementa.
package domain

import (
	"dvra-api/internal/app/dtos"
	"dvra-api/internal/app/models"
)

// StaffingClientRepository es el puerto de salida hacia la persistencia de clientes finales.
type StaffingClientRepository interface {
	GetByCompanyID(companyID uint, status string) ([]models.StaffingClient, error)
	GetByID(id uint) (*models.StaffingClient, error)
	ExistsBySlug(companyID uint, slug string, excludeID uint) (bool, error)
	Create(client *models.StaffingClient) (*models.StaffingClient, error)
	Update(client *models.StaffingClient) (*models.StaffingClient, error)
	Delete(id uint) error
}

// PlacementRepository es el puerto de salida hacia la persistencia de colocaciones.
type PlacementRepository interface {
	GetByCompanyID(companyID uint, filters dtos.PlacementFilters) ([]models.Placement, error)
	GetByID(id uint) (*models.Placement, error)
	ExistsByApplicationID(applicationID uint) (bool, error)
	Create(placement *models.Placement) (*models.Placement, error)
	Update(placement *models.Placement) (*models.Placement, error)
	Delete(id uint) error
}

// HiredApplication es la vista mínima que staffing necesita de una Application
// (que vive en el módulo recruitment). Es un PUERTO DEFINIDO POR EL CONSUMIDOR:
// staffing NO importa recruitment; el composition root inyecta un adaptador.
type HiredApplication struct {
	ID          uint
	CompanyID   uint
	CandidateID uint
	JobID       uint
	Stage       string
}

// ApplicationFinder resuelve una Application por ID. Lo implementa un adaptador
// (en el composition root) que envuelve el repositorio de recruitment.
type ApplicationFinder interface {
	FindByID(id uint) (*HiredApplication, error)
}
