# Bitácora de Actividades — Dvra API

> Registro cronológico de implementaciones, auditorías y decisiones técnicas.
> Consolida los antiguos `MULTI_COMPANY_IMPLEMENTATION.md`, `SECURITY_AUDIT.md` y `SUPERADMIN_IMPLEMENTATION.md`.
> Las entradas más recientes van al inicio. Para detalle vigente de cada módulo, ver `docs/md/tecnico/`.

---

## 2026-06-20 — Piloto monolito modular + hexagonal-lite: módulo `staffing`

**Contexto:** primer paso del ADR-001. Se migró `staffing` de la organización por capa técnica (`internal/app/{models,dtos,...}`) a un módulo autocontenido en `internal/modules/staffing/`, como piloto para validar el patrón antes de replicarlo.

**Estructura resultante:**
- `domain/ports.go` — interfaces (StaffingClientRepository, PlacementRepository) + el puerto cross-módulo `ApplicationFinder`/`HiredApplication`. No importa gin ni gorm.
- `service/` — casos de uso (StaffingClientService, PlacementService). Dependen solo de `domain`, `models`, `apperr`.
- `repository/` — adaptadores GORM con `*gorm.DB` **inyectado** (ya no usan el global `database.DB` → testeables).
- `transport/` — handlers HTTP + `routes.go` (adaptador de entrada).
- `module.go` — wiring del módulo + `RegisterRoutes`.

**Decisiones clave (del ADR-001):**
- **Modelos compartidos.** `models.StaffingClient`/`Placement` y los DTOs **se quedaron** en `internal/app/{models,dtos}`. Separarlos provocaría ciclos de importación por las relaciones GORM mutuas (`StaffingClient↔Job`, `Placement↔Application`). Es la 1ª pasada pragmática.
- **Puertos definidos por el consumidor para romper ciclos cross-módulo:**
  - `staffing` necesita datos de `Application` (recruitment) → define `domain.ApplicationFinder`; el composition root (`platform/server/wiring_adapters.go`) inyecta un adaptador sobre `ApplicationRepository`.
  - `recruitment` (job_service) necesita validar un cliente final → ahora depende de una interfaz local `staffingClientReader`; el composition root le pasa `staffingModule.ClientRepo`.
  - Resultado: **ningún módulo importa al otro**. Grafo acíclico verificado por `go build`. El único consumidor de `internal/modules/staffing` es el composition root.
- **Acople aceptado:** el módulo importa `app/services` solo por el *tipo* `PlanService` (entitlement, para `RequireFeature`). Podría convertirse en otro puerto del consumidor más adelante.

**Archivos eliminados** (reemplazados por el módulo): `services/{staffing_client,placement}_service.go`, `repositories/{staffing_client,placement}_repository.go`, `handlers/{staffing_client,placement}_handler.go`.

**Verificado:** `go build ./...`, `go vet ./...`, `gofmt -l` (limpio), `go test ./...` y build de ambos binarios pasan. Sin cambio de comportamiento ni de contrato de API (mismas rutas, mismos códigos).

**Pendientes / siguiente paso (del ADR-001):**
- [ ] Tests del módulo `staffing` (ya es trivial: repos con `db` inyectado + services mockeables vía puertos).
- [ ] Replicar el patrón al resto de dominios (recruitment, iam, billing, platform), uno por uno con build/tests verdes.

**Referencia vigente:** `docs/md/tecnico/ADR-001-arquitectura-modular.md`, `internal/modules/staffing/`

---

## 2026-06-20 — Unificación de `auth` y `plan` a `apperr`

**Contexto:** en la pasada anterior se dejaron `auth` y `plan` fuera porque ya mapeaban códigos correctamente vía errores centinela + ladders `if err == services.ErrX`. Para uniformar el manejo de errores en TODO el dominio, ahora también usan `apperr`.

