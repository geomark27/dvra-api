# Dvra — Flujo Completo de Uso de la Aplicación

> **Cómo se usa el sistema de punta a punta, por actor y por proceso**
> Consolidado: Junio 2026 | Reemplaza a: `CLIENT_FLUJO_COMPLETO.md`, `SUPERADMIN_FLUJO_COMPLETO.md`, `SUPERADMIN.md`, `LOCATION_USAGE.md`
>
> 📌 Documentos relacionados:
> - Reglas de negocio que gobiernan estos flujos → [01_LOGICA_DE_NEGOCIO.md](./01_LOGICA_DE_NEGOCIO.md)
> - Plan operativo → [02_PLAN_DE_NEGOCIO.md](./02_PLAN_DE_NEGOCIO.md)
> - Especificación técnica de cada endpoint → [04_DOCUMENTACION_TECNICA_API.md](./04_DOCUMENTACION_TECNICA_API.md)

---

## 1. Actores del Sistema

| Actor | ¿Quién es? | Contexto en el sistema |
|---|---|---|
| **Cliente (usuario de empresa)** | Cualquier usuario con membresía en una o más empresas: Admin, Recruiter, Hiring Manager o User | JWT **con** `company_id`; ve solo datos de su empresa activa |
| **Candidato (visitante público)** | Persona que aplica a una vacante desde la career page | Sin cuenta ni autenticación; usa rutas `/public/*` |
| **SuperAdmin** | Operador de la plataforma (gestiona empresas, planes, suspensiones, analytics globales) | JWT **sin** `company_id`. ⚠️ Desde el commit "separación del superadmin", su panel vive en un **servicio separado**; este API conserva solo los checks de rol residuales |

### Cliente vs SuperAdmin — referencia rápida

| Característica | Cliente | SuperAdmin |
|---|---|---|
| `company_id` en JWT | ✅ Requerido | ❌ NULL |
| Alcance de datos | Solo su empresa | Todas las empresas |
| Crear empresas / gestionar planes | ❌ | ✅ |
| Multi-empresa (switch) | ✅ | N/A |
| Límites de plan | ✅ Aplican | ❌ |
| Crear jobs/candidatos/aplicaciones | ✅ | ❌ (no tiene empresa) |

---

## 2. Flujo 1 — Registro y Onboarding de una Empresa

### 2.1 Auto-registro (flujo principal)

```http
POST /api/v1/auth/register-company
{
  "company_name": "Mi Startup Tech",
  "company_slug": "mi-startup-tech",
  "admin_email": "ceo@mistartup.com",
  "admin_password": "SecurePass123!",
  "admin_first_name": "Juan",
  "admin_last_name": "Pérez",
  "timezone": "America/Bogota"
}
```

**Qué hace el sistema (en una sola transacción):**
1. Valida que el email no exista y que el plan `free` exista y esté activo.
2. Crea la **Company** (plan `free`, trial de 1 mes).
3. Crea el **User** administrador (password con bcrypt).
4. Crea la **Membership** (user → company, rol `admin`, `is_default = true`).
5. Devuelve `access_token` + `refresh_token` ya con `company_id` → el usuario puede trabajar de inmediato.

Si cualquier paso falla, se hace rollback completo (no quedan empresas o usuarios huérfanos).

> El endpoint `POST /auth/register` (usuario sin empresa) está **DEPRECADO**: crea un usuario que no puede acceder a nada hasta que un SuperAdmin lo asigne a una empresa.

### 2.2 Onboarding visual previsto (frontend, roadmap Q1)

1. Signup → empresa + admin creados.
2. Wizard: nombre/logo/industria → tamaño/ubicación/timezone → selección de plan (trial 14 días) → método de pago (Stripe).
3. Tutorial guiado: "crea tu primer job".
4. Invitar al equipo (opcional).
5. Dashboard con quick tips. **Tiempo estimado: 5–10 minutos.**

