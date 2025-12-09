# Implementación Multi-Company - Resumen

> **Fecha:** 8 de Diciembre, 2025  
> **Versión:** v1.3.0

---

## 🎯 **Objetivo**

Implementar sistema multi-company que permite:
1. Registro de empresa con su primer admin
2. Usuarios pertenecientes a múltiples empresas
3. Cambio de contexto entre empresas (switch company)
4. Mismo usuario con mismas credenciales en múltiples empresas

---

## ✅ **Cambios Implementados**

### **1. DTOs Nuevos** (`internal/app/dtos/auth_dto.go`)

```go
// Registro de empresa con admin
type RegisterCompanyDTO struct {
    CompanyName    string
    CompanySlug    string
    AdminEmail     string
    AdminPassword  string
    AdminFirstName string
    AdminLastName  string
    Timezone       string
}

// Response con empresa + admin + tokens
type RegisterCompanyResponseDTO struct {
    AccessToken  string
    RefreshToken string
    Company      CompanyResponse
    Admin        UserResponse
}

// Login response con lista de empresas
type LoginResponseWithCompaniesDTO struct {
    AccessToken  string
    RefreshToken string
    User         UserResponse
    Companies    []CompanyResponse // ← NUEVO
}

// Switch company request
type SwitchCompanyDTO struct {
    CompanyID uint
}

// Switch company response
type SwitchCompanyResponseDTO struct {
    AccessToken string
    Company     CompanyResponse
}
```

---

### **2. AuthService - Nuevos Métodos** (`internal/app/services/auth_service.go`)

#### **RegisterCompany** (Transaction)
```go
func (s *AuthService) RegisterCompany(dto *RegisterCompanyDTO) (*RegisterCompanyResponseDTO, error)
```

**Operaciones en transacción:**
1. Crear `Company` (plan: free)
2. Crear `User` (admin) con password hasheado
3. Crear `Membership` (user → company, role: admin, is_default: true)
4. Commit transaction
5. Generar tokens con `company_id`

**Resultado:** Empresa + Admin + Token funcional en una sola operación.

---

#### **LoginWithCompanies**
```go
func (s *AuthService) LoginWithCompanies(dto *LoginDTO) (*LoginResponseWithCompaniesDTO, error)
```

**Cambios:**
- Busca membership por defecto (`is_default=true`)
- Obtiene **todas las empresas** del usuario
- Retorna tokens + lista de empresas

**Response:**
```json
{
  "access_token": "...",
  "companies": [
    { "id": 1, "name": "Azentic Sys" },
    { "id": 2, "name": "DevCorp" }
  ]
}
```

---

#### **SwitchCompany**
```go
func (s *AuthService) SwitchCompany(userID uint, dto *SwitchCompanyDTO) (*SwitchCompanyResponseDTO, error)
```

**Operaciones:**
1. Valida que user tenga membership en esa empresa
2. Obtiene role del membership
3. Genera **nuevo token** con diferente `company_id`
4. Retorna nuevo token + info de empresa

**Uso:** Cambiar contexto sin re-login.

---

#### **GetUserCompanies**
```go
func (s *AuthService) GetUserCompanies(userID uint) ([]CompanyResponse, error)
```

**Operaciones:**
- Busca todas las memberships activas del usuario
- Carga empresas relacionadas (Preload)
- Retorna lista de empresas

---

### **3. AuthHandler - Nuevos Endpoints** (`internal/app/handlers/auth_handler.go`)

#### **POST /auth/register-company**
```go
func (h *AuthHandler) RegisterCompany(c *gin.Context)
```
- **Acceso:** Público
- **Crea:** Company + Admin + Membership
- **Response:** Tokens + Company + Admin

---

#### **POST /auth/login** (Modificado)
```go
func (h *AuthHandler) Login(c *gin.Context)
```
- **Acceso:** Público
- **Ahora retorna:** Lista de empresas del usuario
- **Token:** Usa empresa por defecto

---

#### **POST /auth/switch-company**
```go
func (h *AuthHandler) SwitchCompany(c *gin.Context)
```
- **Acceso:** Protegido (requiere auth)
- **Genera:** Nuevo token con diferente company_id
- **Response:** Nuevo token + info de empresa

---

#### **GET /auth/my-companies**
```go
func (h *AuthHandler) GetMyCompanies(c *gin.Context)
```
- **Acceso:** Protegido (requiere auth)
- **Response:** Lista de empresas del usuario
- **Uso:** Mostrar selector de empresas en UI

---

### **4. Routes** (`internal/platform/server/routes.go`)

**Rutas públicas:**
```go
POST /api/v1/auth/register-company  // NUEVO - Principal
POST /api/v1/auth/register          // DEPRECATED
POST /api/v1/auth/login
POST /api/v1/auth/refresh
```

**Rutas protegidas:**
```go
GET  /api/v1/auth/me
POST /api/v1/auth/change-password
POST /api/v1/auth/logout
POST /api/v1/auth/switch-company    // NUEVO
GET  /api/v1/auth/my-companies      // NUEVO
```

