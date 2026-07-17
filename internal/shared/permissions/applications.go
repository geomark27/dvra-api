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
	define(ApplicationsView, ApplicationsCreate, ApplicationsUpdate, ApplicationsMove, ApplicationsRate, ApplicationsDelete)
	grant(RoleAdmin, ApplicationsView, ApplicationsCreate, ApplicationsUpdate, ApplicationsMove, ApplicationsRate, ApplicationsDelete)
	grant(RoleRecruiter, ApplicationsView, ApplicationsCreate, ApplicationsUpdate, ApplicationsMove, ApplicationsRate)
	// El repositorio limita al hiring_manager a postulaciones de sus jobs.
	grant(RoleHiringManager, ApplicationsView, ApplicationsRate)
}
