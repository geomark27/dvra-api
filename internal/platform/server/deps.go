package server

import (
	"dvra-api/internal/modules/billing"
	billingrepo "dvra-api/internal/modules/billing/repository"
	"dvra-api/internal/modules/catalog"
	"dvra-api/internal/modules/dashboard"
	"dvra-api/internal/modules/iam"
	"dvra-api/internal/modules/locations"
	"dvra-api/internal/modules/platformsettings"
	"dvra-api/internal/modules/publicapi"
	"dvra-api/internal/modules/recruitment"
	"dvra-api/internal/modules/staffing"
	"dvra-api/internal/platform/config"
	"dvra-api/internal/platform/health"
	"dvra-api/internal/platform/httpx"
	storagepkg "dvra-api/internal/platform/storage"
	"dvra-api/internal/shared/middleware"

	"gorm.io/gorm"
)

// deps es todo lo que el registro de rutas necesita: la lista de módulos y las
// políticas transversales que comparten.
//
// Ya no hay handlers sueltos: la migración por dominios terminó y cada módulo
// monta sus propias rutas.
type deps struct {
	Modules []httpx.Module

	SessionAuth middleware.Authenticator
	RBAC        middleware.PermissionChecker
	Plans       middleware.FeatureChecker
}

// buildDeps construye los módulos y los conecta entre sí.
//
// Es el composition root: el único punto que conoce a todos los módulos y que
// resuelve los puertos que cruzan sus fronteras. Ningún módulo importa a otro.
//
// El orden importa solo por las dependencias entre puertos:
// billing necesita las empresas de iam, iam necesita los planes de billing (por
// eso recibe el repositorio, no el service), y recruitment necesita ambos más
// los clientes finales de staffing.
func buildDeps(cfg *config.Config, db *gorm.DB) (*deps, error) {
	fileStorage, err := storagepkg.NewFilesystem(cfg.StorageRoot)
	if err != nil {
		return nil, err
	}

	// Repositorio de planes: lo consume iam para resolver el plan inicial al dar
	// de alta una empresa, sin depender del service de billing.
	billingPlans := billingrepo.NewPlanRepository(db)

	iamModule := iam.New(db, cfg.JWTSecret, cfg.JWTRefreshSecret, iam.Deps{Plans: billingPlans})
	billingModule := billing.New(db, iamModule.Companies)
	staffingModule := staffing.New(db, nil) // el puerto hacia recruitment se cablea abajo

	recruitmentModule := recruitment.New(db, recruitment.Deps{
		Plans:           billingPlans,
		PlanLimit:       billingModule.Plans(),
		StaffingClients: staffingModule.ClientRepo,
		Storage:         fileStorage,
	})

	// staffing necesita leer una Application de recruitment para derivar una
	// colocación. Se cablea después de construir recruitment, con un adaptador
	// que traduce al puerto que staffing declaró: ninguno importa al otro.
	staffingModule.BindApplications(staffingAppFinder{repo: recruitmentModule.Applications})

	publicModule := publicapi.New(publicapi.Deps{
		Companies:    iamModule.Companies,
		Jobs:         recruitmentModule.Jobs,
		Candidates:   recruitmentModule.Candidates,
		Applications: recruitmentModule.Applications,
		Resumes:      resumeStore{svc: recruitmentModule.Resumes()},
	})

	return &deps{
		Modules: []httpx.Module{
			iamModule,
			billingModule,
			recruitmentModule,
			staffingModule,
			catalog.New(db),
			locations.New(db),
			platformsettings.New(db),
			dashboard.New(db),
			publicModule,
			health.New(gormPinger{db: db}),
		},

		SessionAuth: iamModule.Auth(),
		RBAC:        iamModule.RBAC(),
		Plans:       billingModule.Plans(),
	}, nil
}