---

### **5. Documentación** (`docs/API_ENDPOINTS.md`)

**Actualizado:**
- ✅ Endpoint `POST /auth/register-company` con ejemplos
- ✅ Endpoint `POST /auth/login` ahora retorna lista de empresas
- ✅ Endpoint `POST /auth/switch-company` documentado
- ✅ Endpoint `GET /auth/my-companies` documentado
- ✅ Marcado `/auth/register` como DEPRECATED

---

## 📊 **Flujos Implementados**

### **Flujo 1: Registro de Nueva Empresa**

```
1. Usuario visita sitio
   ↓
2. POST /auth/register-company
   {
     "company_name": "Azentic Sys",
     "admin_email": "admin@azentic.com",
     "admin_password": "Admin123!",
     ...
   }
   ↓
3. Sistema (Transaction):
   - Crea Company (id: 1)
   - Crea User (id: 1)
   - Crea Membership (user: 1, company: 1, role: admin)
   ↓
4. Response:
   - access_token (con company_id: 1)
   - refresh_token
   - Company info
   - Admin info
   ↓
5. Usuario puede acceder inmediatamente a:
   - /jobs, /candidates, /applications, etc.
```

---

### **Flujo 2: Usuario con Múltiples Empresas**

#### **Escenario:**
Marcos tiene 2 empresas:
- Azentic Sys (id: 1)
- DevCorp (id: 2)

#### **Paso 1: Login**
```bash
POST /auth/login
{
  "email": "marcos@email.com",
  "password": "Password123!"
}
```

**Response:**
```json
{
  "access_token": "eyJ...",  // company_id: 1 (default)
  "companies": [
    { "id": 1, "name": "Azentic Sys" },
    { "id": 2, "name": "DevCorp" }
  ]
}
```

#### **Paso 2: Trabajar en Azentic**
```bash
GET /jobs
Authorization: Bearer <token_empresa_1>

# Ve solo jobs de Azentic Sys
```

#### **Paso 3: Cambiar a DevCorp**
```bash
POST /auth/switch-company
Authorization: Bearer <token_empresa_1>
{
  "company_id": 2
}
```

**Response:**
```json
{
  "access_token": "eyJ...",  // Nuevo token con company_id: 2
  "company": {
    "id": 2,
    "name": "DevCorp"
  }
}
```

#### **Paso 4: Trabajar en DevCorp**
```bash
GET /jobs
Authorization: Bearer <token_empresa_2>

# Ahora ve solo jobs de DevCorp
```

---

### **Flujo 3: Agregar Usuario a Segunda Empresa (MVP - SuperAdmin)**

#### **Situación:**
- Marcos ya tiene cuenta (empresa 1)
- Ahora quiere agregar empresa 2

#### **SuperAdmin asigna a Marcos:**
```bash
POST /admin/memberships
Authorization: Bearer <superadmin_token>
{
  "user_id": 1,
  "company_id": 2,
  "role": "recruiter"
}
```

**Restricción MVP:**
- **SOLO SuperAdmin** puede crear memberships
- Clientes intentando POST /memberships reciben: `403 Forbidden`
- Mensaje: "Only superadmin can assign users to companies. Regular users should create new users instead."

**Sistema:**
1. Valida que role == "superadmin"
2. Crea nuevo Membership (user: 1, company: 2, role: recruiter)

**Resultado:**
- Marcos ahora tiene 2 memberships
- Mismo email, misma password
- Puede hacer switch entre empresas

**Fase 2 (Futuro):**
- Sistema de invitación por email
- Admin envía invitación → Usuario acepta → Membership creado

---

## 🔐 **JWT Claims por Contexto**

### **Token para Empresa 1:**
```json
{
  "user_id": 1,
  "company_id": 1,
  "email": "marcos@email.com",
  "role": "admin",
  "exp": 1733687400
}
```

### **Token para Empresa 2:**
```json
{
  "user_id": 1,
  "company_id": 2,      // ← Diferente empresa
  "email": "marcos@email.com",
  "role": "recruiter",  // ← Puede tener diferente rol
  "exp": 1733687500
}
```

**Importante:** Mismo usuario, diferente contexto empresarial.

---

## 🗄️ **Estructura de Base de Datos**

### **Tabla `users`**
| ID | Email | PasswordHash | FirstName | LastName |
|----|-------|--------------|-----------|----------|
| 1 | marcos@email.com | $2a$10$... | Marcos | Ramos |

### **Tabla `companies`**
| ID | Name | Slug | PlanTier |
|----|------|------|----------|
| 1 | Azentic Sys | azentic-sys | professional |
| 2 | DevCorp | devcorp | enterprise |

### **Tabla `memberships`**
| ID | UserID | CompanyID | Role | IsDefault | Status |
|----|--------|-----------|------|-----------|--------|
| 1 | 1 | 1 | admin | ✅ true | active |
| 2 | 1 | 2 | recruiter | ❌ false | active |

