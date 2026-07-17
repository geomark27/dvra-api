package permissions

import "testing"

func TestSuperAdminPermissionsArePlatformOnly(t *testing.T) {
	if Can(RoleSuperAdmin, JobsCreate) {
		t.Fatal("platform superadmin inherited tenant permissions")
	}
	for _, permission := range []string{PlatformCompaniesManage, PlatformPlansManage, PlatformSettingsManage} {
		if !PlatformCan(RoleSuperAdmin, permission) {
			t.Errorf("superadmin should have platform permission %q", permission)
		}
	}
}

func TestMatrixPorRol(t *testing.T) {
	cases := []struct {
		role string
		perm string
		want bool
	}{
		// admin: catálogo tenant completo. Las rutas landlord mantienen su
		// autorización platform independiente.
		{RoleAdmin, UsersCreate, true},
		{RoleAdmin, UsersUpdate, true},
		{RoleAdmin, CompaniesUpdate, true},
		{RoleAdmin, CompaniesCreate, true},
		{RoleAdmin, CompaniesDelete, true},
		{RoleAdmin, MembershipsCreate, true},

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

		// interviewer: conserva el acceso de lectura del perfil legacy user
		// hasta que exista asignación explícita de entrevistas y scorecards.
		{RoleInterviewer, JobsView, true},
		{RoleInterviewer, DashboardView, true},
		{RoleInterviewer, CandidatesView, false},
		{RoleInterviewer, ApplicationsView, false},
		{RoleInterviewer, JobsCreate, false},
		{RoleInterviewer, ApplicationsRate, false},
	}

	for _, tc := range cases {
		if got := Can(tc.role, tc.perm); got != tc.want {
			t.Errorf("Can(%q, %q) = %v, se esperaba %v", tc.role, tc.perm, got, tc.want)
		}
	}
}

func TestAdminHasEveryTenantPermission(t *testing.T) {
	all := AllTenantPermissions()
	if len(all) == 0 {
		t.Fatal("tenant permission catalog is empty")
	}
	for _, permission := range all {
		if !Can(RoleAdmin, permission) {
			t.Errorf("admin is missing %q", permission)
		}
	}
	if got := For(RoleAdmin); len(got) != len(all) {
		t.Fatalf("admin has %d permissions, catalog has %d", len(got), len(all))
	}
}

func TestTenantRoles(t *testing.T) {
	want := map[string]bool{
		RoleAdmin: true, RoleRecruiter: true,
		RoleHiringManager: true, RoleInterviewer: true,
	}
	for _, role := range TenantRoles() {
		if !want[role] {
			t.Errorf("unexpected tenant role %q", role)
		}
		delete(want, role)
	}
	if len(want) > 0 {
		t.Fatalf("missing tenant roles: %v", want)
	}
	for _, role := range []string{"", "user", RoleSuperAdmin, "intruso"} {
		if IsTenantRole(role) {
			t.Errorf("%q should not be a tenant role", role)
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
