# Dvra API — Documentación Técnica del Backend

> **Especificación técnica del API REST (solo backend)** — verificada contra el código en `v1.4.0` (commit `e04aafd`, "separación del superadmin")
> Consolidado: Junio 2026 | Reemplaza a: `ARCHITECTURE.md`, `API_ENDPOINTS.md`, `PLANS_MODULE.md` (técnico), `LOCATION_MODULE.md`, `LOCATION_VERIFICATION.md`, `MULTI_COMPANY_IMPLEMENTATION.md`, `SUPERADMIN_IMPLEMENTATION.md`, `SECURITY_AUDIT.md`, `SWAGGER_GUIDE.md`
>
> 📌 Documentos relacionados: [01_LOGICA_DE_NEGOCIO.md](./01_LOGICA_DE_NEGOCIO.md) · [02_PLAN_DE_NEGOCIO.md](./02_PLAN_DE_NEGOCIO.md) · [03_FLUJO_DE_USO.md](./03_FLUJO_DE_USO.md)

---

## 1. Stack Técnico

| Componente | Tecnología | Versión |
|---|---|---|
| Lenguaje | Go | 1.24.0 |
| Framework HTTP | Gin | 1.11.0 |
| ORM | GORM (driver PostgreSQL) | 1.31.1 |
| Base de datos | PostgreSQL | 16-alpine |
| Autenticación | golang-jwt/jwt | v5.3.0 |
| Passwords | golang.org/x/crypto (bcrypt) | 0.46.0 |
| CLI de consola | spf13/cobra | 1.9.1 |
| Env vars | joho/godotenv | 1.5.1 |
| Docs API | swaggo/gin-swagger | 1.6.1 |
| Scaffolding/helpers | geomark27/loom-go | 1.1.3 |

**Puertos:** API `8080` (configurable vía `PORT`); PostgreSQL `5433` en dev local / `5432` en Docker.

**Variables de entorno** (`.env.example`): `PORT`, `ENVIRONMENT`, `LOG_LEVEL`, `CORS_ALLOWED_ORIGINS`, `DB_HOST/PORT/USER/PASSWORD/NAME`, `JWT_SECRET`, `JWT_REFRESH_SECRET`.

---

## 2. Arquitectura

### 2.1 Capas

```
HTTP (Gin Router + Middleware: CORS, Auth JWT)
  └── Handlers   (internal/app/handlers)    — parseo HTTP, extracción de claims, validación de permisos
        └── Services (internal/app/services)     — lógica de negocio, validaciones, transacciones
              └── Repositories (internal/app/repositories) — acceso a datos con GORM
                    └── Models (internal/app/models)            — entidades de dominio (tags GORM/JSON)
                          └── PostgreSQL 16 (GORM, AutoMigrate, soft deletes)
```

### 2.2 Estructura del proyecto

```
dvra-api/
├── cmd/
│   ├── dvra-api/main.go        # Entry point del servidor HTTP (+ anotaciones Swagger raíz)
│   └── console/main.go         # CLI Cobra: migrate / seed
├── internal/
│   ├── app/
│   │   ├── handlers/           # auth, user, company, membership, job, candidate,
│   │   │                       # application, plan, system_value, location,
│   │   │                       # dashboard, public, platform_settings, health
│   │   ├── services/           # lógica de negocio + jwt_service
│   │   ├── repositories/       # acceso a datos (incluye dashboard_repository)
│   │   ├── dtos/               # request/response objects con binding tags
│   │   └── models/             # entidades + constantes de roles/estados
│   ├── database/
│   │   ├── (init / models_all) # InitDB, AutoMigrate, pool de conexiones
│   │   └── seeders/            # role, plan, system_value, platform_settings, user, company
│   ├── platform/
│   │   ├── config/config.go    # Load() desde env, helpers IsDevelopment/IsProduction
│   │   └── server/             # server.go (DI manual + CORS) y routes.go (registro de rutas)
│   └── shared/middleware/      # auth_middleware.go (AuthMiddleware, RequireRole, RequireCompany, OptionalAuth)
├── docs/                       # esta documentación + swagger generado (docs.go/swagger.json/yaml)
├── scripts/                    # SQL auxiliar (carga masiva de ubicaciones)
├── Dockerfile                  # multi-stage (builder Go → alpine, usuario no-root)
├── docker-compose.yml          # api + postgres 16 (healthcheck, volumen persistente)
└── Makefile                    # run, build, swagger, db-migrate, db-seed, db-fresh, db-location...
```

