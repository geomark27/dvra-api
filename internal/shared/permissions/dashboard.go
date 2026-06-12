package permissions

// Permisos del módulo Dashboard / Analytics
const (
	DashboardView = "dashboard.view"
)

func init() {
	// Ver analytics es para todos los roles (matriz 3.2)
	grant(RoleAdmin, DashboardView)
	grant(RoleRecruiter, DashboardView)
	grant(RoleHiringManager, DashboardView)
	grant(RoleUser, DashboardView)
}
