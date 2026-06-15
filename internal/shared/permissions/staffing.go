package permissions

// Permisos del módulo Staffing: clientes finales (StaffingClient) y colocaciones
// (Placement). El acceso al módulo además está gateado por el plan vía
// RequireFeature("staffing"); esta matriz solo define qué rol puede operarlo.
const (
	StaffingClientsView   = "staffing_clients.view"
	StaffingClientsCreate = "staffing_clients.create"
	StaffingClientsUpdate = "staffing_clients.update"
	StaffingClientsDelete = "staffing_clients.delete"

	PlacementsView   = "placements.view"
	PlacementsCreate = "placements.create"
	PlacementsUpdate = "placements.update"
	PlacementsDelete = "placements.delete"
)

func init() {
	// admin: control total del módulo (delete reservado a admin, igual que jobs).
	grant(RoleAdmin,
		StaffingClientsView, StaffingClientsCreate, StaffingClientsUpdate, StaffingClientsDelete,
		PlacementsView, PlacementsCreate, PlacementsUpdate, PlacementsDelete,
	)
	// recruiter: opera clientes y placements pero no elimina.
	grant(RoleRecruiter,
		StaffingClientsView, StaffingClientsCreate, StaffingClientsUpdate,
		PlacementsView, PlacementsCreate, PlacementsUpdate,
	)
	// hiring_manager: lectura (los placements incluyen datos de billing).
	grant(RoleHiringManager, StaffingClientsView, PlacementsView)
	// user: solo ve el catálogo de clientes, no los placements/billing.
	grant(RoleUser, StaffingClientsView)
}
