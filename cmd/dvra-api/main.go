package main

import (
	"log"

	"dvra-api/internal/database"
	"dvra-api/internal/platform/config"
	"dvra-api/internal/platform/server"

	"github.com/joho/godotenv"
)

func main() {
	// Cargar variables de entorno desde .env
	_ = godotenv.Load()

	// Cargar configuración
	cfg := config.Load()

	// Inicializar base de datos
	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatal("Error conectando a la base de datos:", err)
	}
	defer func() {
		if err := database.CloseDB(); err != nil {
			log.Printf("Error cerrando la base de datos: %v", err)
		}
	}()

	// NOTA: Las migraciones se ejecutan manualmente con:
	// go run cmd/console/main.go migrate
	// O con: loom db:migrate
	//
	// Si err := database.AutoMigrate(db); err != nil {
	// 	log.Fatal("Error ejecutando migraciones:", err)
	// }

	// Crear servidor
	srv := server.New(cfg, db)

	// Mensaje de inicio
	log.Printf("🚀 Servidor %s iniciado en http://localhost:%s", "dvra-api", cfg.Port)
	log.Printf("✨ Proyecto generado con Loom")

	// Iniciar servidor
	if err := srv.Start(); err != nil {
		log.Fatal("Error iniciando servidor:", err)
	}
}
