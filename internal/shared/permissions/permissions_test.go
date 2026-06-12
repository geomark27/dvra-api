package permissions

import "testing"

func TestSuperAdminCanEverything(t *testing.T) {
	for _, perm := range []string{JobsCreate, CompaniesDelete, MembershipsCreate, UsersDelete} {
		if !Can(RoleSuperAdmin, perm) {
			t.Errorf("superadmin debería tener %q", perm)
		}
	}
}

func TestMatrixPorRol(t *testing.T) {
	cases := []struct {
		role string
		perm string
		want bool
	}{
		// admin: control total de su empresa, pero no crear/eliminar empresas
		{RoleAdmin, UsersCreate, true},
		{RoleAdmin, CompaniesUpdate, true},
		{RoleAdmin, CompaniesCreate, false},
		{RoleAdmin, CompaniesDelete, false},
		{RoleAdmin, MembershipsCreate, false}, // RN-MEMB-004: solo SuperAdmin en MVP

		// recruiter: gestiona jobs y candidatos, no usuarios ni billing
		{RoleRecruiter, JobsCreate, true},
		{RoleRecruiter, JobsPublish, true},
		{RoleRecruiter, CandidatesCreate, true},
		{RoleRecruiter, ApplicationsMove, true},
		{RoleRecruiter, UsersCreate, false},
		{RoleRecruiter, CompaniesView, false},
		{RoleRecruiter, JobsDelete, false},

		// hiring_manager: lectura + calificar
		{RoleHiringManager, JobsView, true},
		{RoleHiringManager, ApplicationsRate, true},
		{RoleHiringManager, JobsCreate, false},
		{RoleHiringManager, ApplicationsMove, false},
		{RoleHiringManager, CandidatesCreate, false},

		// user: solo lectura
		{RoleUser, JobsView, true},
		{RoleUser, DashboardView, true},
		{RoleUser, JobsCreate, false},
		{RoleUser, ApplicationsRate, false},
	}

	for _, tc := range cases {
		if got := Can(tc.role, tc.perm); got != tc.want {
			t.Errorf("Can(%q, %q) = %v, se esperaba %v", tc.role, tc.perm, got, tc.want)
		}
	}
}

func TestRolDesconocidoNoTienePermisos(t *testing.T) {
	if Can("intruso", JobsView) {
		t.Error("un rol desconocido no debe tener ningún permiso")
	}
	if Can("", JobsView) {
		t.Error("rol vacío no debe tener ningún permiso")
	}
}

func TestForDevuelveListaOrdenada(t *testing.T) {
	perms := For(RoleRecruiter)
	if len(perms) == 0 {
		t.Fatal("recruiter debe tener permisos")
	}
	for i := 1; i < len(perms); i++ {
		if perms[i-1] >= perms[i] {
			t.Fatalf("lista no ordenada o con duplicados: %v", perms)
		}
	}
}
