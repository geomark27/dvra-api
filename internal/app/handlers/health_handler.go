package handlers

import (
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"

	"gorm.io/gorm"
)

// HealthHandler handles health check routes
type HealthHandler struct{ db *gorm.DB }

// NewHealthHandler creates a new HealthHandler instance
func NewHealthHandler(db *gorm.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

// Health performs a basic health check
func (h *HealthHandler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "dvra-api",
		"version": "v1.1.0",
	})
}

// Ready checks if the service is ready to accept traffic
func (h *HealthHandler) Ready(c *gin.Context) {
	sqlDB, err := h.db.DB()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "not_ready", "checks": gin.H{"database": "error"}})
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
	defer cancel()
	if err := sqlDB.PingContext(ctx); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "not_ready", "checks": gin.H{"database": "error"}})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "ready",
		"checks": gin.H{
			"database": "ok",
		},
	})
}
