package server

import (
	"context"
	"net/http"
	"time"

	"dvra-api/internal/app/handlers"
	"dvra-api/internal/app/repositories"
	"dvra-api/internal/app/services"
	"dvra-api/internal/modules/staffing"
	"dvra-api/internal/platform/config"
	storagepkg "dvra-api/internal/platform/storage"

	_ "dvra-api/docs" // Importar documentación generada por Swagger

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
func New(cfg *config.Config, db *gorm.DB) (*Server, error) {
	// Set Gin mode based on environment
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create JWT service
	jwtService := services.NewJWTService(cfg.JWTSecret, cfg.JWTRefreshSecret)
	fileStorage, err := storagepkg.NewFilesystem(cfg.StorageRoot)
	if err != nil {
		return nil, err
	}

	// Create repositories (injecting DB connection)
	userRepo := repositories.NewUserRepository(db)
	companyRepo := repositories.NewCompanyRepository(db)
	membershipRepo := repositories.NewMembershipRepository(db)
	candidateRepo := repositories.NewCandidateRepository(db)
	applicationRepo := repositories.NewApplicationRepository(db)
	jobRepo := repositories.NewJobRepository(db)
	planRepo := repositories.NewPlanRepository(db)
	systemValueRepo := repositories.NewSystemValueRepository(db)
	locationRepo := repositories.NewLocationRepository()
	dashboardRepo := repositories.NewDashboardRepository()
	platformSettingsRepo := repositories.NewPlatformSettingsRepository(db)
	rbacService := services.NewTenantRBACService(db)

	// Create services (injecting repositories)
	authService := services.NewAuthService(userRepo, planRepo, jwtService, rbacService, db)
	userService := services.NewUserService(userRepo, db, rbacService)
	companyService := services.NewCompanyService(companyRepo)
	membershipService := services.NewMembershipService(membershipRepo, rbacService)
	candidateService := services.NewCandidateService(candidateRepo, planRepo, db, fileStorage)
	applicationService := services.NewApplicationService(applicationRepo, jobRepo, candidateRepo)
	// Módulo staffing (monolito modular + hexagonal-lite, ver ADR-001). Se cablea
	// con db + un adaptador del repo de applications hacia su puerto ApplicationFinder.
	staffingModule := staffing.New(db, staffingAppFinder{repo: applicationRepo})

	planService := services.NewPlanService(planRepo, companyRepo, db)
	jobService := services.NewJobService(jobRepo, staffingModule.ClientRepo, planService)
	systemValueService := services.NewSystemValueService(systemValueRepo)
	locationService := services.NewLocationService(locationRepo)
	dashboardService := services.NewDashboardService(dashboardRepo)
	publicService := services.NewPublicService(companyRepo, jobRepo, candidateRepo, applicationRepo, candidateService)
	platformSettingsService := services.NewPlatformSettingsService(platformSettingsRepo)

	// Create handlers (injecting services)
	healthHandler := handlers.NewHealthHandler(db)
	authHandler := handlers.NewAuthHandler(authService)
	userHandler := handlers.NewUserHandler(userService)
	companyHandler := handlers.NewCompanyHandler(companyService)
	membershipHandler := handlers.NewMembershipHandler(membershipService)
	candidateHandler := handlers.NewCandidateHandler(candidateService)
	applicationHandler := handlers.NewApplicationHandler(applicationService)
	jobHandler := handlers.NewJobHandler(jobService)
	planHandler := handlers.NewPlanHandler(planService)
	systemValueHandler := handlers.NewSystemValueHandler(systemValueService)
	locationHandler := handlers.NewLocationHandler(locationService)
	dashboardHandler := handlers.NewDashboardHandler(dashboardService)
	publicHandler := handlers.NewPublicHandler(publicService)
	platformSettingsHandler := handlers.NewPlatformSettingsHandler(platformSettingsService)
	tenantRoleHandler := handlers.NewTenantRoleHandler(rbacService)

	// Create Gin router
	router := gin.Default()

	// Configure CORS middleware
	router.Use(corsMiddleware(cfg.CorsAllowedOrigins))

	// Register routes (passing config for dynamic Swagger host)
	registerRoutes(router, healthHandler, authHandler, userHandler, companyHandler, membershipHandler, candidateHandler, applicationHandler, jobHandler, staffingModule, planHandler, planService, systemValueHandler, locationHandler, dashboardHandler, publicHandler, platformSettingsHandler, tenantRoleHandler, jwtService, rbacService, cfg)

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
	}, nil
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
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-Company-ID")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
