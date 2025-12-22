# 🔐 SuperAdmin - Flujo Completo del Sistema

> **Documento de Referencia Técnica**  
> Versión: 1.0 | Fecha: 9 de Diciembre, 2025  
> Describe en detalle qué puede y NO puede hacer el SuperAdmin

---

## 📋 Índice

1. [¿Qué es el SuperAdmin?](#qué-es-el-superadmin)
2. [Características Únicas del SuperAdmin](#características-únicas-del-superadmin)
3. [Autenticación y Acceso](#autenticación-y-acceso)
4. [Endpoints Disponibles](#endpoints-disponibles)
5. [Permisos COMPLETOS por Módulo](#permisos-completos-por-módulo)
6. [Restricciones y Limitaciones](#restricciones-y-limitaciones)
7. [Validaciones de Seguridad](#validaciones-de-seguridad)
8. [Casos de Uso Reales](#casos-de-uso-reales)

---

## 1. ¿Qué es el SuperAdmin?

### **Definición**

El **SuperAdmin** es un rol especial con acceso global a TODA la plataforma Dvra ATS. Es el único usuario que NO pertenece a ninguna empresa específica y puede gestionar todas las empresas del sistema.

### **Características Principales**

```
┌─────────────────────────────────────────────────────────────┐
│                     SUPERADMIN                               │
├─────────────────────────────────────────────────────────────┤
│  • Acceso GLOBAL (sin company_id)                           │
│  • Puede ver/editar TODAS las empresas                      │
│  • Gestiona planes de suscripción                           │
│  • Asigna usuarios a empresas                               │
│  • Suspende/activa empresas                                 │
│  • Ve analytics globales                                    │
│  • NO aparece en listados de team members                   │
└─────────────────────────────────────────────────────────────┘
```

### **Diferencias con Admin de Empresa**

| Característica | SuperAdmin | Admin (Empresa) |
|---------------|------------|-----------------|
| **company_id en JWT** | ❌ NULL | ✅ Requerido |
| **Alcance** | Global (todas las empresas) | Solo su empresa |
| **Ver datos de otras empresas** | ✅ Sí | ❌ No |
| **Gestionar planes** | ✅ Sí | ❌ No |
| **Crear empresas** | ✅ Sí | ❌ No |
| **Suspender empresas** | ✅ Sí | ❌ No |
| **Crear memberships** | ✅ Sí | ❌ No (MVP) |
| **Ver analytics globales** | ✅ Sí | ❌ Solo su empresa |

---

## 2. Características Únicas del SuperAdmin

### **2.1 Sin Contexto de Empresa**

```go
// JWT Token de SuperAdmin
{
  "user_id": 1,
  "company_id": null,        // ← SIN empresa
  "email": "admin@dvra.com",
  "role": "superadmin",
  "exp": 1733687400
}
```

**Validación en Middleware:**
```go
// internal/shared/middleware/auth_middleware.go

func RequireSuperAdmin() gin.HandlerFunc {
    return func(c *gin.Context) {
        role, _ := c.Get("role")
        companyID, hasCompany := c.Get("company_id")
        
        // ✅ Debe tener role = "superadmin"
        if role != "superadmin" {
            return 403 "SuperAdmin access required"
        }
        
        // ✅ NO debe tener company_id (acceso global)
        if hasCompany && companyID != nil {
            return 403 "SuperAdmin routes require no company context"
        }
        
        c.Next()
    }
}
```

### **2.2 Acceso a Rutas Exclusivas**

Todas las rutas bajo `/api/v1/admin/*` requieren ser SuperAdmin:

```go
// internal/platform/server/routes.go

admin := api.Group("/admin")
admin.Use(middleware.AuthMiddleware(jwtService))
admin.Use(middleware.RequireSuperAdmin())  // ← Validación estricta
{
    // Solo SuperAdmin puede acceder aquí
}
```

---

## 3. Autenticación y Acceso

### **3.1 Login de SuperAdmin**

**Endpoint Exclusivo:**
```http
POST /api/v1/auth/superadmin/login
```

**Request Body:**
```json
{
  "email": "admin@dvra.com",
  "password": "SuperSecurePassword123!"
}
```

**Response (200):**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "email": "admin@dvra.com",
    "first_name": "Super",
    "last_name": "Admin",
    "role": "superadmin",
    "company_id": null  // ← Importante: sin empresa
  }
}
```

**Validaciones Internas:**
```go
// internal/app/services/auth_service.go

func (s *AuthService) SuperAdminLogin(dto *SuperAdminLoginDTO) (*LoginResponseDTO, error) {
    // 1. Buscar user por email
    user, err := s.userRepo.FindByEmail(dto.Email)
    
    // 2. Verificar password
    bcrypt.CompareHashAndPassword(user.PasswordHash, dto.Password)
    
    // 3. Buscar membership con role = superadmin
    membership, err := s.membershipRepo.FindByUserID(user.ID)
    
    // 4. ✅ Validar que role = "superadmin"
    if membership.Role != models.RoleSuperAdmin {
        return errors.New("user is not a superadmin")
    }
    
    // 5. ✅ Validar que company_id = NULL
    if membership.CompanyID != nil {
        return errors.New("superadmin cannot have company association")
    }
    
    // 6. Generar tokens SIN company_id
    accessToken := jwt.Generate({
        user_id: user.ID,
        company_id: nil,  // ← Sin empresa
        role: "superadmin"
    })
    
    return accessToken
}
```

### **3.2 Creación del SuperAdmin**

**Seeder Automático:**
```go
// internal/database/seeders/superadmin_seeder.go

func SeedSuperAdmin(db *gorm.DB) error {
    // 1. Crear user
    superAdmin := models.User{
        Email:        "admin@dvra.com",
        PasswordHash: bcrypt.Hash("SuperSecurePassword123!"),
        FirstName:    "Super",
        LastName:     "Admin",
        IsActive:     true,
    }
    db.Create(&superAdmin)
    
    // 2. Crear membership SIN company_id
    membership := models.Membership{
        UserID:    superAdmin.ID,
        CompanyID: nil,              // ← NULL
        Role:      "superadmin",     // ← Rol 100
        Status:    "active",
        IsDefault: true,
    }
    db.Create(&membership)
    
    return nil
}
```

**Ejecutar Seeder:**
```bash
go run cmd/console/main.go seed
# ✅ SuperAdmin creado: admin@dvra.com
```

---

## 4. Endpoints Disponibles

### **4.1 Gestión de Empresas (Companies)**

#### **📋 Listar Todas las Empresas**
```http
GET /api/v1/admin/companies
Authorization: Bearer {superadmin_token}

Query Params:
  - page (int): Página actual (default: 1)
  - limit (int): Items por página (default: 20)
  - search (string): Buscar por nombre o slug
  - plan_tier (string): Filtrar por plan (free, starter, professional, enterprise)
```

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
      "created_at": "2025-01-15T10:00:00Z",
      "trial_ends_at": "2025-02-15T10:00:00Z"
    },
    {
      "id": 2,
      "name": "DevCorp",
      "slug": "devcorp",
      "plan_tier": "enterprise",
      "status": "active",
      "user_count": 25,
      "job_count": 50,
      "created_at": "2025-02-01T14:30:00Z",
      "trial_ends_at": null
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 45
  }
}
```

**Capacidades:**
- ✅ Ve TODAS las empresas del sistema
- ✅ Filtros avanzados (plan, búsqueda)
- ✅ Paginación eficiente
- ✅ Estadísticas por empresa (users, jobs)

---

#### **➕ Crear Empresa con Admin**
```http
POST /api/v1/admin/companies
Authorization: Bearer {superadmin_token}
Content-Type: application/json
```

**Request Body:**
```json
{
  "company_name": "TechCorp SA",
  "company_slug": "techcorp-sa",
  "admin_email": "admin@techcorp.com",
  "admin_password": "Admin123!",
  "admin_first_name": "Carlos",
  "admin_last_name": "Martínez",
  "plan_slug": "starter"  // Opcional: free, starter, professional, enterprise
}
```

**Response (201):**
```json
{
  "company": {
    "id": 46,
    "name": "TechCorp SA",
    "slug": "techcorp-sa",
    "plan_tier": "starter",
    "trial_ends_at": "2026-01-09T10:00:00Z"
  },
  "admin": {
    "id": 120,
    "email": "admin@techcorp.com",
    "first_name": "Carlos",
    "last_name": "Martínez",
    "is_active": true
  },
  "message": "Company and admin created successfully"
}
```

**Proceso Interno (Transaction):**
```go
// internal/app/services/admin/superadmin_companies_service.go

func (s *SuperAdminCompaniesService) CreateCompanyWithAdmin(dto) {
    tx := s.db.Begin()
    
    // 1. Validar que plan existe y está activo
    plan, err := s.planRepo.FindActiveBySlug(dto.PlanSlug)
    if err != nil {
        return "plan 'X' is not available or inactive"
    }
    
    // 2. Crear empresa con plan validado
    company := models.Company{
        Name:        dto.CompanyName,
        Slug:        dto.CompanySlug,
        PlanTier:    plan.Slug,        // ✅ Slug validado
        TrialEndsAt: now + 1 month,
    }
    tx.Create(&company)
    
    // 3. Crear admin user
    admin := models.User{
        Email:        dto.AdminEmail,
        PasswordHash: bcrypt.Hash(dto.AdminPassword),
        FirstName:    dto.AdminFirstName,
        LastName:     dto.AdminLastName,
        IsActive:     true,
    }
    tx.Create(&admin)
    
    // 4. Crear membership (admin → empresa)
    membership := models.Membership{
        UserID:    admin.ID,
        CompanyID: &company.ID,
        Role:      "admin",
        Status:    "active",
        IsDefault: true,
    }
    tx.Create(&membership)
    
    tx.Commit()
    return company, admin, nil
}
```

**Validaciones:**
- ✅ Email del admin no debe existir
- ✅ Plan debe existir y estar activo (`FindActiveBySlug()`)
- ✅ Slug de empresa único
- ✅ Transacción completa (rollback si falla)

---

#### **🔄 Cambiar Plan de Empresa**
```http
PUT /api/v1/admin/companies/:id/plan
Authorization: Bearer {superadmin_token}
Content-Type: application/json
```

**Request Body:**
```json
{
  "new_plan": "enterprise"
}
```

**Response (200):**
```json
{
  "message": "Plan updated successfully",
  "company_id": 46,
  "new_plan": "enterprise"
}
```

**Proceso:**
```go
func (s *SuperAdminCompaniesService) ChangeCompanyPlan(companyID, newPlan) {
    db.Model(&models.Company{}).
       Where("id = ?", companyID).
       Update("plan_tier", newPlan)
}
```

**Casos de Uso:**
- Upgrade manual (cliente pagó fuera del sistema)
- Downgrade por falta de pago
- Cambio a plan custom negociado
- Promociones especiales

---

#### **⛔ Suspender Empresa**
```http
POST /api/v1/admin/companies/:id/suspend
Authorization: Bearer {superadmin_token}
Content-Type: application/json
```

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
  "company_id": 5,
  "reason": "Falta de pago - 3 meses vencidos"
}
```

**Efecto:**
```go
// Cambia plan_tier a "suspended"
db.Model(&models.Company{}).
   Where("id = ?", companyID).
   Update("plan_tier", "suspended")

// Los usuarios de la empresa NO podrán hacer login
// Middleware valida plan_tier != "suspended"
```

**Validaciones en Login:**
```go
func (s *AuthService) Login(dto) {
    user := FindByEmail(dto.Email)
    membership := FindByUserID(user.ID)
    company := FindByID(membership.CompanyID)
    
    // ✅ Verificar que empresa NO esté suspendida
    if company.PlanTier == "suspended" {
        return 403 "Company is suspended. Contact support."
    }
}
```

---

#### **👥 Ver Usuarios de Empresa**
```http
GET /api/v1/admin/companies/:id/users
Authorization: Bearer {superadmin_token}
```

**Response (200):**
```json
{
  "company_id": 1,
  "users": [
    {
      "id": 5,
      "email": "admin@azentic.com",
      "first_name": "Marcos",
      "last_name": "Ramos",
      "is_active": true,
      "role": "admin"
    },
    {
      "id": 8,
      "email": "recruiter@azentic.com",
      "first_name": "Ana",
      "last_name": "Gómez",
      "is_active": true,
      "role": "recruiter"
    }
  ],
  "count": 2
}
```

**Query:**
```sql
SELECT users.* 
FROM users
JOIN memberships ON memberships.user_id = users.id
WHERE memberships.company_id = :company_id
```

---

### **4.2 Gestión de Planes (Plans)**

#### **📋 Listar Todos los Planes (Admin)**
```http
GET /api/v1/admin/plans
Authorization: Bearer {superadmin_token}
```

**Response:**
```json
{
  "status": "success",
  "data": {
    "plans": [
      {
        "id": 1,
        "name": "Free",
        "slug": "free",
        "price": 0,
        "is_active": true,
        "is_public": true,
        "max_users": 2,
        "max_jobs": 3
      },
      {
        "id": 2,
        "name": "Starter",
        "slug": "starter",
        "price": 29.99,
        "is_active": true,
        "is_public": true,
        "trial_days": 14
      }
    ],
    "count": 4
  }
}
```

**Capacidades:**
- ✅ Ve planes activos E inactivos
- ✅ Ve planes públicos y privados
- ✅ Acceso completo a configuración

---

#### **➕ Crear Plan Personalizado**
```http
POST /api/v1/admin/plans
Authorization: Bearer {superadmin_token}
Content-Type: application/json
```

**Request Body:**
```json
{
  "name": "Corporate Custom",
  "slug": "corporate-custom",
  "description": "Plan personalizado para empresa X",
  "price": 599.99,
  "currency": "USD",
  "billing_cycle": "monthly",
  "is_active": true,
  "is_public": false,
  "trial_days": 30,
  "max_users": 50,
  "max_jobs": 200,
  "max_candidates": -1,
  "max_applications": -1,
  "max_storage_gb": 100,
  "can_export_data": true,
  "can_use_custom_brand": true,
  "can_use_api": true,
  "can_use_integrations": true,
  "support_level": "dedicated"
}
```

**Response (201):**
```json
{
  "status": "success",
  "data": {
    "plan": {
      "id": 5,
      "name": "Corporate Custom",
      "slug": "corporate-custom",
      "price": 599.99
    }
  },
  "message": "Plan created successfully"
}
```

**Validaciones:**
- ✅ Slug único
- ✅ Precio >= 0
- ✅ Moneda válida (3 caracteres)
- ✅ Billing cycle: `monthly` o `yearly`
- ✅ Support level: `email`, `priority`, `dedicated`

---

#### **✏️ Actualizar Plan**
```http
PUT /api/v1/admin/plans/:id
Authorization: Bearer {superadmin_token}
```

**Request (Actualización Parcial):**
```json
{
  "price": 99.99,
  "max_users": 20,
  "is_public": true
}
```

**Capacidades:**
- ✅ Actualización parcial (solo campos enviados)
- ✅ Puede cambiar límites en cualquier momento
- ✅ Afecta a nuevas empresas, no a existentes

---

#### **🔄 Activar/Desactivar Plan**
```http
PATCH /api/v1/admin/plans/:id/toggle
Authorization: Bearer {superadmin_token}
```

**Request:**
```json
{
  "is_active": false
}
```

**Efecto:**
- ✅ Plan ya NO aparece en `/plans` (público)
- ✅ Empresas que YA lo tienen siguen usándolo
- ✅ NO se pueden asignar nuevas empresas a este plan

---

#### **🔗 Asignar Plan a Empresa**
```http
POST /api/v1/admin/plans/assign
Authorization: Bearer {superadmin_token}
```

**Request:**
```json
{
  "company_id": 5,
  "plan_id": 3
}
```

**Proceso:**
```go
func (s *PlanService) AssignPlanToCompany(dto) {
    // 1. Verificar que plan existe y está activo
    plan, err := s.planRepo.FindByID(dto.PlanID)
    if !plan.IsActive {
        return "Plan is not active"
    }
    
    // 2. Verificar que empresa existe
    company, err := s.companyRepo.FindByID(dto.CompanyID)
    
    // 3. Actualizar plan_tier de la empresa
    company.PlanTier = plan.Slug
    s.companyRepo.Update(company)
}
```

---

#### **🗑️ Eliminar Plan**
```http
DELETE /api/v1/admin/plans/:id
Authorization: Bearer {superadmin_token}
```

**Validaciones:**
```go
func (s *PlanService) DeletePlan(planID) {
    // ✅ Verificar que NO esté en uso
    var count int64
    db.Model(&models.Company{}).
       Where("plan_tier = ?", plan.Slug).
       Count(&count)
    
    if count > 0 {
        return 409 "Cannot delete plan that is in use by companies"
    }
    
    // Soft delete
    db.Delete(&models.Plan{}, planID)
}
```

---

### **4.3 Gestión de Memberships**

#### **➕ Crear Membership (Asignar Usuario a Empresa)**
```http
POST /api/v1/admin/memberships
Authorization: Bearer {superadmin_token}
Content-Type: application/json
```

**Request:**
```json
{
  "user_id": 10,
  "company_id": 2,
  "role": "recruiter"
}
```

**Response (201):**
```json
{
  "id": 25,
  "user_id": 10,
  "company_id": 2,
  "role": "recruiter",
  "status": "active"
}
```

**Proceso:**
```go
func (h *MembershipHandler) CreateMembership(c *gin.Context) {
    role, _ := c.Get("role")
    
    // ✅ SOLO SuperAdmin puede crear memberships
    if role != "superadmin" {
        return 403 "Only superadmin can assign users to companies"
    }
    
    // ✅ company_id es requerido
    if dto.CompanyID == nil {
        return 400 "company_id is required"
    }
    
    // Crear membership
    membership := models.Membership{
        UserID:    dto.UserID,
        CompanyID: dto.CompanyID,
        Role:      dto.Role,
        Status:    "active",
    }
    db.Create(&membership)
}
```

**Casos de Uso:**
- Usuario existente se une a segunda empresa
- Admin pide agregar freelancer temporalmente
- Migración de usuarios entre empresas

**Restricción MVP:**
- ❌ Clients NO pueden crear memberships
- ✅ Solo SuperAdmin tiene este poder
- 🔜 Fase 2: Sistema de invitaciones por email

---

### **4.4 Analytics Globales**

#### **📊 Ver Analytics del Sistema**
```http
GET /api/v1/admin/analytics
Authorization: Bearer {superadmin_token}
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
  "churn_rate": 0.05,
  "companies_by_plan": {
    "free": 20,
    "starter": 12,
    "professional": 8,
    "enterprise": 2
  }
}
```

**Cálculos:**
```go
func (s *SuperAdminCompaniesService) GetGlobalAnalytics() {
    // Contadores totales
    db.Model(&models.Company{}).Count(&totalCompanies)
    db.Model(&models.Company{}).
       Where("plan_tier != ?", "suspended").
       Count(&activeCompanies)
    
    db.Model(&models.User{}).Count(&totalUsers)
    db.Model(&models.Job{}).Count(&totalJobs)
    db.Model(&models.Application{}).Count(&totalApplications)
    
    // MRR (Monthly Recurring Revenue)
    planPrices := map[string]float64{
        "free":         0,
        "starter":      29.99,
        "professional": 89.99,
        "enterprise":   149.99,
    }
    
    var companies []models.Company
    db.Where("plan_tier != ?", "suspended").Find(&companies)
    
    mrr := 0.0
    for _, company := range companies {
        mrr += planPrices[company.PlanTier]
    }
    
    analytics.MonthlyRevenue = mrr
}
```

---

## 5. Permisos COMPLETOS por Módulo

### **5.1 Companies (Empresas)**

| Operación | SuperAdmin | Admin Empresa | Notas |
|-----------|------------|---------------|-------|
| **GET** /companies | ✅ Ve TODAS | ✅ Solo la suya | SuperAdmin sin filtro |
| **GET** /companies/:id | ✅ Cualquiera | ✅ Solo la suya | Validación de company_id |
| **POST** /companies | ✅ Sí | ❌ No | Solo SuperAdmin crea empresas |
| **PUT** /companies/:id | ✅ Cualquiera | ✅ Solo la suya | Admin puede editar su empresa |
| **DELETE** /companies/:id | ✅ Sí | ❌ No | Solo SuperAdmin elimina |
| **POST** /admin/companies | ✅ Sí | ❌ No | Crear empresa con admin |
| **PUT** /admin/companies/:id/plan | ✅ Sí | ❌ No | Cambiar plan manualmente |
| **POST** /admin/companies/:id/suspend | ✅ Sí | ❌ No | Suspender empresa |
| **GET** /admin/companies/:id/users | ✅ Sí | ❌ No | Ver usuarios de cualquier empresa |

**Código de Validación:**
```go
// internal/app/handlers/company_handler.go

func (h *CompanyHandler) GetCompanies(c *gin.Context) {
    role, _ := c.Get("role")
    
    if role == "superadmin" {
        // ✅ SuperAdmin ve TODAS las empresas
        companies := h.service.GetAllCompanies()
        return companies
    }
    
    // Cliente: solo su empresa
    companyID, _ := c.Get("company_id")
    company := h.service.GetCompanyByID(companyID)
    return [company]  // Array con 1 elemento
}
```

---

### **5.2 Users (Usuarios)**

| Operación | SuperAdmin | Cliente | Filtro |
|-----------|------------|---------|--------|
| **GET** /users | ✅ TODOS | ✅ Solo de mi empresa | JOIN memberships |
| **GET** /users/:id | ✅ Cualquiera | ✅ Si pertenece a mi empresa | Valida membership |
| **POST** /users | ✅ Sí | ✅ Sí (para su empresa) | Fuerza company_id |
| **PUT** /users/:id | ✅ Cualquiera | ✅ Solo de mi empresa | Valida membership |
| **DELETE** /users/:id | ✅ Cualquiera | ✅ Solo de mi empresa | Valida membership |

**Código:**
```go
func (h *UserHandler) GetUsers(c *gin.Context) {
    role, _ := c.Get("role")
    
    if role == "superadmin" {
        // ✅ Todos los usuarios del sistema
        users := h.service.GetAllUsers()
        return users
    }
    
    // Cliente: usuarios de su empresa
    companyID, _ := c.Get("company_id")
    users := h.service.GetUsersByCompanyID(companyID)
    return users
}
```

---

### **5.3 Jobs (Ofertas de Trabajo)**

| Operación | SuperAdmin | Cliente | Validación |
|-----------|------------|---------|------------|
| **GET** /jobs | ✅ TODOS | ✅ Solo de mi empresa | WHERE company_id = ? |
| **GET** /jobs/:id | ✅ Cualquiera | ✅ Solo de mi empresa | Verifica job.company_id |
| **POST** /jobs | ❌ No* | ✅ Sí | Fuerza company_id del token |
| **PUT** /jobs/:id | ✅ Cualquiera | ✅ Solo de mi empresa | Verifica job.company_id |
| **DELETE** /jobs/:id | ✅ Cualquiera | ✅ Solo de mi empresa | Verifica job.company_id |

**Nota:** SuperAdmin NO crea jobs porque NO tiene company_id. Solo puede ver/editar/eliminar.

**Código:**
```go
func (h *JobHandler) GetJobs(c *gin.Context) {
    role, _ := c.Get("role")
    
    if role == "superadmin" {
        // ✅ Todos los jobs de todas las empresas
        jobs := h.service.GetAllJobs()
        return jobs
    }
    
    // Cliente: jobs de su empresa
    companyID, _ := c.Get("company_id")
    jobs := h.service.GetJobsByCompanyID(companyID)
    return jobs
}
```

---

### **5.4 Candidates (Candidatos)**

| Operación | SuperAdmin | Cliente | Filtro |
|-----------|------------|---------|--------|
| **GET** /candidates | ✅ TODOS | ✅ Solo de mi empresa | WHERE company_id = ? |
| **GET** /candidates/:id | ✅ Cualquiera | ✅ Solo de mi empresa | Verifica candidate.company_id |
| **POST** /candidates | ❌ No* | ✅ Sí | Fuerza company_id |
| **PUT** /candidates/:id | ✅ Cualquiera | ✅ Solo de mi empresa | Valida pertenencia |
| **DELETE** /candidates/:id | ✅ Cualquiera | ✅ Solo de mi empresa | Valida pertenencia |

---

### **5.5 Applications (Aplicaciones)**

| Operación | SuperAdmin | Cliente | Filtro |
|-----------|------------|---------|--------|
| **GET** /applications | ✅ TODAS | ✅ Solo de mi empresa | WHERE company_id = ? |
| **GET** /applications/:id | ✅ Cualquiera | ✅ Solo de mi empresa | Verifica application.company_id |
| **POST** /applications | ❌ No* | ✅ Sí | Fuerza company_id |
| **PUT** /applications/:id | ✅ Cualquiera | ✅ Solo de mi empresa | Valida pertenencia |
| **DELETE** /applications/:id | ✅ Cualquiera | ✅ Solo de mi empresa | Valida pertenencia |

---

### **5.6 Memberships (Membresías)**

| Operación | SuperAdmin | Cliente | Notas |
|-----------|------------|---------|-------|
| **GET** /memberships | ✅ TODAS | ✅ Solo de mi empresa | Filtrado por company_id |
| **GET** /memberships/:id | ✅ Cualquiera | ✅ Solo de mi empresa | Valida pertenencia |
| **POST** /admin/memberships | ✅ Sí | ❌ No | **SOLO SuperAdmin (MVP)** |
| **PUT** /memberships/:id | ✅ Cualquiera | ✅ Solo de mi empresa | Cambiar role |
| **DELETE** /memberships/:id | ✅ Cualquiera | ✅ Solo de mi empresa | Remover usuario |

**Código Crítico:**
```go
func (h *MembershipHandler) CreateMembership(c *gin.Context) {
    role, _ := c.Get("role")
    
    // ✅ RESTRICCIÓN MVP: Solo SuperAdmin
    if role != "superadmin" {
        c.JSON(403, gin.H{
            "error": "Only superadmin can assign users to companies. Regular users should create new users instead.",
        })
        return
    }
    
    // SuperAdmin debe especificar company_id explícitamente
    if dto.CompanyID == nil {
        return 400 "company_id is required"
    }
    
    // Crear membership
    membership := h.service.CreateMembership(dto)
    return 201 membership
}
```

---

### **5.7 Plans (Planes)**

| Operación | SuperAdmin | Cliente | Público |
|-----------|------------|---------|---------|
| **GET** /plans | - | - | ✅ Sí (planes públicos) |
| **GET** /plans/:slug | - | - | ✅ Sí (detalle de plan) |
| **GET** /admin/plans | ✅ Sí | ❌ No | Todos los planes (activos + inactivos) |
| **POST** /admin/plans | ✅ Sí | ❌ No | Crear plan |
| **PUT** /admin/plans/:id | ✅ Sí | ❌ No | Editar plan |
| **PATCH** /admin/plans/:id/toggle | ✅ Sí | ❌ No | Activar/desactivar |
| **DELETE** /admin/plans/:id | ✅ Sí | ❌ No | Eliminar plan |
| **POST** /admin/plans/assign | ✅ Sí | ❌ No | Asignar plan a empresa |

**Rutas Públicas (Sin Auth):**
```go
// Cualquiera puede ver los planes públicos (pricing page)
plans := api.Group("/plans")
{
    plans.GET("", planHandler.GetPublicPlans)
    plans.GET("/:slug", planHandler.GetPlanBySlug)
}
```

---

## 6. Restricciones y Limitaciones

### **6.1 Lo que SuperAdmin NO Puede Hacer**

#### **❌ Crear Recursos Asociados a Empresa**

SuperAdmin NO puede crear:
- Jobs
- Candidates
- Applications

**Razón:** Estos recursos requieren `company_id`, y SuperAdmin NO tiene empresa.

**Código:**
```go
func (h *JobHandler) CreateJob(c *gin.Context) {
    companyID, exists := c.Get("company_id")
    
    if !exists || companyID == nil {
        return 400 "company_id is required to create a job"
    }
    
    // Forzar company_id del token
    dto.CompanyID = companyID
    
    job := h.service.CreateJob(dto)
}
```

**Workaround:**
Si SuperAdmin necesita crear un job para una empresa:
1. Usar endpoint de cliente `/jobs`
2. O hacer switch-company temporalmente (NO implementado para SuperAdmin)

---

#### **❌ Pertenecer a una Empresa**

SuperAdmin **NUNCA** debe tener `company_id`:

```go
// Validación en RequireSuperAdmin middleware
if hasCompany && companyID != nil {
    return 403 "SuperAdmin routes require no company context"
}
```

Si un SuperAdmin tiene `company_id`, pierde acceso global.

---

#### **❌ Hacer Login en Endpoint Normal**

SuperAdmin debe usar su endpoint exclusivo:

```bash
# ❌ INCORRECTO
POST /api/v1/auth/login

# ✅ CORRECTO
POST /api/v1/auth/superadmin/login
```

**Validación:**
```go
func (s *AuthService) Login(dto) {
    // Este endpoint es para clientes (con empresa)
    membership := FindByUserID(user.ID)
    
    if membership.CompanyID == nil {
        return 403 "SuperAdmin users must use /auth/superadmin/login"
    }
}
```

---

### **6.2 Validaciones de Seguridad Implementadas**

#### **✅ Middleware RequireSuperAdmin**

```go
// Aplicado a TODAS las rutas /admin/*
admin := api.Group("/admin")
admin.Use(middleware.RequireSuperAdmin())

// Valida:
// 1. Usuario autenticado
// 2. role = "superadmin"
// 3. company_id = NULL
```

#### **✅ Validación en Handlers**

```go
// Patrón repetido en todos los handlers
func (h *Handler) GetResources(c *gin.Context) {
    role, _ := c.Get("role")
    
    if role == "superadmin" {
        // Sin filtro: ver TODO
        return GetAll()
    }
    
    // Filtrar por empresa
    companyID, _ := c.Get("company_id")
    return GetByCompanyID(companyID)
}
```

#### **✅ Validación en Servicios**

```go
func (s *Service) DeleteCompany(companyID) {
    // Solo SuperAdmin puede eliminar empresas
    // (ya validado por middleware)
    
    // Verificaciones adicionales:
    // - Empresa existe
    // - No tiene datos críticos
    // - Soft delete
}
```

---

## 7. Casos de Uso Reales

### **Caso 1: Onboarding de Cliente Nuevo**

**Flujo:**
```bash
# 1. SuperAdmin crea empresa con admin
POST /api/v1/admin/companies
{
  "company_name": "StartupX",
  "company_slug": "startupx",
  "admin_email": "ceo@startupx.com",
  "admin_password": "SecurePass123!",
  "admin_first_name": "Juan",
  "admin_last_name": "Pérez",
  "plan_slug": "starter"
}

# 2. Sistema crea:
# - Empresa (id: 50)
# - Usuario admin (id: 150)
# - Membership (user: 150 → company: 50, role: admin)

# 3. Admin recibe email con credenciales
# 4. Admin hace login
POST /api/v1/auth/login
{
  "email": "ceo@startupx.com",
  "password": "SecurePass123!"
}

# 5. Admin empieza a usar el ATS
```

---

### **Caso 2: Cliente Solicita Upgrade**

**Flujo:**
```bash
# 1. Cliente contacta soporte: "Queremos Enterprise"
# 2. SuperAdmin recibe pago offline
# 3. SuperAdmin actualiza plan

PUT /api/v1/admin/companies/50/plan
{
  "new_plan": "enterprise"
}

# 4. Empresa #50 ahora tiene plan Enterprise
# 5. Usuarios ven nuevos límites (usuarios ilimitados, etc.)
```

---

### **Caso 3: Falta de Pago - Suspensión**

**Flujo:**
```bash
# 1. Empresa #25 no paga hace 3 meses
# 2. SuperAdmin suspende empresa

POST /api/v1/admin/companies/25/suspend
{
  "reason": "Falta de pago - 90 días vencidos"
}

# 3. Empresa marcada como "suspended"
# 4. Usuarios NO pueden hacer login
# 5. Cuando pagan, SuperAdmin reactiva:

PUT /api/v1/admin/companies/25/plan
{
  "new_plan": "professional"
}
```

---

### **Caso 4: Usuario Freelance en Múltiples Empresas**

**Flujo:**
```bash
# 1. Freelancer "Ana" trabaja para Empresa A (id: 10)
# 2. Empresa B (id: 20) contrata a Ana
# 3. SuperAdmin crea membership

POST /api/v1/admin/memberships
{
  "user_id": 75,      // Ana
  "company_id": 20,   // Empresa B
  "role": "recruiter"
}

# 4. Ana ahora tiene 2 memberships:
# - Empresa A (admin)
# - Empresa B (recruiter)

# 5. Ana puede hacer switch entre empresas
POST /api/v1/auth/switch-company
{
  "company_id": 20
}

# 6. Recibe nuevo token con company_id = 20
```

---

### **Caso 5: Crear Plan Custom para Enterprise**

**Flujo:**
```bash
# 1. Cliente Enterprise negocia plan personalizado
# 2. SuperAdmin crea plan privado

POST /api/v1/admin/plans
{
  "name": "Acme Corp Custom",
  "slug": "acme-custom",
  "price": 999.99,
  "is_public": false,     // ← No aparece en pricing page
  "max_users": 100,
  "max_jobs": 500,
  "max_candidates": -1,
  "can_use_api": true,
  "support_level": "dedicated"
}

# 3. Asignar plan a empresa
POST /api/v1/admin/plans/assign
{
  "company_id": 30,
  "plan_id": 5
}

# 4. Empresa #30 ahora usa plan custom
```

---

### **Caso 6: Migración Masiva de Usuarios**

**Flujo:**
```bash
# Empresa A (id: 5) se fusiona con Empresa B (id: 10)
# Necesitamos mover 20 usuarios de A → B

# SuperAdmin crea script:
users = [15, 16, 17, 18, 19, ...]  # IDs de usuarios

for user_id in users:
    POST /api/v1/admin/memberships
    {
        "user_id": user_id,
        "company_id": 10,
        "role": "recruiter"
    }

# Usuarios ahora están en ambas empresas
# Luego eliminar memberships de Empresa A si es necesario
```

---

### **Caso 7: Analytics para Inversionistas**

**Flujo:**
```bash
# SuperAdmin genera reporte mensual
GET /api/v1/admin/analytics

# Response:
{
  "total_companies": 150,
  "active_companies": 145,
  "monthly_revenue": 15250.00,
  "companies_by_plan": {
    "free": 80,
    "starter": 40,
    "professional": 20,
    "enterprise": 5
  },
  "total_users": 980,
  "total_jobs": 650,
  "total_applications": 12500
}

# SuperAdmin exporta a Excel para board meeting
```

---

## 8. Resumen Ejecutivo

### **✅ Lo que SuperAdmin PUEDE Hacer**

| Categoría | Capacidades |
|-----------|-------------|
| **Empresas** | Ver todas, crear, editar, suspender, cambiar plan |
| **Usuarios** | Ver todos (de todas las empresas), asignar a empresas |
| **Planes** | Crear, editar, activar/desactivar, asignar |
| **Memberships** | Crear (asignar usuarios a empresas) |
| **Analytics** | Ver métricas globales del sistema |
| **Lectura** | Acceso completo a jobs, candidates, applications de TODAS las empresas |

### **❌ Lo que SuperAdmin NO PUEDE Hacer**

| Restricción | Razón |
|-------------|-------|
| Crear jobs, candidates, applications | No tiene `company_id` |
| Pertenecer a una empresa | Rol especial sin contexto de empresa |
| Usar endpoints normales de auth | Debe usar `/auth/superadmin/login` |
| Ver memberships en listado de team | No tiene empresa |

### **🔒 Seguridad Implementada**

1. ✅ Middleware `RequireSuperAdmin()` en todas las rutas `/admin/*`
2. ✅ Validación `role = "superadmin"` en cada handler
3. ✅ Validación `company_id = NULL` para evitar contexto de empresa
4. ✅ Endpoint de login exclusivo
5. ✅ Tokens JWT sin `company_id`
6. ✅ Filtrado condicional en handlers (SuperAdmin → todos, Cliente → filtrado)

### **📊 Estadísticas de Implementación**

- **Rutas exclusivas SuperAdmin:** 13
- **Handlers con validación SuperAdmin:** 8
- **Servicios admin:** 1 (`SuperAdminCompaniesService`)
- **Middlewares:** 2 (`AuthMiddleware` + `RequireSuperAdmin`)
- **DTOs específicos:** 6

---

**Documento generado automáticamente**  
**Basado en código real de:** `/home/ramosmg/go/src/dvra-api`  
**Fecha:** 9 de Diciembre, 2025  
**Versión API:** v1.2.0