**Relación:** 
- 1 User → N Memberships
- 1 Company → N Memberships
- 1 Membership → 1 User + 1 Company

---

## 📝 **Casos de Uso Soportados**

### ✅ **Caso 1: Empresa Nueva**
```
Usuario → Registro empresa → Admin automático → Acceso inmediato
```

### ✅ **Caso 2: Multi-Empresa**
```
Usuario → Login → Ve lista empresas → Switch company → Trabaja en empresa 2
```

### ✅ **Caso 3: Mismo Usuario, Múltiples Empresas**
```
User existente → Admin 2 crea membership → User hace switch → Accede a ambas
```

### ✅ **Caso 4: Roles Diferentes**
```
Marcos: Admin en Azentic, Recruiter en DevCorp
Token refleja role según empresa activa
```

---

## 🚀 **Endpoints Disponibles**

| Endpoint | Método | Acceso | Descripción |
|----------|--------|--------|-------------|
| `/auth/register-company` | POST | Público | Crear empresa + admin |
| `/auth/register` | POST | Público | ⚠️ DEPRECATED |
| `/auth/login` | POST | Público | Login con lista empresas |
| `/auth/refresh` | POST | Público | Renovar token |
| `/auth/me` | GET | Protegido | Info usuario |
| `/auth/change-password` | POST | Protegido | Cambiar password |
| `/auth/logout` | POST | Protegido | Logout |
| `/auth/switch-company` | POST | Protegido | Cambiar empresa |
| `/auth/my-companies` | GET | Protegido | Listar empresas |

**Total:** 9 endpoints de autenticación

---

## 🔧 **Cambios Técnicos**

### **Errores Nuevos:**
```go
ErrCompanyNotFound  = errors.New("company not found")
ErrNoMembership     = errors.New("user does not belong to this company")
```

### **Métodos de Servicio:**
- `RegisterCompany()` - Transaction completa
- `LoginWithCompanies()` - Login mejorado
- `SwitchCompany()` - Cambio de contexto
- `GetUserCompanies()` - Lista empresas

### **Handlers:**
- `RegisterCompany()` - Endpoint público
- `SwitchCompany()` - Endpoint protegido
- `GetMyCompanies()` - Endpoint protegido

---

## 📈 **Métricas de Implementación**

| Métrica | Valor |
|---------|-------|
| **Archivos Modificados** | 5 |
| **Archivos Creados** | 0 (usamos existentes) |
| **DTOs Nuevos** | 4 |
| **Métodos Service** | 4 |
| **Handlers Nuevos** | 3 |
| **Endpoints Nuevos** | 3 |
| **Líneas de Código** | ~250 |

---

## ✅ **Testing Manual**

### **1. Registrar Empresa**
```bash
curl -X POST http://localhost:8001/api/v1/auth/register-company \
  -H "Content-Type: application/json" \
  -d '{
    "company_name": "Test Corp",
    "company_slug": "test-corp",
    "admin_email": "admin@test.com",
    "admin_password": "Admin123!",
    "admin_first_name": "Admin",
    "admin_last_name": "Test"
  }'
```

### **2. Login**
```bash
curl -X POST http://localhost:8001/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@test.com",
    "password": "Admin123!"
  }'
```

### **3. Ver Mis Empresas**
```bash
curl http://localhost:8001/api/v1/auth/my-companies \
  -H "Authorization: Bearer <token>"
```

### **4. Cambiar de Empresa**
```bash
curl -X POST http://localhost:8001/api/v1/auth/switch-company \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{"company_id": 2}'
```

---

## 🎯 **Próximos Pasos**

### **Fase 2: Sistema de Invitaciones** (PRIORITARIO)
- `POST /memberships/invite` - Admin invita por email
- `POST /auth/accept-invite` - Usuario acepta invitación
- Sistema de tokens de invitación con expiración
- **Razón:** Actualmente solo SuperAdmin puede crear memberships, necesitamos que admins de empresa puedan invitar usuarios de forma segura

### **Gestión de Memberships** (✅ Implementado)
- `GET /memberships` - Ver memberships de mi empresa ✅
- `GET /memberships/:id` - Ver detalle ✅
- `PUT /memberships/:id` - Cambiar role ✅
- `DELETE /memberships/:id` - Remover usuario de empresa ✅
- `POST /admin/memberships` - SuperAdmin asigna usuarios ✅

### **Fase 3: Onboarding**
- Tutorial post-registro
- Configuración inicial de empresa
- Importación de datos

---

## 📚 **Documentación Relacionada**

- [API_ENDPOINTS.md](./API_ENDPOINTS.md) - Referencia completa de endpoints
- [SUPERADMIN.md](./SUPERADMIN.md) - Gestión global de empresas
- [ARCHITECTURE.md](./ARCHITECTURE.md) - Arquitectura técnica

---

**Última actualización:** 8 de Diciembre, 2025  
**Versión:** v1.3.0  
**Status:** ✅ Implementado y funcional
