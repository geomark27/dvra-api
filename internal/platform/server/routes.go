package server

import (
	"dvra-api/internal/app/handlers"
	"dvra-api/internal/app/services"
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
	membershipHandler *handlers.MembershipHandler,
	candidateHandler *handlers.CandidateHandler,
	applicationHandler *handlers.ApplicationHandler,
	jobHandler *handlers.JobHandler,
	staffingClientHandler *handlers.StaffingClientHandler,
	placementHandler *handlers.PlacementHandler,
	planHandler *handlers.PlanHandler,
	planService services.PlanService,
	systemValueHandler *handlers.SystemValueHandler,
	locationHandler *handlers.LocationHandler,
	dashboardHandler *handlers.DashboardHandler,
	publicHandler *handlers.PublicHandler,
	platformSettingsHandler *handlers.PlatformSettingsHandler,
	jwtService services.JWTService,
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
			authProtected.Use(middleware.AuthMiddleware(jwtService))
			{
				authProtected.GET("/me", authHandler.GetMe)
				authProtected.POST("/change-password", authHandler.ChangePassword)
				authProtected.POST("/logout", authHandler.Logout)
				authProtected.POST("/switch-company", authHandler.SwitchCompany)
				authProtected.GET("/my-companies", authHandler.GetMyCompanies)
			}
		}

		// Protected routes (require authentication)
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(jwtService))
		{
			// User routes
			users := protected.Group("/users")
			{
				users.GET("", middleware.RequirePermission(permissions.UsersView), userHandler.GetUsers)
				users.POST("", middleware.RequirePermission(permissions.UsersCreate), userHandler.CreateUser)
				users.GET("/:id", middleware.RequirePermission(permissions.UsersView), userHandler.GetUser)
				users.PUT("/:id", middleware.RequirePermission(permissions.UsersUpdate), userHandler.UpdateUser)
				users.DELETE("/:id", middleware.RequirePermission(permissions.UsersDelete), userHandler.DeleteUser)
			}

			// Company routes
			companies := protected.Group("/companies")
			{
				companies.GET("", middleware.RequirePermission(permissions.CompaniesView), companyHandler.GetCompanies)
				companies.POST("", middleware.RequirePermission(permissions.CompaniesCreate), companyHandler.CreateCompany)
				companies.GET("/:id", middleware.RequirePermission(permissions.CompaniesView), companyHandler.GetCompany)
				companies.PUT("/:id", middleware.RequirePermission(permissions.CompaniesUpdate), companyHandler.UpdateCompany)
				companies.DELETE("/:id", middleware.RequirePermission(permissions.CompaniesDelete), companyHandler.DeleteCompany)
			}

			// Membership routes
			memberships := protected.Group("/memberships")
			{
				memberships.GET("", middleware.RequirePermission(permissions.MembershipsView), membershipHandler.GetMemberships)
				memberships.POST("", middleware.RequirePermission(permissions.MembershipsCreate), membershipHandler.CreateMembership)
				memberships.GET("/:id", middleware.RequirePermission(permissions.MembershipsView), membershipHandler.GetMembership)
				memberships.PUT("/:id", middleware.RequirePermission(permissions.MembershipsUpdate), membershipHandler.UpdateMembership)
				memberships.DELETE("/:id", middleware.RequirePermission(permissions.MembershipsDelete), membershipHandler.DeleteMembership)
			}

			// Job routes
			jobs := protected.Group("/jobs")
			{
				jobs.GET("", middleware.RequirePermission(permissions.JobsView), jobHandler.GetJobs)
				jobs.POST("", middleware.RequirePermission(permissions.JobsCreate), jobHandler.CreateJob)
				jobs.GET("/:id", middleware.RequirePermission(permissions.JobsView), jobHandler.GetJob)
				jobs.PUT("/:id", middleware.RequirePermission(permissions.JobsUpdate), jobHandler.UpdateJob)
				jobs.DELETE("/:id", middleware.RequirePermission(permissions.JobsDelete), jobHandler.DeleteJob)
				jobs.PATCH("/:id/publish", middleware.RequirePermission(permissions.JobsPublish), jobHandler.PublishJob)
				jobs.PATCH("/:id/close", middleware.RequirePermission(permissions.JobsClose), jobHandler.CloseJob)
			}

			// Staffing client routes (gateado por plan vía RequireFeature)
			staffingClients := protected.Group("/staffing-clients")
			staffingClients.Use(middleware.RequireFeature(planService, "staffing"))
			{
				staffingClients.GET("", middleware.RequirePermission(permissions.StaffingClientsView), staffingClientHandler.GetStaffingClients)
				staffingClients.POST("", middleware.RequirePermission(permissions.StaffingClientsCreate), staffingClientHandler.CreateStaffingClient)
				staffingClients.GET("/:id", middleware.RequirePermission(permissions.StaffingClientsView), staffingClientHandler.GetStaffingClient)
				staffingClients.PUT("/:id", middleware.RequirePermission(permissions.StaffingClientsUpdate), staffingClientHandler.UpdateStaffingClient)
				staffingClients.DELETE("/:id", middleware.RequirePermission(permissions.StaffingClientsDelete), staffingClientHandler.DeleteStaffingClient)
			}

			// Placement routes (gateado por plan vía RequireFeature)
			placements := protected.Group("/placements")
			placements.Use(middleware.RequireFeature(planService, "staffing"))
			{
				placements.GET("", middleware.RequirePermission(permissions.PlacementsView), placementHandler.GetPlacements)
				placements.POST("", middleware.RequirePermission(permissions.PlacementsCreate), placementHandler.CreatePlacement)
				placements.GET("/:id", middleware.RequirePermission(permissions.PlacementsView), placementHandler.GetPlacement)
				placements.PUT("/:id", middleware.RequirePermission(permissions.PlacementsUpdate), placementHandler.UpdatePlacement)
				placements.DELETE("/:id", middleware.RequirePermission(permissions.PlacementsDelete), placementHandler.DeletePlacement)
			}

			// Candidate routes
			candidates := protected.Group("/candidates")
			{
				candidates.GET("", middleware.RequirePermission(permissions.CandidatesView), candidateHandler.GetCandidates)
				candidates.POST("", middleware.RequirePermission(permissions.CandidatesCreate), candidateHandler.CreateCandidate)
				candidates.GET("/:id", middleware.RequirePermission(permissions.CandidatesView), candidateHandler.GetCandidate)
				candidates.PUT("/:id", middleware.RequirePermission(permissions.CandidatesUpdate), candidateHandler.UpdateCandidate)
				candidates.DELETE("/:id", middleware.RequirePermission(permissions.CandidatesDelete), candidateHandler.DeleteCandidate)
				candidates.POST("/:id/upload-resume", middleware.RequirePermission(permissions.CandidatesUploadResume), candidateHandler.UploadResume)
			}

			// Application routes
			applications := protected.Group("/applications")
			{
				applications.GET("", middleware.RequirePermission(permissions.ApplicationsView), applicationHandler.GetApplications)
				applications.GET("/by-stage", middleware.RequirePermission(permissions.ApplicationsView), applicationHandler.GetApplicationsByStage)
				applications.POST("", middleware.RequirePermission(permissions.ApplicationsCreate), applicationHandler.CreateApplication)
				applications.GET("/:id", middleware.RequirePermission(permissions.ApplicationsView), applicationHandler.GetApplication)
				applications.PUT("/:id", middleware.RequirePermission(permissions.ApplicationsUpdate), applicationHandler.UpdateApplication)
				applications.PATCH("/:id/move", middleware.RequirePermission(permissions.ApplicationsMove), applicationHandler.MoveApplication)
				applications.PATCH("/:id/rate", middleware.RequirePermission(permissions.ApplicationsRate), applicationHandler.RateApplication)
				applications.DELETE("/:id", middleware.RequirePermission(permissions.ApplicationsDelete), applicationHandler.DeleteApplication)
			}

			// Dashboard routes
			dashboard := protected.Group("/dashboard")
			{
				dashboard.GET("/stats", middleware.RequirePermission(permissions.DashboardView), dashboardHandler.GetStats)
			}

			// System Values routes (read only)
			systemValues := api.Group("/system-values")
			{
				systemValues.GET("/:category", systemValueHandler.GetByCategory)
			}
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
