package permissions

// Permisos del módulo Memberships
const (
	MembershipsView   = "memberships.view"
	MembershipsUpdate = "memberships.update"
	MembershipsDelete = "memberships.delete"
	// La acción se cataloga para el perfil admin completo. En el MVP no existe
	// una ruta tenant para crear memberships; el alta directa sigue siendo una
	// operación platform y la creación tenant ocurre vía POST /users.
	MembershipsCreate = "memberships.create"
)

func init() {
	define(MembershipsView, MembershipsCreate, MembershipsUpdate, MembershipsDelete)
	grant(RoleAdmin, MembershipsView, MembershipsUpdate, MembershipsDelete)
}
