# Módulo SuperAdmin

## Descripción

El módulo SuperAdmin proporciona acceso global al sistema sin restricción de empresa (company_id). Permite gestionar todas las empresas, usuarios, planes y obtener analytics del sistema completo.

## Arquitectura

```
internal/app/
├── handlers/admin/
│   └── superadmin_companies_handler.go    # Endpoints HTTP
├── services/admin/
│   └── superadmin_companies_service.go    # Lógica de negocio
├── dtos/
│   └── superadmin_dto.go                  # DTOs específicos
└── middleware/
    └── auth_middleware.go                  # RequireSuperAdmin()
```

## Autenticación

El SuperAdmin debe cumplir:
1. ✅ Token JWT válido
2. ✅ Role = `"superadmin"`
3. ✅ **Sin `company_id`** en el token (acceso global)

```json
// JWT Claims para SuperAdmin
{
  "user_id": 1,
  "email": "superadmin@dvra.com",
  "role": "superadmin",
  "company_id": null  // ← Sin empresa (global access)
}
```

## Endpoints

### Base URL
```
/api/v1/admin/*
```

### 1. Listar Todas las Empresas
```http
GET /api/v1/admin/companies
Authorization: Bearer <superadmin_token>
```

**Query Params:**
- `page` (int): Página actual (default: 1)
- `limit` (int): Items por página (default: 20)
- `search` (string): Buscar por nombre o slug
- `plan_tier` (string): Filtrar por plan (`free`, `professional`, `enterprise`)

**Response:**
```json
{
  "companies": [
    {
      "id": 1,
      "name": "Azentic Sys",
      "slug": "azentic-sys",
      "plan_tier": "professional",
      "status": "active",
      "user_count": 5,
      "job_count": 12,
      "created_at": "2026-01-15T10:00:00Z",
      "trial_ends_at": "2026-02-15T10:00:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 45
  }
}
```

---

### 2. Crear Empresa con Admin
```http
POST /api/v1/admin/companies
Authorization: Bearer <superadmin_token>
Content-Type: application/json
```

**Body:**
```json
{
  "company_name": "TechCorp SA",
  "company_slug": "techcorp-sa",
  "admin_email": "admin@techcorp.com",
  "admin_password": "SecurePass123!",
  "admin_first_name": "Juan",
  "admin_last_name": "Pérez"
}
```

**Response (201):**
```json
{
  "company": {
    "id": 46,
    "name": "TechCorp SA",
    "slug": "techcorp-sa",
    "plan_tier": "professional",
    "trial_ends_at": "2026-02-08T12:00:00Z"
  },
  "admin": {
    "id": 120,
    "email": "admin@techcorp.com",
    "first_name": "Juan",
    "last_name": "Pérez"
  },
  "message": "Company and admin created successfully"
}
```

---

### 3. Cambiar Plan de Empresa
```http
PUT /api/v1/admin/companies/:id/plan
Authorization: Bearer <superadmin_token>
Content-Type: application/json
```

**Body:**
```json
{
  "new_plan": "enterprise"
}
```

**Valores válidos:** `free`, `professional`, `enterprise`

**Response (200):**
```json
{
  "message": "Plan updated successfully",
  "company_id": 46,
  "new_plan": "enterprise"
}
```

---

### 4. Suspender Empresa
```http
POST /api/v1/admin/companies/:id/suspend
Authorization: Bearer <superadmin_token>
Content-Type: application/json
```

**Body:**
```json
{
  "reason": "Falta de pago - 3 meses vencidos"
}
```

**Response (200):**
```json
{
  "message": "Company suspended successfully",
  "company_id": 46,
  "reason": "Falta de pago - 3 meses vencidos"
}
```

**Efecto:** 
- `plan_tier` se actualiza a `"suspended"`
- Los usuarios de esa empresa no podrán acceder al sistema

---

### 5. Ver Usuarios de una Empresa
```http
GET /api/v1/admin/companies/:id/users
Authorization: Bearer <superadmin_token>
```

**Response (200):**
```json
{
  "company_id": 1,
  "users": [
    {
      "id": 5,
      "email": "admin@azentic.com",
      "first_name": "Admin",
      "last_name": "Azentic",
      "is_active": true
    },
    {
      "id": 8,
      "email": "recruiter@azentic.com",
      "first_name": "María",
      "last_name": "González",
      "is_active": true
    }
  ],
  "count": 2
}
```

---

### 6. Analytics Globales
```http
GET /api/v1/admin/analytics
Authorization: Bearer <superadmin_token>
```

**Response (200):**
```json
{
  "total_companies": 45,
  "active_companies": 42,
  "suspended_companies": 3,
  "total_users": 230,
  "total_jobs": 156,
  "total_applications": 1842,
  "monthly_revenue": 6253.00,
  "churn_rate": 0.05
}
```

---

## Diferencias: Admin vs SuperAdmin