### 2.3 Inicialización (cmd/dvra-api/main.go)

1. `godotenv.Load()` → carga `.env`.
2. `config.Load()` → struct `Config`.
3. `database.InitDB(cfg)` → conexión GORM + AutoMigrate de `AllModels` + pool (10 idle / 100 max open / 1h lifetime). Expone singleton `database.DB`.
4. `server.New(cfg, db)` → **inyección de dependencias manual**: instancia repositorios → services → handlers y los pasa a `registerRoutes()`.
5. `srv.Start()` → escucha en `:PORT`.

### 2.4 CORS

Middleware propio en `server.go`: orígenes desde `CORS_ALLOWED_ORIGINS`; métodos GET/POST/PUT/DELETE/PATCH/OPTIONS; headers `Authorization`, `Content-Type`, `X-Company-ID`; maneja preflight.

---

## 3. Modelos de Datos

Todos usan `gorm.Model` (ID, CreatedAt, UpdatedAt, DeletedAt → **soft delete global**).

### 3.1 Identidad y multi-tenancy

**`users`** — `Email` (unique, not null), `PasswordHash` (bcrypt), `FirstName`, `LastName`, `AvatarURL`, `EmailVerified` (default false), `LastLoginAt`, `IsActive` (default true). Relación: `Memberships` 1:N.

**`companies`** (tenant) — `Name`, `Slug` (unique), `LogoURL`, `PlanTier` (default `'free'`, referencia el slug del plan), `TrialEndsAt`, `Timezone` (default `'America/Bogota'`). Método `IsTrialActive()`. Relación: `Memberships` 1:N.

**`memberships`** — ⭐ pieza central del multi-tenancy:

| Campo | Detalle |
|---|---|
| `UserID` | FK not null (índice compuesto con CompanyID) |
| `CompanyID` | FK **nullable** → `NULL` = SuperAdmin |
| `Role` | `admin` / `recruiter` / `hiring_manager` / `user` (/ `superadmin`) |
| `Status` | default `active` (pending/suspended/removed previstos) |
| `IsDefault` | empresa usada al hacer login |
| `InvitedBy`, `InvitedAt`, `JoinedAt` | tracking de invitación |

**`roles`** — catálogo: `Name`, `Slug` (unique), `Level` (admin=50, recruiter=30, user=10), `IsSystem` (no eliminables). Constantes en `models/role.go`.

### 3.2 Core ATS

**`jobs`** — `CompanyID` (FK, índice compuesto con Status), `Title` (not null), `Description`, `SalaryMin/Max`, `Requirements`, `Benefits`, `Status` (`draft`/`published`/`closed`), `LocationType` (`remote`/`onsite`/`hybrid`), `CityID` (FK nullable → catálogo de ubicaciones), `AssignedRecruiter`, `HiringManager` (FK nullable → users). Relaciones: Company, City, Applications 1:N.

**`candidates`** — `CompanyID` + `Email` con **unique compuesto** (mismo email puede existir en otra empresa), `FirstName`, `LastName`, `Phone`, `ResumeURL`, `GithubURL`, `LinkedinURL`, `Source` (linkedin/referral/direct_apply/agency/...). Relaciones: Company, Applications 1:N.

**`applications`** — `JobID`, `CandidateID`, `CompanyID` (redundante a propósito: índice compuesto con Stage para queries de pipeline), `Stage` (`applied`/`screening`/`technical`/`interview`/`offer`/`hired`/`rejected`), `Rating` (1–5 nullable), `Notes`, `AppliedAt`, `RejectedAt`, `HiredAt`.

### 3.3 Monetización y plataforma

**`plans`** — `Name`, `Slug` (unique), `Description`, `Price` decimal, `Currency` (default USD), `BillingCycle` (`monthly`/`yearly`), `IsActive`, `IsPublic`, `TrialDays`, `DisplayOrder`; límites `MaxUsers/MaxJobs/MaxCandidates/MaxApplications/MaxStorageGB` (**-1 = ilimitado**); features `CanExportData/CanUseCustomBrand/CanUseAPI/CanUseIntegrations`; `SupportLevel` (`email`/`priority`/`dedicated`). Métodos: `IsUnlimited(limitType)`, `HasFeature(feature)`.

