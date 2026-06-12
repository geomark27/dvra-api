# 📋 Módulo de Planes (Plans) - Documentación Completa

## 🎯 Descripción General

Módulo completo de gestión de planes de suscripción para Dvra ATS. Permite al SuperAdmin crear, editar, activar/desactivar planes y asignarlos a empresas. Los clientes pueden ver los planes públicos disponibles.

---

## 🏗️ Arquitectura del Módulo

### **Capas Implementadas**

```
┌─────────────────────────────────────────────────────────────┐
│                     PLAN MODULE                              │
├─────────────────────────────────────────────────────────────┤
│                                                              │
│  1. MODEL (models/plan.go)                                  │
│     • Plan struct con 25+ campos                            │
│     • Métodos: IsUnlimited(), HasFeature()                  │
│                                                              │
│  2. DTO (dtos/plan_dto.go)                                  │
│     • PlanDTO (create)                                      │
│     • UpdatePlanDTO (partial update)                        │
│     • PlanResponse                                          │
│     • AssignPlanToCompanyDTO                                │
│     • TogglePlanStatusDTO                                   │
│                                                              │
│  3. REPOSITORY (repositories/plan_repository.go)            │
│     • CRUD completo + queries especiales                    │
│     • FindActive(), FindPublic(), ExistsBySlug()            │
│                                                              │
│  4. SERVICE (services/plan_service.go)                      │
│     • Lógica de negocio                                     │
│     • Validaciones                                          │
│     • AssignPlanToCompany()                                 │
│                                                              │
│  5. HANDLER (handlers/plan_handler.go)                      │
│     • 9 endpoints HTTP                                      │
│     • Validación de permisos (SuperAdmin)                   │
│                                                              │
│  6. ROUTES (platform/server/routes.go)                      │
│     • /api/v1/plans (público)                               │
│     • /api/v1/admin/plans (SuperAdmin)                      │
│                                                              │
│  7. SEEDER (database/seeders/plan_seeder.go)                │
│     • 4 planes predefinidos (Free, Starter, Pro, Enterprise)│
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

---

## 📊 Planes Predefinidos

### **1. Free Plan**
```json
{
  "name": "Free",
  "slug": "free",
  "price": 0.00,
  "max_users": 2,
  "max_jobs": 3,
  "max_candidates": 50,
  "max_applications": 100,
  "max_storage_gb": 1,
  "can_export_data": false,
  "can_use_custom_brand": false,
  "can_use_api": false,
  "can_use_integrations": false,
  "support_level": "email"
}
```

### **2. Starter Plan**
```json
{
  "name": "Starter",
  "slug": "starter",
  "price": 39.99,
  "trial_days": 14,
  "max_users": 5,
  "max_jobs": 10,
  "max_candidates": 200,
  "max_applications": 500,
  "max_storage_gb": 5,
  "can_export_data": true,
  "support_level": "email"
}
```

### **3. Professional Plan**
```json
{
  "name": "Professional",
  "slug": "professional",
  "price": 79.99,
  "trial_days": 14,
  "max_users": 15,
  "max_jobs": 50,
  "max_candidates": 1000,
  "max_applications": 5000,
  "max_storage_gb": 20,
  "can_export_data": true,
  "can_use_custom_brand": true,
  "can_use_api": true,
  "can_use_integrations": true,
  "support_level": "priority"
}
```

### **4. Enterprise Plan**
```json
{
  "name": "Enterprise",
  "slug": "enterprise",
  "price": 159.99,
  "trial_days": 30,
  "max_users": -1,           // Unlimited
  "max_jobs": -1,            // Unlimited
  "max_candidates": -1,      // Unlimited
  "max_applications": -1,    // Unlimited
  "max_storage_gb": -1,      // Unlimited
  "can_export_data": true,
  "can_use_custom_brand": true,
  "can_use_api": true,
  "can_use_integrations": true,
  "support_level": "dedicated"
}
```

---

## 🔌 Endpoints Implementados

### **RUTAS PÚBLICAS (Sin Autenticación)**

#### 1. Obtener Planes Públicos (Pricing Page)
```http
GET /api/v1/plans
```
**Descripción:** Devuelve todos los planes activos y públicos para mostrar en la página de precios.

**Response:**
```json
{
  "status": "success",
  "data": {
    "plans": [...],
    "count": 4
  }
}
```

#### 2. Obtener Plan por Slug
```http
GET /api/v1/plans/{slug}
```
**Ejemplo:** `GET /api/v1/plans/professional`

---

### **RUTAS DE SUPERADMIN** (Requieren `role=superadmin`)

#### 3. Crear Plan
```http
POST /api/v1/admin/plans
Authorization: Bearer {superadmin_token}
Content-Type: application/json

