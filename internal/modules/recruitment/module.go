// Package recruitment es el punto de ensamblaje del núcleo del ATS:
// departamentos, vacantes, candidatos y el pipeline de postulaciones.
//
// Es el módulo con más superficie de datos de tenant, así que toda su
// persistencia va acotada por company_id (ESTANDARES §4).
package recruitment

import (
	"dvra-api/internal/modules/recruitment/repository"
	"dvra-api/internal/modules/recruitment/service"
	"dvra-api/internal/modules/recruitment/transport"
	"dvra-api/internal/platform/httpx"

	"gorm.io/gorm"
)

// Deps son las dependencias cross-módulo que recruitment necesita.
type Deps struct {
	// Plans resuelve el plan para la cuota de almacenamiento de CVs (billing).
	Plans service.PlanReader
	// PlanLimit comprueba el límite de vacantes publicadas (billing).
	PlanLimit service.PlanGuard
	// StaffingClients resuelve el cliente final de una vacante en modo agencia
	// (staffing).
	StaffingClients service.StaffingClientReader
	// Storage guarda los CVs.
	Storage service.FileStorage
}

// Module agrupa las dependencias ya cableadas del módulo.
type Module struct {
	departments service.DepartmentService
	jobs        service.JobService
	candidates  service.CandidateService
	apps        service.ApplicationService

	// Repositorios expuestos para que otros módulos los consuman a través de un
	// puerto que ellos declaren. Los usa publicapi para la career page y el
	// composition root para el puerto ApplicationFinder de staffing.
	Jobs         *repository.JobRepository
	Candidates   *repository.CandidateRepository
	Applications *repository.ApplicationRepository
}

// New construye el módulo con la conexión inyectada.
func New(db *gorm.DB, deps Deps) *Module {
	jobRepo := repository.NewJobRepository(db)
	candidateRepo := repository.NewCandidateRepository(db)
	applicationRepo := repository.NewApplicationRepository(db)
	departmentRepo := repository.NewDepartmentRepository(db)

	return &Module{
		departments:  service.NewDepartmentService(departmentRepo),
		jobs:         service.NewJobService(jobRepo, deps.StaffingClients, deps.PlanLimit, departmentRepo),
		candidates:   service.NewCandidateService(candidateRepo, deps.Plans, db, deps.Storage),
		apps:         service.NewApplicationService(applicationRepo, jobRepo, candidateRepo),
		Jobs:         jobRepo,
		Candidates:   candidateRepo,
		Applications: applicationRepo,
	}
}

// Resumes expone el service de candidatos, que aplica la cuota del plan al
// guardar un CV. Lo consume publicapi por su puerto ResumeStore.
func (m *Module) Resumes() service.CandidateService { return m.candidates }

// Name identifica al módulo en logs y en la auditoría de rutas.
func (m *Module) Name() string { return "recruitment" }

// RegisterRoutes monta las rutas de tenant.
func (m *Module) RegisterRoutes(r httpx.Router) {
	transport.RegisterRoutes(r.Protected, m.departments, m.jobs, m.candidates, m.apps, r.RBAC)
}
