package server

import (
	"fmt"
	"net/http"

	"dvra-api/internal/platform/config"
	"dvra-api/internal/platform/httpx"
	"dvra-api/internal/shared/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// registerRoutes arma los tres grupos de rutas y deja que cada módulo monte las
// suyas.
//
// Ya no monta ningún endpoint de dominio: tras la migración por dominios, esto
// es solo el esqueleto. Añadir un módulo cuesta una línea en buildDeps.
func registerRoutes(router *gin.Engine, d *deps, cfg *config.Config) {
	router.GET("/", rootHandler)

	// Swagger con host dinámico
	swaggerURL := ginSwagger.URL(fmt.Sprintf("http://localhost:%s/swagger/doc.json", cfg.Port))
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, swaggerURL))

	api := router.Group("/api/v1")

	// Grupo protegido: exige sesión válida y resuelve el rol vigente en cada
	// petición.
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware(d.SessionAuth))

	// Grupo de landlord: misma autenticación, pero sus rutas autorizan contra la
	// matriz de plataforma, no contra el rol de empresa (ADR-002).
	platform := api.Group("/platform")
	platform.Use(middleware.AuthMiddleware(d.SessionAuth))

	moduleRouter := httpx.Router{
		Public:    api,
		Protected: protected,
		Platform:  platform,
		RBAC:      d.RBAC,
		Plans:     d.Plans,
	}
	for _, module := range d.Modules {
		module.RegisterRoutes(moduleRouter)
	}
}

// rootHandler describe la API en su raíz.
func rootHandler(c *gin.Context) {
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
}