{
  "name": "Custom Plan",
  "slug": "custom",
  "description": "Plan personalizado para necesidades especiales",
  "price": 149.99,
  "currency": "USD",
  "billing_cycle": "monthly",
  "is_active": true,
  "is_public": true,
  "trial_days": 7,
  "display_order": 5,
  "max_users": 10,
  "max_jobs": 25,
  "max_candidates": 500,
  "max_applications": 2000,
  "max_storage_gb": 10,
  "can_export_data": true,
  "can_use_custom_brand": false,
  "can_use_api": true,
  "can_use_integrations": false,
  "support_level": "priority"
}
```

**Response (201):**
```json
{
  "status": "success",
  "data": {
    "plan": {
      "id": 5,
      "name": "Custom Plan",
      "slug": "custom",
      ...
    }
  }
}
```

#### 4. Obtener Todos los Planes
```http
GET /api/v1/admin/plans
Authorization: Bearer {superadmin_token}
```
**Descripción:** SuperAdmin ve todos los planes (activos e inactivos, públicos y privados).

#### 5. Obtener Plan por ID
```http
GET /api/v1/admin/plans/{id}
Authorization: Bearer {superadmin_token}
```

#### 6. Actualizar Plan
```http
PUT /api/v1/admin/plans/{id}
Authorization: Bearer {superadmin_token}
Content-Type: application/json

{
  "price": 119.99,
  "max_users": 20,
  "can_use_integrations": true
}
```
**Nota:** Actualización parcial, solo envía los campos a modificar.

#### 7. Habilitar/Deshabilitar Plan
```http
PATCH /api/v1/admin/plans/{id}/toggle
Authorization: Bearer {superadmin_token}
Content-Type: application/json

{
  "is_active": false
}
```

**Response:**
```json
{
  "status": "success",
  "message": "Plan disabled successfully",
  "data": {
    "plan": {...}
  }
}
```

#### 8. Eliminar Plan
```http
DELETE /api/v1/admin/plans/{id}
Authorization: Bearer {superadmin_token}
```

**Validación:**
- ❌ No se puede eliminar si hay empresas usando el plan
- ✅ Soft delete (no borra físicamente)

**Errores:**
- `409 Conflict` - Plan en uso por empresas

#### 9. Asignar Plan a Empresa
```http
POST /api/v1/admin/plans/assign
Authorization: Bearer {superadmin_token}
Content-Type: application/json

{
  "company_id": 5,
  "plan_id": 3
}
```

**Response:**
```json
{
  "status": "success",
  "message": "Plan assigned to company successfully"
}
```

---

## 🔒 Control de Acceso

### **Matriz de Permisos**

| Endpoint | SuperAdmin | Cliente |
|----------|-----------|---------|
| `GET /plans` | ✅ Todos | ✅ Públicos |
| `GET /plans/{slug}` | ✅ Sí | ✅ Sí |
| `GET /admin/plans` | ✅ Todos | ❌ 403 |
| `POST /admin/plans` | ✅ Sí | ❌ 403 |
| `PUT /admin/plans/{id}` | ✅ Sí | ❌ 403 |
| `PATCH /admin/plans/{id}/toggle` | ✅ Sí | ❌ 403 |
| `DELETE /admin/plans/{id}` | ✅ Sí | ❌ 403 |
| `POST /admin/plans/assign` | ✅ Sí | ❌ 403 |

### **Validación en Middleware**

```go
// Todas las rutas /admin/* pasan por:
admin.Use(middleware.AuthMiddleware(jwtService))
admin.Use(middleware.RequireSuperAdmin())

