package permissions

// Permisos del módulo Jobs
const (
	JobsView    = "jobs.view"
	JobsCreate  = "jobs.create"
	JobsUpdate  = "jobs.update"
	JobsDelete  = "jobs.delete"
	JobsPublish = "jobs.publish"
	JobsClose   = "jobs.close"
)

func init() {
	define(JobsView, JobsCreate, JobsUpdate, JobsDelete, JobsPublish, JobsClose)
	grant(RoleAdmin, JobsView, JobsCreate, JobsUpdate, JobsDelete, JobsPublish, JobsClose)
	grant(RoleRecruiter, JobsView, JobsCreate, JobsUpdate, JobsPublish, JobsClose)
	// El repositorio limita la lectura del hiring_manager a sus jobs asignados.
	grant(RoleHiringManager, JobsView)
	grant(RoleInterviewer, JobsView)
}
