package permissions

// Permisos del módulo Applications (pipeline)
const (
	ApplicationsView   = "applications.view"
	ApplicationsCreate = "applications.create"
	ApplicationsUpdate = "applications.update"
	ApplicationsMove   = "applications.move"
	ApplicationsRate   = "applications.rate"
	ApplicationsDelete = "applications.delete"
)

func init() {
	grant(RoleAdmin, ApplicationsView, ApplicationsCreate, ApplicationsUpdate, ApplicationsMove, ApplicationsRate, ApplicationsDelete)
	grant(RoleRecruiter, ApplicationsView, ApplicationsCreate, ApplicationsUpdate, ApplicationsMove, ApplicationsRate)
	// hiring_manager: califica y comenta; mover stage solo en sus jobs (matriz 3.2),
	// pendiente de chequeo a nivel de recurso (RN-MEMB-007).
	grant(RoleHiringManager, ApplicationsView, ApplicationsRate)
	grant(RoleUser, ApplicationsView)
}
