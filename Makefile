# dvra-api - Makefile

# Cargar variables de entorno desde .env
ifneq (,$(wildcard .env))
	include .env
	export
endif

.PHONY: build run test clean fmt fmt-check vet lint check hooks-install deps help

# Variables
APP_NAME=dvra-api
BUILD_DIR=build
CMD_DIR=cmd/$(APP_NAME)
BRANCH := $(shell git branch --show-current)

# Comandos principales
help: ## Muestra esta ayuda
	@echo "================================================================================"
	@echo "                         DVRA-API - Makefile Help"
	@echo "================================================================================"
	@echo ""
	@echo "  Build & Run:"
	@echo "    build              Compila la aplicacion"
	@echo "    run                Ejecuta la aplicacion"
	@echo "    dev                Modo desarrollo con hot reload (requiere air)"
	@echo "    dev-full           Setup completo (DB + migrate + seed + run)"
	@echo ""
	@echo "  PostgreSQL (Desarrollo - solo DB):"
	@echo "    db-up              Inicia solo PostgreSQL"
	@echo "    db-up fresh=1      Inicia PostgreSQL + reset DB + seed"
	@echo "    db-up fresh=1 location=1"
	@echo "                       Inicia PostgreSQL + reset DB + seed + ubicaciones"
	@echo "    db-fresh           Alias: db-up fresh=1"
	@echo "    db-fresh-full      Alias: db-up fresh=1 location=1"
	@echo "    db-down            Detiene PostgreSQL"
	@echo "    db-restart         Reinicia PostgreSQL"
	@echo "    db-logs            Muestra logs de PostgreSQL"
	@echo "    db-clean           Elimina PostgreSQL y volumenes"
	@echo "    db-shell           Accede a psql en el contenedor"
	@echo ""
	@echo "  Docker (Produccion - API + DB):"
	@echo "    up                 Levanta toda la aplicacion (API + DB)"
	@echo "    down               Detiene toda la aplicacion"
	@echo "    restart            Reinicia toda la aplicacion"
	@echo "    logs               Muestra logs de todos los servicios"
	@echo "    logs-api           Muestra logs solo de la API"
	@echo "    rebuild            Reconstruye y levanta la API"
	@echo ""
	@echo "  Database (LOOM):"
	@echo "    db-migrate         Ejecuta migraciones"
	@echo "    db-seed            Ejecuta seeders"
	@echo "    db-location        Pobla datos de ubicaciones (countries, cities, etc)"
	@echo "    fresh              Reset completo (clean + up + migrate + seed)"
	@echo "    fresh location=1   Reset completo + datos de ubicaciones"
	@echo ""
	@echo "  Testing & Calidad:"
	@echo "    test               Ejecuta los tests"
	@echo "    test-coverage      Ejecuta tests con cobertura"
	@echo "    fmt                Formatea el codigo"
	@echo "    fmt-check          Verifica formato (falla si hay archivos sin gofmt)"
	@echo "    vet                Ejecuta go vet"
	@echo "    lint               Ejecuta golangci-lint"
	@echo "    check              Gate de calidad: fmt-check + vet + lint + test + build"
	@echo ""
	@echo "  Git (rama: $(BRANCH)):"
	@echo "    push m='msg'       Add + Commit + Push a $(BRANCH)"
	@echo "    pull               Pull desde origin/$(BRANCH)"
	@echo "    status             Ver estado de git"
	@echo "    sync m='msg'       Pull + Add + Commit + Push (sincronizar)"
	@echo ""
	@echo "  Utilidades:"
	@echo "    clean              Limpia archivos generados"
	@echo "    deps               Descarga las dependencias"
	@echo "    install-tools      Instala herramientas de desarrollo (air, golangci-lint, swag)"
	@echo "    hooks-install      Instala el git hook de pre-commit"
	@echo "    swagger            Genera documentacion Swagger"
	@echo ""
	@echo "================================================================================"
	@echo "                              Ejemplos de Uso"
	@echo "================================================================================"
	@echo ""
	@echo "  Desarrollo rapido:"
	@echo "    make db-up                        Solo inicia PostgreSQL"
	@echo "    make db-up fresh=1                PostgreSQL + reset DB + seed"
	@echo "    make db-up fresh=1 location=1    PostgreSQL + reset + ubicaciones"
	@echo "    make db-fresh                     Atajo para db-up fresh=1"
	@echo "    make db-fresh-full                Atajo para db-up fresh=1 location=1"
	@echo ""
	@echo "  Workflow tipico:"
	@echo "    make db-fresh-full && make run    Reset completo + ejecutar API"
	@echo "    make dev-full                     Todo automatico (DB + API)"
	@echo ""
	@echo "  Git rapido:"
	@echo "    make push m='fix: corregido bug en login'"
	@echo "    make sync m='feat: nueva funcionalidad'"
	@echo ""
	@echo "================================================================================"

