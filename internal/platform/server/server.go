package server

import (
	"context"
	"net/http"
	"time"

	"dvra-api/internal/app/handlers"
	"dvra-api/internal/app/repositories"
	"dvra-api/internal/app/services"
	"dvra-api/internal/platform/config"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Server represents the HTTP server
type Server struct {
	config     *config.Config
	router     *gin.Engine
	httpServer *http.Server
	db         *gorm.DB
}

// New creates a new server instance with all dependencies injected
func New(cfg *config.Config, db *gorm.DB) *Server {
	// Set Gin mode based on environment
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create repositories (injecting DB connection)
	userRepo := repositories.NewUserRepository()
	companyRepo := repositories.NewCompanyRepository()
	membershipRepo := repositories.NewMembershipRepository()
	candidateRepo := repositories.NewCandidateRepository()
	applicationRepo := repositories.NewApplicationRepository()
	jobRepo := repositories.NewJobRepository()

	// Create services (injecting repositories)
	userService := services.NewUserService(userRepo)
	companyService := services.NewCompanyService(companyRepo)
	membershipService := services.NewMembershipService(membershipRepo)
	candidateService := services.NewCandidateService(candidateRepo)
	applicationService := services.NewApplicationService(applicationRepo)
	jobService := services.NewJobService(jobRepo)

	// Create handlers (injecting services)
	healthHandler := handlers.NewHealthHandler()
	userHandler := handlers.NewUserHandler(userService)
	companyHandler := handlers.NewCompanyHandler(companyService)
	membershipHandler := handlers.NewMembershipHandler(membershipService)
	candidateHandler := handlers.NewCandidateHandler(candidateService)
	applicationHandler := handlers.NewApplicationHandler(applicationService)
	jobHandler := handlers.NewJobHandler(jobService)

	// Create Gin router
	router := gin.Default()

	// Configure CORS middleware
	router.Use(corsMiddleware(cfg.CorsAllowedOrigins))

	// Register routes
	registerRoutes(router, healthHandler, userHandler, companyHandler, membershipHandler, candidateHandler, applicationHandler, jobHandler)

	// Configure HTTP server
	httpServer := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &Server{
		config:     cfg,
		router:     router,
		httpServer: httpServer,
		db:         db,
	}
}

// Start starts the HTTP server
func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

// corsMiddleware returns a Gin middleware for CORS
func corsMiddleware(allowedOrigins []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		// Check if origin is allowed
		allowed := false
		for _, allowedOrigin := range allowedOrigins {
			if allowedOrigin == "*" || allowedOrigin == origin {
				allowed = true
				break
			}
		}

		if allowed {
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
