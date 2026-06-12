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
	grant(RoleAdmin, JobsView, JobsCreate, JobsUpdate, JobsDelete, JobsPublish, JobsClose)
	grant(RoleRecruiter, JobsView, JobsCreate, JobsUpdate, JobsPublish, JobsClose)
	// hiring_manager edita solo sus jobs asignados (RN: matriz 3.2) — requiere
	// chequeo a nivel de recurso, pendiente de la tabla de asignaciones (RN-MEMB-007).
	// Hasta entonces solo tiene lectura.
	grant(RoleHiringManager, JobsView)
	grant(RoleUser, JobsView)
}
