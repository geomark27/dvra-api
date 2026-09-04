# dvra-api - Makefile

# Cargar variables de entorno desde .env
ifneq (,$(wildcard .env))
	include .env
	export
endif

.PHONY: help build run test test-coverage fmt vet fmt-check lint check hooks-install \
	deps clean dev dev-full fresh install-tools up down restart logs logs-api \
	ps rebuild db-up db-down db-restart db-logs db-clean clean-all db-shell \
	db-fresh db-fresh-full db-migrate db-seed db-location push pull status sync swagger

# Variables
APP_NAME=dvra-api
BUILD_DIR=build
CMD_DIR=cmd/$(APP_NAME)
BRANCH := $(shell git branch --show-current)
DOCKER ?= docker
COMPOSE := $(DOCKER) compose
COMPOSE_PROJECT := $(if $(COMPOSE_PROJECT_NAME),$(COMPOSE_PROJECT_NAME),dvra-ats)

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
	@echo "    db-clean           Elimina PostgreSQL; conserva uploads"
	@echo "    clean-all          Elimina DB, uploads y todos los contenedores"
	@echo "    db-shell           Accede a psql en el contenedor"
	@echo ""
	@echo "  Docker Compose (API + DB):"
	@echo "    up                 Levanta toda la aplicacion (API + DB)"
	@echo "    down               Detiene toda la aplicacion"
	@echo "    restart            Reinicia toda la aplicacion"
	@echo "    logs               Muestra logs de todos los servicios"
	@echo "    logs-api           Muestra logs solo de la API"
	@echo "    ps                 Muestra el estado de los servicios"
	@echo "    rebuild            Reconstruye y levanta la API"
	@echo ""
	@echo "  Database:"
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
	@echo "    hooks-install      Instala los git hooks (pre-commit y commit-msg)"
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
	@echo "Nota: Para migraciones usa 'make db-migrate'"
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

hooks-install: ## Instala los git hooks (pre-commit y commit-msg)
	@cp scripts/git-hooks/pre-commit .git/hooks/pre-commit
	@chmod +x .git/hooks/pre-commit
	@cp scripts/git-hooks/commit-msg .git/hooks/commit-msg
	@chmod +x .git/hooks/commit-msg
	@echo "Hooks instalados: .git/hooks/pre-commit, .git/hooks/commit-msg"

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
	@go run ./cmd/console migrate --seed
	@echo "Ready! Starting API..."
	@go run $(CMD_DIR)/main.go

fresh: ## Reset completo (clean DB + migrate + seed) [location=1 para ubicaciones]
	@echo "Fresh install..."
	@$(MAKE) db-clean
	@$(MAKE) db-up
	@go run ./cmd/console migrate --fresh --seed
	@if [ "$(location)" = "1" ]; then \
		echo "Poblando ubicaciones..."; \
		$(MAKE) db-location; \
	fi
	@echo "Database fresh and seeded!"

install-tools: ## Instala herramientas de desarrollo
	@echo "Instalando herramientas de desarrollo..."
	@go install github.com/cosmtrek/air@latest
	@go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2
	@# Fijado a la version de go.mod: con @latest, regenerar el Swagger renombra
	@# todas las definiciones y produce un diff global espurio.
	@# Ojo: v1.16.6 reporta "v1.16.4" en --version (constante mal puesta upstream).
	@go install github.com/swaggo/swag/cmd/swag@v1.16.6

# ============================================
# COMANDOS DOCKER (API + DB) - PRODUCCION
# ============================================

up: ## Levanta toda la aplicacion (API + PostgreSQL)
	@echo "Starting DVRA API..."
	@$(COMPOSE) up -d --wait postgres
	@$(MAKE) db-migrate
	@$(COMPOSE) up -d api
	@echo "API running on http://localhost:$(PORT)"
	@echo "Swagger: http://localhost:$(PORT)/swagger/index.html"

