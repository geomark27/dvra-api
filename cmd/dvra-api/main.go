package main

import (
	"log"

	"dvra-api/internal/platform/config"
	"dvra-api/internal/platform/server"

	"github.com/joho/godotenv"
)

func main() {
	// Cargar variables de entorno desde .env
	_ = godotenv.Load()

	// Cargar configuración
	cfg := config.Load()

	// Crear servidor
	srv := server.New(cfg)

	// Mensaje de inicio
	log.Printf("🚀 Servidor %s iniciado en http://localhost:%s", "dvra-api", cfg.Port)
	log.Printf("✨ Proyecto generado con Loom")
	log.Printf("📖 Documentación disponible en: docs/API.md")

	// Iniciar servidor
	if err := srv.Start(); err != nil {
		log.Fatal("Error iniciando servidor:", err)
	}
}