**Qué se hizo:**
- **Nuevo constructor `apperr.Unauthorized` (401)** para los errores de credenciales de auth.
- **Centinelas convertidos a valores `*apperr.AppError`** (en `plan_service.go` y `auth_service.go`): siguen siendo variables exportadas comparables con `errors.Is`/`==` (compatibilidad total), pero ahora llevan su código HTTP. Mapeo: `ErrPlanNotFound`/`ErrUserNotFound`/`ErrCompanyNotFound` → 404; `ErrPlanSlugExists`/`ErrPlanInUse`/`ErrEmailExists` → 409; `ErrInvalidCredentials`/`ErrInvalidPassword` → 401; `ErrNoMembership` → 403; `ErrInvalidPlanData` → 400. Inline: `account is inactive` → 403, `invalid refresh token` → 401, `plan is not active` → 400.
- **Handlers simplificados** (`plan_handler`, `auth_handler`): se eliminaron los ladders `if err == services.ErrX { ... }` (≈12 bloques) a favor de `apperr.StatusCode(err)`. Menos código y misma (o más correcta) semántica de códigos.
- **Preservado por seguridad/config:**
  - `RefreshToken` mantiene su **401 genérico** para cualquier fallo (decisión deliberada: el frontend trata 401 en refresh como "cerrar sesión"; un 403/404 podría no disparar ese flujo).
  - `free plan is not available` sigue como error plano → 500 (es un problema de configuración del servidor, no del cliente).
  - Los mensajes de credenciales siguen siendo genéricos ("invalid email or password") — no revelan si el email existe.

**Verificado:** `go build ./...`, `go vet ./...`, `gofmt -l` (limpio), `go test ./...` y build de ambos binarios (`dvra-api`, `console`) pasan. Resultado: 14 services y 14 handlers usan `apperr` (solo quedan fuera `dashboard`/`health`, sin errores de negocio).

**Referencia vigente:** `internal/shared/apperr/errors.go`, `internal/app/services/{auth,plan}_service.go`, `internal/app/handlers/{auth,plan,platform_settings}_handler.go`

---

## 2026-06-20 — Adopción de `apperr` en todo el dominio (códigos HTTP correctos)

**Contexto:** tras introducir el paquete `internal/shared/apperr` en el módulo staffing, se extendió el mismo patrón al resto de la API. El problema transversal: los services devolvían errores genéricos (`fmt.Errorf`/`errors.New`) y los handlers respondían **HTTP 500 para errores que son del cliente** (no encontrado, duplicado, validación). Además había dos antipatrones: detección de tipo de error por comparación de strings (`err.Error() == "company not found"`, `strings.Contains(err, "already applied")`) y `strconv.ParseUint` ignorando el error de parseo.

**Qué se hizo:**
- **Services convertidos a `apperr`** (`NotFound`/`Conflict`/`BadRequest`/`Forbidden`): `job`, `candidate`, `company`, `user`, `application`, `membership`, `public`, `platform_settings`, `system_value`, `location`. Los errores de negocio ahora llevan su código HTTP; los fallos internos reales (wraps con `%w`: creación de directorios/candidato/aplicación) se mantienen como errores planos → siguen siendo 500.
- **Handlers actualizados** para usar `apperr.StatusCode(err)` en lugar de códigos hardcodeados tras llamadas al service: `job`, `candidate`, `company`, `user`, `application`, `membership`, `public`, `system_value`, `location` (+ los `staffing`/`placement` ya hechos). Resultado: slug/email duplicado → 409, recurso inexistente → 404, validación → 400.
- **Eliminado el matching frágil de strings**: `company_handler` y `user_handler` ya no comparan `err.Error() == "..."`; `public_handler` ya no usa `strings.Contains` para clasificar el error de postulación.
- **`strconv.ParseUint` validado** en los handlers que lo ignoraban (`job` ×5, `candidate` ×4, `application` ×5): ID no numérico → 400 en vez de un 404 confuso.
- **Mensaje genérico preservado donde corresponde**: en el endpoint público `ApplyToJob`, los errores internos (500) siguen devolviendo "Failed to submit application" para no filtrar detalles; solo los errores de negocio exponen su mensaje.

**Decisiones de alcance (lo que NO se tocó y por qué):**
- **`auth` y `plan` se dejaron intactos**: ya mapeaban códigos correctamente vía errores centinela (`errors.Is`/`==`). Además `auth` usa mensajes genéricos deliberados ("Invalid email or password") por seguridad; reescribirlos arriesgaba filtrar si un email existe. No tenían el bug del 500.
- **No se refactorizó el scoping multi-tenant** (`c.Get("company_id")` → `authctx`, ítems **C1/M1** de la auditoría) ni se eliminaron los double-lookups de los handlers genéricos: son cambios mayores y ortogonales a "aplicar apperr".

