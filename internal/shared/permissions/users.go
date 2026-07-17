package permissions

// Permisos del módulo Users (team members)
const (
	UsersView   = "users.view"
	UsersCreate = "users.create"
	UsersUpdate = "users.update"
	UsersDelete = "users.delete"
)

func init() {
	define(UsersView, UsersCreate, UsersUpdate, UsersDelete)
	grant(RoleAdmin, UsersView, UsersCreate, UsersDelete)
	// Ver team members es para todos los roles (matriz 3.2)
	grant(RoleRecruiter, UsersView)
	grant(RoleHiringManager, UsersView)
	grant(RoleInterviewer, UsersView)
}