down: ## Detiene toda la aplicacion
	@echo "Stopping DVRA..."
	@$(COMPOSE) down

restart: ## Reinicia toda la aplicacion
	@echo "Restarting DVRA..."
	@$(COMPOSE) restart

logs: ## Muestra logs de todos los servicios
	@$(COMPOSE) logs -f

logs-api: ## Muestra logs solo de la API
	@$(COMPOSE) logs -f api

ps: ## Muestra el estado de los servicios
	@$(COMPOSE) ps --all

rebuild: ## Reconstruye y levanta la API
	@echo "Rebuilding DVRA API..."
	@$(COMPOSE) up -d --wait postgres
	@$(MAKE) db-migrate
	@$(COMPOSE) up -d --build api
	@echo "API rebuilt and running!"

# ============================================
# COMANDOS DOCKER (solo PostgreSQL) - DESARROLLO
# ============================================

db-up: ## Inicia PostgreSQL [fresh=1 para reset] [location=1 para ubicaciones]
	@echo "Starting PostgreSQL (DEV MODE)..."
	@$(COMPOSE) up -d --wait postgres
	@echo "PostgreSQL running on localhost:$(DB_PORT)"
	@if [ "$(fresh)" = "1" ]; then \
		echo ""; \
		echo "Fresh flag detected! Running database reset..."; \
		go run ./cmd/console migrate --fresh --seed; \
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
	@$(COMPOSE) stop postgres

db-restart: ## Reinicia PostgreSQL
	@echo "Restarting PostgreSQL..."
	@$(COMPOSE) restart postgres

db-logs: ## Muestra logs de PostgreSQL
	@$(COMPOSE) logs -f postgres

db-clean: ## Elimina PostgreSQL y su volumen; conserva uploads
	@echo "Cleaning database..."
	@$(COMPOSE) stop api postgres 2>/dev/null || true
	@$(COMPOSE) rm -f api postgres 2>/dev/null || true
	@$(DOCKER) volume rm $(COMPOSE_PROJECT)_postgres_data 2>/dev/null || true
	@echo "Database cleaned; uploads volume preserved"

clean-all: ## Elimina contenedores y todos los volumenes (DB + uploads)
	@echo "WARNING: deleting database and uploaded files..."
	@$(COMPOSE) down -v
	@echo "All Docker data cleaned"

db-shell: ## Accede a psql en el contenedor
	@$(COMPOSE) exec postgres psql -U $(DB_USER) -d $(DB_NAME)

db-fresh: ## Alias: db-up con fresh automatico
	@$(MAKE) db-up fresh=1

db-fresh-full: ## Alias: db-up + fresh + locations
	@$(MAKE) db-up fresh=1 location=1

# ============================================
# COMANDOS BASE DE DATOS
# ============================================

db-migrate: ## Ejecuta migraciones de esquema y datos
	@echo "Running migrations..."
	@go run ./cmd/console migrate

db-seed: ## Ejecuta seeders
	@echo "Running seeders..."
	@go run ./cmd/console seed

db-location: ## Pobla datos de ubicaciones (countries, cities, etc)
	@echo "Poblando datos de ubicaciones..."
	@$(COMPOSE) exec -T postgres psql -U $(DB_USER) -d $(DB_NAME) < scripts/location.sql
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
	@if [ "$(BRANCH)" = "main" ]; then \
		echo "Error: no se commitea directo a main."; \
		echo "   Crea una rama: git switch -c feat/<scope>-<slug>"; \
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
	@if [ "$(BRANCH)" = "main" ]; then \
		echo "Error: no se commitea directo a main."; \
		echo "   Crea una rama: git switch -c feat/<scope>-<slug>"; \
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
	@~/go/bin/swag init -g cmd/dvra-api/main.go -o docs --parseDependency --parseInternal
	@echo "Documentacion generada en docs/"
	@echo "Accede a http://localhost:$(PORT)/swagger/index.html"