**Nota de comportamiento:** los handlers que ya devolvían `err.Error()` en 500 (la mayoría) lo siguen haciendo; ahora también lo hacen `company`/`user` (antes usaban mensajes genéricos). Es consistente con la norma del código, pero expone el texto del error en fallos 5xx — pendiente de endurecer si se requiere.

**Verificado:** `go build ./...`, `go vet ./...`, `gofmt -l` (limpio) y `go test ./...` (permissions) pasan.

**Referencia vigente:** `internal/shared/apperr/errors.go` y todos los `*_service.go` / `*_handler.go` listados arriba.

---

## 2026-06-20 — Correcciones al módulo Staffing / Placement (auditoría de calidad)

**Contexto:** revisión completa del módulo staffing/placement (modelos, DTOs, repositorios, servicios, handlers y rutas) para corregir bugs reales, inconsistencias de diseño y malos códigos HTTP detectados en auditoría. El código compilaba pero tenía errores de comportamiento silenciosos.

**Bugs corregidos:**

- **Error de compilación encubierto:** `placement_service.go` tenía un método `GetAll` con cuerpo vacío que no compilaría. Se eliminó.
- **Código HTTP incorrecto en errores de negocio:** `CreateStaffingClient` y `UpdateStaffingClient` devolvían HTTP 500 cuando el slug ya existía (conflicto de datos). Ahora devuelven 409 Conflict. El 500 indica fallo del servidor; el 409 indica que el cliente envió datos que chocan con datos existentes.
- **ID inválido en URL no se validaba:** si alguien llamaba `GET /staffing-clients/abc`, el sistema ignoraba el error de parseo, buscaba el registro con ID=0 y devolvía un 404 sin explicación. Ahora devuelve 400 Bad Request con mensaje claro, en los 6 handlers afectados (get, update, delete de ambos módulos).
- **Placement duplicado permitido:** se podía llamar `POST /placements` dos veces con la misma `application_id` y se creaban dos registros. Ahora el service verifica unicidad y devuelve 409 si ya existe un placement para esa application.
- **Método fantasma en repositorio:** `PlacementRepository` declaraba `GetAll` y `GetByCandidate` en su interface pero ningún service ni handler los llamaba. Código muerto en la interface confunde y obliga a las implementaciones futuras a cumplir contratos inútiles. Ambos fueron eliminados de la interface y sus implementaciones.
- **`deleted_at IS NULL` manual y redundante:** `GetAll` del repositorio filtraba explícitamente `deleted_at IS NULL`, algo que GORM ya hace solo con soft delete. Era redundante e inconsistente con el resto del repo.

**Inconsistencias de diseño corregidas:**

- **Doble consulta a BD en Update y Delete:** los handlers de `UpdateStaffingClient`, `DeleteStaffingClient`, `UpdatePlacement` y `DeletePlacement` hacían dos consultas separadas — una en el handler para verificar el tenant, y otra dentro del service para traer el registro. Ahora el service recibe el `companyID` directamente y hace una sola consulta que verifica todo. Se agregó la convención: `companyID = 0` significa SuperAdmin (omite validación de tenant).
- **Modelos crudos en el DTO de respuesta:** `PlacementResponseDTO` incluía `*models.StaffingClient` y `*models.Candidate` directamente — los modelos de base de datos expuestos en la API. Esto puede filtrar campos internos no deseados. Ahora usa `*StaffingClientResponseDTO` y `*CandidateResponseDTO` (DTOs de presentación).
- **`JobID` y `ApplicationID` nullable sin razón:** en el modelo `Placement`, ambos campos eran punteros (`*uint`, nullable en BD) aunque la lógica de negocio siempre los rellena al copiarlos de una `Application` (que los tiene como `not null`). Cambiados a `uint not null`. Esto fuerza la constraint en BD vía AutoMigrate y simplifica el código del service (ya no necesita `&jobID`).

**Nuevo paquete:**

- **`internal/shared/apperr`**: tipos de error con código HTTP asociado (`NotFound`, `Conflict`, `BadRequest`, `Forbidden` y `StatusCode(err)`). Permite que los services señalen el tipo de error y los handlers elijan el status code correcto sin inspeccionar strings de mensajes. Resuelve el patrón repetido de "todos los errores del service → 500".

