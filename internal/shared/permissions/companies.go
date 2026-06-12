package permissions

// Permisos del módulo Companies
const (
	CompaniesView   = "companies.view"
	CompaniesUpdate = "companies.update"
	// CompaniesCreate y CompaniesDelete no se asignan a ningún rol:
	// crear empresas arbitrarias y eliminarlas es exclusivo del SuperAdmin
	// (auditoría de seguridad 2025-12-08). El alta normal de empresas es
	// vía POST /auth/register-company.
	CompaniesCreate = "companies.create"
	CompaniesDelete = "companies.delete"
)

func init() {
	grant(RoleAdmin, CompaniesView, CompaniesUpdate)
}
