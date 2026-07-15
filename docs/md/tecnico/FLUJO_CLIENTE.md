# 👥 Cliente (Admin de Empresa) - Flujo Completo del Sistema

> **Documento de Referencia Técnica**  
> Versión: 1.1 | Actualizado: 15 de Julio, 2026
> Describe en detalle qué puede y NO puede hacer un Cliente en el ATS

> La definición canónica de los contextos está en
> [ADR-002](ADR-002-contextos-landlord-tenant.md). Los ejemplos antiguos que
> contradigan ese ADR deben considerarse históricos.

---

## 📋 Índice

1. [¿Qué es un Cliente?](#qué-es-un-cliente)
2. [Jerarquía de Roles](#jerarquía-de-roles)
3. [Registro y Onboarding](#registro-y-onboarding)
4. [Autenticación y Multi-Empresa](#autenticación-y-multi-empresa)
5. [Endpoints Disponibles por Módulo](#endpoints-disponibles-por-módulo)
6. [Permisos Detallados por Rol](#permisos-detallados-por-rol)
7. [Aislamiento Multi-Tenant](#aislamiento-multi-tenant)
8. [Flujo Completo de Uso](#flujo-completo-de-uso)
9. [Casos de Uso Reales](#casos-de-uso-reales)
10. [Restricciones y Límites](#restricciones-y-límites)

---

## 1. ¿Qué es un Cliente?

### **Definición**

Un **Cliente** es cualquier usuario que pertenece a una o más **empresas** en el sistema Dvra ATS. A diferencia del SuperAdmin, los clientes tienen su acceso **limitado y aislado** a los datos de su(s) empresa(s).

### **Características Principales**

```
┌─────────────────────────────────────────────────────────────┐
│                    CLIENTE (Company User)                    │
├─────────────────────────────────────────────────────────────┤
│  • Pertenece a UNA o MÁS empresas                           │
│  • Tiene company_id en JWT token                            │
│  • Ve SOLO datos de su empresa actual                       │
│  • Puede tener diferentes roles por empresa                 │
│  • Puede hacer switch entre empresas                        │
│  • Sujeto a límites del plan de su empresa                  │
└─────────────────────────────────────────────────────────────┘
```

### **Diferencias con SuperAdmin**

| Característica | Cliente | SuperAdmin |
|---------------|---------|------------|
| **Contexto JWT** | `tenant` con `company_id` | `platform` sin `company_id` |
| **Alcance de datos** | Solo su empresa | Todas las empresas |
| **Puede crear empresas** | ❌ No | ✅ Sí |
| **Gestionar planes** | ❌ No | ✅ Sí |
| **Ver otras empresas** | ❌ No | ✅ Sí |
| **Crear memberships** | Según permisos tenant vigentes | ✅ Desde `/platform/memberships` |
| **Multi-empresa** | ✅ Sí (via switch) | Puede cambiar a tenant si también posee membership |
| **Límites de plan** | ✅ Aplican | ❌ Ilimitado |

---

## 2. Jerarquía de Roles

### **2.1 Roles Disponibles para Clientes**

```go
// internal/shared/permissions/permissions.go

const (
    RoleAdmin          = "admin"           // Level 50
    RoleRecruiter      = "recruiter"       // Level 30
    RoleHiringManager  = "hiring_manager"  // Level 20
    RoleUser           = "user"            // Level 10
)
```

### **2.2 Niveles de Acceso**

```
┌─────────────────────────────────────────────────────────────┐
│                    JERARQUÍA DE ROLES                        │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  🔴 ADMIN (Level 50)                                         │
│     • Control total de la empresa                           │
│     • Gestiona usuarios, billing, configuración             │
│     • Crea/edita/elimina jobs, candidatos, aplicaciones     │
│     • Ve todos los datos de la empresa                      │
│                                                              │
│  🟡 RECRUITER (Level 30)                                     │
│     • Gestiona proceso de reclutamiento                     │
│     • Crea/edita jobs y candidatos                          │
│     • Gestiona aplicaciones y pipeline                      │
│     • NO puede invitar usuarios ni cambiar billing          │
│                                                              │
│  🟢 HIRING_MANAGER (Level 20)                                │
│     • Ve candidatos de jobs asignados                       │
│     • Puede comentar y calificar aplicaciones               │
│     • NO puede crear jobs ni candidatos                     │
│     • Acceso limitado de solo lectura                       │
│                                                              │
│  ⚪ USER (Level 10)                                          │
│     • Solo lectura                                          │
│     • Ve reportes básicos                                   │
│     • NO puede modificar nada                               │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

### **2.3 Matriz de Permisos por Rol**

| Acción | Admin | Recruiter | Hiring Mgr | User |
|--------|-------|-----------|------------|------|
| **Company Settings** |
| Ver configuración empresa | ✅ | ❌ | ❌ | ❌ |
| Editar empresa | ✅ | ❌ | ❌ | ❌ |
| Ver billing/plan | ✅ | ❌ | ❌ | ❌ |
| **Gestión de Usuarios** |
| Ver team members | ✅ | ✅ | ✅ | ✅ |
| Crear usuarios | ✅ | ❌ | ❌ | ❌ |
| Editar roles | ✅ | ❌ | ❌ | ❌ |
| Remover usuarios | ✅ | ❌ | ❌ | ❌ |
| **Jobs** |
| Ver todos los jobs | ✅ | ✅ | Solo asignados | ✅ |
| Crear job | ✅ | ✅ | ❌ | ❌ |
| Editar job | ✅ | ✅ | Solo asignados | ❌ |
| Eliminar job | ✅ | ✅ | ❌ | ❌ |
| Publicar/cerrar job | ✅ | ✅ | ❌ | ❌ |
| **Candidates** |
| Ver todos candidatos | ✅ | ✅ | Solo de sus jobs | Solo de sus jobs |
| Crear candidato | ✅ | ✅ | ❌ | ❌ |
| Editar candidato | ✅ | ✅ | ❌ | ❌ |
| Eliminar candidato | ✅ | ✅ | ❌ | ❌ |
| **Applications** |
| Ver aplicaciones | ✅ | ✅ | Solo de sus jobs | Solo de sus jobs |
| Cambiar stage | ✅ | ✅ | Solo de sus jobs | ❌ |
| Calificar candidato | ✅ | ✅ | ✅ | ❌ |
| Agregar notas | ✅ | ✅ | ✅ | ❌ |
| Rechazar/contratar | ✅ | ✅ | Solo de sus jobs | ❌ |
| **Memberships** |
| Ver memberships | ✅ | ✅ | ✅ | ✅ |
| Cambiar roles | ✅ | ❌ | ❌ | ❌ |
| Remover de empresa | ✅ | ❌ | ❌ | ❌ |
| Crear membership | ❌ (MVP) | ❌ | ❌ | ❌ |

---

## 3. Registro y Onboarding

### **3.1 Registro de Nueva Empresa (Flujo Principal)**

**Endpoint:**
```http
POST /api/v1/auth/register-company
Content-Type: application/json
```

**Request Body:**
```json
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

**Response (201):**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "company": {
    "id": 5,
    "name": "Mi Startup Tech",
    "slug": "mi-startup-tech",
    "plan_tier": "free",
    "trial_ends_at": "2026-01-09T10:00:00Z"
  },
  "admin": {
    "id": 25,
    "email": "ceo@mistartup.com",
    "first_name": "Juan",
    "last_name": "Pérez",
    "is_active": true
  }
}
```

### **Proceso Interno (Transaction Completa):**

```go
// internal/app/services/auth_service.go

func (s *AuthService) RegisterCompany(dto) (*RegisterCompanyResponseDTO, error) {
    tx := s.db.Begin()
    
    // 1. Validar que email no exista
    existingUser, _ := s.userRepo.FindByEmail(dto.AdminEmail)
    if existingUser != nil {
        return errors.New("Email already exists")
    }
    
    // 2. Validar que plan "free" existe y está activo
    freePlan, err := s.planRepo.FindActiveBySlug("free")
    if err != nil {
        return errors.New("Free plan is not available, please contact support")
    }
    
    // 3. Crear empresa con plan "free"
    trialEnds := time.Now().AddDate(0, 1, 0) // 1 mes de trial
    company := models.Company{
        Name:        dto.CompanyName,
        Slug:        dto.CompanySlug,
        PlanTier:    freePlan.Slug,  // ✅ Plan validado
        TrialEndsAt: &trialEnds,
        Timezone:    dto.Timezone,
    }
    tx.Create(&company)
    
    // 4. Crear usuario admin
    hashedPassword := bcrypt.GenerateFromPassword(dto.AdminPassword)
    admin := models.User{
        Email:        dto.AdminEmail,
        PasswordHash: string(hashedPassword),
        FirstName:    dto.AdminFirstName,
        LastName:     dto.AdminLastName,
        IsActive:     true,
    }
    tx.Create(&admin)
    
    // 5. Crear membership (user → company con role "admin")
    membership := models.Membership{
        UserID:    admin.ID,
        CompanyID: &company.ID,
        Role:      models.RoleAdmin,  // ✅ Admin automáticamente
        Status:    "active",
        IsDefault: true,  // ✅ Empresa por defecto
        JoinedAt:  &now,
    }
    tx.Create(&membership)
    
    // 6. Generar tokens JWT
    accessToken := s.jwtService.GenerateAccessToken(JWTClaims{
        UserID:    admin.ID,
        CompanyID: company.ID,  // ✅ Con empresa
        Email:     admin.Email,
        Role:      models.RoleAdmin,
    })
    
    refreshToken := s.jwtService.GenerateRefreshToken(admin.ID)
    
    tx.Commit()
    
    return &RegisterCompanyResponseDTO{
        AccessToken:  accessToken,
        RefreshToken: refreshToken,
        Company:      company,
        Admin:        admin,
    }
}
```

**Validaciones Implementadas:**
- ✅ Email único (no duplicados)
- ✅ Plan "free" existe y está activo
- ✅ Slug de empresa único
- ✅ Password hasheado con bcrypt
- ✅ Transaction completa (rollback si falla)
- ✅ Membership con role "admin" automático
- ✅ Token JWT generado con company_id

---

### **3.2 Registro Alternativo (Usuario sin Empresa - DEPRECATED)**

```http
POST /api/v1/auth/register
```

⚠️ **DEPRECADO:** Este endpoint crea un usuario sin empresa, requiere que SuperAdmin lo asigne después.

**Flujo Recomendado:** Siempre usar `/auth/register-company`

---

## 4. Autenticación y Multi-Empresa

### **4.1 Login de Cliente**

**Endpoint:**
```http
POST /api/v1/auth/login
Content-Type: application/json
```

**Request:**
```json
{
  "email": "ceo@mistartup.com",
  "password": "SecurePass123!"
}
```

**Response (200) - Usuario con 1 Empresa:**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 25,
    "email": "ceo@mistartup.com",
    "first_name": "Juan",
    "last_name": "Pérez"
  },
  "company": {
    "id": 5,
    "name": "Mi Startup Tech",
    "slug": "mi-startup-tech",
    "plan_tier": "free"
  },
  "role": "admin"
}
```

**Response (200) - Usuario con Múltiples Empresas:**
```json
{
  "user": {
    "id": 25,
    "email": "ceo@mistartup.com",
    "first_name": "Juan",
    "last_name": "Pérez"
  },
  "companies": [
    {
      "id": 5,
      "name": "Mi Startup Tech",
      "slug": "mi-startup-tech",
      "role": "admin",
      "is_default": true
    },
    {
      "id": 12,
      "name": "DevCorp",
      "slug": "devcorp",
      "role": "recruiter",
      "is_default": false
    }
  ],
  "message": "Please select a company",
  "requires_company_selection": true
}
```

### **JWT Token Structure**

```go
// Token de cliente (con empresa)
{
  "user_id": 25,
  "company_id": 5,        // ✅ Empresa del contexto actual
  "email": "ceo@mistartup.com",
  "role": "admin",        // ✅ Rol en esta empresa
  "exp": 1733687400
}
```

---

### **4.2 Multi-Empresa: Switch Company**

**Endpoint:**
```http
POST /api/v1/auth/switch-company
Authorization: Bearer {token}
Content-Type: application/json
```

**Request:**
```json
{
  "company_id": 12
}
```

**Response (200):**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "company": {
    "id": 12,
    "name": "DevCorp",
    "slug": "devcorp",
    "plan_tier": "professional"
  },
  "role": "recruiter"
}
```

**Nuevo Token Generado:**
```go
{
  "user_id": 25,
  "company_id": 12,       // ← CAMBIÓ a empresa 12
  "email": "ceo@mistartup.com",
  "role": "recruiter",    // ← Rol diferente en esta empresa
  "exp": 1733687500
}
```

**Validaciones:**
```go
func (s *AuthService) SwitchCompany(userID, companyID) {
    // 1. Verificar que membership existe
    membership := FindMembership(userID, companyID)
    if membership == nil {
        return 404 "You are not a member of this company"
    }
    
    // 2. Verificar que membership está activa
    if membership.Status != "active" {
        return 403 "Your membership is not active"
    }
    
    // 3. Generar nuevo token con nuevo company_id y role
    newToken := GenerateToken({
        user_id: userID,
        company_id: companyID,
        role: membership.Role  // Puede ser diferente
    })
}
```

---

### **4.3 Ver Mis Empresas**

**Endpoint:**
```http
GET /api/v1/auth/my-companies
Authorization: Bearer {token}
```

**Response (200):**
```json
{
  "companies": [
    {
      "id": 5,
      "name": "Mi Startup Tech",
      "slug": "mi-startup-tech",
      "plan_tier": "free",
      "role": "admin",
      "status": "active",
      "is_default": true
    },
    {
      "id": 12,
      "name": "DevCorp",
      "slug": "devcorp",
      "plan_tier": "professional",
      "role": "recruiter",
      "status": "active",
      "is_default": false
    }
  ],
  "count": 2
}
```

---

## 5. Endpoints Disponibles por Módulo

### **5.1 Gestión de Empresa (Company)**

#### **📋 Ver Mi Empresa**
```http
GET /api/v1/companies
Authorization: Bearer {token}
```

**Response:**
```json
{
  "status": "success",
  "message": "Companies retrieved successfully",
  "data": {
    "companies": [
      {
        "id": 5,
        "name": "Mi Startup Tech",
        "slug": "mi-startup-tech",
        "plan_tier": "free",
        "timezone": "America/Bogota",
        "created_at": "2025-12-01T10:00:00Z",
        "trial_ends_at": "2026-01-01T10:00:00Z"
      }
    ],
    "count": 1
  }
}
```

**Nota:** Cliente solo ve SU empresa, no otras.

---

#### **🔍 Ver Detalles de Mi Empresa**
```http
GET /api/v1/companies/:id
Authorization: Bearer {token}
```

**Validación:**
```go
func (h *CompanyHandler) GetCompany(c *gin.Context) {
    id := c.Param("id")
    companyID, _ := c.Get("company_id")
    
    // ✅ Solo puede ver su propia empresa
    if id != companyID {
        return 403 "Access denied"
    }
    
    company := h.service.GetCompanyByID(id)
    return company
}
```

---

#### **✏️ Actualizar Mi Empresa (Admin only)**
```http
PUT /api/v1/companies/:id
Authorization: Bearer {token}
Content-Type: application/json
```

**Request:**
```json
{
  "name": "Mi Startup Tech SAS",
  "timezone": "America/Bogota",
  "website": "https://mistartup.com"
}
```

**Validaciones:**
- ✅ Solo Admin puede editar
- ✅ Solo puede editar su propia empresa
- ✅ NO puede cambiar `plan_tier` (solo SuperAdmin)

---

### **5.2 Gestión de Usuarios (Users)**

#### **📋 Listar Usuarios de Mi Empresa**
```http
GET /api/v1/users
Authorization: Bearer {token}
```

**Response:**
```json
{
  "status": "success",
  "message": "Users retrieved successfully",
  "data": {
    "users": [
      {
        "id": 25,
        "email": "ceo@mistartup.com",
        "first_name": "Juan",
        "last_name": "Pérez",
        "is_active": true,
        "role": "admin"
      },
      {
        "id": 30,
        "email": "recruiter@mistartup.com",
        "first_name": "Ana",
        "last_name": "Gómez",
        "is_active": true,
        "role": "recruiter"
      }
    ],
    "count": 2
  }
}
```

**Query Interna:**
```sql
SELECT users.* 
FROM users
JOIN memberships ON memberships.user_id = users.id
WHERE memberships.company_id = :company_id_from_token
  AND memberships.status = 'active'
```

---

#### **➕ Crear Usuario (Admin only)**
```http
POST /api/v1/users
Authorization: Bearer {token}
Content-Type: application/json
```

**Request:**
```json
{
  "email": "developer@mistartup.com",
  "password": "TempPass123!",
  "first_name": "Carlos",
  "last_name": "Ramírez",
  "role": "user"
}
```

**Proceso Interno:**
```go
func (h *UserHandler) CreateUser(c *gin.Context) {
    // 1. Validar que solo Admin puede crear usuarios
    role, _ := c.Get("role")
    if role != "admin" {
        return 403 "Only admin can create users"
    }
    
    // 2. Forzar company_id del token
    companyID, _ := c.Get("company_id")
    
    tx.Begin()
    
    // 3. Crear usuario
    user := models.User{
        Email:        dto.Email,
        PasswordHash: bcrypt.Hash(dto.Password),
        FirstName:    dto.FirstName,
        LastName:     dto.LastName,
        IsActive:     true,
    }
    tx.Create(&user)
    
    // 4. Crear membership automáticamente
    membership := models.Membership{
        UserID:    user.ID,
        CompanyID: &companyID,  // ✅ Empresa del token
        Role:      dto.Role,
        Status:    "active",
        IsDefault: true,
    }
    tx.Create(&membership)
    
    tx.Commit()
}
```

**Response (201):**
```json
{
  "status": "success",
  "data": {
    "user": {
      "id": 35,
      "email": "developer@mistartup.com",
      "first_name": "Carlos",
      "last_name": "Ramírez",
      "is_active": true
    },
    "membership": {
      "id": 15,
      "role": "user",
      "status": "active"
    }
  }
}
```

---

#### **✏️ Actualizar Usuario (Admin only)**
```http
PUT /api/v1/users/:id
Authorization: Bearer {token}
```

**Validaciones:**
- ✅ Solo Admin puede actualizar
- ✅ Usuario debe pertenecer a la empresa
- ✅ NO puede cambiar su propio rol (prevención)

---

#### **🗑️ Eliminar Usuario (Admin only)**
```http
DELETE /api/v1/users/:id
Authorization: Bearer {token}
```

**Efecto:** Soft delete del usuario + eliminación de membership

---

### **5.3 Gestión de Jobs (Ofertas de Trabajo)**

#### **📋 Listar Jobs de Mi Empresa**
```http
GET /api/v1/jobs
Authorization: Bearer {token}
```

**Response:**
```json
{
  "status": "success",
  "data": {
    "jobs": [
      {
        "id": 10,
        "title": "Senior Full Stack Developer",
        "department": "Engineering",
        "location": "Remote",
        "employment_type": "full-time",
        "status": "published",
        "company_id": 5,
        "created_at": "2025-12-01T10:00:00Z"
      },
      {
        "id": 15,
        "title": "Product Designer",
        "department": "Design",
        "location": "Bogotá",
        "employment_type": "full-time",
        "status": "draft",
        "company_id": 5,
        "created_at": "2025-12-05T14:30:00Z"
      }
    ],
    "count": 2
  }
}
```

**Query:**
```sql
SELECT * FROM jobs 
WHERE company_id = :company_id_from_token
ORDER BY created_at DESC
```

---

#### **➕ Crear Job (Admin/Recruiter)**
```http
POST /api/v1/jobs
Authorization: Bearer {token}
Content-Type: application/json
```

**Request:**
```json
{
  "title": "Backend Developer",
  "description": "Buscamos desarrollador backend con experiencia en Go...",
  "department": "Engineering",
  "location": "Remote",
  "employment_type": "full-time",
  "salary_min": 50000,
  "salary_max": 80000,
  "salary_currency": "USD",
  "requirements": "5+ años de experiencia, Go, PostgreSQL...",
  "status": "draft"
}
```

**Validación de Seguridad:**
```go
func (h *JobHandler) CreateJob(c *gin.Context) {
    role, _ := c.Get("role")
    
    // ✅ Solo Admin o Recruiter pueden crear jobs
    if role != "admin" && role != "recruiter" {
        return 403 "Insufficient permissions"
    }
    
    // ✅ Forzar company_id del token (previene manipulación)
    companyID, _ := c.Get("company_id")
    dto.CompanyID = companyID  // ← Ignora cualquier company_id enviado
    
    job := h.service.CreateJob(dto)
    return 201 job
}
```

**Response (201):**
```json
{
  "status": "success",
  "data": {
    "id": 20,
    "title": "Backend Developer",
    "company_id": 5,
    "status": "draft",
    "created_at": "2025-12-09T10:00:00Z"
  }
}
```

---

#### **✏️ Actualizar Job**
```http
PUT /api/v1/jobs/:id
Authorization: Bearer {token}
```

**Validaciones:**
- ✅ Job debe pertenecer a mi empresa
- ✅ Solo Admin/Recruiter pueden editar
- ✅ Hiring Manager solo si es su job asignado

---

#### **🗑️ Eliminar Job (Soft Delete)**
```http
DELETE /api/v1/jobs/:id
Authorization: Bearer {token}
```

**Efecto:** 
- Marca `deleted_at` (soft delete)
- Mantiene historial de aplicaciones
- Solo Admin/Recruiter

---

### **5.4 Gestión de Candidatos (Candidates)**

#### **📋 Listar Candidatos**
```http
GET /api/v1/candidates
Authorization: Bearer {token}
```

**Filtrado por Rol:**
```go
func (h *CandidateHandler) GetCandidates(c *gin.Context) {
    companyID, _ := c.Get("company_id")
    role, _ := c.Get("role")
    
    if role == "admin" || role == "recruiter" {
        // ✅ Ven TODOS los candidatos de la empresa
        candidates := h.service.GetCandidatesByCompanyID(companyID)
    } else {
        // Hiring Manager/User: solo de sus jobs asignados
        userID, _ := c.Get("user_id")
        candidates := h.service.GetCandidatesForUser(userID, companyID)
    }
}
```

**Response:**
```json
{
  "status": "success",
  "data": {
    "candidates": [
      {
        "id": 50,
        "email": "candidato@example.com",
        "first_name": "María",
        "last_name": "López",
        "phone": "+57 300 1234567",
        "location": "Medellín, Colombia",
        "linkedin_url": "https://linkedin.com/in/marialopez",
        "github_url": "https://github.com/marialopez",
        "resume_url": "https://s3.../resume.pdf",
        "source": "linkedin",
        "company_id": 5,
        "created_at": "2025-12-02T14:00:00Z"
      }
    ],
    "count": 15
  }
}
```

---

#### **➕ Crear Candidato (Admin/Recruiter)**
```http
POST /api/v1/candidates
Authorization: Bearer {token}
Content-Type: application/json
```

**Request:**
```json
{
  "email": "nuevo@candidate.com",
  "first_name": "Pedro",
  "last_name": "Sánchez",
  "phone": "+57 310 9876543",
  "location": "Bogotá, Colombia",
  "linkedin_url": "https://linkedin.com/in/pedrosanchez",
  "source": "referral",
  "source_details": "Referido por Juan Pérez"
}
```

**Validación:**
```go
func (h *CandidateHandler) CreateCandidate(c *gin.Context) {
    role, _ := c.Get("role")
    
    // ✅ Solo Admin/Recruiter
    if role != "admin" && role != "recruiter" {
        return 403 "Insufficient permissions"
    }
    
    // ✅ Forzar company_id
    companyID, _ := c.Get("company_id")
    dto.CompanyID = companyID
    
    // ✅ Verificar email único dentro de la empresa
    existing := h.service.FindByEmailAndCompany(dto.Email, companyID)
    if existing != nil {
        return 409 "Candidate already exists in your company"
    }
    
    candidate := h.service.CreateCandidate(dto)
}
```

---

### **5.5 Gestión de Aplicaciones (Applications)**

#### **📋 Listar Aplicaciones**
```http
GET /api/v1/applications
Authorization: Bearer {token}

Query Params:
  - job_id (int): Filtrar por job
  - stage (string): Filtrar por etapa (applied, screening, technical, offer, hired, rejected)
  - candidate_id (int): Filtrar por candidato
```

**Response:**
```json
{
  "status": "success",
  "data": {
    "applications": [
      {
        "id": 100,
        "job_id": 10,
        "candidate_id": 50,
        "stage": "technical",
        "rating": 4,
        "notes": "Buen desempeño en prueba técnica",
        "applied_at": "2025-12-01T10:00:00Z",
        "company_id": 5,
        "job": {
          "id": 10,
          "title": "Senior Full Stack Developer"
        },
        "candidate": {
          "id": 50,
          "first_name": "María",
          "last_name": "López",
          "email": "candidato@example.com"
        }
      }
    ],
    "count": 25
  }
}
```

---

#### **➕ Crear Aplicación**
```http
POST /api/v1/applications
Authorization: Bearer {token}
Content-Type: application/json
```

**Request:**
```json
{
  "job_id": 10,
  "candidate_id": 50,
  "notes": "Aplicación recibida via LinkedIn"
}
```

**Validaciones:**
- ✅ Job y Candidate deben pertenecer a mi empresa
- ✅ Candidato no puede tener aplicación duplicada al mismo job
- ✅ Job debe estar en estado `published`

---

#### **✏️ Actualizar Aplicación (Cambiar Stage)**
```http
PUT /api/v1/applications/:id
Authorization: Bearer {token}
Content-Type: application/json
```

**Request:**
```json
{
  "stage": "offer",
  "rating": 5,
  "notes": "Excelente candidato, preparar oferta"
}
```

**Pipeline de Stages:**
```
applied → screening → technical → offer → hired
   ↓         ↓           ↓         ↓
rejected  rejected    rejected  rejected
```

**Validaciones:**
- ✅ Solo Admin/Recruiter pueden cambiar stage
- ✅ Hiring Manager solo de sus jobs
- ✅ Timestamps automáticos (hired_at, rejected_at)

---

### **5.6 Gestión de Memberships**

#### **📋 Ver Memberships de Mi Empresa**
```http
GET /api/v1/memberships
Authorization: Bearer {token}
```

**Response:**
```json
{
  "status": "success",
  "data": {
    "memberships": [
      {
        "id": 10,
        "user_id": 25,
        "company_id": 5,
        "role": "admin",
        "status": "active",
        "is_default": true,
        "joined_at": "2025-12-01T10:00:00Z",
        "user": {
          "id": 25,
          "email": "ceo@mistartup.com",
          "first_name": "Juan",
          "last_name": "Pérez"
        }
      },
      {
        "id": 15,
        "user_id": 30,
        "company_id": 5,
        "role": "recruiter",
        "status": "active",
        "is_default": true,
        "joined_at": "2025-12-03T14:00:00Z",
        "user": {
          "id": 30,
          "email": "recruiter@mistartup.com",
          "first_name": "Ana",
          "last_name": "Gómez"
        }
      }
    ],
    "count": 2
  }
}
```

---

#### **✏️ Actualizar Membership (Cambiar Rol - Admin only)**
```http
PUT /api/v1/memberships/:id
Authorization: Bearer {token}
Content-Type: application/json
```

**Request:**
```json
{
  "role": "admin",
  "status": "active"
}
```

**Validaciones:**
- ✅ Solo Admin puede cambiar roles
- ✅ Membership debe ser de mi empresa
- ✅ NO puede cambiar su propio rol (prevención)

---

#### **🗑️ Remover Usuario de Empresa (Admin only)**
```http
DELETE /api/v1/memberships/:id
Authorization: Bearer {token}
```

**Efecto:**
- Usuario pierde acceso a la empresa
- Si tiene múltiples empresas, sigue en las demás
- Si es su única empresa, debe ser re-invitado

---

#### **❌ Crear Membership (RESTRINGIDO en MVP)**

```http
POST /api/v1/memberships  ← NO DISPONIBLE
```

**Restricción MVP:**
```go
func (h *MembershipHandler) CreateMembership(c *gin.Context) {
    role, _ := c.Get("role")
    
    if role != "superadmin" {
        c.JSON(403, gin.H{
            "error": "Only superadmin can assign users to companies. Regular users should create new users instead.",
        })
        return
    }
}
```

**Alternativa:**
- Admin crea nuevo usuario con POST `/users`
- Usuario se crea automáticamente en la empresa del admin
- Para agregar usuario existente → contactar SuperAdmin

---

## 6. Permisos Detallados por Rol

### **6.1 Admin (Level 50)**

**Control Total de la Empresa**

```
✅ PUEDE HACER:
├── Company
│   ├── Ver configuración completa
│   ├── Editar nombre, timezone, website
│   ├── Ver plan actual y límites
│   └── Ver billing (cuando esté implementado)
├── Users
│   ├── Ver todos los usuarios
│   ├── Crear nuevos usuarios
│   ├── Editar usuarios
│   ├── Cambiar roles de otros
│   └── Eliminar usuarios
├── Jobs
│   ├── Ver, crear, editar, eliminar
│   ├── Publicar/cerrar jobs
│   └── Asignar recruiters/hiring managers
├── Candidates
│   ├── Ver, crear, editar, eliminar
│   └── Acceso completo a todos los candidatos
├── Applications
│   ├── Ver todas las aplicaciones
│   ├── Cambiar stages (pipeline)
│   ├── Calificar candidatos
│   ├── Agregar notas
│   └── Contratar/rechazar
└── Memberships
    ├── Ver memberships
    ├── Cambiar roles
    └── Remover usuarios

❌ NO PUEDE HACER:
├── Cambiar plan de suscripción (solo SuperAdmin)
├── Ver datos de otras empresas
├── Crear memberships (MVP - solo SuperAdmin)
└── Exceder límites del plan
```

---

### **6.2 Recruiter (Level 30)**

**Enfoque en Reclutamiento**

```
✅ PUEDE HACER:
├── Jobs
│   ├── Ver, crear, editar
│   ├── Publicar/cerrar jobs
│   └── Gestionar jobs asignados
├── Candidates
│   ├── Ver, crear, editar
│   └── Gestionar base de datos de candidatos
├── Applications
│   ├── Ver aplicaciones
│   ├── Cambiar stages
│   ├── Calificar candidatos
│   └── Agregar notas
├── Users
│   └── Ver team members (solo lectura)
└── Memberships
    └── Ver memberships (solo lectura)

❌ NO PUEDE HACER:
├── Editar configuración de empresa
├── Crear/editar/eliminar usuarios
├── Cambiar roles de memberships
├── Ver billing
└── Eliminar jobs creados por admin
```

---

### **6.3 Hiring Manager (Level 20)**

**Acceso Limitado a Jobs Asignados**

```
✅ PUEDE HACER:
├── Jobs
│   ├── Ver jobs asignados a él
│   └── Editar solo sus jobs
├── Candidates
│   └── Ver candidatos de sus jobs
├── Applications
│   ├── Ver aplicaciones de sus jobs
│   ├── Calificar candidatos
│   ├── Agregar notas
│   └── Cambiar stage (con permisos)
└── Users
    └── Ver team members

❌ NO PUEDE HACER:
├── Ver todos los jobs (solo asignados)
├── Crear jobs
├── Crear candidatos
├── Ver candidatos de otros jobs
├── Cambiar configuración de empresa
└── Gestionar usuarios/memberships
```

---

### **6.4 User (Level 10)**

**Solo Lectura**

```
✅ PUEDE HACER:
├── Ver jobs publicados
├── Ver candidatos (limitado)
├── Ver reportes básicos
└── Ver team members

❌ NO PUEDE HACER:
├── Crear/editar/eliminar nada
├── Cambiar stages de aplicaciones
├── Calificar candidatos
├── Ver datos sensibles
└── Gestionar configuración
```

---

## 7. Aislamiento Multi-Tenant

### **7.1 Validación en Cada Request**

**Middleware de Autenticación:**
```go
// internal/shared/middleware/auth_middleware.go

func AuthMiddleware(jwtService) gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. Extraer token del header
        token := c.GetHeader("Authorization")
        
        // 2. Validar y parsear JWT
        claims, err := jwtService.ValidateToken(token)
        
        // 3. Inyectar datos en contexto
        c.Set("user_id", claims.UserID)
        c.Set("company_id", claims.CompanyID)  // ✅ Empresa del usuario
        c.Set("email", claims.Email)
        c.Set("role", claims.Role)
        
        c.Next()
    }
}
```

**Validación en Handlers:**
```go
func (h *Handler) GetResources(c *gin.Context) {
    companyID, exists := c.Get("company_id")
    if !exists {
        return 403 "No company context"
    }
    
    // ✅ Filtrar por empresa del token
    resources := h.service.GetByCompanyID(companyID)
    return resources
}
```

---

### **7.2 Queries con Filtrado Automático**

**Patrón en Repositories:**
```go
// internal/app/repositories/job_repository.go

func (r *jobRepository) GetJobsByCompanyID(companyID uint) ([]models.Job, error) {
    var jobs []models.Job
    
    // ✅ WHERE company_id = ? (aislamiento)
    err := r.db.Where("company_id = ?", companyID).
              Order("created_at DESC").
              Find(&jobs).Error
    
    return jobs, err
}
```

**Validación en Operaciones Individuales:**
```go
func (s *JobService) UpdateJob(jobID uint, dto UpdateJobDTO, companyID uint) error {
    // 1. Obtener job
    job, err := s.repo.FindByID(jobID)
    
    // 2. ✅ Validar que pertenece a la empresa
    if job.CompanyID != companyID {
        return errors.New("Access denied: job does not belong to your company")
    }
    
    // 3. Actualizar
    job.Title = dto.Title
    job.Description = dto.Description
    s.repo.Update(job)
}
```

---

### **7.3 Prevención de Manipulación de company_id**

**Forzar company_id del Token:**
```go
func (h *JobHandler) CreateJob(c *gin.Context) {
    var dto dtos.CreateJobDTO
    c.ShouldBindJSON(&dto)
    
    // ✅ IGNORAR cualquier company_id enviado en el body
    companyID, _ := c.Get("company_id")
    dto.CompanyID = companyID  // ← Forzar del token
    
    // Ahora es seguro crear
    job := h.service.CreateJob(dto)
}
```

**Intento de Ataque (Bloqueado):**
```bash
# ❌ Atacante intenta crear job para otra empresa
POST /api/v1/jobs
Authorization: Bearer {token_company_5}
{
  "title": "Hacker Job",
  "company_id": 999  # ← Intento de manipulación
}

# ✅ Sistema ignora company_id=999 y fuerza company_id=5 del token
# Job creado con company_id=5 (seguro)
```

---

## 8. Flujo Completo de Uso

### **Caso 1: Onboarding de Empresa Nueva**

```bash
# Día 1: Registro
POST /api/v1/auth/register-company
{
  "company_name": "TechStartup",
  "company_slug": "techstartup",
  "admin_email": "ceo@techstartup.com",
  "admin_password": "SecurePass123!",
  "admin_first_name": "Laura",
  "admin_last_name": "Martínez"
}

# ✅ Sistema crea:
# - Empresa (TechStartup, plan=free, trial 1 mes)
# - Usuario Admin (Laura)
# - Membership (Laura → TechStartup, role=admin)
# - Token JWT con company_id

# Día 1: Laura crea primer job
POST /api/v1/jobs
Authorization: Bearer {token}
{
  "title": "Senior Developer",
  "description": "...",
  "status": "published"
}

# Día 2: Laura invita recruiter
POST /api/v1/users
{
  "email": "recruiter@techstartup.com",
  "password": "TempPass123!",
  "first_name": "Carlos",
  "last_name": "Gómez",
  "role": "recruiter"
}

# Día 3: Carlos (recruiter) hace login
POST /api/v1/auth/login
{
  "email": "recruiter@techstartup.com",
  "password": "TempPass123!"
}

# Día 3: Carlos crea candidato
POST /api/v1/candidates
Authorization: Bearer {carlos_token}
{
  "email": "candidate@example.com",
  "first_name": "Ana",
  "last_name": "Rodríguez"
}

# Día 4: Carlos crea aplicación
POST /api/v1/applications
{
  "job_id": 10,
  "candidate_id": 50
}

# Día 10: Laura revisa candidatos
GET /api/v1/applications?job_id=10

# Día 15: Carlos mueve candidato a technical
PUT /api/v1/applications/100
{
  "stage": "technical",
  "notes": "Pasó screening, agendar entrevista técnica"
}

# Día 20: Laura hace oferta
PUT /api/v1/applications/100
{
  "stage": "offer",
  "rating": 5
}

# Día 25: Candidato acepta
PUT /api/v1/applications/100
{
  "stage": "hired"
}

# ✅ Proceso completo de hiring completado
```

---

### **Caso 2: Usuario Freelance en Múltiples Empresas**

```bash
# Ana trabaja para dos empresas

# Login inicial
POST /api/v1/auth/login
{
  "email": "ana@freelance.com",
  "password": "Pass123!"
}

# Response: Tiene 2 empresas
{
  "companies": [
    {"id": 5, "name": "CompanyA", "role": "recruiter"},
    {"id": 12, "name": "CompanyB", "role": "admin"}
  ],
  "requires_company_selection": true
}

# Ana selecciona CompanyA
POST /api/v1/auth/switch-company
{
  "company_id": 5
}

# Token nuevo con company_id=5, role=recruiter

# Ana trabaja en CompanyA
GET /api/v1/jobs
# Ve solo jobs de CompanyA

# Después Ana quiere cambiar a CompanyB
POST /api/v1/auth/switch-company
{
  "company_id": 12
}

# Token nuevo con company_id=12, role=admin

# Ahora Ana ve datos de CompanyB
GET /api/v1/jobs
# Ve solo jobs de CompanyB

# Ana tiene más permisos aquí (es admin)
POST /api/v1/users
# Puede crear usuarios en CompanyB
```

---

### **Caso 3: Alcanzar Límites del Plan**

```bash
# Empresa en plan Free (max 3 jobs)

# Admin crea 3 jobs
POST /api/v1/jobs  # Job 1 ✅
POST /api/v1/jobs  # Job 2 ✅
POST /api/v1/jobs  # Job 3 ✅

# Intenta crear 4to job
POST /api/v1/jobs  # Job 4 ❌

# Response:
{
  "error": "Job limit reached",
  "message": "Your Free plan allows maximum 3 active jobs. Upgrade to create more.",
  "current_limit": 3,
  "current_usage": 3,
  "upgrade_url": "/billing/upgrade"
}

# Admin debe:
# 1. Cerrar/eliminar un job existente, O
# 2. Upgradear plan (contactar SuperAdmin en MVP)
```

---

## 9. Casos de Uso Reales

### **Caso 1: Pipeline Completo de Hiring**

```
1. Admin crea job "Backend Developer" (status=published)
2. Recruiter busca candidatos en LinkedIn
3. Recruiter crea 10 candidatos en el sistema
4. Recruiter crea aplicaciones para el job
5. Applications en stage "applied"
6. Recruiter hace screening telefónico
7. Cambia stage → "screening" para 5 mejores
8. Asigna Hiring Manager al job
9. Hiring Manager revisa y califica (rating 1-5)
10. Mejores 2 → stage "technical"
11. Tech lead hace entrevista técnica
12. 1 candidato destaca → stage "offer"
13. Admin prepara oferta salarial
14. Candidato acepta → stage "hired"
15. Job se cierra (status=closed)
```

---

### **Caso 2: Gestión de Equipo**

```bash
# Admin crea estructura de equipo

# 1. Crear recruiters
POST /api/v1/users
{"email": "recruiter1@...", "role": "recruiter"}
POST /api/v1/users
{"email": "recruiter2@...", "role": "recruiter"}

# 2. Crear hiring managers
POST /api/v1/users
{"email": "tech-lead@...", "role": "hiring_manager"}
POST /api/v1/users
{"email": "cto@...", "role": "hiring_manager"}

# 3. Asignar jobs específicos
PUT /api/v1/jobs/10
{
  "assigned_recruiter_id": 30,
  "hiring_manager_id": 35
}

# 4. Recruiter 1 trabaja solo en sus jobs asignados
# 5. Tech Lead solo ve candidatos de sus jobs
```

---

### **Caso 3: Exportar Datos**

```bash
# Feature: exportar lista de candidatos (si plan lo permite)

GET /api/v1/candidates/export?format=csv
Authorization: Bearer {token}

# Validación:
# 1. Verificar plan permite exportar (can_export_data=true)
# 2. Filtrar solo candidatos de la empresa
# 3. Generar CSV con datos

# Response:
# - Plan Free: 403 "Upgrade to export data"
# - Plan Starter+: CSV descargable
```

---

## 10. Restricciones y Límites

### **10.1 Límites por Plan**

**Free Plan:**
```
✅ Límites:
├── Max Users: 2
├── Max Jobs: 3
├── Max Candidates: 50
├── Max Applications: 100
├── Storage: 1 GB
└── Support: Email

❌ Restricciones:
├── NO exportar datos
├── NO custom branding
├── NO API access
└── NO integrations
```

**Starter Plan ($29.99/mes):**
```
✅ Límites:
├── Max Users: 5
├── Max Jobs: 10
├── Max Candidates: 200
├── Max Applications: 500
├── Storage: 5 GB
├── Trial: 14 días
└── Support: Email

✅ Features:
└── Exportar datos ✅
```

**Professional Plan ($89.99/mes):**
```
✅ Límites:
├── Max Users: 15
├── Max Jobs: 50
├── Max Candidates: 1000
├── Max Applications: 5000
├── Storage: 20 GB
└── Support: Priority

✅ Features:
├── Exportar datos ✅
├── Custom branding ✅
├── API access ✅
└── Integrations ✅
```

**Enterprise Plan ($149.99/mes):**
```
✅ TODO ILIMITADO:
├── Users: -1 (unlimited)
├── Jobs: -1
├── Candidates: -1
├── Applications: -1
├── Storage: -1
└── Support: Dedicated

✅ Features:
├── Todo de Professional +
└── Soporte dedicado
```

---

### **10.2 Validación de Límites (Ejemplo)**

```go
// internal/app/services/job_service.go

func (s *JobService) CreateJob(dto CreateJobDTO) (*models.Job, error) {
    // 1. Obtener empresa
    company, _ := s.companyRepo.FindByID(dto.CompanyID)
    
    // 2. Obtener plan
    plan, _ := s.planRepo.FindBySlug(company.PlanTier)
    
    // 3. Contar jobs activos
    activeJobsCount, _ := s.jobRepo.CountActiveByCompany(company.ID)
    
    // 4. ✅ Validar límite
    if !plan.IsUnlimited("max_jobs") && activeJobsCount >= plan.MaxJobs {
        return nil, errors.New("Job limit reached. Upgrade your plan to create more jobs.")
    }
    
    // 5. Crear job
    job := models.Job{
        Title:       dto.Title,
        CompanyID:   dto.CompanyID,
        Status:      dto.Status,
    }
    s.jobRepo.Create(&job)
    
    return &job, nil
}
```

---

### **10.3 Validación de Trial**

```go
// Middleware futuro (Fase 2)
func CheckTrialExpired() gin.HandlerFunc {
    return func(c *gin.Context) {
        companyID, _ := c.Get("company_id")
        company, _ := companyRepo.FindByID(companyID)
        
        // ✅ Verificar si trial expiró
        if company.TrialEndsAt != nil && time.Now().After(*company.TrialEndsAt) {
            // Opciones:
            // A) Auto-downgrade a Free
            // B) Bloquear escritura (solo lectura)
            // C) Mostrar banner pero permitir uso
            
            c.JSON(403, gin.H{
                "error": "Trial expired",
                "message": "Your trial has ended. Please upgrade to continue.",
                "trial_ended": company.TrialEndsAt,
            })
            c.Abort()
            return
        }
        
        c.Next()
    }
}
```

---

## 📊 Resumen Ejecutivo

### **Cliente en Resumen**

```
┌─────────────────────────────────────────────────────────────┐
│  CLIENTE = Usuario con Contexto de Empresa                  │
├─────────────────────────────────────────────────────────────┤
│  ✅ TIENE: company_id en JWT                                │
│  ✅ VE: Solo datos de su empresa                            │
│  ✅ PUEDE: Crear/editar recursos de su empresa              │
│  ✅ ROLES: admin, recruiter, hiring_manager, user           │
│  ✅ MULTI-EMPRESA: Switch entre empresas                    │
│  ✅ LÍMITES: Según plan (free, starter, pro, enterprise)    │
│                                                              │
│  ❌ NO VE: Datos de otras empresas                          │
│  ❌ NO PUEDE: Cambiar plan, ver otras empresas              │
│  ❌ NO PUEDE: Crear memberships (MVP)                       │
└─────────────────────────────────────────────────────────────┘
```

### **Flujo Típico:**

```bash
1. Registro → Empresa + Admin creados
2. Login → Token con company_id
3. Crear usuarios → Team building
4. Crear jobs → Publicar posiciones
5. Agregar candidatos → Database
6. Crear aplicaciones → Tracking
7. Mover pipeline → Stages
8. Contratar → ¡Success!
```

### **Seguridad Multi-Tenant:**

- ✅ company_id en TODAS las queries
- ✅ Validación en handlers
- ✅ Forzar company_id del token (no del body)
- ✅ Soft deletes (mantener historial)
- ✅ Aislamiento completo entre empresas

---

**Documento generado automáticamente**  
**Basado en código real de:** `/home/ramosmg/go/src/dvra-api`  
**Fecha:** 9 de Diciembre, 2025  
**Versión API:** v1.2.0
