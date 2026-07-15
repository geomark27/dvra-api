package permissions

const (
	PlatformCompaniesView     = "platform.companies.view"
	PlatformCompaniesManage   = "platform.companies.manage"
	PlatformUsersView         = "platform.users.view"
	PlatformMembershipsView   = "platform.memberships.view"
	PlatformMembershipsManage = "platform.memberships.manage"
	PlatformPlansManage       = "platform.plans.manage"
	PlatformSettingsManage    = "platform.settings.manage"
)

var platformRolePermissions = map[string]map[string]bool{
	RoleSuperAdmin: {
		PlatformCompaniesView: true, PlatformCompaniesManage: true,
		PlatformUsersView: true, PlatformMembershipsView: true,
		PlatformMembershipsManage: true, PlatformPlansManage: true,
		PlatformSettingsManage: true,
	},
}

func PlatformCan(role, permission string) bool {
	return platformRolePermissions[role][permission]
}

func ForPlatform(role string) []string {
	return sortedKeys(platformRolePermissions[role])
}
