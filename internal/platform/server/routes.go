package server

import (
	"dvra-api/internal/app/handlers"
	adminHandlers "dvra-api/internal/app/handlers/admin"
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
	superAdminHandler *adminHandlers.SuperAdminCompaniesHandler,
	jwtService services.JWTService,
	cfg *config.Config,
) {
	// Root route
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message":        "Welcome to dvra-api!",
			"status":         "success",
			"version":        "v1.2.0",
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
				"plans":        "/api/v1/plans (public)",
				"locations":    "/api/v1/locations (public)",
				"admin":        "/api/v1/admin (SuperAdmin only)",
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

			// SuperAdmin login (separate endpoint)
			auth.POST("/superadmin/login", authHandler.SuperAdminLogin)

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

			// Membership routes (READ-ONLY for clients)
			memberships := protected.Group("/memberships")
			{
				memberships.GET("", membershipHandler.GetMemberships)          // Ver memberships de mi empresa
				memberships.GET("/:id", membershipHandler.GetMembership)       // Ver detalle
				memberships.PUT("/:id", membershipHandler.UpdateMembership)    // Actualizar roles
				memberships.DELETE("/:id", membershipHandler.DeleteMembership) // Remover de empresa
			}

			// Job routes
			jobs := protected.Group("/jobs")
			{
				jobs.GET("", jobHandler.GetJobs)
				jobs.POST("", jobHandler.CreateJob)
				jobs.GET("/:id", jobHandler.GetJob)
				jobs.PUT("/:id", jobHandler.UpdateJob)
				jobs.DELETE("/:id", jobHandler.DeleteJob)
			}

			// Candidate routes
			candidates := protected.Group("/candidates")
			{
				candidates.GET("", candidateHandler.GetCandidates)
				candidates.POST("", candidateHandler.CreateCandidate)
				candidates.GET("/:id", candidateHandler.GetCandidate)
				candidates.PUT("/:id", candidateHandler.UpdateCandidate)
				candidates.DELETE("/:id", candidateHandler.DeleteCandidate)
			}

			// Application routes
			applications := protected.Group("/applications")
			{
				applications.GET("", applicationHandler.GetApplications)
				applications.POST("", applicationHandler.CreateApplication)
				applications.GET("/:id", applicationHandler.GetApplication)
				applications.PUT("/:id", applicationHandler.UpdateApplication)
				applications.DELETE("/:id", applicationHandler.DeleteApplication)
			}

			// System Values routes (public - read only)
			systemValues := api.Group("/system-values")
			{
				systemValues.GET("/:category", systemValueHandler.GetByCategory) // Get values by category
			}
		}

		// Public Plans routes (no auth required)
		plans := api.Group("/plans")
		{
			plans.GET("", planHandler.GetPublicPlans)      // Public pricing page
			plans.GET("/:slug", planHandler.GetPlanBySlug) // Get plan details by slug
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

		// SuperAdmin routes (Global - No company scope required)
		admin := api.Group("/admin")
		admin.Use(middleware.AuthMiddleware(jwtService))
		admin.Use(middleware.RequireSuperAdmin())
		{
			// Plan management (SuperAdmin only)
			adminPlans := admin.Group("/plans")
			{
				adminPlans.GET("", planHandler.GetPlans)
				adminPlans.POST("", planHandler.CreatePlan)
				adminPlans.GET("/:id", planHandler.GetPlanByID)
				adminPlans.PUT("/:id", planHandler.UpdatePlan)
				adminPlans.PATCH("/:id/toggle", planHandler.TogglePlanStatus)
				adminPlans.DELETE("/:id", planHandler.DeletePlan)
				adminPlans.POST("/assign", planHandler.AssignPlanToCompany)
			}

			// Company management
			admin.GET("/companies", superAdminHandler.GetAllCompanies)
			admin.POST("/companies", superAdminHandler.CreateCompany)
			admin.PUT("/companies/:id/plan", superAdminHandler.ChangeCompanyPlan)
			admin.POST("/companies/:id/suspend", superAdminHandler.SuspendCompany)
			admin.GET("/companies/:id/users", superAdminHandler.GetCompanyUsers)

			// Membership management (Assign users to companies)
			adminMemberships := admin.Group("/memberships")
			{
				adminMemberships.POST("", membershipHandler.CreateMembership) // Asignar usuario a empresa
			}

			// System Values management (CRUD for SuperAdmin)
			adminSystemValues := admin.Group("/system-values")
			{
				adminSystemValues.GET("", systemValueHandler.GetAll)
				adminSystemValues.POST("", systemValueHandler.Create)
				adminSystemValues.PUT("/:id", systemValueHandler.Update)
				adminSystemValues.DELETE("/:id", systemValueHandler.Delete)
			}

			// Location management (CRUD for SuperAdmin)
			adminLocations := admin.Group("/locations")
			{
				// Regions
				adminLocations.POST("/regions", locationHandler.CreateRegion)
				adminLocations.PUT("/regions/:id", locationHandler.UpdateRegion)
				adminLocations.DELETE("/regions/:id", locationHandler.DeleteRegion)

				// Subregions
				adminLocations.POST("/subregions", locationHandler.CreateSubregion)
				adminLocations.PUT("/subregions/:id", locationHandler.UpdateSubregion)
				adminLocations.DELETE("/subregions/:id", locationHandler.DeleteSubregion)

				// Countries
				adminLocations.POST("/countries", locationHandler.CreateCountry)
				adminLocations.PUT("/countries/:id", locationHandler.UpdateCountry)
				adminLocations.DELETE("/countries/:id", locationHandler.DeleteCountry)

				// States
				adminLocations.POST("/states", locationHandler.CreateState)
				adminLocations.PUT("/states/:id", locationHandler.UpdateState)
				adminLocations.DELETE("/states/:id", locationHandler.DeleteState)

				// Cities
				adminLocations.POST("/cities", locationHandler.CreateCity)
				adminLocations.PUT("/cities/:id", locationHandler.UpdateCity)
				adminLocations.DELETE("/cities/:id", locationHandler.DeleteCity)
			}

			// Analytics and reports
			admin.GET("/analytics", superAdminHandler.GetGlobalAnalytics)
		}
	}
}
