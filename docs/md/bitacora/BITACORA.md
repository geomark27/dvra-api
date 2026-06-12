# Bitácora de Actividades — Dvra API

> Registro cronológico de implementaciones, auditorías y decisiones técnicas.
> Consolida los antiguos `MULTI_COMPANY_IMPLEMENTATION.md`, `SECURITY_AUDIT.md` y `SUPERADMIN_IMPLEMENTATION.md`.
> Las entradas más recientes van al inicio. Para detalle vigente de cada módulo, ver `docs/md/tecnico/`.

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
