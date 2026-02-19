package server

import (
	"dvra-api/internal/app/handlers"
	"dvra-api/internal/app/services"
	"dvra-api/internal/platform/config"
	"dvra-api/internal/shared/middleware"
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
	planHandler *handlers.PlanHandler,
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
				"health":       "/api/v1/health",
				"auth":         "/api/v1/auth",
				"users":        "/api/v1/users",
				"companies":    "/api/v1/companies",
				"memberships":  "/api/v1/memberships",
				"jobs":         "/api/v1/jobs",
				"candidates":   "/api/v1/candidates",
				"applications": "/api/v1/applications",
				"dashboard":    "/api/v1/dashboard",
				"plans":        "/api/v1/plans (public)",
				"locations":    "/api/v1/locations (public)",
				"public":       "/api/v1/public (career page)",
				"swagger":      "/swagger/index.html",
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
				users.GET("", userHandler.GetUsers)
				users.POST("", userHandler.CreateUser)
				users.GET("/:id", userHandler.GetUser)
				users.PUT("/:id", userHandler.UpdateUser)
				users.DELETE("/:id", userHandler.DeleteUser)
			}

			// Company routes
			companies := protected.Group("/companies")
			{
				companies.GET("", companyHandler.GetCompanies)
				companies.POST("", companyHandler.CreateCompany)
				companies.GET("/:id", companyHandler.GetCompany)
				companies.PUT("/:id", companyHandler.UpdateCompany)
				companies.DELETE("/:id", companyHandler.DeleteCompany)
			}

			// Membership routes
			memberships := protected.Group("/memberships")
			{
				memberships.GET("", membershipHandler.GetMemberships)
				memberships.POST("", membershipHandler.CreateMembership)
				memberships.GET("/:id", membershipHandler.GetMembership)
				memberships.PUT("/:id", membershipHandler.UpdateMembership)
				memberships.DELETE("/:id", membershipHandler.DeleteMembership)
			}

			// Job routes
			jobs := protected.Group("/jobs")
			{
				jobs.GET("", jobHandler.GetJobs)
				jobs.POST("", jobHandler.CreateJob)
				jobs.GET("/:id", jobHandler.GetJob)
				jobs.PUT("/:id", jobHandler.UpdateJob)
				jobs.DELETE("/:id", jobHandler.DeleteJob)
				jobs.PATCH("/:id/publish", jobHandler.PublishJob)
				jobs.PATCH("/:id/close", jobHandler.CloseJob)
			}

			// Candidate routes
			candidates := protected.Group("/candidates")
			{
				candidates.GET("", candidateHandler.GetCandidates)
				candidates.POST("", candidateHandler.CreateCandidate)
				candidates.GET("/:id", candidateHandler.GetCandidate)
				candidates.PUT("/:id", candidateHandler.UpdateCandidate)
				candidates.DELETE("/:id", candidateHandler.DeleteCandidate)
				candidates.POST("/:id/upload-resume", candidateHandler.UploadResume)
			}

			// Application routes
			applications := protected.Group("/applications")
			{
				applications.GET("", applicationHandler.GetApplications)
				applications.GET("/by-stage", applicationHandler.GetApplicationsByStage)
				applications.POST("", applicationHandler.CreateApplication)
				applications.GET("/:id", applicationHandler.GetApplication)
				applications.PUT("/:id", applicationHandler.UpdateApplication)
				applications.PATCH("/:id/move", applicationHandler.MoveApplication)
				applications.PATCH("/:id/rate", applicationHandler.RateApplication)
				applications.DELETE("/:id", applicationHandler.DeleteApplication)
			}

			// Dashboard routes
			dashboard := protected.Group("/dashboard")
			{
				dashboard.GET("/stats", dashboardHandler.GetStats)
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