---

## 3. Flujo 2 — Autenticación y Multi-Empresa

### 3.1 Login

```http
POST /api/v1/auth/login
{ "email": "ceo@mistartup.com", "password": "SecurePass123!" }
```

- Verifica credenciales y que la cuenta esté activa.
- Busca la membresía por defecto (`is_default = true`) y genera el token con ese `company_id` y el rol correspondiente.
- Devuelve además **la lista de todas las empresas del usuario** (para mostrar selector en el frontend).
- Actualiza `last_login_at`.

**Estructura del token de un cliente:**
```json
{ "user_id": 25, "company_id": 5, "email": "ceo@mistartup.com", "role": "admin", "exp": ... }
```

### 3.2 Cambiar de empresa (switch company)

Un usuario freelance puede ser Admin en la empresa A y Recruiter en la B, con las mismas credenciales:

```http
POST /api/v1/auth/switch-company        { "company_id": 12 }
```

1. Valida que exista membresía activa del usuario en esa empresa (si no: 403/404).
2. Genera un **nuevo token** con el nuevo `company_id` y el **rol de esa empresa** (puede ser distinto).
3. A partir de ahí, todos los requests ven solo datos de la empresa 12.

Endpoints de apoyo:
- `GET /auth/my-companies` — lista de empresas del usuario (selector de empresa).
- `GET /auth/me` — datos del usuario autenticado.
- `POST /auth/refresh` — renovar access token (refresh token dura 30 días; access 1 hora).
- `POST /auth/change-password`, `POST /auth/logout` (JWT stateless: el cliente descarta los tokens).

---

## 4. Flujo 3 — Gestión del Equipo

### 4.1 Crear usuarios del equipo (Admin)

```http
POST /api/v1/users
{ "email": "recruiter@mistartup.com", "password": "TempPass123!", "first_name": "Ana", "last_name": "Gómez", "role": "recruiter" }
```

- El sistema crea el User **y su Membership automáticamente** en la empresa del token (el `company_id` se fuerza desde el JWT; cualquier valor enviado en el body se ignora).
- El nuevo usuario hace login con sus credenciales y ya trabaja en la empresa.

### 4.2 Gestionar membresías

| Acción | Endpoint | Quién |
|---|---|---|
| Ver miembros de la empresa | `GET /memberships`, `GET /users` | Todos los roles |
| Cambiar rol de un miembro | `PUT /memberships/:id` | Admin (no puede cambiar su propio rol) |
| Remover de la empresa | `DELETE /memberships/:id` | Admin |
| **Asignar usuario EXISTENTE a una empresa** | `POST /memberships` | ❌ Bloqueado para clientes (403). **Solo SuperAdmin** (MVP). Fase 2: invitaciones por email |

Al remover una membresía: el usuario pierde acceso a esa empresa; si tiene otras empresas, sigue en ellas.

### 4.3 Estructura típica de un equipo

```
1. Admin crea recruiters       → role: "recruiter"
2. Admin crea hiring managers  → role: "hiring_manager" (tech leads, CTO)
3. Admin asigna jobs:          PUT /jobs/:id { "assigned_recruiter": 30, "hiring_manager": 35 }
4. Cada quien ve/gestiona según su rol (ver matriz en 01_LOGICA_DE_NEGOCIO.md §3.2)
```

---

## 5. Flujo 4 — Ciclo de Vida de un Job

```
        crear                publicar              cerrar
  ───► [draft] ──────────► [published] ─────────► [closed]
         POST /jobs       PATCH /jobs/:id/publish  PATCH /jobs/:id/close
```