**`platform_settings`** — singleton (1 fila): branding (`PlatformName`, `Tagline`, logos, `PrimaryColor`), contacto (`SupportEmail`, `SalesEmail`), URLs (marketing/docs/terms/privacy), defaults de negocio (`DefaultTrialDays=14`, `DefaultPlanTier`), datos legales y redes sociales.

**`system_values`** — catálogos dinámicos: `Category` + `Value` (índice compuesto), `Label`, `Description`, `DisplayOrder`, `IsActive`, `CompanyID` **nullable** (NULL = global; con valor = override por empresa). Categorías sembradas: `job_status`, `application_status`, `contract_type`, `work_mode`, `experience_level`, `priority`, `candidate_source`.

### 3.4 Ubicaciones (jerarquía geográfica)

```
regions ──1:N── subregions ──1:N── countries ──1:N── states ──1:N── cities
```

| Tabla | Campos clave |
|---|---|
| `regions` | Name, IsActive |
| `subregions` | Name, RegionID (FK) |
| `countries` | Name, **Iso2/Iso3 (unique)**, NumericCode, PhoneCode, Timezones, SubregionID (FK *nullable*) |
| `states` | Name, CountryID (FK), CountryCode |
| `cities` | Name, StateID (FK), lat/lng |

Relaciones GORM bidireccionales (`belongs-to` + `has-many`) → permiten preload en cascada en ambas direcciones, p. ej. `db.Preload("State.Country.Subregion.Region")`. El orden en `AllModels` respeta las dependencias para el AutoMigrate.

---

## 4. Autenticación y Autorización

### 4.1 JWT

```go
type JWTClaims struct {
    UserID    uint
    CompanyID *uint   // nil = sin contexto de empresa (SuperAdmin)
    Email     string
    Role      string
    jwt.RegisteredClaims
}
```

- **Access token:** 1 hora · **Refresh token:** 30 días.
- Firmas HS256 con `JWT_SECRET` / `JWT_REFRESH_SECRET`.
- `JWTService`: `GenerateAccessToken(userID, companyID, email, role)`, `GenerateRefreshToken(userID)`, `ValidateToken(token)`.
- El token **lleva el contexto completo**: cambiar de empresa = emitir token nuevo (`/auth/switch-company`), nunca mutar el actual.

### 4.2 Middleware (`internal/shared/middleware/auth_middleware.go`)

| Middleware | Función |
|---|---|
| `AuthMiddleware(jwtService)` | Valida `Authorization: Bearer <token>`; inyecta en el contexto Gin: `user_id`, `email`, `role`, `company_id` (si existe). 401 si inválido/expirado |
| `RequireRole(minLevel)` | Jerarquía: admin=50, recruiter=30, hiring_manager=20, user=10. 403 si insuficiente |
| `RequireCompany()` | Exige `company_id` en contexto. 403 si falta |
| `OptionalAuth(jwtService)` | Valida token si está presente; continúa sin él (rutas públicas con contexto opcional) |

### 4.3 Roles

Jerarquía completa (`getRoleLevel`): `superadmin`=100 > `admin`=50 > `recruiter`=30 > `hiring_manager`=20 > `user`=10.

> ⚠️ **Separación del SuperAdmin (commit `e04aafd`):** este API **ya no expone** `/auth/superadmin/login` ni el grupo `/api/v1/admin/*` (handlers `admin/superadmin_companies_handler.go`, service `admin/superadmin_companies_service.go` y `superadmin_dto.go` fueron **eliminados**; el panel vive en un servicio aparte). Lo que **permanece** en este API:
> - `Membership.CompanyID` nullable y `JWTClaims.CompanyID *uint` (soporte de tokens sin empresa).
> - Checks `if role == "superadmin"` en los handlers: un token superadmin válido obtiene lectura global (sin filtro por empresa) y es el único que puede `POST /memberships`.
> - El seeder aún crea el usuario `superadmin@dvra.com` / `SuperAdmin123!` (⚠️ cambiar en producción).

---

## 5. Referencia de Endpoints

**Base URL:** `/api/v1` · **Swagger UI:** `/swagger/index.html` · **Root `GET /`:** metadata de la API (versión, índice de endpoints).