**Versión:**

- Swagger `@version` en `main.go` sincronizado a `1.4.0` (estaba en `1.2.0`; la ruta raíz ya devolvía `v1.4.0`).

**Verificado:** `go build ./...` y `go vet ./...` pasan sin errores ni advertencias.

**Pendientes que siguen abiertos (de auditoría anterior):**
- [ ] Regenerar Swagger: `make swagger`
- [ ] BD ya poblada: planes Professional/Enterprise necesitan `can_use_staffing=true` manual en BD existente
- [ ] Tests del módulo staffing (ahora hay más superficies que probar: unicidad de application_id, apperr codes)
- [ ] Ítem de auditoría **A1** sigue abierto: 47 tags `validate` que no se ejecutan

**Referencia vigente:** `internal/shared/apperr/errors.go`, `internal/app/models/placement.go`, `internal/app/repositories/placement_repository.go`, `internal/app/services/{staffing_client,placement}_service.go`, `internal/app/handlers/{staffing_client,placement}_handler.go`, `internal/app/dtos/placement_dto.go`

---

## 2026-06-14 — Estándares de Desarrollo: documento + skills + enforcement + auditoría

**Contexto:** formalizar los estándares de Dvra (cómo se escribe y estructura el código) no solo por escrito, sino **operacionalizados** como skills de Claude que validan a la entrada y a la salida de cada tarea.

**Qué se hizo:**
- **Documento fuente de verdad** `docs/md/tecnico/ESTANDARES_DESARROLLO.md`: stack, cómo se escribe/estructura/verifica/documenta el código, **regla de oro multi-tenant** (`company_id` única frontera), autorización vs entitlement, proceso **HU → HT** y **Definición de Done**.
- **Skills** (`.claude/skills/`): `hu-a-ht` (puerta de entrada: traduce Historia de Usuario funcional a Historia Técnica validando estándares) y `validar-estandares` (puerta de salida: valida el diff contra la Definición de Done). Ambas anclan al documento como fuente de verdad.
- **Enforcement automático**: `.golangci.yml` (v1.64.8: errcheck, govet, ineffassign, staticcheck, unused, gofmt — sin `misspell` por falsos positivos en español, sin `goimports` local-prefixes por churn), `.editorconfig`, CI en `.github/workflows/ci.yml` (gofmt+vet+build+test+lint), hook `scripts/git-hooks/pre-commit`, y targets `make check`/`fmt-check`/`hooks-install` (+ `swag` en `install-tools`).
- **Baseline verde**: se cerró deuda mínima pre-existente — 4 `errcheck` (`CloseDB` en CLI, `UpdateLastLogin` best-effort) y 4 archivos sin `gofmt`. `make check` pasa entero.
- **Auditoría de brechas** `docs/md/tecnico/AUDITORIA_ESTANDARES.md`: hallazgos priorizados con evidencia numérica.

**Pendientes (del plan de remediación de la auditoría):**
- [ ] **A1** — Unificar validación: 47 tags `validate` en 6 DTOs no se ejecutan; migrar a `binding`.
- [ ] **C1 (Fase 4)** — Centralizar el scoping multi-tenant (41 usos de `IsSuperAdmin` repetidos); riesgo de fuga cross-tenant. Candidato a estrenar las skills.
- [ ] **M1** — Migrar 7 handlers de `c.Get("company_id")` a `authctx`.
- [ ] **M2** — Unificar repos hacia `db` inyectado (10 usan global `database.DB`); habilita tests.
- [ ] **A2** — Tests de services con validaciones (hoy 1 test / 16 services).
- [ ] **B1** — Check en CI de que Swagger esté regenerado.

**Referencia vigente:** `docs/md/tecnico/ESTANDARES_DESARROLLO.md`, `docs/md/tecnico/AUDITORIA_ESTANDARES.md`, `.claude/skills/`

---

## 2026-06-14 — Módulo Staffing / Outsourcing (StaffingClient + Placement)

**Contexto:** habilitar a Dvra para firmas de staffing/outsourcing (que gestionan talento para clientes finales), no solo empresas que reclutan para sí mismas. Caso motivador: una firma que presta personal a un cliente final. Se modeló como segmento objetivo sin romper el multi-tenancy existente.

