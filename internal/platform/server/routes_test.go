package server

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"dvra-api/internal/platform/config"

	"github.com/gin-gonic/gin"
)

// registerTestRoutes monta el árbol completo de rutas sin base de datos.
//
// Funciona porque los repositorios reciben *gorm.DB inyectado y no consultan
// nada al construirse: basta con pasar nil. Si algún día un constructor toca la
// base, este test lo detecta de inmediato.
func registerTestRoutes(t *testing.T, router *gin.Engine) {
	t.Helper()

	cfg := &config.Config{
		Environment:      "test",
		JWTSecret:        "test-access-secret-with-more-than-32-chars",
		JWTRefreshSecret: "test-refresh-secret-with-more-than-32-chars",
		StorageRoot:      t.TempDir(),
	}

	d, err := buildDeps(cfg, nil)
	if err != nil {
		t.Fatalf("buildDeps: %v", err)
	}

	registerRoutes(router, d, cfg)
}

// TestRouteContract congela el conjunto exacto de rutas que expone la API.
//
// Es la red de seguridad de la migración a módulos: mover un handler de paquete
// NUNCA debe cambiar esta lista. Si una fase de migración hace fallar este test,
// es un bug de esa fase, no un cambio legítimo.
//
// Para actualizar el golden a propósito (una ruta nueva de verdad):
//
//	UPDATE_GOLDEN=1 go test ./internal/platform/server/ -run TestRouteContract
//
// y revisa el diff del .golden en el PR: ahí se ve el cambio de contrato.
func TestRouteContract(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	registerTestRoutes(t, router)

	got := make([]string, 0, len(router.Routes()))
	for _, route := range router.Routes() {
		got = append(got, route.Method+" "+route.Path)
	}
	sort.Strings(got)

	goldenPath := filepath.Join("testdata", "routes.golden")

	if os.Getenv("UPDATE_GOLDEN") == "1" {
		if err := os.MkdirAll("testdata", 0o755); err != nil {
			t.Fatal(err)
		}
		content := strings.Join(got, "\n") + "\n"
		if err := os.WriteFile(goldenPath, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
		t.Logf("golden actualizado: %d rutas", len(got))
		return
	}

	raw, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Fatalf("no se pudo leer %s: %v\nGenéralo con: UPDATE_GOLDEN=1 go test ./internal/platform/server/ -run TestRouteContract", goldenPath, err)
	}

	want := strings.Split(strings.TrimSpace(string(raw)), "\n")
	sort.Strings(want)

	if diff := diffRoutes(want, got); diff != "" {
		t.Fatalf("el contrato de rutas cambió:\n%s", diff)
	}
}

// diffRoutes describe qué rutas sobran y cuáles faltan.
func diffRoutes(want, got []string) string {
	inWant := make(map[string]bool, len(want))
	for _, route := range want {
		inWant[route] = true
	}
	inGot := make(map[string]bool, len(got))
	for _, route := range got {
		inGot[route] = true
	}

	var b strings.Builder
	for _, route := range got {
		if !inWant[route] {
			b.WriteString("  + " + route + "  (ruta nueva, no está en el golden)\n")
		}
	}
	for _, route := range want {
		if !inGot[route] {
			b.WriteString("  - " + route + "  (ruta perdida)\n")
		}
	}
	return b.String()
}
