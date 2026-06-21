// Package staffing es el punto de ensamblaje del módulo: arma repos -> services y
// expone el registro de rutas. Nadie importa este paquete salvo el composition root.
package staffing

import (
	"dvra-api/internal/app/services"
	"dvra-api/internal/modules/staffing/domain"
	"dvra-api/internal/modules/staffing/repository"
	"dvra-api/internal/modules/staffing/service"
	"dvra-api/internal/modules/staffing/transport"

	"github.com/gin-gonic/gin"
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

// New construye el módulo. appFinder es el puerto cross-módulo hacia recruitment
// (lo provee el composition root como adaptador sobre el repo de applications).
func New(db *gorm.DB, appFinder domain.ApplicationFinder) *Module {
	clientRepo := repository.NewStaffingClientRepository(db)
	placementRepo := repository.NewPlacementRepository(db)

	return &Module{
		ClientRepo:   clientRepo,
		clientSvc:    service.NewStaffingClientService(clientRepo),
		placementSvc: service.NewPlacementService(placementRepo, clientRepo, appFinder),
	}
}

// RegisterRoutes monta las rutas del módulo bajo el grupo protegido.
func (m *Module) RegisterRoutes(rg *gin.RouterGroup, planService services.PlanService) {
	transport.RegisterRoutes(rg, m.clientSvc, m.placementSvc, planService)
}
