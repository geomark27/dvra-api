package server

import (
	"dvra-api/internal/app/handlers"
	"net/http"

	"github.com/gin-gonic/gin"
)

// registerRoutes registers all application routes
func registerRoutes(
	router *gin.Engine,
	healthHandler *handlers.HealthHandler,
	userHandler *handlers.UserHandler,
	companyHandler *handlers.CompanyHandler,
	membershipHandler *handlers.MembershipHandler,
	candidateHandler *handlers.CandidateHandler,
	applicationHandler *handlers.ApplicationHandler,
	jobHandler *handlers.JobHandler,
) {
	// Root route
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message":        "Welcome to dvra-api!",
			"status":         "success",
			"version":        "v1.1.0",
			"generated_with": "Loom",
			"endpoints": gin.H{
				"health":       "/api/v1/health",
				"users":        "/api/v1/users",
				"companies":    "/api/v1/companies",
				"memberships":  "/api/v1/memberships",
				"jobs":         "/api/v1/jobs",
				"candidates":   "/api/v1/candidates",
				"applications": "/api/v1/applications",
			},
		})
	})

	// API v1 group
	api := router.Group("/api/v1")
	{
		// Health routes
		api.GET("/health", healthHandler.Health)
		api.GET("/health/ready", healthHandler.Ready)

		// User routes
		users := api.Group("/users")
		{
			users.GET("", userHandler.GetUsers)
			users.POST("", userHandler.CreateUser)
			users.GET("/:id", userHandler.GetUser)
			users.PUT("/:id", userHandler.UpdateUser)
			users.DELETE("/:id", userHandler.DeleteUser)
		}

		// Company routes
		companies := api.Group("/companies")
		{
			companies.GET("", companyHandler.GetCompanies)
			companies.POST("", companyHandler.CreateCompany)
			companies.GET("/:id", companyHandler.GetCompany)
			companies.PUT("/:id", companyHandler.UpdateCompany)
			companies.DELETE("/:id", companyHandler.DeleteCompany)
		}

		// Membership routes
		memberships := api.Group("/memberships")
		{
			memberships.GET("", membershipHandler.GetMemberships)
			memberships.POST("", membershipHandler.CreateMembership)
			memberships.GET("/:id", membershipHandler.GetMembership)
			memberships.PUT("/:id", membershipHandler.UpdateMembership)
			memberships.DELETE("/:id", membershipHandler.DeleteMembership)
		}

		// Job routes
		jobs := api.Group("/jobs")
		{
			jobs.GET("", jobHandler.GetJobs)
			jobs.POST("", jobHandler.CreateJob)
			jobs.GET("/:id", jobHandler.GetJob)
			jobs.PUT("/:id", jobHandler.UpdateJob)
			jobs.DELETE("/:id", jobHandler.DeleteJob)
		}

		// Candidate routes
		candidates := api.Group("/candidates")
		{
			candidates.GET("", candidateHandler.GetCandidates)
			candidates.POST("", candidateHandler.CreateCandidate)
			candidates.GET("/:id", candidateHandler.GetCandidate)
			candidates.PUT("/:id", candidateHandler.UpdateCandidate)
			candidates.DELETE("/:id", candidateHandler.DeleteCandidate)
		}

		// Application routes
		applications := api.Group("/applications")
		{
			applications.GET("", applicationHandler.GetApplications)
			applications.POST("", applicationHandler.CreateApplication)
			applications.GET("/:id", applicationHandler.GetApplication)
			applications.PUT("/:id", applicationHandler.UpdateApplication)
			applications.DELETE("/:id", applicationHandler.DeleteApplication)
		}
	}
}