// RequireSuperAdmin valida:
role, _ := c.Get("role")
companyID, _ := c.Get("company_id")

if role != "superadmin" || companyID != nil {
    return 403 Forbidden
}
```

---

## 📦 Modelo de Datos

### **Tabla: plans**

| Campo | Tipo | Descripción |
|-------|------|-------------|
| `id` | uint | Primary key |
| `name` | string | Nombre del plan |
| `slug` | string | Identificador único (URL-friendly) |
| `description` | text | Descripción detallada |
| `price` | decimal(10,2) | Precio mensual/anual |
| `currency` | varchar(3) | Moneda (USD, EUR, etc.) |
| `billing_cycle` | varchar(20) | `monthly` o `yearly` |
| `is_active` | boolean | Plan activo/inactivo |
| `is_public` | boolean | Visible en página de precios |
| `trial_days` | int | Días de prueba gratis |
| `display_order` | int | Orden de visualización |
| `max_users` | int | Límite de usuarios (-1 = ilimitado) |
| `max_jobs` | int | Límite de vacantes (-1 = ilimitado) |
| `max_candidates` | int | Límite de candidatos (-1 = ilimitado) |
| `max_applications` | int | Límite de aplicaciones (-1 = ilimitado) |
| `max_storage_gb` | int | Almacenamiento en GB (-1 = ilimitado) |
| `can_export_data` | boolean | Puede exportar datos |
| `can_use_custom_brand` | boolean | Puede personalizar marca |
| `can_use_api` | boolean | Acceso a API |
| `can_use_integrations` | boolean | Integraciones de terceros |
| `support_level` | varchar(50) | `email`, `priority`, `dedicated` |
| `created_at` | timestamp | Fecha de creación |
| `updated_at` | timestamp | Última modificación |
| `deleted_at` | timestamp | Soft delete |

---

## 🧪 Casos de Uso

### **Caso 1: Landing Page - Mostrar Precios**

```javascript
// Frontend (sin autenticación)
fetch('http://localhost:8001/api/v1/plans')
  .then(res => res.json())
  .then(data => {
    // Renderizar pricing cards
    data.data.plans.forEach(plan => {
      renderPlanCard(plan);
    });
  });
```

### **Caso 2: SuperAdmin - Crear Plan Custom**

```bash
# 1. Login como SuperAdmin
POST /auth/superadmin/login
{
  "email": "superadmin@dvra.com",
  "password": "SuperSecret123!"
}

Response: { "access_token": "eyJ..." }

# 2. Crear plan
POST /admin/plans
Authorization: Bearer eyJ...
{
  "name": "Agency Plan",
  "slug": "agency",
  "price": 199.99,
  ...
}
```

### **Caso 3: SuperAdmin - Asignar Plan a Empresa**

```bash
# Cambiar empresa de Free a Professional
POST /admin/plans/assign
{
  "company_id": 5,
  "plan_id": 3
}

# La empresa #5 ahora tiene plan_tier = "professional"
```

### **Caso 4: SuperAdmin - Deshabilitar Plan**

```bash
# Ocultar un plan temporalmente
PATCH /admin/plans/2/toggle
{
  "is_active": false
}

# El plan Starter ya no aparece en GET /plans
# Pero las empresas que lo tienen siguen usándolo
```

---

## 🚀 Inicialización del Sistema

### **1. Migración Automática**

```go
// internal/database/models_all.go
var AllModels = []interface{}{
    ...,
    &models.Plan{},  // ← Agregado
}

// Al iniciar el servidor:
db.AutoMigrate(AllModels...)
```

### **2. Seeder de Planes**

```bash
# Ejecutar seeders
go run cmd/console/main.go seed

