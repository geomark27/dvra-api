package permissions

// Permisos del módulo Memberships
const (
	MembershipsView   = "memberships.view"
	MembershipsUpdate = "memberships.update"
	MembershipsDelete = "memberships.delete"
	// MembershipsCreate no se asigna a ningún rol: en el MVP solo SuperAdmin
	// crea memberships (RN-MEMB-004). Cuando exista el sistema de invitaciones
	// por email (Fase 2) se otorgará a admin un permiso memberships.invite.
	MembershipsCreate = "memberships.create"
)

func init() {
	grant(RoleAdmin, MembershipsView, MembershipsUpdate, MembershipsDelete)
}