### 5.1 Salud (público)

| Método | Ruta | Descripción |
|---|---|---|
| GET | `/health` | Liveness |
| GET | `/health/ready` | Readiness (incluye estado de la BD) |

### 5.2 Autenticación

| Método | Ruta | Acceso | Descripción |
|---|---|---|---|
| POST | `/auth/register-company` | Público | **Flujo principal**: crea Company + User admin + Membership en transacción; devuelve tokens con `company_id` |
| POST | `/auth/register` | Público | ⚠️ DEPRECATED — usuario sin empresa |
| POST | `/auth/login` | Público | Valida credenciales; token con la empresa default; devuelve lista de empresas del usuario |
| POST | `/auth/refresh` | Público | Nuevo access token a partir del refresh token |
| GET | `/auth/me` | JWT | Usuario autenticado |
| POST | `/auth/change-password` | JWT | Valida password anterior, re-hashea |
| POST | `/auth/logout` | JWT | Stateless: el cliente descarta los tokens |
| POST | `/auth/switch-company` | JWT | Valida membresía activa; emite token con nuevo `company_id` + rol de esa empresa |
| GET | `/auth/my-companies` | JWT | Empresas del usuario (selector multi-empresa) |

### 5.3 Recursos protegidos (JWT + scope de empresa)

| Recurso | Endpoints |
|---|---|
| **Users** | `GET /users` · `POST /users` (crea User + Membership en la empresa del token) · `GET/PUT/DELETE /users/:id` |
| **Companies** | `GET /companies` (cliente: solo la suya) · `POST /companies` · `GET/PUT/DELETE /companies/:id` |
| **Memberships** | `GET /memberships` · `POST /memberships` (**403 salvo superadmin**) · `GET/PUT/DELETE /memberships/:id` |
| **Jobs** | `GET /jobs` · `POST /jobs` (nace `draft`) · `GET/PUT/DELETE /jobs/:id` · `PATCH /jobs/:id/publish` · `PATCH /jobs/:id/close` |
| **Candidates** | `GET /candidates` · `POST /candidates` (email único por empresa) · `GET/PUT/DELETE /candidates/:id` · `POST /candidates/:id/upload-resume` (multipart) |
| **Applications** | `GET /applications` · `GET /applications/by-stage` (agrupado para Kanban) · `POST /applications` · `GET/PUT/DELETE /applications/:id` · `PATCH /applications/:id/move` (cambia stage + timestamps automáticos) · `PATCH /applications/:id/rate` (1–5) |
| **Dashboard** | `GET /dashboard/stats` (estadísticas completas de la empresa, ver §7.7) |

### 5.4 Rutas públicas (sin autenticación)

| Grupo | Endpoints |
|---|---|
| **Plans** | `GET /plans` (activos + públicos, pricing page) · `GET /plans/:slug` |
| **System Values** | `GET /system-values/:category` (header opcional `X-Company-ID` para incluir overrides de empresa) |
| **Career page** | `GET /public/platform-settings` · `GET /public/companies/:slug` · `GET /public/companies/:slug/jobs` (solo `published`) · `GET /public/jobs/:id` · `POST /public/jobs/:id/apply` (crea Candidate + Application) |
| **Locations** (read-only) | `GET /locations/regions[/:id]` · `/subregions[/:id]` · `/countries[/:id]` · `/countries/iso/:iso` · `/states[/:id]` · `/cities[/:id]` · `/hierarchy/:id` · `/search?q=` — filtros: `region_id`, `subregion_id`, `country_id`, `state_id`, `search`, `include_*=true` para preload |

### 5.5 Códigos de estado y convenciones

| Código | Uso |
|---|---|
| 200/201 | OK / creado |
| 400 | Body inválido (binding tags) |
| 401 | Token inválido/expirado, credenciales incorrectas |
| 403 | Sin permisos o recurso de otra empresa |
| 404 | No encontrado |
| 409 | Conflicto (email/slug duplicado, plan en uso) |
| 500 | Error interno |

Paginación (donde aplica): `page` (default 1), `limit` (default 20, max 100).

---

## 6. Multi-Tenancy: Implementación

**Modelo:** base de datos compartida con aislamiento por fila (`company_id` scoping). Índices compuestos en `jobs(company_id, status)`, `candidates(company_id, email)`, `applications(company_id, stage)`.