**Decisión de arquitectura:** el cliente final NO es un tenant. `company_id` sigue siendo la ÚNICA frontera de aislamiento; el cliente final es una dimensión interna del tenant.

**Qué se hizo:**
- Nuevos modelos `StaffingClient` (cliente final dentro del tenant) y `Placement` (colocación de un candidato en un cliente final, con datos de contrato y billing). Reescritos para embeber `BaseModel` (antes usaban `gorm.Model`, sin tags JSON ni `TableName()`) y registrados en `AllModels` para AutoMigrate (antes estaban huérfanos: las tablas nunca se creaban).
- `Job.StaffingClientID *uint` opcional: `nil` = reclutamiento propio; con valor = modo staffing. Filtro `GET /jobs?staffing_client_id=` y exposición en respuestas.
- `Placement` **deriva de una `Application` en etapa `hired`**: `candidate_id`/`job_id` se copian de ella (no se confían del body). El service valida integridad cross-tenant: application y cliente final deben pertenecer al mismo `company_id`. Igual validación al asignar `Job.StaffingClientID`.
- CRUD completo de ambos (dto/repo/service/handler + rutas `/api/v1/staffing-clients` y `/api/v1/placements`), con scoping manual por `company_id` (mismo patrón que jobs) y slug único por empresa validado en service (no índice único, para no chocar con soft delete).
- **Entitlement por plan:** nuevo flag `Plan.CanUseStaffing` + `Plan.HasFeature("staffing")` + método `PlanService.CompanyHasFeature()`. Nuevo middleware `RequireFeature(planService, "staffing")` (ortogonal a `RequirePermission`: la acción exige plan que lo incluya Y rol que lo permita). Seeder: Professional y Enterprise lo traen activo.
- Nuevos permisos en `permissions/staffing.go` (`staffing_clients.*`, `placements.*`) con grants por rol (delete solo admin, igual que jobs; `user` no ve placements por el billing).
- Anotaciones Swagger en los 10 endpoints nuevos.
- Verificado con `go build ./...`, `go vet ./...` y tests (`permissions` pasa).

**Pendientes:**
- [ ] Regenerar Swagger: `make swagger` (requiere instalar `swag`; no estaba en el entorno al implementar)
- [ ] BD ya poblada: `SeedPlans` omite planes existentes, así que las filas Professional/Enterprise NO obtienen `can_use_staffing=true` solas (AutoMigrate agrega la columna con default `false`). Actualizar vía `PUT /plans/:id` o SQL. En BD fresca (`make fresh`) quedan correctas.
- [ ] `PUT /jobs/:id` no permite desasignar el cliente (mandar `staffing_client_id: null` se interpreta como "no enviado") — requiere flag explícito si se necesita.
- [ ] Tests unitarios del flujo staffing (validaciones de integridad y entitlement)
- [ ] Evaluar control field-level sobre pay/bill rate (datos sensibles) por rol

**Referencia vigente:** `internal/app/models/{staffing_client,placement}.go`, `internal/shared/permissions/staffing.go`, `internal/shared/middleware/feature_middleware.go`

---

## 2026-06-12 — Sistema de Permisos por Acción (RBAC en código)

**Qué se hizo:**
- Nuevo paquete `internal/shared/permissions`: matriz rol → permisos en código, organizada por módulos (jobs.go, candidates.go, applications.go, users.go, companies.go, memberships.go, dashboard.go), análogo a un seeder de permisos por módulo pero verificado en compile-time. Punto único de consulta: `permissions.Can(role, permission)` — si a futuro se necesitan roles personalizados por empresa, solo se cambia la implementación de `Can()` a BD con caché, sin tocar rutas ni constantes.
- Nuevo middleware `RequirePermission(perm)` aplicado a TODAS las rutas protegidas en `routes.go` (estilo declarativo, como `can:` de Laravel). Antes ningún endpoint distinguía roles dentro de la empresa: un rol `user` (solo lectura) podía crear jobs o mover el pipeline.
- Permisos no asignados a ningún rol = exclusivos de SuperAdmin (companies.create/delete, memberships.create según RN-MEMB-004).
- Nuevo paquete `internal/shared/authctx`: accesores tipados al contexto (`Role`, `IsSuperAdmin`, `UserID`, `CompanyID`). Refactor de los 34 `if role == "superadmin"` dispersos en 8 handlers → `authctx.IsSuperAdmin(c)`; el literal "superadmin" ya no existe fuera de `permissions`/`authctx`.
- Fix: `getRoleLevel` no incluía `superadmin` (habría recibido nivel 0 con `RequireRole`); ahora nivel 100.
- `GET /auth/me` ahora devuelve `role` y `permissions` de la empresa activa para que el frontend muestre/oculte acciones (mismo contrato que con Spatie).
- Primeros tests unitarios del proyecto: `permissions_test.go` valida la matriz por rol.