# Output:
🌱 Running seeders...
✅ Roles seeded
✅ Plans seeded (4 planes creados)
✅ Users seeded
✅ Companies seeded
```

**Planes creados automáticamente** (fuente de verdad: `internal/database/seeders/plan_seeder.go`):
1. Free ($0)
2. Starter ($39.99)
3. Professional ($79.99)
4. Enterprise ($159.99)

---

## 🛡️ Validaciones Implementadas

### **En Service Layer**

1. **CreatePlan:**
   - ✅ Slug único (no duplicados)
   - ✅ Precio >= 0
   - ✅ Moneda válida (3 caracteres)
   - ✅ Billing cycle válido (`monthly` o `yearly`)
   - ✅ Support level válido (`email`, `priority`, `dedicated`)

2. **UpdatePlan:**
   - ✅ Plan existe
   - ✅ Actualización parcial (solo campos enviados)

3. **DeletePlan:**
   - ✅ Plan existe
   - ❌ No se puede eliminar si está en uso
   - ✅ Soft delete (recuperable)

4. **AssignPlanToCompany:**
   - ✅ Plan existe y está activo (usa `FindActiveBySlug()`)
   - ✅ Empresa existe
   - ✅ Actualiza `company.plan_tier` con el slug del plan

5. **Validación en Creación de Empresas:**
   - ✅ `RegisterCompany()` valida que el plan "free" exista y esté activo antes de asignar
   - ✅ `CreateCompanyWithAdmin()` (SuperAdmin) valida plan elegido con `FindActiveBySlug()`
   - ✅ Si plan no existe o está inactivo: error "plan is not available"
   - ✅ NO usa strings hardcodeados, siempre valida contra la base de datos

### **En Handler Layer**

```go
// Validación automática con binding tags
type PlanDTO struct {
    Name      string  `json:"name" binding:"required"`
    Price     float64 `json:"price" binding:"required,min=0"`
    Currency  string  `json:"currency" binding:"required,len=3"`
    // ...
}
```

---

## 📝 Ejemplos de Respuestas

### **Success Response**
```json
{
  "status": "success",
  "data": {
    "plan": {
      "id": 1,
      "name": "Free",
      "slug": "free",
      "description": "Perfect for trying out Dvra ATS",
      "price": 0,
      "currency": "USD",
      "billing_cycle": "monthly",
      "is_active": true,
      "is_public": true,
      "trial_days": 0,
      "max_users": 2,
      "max_jobs": 3,
      "can_export_data": false,
      "support_level": "email",
      "created_at": "2025-12-09T10:30:00Z",
      "updated_at": "2025-12-09T10:30:00Z"
    }
  }
}
```

### **Error Responses**

**Plan no encontrado (404):**
```json
{
  "error": "Plan not found"
}
```

**Slug duplicado (409):**
```json
{
  "error": "Plan with this slug already exists"
}
```

**Plan en uso (409):**
```json
{
  "error": "Cannot delete plan that is in use by companies"
}
```

**Sin permisos (403):**
```json
{
  "error": "superadmin routes require no company context"
}
```

---

## 🔄 Integración con Companies

### **Relación Company → Plan**

```go
// models/company.go
type Company struct {
    ...
    PlanTier string `json:"plan_tier"`  // "free", "starter", "professional", "enterprise"
}
```

### **Flujo de Asignación**

```
┌─────────────────────────────────────────────────────────────┐
│  SuperAdmin asigna plan a empresa                            │
└─────────────────────────────────────────────────────────────┘

POST /admin/plans/assign
{
  "company_id": 5,
  "plan_id": 3
}

Backend:
1. Buscar plan por ID → plan.slug = "professional"
2. Validar que plan.is_active = true
3. Buscar empresa por ID
4. Actualizar company.plan_tier = "professional"
5. Guardar cambios

✅ Empresa actualizada

┌─────────────────────────────────────────────────────────────┐
│  Cliente ve sus límites del plan                            │
└─────────────────────────────────────────────────────────────┘

GET /my-company

Response:
{
  "company": {
    "id": 5,
    "name": "Azentic Sys",
    "plan_tier": "professional"
  }
}

Frontend consulta:
GET /plans/professional

Response:
{
  "max_users": 15,
  "max_jobs": 50,
  "can_use_api": true,
  ...
}
```

---

## 🎨 Frontend Integration Guide

### **Pricing Page Component**

```javascript
// components/PricingPage.jsx
import { useState, useEffect } from 'react';