build: ## Compila la aplicacion
	@echo "Compilando $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(APP_NAME) $(CMD_DIR)/main.go
	@echo "Compilacion exitosa: $(BUILD_DIR)/$(APP_NAME)"

run: ## Ejecuta la aplicacion (sin migraciones)
	@echo "Ejecutando $(APP_NAME)..."
	@echo "Nota: Para migraciones usa 'make db-migrate' o 'loom db:migrate'"
	@go run $(CMD_DIR)/main.go

test: ## Ejecuta los tests
	@echo "Ejecutando tests..."
	@go test -v ./...

test-coverage: ## Ejecuta tests con cobertura
	@echo "Ejecutando tests con cobertura..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Reporte de cobertura generado: coverage.html"

fmt: ## Formatea el codigo
	@echo "Formateando codigo..."
	@go fmt ./...

vet: ## Ejecuta go vet
	@echo "Analizando codigo..."
	@go vet ./...

fmt-check: ## Verifica formato sin modificar (falla si hay archivos sin gofmt)
	@echo "Verificando formato..."
	@unformatted=$$(gofmt -l .); \
	if [ -n "$$unformatted" ]; then \
		echo "Archivos sin gofmt:"; echo "$$unformatted"; \
		echo "Corrige con: gofmt -w ."; exit 1; \
	fi; \
	echo "Formato OK"

lint: ## Ejecuta golangci-lint (requiere instalacion)
	@echo "Ejecutando linter..."
	@golangci-lint run

check: ## Gate de calidad: fmt-check + vet + lint + test + build
	@$(MAKE) fmt-check
	@$(MAKE) vet
	@$(MAKE) lint
	@$(MAKE) test
	@$(MAKE) build
	@echo "Todos los checks pasaron"

hooks-install: ## Instala el git hook de pre-commit
	@cp scripts/git-hooks/pre-commit .git/hooks/pre-commit
	@chmod +x .git/hooks/pre-commit
	@echo "Hook pre-commit instalado en .git/hooks/pre-commit"

deps: ## Descarga las dependencias
	@echo "Descargando dependencias..."
	@go mod download
	@go mod tidy

clean: ## Limpia archivos generados
	@echo "Limpiando archivos generados..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@go clean

dev: ## Modo desarrollo con hot reload (requiere air)
	@echo "Iniciando en modo desarrollo..."
	@air

dev-full: ## Setup completo desarrollo (DB + migrate + seed + run)
	@echo "Starting development environment..."
	@$(MAKE) db-up
	@echo "Waiting for PostgreSQL..."
	@sleep 3
	@loom db:migrate --seed
	@echo "Ready! Starting API..."
	@go run $(CMD_DIR)/main.go

fresh: ## Reset completo (clean DB + migrate + seed) [location=1 para ubicaciones]
	@echo "Fresh install..."
	@$(MAKE) db-clean
	@$(MAKE) db-up
	@echo "Waiting for PostgreSQL..."
	@sleep 3
	@loom db:fresh --seed
	@if [ "$(location)" = "1" ]; then \
		echo "Poblando ubicaciones..."; \
		$(MAKE) db-location; \
	fi
	@echo "Database fresh and seeded!"

install-tools: ## Instala herramientas de desarrollo
	@echo "Instalando herramientas de desarrollo..."
	@go install github.com/cosmtrek/air@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/swaggo/swag/cmd/swag@latest

# ============================================
# COMANDOS DOCKER (API + DB) - PRODUCCION
# ============================================

up: ## Levanta toda la aplicacion (API + PostgreSQL)
	@echo "Starting DVRA API..."
	@docker compose up -d
	@echo "API running on http://localhost:$(PORT)"
	@echo "Swagger: http://localhost:$(PORT)/swagger/index.html"

down: ## Detiene toda la aplicacion
	@echo "Stopping DVRA..."
	@docker compose down

restart: ## Reinicia toda la aplicacion
	@echo "Restarting DVRA..."
	@docker compose restart

logs: ## Muestra logs de todos los servicios
	@docker compose logs -f

logs-api: ## Muestra logs solo de la API
	@docker compose logs -f api

rebuild: ## Reconstruye y levanta la API
	@echo "Rebuilding DVRA API..."
	@docker compose build --no-cache api
	@docker compose up -d api
	@echo "API rebuilt and running!"

