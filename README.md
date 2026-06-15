# Dvra API

Backend del **Dvra ATS** — plataforma SaaS multi-tenant de reclutamiento (Applicant Tracking System) para startups tech en LATAM, con marketplace de talento (Red Dvra) en el roadmap.

API REST construida en **Go 1.24 + Gin + GORM + PostgreSQL 16**, con autenticación JWT, multi-empresa (switch de contexto), pipeline de candidatos, dashboard analítico y career page pública.

## 📚 Documentación

Toda la documentación está consolidada en 4 documentos:

| Documento | Contenido |
|---|---|
| [docs/01_LOGICA_DE_NEGOCIO.md](docs/01_LOGICA_DE_NEGOCIO.md) | Visión del producto, modelo de negocio dual (ATS + Marketplace), roles y permisos, reglas de negocio, pipeline, pricing y límites por plan, compliance |
| [docs/02_PLAN_DE_NEGOCIO.md](docs/02_PLAN_DE_NEGOCIO.md) | Plan operativo y financiero del Año 1: objetivos, timeline mensual, go-to-market, modelo financiero, KPIs, riesgos y visión multi-año |
| [docs/03_FLUJO_DE_USO.md](docs/03_FLUJO_DE_USO.md) | Flujo completo de uso de la aplicación por actor: registro/onboarding, multi-empresa, equipo, jobs, candidatos, pipeline, career page, SuperAdmin |
| [docs/04_DOCUMENTACION_TECNICA_API.md](docs/04_DOCUMENTACION_TECNICA_API.md) | Detalle técnico del backend: stack, arquitectura, modelos de datos, endpoints, JWT/middleware, multi-tenancy, seeders, Swagger, despliegue y deuda técnica |

## 🏃 Inicio rápido

```bash
cp .env.example .env
go mod tidy

# Base de datos + seeders
make db-fresh          # drop + migrate + seed (~7 s)

# Ejecutar
make run               # API en http://localhost:8080

# O con Docker
docker compose up -d
```

- **Swagger UI:** http://localhost:8080/swagger/index.html (regenerar con `make swagger`)
- **Health check:** `GET /api/v1/health`
- **Comandos disponibles:** `make help`

## 🏗️ Arquitectura (resumen)

```
cmd/dvra-api        → entry point del servidor HTTP
cmd/console         → CLI (migrate / seed)
internal/app        → handlers → services → repositories → models (+ dtos)
internal/platform   → config y server (DI + rutas)
internal/shared     → middleware (auth JWT, roles)
internal/database   → init + seeders
```

Detalle completo en [docs/04_DOCUMENTACION_TECNICA_API.md](docs/04_DOCUMENTACION_TECNICA_API.md).

---

Proyecto generado inicialmente con [Loom](https://github.com/geomark27/loom-go).
