// Package permissions define la matriz de autorización rol → permisos del sistema.
//
// La matriz vive en código (no en BD) porque los roles son parte del producto:
// cambian con releases, no en runtime. Cada módulo declara sus permisos en su
// propio archivo (jobs.go, candidates.go, ...) y los asigna a roles vía grant()
// en su init(), análogo a un seeder de permisos por módulo.
//
// Todo el proyecto consulta la matriz a través de Can(); si en el futuro se
// necesitan roles personalizados por empresa, basta con cambiar la
// implementación de Can() para leer de BD (con caché) sin tocar rutas,
// middleware ni constantes.
package permissions

import "sort"

// Roles del sistema (deben coincidir con memberships.role y el claim del JWT)
const (
	RoleSuperAdmin    = "superadmin"
	RoleAdmin         = "admin"
	RoleRecruiter     = "recruiter"
	RoleHiringManager = "hiring_manager"
	RoleUser          = "user"
)

// rolePermissions es la matriz rol → set de permisos.
// Se llena desde los init() de cada archivo de módulo.
var rolePermissions = map[string]map[string]bool{}

// grant asigna permisos a un rol. Uso exclusivo de los init() de este paquete.
func grant(role string, perms ...string) {
	set, ok := rolePermissions[role]
	if !ok {
		set = make(map[string]bool)
		rolePermissions[role] = set
	}
	for _, p := range perms {
		set[p] = true
	}
}

// Can reporta si el rol tiene el permiso indicado.
// SuperAdmin tiene acceso total por definición (RN-MEMB-003).
func Can(role, permission string) bool {
	if role == RoleSuperAdmin {
		return true
	}
	return rolePermissions[role][permission]
}

// For devuelve la lista ordenada de permisos de un rol.
// Pensado para exponerse al frontend (GET /auth/me) y que la UI
// muestre/oculte acciones sin duplicar la matriz.
func For(role string) []string {
	if role == RoleSuperAdmin {
		all := make(map[string]bool)
		for _, set := range rolePermissions {
			for p := range set {
				all[p] = true
			}
		}
		return sortedKeys(all)
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