1. **Crear** — `POST /jobs` (Admin/Recruiter). Nace en `draft`. Campos: título, descripción, salario min/max, requisitos, beneficios, modalidad (`remote`/`onsite`/`hybrid`), ciudad (catálogo de ubicaciones), recruiter y hiring manager asignados.
2. **Publicar** — `PATCH /jobs/:id/publish`. Valida que tenga título y descripción. El job queda **visible en la career page pública** de la empresa.
3. **Recibir aplicaciones** — manuales (recruiter) o públicas (candidatos desde la career page).
4. **Cerrar** — `PATCH /jobs/:id/close` cuando la posición se llena. Deja de aceptar aplicaciones.
5. **Eliminar** — `DELETE /jobs/:id` es **soft delete**: el historial de aplicaciones se conserva siempre.

Listado y filtros: `GET /jobs` (de la empresa), por estado vía service.

---

## 6. Flujo 5 — Candidatos

### 6.1 Alta manual (Recruiter/Admin)

```http
POST /api/v1/candidates
{
  "email": "maria@example.com", "first_name": "María", "last_name": "López",
  "phone": "+57 300 1234567", "linkedin_url": "...", "github_url": "...",
  "source": "referral"
}
```

- Email único **dentro de la empresa** (409 si ya existe).
- `source` alimenta la analítica de canales (linkedin, website, referral, job_board, direct, other).

### 6.2 Subir CV

```http
POST /api/v1/candidates/:id/upload-resume   (multipart/form-data)
```

El archivo se guarda en el storage de la empresa (`uploads/companies/{slug}/`).

### 6.3 Alta automática (career page)

Cuando un candidato aplica desde el portal público, el sistema crea el Candidate (o lo reutiliza si el email ya existe en esa empresa) + la Application. Ver Flujo 7.

---

## 7. Flujo 6 — Pipeline de Aplicaciones

### 7.1 Crear aplicación

```http
POST /api/v1/applications     { "job_id": 10, "candidate_id": 50 }
```

Validaciones: job y candidato deben pertenecer a la empresa; sin duplicados al mismo job. Nace en stage `applied` con `applied_at` automático.

### 7.2 Mover por el pipeline

```http
PATCH /api/v1/applications/:id/move    { "stage": "technical" }
```

```
applied → screening → technical → interview → offer → hired
   └────────┴────────────┴───────────┴──────────┴──► rejected
```

- Al mover a `rejected` → se setea `rejected_at` automáticamente.
- Al mover a `hired` → se setea `hired_at` automáticamente. 🎉

### 7.3 Calificar y anotar

```http
PATCH /api/v1/applications/:id/rate    { "rating": 5 }      # 1–5 estrellas
PUT   /api/v1/applications/:id         { "notes": "Excelente comunicación, 5 años en Go" }
```

### 7.4 Vista Kanban

```http
GET /api/v1/applications/by-stage
```

Devuelve las aplicaciones de la empresa agrupadas por stage — es la fuente del tablero Kanban del frontend.

### 7.5 Proceso de screening típico (recruiter)

1. Abre el pipeline filtrado por `stage = applied`, ordenado por fecha.
2. Revisa perfil: CV, LinkedIn, GitHub, notas previas.
3. Agenda y realiza screening call; toma notas en la app.
4. Decisión: rating 3–5 → avanza a `screening`/`technical`; rating 1–2 → `rejected`.
5. (Roadmap) Email automático al candidato según avance o rechazo.

### 7.6 Decisión de contratación

1. Candidato llega a `interview`/`offer`: hiring manager + recruiter entrevistan; todos califican y anotan.
2. Se prepara la oferta (salario, beneficios, fecha de inicio).
3. Candidato acepta → `hired` (timestamp automático) | rechaza → `rejected`.
4. Si era la única posición, el job se cierra (`PATCH /jobs/:id/close`).

---

## 8. Flujo 7 — Candidato Aplica desde la Career Page (público, sin login)

