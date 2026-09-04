package health

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"dvra-api/internal/platform/httpx"

	"github.com/gin-gonic/gin"
)

type fakePinger struct{ err error }

func (f fakePinger) Ping(context.Context) error { return f.err }

func routerFor(p Pinger) *gin.Engine {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	New(p).RegisterRoutes(httpx.Router{Public: engine.Group("")})
	return engine
}

func status(t *testing.T, engine *gin.Engine, path string) int {
	t.Helper()
	response := httptest.NewRecorder()
	engine.ServeHTTP(response, httptest.NewRequest(http.MethodGet, path, nil))
	return response.Code
}

// La sonda de vida no depende de la base: si fallara por ella, el orquestador
// reiniciaría el servicio cuando el problema está en otro sitio.
func TestLiveNoDependeDeLaBase(t *testing.T) {
	engine := routerFor(fakePinger{err: errors.New("db caída")})

	if code := status(t, engine, "/health"); code != http.StatusOK {
		t.Fatalf("/health devolvió %d, se esperaba 200", code)
	}
}

func TestReadyReflejaElEstadoDeLaBase(t *testing.T) {
	if code := status(t, routerFor(fakePinger{}), "/health/ready"); code != http.StatusOK {
		t.Fatalf("con base sana devolvió %d, se esperaba 200", code)
	}

	engine := routerFor(fakePinger{err: errors.New("db caída")})
	if code := status(t, engine, "/health/ready"); code != http.StatusServiceUnavailable {
		t.Fatalf("con base caída devolvió %d, se esperaba 503", code)
	}
}
