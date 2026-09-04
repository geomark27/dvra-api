package server

import (
	"context"
	"net/http"
	"time"

	"dvra-api/internal/platform/config"

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

	// Composition root: repositorios -> servicios -> módulos -> handlers.
	d, err := buildDeps(cfg, db)
	if err != nil {
		return nil, err
	}

	// Create Gin router
	router := gin.Default()

	// Configure CORS middleware
	router.Use(corsMiddleware(cfg.CorsAllowedOrigins))

	// Register routes (passing config for dynamic Swagger host)
	registerRoutes(router, d, cfg)

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