### 6.1 Las 4 garantías (auditadas)

**1. Listados filtrados por el token:**
```go
role, _ := c.Get("role")
if role == "superadmin" {
    return h.service.GetAll()                       // lectura global
}
companyID, _ := c.Get("company_id")
return h.service.GetByCompanyID(companyID.(uint))    // WHERE company_id = ?
```

**2. CREATE fuerza el company_id del token** (anti-manipulación):
```go
companyID, _ := c.Get("company_id")
dto.CompanyID = companyID.(uint)   // ignora cualquier company_id del body
```

**3. GET/PUT/DELETE individual valida pertenencia:**
```go
if resource.CompanyID != tokenCompanyID { return 403 }
```

**4. Switch de contexto = token nuevo** — imposible mezclar empresas en una sesión.

### 6.2 Cobertura por handler (resultado de la auditoría de seguridad)

| Handler | Listado filtrado | Create fuerza company | Update/Delete validan pertenencia |
|---|---|---|---|
| Jobs | ✅ | ✅ | ✅ |
| Candidates | ✅ | ✅ | ✅ |
| Applications | ✅ | ✅ | ✅ |
| Memberships | ✅ | ✅ (POST solo superadmin) | ✅ |
| Companies | ✅ (cliente solo la suya; GET /:id valida id == token) | POST abierta a autenticados | DELETE restringido |
| Users | ✅ (JOIN memberships) | ✅ | ⚠️ PUT/DELETE `/users/:id` pendiente de validar membership |

---

## 7. Servicios y Lógica de Negocio

### 7.1 AuthService
- **`RegisterCompany`** — transacción: valida email único → valida plan `free` activo (`FindActiveBySlug`, sin strings a ciegas) → crea Company (trial +1 mes) → User (bcrypt) → Membership (admin, default) → tokens. Rollback total ante error.
- **`Login` / `LoginWithCompanies`** — bcrypt compare → membresía default → token con `company_id` → lista de empresas → actualiza `LastLoginAt`.
- **`SwitchCompany`** — valida membresía activa en la empresa destino → token nuevo con el rol de esa membresía. Errores tipados: `ErrCompanyNotFound`, `ErrNoMembership`.
- **`RefreshToken`**, **`ChangePassword`** (valida la anterior), **`GetMe`**, **`GetUserCompanies`** (memberships activas con Preload de Company).

### 7.2 JobService
- `CreateJob` (status default `draft`), `GetJobsByCompanyID`, `GetJobsByStatus`, `UpdateJob`, `DeleteJob` (soft).
- **`PublishJob`** — exige `title` y `description` antes de pasar a publicado.
- **`CloseJob`** — transición a `closed`.

### 7.3 CandidateService
- CRUD scoped por empresa; unicidad `(company_id, email)`.
- **`UploadResume`** — guarda el archivo bajo `uploads/companies/{slug}/` (directorios creados por CompanyService al crear la empresa).

### 7.4 ApplicationService
- `CreateApplication` (stage inicial `applied`, `applied_at` auto).
- **`MoveToStage`** — setea automáticamente `RejectedAt` al pasar a `rejected` y `HiredAt` al pasar a `hired`.
- **`RateApplication`** (1–5), `GetApplicationsGroupedByStage` (Kanban), filtros por job/empresa/stage.

### 7.5 PlanService
- CRUD con validaciones: **slug único** (409), precio ≥ 0, currency de 3 letras, `billing_cycle ∈ {monthly, yearly}`, `support_level ∈ {email, priority, dedicated}`.
- `GetPublicPlans` (pricing page: `is_public AND is_active`), `GetActivePlans`, `GetAllPlans`.
- `TogglePlanStatus` — desactivar oculta el plan de `/plans` pero **no afecta** a empresas que ya lo tienen.
- `DeletePlan` — bloqueado (409) si alguna empresa usa el plan; soft delete.
- `AssignPlanToCompany` — valida plan activo → `company.plan_tier = plan.slug`.

### 7.6 PublicService (career page)
- `GetCompanyBySlug`, `GetPublishedJobsByCompanySlug`, `GetPublishedJobByID` (solo jobs `published`).
- **`ApplyToJob`** — crea/reusa Candidate por email dentro de la empresa + crea Application sin autenticación.

