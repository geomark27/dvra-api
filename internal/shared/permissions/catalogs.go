package permissions

const (
	CatalogsView   = "catalogs.view"
	CatalogsManage = "catalogs.manage"
)

func init() {
	define(CatalogsView, CatalogsManage)
	grant(RoleRecruiter, CatalogsView)
	grant(RoleHiringManager, CatalogsView)
	grant(RoleInterviewer, CatalogsView)
}
