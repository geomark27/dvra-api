# API Endpoints - Dvra ATS

> **Versión:** v1.3.0  
> **Base URL:** `http://localhost:8001/api/v1`  
> **Última Actualización:** 21 de Diciembre, 2025

---

## 📋 Índice

- [Autenticación](#autenticación)
- [Health Check](#health-check)
- [Usuarios](#usuarios)
- [Empresas](#empresas)
- [Memberships](#memberships)
- [Jobs (Ofertas de Trabajo)](#jobs-ofertas-de-trabajo)
- [Candidatos](#candidatos)
- [Aplicaciones](#aplicaciones)
- [System Values (Catálogos)](#system-values-catálogos)
- [SuperAdmin](#superadmin)

---

## 🔐 Autenticación

Todos los endpoints protegidos requieren un token JWT en el header:

```http
Authorization: Bearer <access_token>
```

### Tipos de Acceso

| Tipo | Company ID | Descripción |
|------|-----------|-------------|
| **Público** | - | Sin autenticación requerida |
| **Protegido** | ✅ Requerido | Usuario debe tener company_id |
| **SuperAdmin** | ❌ Sin company | Acceso global sin empresa |

---

## 🔑 Autenticación

### Registro de Empresa (Recomendado)

```http
POST /api/v1/auth/register-company
```

**Acceso:** Público

**Descripción:** Crea una nueva empresa con su primer usuario administrador. Este es el flujo principal de registro. La empresa se crea con el plan "free" por defecto, validando que el plan exista y esté activo en la base de datos.

**Request Body:**
```json
{
  "company_name": "Azentic Sys",
  "company_slug": "azentic-sys",
  "admin_email": "admin@azentic.com",
  "admin_password": "Admin123!",
  "admin_first_name": "Marcos",
  "admin_last_name": "Ramos",
  "timezone": "America/Bogota"
}
```

**Response (201):**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "company": {
    "id": 1,
    "name": "Azentic Sys",
    "slug": "azentic-sys",
    "plan_tier": "free"
  },
  "admin": {
    "id": 1,
    "email": "admin@azentic.com",
    "first_name": "Marcos",
    "last_name": "Ramos",
    "is_active": true
  }
}
```

**Errores:**
- `409` - Email ya existe

**Nota:** El token generado ya incluye `company_id`, por lo que el usuario puede acceder inmediatamente a todas las funcionalidades.

---

### Registro de Usuario (Deprecated)

```http
POST /api/v1/auth/register
```

**Acceso:** Público

**⚠️ DEPRECADO:** Este endpoint crea un usuario sin empresa. Usa `/auth/register-company` en su lugar.

**Request Body:**
```json
{
  "email": "usuario@ejemplo.com",
  "password": "Password123!",
  "first_name": "Juan",
  "last_name": "Pérez"
}
```

**Response (201):**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "email": "usuario@ejemplo.com",
    "first_name": "Juan",
    "last_name": "Pérez",
    "is_active": true
  }
}
```

**Limitación:** El token no tiene `company_id`, por lo que el usuario no podrá acceder a rutas protegidas hasta que sea agregado a una empresa.

---

### Login de Usuario

```http
POST /api/v1/auth/login
```

**Acceso:** Público

**Descripción:** Autentica usuarios normales (admin, users, etc.) vinculados a empresas y retorna tokens junto con la lista de empresas a las que pertenece.

**Request Body:**
```json
{
  "email": "admin@azentic.com",
  "password": "Admin123!"
}
```

**Response (200):**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "email": "admin@azentic.com",
    "first_name": "Marcos",
    "last_name": "Ramos",
    "is_active": true
  },
  "companies": [
    {
      "id": 1,
      "name": "Azentic Sys",
      "slug": "azentic-sys",
      "plan_tier": "professional"
    },
    {
      "id": 2,
      "name": "DevCorp",
      "slug": "devcorp",
      "plan_tier": "enterprise"
    }
  ]
}
```

**Errores:**
- `401` - Email o contraseña inválidos
- `401` - Cuenta inactiva

**Nota:** El token generado usa la empresa por defecto (primera membership con `is_default=true`). Si el usuario pertenece a múltiples empresas, puede cambiar el contexto con `/auth/switch-company`.

---

### Login de SuperAdmin

```http
POST /api/v1/auth/superadmin/login
```

**Acceso:** Público

**Descripción:** Autentica usuarios SuperAdmin. Este es un endpoint separado exclusivo para SuperAdmins que genera tokens sin `company_id`.

**Request Body:**
```json
{
  "email": "superadmin@dvra.com",
  "password": "SuperAdmin123!"
}
```

**Response (200):**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "email": "superadmin@dvra.com",
    "first_name": "Super",
    "last_name": "Admin",
    "is_active": true
  },
  "is_superadmin": true
}
```

**Errores:**
- `401` - Email o contraseña inválidos
- `403` - Usuario no es SuperAdmin
- `401` - Cuenta inactiva

**Diferencias con login normal:**
- ✅ Valida que el usuario tenga `is_superadmin = true`
- ✅ Genera token sin `company_id` (acceso global)
- ✅ Rol en token: `"superadmin"`
- ❌ No retorna lista de empresas

---

### Refresh Token

```http
POST /api/v1/auth/refresh
```

**Acceso:** Público

**Request Body:**
```json
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Response (200):**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

---

### Obtener Usuario Actual

```http
GET /api/v1/auth/me
```

**Acceso:** Protegido  
**Headers:** `Authorization: Bearer <token>`

**Response (200):**
```json
{
  "id": 1,
  "email": "usuario@ejemplo.com",
  "first_name": "Juan",
  "last_name": "Pérez",
  "is_active": true
}
```

---

### Cambiar Contraseña

```http
POST /api/v1/auth/change-password
```

**Acceso:** Protegido  
**Headers:** `Authorization: Bearer <token>`

**Request Body:**
```json
{
  "old_password": "Password123!",
  "new_password": "NewPassword456!"
}
```

**Response (200):**
```json
{
  "message": "Password changed successfully"
}
```

**Errores:**
- `401` - Contraseña antigua incorrecta

---

### Logout

```http
POST /api/v1/auth/logout
```

**Acceso:** Protegido  
**Headers:** `Authorization: Bearer <token>`

**Response (200):**
```json
{
  "message": "Logged out successfully"
}
```

**Nota:** JWT es stateless. El cliente debe descartar los tokens.

---

### Cambiar de Empresa (Switch Company)

```http
POST /api/v1/auth/switch-company
```

**Acceso:** Protegido  
**Headers:** `Authorization: Bearer <token>`

**Descripción:** Cambia el contexto del usuario a otra empresa. Genera un nuevo token con el `company_id` de la empresa seleccionada.

**Request Body:**
```json
{
  "company_id": 2
}
```

**Response (200):**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "company": {
    "id": 2,
    "name": "DevCorp",
    "slug": "devcorp",
    "plan_tier": "enterprise"
  }
}
```

**Errores:**
- `403` - El usuario no pertenece a esa empresa
- `404` - Empresa no encontrada

**Uso:**
```bash
# 1. Login (recibe empresa por defecto)
POST /auth/login

# 2. Cambiar a otra empresa
POST /auth/switch-company "company_id": 2 }

# 3. Usar nuevo token para acceder a datos de la empresa 2
GET /jobs
Authorization: Bearer <nuevo_token>
```

---

### Obtener Mis Empresas

```http
GET /api/v1/auth/my-companies
```

**Acceso:** Protegido  
**Headers:** `Authorization: Bearer <token>`

**Descripción:** Retorna la lista de todas las empresas a las que pertenece el usuario autenticado.

**Response (200):**
```json
{
  "companies": [
    {
      "id": 1,
      "name": "Azentic Sys",
      "slug": "azentic-sys",
      "plan_tier": "professional"
    },
    {
      "id": 2,
      "name": "DevCorp",
      "slug": "devcorp",
      "plan_tier": "enterprise"
    }
  ]
}
```

**Caso de Uso:** Útil para mostrar un selector de empresas en el frontend cuando el usuario tiene múltiples memberships.

---

## 🏥 Health Check

### Verificar Estado del Servidor

```http
GET /api/v1/health
```

**Acceso:** Público

**Response (200):**
```json
{
  "status": "ok",
  "timestamp": "2025-12-08T15:30:00Z"
}
```

---

### Verificar Disponibilidad

```http
GET /api/v1/health/ready
```

**Acceso:** Público

**Response (200):**
```json
{
  "status": "ready",
  "database": "connected"
}
```

---

## 👥 Usuarios

### Listar Usuarios

```http
GET /api/v1/users
```

**Acceso:** Protegido (Company-scoped)  
**Headers:** `Authorization: Bearer <token>`

**Response (200):**
```json
[
  {
    "id": 1,
    "first_name": "Juan",
    "last_name": "Pérez",
    "email": "juan@empresa.com"
  }
]
```

---

### Obtener Usuario por ID

```http
GET /api/v1/users/:id
```

**Acceso:** Protegido (Company-scoped)  
**Headers:** `Authorization: Bearer <token>`

**Response (200):**
```json
{
  "id": 1,
  "first_name": "Juan",
  "last_name": "Pérez",
  "email": "juan@empresa.com",
  "is_active": true
}
```

---

### Crear Usuario

```http
POST /api/v1/users
```

**Acceso:** Protegido (Company-scoped)  
**Headers:** `Authorization: Bearer <token>`

**Request Body:**
```json
{
  "first_name": "María",
  "last_name": "González",
  "email": "maria@empresa.com"
}
```

**Response (201):**
```json
{
  "id": 2,
  "first_name": "María",
  "last_name": "González",
  "email": "maria@empresa.com"
}
```

---

### Actualizar Usuario

```http
PUT /api/v1/users/:id
```

**Acceso:** Protegido (Company-scoped)  
**Headers:** `Authorization: Bearer <token>`

**Request Body:**
```json
{
  "first_name": "María Fernanda",
  "last_name": "González",
  "email": "mariaf@empresa.com"
}
```

**Response (200):**
```json
{
  "id": 2,
  "first_name": "María Fernanda",
  "last_name": "González",
  "email": "mariaf@empresa.com"
}
```

---

### Eliminar Usuario

```http
DELETE /api/v1/users/:id
```

**Acceso:** Protegido (Company-scoped)  
**Headers:** `Authorization: Bearer <token>`

**Response (200):**
```json
{
  "message": "User deleted successfully"
}
```

---

## 🏢 Empresas

### Listar Empresas

```http
GET /api/v1/companies
```

**Acceso:** Protegido (Company-scoped)  
**Headers:** `Authorization: Bearer <token>`

**Response (200):**
```json
[
  {
    "id": 1,
    "name": "Azentic Sys",
    "slug": "azentic-sys",
    "plan_tier": "professional"
  }
]
```

---

### Obtener Empresa por ID

```http
GET /api/v1/companies/:id
```

**Acceso:** Protegido (Company-scoped)  
**Headers:** `Authorization: Bearer <token>`

**Response (200):**
```json
{
  "id": 1,
  "name": "Azentic Sys",
  "slug": "azentic-sys",
  "logo_url": "",
  "plan_tier": "professional",
  "timezone": "America/Bogota"
}
```

---

### Crear Empresa

```http
POST /api/v1/companies
```

**Acceso:** Protegido  
**Headers:** `Authorization: Bearer <token>`

**Request Body:**
```json
{
  "name": "TechCorp",
  "slug": "techcorp",
  "timezone": "America/Bogota"
}
```

**Response (201):**
```json
{
  "id": 2,
  "name": "TechCorp",
  "slug": "techcorp",
  "plan_tier": "free"
}
```

---

### Actualizar Empresa

```http
PUT /api/v1/companies/:id
```

**Acceso:** Protegido (Company-scoped)  
**Headers:** `Authorization: Bearer <token>`

**Request Body:**
```json
{
  "name": "TechCorp International",
  "logo_url": "https://example.com/logo.png"
}
```

**Response (200):**
```json
{
  "id": 2,
  "name": "TechCorp International",
  "slug": "techcorp",
  "logo_url": "https://example.com/logo.png"
}
```

---

### Eliminar Empresa

```http
DELETE /api/v1/companies/:id
```

**Acceso:** Protegido (Admin only)  
**Headers:** `Authorization: Bearer <token>`

**Response (200):**
```json
{
  "message": "Company deleted successfully"
}
```

---

## 👨‍💼 Memberships

### Listar Memberships

```http
GET /api/v1/memberships
```

**Acceso:** Protegido (Company-scoped)  
**Headers:** `Authorization: Bearer <token>`

**Response (200):**
```json
[
  {
    "id": 1,
    "user_id": 1,
    "company_id": 1,
    "role": "admin",
    "status": "active"
  }
]
```

---

### Obtener Membership por ID

```http
GET /api/v1/memberships/:id
```

**Acceso:** Protegido (Company-scoped)  
**Headers:** `Authorization: Bearer <token>`

**Response (200):**
```json
{
  "id": 1,
  "user_id": 1,
  "company_id": 1,
  "role": "admin",
  "status": "active",
  "is_default": true
}
```

---

> **NOTA MVP:** La creación de memberships (asignar usuarios a empresas) está **restringida a SuperAdmin** únicamente. Los administradores de empresa solo pueden ver, actualizar roles y eliminar memberships de su propia empresa. Ver [sección SuperAdmin](#crear-membership-asignar-usuario-a-empresa) para crear memberships.

---

### Actualizar Membership

```http
PUT /api/v1/memberships/:id
```

**Acceso:** Protegido (Admin only)  
**Headers:** `Authorization: Bearer <token>`

**Request Body:**
```json
{
  "role": "admin",
  "status": "active"
}
```

**Response (200):**
```json
{
  "id": 2,
  "role": "admin",
  "status": "active"
}
```

---

### Eliminar Membership

```http
DELETE /api/v1/memberships/:id
```

**Acceso:** Protegido (Admin only)  
**Headers:** `Authorization: Bearer <token>`

**Response (200):**
```json
{
  "message": "Membership deleted successfully"
}
```

---

## 💼 Jobs (Ofertas de Trabajo)

### Listar Jobs

```http
GET /api/v1/jobs
```

**Acceso:** Protegido (Company-scoped)  
**Headers:** `Authorization: Bearer <token>`

**Response (200):**
```json
[
  {
    "id": 1,
    "company_id": 1,
    "title": "Backend Developer",
    "status": "published"
  }
]
```

---

### Obtener Job por ID

```http
GET /api/v1/jobs/:id
```

**Acceso:** Protegido (Company-scoped)  
**Headers:** `Authorization: Bearer <token>`

**Response (200):**
```json
{
  "id": 1,
  "company_id": 1,
  "title": "Backend Developer",
  "description": "Buscamos desarrollador con experiencia en Go...",
  "status": "published",
  "remote_type": "hybrid"
}
```

---

### Crear Job

```http
POST /api/v1/jobs
```

**Acceso:** Protegido (Company-scoped)  
**Headers:** `Authorization: Bearer <token>`

**Request Body:**
```json
{
  "title": "Frontend Developer",
  "description": "Experiencia en React y TypeScript...",
  "remote_type": "remote",
  "status": "draft"
}
```

**Response (201):**
```json
{
  "id": 2,
  "company_id": 1,
  "title": "Frontend Developer",
  "status": "draft"
}
```

---

### Actualizar Job

```http
PUT /api/v1/jobs/:id
```

**Acceso:** Protegido (Company-scoped)  
**Headers:** `Authorization: Bearer <token>`

**Request Body:**
```json
{
  "title": "Senior Frontend Developer",
  "status": "published"
}
```

**Response (200):**
```json
{
  "id": 2,
  "title": "Senior Frontend Developer",
  "status": "published"
}
```

---

### Eliminar Job

```http
DELETE /api/v1/jobs/:id
```

**Acceso:** Protegido (Company-scoped)  
**Headers:** `Authorization: Bearer <token>`

**Response (200):**
```json
{
  "message": "Job deleted successfully"
}
```

---

## 👤 Candidatos

### Listar Candidatos

```http
GET /api/v1/candidates
```

**Acceso:** Protegido (Company-scoped)  
**Headers:** `Authorization: Bearer <token>`

**Response (200):**
```json
[
  {
    "id": 1,
    "email": "candidato@example.com",
    "first_name": "Carlos",
    "last_name": "Ramírez"
  }
]
```

---

### Obtener Candidato por ID

```http
GET /api/v1/candidates/:id
```

**Acceso:** Protegido (Company-scoped)  
**Headers:** `Authorization: Bearer <token>`

**Response (200):**
```json
{
  "id": 1,
  "email": "candidato@example.com",
  "first_name": "Carlos",
  "last_name": "Ramírez",
  "phone": "+57 300 1234567"
}
```

---

### Crear Candidato

```http
POST /api/v1/candidates
```

**Acceso:** Protegido (Company-scoped)  
**Headers:** `Authorization: Bearer <token>`

**Request Body:**
```json
{
  "email": "nuevo@example.com",
  "first_name": "Ana",
  "last_name": "Martínez",
  "phone": "+57 300 9876543"
}
```

**Response (201):**
```json
{
  "id": 2,
  "email": "nuevo@example.com",
  "first_name": "Ana",
  "last_name": "Martínez"
}
```

---

### Actualizar Candidato

```http
PUT /api/v1/candidates/:id
```

**Acceso:** Protegido (Company-scoped)  
**Headers:** `Authorization: Bearer <token>`

**Request Body:**
```json
{
  "phone": "+57 300 1111111",
  "linkedin_url": "https://linkedin.com/in/anamartinez"
}
```

**Response (200):**
```json
{
  "id": 2,
  "phone": "+57 300 1111111",
  "linkedin_url": "https://linkedin.com/in/anamartinez"
}
```

---

### Eliminar Candidato

```http
DELETE /api/v1/candidates/:id
```

**Acceso:** Protegido (Company-scoped)  
**Headers:** `Authorization: Bearer <token>`

**Response (200):**
```json
{
  "message": "Candidate deleted successfully"
}
```

---

## 📝 Aplicaciones

### Listar Aplicaciones

```http
GET /api/v1/applications
```

**Acceso:** Protegido (Company-scoped)  
**Headers:** `Authorization: Bearer <token>`

**Response (200):**
```json
[
  {
    "id": 1,
    "job_id": 1,
    "candidate_id": 1,
    "status": "pending"
  }
]
```

---

### Obtener Aplicación por ID

```http
GET /api/v1/applications/:id
```

**Acceso:** Protegido (Company-scoped)  
**Headers:** `Authorization: Bearer <token>`

**Response (200):**
```json
{
  "id": 1,
  "job_id": 1,
  "candidate_id": 1,
  "status": "pending",
  "applied_at": "2025-12-08T10:00:00Z"
}
```

---

### Crear Aplicación

```http
POST /api/v1/applications
```

**Acceso:** Protegido (Company-scoped)  
**Headers:** `Authorization: Bearer <token>`

**Request Body:**
```json
{
  "job_id": 1,
  "candidate_id": 2
}
```

**Response (201):**
```json
{
  "id": 2,
  "job_id": 1,
  "candidate_id": 2,
  "status": "pending"
}
```

---

### Actualizar Aplicación

```http
PUT /api/v1/applications/:id
```

**Acceso:** Protegido (Company-scoped)  
**Headers:** `Authorization: Bearer <token>`

**Request Body:**
```json
{
  "status": "reviewing"
}
```

**Response (200):**
```json
{
  "id": 2,
  "status": "reviewing"
}
```

**Estados disponibles:**
- `pending` - Pendiente de revisión
- `reviewing` - En revisión
- `interviewed` - Entrevistado
- `rejected` - Rechazado
- `hired` - Contratado

---

### Eliminar Aplicación

```http
DELETE /api/v1/applications/:id
```

**Acceso:** Protegido (Company-scoped)  
**Headers:** `Authorization: Bearer <token>`

**Response (200):**
```json
{
  "message": "Application deleted successfully"
}
```

---

## 🗂️ System Values (Catálogos)

Los **System Values** son valores del sistema configurables que alimentan los selects y dropdowns del frontend. Reemplazan valores hardcodeados permitiendo flexibilidad y personalización por empresa.

### Categorías Disponibles

Después del seed inicial, se incluyen las siguientes categorías:

| Categoría | Descripción | Valores |
|-----------|-------------|---------|
| `job_status` | Estados de trabajos | draft, published, closed |
| `application_status` | Estados de aplicaciones | applied, screening, technical, interview, offer, hired, rejected |
| `contract_type` | Tipos de contrato | full_time, part_time, contractor, internship, temporary |
| `work_mode` | Modalidad de trabajo | remote, onsite, hybrid |
| `experience_level` | Nivel de experiencia | entry, mid, senior, lead |
| `priority` | Prioridades | low, medium, high, urgent |
| `candidate_source` | Fuente de candidatos | linkedin, website, referral, job_board, direct, other |

### Obtener Valores por Categoría

```http
GET /api/v1/system-values/:category
```

**Acceso:** Público  
**Headers:** `X-Company-ID` (opcional) - Si se envía, incluye valores globales + específicos de esa empresa

**Descripción:** Retorna todos los valores activos para una categoría específica. Los valores globales (`company_id = null`) están disponibles para todos. Si se envía `X-Company-ID`, también incluye valores personalizados de esa empresa.

**Ejemplo:**
```bash
GET /api/v1/system-values/job_status
```

**Response (200):**
```json
{
  "status": "success",
  "data": [
    {
      "id": 1,
      "category": "job_status",
      "value": "draft",
      "label": "Borrador (no visible para candidatos)",
      "description": null,
      "display_order": 1,
      "is_active": true,
      "company_id": null
    },
    {
      "id": 2,
      "category": "job_status",
      "value": "published",
      "label": "Publicada (visible para candidatos)",
      "description": null,
      "display_order": 2,
      "is_active": true,
      "company_id": null
    },
    {
      "id": 3,
      "category": "job_status",
      "value": "closed",
      "label": "Cerrada",
      "description": null,
      "display_order": 3,
      "is_active": true,
      "company_id": null
    }
  ]
}
```

**Uso en Frontend:**
```typescript
// lib/hooks/useSystemValues.ts
export const useSystemValues = (category: string) => {
  return useQuery({
    queryKey: ['system-values', category],
    queryFn: () => api.get(`/system-values/${category}`),
  });
};

// En componente:
const { data: jobStatuses } = useSystemValues('job_status');

<Select>
  {jobStatuses?.data.map(status => (
    <SelectItem key={status.value} value={status.value}>
      {status.label}
    </SelectItem>
  ))}
</Select>
```

---

### Ejemplos de Categorías

#### Estados de Trabajo (job_status)
```bash
GET /api/v1/system-values/job_status
```

#### Estados de Aplicación (application_status)
```bash
GET /api/v1/system-values/application_status
```

#### Tipos de Contrato (contract_type)
```bash
GET /api/v1/system-values/contract_type
```

#### Modalidad de Trabajo (work_mode)
```bash
GET /api/v1/system-values/work_mode
```

#### Nivel de Experiencia (experience_level)
```bash
GET /api/v1/system-values/experience_level
```

#### Prioridades (priority)
```bash
GET /api/v1/system-values/priority
```

#### Fuente de Candidatos (candidate_source)
```bash
GET /api/v1/system-values/candidate_source
```

---

### Eliminar Aplicación

```http
DELETE /api/v1/applications/:id
```

**Acceso:** Protegido (Company-scoped)  
**Headers:** `Authorization: Bearer <token>`

**Response (200):**
```json
{
  "message": "Application deleted successfully"
}
```

---

## 🔒 SuperAdmin

> **IMPORTANTE:** Estos endpoints solo están disponibles para usuarios con role `superadmin` y **sin company_id** (acceso global).

### Listar Todas las Empresas

```http
GET /api/v1/admin/companies
```

**Acceso:** SuperAdmin  
**Headers:** `Authorization: Bearer <superadmin_token>`

**Query Params:**
- `page` (int) - Página actual (default: 1)
- `limit` (int) - Items por página (default: 20)
- `search` (string) - Buscar por nombre o slug
- `plan_tier` (string) - Filtrar por plan (`free`, `professional`, `enterprise`)

**Response (200):**
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

### Crear Empresa con Admin

```http
POST /api/v1/admin/companies
```

**Acceso:** SuperAdmin  
**Headers:** `Authorization: Bearer <superadmin_token>`

**Request Body:**
```json
{
  "company_name": "TechCorp SA",
  "company_slug": "techcorp-sa",
  "admin_email": "admin@techcorp.com",
  "admin_password": "SecurePass123!",
  "admin_first_name": "Juan",
  "admin_last_name": "Pérez",
  "plan_slug": "starter"  // Opcional: free, starter, professional, enterprise (default: free)
}
```

**Validación:** El sistema valida que el plan elegido exista en la base de datos y esté activo usando `FindActiveBySlug()`. Si no se especifica plan_slug, se asigna "free" por defecto.

**Response (201):**
```json
{
  "company": {
    "id": 46,
    "name": "TechCorp SA",
    "slug": "techcorp-sa",
    "plan_tier": "starter"
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

### Cambiar Plan de Empresa

```http
PUT /api/v1/admin/companies/:id/plan
```

**Acceso:** SuperAdmin  
**Headers:** `Authorization: Bearer <superadmin_token>`

**Request Body:**
```json
{
  "new_plan": "enterprise"
}
```

**Planes válidos:** `free`, `professional`, `enterprise`

**Response (200):**
```json
{
  "message": "Plan updated successfully",
  "company_id": 46,
  "new_plan": "enterprise"
}
```

---

### Suspender Empresa

```http
POST /api/v1/admin/companies/:id/suspend
```

**Acceso:** SuperAdmin  
**Headers:** `Authorization: Bearer <superadmin_token>`

**Request Body:**
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

---

### Ver Usuarios de Empresa

```http
GET /api/v1/admin/companies/:id/users
```

**Acceso:** SuperAdmin  
**Headers:** `Authorization: Bearer <superadmin_token>`

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
    }
  ],
  "count": 5
}
```

---

### Crear Membership (Asignar Usuario a Empresa)

```http
POST /api/v1/admin/memberships
```

**Acceso:** SuperAdmin ÚNICAMENTE  
**Headers:** `Authorization: Bearer <superadmin_token>`

**Descripción:** Asigna un usuario existente a una empresa. Solo SuperAdmin puede realizar esta operación para evitar manipulación cross-company.

**Request Body:**
```json
{
  "user_id": 2,
  "company_id": 1,
  "role": "recruiter"
}
```

**Response (201):**
```json
{
  "id": 5,
  "user_id": 2,
  "company_id": 1,
  "role": "recruiter",
  "status": "active"
}
```

**Roles disponibles:**
- `superadmin` (100) - Acceso global
- `admin` (50) - Administrador de empresa
- `recruiter` (30) - Reclutador
- `hiring_manager` (20) - Gerente de contratación
- `user` (10) - Usuario básico

**Errores:**
- `403` - Usuario no es SuperAdmin
- `400` - company_id es requerido
- `404` - Usuario o empresa no encontrados

**Nota:** Los clientes (admins de empresa) deben crear nuevos usuarios en lugar de asignar usuarios existentes. Sistema de invitación por email vendrá en Fase 2.

---

### Gestión de System Values (SuperAdmin)

> **IMPORTANTE:** Solo SuperAdmin puede crear, editar o eliminar System Values. Los usuarios normales solo pueden consultarlos.

#### Listar Todos los System Values

```http
GET /api/v1/admin/system-values
```

**Acceso:** SuperAdmin  
**Headers:** `Authorization: Bearer <superadmin_token>`

**Descripción:** Retorna todos los system values del sistema (globales y específicos de empresas).

**Response (200):**
```json
{
  "status": "success",
  "data": [
    {
      "id": 1,
      "category": "job_status",
      "value": "draft",
      "label": "Borrador (no visible para candidatos)",
      "description": null,
      "display_order": 1,
      "is_active": true,
      "company_id": null,
      "created_at": "2025-12-21T10:00:00Z",
      "updated_at": "2025-12-21T10:00:00Z"
    }
  ]
}
```

---

#### Crear System Value

```http
POST /api/v1/admin/system-values
```

**Acceso:** SuperAdmin  
**Headers:** `Authorization: Bearer <superadmin_token>`

**Descripción:** Crea un nuevo valor del sistema. Puede ser global (`company_id: null`) o específico de una empresa.

**Request Body:**
```json
{
  "category": "priority",
  "value": "critical",
  "label": "Crítico",
  "description": "Para casos extremadamente urgentes",
  "display_order": 5,
  "company_id": null
}
```

**Campos:**
- `category` (string, required) - Categoría del valor (job_status, priority, etc.)
- `value` (string, required) - Valor técnico (usado en código)
- `label` (string, required) - Etiqueta mostrada al usuario
- `description` (string, optional) - Descripción adicional
- `display_order` (int, optional) - Orden de visualización (default: 0)
- `company_id` (int, optional) - NULL para global, ID para específico de empresa

**Response (201):**
```json
{
  "status": "success",
  "message": "System value created successfully",
  "data": {
    "id": 35,
    "category": "priority",
    "value": "critical",
    "label": "Crítico",
    "description": "Para casos extremadamente urgentes",
    "display_order": 5,
    "is_active": true,
    "company_id": null
  }
}
```

**Errores:**
- `400` - Datos inválidos
- `409` - System value already exists for this category and company

---

#### Actualizar System Value

```http
PUT /api/v1/admin/system-values/:id
```

**Acceso:** SuperAdmin  
**Headers:** `Authorization: Bearer <superadmin_token>`

**Descripción:** Actualiza un valor del sistema existente. No se puede cambiar `category` ni `value`.

**Request Body:**
```json
{
  "label": "Crítico - Máxima Prioridad",
  "description": "Urgencia máxima, atender inmediatamente",
  "display_order": 10,
  "is_active": true
}
```

**Response (200):**
```json
{
  "status": "success",
  "message": "System value updated successfully",
  "data": {
    "id": 35,
    "category": "priority",
    "value": "critical",
    "label": "Crítico - Máxima Prioridad",
    "description": "Urgencia máxima, atender inmediatamente",
    "display_order": 10,
    "is_active": true
  }
}
```

**Errores:**
- `404` - System value not found

---

#### Eliminar System Value

```http
DELETE /api/v1/admin/system-values/:id
```

**Acceso:** SuperAdmin  
**Headers:** `Authorization: Bearer <superadmin_token>`

**Descripción:** Elimina (soft delete) un valor del sistema. Recomendado desactivar (`is_active: false`) en lugar de eliminar.

**Response (200):**
```json
{
  "status": "success",
  "message": "System value deleted successfully"
}
```

**Errores:**
- `404` - System value not found

**Nota:** Usar soft delete permite mantener historial. Considera desactivar en lugar de eliminar.

---

### Analytics Globales

```http
GET /api/v1/admin/analytics
```

**Acceso:** SuperAdmin  
**Headers:** `Authorization: Bearer <superadmin_token>`

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

## 🔧 Códigos de Estado HTTP

| Código | Descripción |
|--------|-------------|
| `200` | OK - Solicitud exitosa |
| `201` | Created - Recurso creado exitosamente |
| `400` | Bad Request - Datos inválidos |
| `401` | Unauthorized - Token inválido o expirado |
| `403` | Forbidden - Sin permisos suficientes |
| `404` | Not Found - Recurso no encontrado |
| `409` | Conflict - Recurso ya existe |
| `500` | Internal Server Error - Error del servidor |

---

## 📝 Notas

### Tokens JWT

Los tokens tienen una duración de:
- **Access Token:** 1 hora
- **Refresh Token:** 30 días

### Paginación

Los endpoints que soportan paginación aceptan:
- `page` (int) - Número de página (default: 1)
- `limit` (int) - Items por página (default: 20, max: 100)

### Company Scoping

La mayoría de los endpoints están "scoped" a la empresa del usuario. Esto significa que solo verás/modificarás datos de tu propia empresa.

### Roles y Permisos

Jerarquía de roles (del mayor al menor):
1. **superadmin** (100) - Acceso global sin empresa
2. **admin** (50) - Administrador de empresa
3. **recruiter** (30) - Reclutador
4. **hiring_manager** (20) - Gerente de contratación
5. **user** (10) - Usuario básico

---

## 🚀 Ejemplos con cURL

### Login
```bash
curl -X POST http://localhost:8001/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@azentic.com",
    "password": "Admin123!"
  }'
```

### Listar Jobs
```bash
curl http://localhost:8001/api/v1/jobs \
  -H "Authorization: Bearer <your_token>"
```

### Crear Candidato
```bash
curl -X POST http://localhost:8001/api/v1/candidates \
  -H "Authorization: Bearer <your_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "candidato@example.com",
    "first_name": "Carlos",
    "last_name": "Ramírez",
    "phone": "+57 300 1234567"
  }'
```

### Obtener System Values por Categoría
```bash
curl http://localhost:8001/api/v1/system-values/job_status

# Con company_id para incluir valores personalizados
curl http://localhost:8001/api/v1/system-values/job_status \
  -H "X-Company-ID: 1"
```

### Crear System Value (SuperAdmin)
```bash
curl -X POST http://localhost:8001/api/v1/admin/system-values \
  -H "Authorization: Bearer <superadmin_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "category": "priority",
    "value": "critical",
    "label": "Crítico",
    "display_order": 5
  }'
```

---

## 📚 Recursos Adicionales

- [Flujo SuperAdmin](./FLUJO_SUPERADMIN.md)
- [Flujo Cliente](./FLUJO_CLIENTE.md)
- [Arquitectura del Sistema](./ARQUITECTURA.md)
- [Plan de Negocio Año 1](../negocio/PLAN_DE_NEGOCIO_ANO1.md)

---

**Última actualización:** 21 de Diciembre, 2025  
**Versión API:** v1.3.0