### 7.7 DashboardService (+ `dashboard_repository.go`, ~246 líneas de queries)
`GET /dashboard/stats` devuelve: totales de jobs por estado, total de candidatos, aplicaciones por stage, métricas del mes (nuevos candidatos/aplicaciones/contratados), **time-to-hire promedio**, **conversion rate**, tendencias diarias de 30 días, top jobs por aplicaciones y distribución por fuente.

### 7.8 Otros
- **SystemValueService** — `GetByCategory` y `GetByCategoryAndCompanyID` (globales + específicos de empresa vía `X-Company-ID`).
- **PlatformSettingsService** — lectura del singleton para branding público.
- **LocationService** — lecturas jerárquicas con preload selectivo (`include_states=true`...), búsqueda ILIKE case-insensitive, `GetLocationHierarchy`, `GetCountryByISO` (iso2/iso3). Tiempos típicos: países ~50ms, estados ~10ms, jerarquía completa ~150ms.
- **CompanyService** — CRUD + creación de directorios de uploads + `GetCompanyWithMembers`.

---

## 8. Base de Datos, Seeders y Consola

### 8.1 Migraciones
`AutoMigrate(AllModels...)` al iniciar el servidor (orden de modelos respeta FKs: Region → Subregion → Country → State → City, etc.). Soft deletes en todas las tablas vía `gorm.Model`.

### 8.2 Seeders (`internal/database/seeders/`, orquestados por `DatabaseSeeder.Run`)

| Seeder | Siembra |
|---|---|
| `role_seeder` | admin (50), recruiter (30), user (10) — roles de sistema |
| `plan_seeder` | **Free $0** (2u/3j/50c/100a/1GB) · **Starter $39.99** (trial 14, 5u/10j/200c/500a/5GB, export) · **Professional $79.99** (trial 14, 15u/50j/1000c/5000a/20GB, +brand/API/integraciones) · **Enterprise $159.99** (trial 30, todo -1, soporte dedicado) |
| `system_value_seeder` | 7 categorías de catálogos (~30 valores) |
| `platform_settings_seeder` | Singleton: "DVRA ATS", color sky-500, support@dvra.io, https://dvra.io |
| `user_seeder` | `superadmin@dvra.com` / `SuperAdmin123!` (sin empresa) y admin demo `admin@azentic.com` / `Admin123!` |
| `company_seeder` | Empresa demo (Azentic Sys) |

### 8.3 Consola y Makefile

```bash
go run cmd/console/main.go migrate [--seed] [--fresh]
go run cmd/console/main.go seed

make run          # servidor en :8080
make build        # binario
make swagger      # regenerar docs (swag init -g cmd/dvra-api/main.go -o docs)
make db-migrate   # migraciones
make db-seed      # seeders
make db-fresh     # drop + migrate + seed (~7 s)
make db-location  # carga masiva de ubicaciones (~157k filas, 3–5 min, script SQL)
make fmt / vet / test
```

---

## 9. Seguridad

### 9.1 Implementado
- bcrypt para passwords (default cost).
- JWT firmado (HS256) con secrets separados para access/refresh.
- Aislamiento multi-tenant auditado (ver §6 — corrigió una vulnerabilidad crítica donde los listados devolvían datos de todas las empresas).
- Forzado de `company_id` desde el token en todas las creaciones.
- CORS restringido por configuración.
- Soft deletes (sin pérdida de historial; recuperación posible).
- Docker con usuario no-root y build multi-stage.
- Suspensión de empresas (`plan_tier = "suspended"`) bloquea el login de sus usuarios.

### 9.2 Pendiente (recomendaciones de la auditoría)
- `PUT/DELETE /users/:id`: validar memberships de la empresa antes de operar.
- Logging de intentos de acceso cross-company (warn con user/empresa solicitada).
- Rate limiting (especialmente en rutas públicas y login).
- Tests unitarios de los filtros multi-tenant.
- Rotar credenciales seed (`superadmin@dvra.com`, `admin@azentic.com`) fuera de desarrollo.

---

## 10. Swagger / OpenAPI

