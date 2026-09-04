// Package health expone las sondas de vida y disponibilidad del servicio.
//
// Vive en platform/ y no en modules/ porque no es un dominio de negocio: es
// infraestructura que responde a un balanceador o a un orquestador.
package health

import (
	"context"
	"net/http"
	"time"

	"dvra-api/internal/platform/httpx"

	"github.com/gin-gonic/gin"
)

// pingTimeout acota la sonda: si la base no responde en 2 s, el servicio no
// está listo. Sin límite, un readiness colgado deja al orquestador sin señal.
const pingTimeout = 2 * time.Second

// Pinger es lo mínimo que la sonda necesita de la base de datos.
//
// Es un puerto y no *gorm.DB para que el handler no conozca la persistencia:
// antes recibía la conexión cruda, saltándose la capa de datos.
type Pinger interface {
	Ping(ctx context.Context) error
}

// Module expone las sondas.
type Module struct {
	db Pinger
}

func New(db Pinger) *Module { return &Module{db: db} }

// Name identifica al módulo en logs y en la auditoría de rutas.
func (m *Module) Name() string { return "health" }

// RegisterRoutes monta las sondas sin autenticación: quien las consulta es la
// infraestructura, que no tiene sesión.
func (m *Module) RegisterRoutes(r httpx.Router) {
	r.Public.GET("/health", m.live)
	r.Public.GET("/health/ready", m.ready)
}

// live responde si el proceso está en pie. No toca dependencias externas: una
// sonda de vida que falla por la base provocaría que reinicien el servicio
// cuando el problema está en otro sitio.
func (m *Module) live(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "dvra-api",
		"version": "v1.1.0",
	})
}

// ready responde si el servicio puede atender tráfico, lo que exige base de
// datos.
func (m *Module) ready(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), pingTimeout)
	defer cancel()

	if err := m.db.Ping(ctx); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status": "not_ready",
			"checks": gin.H{"database": "error"},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "ready",
		"checks": gin.H{"database": "ok"},
	})
}
