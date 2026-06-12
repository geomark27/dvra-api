# Bitácora de Actividades — Dvra API

> Registro cronológico de implementaciones, auditorías y decisiones técnicas.
> Consolida los antiguos `MULTI_COMPANY_IMPLEMENTATION.md`, `SECURITY_AUDIT.md` y `SUPERADMIN_IMPLEMENTATION.md`.
> Las entradas más recientes van al inicio. Para detalle vigente de cada módulo, ver `docs/md/tecnico/`.

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