| Característica | Admin Regular | SuperAdmin |
|----------------|---------------|------------|
| **Company Scope** | ✅ Tiene `company_id` | ❌ Sin `company_id` (global) |
| **Ver empresas** | Solo la suya | Todas |
| **Crear empresas** | No | ✅ Sí |
| **Cambiar planes** | No | ✅ Sí |
| **Suspender empresas** | No | ✅ Sí |
| **Ver analytics globales** | No | ✅ Sí |
| **Gestionar usuarios** | Solo de su empresa | De cualquier empresa |
| **Crear jobs** | ✅ De su empresa | N/A (no tiene empresa) |
| **Ver candidatos** | ✅ De su empresa | N/A |

---

## Flujo de Uso

### 1. Login como SuperAdmin
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "superadmin@dvra.com",
    "password": "SuperAdmin123!"
  }'
```

**Response:**
```json
{
  "access_token": "eyJhbGc...",
  "refresh_token": "eyJhbGc...",
  "user": {
    "id": 1,
    "email": "superadmin@dvra.com",
    "first_name": "Super",
    "last_name": "Admin"
  }
}
```

---

### 2. Listar Empresas Activas
```bash
curl http://localhost:8080/api/v1/admin/companies?plan_tier=professional \
  -H "Authorization: Bearer eyJhbGc..."
```

---

### 3. Crear Nueva Empresa
```bash
curl -X POST http://localhost:8080/api/v1/admin/companies \
  -H "Authorization: Bearer eyJhbGc..." \
  -H "Content-Type: application/json" \
  -d '{
    "company_name": "InnovateLabs",
    "company_slug": "innovate-labs",
    "admin_email": "ceo@innovatelabs.io",
    "admin_password": "Welcome2024!",
    "admin_first_name": "Carlos",
    "admin_last_name": "Ramírez"
  }'
```

---

### 4. Upgrade Plan
```bash
curl -X PUT http://localhost:8080/api/v1/admin/companies/5/plan \
  -H "Authorization: Bearer eyJhbGc..." \
  -H "Content-Type: application/json" \
  -d '{"new_plan": "enterprise"}'
```

---

### 5. Suspender Empresa Morosa
```bash
curl -X POST http://localhost:8080/api/v1/admin/companies/12/suspend \
  -H "Authorization: Bearer eyJhbGc..." \
  -H "Content-Type: application/json" \
  -d '{"reason": "Impago - 90 días vencidos"}'
```

---

### 6. Ver Analytics del Sistema
```bash
curl http://localhost:8080/api/v1/admin/analytics \
  -H "Authorization: Bearer eyJhbGc..."
```

---

## Seguridad

### Middleware de Validación
```go
// RequireSuperAdmin() valida:
1. Token JWT válido
2. Role = "superadmin"
3. company_id = nil (sin contexto de empresa)
```

### Credenciales por Defecto
```
Email: superadmin@dvra.com
Password: SuperAdmin123!
```

⚠️ **IMPORTANTE:** Cambiar contraseña en producción.

---

## Casos de Uso

### 1. Onboarding de Nueva Empresa
1. Cliente firma contrato
2. SuperAdmin crea empresa con admin inicial
3. Admin recibe credenciales por email
4. Admin ingresa y configura su empresa

### 2. Gestión de Morosidad
1. Sistema detecta pago vencido
2. SuperAdmin revisa y suspende empresa
3. Usuarios no pueden acceder
4. Cliente paga
5. SuperAdmin reactiva empresa

### 3. Reportes Ejecutivos
1. SuperAdmin accede a analytics
2. Genera reporte mensual:
   - MRR (Monthly Recurring Revenue)
   - Churn rate
   - Empresas activas vs suspendidas
   - Crecimiento de usuarios

### 4. Soporte Técnico
1. Cliente reporta problema
2. SuperAdmin accede a `/admin/companies/:id/users`
3. Verifica configuración de la empresa
4. Resuelve issue sin necesidad de login del cliente

---

## Mejoras Futuras

### 1. Impersonation
```http
POST /api/v1/admin/impersonate/:user_id
```
Permitir login como cualquier usuario para debugging.

### 2. Audit Logs
Registrar todas las acciones del SuperAdmin:
```json
{
  "action": "suspend_company",
  "superadmin_id": 1,
  "target_company_id": 12,
  "reason": "Impago",
  "timestamp": "2026-01-08T14:30:00Z"
}
```

### 3. Reportes Avanzados
- Churn rate real (tracking de bajas)
- Cohort analysis
- Exportar reportes a CSV/PDF

### 4. Notificaciones
Alertas automáticas cuando:
- Nueva empresa registrada
- Empresa suspendida
- Plan upgrade/downgrade
- Empresa alcanza límite de trial

---

## Testing

### Test de Acceso Denegado
```bash
# Intentar acceso con token de Admin regular (debe fallar)
curl http://localhost:8080/api/v1/admin/companies \
  -H "Authorization: Bearer <admin_regular_token>"

# Response esperado: 403 Forbidden
{
  "error": "SuperAdmin access required"
}
```

### Test de Token sin Empresa
```bash
# Token con company_id presente (debe fallar)
curl http://localhost:8080/api/v1/admin/companies \
  -H "Authorization: Bearer <token_with_company_id>"

# Response esperado: 403 Forbidden
{
  "error": "SuperAdmin routes require no company context"
}
```

---

## Contacto

Si tienes dudas sobre el módulo SuperAdmin, contacta al equipo de desarrollo.
