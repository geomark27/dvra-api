package permissions

// Permisos del módulo Companies
const (
	CompaniesView   = "companies.view"
	CompaniesUpdate = "companies.update"
	// Las acciones están catalogadas para que admin represente acceso tenant
	// completo. Crear/eliminar empresas sigue siendo una ruta platform y exige
	// autorización landlord independiente.
	CompaniesCreate = "companies.create"
	CompaniesDelete = "companies.delete"
)

func init() {
	define(CompaniesView, CompaniesUpdate, CompaniesCreate, CompaniesDelete)
	grant(RoleAdmin, CompaniesView, CompaniesUpdate)
}