```
1. Descubrimiento   GET /public/companies/:slug          → branding e info de la empresa
                    GET /public/companies/:slug/jobs     → vacantes publicadas
2. Ver detalle      GET /public/jobs/:id                 → título, descripción, requisitos, salario
3. Aplicar          POST /public/jobs/:id/apply
                    { "email", "first_name", "last_name", "phone", "linkedin_url", ... }
4. Sistema          → crea/reusa Candidate + crea Application (stage: applied)
5. (Roadmap)        → email de confirmación al candidato + notificación al recruiter
```

También es público `GET /public/platform-settings` (branding global de la plataforma: nombre, logo, colores, emails de contacto) para renderizar la página sin autenticación.

**Tiempo estimado para el candidato: 3–5 minutos.**

---

## 9. Flujo 8 — Dashboard y Seguimiento

```http
GET /api/v1/dashboard/stats
```

Una sola llamada entrega para la empresa activa:

- Totales: jobs (total/activos/draft/cerrados), candidatos.
- Aplicaciones por stage (applied, screening, technical, offer, hired…).
- Métricas del mes: nuevos candidatos, nuevas aplicaciones, contratados.
- **Time-to-hire promedio** y **conversion rate**.
- Tendencias diarias de los últimos 30 días.
- Top jobs por número de aplicaciones y distribución por fuente de candidatos.

---

## 10. Flujo 9 — Catálogos y Ubicaciones (soporte de formularios)

### 10.1 System Values (selects dinámicos)

```http
GET /api/v1/system-values/:category        # público
```

Categorías sembradas: `job_status`, `application_status`, `contract_type`, `work_mode`, `experience_level`, `priority`, `candidate_source`. Los valores globales tienen `company_id = null`; pueden existir valores específicos por empresa.

**Uso en frontend:** alimentar todos los dropdowns sin hardcodear valores.

### 10.2 Ubicaciones (jerarquía geográfica, público y read-only)

```
Region → Subregion → Country → State → City
```

Patrón típico — **select en cascada** en formularios de registro/jobs:

```http
GET /locations/countries                      # 1. cargar países
GET /locations/states?country_id=142          # 2. estados del país elegido
GET /locations/cities?state_id=3456           # 3. ciudades del estado elegido
```

Extras: búsqueda global (`GET /locations/search?q=london`), país por ISO (`GET /locations/countries/iso/MX`), jerarquía completa (`GET /locations/hierarchy/:id`). Datos base sembrados; carga masiva opcional de ~157,000 ubicaciones con `make db-location`.

**Recomendaciones frontend:** cachear países (no cambian), debounce de 300ms en búsquedas, manejar países sin estados/ciudades.

---

## 11. Flujo 10 — Planes y Límites

### 11.1 Pricing page (público)

```http
GET /api/v1/plans            # planes activos y públicos
GET /api/v1/plans/:slug      # detalle de un plan
```

### 11.2 Cómo viven los límites en la operación

1. La empresa tiene `plan_tier` (slug del plan: free, starter, professional, enterprise).
2. El frontend consulta `GET /plans/:slug` para mostrar límites y features del plan de la empresa.
3. Comportamiento esperado al alcanzar un límite (ej. plan Free = 3 jobs):

```
POST /jobs (4to job) →
{ "error": "Job limit reached",
  "message": "Your Free plan allows maximum 3 active jobs. Upgrade to create more." }
```

> ⚠️ **Estado actual:** este enforcement está definido como regla de negocio pero **todavía no implementado en el API** (ver deuda técnica en [04_DOCUMENTACION_TECNICA_API.md](./04_DOCUMENTACION_TECNICA_API.md) §12).

4. Para hacer upgrade en el MVP: la empresa contacta a soporte y el **SuperAdmin** cambia el plan (pago offline). Stripe self-service es parte del roadmap Q1.

---

## 12. Flujo 11 — Operación del SuperAdmin