- **UI:** `http://localhost:8080/swagger/index.html` (URL del JSON dinámica según `PORT`).
- **Regenerar:** `make swagger` cada vez que se agregan/modifican endpoints.
- **Anotar handlers** con comentarios Swaggo (`@Summary`, `@Tags`, `@Param`, `@Success`, `@Failure`, `@Security BearerAuth`, `@Router`). Los DTOs se documentan solos a partir de sus tags `json` y `binding`.
- **Probar endpoints autenticados:** login en la UI → botón *Authorize* → `Bearer <token>`.
- Archivos generados (`docs/docs.go`, `swagger.json`, `swagger.yaml`) — no editar a mano.
- Metadata general (título, versión, host, securityDefinitions) en las anotaciones raíz de `cmd/dvra-api/main.go`.

---

## 11. Despliegue

### 11.1 Docker (actual)

- **Dockerfile** multi-stage: `golang:1.24-alpine` (build, `CGO_ENABLED=0`) → `alpine` runtime con usuario no-root. Expone 8080.
- **docker-compose:** servicios `api` + `postgres:16-alpine` con healthcheck y volumen persistente `postgres_data`.

```bash
docker compose up -d
```

### 11.2 Infraestructura objetivo (AWS, roadmap)

```
Route53 → CloudFront (SSL/CDN) → ALB (multi-AZ) → ECS Fargate ×2 → RDS PostgreSQL (Multi-AZ)
S3: dvra-resumes (CVs) + dvra-assets | SendGrid/SES | CloudWatch + Sentry | Secrets Manager
```

Costo estimado inicial: ~$95–215/mes. Estrategia: Lightsail/EC2 al inicio → ECS al escalar. Si se llega a 1000+ empresas: sharding por `company_id` o BD dedicada para Enterprise.

### 11.3 Observabilidad (roadmap)
Logging estructurado (zap), Sentry para errores, métricas Prometheus (`http_requests_total`), backups diarios de BD, status page público. Caching con Redis para sesiones/planes públicos cuando el volumen lo justifique.

---

## 12. Estado Actual y Deuda Técnica

### ✅ Implementado y funcional
Multi-tenancy con aislamiento auditado · JWT + refresh + switch-company · registro transaccional de empresas · CRUD completo de todos los recursos · pipeline con timestamps automáticos y rating · dashboard analítico · career page pública con aplicación anónima · planes con features granulares · catálogos dinámicos · ubicaciones jerárquicas con búsqueda · seeders · Swagger · Docker.

### ⚠️ Deuda técnica priorizada

| # | Pendiente | Detalle |
|---|---|---|
| 1 | **Enforcement de límites de plan** | `MaxUsers/MaxJobs/...` existen en el modelo pero **ningún service los valida** al crear recursos. Falta un `PlanService.CheckLimit(companyID, resource)` invocado desde Create*/Publish |
| 2 | **Validación de trial expirado** | `TrialEndsAt` + `IsTrialActive()` existen; falta el middleware que actúe al expirar (downgrade a free / solo lectura) |
| 3 | **Validación de uploads** | Sin límite de tamaño/tipo de archivo ni tracking de quota de storage |
| 4 | **`PUT/DELETE /users/:id`** | Falta validar pertenencia vía memberships (ver §9.2) |
| 5 | **Rate limiting** | No implementado |
| 6 | **Transacciones** | Solo `RegisterCompany` y asignación de planes son transaccionales |
| 7 | **Tests** | Sin tests unitarios ni de integración |
| 8 | **Emails transaccionales** | SendGrid/SES no integrado (confirmaciones de aplicación, invitaciones, notificaciones de stage) |
| 9 | **Billing** | Stripe no integrado; cambios de plan son manuales (SuperAdmin) |
| 10 | **Sistema de invitaciones** | Fase 2: `POST /memberships/invite` + `POST /auth/accept-invite` con tokens de expiración, para que Admins inviten usuarios existentes sin pasar por SuperAdmin |

### 📝 Notas de coherencia documental
- Los precios de planes citados en documentos antiguos ($29.99/$89.99/$149.99 y $49/$149/$399) **no coinciden** con el seeder actual ($39.99/$79.99/$159.99). La fuente de verdad es `plan_seeder.go`.
- `API_ENDPOINTS.md` (eliminado) listaba estados de aplicación `pending/reviewing/interviewed` — obsoletos; los stages reales son los de §3.2.
- Las rutas `/admin/*` y `/auth/superadmin/login` documentadas antes ya **no existen en este API** (separación del SuperAdmin).
