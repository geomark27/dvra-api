package permissions

const (
	DepartmentsView   = "departments.view"
	DepartmentsCreate = "departments.create"
	DepartmentsUpdate = "departments.update"
	DepartmentsDelete = "departments.delete"
)

func init() {
	define(DepartmentsView, DepartmentsCreate, DepartmentsUpdate, DepartmentsDelete)
	grant(RoleAdmin, DepartmentsView, DepartmentsCreate, DepartmentsUpdate, DepartmentsDelete)
	grant(RoleRecruiter, DepartmentsView)
	grant(RoleHiringManager, DepartmentsView)
	grant(RoleInterviewer, DepartmentsView)
}