> ⚠️ **Importante:** el panel SuperAdmin fue **separado de este API** (commit `e04aafd`). Los flujos siguientes describen la operación de plataforma tal como fue diseñada e implementada originalmente; hoy se ejecutan desde el servicio/proyecto SuperAdmin dedicado. En este API quedan: el seeder del usuario (`superadmin@dvra.com`), el soporte de `company_id = NULL` en JWT/membresías, y los checks `role == "superadmin"` en los handlers (un token superadmin emitido por el servicio correspondiente obtiene lectura global).

### 12.1 Capacidades del SuperAdmin

| Categoría | Puede |
|---|---|
| Empresas | Ver todas (con stats), crear empresa + admin inicial, cambiar plan, suspender |
| Planes | CRUD completo, activar/desactivar, crear planes custom privados, asignar a empresas |
| Memberships | Asignar usuarios existentes a empresas (única vía en MVP) |
| System Values | CRUD de catálogos globales y por empresa |
| Analytics | Métricas globales: empresas totales/activas/suspendidas, usuarios, jobs, aplicaciones, MRR, churn |

Y **no puede**: crear jobs/candidatos/aplicaciones (sin `company_id`), pertenecer a una empresa, ni usar el login normal de clientes.

### 12.2 Casos de operación típicos

**Onboarding gestionado:** cliente firma contrato → SuperAdmin crea empresa con admin y plan elegido → admin recibe credenciales → login y a trabajar.

**Upgrade manual:** cliente paga offline → SuperAdmin cambia `plan_tier` → límites nuevos aplican de inmediato.

**Morosidad:** 90 días sin pago → SuperAdmin suspende (`plan_tier = "suspended"`) → usuarios no pueden hacer login → cliente paga → SuperAdmin reactiva asignando plan.

**Freelancer multi-empresa:** Empresa B contrata a Ana (que ya está en Empresa A) → SuperAdmin crea la membership (user → empresa B, rol recruiter) → Ana hace `switch-company` y trabaja en ambas.

**Plan custom Enterprise:** se negocia un plan privado (`is_public = false`, límites a medida) → se crea y se asigna a esa empresa; no aparece en la pricing page.

**Reporte para inversionistas:** analytics globales → MRR, churn, distribución de empresas por plan.

---

## 13. Caso de Uso End-to-End (resumen ejecutivo)

```
Día 1   Laura registra "TechStartup"        POST /auth/register-company
        → empresa free + Laura admin + token listo
Día 1   Laura crea y publica su primer job  POST /jobs → PATCH /jobs/10/publish
Día 2   Laura crea al recruiter Carlos      POST /users (role: recruiter)
Día 3   Carlos hace login y carga 10        POST /auth/login → POST /candidates × 10
        candidatos sourceados de LinkedIn
Día 4   Carlos crea las aplicaciones        POST /applications × 10  (stage: applied)
Día 5+  Mientras tanto, 5 candidatos más    POST /public/jobs/10/apply (career page)
        aplican solos desde la career page
Día 10  Screening: 5 avanzan                PATCH /applications/:id/move {screening}
Día 15  Pruebas técnicas: 2 avanzan         PATCH /applications/:id/move {technical}
Día 18  Entrevista final: 1 destaca         PATCH /applications/:id/move {offer} + rate 5★
Día 25  Candidata acepta la oferta          PATCH /applications/:id/move {hired} 🎉
Día 25  Laura cierra el job                 PATCH /jobs/10/close
        Dashboard registra: time-to-hire 24 días, conversion 6.7%
```

---

## 14. Garantías de Aislamiento que Sostienen Todos los Flujos

Reglas transversales que el usuario nunca ve pero siempre lo protegen (detalle técnico en doc 04 §9):

1. Toda lectura de listados filtra por el `company_id` **del token** (nunca del request).
2. Toda creación **fuerza** el `company_id` del token e ignora el del body.
3. Toda lectura/edición/borrado individual verifica que el recurso pertenezca a la empresa del token (si no: 403).
4. El switch de empresa emite un token nuevo: no hay forma de "mezclar" contextos.
5. Soft deletes en todo: el historial nunca se pierde.
