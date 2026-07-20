package server

import (
	"dvra-api/internal/app/handlers"
	"dvra-api/internal/app/services"
	"dvra-api/internal/modules/staffing"
	"dvra-api/internal/platform/config"
	"dvra-api/internal/shared/middleware"
	"dvra-api/internal/shared/permissions"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// registerRoutes registers all application routes
func registerRoutes(
	router *gin.Engine,
	healthHandler *handlers.HealthHandler,
	authHandler *handlers.AuthHandler,
	userHandler *handlers.UserHandler,
	companyHandler *handlers.CompanyHandler,
	departmentHandler *handlers.DepartmentHandler,
	membershipHandler *handlers.MembershipHandler,
	candidateHandler *handlers.CandidateHandler,
	applicationHandler *handlers.ApplicationHandler,
	jobHandler *handlers.JobHandler,
	staffingModule *staffing.Module,
	planHandler *handlers.PlanHandler,
	planService services.PlanService,
	systemValueHandler *handlers.SystemValueHandler,
	locationHandler *handlers.LocationHandler,
	dashboardHandler *handlers.DashboardHandler,
	publicHandler *handlers.PublicHandler,
	platformSettingsHandler *handlers.PlatformSettingsHandler,
	companySettingsHandler *handlers.CompanySettingsHandler,
	tenantRoleHandler *handlers.TenantRoleHandler,
	jwtService services.JWTService,
	rbac services.TenantAuthorizer,
	cfg *config.Config,
) {
	// Root route
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message":        "Welcome to dvra-api!",
			"status":         "success",
			"version":        "v1.4.0",
			"generated_with": "Loom",
			"endpoints": gin.H{
				"health":           "/api/v1/health",
				"auth":             "/api/v1/auth",
				"users":            "/api/v1/users",
				"companies":        "/api/v1/companies",
				"departments":      "/api/v1/departments",
				"memberships":      "/api/v1/memberships",
				"jobs":             "/api/v1/jobs",
				"staffing_clients": "/api/v1/staffing-clients",
				"placements":       "/api/v1/placements",
				"candidates":       "/api/v1/candidates",
				"applications":     "/api/v1/applications",
				"dashboard":        "/api/v1/dashboard",
				"plans":            "/api/v1/plans (public)",
				"locations":        "/api/v1/locations (public)",
				"public":           "/api/v1/public (career page)",
				"swagger":          "/swagger/index.html",
			},
		})
	})

	// Swagger documentation with dynamic host
	swaggerURL := ginSwagger.URL(fmt.Sprintf("http://localhost:%s/swagger/doc.json", cfg.Port))
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, swaggerURL))

	// API v1 group
	api := router.Group("/api/v1")
	{
		// Health routes (public)
		api.GET("/health", healthHandler.Health)
		api.GET("/health/ready", healthHandler.Ready)

		// Auth routes (public)
		auth := api.Group("/auth")
		{
			// Public routes
			auth.POST("/register-company", authHandler.RegisterCompany)
			auth.POST("/register", authHandler.Register) // Deprecated: use register-company
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)

			// Protected auth routes
			authProtected := auth.Group("")
			authProtected.Use(middleware.AuthMiddleware(jwtService, rbac))
			{
				authProtected.GET("/me", authHandler.GetMe)
				authProtected.POST("/change-password", authHandler.ChangePassword)
				authProtected.POST("/logout", authHandler.Logout)
				authProtected.POST("/switch-company", authHandler.SwitchCompany)
				authProtected.POST("/switch-platform", authHandler.SwitchToPlatform)
				authProtected.GET("/my-companies", authHandler.GetMyCompanies)
			}
		}

		// Protected routes (require authentication)
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(jwtService, rbac))
		{
			protected.GET("/tenant-roles", middleware.RequirePermission(rbac, permissions.UsersView), tenantRoleHandler.GetTenantRoles)
			protected.GET("/settings", middleware.RequirePermission(rbac, permissions.SettingsView), companySettingsHandler.Get)
			protected.PUT("/settings", middleware.RequirePermission(rbac, permissions.SettingsUpdate), companySettingsHandler.Update)
			protected.GET("/subscription", middleware.RequirePermission(rbac, permissions.BillingView), companySettingsHandler.Subscription)
			protected.GET("/subscription/usage", middleware.RequirePermission(rbac, permissions.BillingView), companySettingsHandler.Usage)

			// User routes
			users := protected.Group("/users")
			{
				users.GET("", middleware.RequirePermission(rbac, permissions.UsersView), userHandler.GetUsers)
				users.POST("", middleware.RequirePermission(rbac, permissions.UsersCreate), userHandler.CreateUser)
				users.GET("/:id", middleware.RequirePermission(rbac, permissions.UsersView), userHandler.GetUser)
				users.DELETE("/:id", middleware.RequirePermission(rbac, permissions.UsersDelete), userHandler.DeleteUser)
			}

			// Company routes
			companies := protected.Group("/companies")
			{
				companies.GET("", middleware.RequirePermission(rbac, permissions.CompaniesView), companyHandler.GetCompanies)
				companies.GET("/:id", middleware.RequirePermission(rbac, permissions.CompaniesView), companyHandler.GetCompany)
				companies.PUT("/:id", middleware.RequirePermission(rbac, permissions.CompaniesUpdate), companyHandler.UpdateCompany)
			}

			// Membership routes
			memberships := protected.Group("/memberships")
			{
				memberships.GET("", middleware.RequirePermission(rbac, permissions.MembershipsView), membershipHandler.GetMemberships)
				memberships.GET("/:id", middleware.RequirePermission(rbac, permissions.MembershipsView), membershipHandler.GetMembership)
				memberships.PUT("/:id", middleware.RequirePermission(rbac, permissions.MembershipsUpdate), membershipHandler.UpdateMembership)
				memberships.DELETE("/:id", middleware.RequirePermission(rbac, permissions.MembershipsDelete), membershipHandler.DeleteMembership)
			}

			departments := protected.Group("/departments")
			{
				departments.GET("", middleware.RequirePermission(rbac, permissions.DepartmentsView), departmentHandler.ListDepartments)
				departments.POST("", middleware.RequirePermission(rbac, permissions.DepartmentsCreate), departmentHandler.CreateDepartment)
				departments.GET("/tree", middleware.RequirePermission(rbac, permissions.DepartmentsView), departmentHandler.GetDepartmentTree)
				departments.GET("/:id", middleware.RequirePermission(rbac, permissions.DepartmentsView), departmentHandler.GetDepartment)
				departments.PUT("/:id", middleware.RequirePermission(rbac, permissions.DepartmentsUpdate), departmentHandler.UpdateDepartment)
				departments.DELETE("/:id", middleware.RequirePermission(rbac, permissions.DepartmentsDelete), departmentHandler.DeactivateDepartment)
			}

			// Job routes
			jobs := protected.Group("/jobs")
			{
				jobs.GET("", middleware.RequirePermission(rbac, permissions.JobsView), jobHandler.GetJobs)
				jobs.POST("", middleware.RequirePermission(rbac, permissions.JobsCreate), jobHandler.CreateJob)
				jobs.GET("/:id", middleware.RequirePermission(rbac, permissions.JobsView), jobHandler.GetJob)
				jobs.PUT("/:id", middleware.RequirePermission(rbac, permissions.JobsUpdate), jobHandler.UpdateJob)
				jobs.DELETE("/:id", middleware.RequirePermission(rbac, permissions.JobsDelete), jobHandler.DeleteJob)
				jobs.PATCH("/:id/publish", middleware.RequirePermission(rbac, permissions.JobsPublish), jobHandler.PublishJob)
				jobs.PATCH("/:id/close", middleware.RequirePermission(rbac, permissions.JobsClose), jobHandler.CloseJob)
			}

			// Módulo staffing (monolito modular, ver ADR-001): registra sus propias
			// rutas (/staffing-clients y /placements), gateadas por plan internamente.
			staffingModule.RegisterRoutes(protected, planService, rbac)

			// Candidate routes
			candidates := protected.Group("/candidates")
			{
				candidates.GET("", middleware.RequirePermission(rbac, permissions.CandidatesView), candidateHandler.GetCandidates)
				candidates.POST("", middleware.RequirePermission(rbac, permissions.CandidatesCreate), candidateHandler.CreateCandidate)
				candidates.GET("/:id", middleware.RequirePermission(rbac, permissions.CandidatesView), candidateHandler.GetCandidate)
				candidates.PUT("/:id", middleware.RequirePermission(rbac, permissions.CandidatesUpdate), candidateHandler.UpdateCandidate)
				candidates.DELETE("/:id", middleware.RequirePermission(rbac, permissions.CandidatesDelete), candidateHandler.DeleteCandidate)
				candidates.POST("/:id/upload-resume", middleware.RequirePermission(rbac, permissions.CandidatesUploadResume), candidateHandler.UploadResume)
				candidates.GET("/:id/resume", middleware.RequirePermission(rbac, permissions.CandidatesView), candidateHandler.DownloadResume)
			}

			// Application routes
			applications := protected.Group("/applications")
			{
				applications.GET("", middleware.RequirePermission(rbac, permissions.ApplicationsView), applicationHandler.GetApplications)
				applications.GET("/by-stage", middleware.RequirePermission(rbac, permissions.ApplicationsView), applicationHandler.GetApplicationsByStage)
				applications.POST("", middleware.RequirePermission(rbac, permissions.ApplicationsCreate), applicationHandler.CreateApplication)
				applications.GET("/:id", middleware.RequirePermission(rbac, permissions.ApplicationsView), applicationHandler.GetApplication)
				applications.PUT("/:id", middleware.RequirePermission(rbac, permissions.ApplicationsUpdate), applicationHandler.UpdateApplication)
				applications.PATCH("/:id/move", middleware.RequirePermission(rbac, permissions.ApplicationsMove), applicationHandler.MoveApplication)
				applications.PATCH("/:id/rate", middleware.RequirePermission(rbac, permissions.ApplicationsRate), applicationHandler.RateApplication)
				applications.DELETE("/:id", middleware.RequirePermission(rbac, permissions.ApplicationsDelete), applicationHandler.DeleteApplication)
			}

			// Dashboard routes
			dashboard := protected.Group("/dashboard")
			{
				dashboard.GET("/stats", middleware.RequirePermission(rbac, permissions.DashboardView), dashboardHandler.GetStats)
			}

			// Tenant catalog routes. The company scope always comes from the JWT.
			systemValues := protected.Group("/system-values")
			{
				systemValues.GET("", middleware.RequirePermission(rbac, permissions.CatalogsView), systemValueHandler.GetCategories)
				systemValues.GET("/:category", middleware.RequirePermission(rbac, permissions.CatalogsView), systemValueHandler.GetByCategory)
				systemValues.GET("/modules/:module/:category", middleware.RequirePermission(rbac, permissions.CatalogsView), systemValueHandler.GetByModuleCategory)
				systemValues.POST("", middleware.RequirePermission(rbac, permissions.CatalogsManage), systemValueHandler.CreateTenant)
				systemValues.PUT("/:id", middleware.RequirePermission(rbac, permissions.CatalogsManage), systemValueHandler.UpdateTenant)
				systemValues.DELETE("/:id", middleware.RequirePermission(rbac, permissions.CatalogsManage), systemValueHandler.DeleteTenant)
				systemValues.PUT("/:id/override", middleware.RequirePermission(rbac, permissions.CatalogsManage), systemValueHandler.UpsertTenantOverride)
			}
		}

		// Landlord routes use an explicit platform context and a permission matrix
		// separate from tenant roles.
		platform := api.Group("/platform")
		platform.Use(middleware.AuthMiddleware(jwtService, rbac))
		{
			platform.GET("/companies", middleware.RequirePlatformPermission(permissions.PlatformCompaniesView), companyHandler.GetCompanies)
			platform.POST("/companies", middleware.RequirePlatformPermission(permissions.PlatformCompaniesManage), companyHandler.CreateCompany)
			platform.GET("/companies/:id", middleware.RequirePlatformPermission(permissions.PlatformCompaniesView), companyHandler.GetCompany)
			platform.PUT("/companies/:id", middleware.RequirePlatformPermission(permissions.PlatformCompaniesManage), companyHandler.UpdateCompany)
			platform.DELETE("/companies/:id", middleware.RequirePlatformPermission(permissions.PlatformCompaniesManage), companyHandler.DeleteCompany)

			platform.GET("/users", middleware.RequirePlatformPermission(permissions.PlatformUsersView), userHandler.GetUsers)
			platform.GET("/users/:id", middleware.RequirePlatformPermission(permissions.PlatformUsersView), userHandler.GetUser)

			platform.GET("/memberships", middleware.RequirePlatformPermission(permissions.PlatformMembershipsView), membershipHandler.GetMemberships)
			platform.POST("/memberships", middleware.RequirePlatformPermission(permissions.PlatformMembershipsManage), membershipHandler.CreateMembership)
			platform.GET("/memberships/:id", middleware.RequirePlatformPermission(permissions.PlatformMembershipsView), membershipHandler.GetMembership)
			platform.PUT("/memberships/:id", middleware.RequirePlatformPermission(permissions.PlatformMembershipsManage), membershipHandler.UpdateMembership)
			platform.DELETE("/memberships/:id", middleware.RequirePlatformPermission(permissions.PlatformMembershipsManage), membershipHandler.DeleteMembership)

			platform.GET("/plans", middleware.RequirePlatformPermission(permissions.PlatformPlansManage), planHandler.GetPlans)
			platform.POST("/plans", middleware.RequirePlatformPermission(permissions.PlatformPlansManage), planHandler.CreatePlan)
			platform.GET("/plans/:id", middleware.RequirePlatformPermission(permissions.PlatformPlansManage), planHandler.GetPlanByID)
			platform.PUT("/plans/:id", middleware.RequirePlatformPermission(permissions.PlatformPlansManage), planHandler.UpdatePlan)
			platform.PATCH("/plans/:id/toggle", middleware.RequirePlatformPermission(permissions.PlatformPlansManage), planHandler.TogglePlanStatus)
			platform.DELETE("/plans/:id", middleware.RequirePlatformPermission(permissions.PlatformPlansManage), planHandler.DeletePlan)
			platform.POST("/plans/assign", middleware.RequirePlatformPermission(permissions.PlatformPlansManage), planHandler.AssignPlanToCompany)

			platform.GET("/settings", middleware.RequirePlatformPermission(permissions.PlatformSettingsManage), platformSettingsHandler.GetFullSettings)
			platform.PUT("/settings", middleware.RequirePlatformPermission(permissions.PlatformSettingsManage), platformSettingsHandler.UpdateSettings)

			platform.GET("/system-values", middleware.RequirePlatformPermission(permissions.PlatformCatalogsManage), systemValueHandler.GetAll)
			platform.POST("/system-values", middleware.RequirePlatformPermission(permissions.PlatformCatalogsManage), systemValueHandler.Create)
			platform.PUT("/system-values/:id", middleware.RequirePlatformPermission(permissions.PlatformCatalogsManage), systemValueHandler.Update)
			platform.DELETE("/system-values/:id", middleware.RequirePlatformPermission(permissions.PlatformCatalogsManage), systemValueHandler.Delete)
		}

		// Public Plans routes (no auth required)
		plans := api.Group("/plans")
		{
			plans.GET("", planHandler.GetPublicPlans)
			plans.GET("/:slug", planHandler.GetPlanBySlug)
		}

		// PUBLIC CAREER PAGE ROUTES (no auth required)
		public := api.Group("/public")
		{
			// Platform settings (public - for branding/config)
			public.GET("/platform-settings", platformSettingsHandler.GetPublicSettings)

			// Company info
			public.GET("/companies/:slug", publicHandler.GetCompanyBySlug)
			public.GET("/companies/:slug/jobs", publicHandler.GetPublishedJobsByCompany)

			// Job details and applications
			public.GET("/jobs/:id", publicHandler.GetPublishedJobByID)
			public.POST("/jobs/:id/apply", publicHandler.ApplyToJob)
		}

		// Public Location routes (no auth required - READ ONLY)
		locations := api.Group("/locations")
		{
			// Regions
			locations.GET("/regions", locationHandler.GetAllRegions)
			locations.GET("/regions/:id", locationHandler.GetRegionByID)

			// Subregions
			locations.GET("/subregions", locationHandler.GetAllSubregions)
			locations.GET("/subregions/:id", locationHandler.GetSubregionByID)

			// Countries
			locations.GET("/countries", locationHandler.GetAllCountries)
			locations.GET("/countries/:id", locationHandler.GetCountryByID)
			locations.GET("/countries/iso/:iso", locationHandler.GetCountryByISO)

			// States
			locations.GET("/states", locationHandler.GetAllStates)
			locations.GET("/states/:id", locationHandler.GetStateByID)

			// Cities
			locations.GET("/cities", locationHandler.GetAllCities)
			locations.GET("/cities/:id", locationHandler.GetCityByID)

			// Helpers
			locations.GET("/hierarchy/:id", locationHandler.GetLocationHierarchy)
			locations.GET("/search", locationHandler.SearchLocations)
		}
	}
}
