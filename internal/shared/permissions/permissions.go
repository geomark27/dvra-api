// Package permissions define las claves y la configuración inicial de permisos.
//
// Cada módulo declara sus claves y grants iniciales. Las migraciones sincronizan
// esta configuración hacia tenant_roles, permissions y tenant_role_permissions;
// la autorización de peticiones siempre se resuelve desde esas tablas.
package permissions

import "sort"

type TenantRoleDefinition struct {
	Key         string
	Name        string
	Description string
	Permissions []string
}

// Claves de roles incorporados por el producto.
const (
	RoleSuperAdmin    = "superadmin"
	RoleAdmin         = "admin"
	RoleRecruiter     = "recruiter"
	RoleHiringManager = "hiring_manager"
	RoleInterviewer   = "interviewer"
)

var tenantRoles = []string{
	RoleAdmin,
	RoleRecruiter,
	RoleHiringManager,
	RoleInterviewer,
}

// rolePermissions es la matriz rol → set de permisos.
// Se llena desde los init() de cada archivo de módulo.
var rolePermissions = map[string]map[string]bool{}
var tenantPermissionCatalog = map[string]bool{}

// define registra acciones tenant aunque todavía no tengan una ruta asignada.
// Admin recibe siempre el catálogo completo.
func define(perms ...string) {
	for _, permission := range perms {
		tenantPermissionCatalog[permission] = true
	}
}

// grant asigna permisos a un rol. Uso exclusivo de los init() de este paquete.
func grant(role string, perms ...string) {
	define(perms...)
	set, ok := rolePermissions[role]
	if !ok {
		set = make(map[string]bool)
		rolePermissions[role] = set
	}
	for _, p := range perms {
		set[p] = true
	}
}

// Can only inspects the release-controlled seed matrix. Runtime authorization
// must use services.TenantAuthorizer, which reads the relational RBAC tables.
func Can(role, permission string) bool {
	if role == RoleAdmin {
		return tenantPermissionCatalog[permission]
	}
	return rolePermissions[role][permission]
}

// IsTenantRole reports whether role is one of DVRA's built-in tenant profiles.
// Platform roles such as superadmin are intentionally excluded.
func IsTenantRole(role string) bool {
	for _, tenantRole := range tenantRoles {
		if tenantRole == role {
			return true
		}
	}
	return false
}

// TenantRoles returns a copy of the valid built-in tenant roles.
func TenantRoles() []string {
	return append([]string(nil), tenantRoles...)
}

// TenantRoleDefinitions is the release-controlled seed for relational tenant
// RBAC. Runtime authorization is resolved from the database, not this map.
func TenantRoleDefinitions() []TenantRoleDefinition {
	metadata := map[string][2]string{
		RoleAdmin:         {"Administrador", "Administración completa del tenant"},
		RoleRecruiter:     {"Recruiter", "Operación diaria de vacantes y candidatos"},
		RoleHiringManager: {"Hiring Manager", "Evaluación de vacantes asignadas"},
		RoleInterviewer:   {"Interviewer", "Lectura básica y futuras entrevistas asignadas"},
	}
	definitions := make([]TenantRoleDefinition, 0, len(tenantRoles))
	for _, key := range tenantRoles {
		definitions = append(definitions, TenantRoleDefinition{
			Key: key, Name: metadata[key][0], Description: metadata[key][1], Permissions: For(key),
		})
	}
	return definitions
}

func AllTenantPermissions() []string {
	return sortedKeys(tenantPermissionCatalog)
}

// For devuelve la lista ordenada de permisos de un rol.
// Pensado para exponerse al frontend (GET /auth/me) y que la UI
// muestre/oculte acciones sin duplicar la matriz.
func For(role string) []string {
	if role == RoleAdmin {
		return AllTenantPermissions()
	}
	return sortedKeys(rolePermissions[role])
}

func sortedKeys(set map[string]bool) []string {
	keys := make([]string, 0, len(set))
	for k := range set {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
