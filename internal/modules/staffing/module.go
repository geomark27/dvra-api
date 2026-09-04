// Package staffing es el punto de ensamblaje del módulo: arma repos -> services y
// expone el registro de rutas. Nadie importa este paquete salvo el composition root.
package staffing

import (
	"dvra-api/internal/modules/staffing/domain"
	"dvra-api/internal/modules/staffing/repository"
	"dvra-api/internal/modules/staffing/service"
	"dvra-api/internal/modules/staffing/transport"
	"dvra-api/internal/platform/httpx"

	"gorm.io/gorm"
)

// Module agrupa las dependencias ya cableadas del módulo staffing.
type Module struct {
	// ClientRepo se expone para que otros módulos (vía un puerto que ellos definan)
	// puedan consultar clientes finales sin importar este paquete. Lo usa recruitment
	// para validar Job.StaffingClientID.
	ClientRepo domain.StaffingClientRepository

	clientSvc    *service.StaffingClientService
	placementSvc *service.PlacementService
}

// New construye el módulo. appFinder es el puerto cross-módulo hacia
// recruitment; puede venir nil y cablearse luego con BindApplications, que es lo
// que hace el composition root para romper el ciclo de construcción entre
// staffing y recruitment (cada uno necesita algo del otro).
func New(db *gorm.DB, appFinder domain.ApplicationFinder) *Module {
	clientRepo := repository.NewStaffingClientRepository(db)
	placementRepo := repository.NewPlacementRepository(db)

	return &Module{
		ClientRepo:   clientRepo,
		clientSvc:    service.NewStaffingClientService(clientRepo),
		placementSvc: service.NewPlacementService(placementRepo, clientRepo, appFinder),
	}
}

// Name identifica al módulo en logs y en la auditoría de rutas.
func (m *Module) Name() string { return "staffing" }

// RegisterRoutes monta las rutas del módulo bajo el grupo protegido.
func (m *Module) RegisterRoutes(r httpx.Router) {
	transport.RegisterRoutes(r.Protected, m.clientSvc, m.placementSvc, r.Plans, r.RBAC)
}

// BindApplications cablea el puerto hacia recruitment después de construir el
// módulo.
//
// Existe porque staffing y recruitment se necesitan mutuamente al arrancar:
// recruitment pide los clientes finales (que solo dependen de la base, así que
// están listos tras New) y staffing pide las postulaciones (que dependen del
// repositorio de recruitment). Construir en dos fases evita el ciclo sin
// recurrir a inicialización perezosa.
func (m *Module) BindApplications(appFinder domain.ApplicationFinder) {
	m.placementSvc.BindApplications(appFinder)
}