function PricingPage() {
  const [plans, setPlans] = useState([]);

  useEffect(() => {
    fetch('http://localhost:8001/api/v1/plans')
      .then(res => res.json())
      .then(data => setPlans(data.data.plans));
  }, []);

  return (
    <div className="pricing-grid">
      {plans.map(plan => (
        <PlanCard key={plan.id} plan={plan} />
      ))}
    </div>
  );
}

function PlanCard({ plan }) {
  return (
    <div className="plan-card">
      <h2>{plan.name}</h2>
      <p className="price">${plan.price}/{plan.billing_cycle}</p>
      <p>{plan.description}</p>
      <ul>
        <li>{plan.max_users} usuarios</li>
        <li>{plan.max_jobs} vacantes</li>
        <li>{plan.max_candidates} candidatos</li>
        {plan.can_export_data && <li>✅ Exportar datos</li>}
        {plan.can_use_api && <li>✅ Acceso API</li>}
      </ul>
      <button>Seleccionar Plan</button>
    </div>
  );
}
```

### **SuperAdmin Dashboard**

```javascript
// components/admin/PlansManagement.jsx
function PlansManagement() {
  const [plans, setPlans] = useState([]);
  const token = localStorage.getItem('superadmin_token');

  useEffect(() => {
    fetch('http://localhost:8001/api/v1/admin/plans', {
      headers: { 'Authorization': `Bearer ${token}` }
    })
      .then(res => res.json())
      .then(data => setPlans(data.data.plans));
  }, []);

  const togglePlan = async (planId, isActive) => {
    await fetch(`http://localhost:8001/api/v1/admin/plans/${planId}/toggle`, {
      method: 'PATCH',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({ is_active: !isActive })
    });
    // Refresh list
  };

  return (
    <table>
      <thead>
        <tr>
          <th>Plan</th>
          <th>Precio</th>
          <th>Estado</th>
          <th>Acciones</th>
        </tr>
      </thead>
      <tbody>
        {plans.map(plan => (
          <tr key={plan.id}>
            <td>{plan.name}</td>
            <td>${plan.price}</td>
            <td>{plan.is_active ? '✅ Activo' : '❌ Inactivo'}</td>
            <td>
              <button onClick={() => togglePlan(plan.id, plan.is_active)}>
                Toggle
              </button>
              <button onClick={() => editPlan(plan.id)}>Editar</button>
            </td>
          </tr>
        ))}
      </tbody>
    </table>
  );
}
```

---

## ✅ Checklist de Implementación

- ✅ Modelo `Plan` con 25+ campos
- ✅ DTOs completos (Create, Update, Response, Assign, Toggle)
- ✅ Repository con 9 métodos
- ✅ Service con lógica de negocio completa
- ✅ Handler con 9 endpoints
- ✅ Rutas públicas y de SuperAdmin
- ✅ Middleware de autorización
- ✅ Seeder con 4 planes predefinidos
- ✅ Migración de base de datos
- ✅ Compilación exitosa
- ✅ Servidor inicia sin errores
- ✅ Documentación completa

---

## 🚀 Próximos Pasos Sugeridos

1. **Testing:**
   - Unit tests para service layer
   - Integration tests para endpoints
   - Test de validaciones

2. **Features Adicionales:**
   - Historial de cambios de plan por empresa
   - Billing automático con Stripe/PayPal
   - Notificaciones cuando empresa excede límites
   - Dashboard de analytics por plan

3. **Mejoras:**
   - Cache de planes públicos (Redis)
   - Versionado de planes (plan_v2, plan_v3)
   - Planes personalizados por empresa
   - Descuentos y promociones

---

## 📞 Soporte

Para más información sobre el módulo de planes, consulta:
- Código fuente: `internal/app/{models,dtos,repositories,services,handlers}/plan*`
- Seeders: `internal/database/seeders/plan_seeder.go`
- Rutas: `internal/platform/server/routes.go`

**Desarrollado para Dvra ATS** 🚀