# ============================================
# COMANDOS DOCKER (solo PostgreSQL) - DESARROLLO
# ============================================

db-up: ## Inicia PostgreSQL [fresh=1 para reset] [location=1 para ubicaciones]
	@echo "Starting PostgreSQL (DEV MODE)..."
	@docker compose up -d postgres
	@echo "PostgreSQL running on localhost:$(DB_PORT)"
	@if [ "$(fresh)" = "1" ]; then \
		echo ""; \
		echo "Fresh flag detected! Running database reset..."; \
		echo "Waiting for PostgreSQL to be ready..."; \
		sleep 3; \
		loom db:fresh --seed; \
		if [ "$(location)" = "1" ]; then \
			echo "Poblando ubicaciones..."; \
			$(MAKE) db-location; \
		fi; \
		echo "Database fresh and seeded!"; \
	else \
		echo ""; \
		echo "TIP: Ejecuta tu API con 'make run' o 'make dev'"; \
		echo "TIP: Para reset completo usa: make db-up fresh=1"; \
		echo "TIP: Para reset + ubicaciones: make db-up fresh=1 location=1"; \
	fi

db-down: ## Detiene PostgreSQL
	@echo "Stopping PostgreSQL..."
	@docker compose stop postgres

db-restart: ## Reinicia PostgreSQL
	@echo "Restarting PostgreSQL..."
	@docker compose restart postgres

db-logs: ## Muestra logs de PostgreSQL
	@docker compose logs -f postgres

db-clean: ## Elimina PostgreSQL y volumenes
	@echo "Cleaning database..."
	@docker compose down -v
	@echo "Database cleaned"

db-shell: ## Accede a psql en el contenedor
	@docker compose exec postgres psql -U $(DB_USER) -d $(DB_NAME)

db-fresh: ## Alias: db-up con fresh automatico
	@$(MAKE) db-up fresh=1

db-fresh-full: ## Alias: db-up + fresh + locations
	@$(MAKE) db-up fresh=1 location=1

# ============================================
# COMANDOS BASE DE DATOS (LOOM)
# ============================================

db-migrate: ## Ejecuta migraciones con LOOM
	@echo "Running migrations..."
	@loom db:migrate

db-seed: ## Ejecuta seeders con LOOM
	@echo "Running seeders..."
	@loom db:seed

db-location: ## Pobla datos de ubicaciones (countries, cities, etc)
	@echo "Poblando datos de ubicaciones..."
	@docker compose exec -T postgres psql -U $(DB_USER) -d $(DB_NAME) < scripts/location.sql
	@echo "Datos de ubicaciones poblados exitosamente"

# ============================================
# COMANDOS GIT
# ============================================

push: ## Push rapido: make push m='mensaje'
	@if [ -z "$(m)" ]; then \
		echo "Error: Debes proporcionar un mensaje"; \
		echo "   Uso: make push m='tu mensaje de commit'"; \
		exit 1; \
	fi
	@echo "Agregando archivos..."
	@git add .
	@echo "Commiteando: $(m)"
	@git commit -m "$(m)"
	@echo "Pusheando a origin/$(BRANCH)..."
	@git push origin $(BRANCH)
	@echo "Push completado exitosamente!"

pull: ## Pull desde origin
	@echo "Pulling desde origin/$(BRANCH)..."
	@git fetch origin
	@git pull origin $(BRANCH)
	@echo "Pull completado!"

status: ## Ver estado de git
	@echo "Estado de Git (rama: $(BRANCH)):"
	@echo ""
	@git status

sync: ## Sincronizar (pull + push): make sync m='mensaje'
	@if [ -z "$(m)" ]; then \
		echo "Error: Debes proporcionar un mensaje"; \
		echo "   Uso: make sync m='tu mensaje de commit'"; \
		exit 1; \
	fi
	@echo "Pulling cambios..."
	@git pull origin $(BRANCH)
	@echo "Agregando archivos..."
	@git add .
	@echo "Commiteando: $(m)"
	@git commit -m "$(m)"
	@echo "Pusheando a origin/$(BRANCH)..."
	@git push origin $(BRANCH)
	@echo "Sincronizacion completada!"

# ============================================
# DOCUMENTACION
# ============================================

swagger: ## Genera documentacion Swagger/OpenAPI
	@echo "Generando documentacion Swagger..."
	@~/go/bin/swag init -g cmd/dvra-api/main.go -o docs
	@echo "Documentacion generada en docs/"
	@echo "Accede a http://localhost:8000/swagger/index.html"