**Pendientes:**
- [ ] Chequeos a nivel de recurso (hiring_manager edita solo sus jobs asignados) — requiere tabla de asignaciones (RN-MEMB-007)
- [ ] Refactorizar `c.Get("company_id")` en handlers a `authctx.CompanyID(c)`
- [ ] Cuando exista invitaciones (Fase 2): otorgar `memberships.invite` a admin

**Referencia vigente:** `negocio/LOGICA_DE_NEGOCIO.md` sección 3.2 (matriz) y `internal/shared/permissions/`

---

## 2026-06-12 — BaseModel unificado + diseño de perfiles internos/externos

**Qué se hizo:**
- Nuevo `models.BaseModel` (`internal/app/models/base.go`): equivalente a `gorm.Model` pero con tags JSON en snake_case (`id`, `created_at`, `updated_at`, `deleted_at`).
- Refactor de TODOS los modelos para embeber `BaseModel`: los 6 que duplicaban los campos inline (User, Company, Job, Candidate, Application, Membership) y los que usaban `gorm.Model` (Plan, Role, ubicaciones, SystemValue, PlatformSettings). No se usó `gorm.Model` directo porque no tiene tags JSON y los primeros 6 se serializan directamente en las respuestas de la API.
- `BaseModel` genera las mismas columnas que `gorm.Model`, así que no hay impacto en BD; excepción: `platform_settings` gana `created_at` y `deleted_at` (AutoMigrate las agrega; el soft delete es inerte porque el singleton nunca se elimina).
- Convención a futuro: todo modelo nuevo embebe `BaseModel`, nunca `gorm.Model` ni campos inline.
- Verificado con `go build ./...`, `go vet ./...` y tests.
- Documentado el diseño de **perfiles internos vs externos** (recruiters freelance/agencias) en `negocio/LOGICA_DE_NEGOCIO.md` sección 3.4: RN-MEMB-006 (`member_type`), RN-MEMB-007 (alcance del externo), RN-MEMB-008 (atribución y reporting por recruiter), y su incorporación al roadmap Fase 2.

**Pendientes:**
- [ ] Implementar invitaciones por email (prerequisito, ya era Fase 2)
- [ ] Implementar `member_type` + asignación vacante ↔ recruiter
- [ ] Historial de pipeline con atribución
- [ ] Reportes por recruiter

**Referencia vigente:** `negocio/LOGICA_DE_NEGOCIO.md` (sección 3.4)

---

## 2025-12-22/23 — Módulo de Ubicaciones (Locations)

**Qué se hizo:**
- Módulo completo de datos geográficos: Region → Subregion → Country → State → City.
- Capas completas: models, DTOs, repository, service, handler y rutas.
- Rutas GET públicas (sin auth); CRUD restringido a SuperAdmin (`/admin/locations/*`).
- Seeder básico incluido en `make fresh` (~7s) y carga masiva opcional con `make db-location` (157k+ registros).
- Comando `make fresh` corregido y Makefile actualizado.

**Referencia vigente:** `tecnico/MODULO_UBICACIONES.md`

---

## 2025-12-09 — Módulo de Planes (Plans)

**Qué se hizo:**
- CRUD completo de planes de suscripción gestionado por SuperAdmin.
- Rutas públicas para pricing page (`GET /plans`, `GET /plans/:slug`).
- Asignación de plan a empresa (`POST /admin/plans/assign`) actualizando `company.plan_tier`.
- Seeder con 4 planes: Free ($0), Starter ($39.99), Professional ($79.99), Enterprise ($159.99).
- Validaciones: slug único, no eliminar planes en uso (soft delete), planes validados contra BD (sin strings hardcodeados) al crear empresas.

