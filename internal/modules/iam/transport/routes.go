package transport

import (
	"dvra-api/internal/modules/iam/service"
	"dvra-api/internal/shared/middleware"
	"dvra-api/internal/shared/permissions"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes monta identidad y acceso en los tres contextos.
//
// Los mismos handlers sirven a tenant y a landlord: la diferencia no está en el
// código sino en la autorización — RequirePermission consulta el rol relacional
// de la empresa, RequirePlatformPermission la matriz de plataforma. Por eso el
// módulo necesita los tres grupos de rutas (ADR-002).
func RegisterRoutes(
	public, protected, platform *gin.RouterGroup,
	auth *service.AuthService,
	users service.UserService,
	companies service.CompanyService,
	memberships service.MembershipService,
	rbac *service.TenantRBACService,
	perms middleware.PermissionChecker,
) {
	authH := NewAuthHandler(auth)
	userH := NewUserHandler(users)
	companyH := NewCompanyHandler(companies)
	membershipH := NewMembershipHandler(memberships)
	tenantRoleH := NewTenantRoleHandler(rbac)

	// --- Sesión ---
	authGroup := public.Group("/auth")
	{
		authGroup.POST("/register-company", authH.RegisterCompany)
		authGroup.POST("/register", authH.Register) // Deprecated: use register-company
		authGroup.POST("/login", authH.Login)
		authGroup.POST("/refresh", authH.RefreshToken)
	}

	authProtected := protected.Group("/auth")
	{
		authProtected.GET("/me", authH.GetMe)
		authProtected.POST("/change-password", authH.ChangePassword)
		authProtected.POST("/logout", authH.Logout)
		authProtected.POST("/switch-company", authH.SwitchCompany)
		authProtected.POST("/switch-platform", authH.SwitchToPlatform)
		authProtected.GET("/my-companies", authH.GetMyCompanies)
	}

	// --- Tenant ---
	protected.GET("/tenant-roles", middleware.RequirePermission(perms, permissions.UsersView), tenantRoleH.GetTenantRoles)

	users_ := protected.Group("/users")
	{
		users_.GET("", middleware.RequirePermission(perms, permissions.UsersView), userH.GetUsers)
		users_.POST("", middleware.RequirePermission(perms, permissions.UsersCreate), userH.CreateUser)
		users_.GET("/:id", middleware.RequirePermission(perms, permissions.UsersView), userH.GetUser)
		users_.DELETE("/:id", middleware.RequirePermission(perms, permissions.UsersDelete), userH.DeleteUser)
	}

	companiesGroup := protected.Group("/companies")
	{
		companiesGroup.GET("", middleware.RequirePermission(perms, permissions.CompaniesView), companyH.GetCompanies)
		companiesGroup.GET("/:id", middleware.RequirePermission(perms, permissions.CompaniesView), companyH.GetCompany)
		companiesGroup.PUT("/:id", middleware.RequirePermission(perms, permissions.CompaniesUpdate), companyH.UpdateCompany)
	}

	membershipsGroup := protected.Group("/memberships")
	{
		membershipsGroup.GET("", middleware.RequirePermission(perms, permissions.MembershipsView), membershipH.GetMemberships)
		membershipsGroup.GET("/:id", middleware.RequirePermission(perms, permissions.MembershipsView), membershipH.GetMembership)
		membershipsGroup.PUT("/:id", middleware.RequirePermission(perms, permissions.MembershipsUpdate), membershipH.UpdateMembership)
		membershipsGroup.DELETE("/:id", middleware.RequirePermission(perms, permissions.MembershipsDelete), membershipH.DeleteMembership)
	}

	// --- Landlord ---
	platform.GET("/companies", middleware.RequirePlatformPermission(permissions.PlatformCompaniesView), companyH.GetCompanies)
	platform.POST("/companies", middleware.RequirePlatformPermission(permissions.PlatformCompaniesManage), companyH.CreateCompany)
	platform.GET("/companies/:id", middleware.RequirePlatformPermission(permissions.PlatformCompaniesView), companyH.GetCompany)
	platform.PUT("/companies/:id", middleware.RequirePlatformPermission(permissions.PlatformCompaniesManage), companyH.UpdateCompany)
	platform.DELETE("/companies/:id", middleware.RequirePlatformPermission(permissions.PlatformCompaniesManage), companyH.DeleteCompany)

	platform.GET("/users", middleware.RequirePlatformPermission(permissions.PlatformUsersView), userH.GetUsers)
	platform.GET("/users/:id", middleware.RequirePlatformPermission(permissions.PlatformUsersView), userH.GetUser)

	platform.GET("/memberships", middleware.RequirePlatformPermission(permissions.PlatformMembershipsView), membershipH.GetMemberships)
	platform.POST("/memberships", middleware.RequirePlatformPermission(permissions.PlatformMembershipsManage), membershipH.CreateMembership)
	platform.GET("/memberships/:id", middleware.RequirePlatformPermission(permissions.PlatformMembershipsView), membershipH.GetMembership)
	platform.PUT("/memberships/:id", middleware.RequirePlatformPermission(permissions.PlatformMembershipsManage), membershipH.UpdateMembership)
	platform.DELETE("/memberships/:id", middleware.RequirePlatformPermission(permissions.PlatformMembershipsManage), membershipH.DeleteMembership)
}