**Referencia vigente:** `tecnico/MODULO_PLANES.md`

---

## 2025-12-08 — Separación Admin / SuperAdmin

**Qué se hizo:**
- Reorganización de la arquitectura para separar Admin regular (scoped a empresa) de SuperAdmin (acceso global, sin `company_id` en el JWT).
- Nuevos archivos: `dtos/superadmin_dto.go`, `services/admin/superadmin_companies_service.go`, `handlers/admin/superadmin_companies_handler.go`.
- Middleware `RequireSuperAdmin()`: valida JWT + role `superadmin` + ausencia de `company_id`.
- Rutas `/api/v1/admin/*`: listar/crear empresas, cambiar plan, suspender, ver usuarios de cualquier empresa, analytics globales (MRR, churn, totales).
- Creación de empresa + admin + membership en una sola transacción.

**Pendientes registrados en su momento:**
- [ ] Testing automatizado de handlers SuperAdmin
- [ ] Audit logs de acciones del SuperAdmin
- [ ] Impersonation (`POST /admin/impersonate/:user_id`)
- [ ] ⚠️ Cambiar credenciales por defecto del SuperAdmin en producción

**Referencia vigente:** `tecnico/FLUJO_SUPERADMIN.md`

---

## 2025-12-08 — Auditoría de Seguridad Multi-Tenancy

**Problema detectado (crítico):** los endpoints devolvían datos de TODAS las empresas sin filtrar por `company_id`; usuarios de una empresa podían ver, crear, modificar y eliminar recursos de otras.

**Correcciones aplicadas:**
- **Lectura:** filtrado por `company_id` del token en listados de memberships, users, companies, jobs, candidates y applications. SuperAdmin (sin `company_id`) conserva acceso global.
- **Lectura individual:** `GET /:id` valida que el recurso pertenezca a la empresa del token (403 si no).
- **Creación:** `POST /jobs|candidates|applications` ignora el `company_id` del body y fuerza el del token.
- **Actualización/Eliminación:** verifican pertenencia del recurso antes de proceder. `DELETE /companies/:id` solo SuperAdmin.
- **Memberships (MVP):** solo SuperAdmin puede crear memberships (`POST /admin/memberships`); admins de empresa reciben 403.
- Métodos nuevos: `UserRepository.GetByCompanyID()`, `UserService.GetUsersByCompanyID()`.

**Pendientes registrados en su momento:**
- [ ] `PUT /users/:id` y `DELETE /users/:id`: falta verificar memberships
- [ ] Tests unitarios de filtrado por empresa
- [ ] Logging de intentos de acceso cross-company
- [ ] Rate limiting contra enumeración

**Referencia vigente:** `tecnico/FLUJO_CLIENTE.md` (sección de aislamiento multi-tenant)

---

## 2025-12-08 — Sistema Multi-Empresa (v1.3.0)

**Qué se hizo:**
- Registro de empresa con su primer admin en una transacción: `POST /auth/register-company` (Company + User + Membership con role admin e `is_default: true`). `POST /auth/register` quedó deprecated.
- Login devuelve lista de empresas del usuario; el token usa la empresa por defecto.
- `POST /auth/switch-company`: genera nuevo token con distinto `company_id` (y el rol correspondiente) sin re-login.
- `GET /auth/my-companies`: lista de empresas del usuario para el selector de la UI.
- Un usuario puede tener N memberships con roles distintos por empresa (ej. admin en una, recruiter en otra).
- Errores nuevos: `ErrCompanyNotFound`, `ErrNoMembership`.

**Pendientes registrados en su momento:**
- [ ] **Fase 2 (prioritario):** sistema de invitaciones por email (`POST /memberships/invite`, `POST /auth/accept-invite`) para que admins de empresa puedan invitar usuarios sin depender del SuperAdmin
- [ ] Fase 3: onboarding post-registro (wizard, configuración inicial, importación de datos)

**Referencia vigente:** `tecnico/FLUJO_CLIENTE.md` y `tecnico/API_ENDPOINTS.md`

---

## Formato para nuevas entradas

```markdown
## AAAA-MM-DD — Título corto

**Qué se hizo:**
- ...

**Pendientes:**
- [ ] ...

**Referencia vigente:** `tecnico/...`
```
