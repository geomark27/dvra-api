# Dvra - ATS + Marketplace de Talento Tech LATAM

> **Documento Maestro de Lógica de Negocio y Arquitectura Estratégica**  
> Versión: 2.0 | Última actualización: Diciembre 8, 2025

---

## 📋 Tabla de Contenidos

1. [Resumen Ejecutivo](#resumen-ejecutivo)
2. [Modelo de Negocio Dual](#modelo-de-negocio-dual)
3. [Arquitectura Multi-Tenant](#arquitectura-multi-tenant)
4. [Sistema de Roles y Permisos](#sistema-de-roles-y-permisos)
5. [Reglas de Negocio Fundamentales](#reglas-de-negocio-fundamentales)
6. [Pipeline de Candidatos](#pipeline-de-candidatos)
7. [Red Dvra - Marketplace de Talento](#red-dvra-marketplace-de-talento)
8. [Evaluación Técnica](#evaluación-técnica)
9. [Pricing y Límites por Tier](#pricing-y-límites-por-tier)
10. [Flujos de Usuario Críticos](#flujos-de-usuario-críticos)
11. [Integraciones Estratégicas](#integraciones-estratégicas)
12. [Roadmap de Implementación](#roadmap-de-implementación)

---

## 1. Resumen Ejecutivo

### Visión del Producto

**Dvra** es una plataforma híbrida única que combina:

1. **ATS (Applicant Tracking System)**: Software SaaS para que empresas gestionen su proceso de reclutamiento interno
2. **Red Dvra (Marketplace de Talento)**: Red curada de desarrolladores LATAM pre-evaluados, accesible desde el ATS

### El Problema que Resolvemos

**Para Empresas:**
- Herramientas enterprise prohibitivamente caras ($15k-50k/año)
- Alternativas económicas genéricas sin especialización tech
- Pierden 25+ horas semanales con herramientas dispersas
- No tienen acceso a talento tech pre-evaluado de LATAM

**Para Desarrolladores:**
- LinkedIn saturado sin curación real
- No existe validación técnica objetiva de skills
- Procesos de aplicación repetitivos y manuales
- Falta de visibilidad para talento LATAM en empresas internacionales

### Diferenciador Único en el Mercado

```
┌─────────────────────────────────────────────────────────┐
│                    COMPETENCIA                          │
├─────────────────────────────────────────────────────────┤
│ Mercor/Turing     → Solo marketplace (sin herramienta)  │
│ Greenhouse/Lever  → Solo ATS (sin red de candidatos)    │
│ LinkedIn          → Tablón sin curación ni herramientas │
└─────────────────────────────────────────────────────────┘
                          ↓
┌─────────────────────────────────────────────────────────┐
│                 DVRA = ATS + MARKETPLACE                │
├─────────────────────────────────────────────────────────┤
│ 1. Empresa usa Dvra ATS ($99/mes)                       │
│ 2. Gestiona proceso de reclutamiento interno            │
│ 3. OPCIÓN: "Buscar en Red Dvra"                         │
│ 4. Acceso a candidatos pre-evaluados                    │
│ 5. Fee solo por contratación exitosa (15-20% salary)   │
└─────────────────────────────────────────────────────────┘
```

**Ventaja Competitiva Sostenible:**
- Empresas ya pagan por el ATS (ingresos recurrentes)
- Marketplace genera ingresos adicionales sin costo marginal
- Network effects: más empresas → más candidatos → más valor
- Sticky: difícil cambiar cuando tienes proceso + candidatos en una plataforma

### Diferenciadores Clave
1. **Modelo híbrido único** - ATS + Marketplace = Solución completa
2. **Pricing 80% más económico** que competencia enterprise
3. **Pre-evaluación técnica** - Candidatos validados antes de contacto
4. **Especialización LATAM** - Timezone, idioma, cultura regional
5. **Doble revenue stream** - SaaS recurrente + fees por hire

### Estado Actual vs Visión

```
IMPLEMENTADO (v0.1 - Fundación):
✅ Arquitectura multi-tenant
✅ Modelos de datos core (Company, User, Job, Candidate, Application)
✅ Sistema de roles básico (SuperAdmin, Admin, Recruiter, User)
✅ Pipeline de candidatos (stages: applied → screening → technical → offer → hired/rejected)
✅ JWT authentication service
✅ CRUD completo de entidades principales

AÑO 1 - MVP ATS (90% esfuerzo):
⏳ Autenticación y autorización completa
⏳ Portal público de aplicación
⏳ CV parsing con AWS Textract
⏳ Email notifications (SendGrid)
⏳ Google Calendar integration
⏳ Analytics dashboard básico
⏳ Tier enforcement (límites por plan)
⏳ Stripe/MercadoPago integration

AÑO 1 - Preparación Marketplace (10% esfuerzo):
⏳ Landing page talento.dvra.com
⏳ Registro de candidatos (formulario simple)
⏳ Evaluación manual de 2-3 candidatos/semana
⏳ Sistema de tags y scoring básico
⏳ Meta: 200 candidatos registrados, 50 evaluados

AÑO 2 - Integración Marketplace (40% esfuerzo):
⏳ Vista "Red Dvra" dentro del ATS
⏳ Búsqueda y filtros de candidatos
⏳ Tracking: Interest → Interview → Offer → Hired
⏳ Sistema de fees automático (15-20% salary)
⏳ Dashboard de conversión marketplace
⏳ Auto-scoring con GitHub/GitLab API

AÑO 3 - Full Híbrido (50/50 esfuerzo):
⏳ ML matching algorithm (tu expertise en grafos)
⏳ Video introductions de candidatos
⏳ Coding challenges integrados
⏳ Badges y gamification
⏳ Mobile app
⏳ API pública + webhooks
```

### Proyecciones Financieras

**Año 1 (ATS Focus):**
- 50 clientes pagantes
- MRR: $5,000 (solo SaaS)
- ARR: $60,000
- Inversión: $25k
- Break-even: Mes 9-10

**Año 2 (Integración Marketplace):**
- 150 clientes
- MRR: $15k SaaS + $3k fees = $18k
- ARR: $216,000
- Red: 500 candidatos evaluados
- 20-30 hires del marketplace

**Año 3 (Full Híbrido):**
- 400 clientes
- MRR: $50k SaaS + $25k fees = $75k
- ARR: $900,000
- Red: 2,000+ candidatos activos
- 100+ hires/año del marketplace
- LTV:CAC = 13:1

---

## 2. Modelo de Negocio Dual

### 2.1 Revenue Streams

**Stream 1: SaaS (ATS) - Recurrente**
```
Professional:  $49/mes  × 12 meses = $588/año
Business:      $149/mes × 12 meses = $1,788/año
Enterprise:    $399/mes × 12 meses = $4,788/año

Características:
- Ingresos predecibles y recurrentes
- Low churn por switching cost
- Gross margin: 70%
```

**Stream 2: Marketplace (Red Dvra) - Transaccional**
```
Fee por hire: 15-20% del annual salary o $3,500-5,000 flat

Ejemplo:
- Developer senior: $60k/año × 15% = $9,000
- Developer mid: $40k/año × 15% = $6,000
- Developer junior: $30k/año o $3,500 flat

Características:
- Ingresos variables pero alto ticket
- Zero marginal cost (candidatos ya evaluados)
- Gross margin: ~95%
- Solo se cobra cuando hay hire exitoso
```

**Blended ARPU (Año 3):**
```
Cliente promedio:
- $125/mes SaaS (plan Business promedio)
- $75/mes fees marketplace (1-2 hires/año distribuido)
= $200/mes ARPU blended

vs. $100/mes solo SaaS → 2x revenue por cliente
```

### 2.2 Fases de Monetización

**Fase 1 (Año 1): Solo SaaS**
```
Goal: Validar ATS, llegar a 50 clientes
Revenue: 100% SaaS
Foco: Product-market fit, retención
```

**Fase 2 (Año 2): Integración Marketplace**
```
Goal: Validar modelo dual, primeros $50k fees
Revenue: 85% SaaS + 15% fees
Foco: Curación de red, matching, conversión
```

**Fase 3 (Año 3+): Full Híbrido**
```
Goal: Escala, $900k ARR
Revenue: 65% SaaS + 35% fees
Foco: Network effects, ML matching, expansión
```

### 2.3 Unit Economics Dual

**Customer Acquisition Cost (CAC):**
- Año 1: $300 (PLG, content, referrals)
- Año 2: $400 (paid ads, partnerships)
- Año 3: $500 (enterprise sales, SDR)

**Lifetime Value (LTV) - Evolución:**

**Solo SaaS (Año 1):**
```
ARPU: $100/mes
Churn: 5%/mes
Lifetime: 20 meses
LTV: $2,000
LTV:CAC = 6.7:1 ✅
```

**SaaS + Marketplace (Año 3):**
```
ARPU: $200/mes ($125 SaaS + $75 fees)
Churn: 3%/mes (más sticky por dual value)
Lifetime: 33 meses
LTV: $6,600
LTV:CAC = 13:1 🚀
```

**Gross Margin Blended:**
```
SaaS COGS: $30/cliente/mes → 70% margin
Marketplace COGS: ~$0 (ya evaluados) → 95% margin
Blended Año 3: 85% gross margin
```

---

## 3. Arquitectura Multi-Tenant

### 2.1 Arquitectura de Tenancy

**Tipo**: Multi-tenant con aislamiento por CompanyID (Tenant Scoping)

```go
// Cada entidad principal tiene company_id
type Job struct {
    CompanyID uint `gorm:"not null;index"`
    // ... otros campos
}

type Candidate struct {
    CompanyID uint `gorm:"not null;index"`
    // ... otros campos
}

// Scoped queries automáticos
db.Where("company_id = ?", currentUser.CompanyID).Find(&jobs)
```

### 2.2 Jerarquía de Entidades

```
SuperAdmin (CompanyID = NULL)
    └── Acceso total a todas las empresas
    └── Gestión de sistema

Company (Tenant)
    ├── Users (via Memberships)
    ├── Jobs
    ├── Candidates
    └── Applications
```

### 2.3 Membership Model (Relación User-Company)

```go
type Membership struct {
    UserID    uint   // Usuario
    CompanyID *uint  // NULL = SuperAdmin global
    Role      string // superadmin, admin, recruiter, hiring_manager, user
    Status    string // active, inactive, pending
    IsDefault bool   // Membership por defecto al login
}
```

**Reglas de Membership**:
1. Un usuario puede pertenecer a múltiples empresas
2. Un usuario tiene UN rol por empresa
3. SuperAdmin NO tiene CompanyID (acceso global)
4. Al invitar un usuario, se crea membership con status=pending

---

## 3. Sistema de Roles y Permisos

### 3.1 Jerarquía de Roles

```
┌─────────────────────────────────────────────────────────┐
│ SuperAdmin (Level 100)                                  │
│ ├─ Acceso a todas las empresas                          │
│ ├─ Gestión de planes y facturación                      │
│ └─ Configuración del sistema                            │
└─────────────────────────────────────────────────────────┘
          │
          v
┌─────────────────────────────────────────────────────────┐
│ Admin (Level 50)                                        │
│ ├─ Administrador de su empresa                          │
│ ├─ Gestión de usuarios (invitar, roles)                 │
│ ├─ Configuración de empresa (logo, plan)                │
│ ├─ Acceso total a jobs, candidatos, aplicaciones        │
│ └─ Ver analytics de empresa                             │
└─────────────────────────────────────────────────────────┘
          │
          v
┌─────────────────────────────────────────────────────────┐
│ Recruiter (Level 30)                                    │
│ ├─ Crear y gestionar jobs asignados                     │
│ ├─ Ver todos los candidatos de la empresa               │
│ ├─ Mover candidatos por pipeline                        │
│ ├─ Crear y editar evaluaciones técnicas                 │
│ └─ NO puede cambiar plan ni invitar admins              │
└─────────────────────────────────────────────────────────┘
          │
          v
┌─────────────────────────────────────────────────────────┐
│ Hiring Manager (Level 20)                               │
│ ├─ Ver jobs donde es HiringManager                      │
│ ├─ Revisar candidatos de sus jobs                       │
│ ├─ Dejar feedback/rating                                │
│ └─ NO puede crear jobs ni mover stages críticos         │
└─────────────────────────────────────────────────────────┘
          │
          v
┌─────────────────────────────────────────────────────────┐
│ User (Level 10)                                         │
│ ├─ Acceso básico (read-only)                            │
│ ├─ Ver listado de candidatos                            │
│ └─ Sin permisos de escritura                            │
└─────────────────────────────────────────────────────────┘
```

### 3.2 Matriz de Permisos

| Recurso | SuperAdmin | Admin | Recruiter | Hiring Manager | User |
|---------|-----------|-------|-----------|----------------|------|
| **Companies** |
| Ver todas | ✅ | ❌ | ❌ | ❌ | ❌ |
| Ver propia | ✅ | ✅ | ✅ | ✅ | ✅ |
| Editar | ✅ | ✅ | ❌ | ❌ | ❌ |
| **Users/Memberships** |
| Invitar usuarios | ✅ | ✅ | ❌ | ❌ | ❌ |
| Cambiar roles | ✅ | ✅ (excepto admin) | ❌ | ❌ | ❌ |
| **Jobs** |
| Crear | ✅ | ✅ | ✅ | ❌ | ❌ |
| Ver todos | ✅ | ✅ | ✅ | Solo asignados | ✅ |
| Editar | ✅ | ✅ | Solo asignados | ❌ | ❌ |
| Eliminar | ✅ | ✅ | ❌ | ❌ | ❌ |
| Asignar recruiter | ✅ | ✅ | ❌ | ❌ | ❌ |
| **Candidates** |
| Crear | ✅ | ✅ | ✅ | ❌ | ❌ |
| Ver todos | ✅ | ✅ | ✅ | Solo de sus jobs | ✅ (limitado) |
| Editar | ✅ | ✅ | ✅ | ❌ | ❌ |
| Eliminar | ✅ | ✅ | ❌ | ❌ | ❌ |
| **Applications** |
| Crear | ✅ | ✅ | ✅ | ❌ | ❌ |
| Mover stage | ✅ | ✅ | ✅ | Limitado* | ❌ |
| Rating (1-5) | ✅ | ✅ | ✅ | ✅ | ❌ |
| Marcar hired/rejected | ✅ | ✅ | ✅ | ❌ | ❌ |

\* Hiring Manager solo puede mover entre: screening ↔️ technical ↔️ offer

### 3.3 Implementación de Permisos (Pendiente)

```go
// middleware/auth.go
func RequireAuth() gin.HandlerFunc {
    return func(c *gin.Context) {
        tokenString := extractToken(c)
        claims, err := jwtService.ValidateToken(tokenString)
        if err != nil {
            c.JSON(401, gin.H{"error": "Unauthorized"})
            c.Abort()
            return
        }
        
        // Inyectar claims en context
        c.Set("user_id", claims.UserID)
        c.Set("company_id", claims.CompanyID)
        c.Set("role", claims.Role)
        c.Next()
    }
}

// middleware/permissions.go
func RequireRole(minLevel int) gin.HandlerFunc {
    return func(c *gin.Context) {
        role := c.GetString("role")
        userLevel := getRoleLevel(role)
        
        if userLevel < minLevel {
            c.JSON(403, gin.H{"error": "Insufficient permissions"})
            c.Abort()
            return
        }
        c.Next()
    }
}

// Uso en routes
api.POST("/jobs", RequireAuth(), RequireRole(30), jobHandler.CreateJob)
```

---

## 4. Reglas de Negocio Fundamentales

### 4.1 Gestión de Companies

**RN-COMP-001: Creación de Company**
- Al crear una empresa, automáticamente se inicia trial de 14 días
- Plan inicial: `free` con límites restrictivos
- TrialEndsAt = Now() + 14 días
- Timezone por defecto: America/Bogota

**RN-COMP-002: Slug Único**
- Cada empresa tiene slug único (usado en URLs: techrecruit.app/azentic-sys)
- Validación: alfanumérico, min 2, max 100 caracteres

**RN-COMP-003: Plan Tiers**
```go
const (
    PlanFree       = "free"       // Trial o gratuito permanente
    PlanStarter    = "starter"    // $0/mes (freemium limitado)
    PlanProfessional = "professional" // $49/mes por usuario
    PlanBusiness   = "business"   // $149/mes por usuario
    PlanEnterprise = "enterprise" // Custom pricing
)
```

### 4.2 Gestión de Jobs

**RN-JOB-001: Estados de Job**
```go
const (
    JobStatusDraft   = "draft"    // Borrador, no visible
    JobStatusActive  = "active"   // Publicado, acepta aplicaciones
    JobStatusOnHold  = "on_hold"  // Pausado temporalmente
    JobStatusClosed  = "closed"   // Cerrado (contratación completada)
)
```

**RN-JOB-002: Asignación de Recruiter**
- Un job DEBE tener AssignedRecruiter (obligatorio al publicar)
- HiringManager es opcional (puede ser el Admin)
- Solo Admin puede asignar/cambiar recruiter

**RN-JOB-003: Publicación de Job**
- Para cambiar de draft → active:
  - Requiere: Title, Description, Location
  - Requiere: AssignedRecruiter asignado
  - Validar límite de posiciones activas según plan

### 4.3 Gestión de Candidates

**RN-CAND-001: Email Único por Company**
- Un candidato puede existir en múltiples companies
- Dentro de una company, email debe ser único
- Índice compuesto: `idx_candidates_company_email`

**RN-CAND-002: Source Tracking**
```go
const (
    SourceLinkedIn    = "linkedin"
    SourceReferral    = "referral"
    SourceDirectApply = "direct_apply"
    SourceAgency      = "agency"
    SourceJobBoard    = "job_board"
)
```

**RN-CAND-003: Datos Obligatorios**
- Mínimo requerido: Email, FirstName, LastName
- Recomendado: ResumeURL (PDF en S3)
- Opcional pero valioso: GithubURL (para evaluación técnica)

### 4.4 Gestión de Applications

**RN-APP-001: Pipeline Stages**
```go
const (
    StageApplied    = "applied"     // Candidato aplicó
    StageScreening  = "screening"   // Screening telefónico/inicial
    StageTechnical  = "technical"   // Evaluación técnica
    StageOffer      = "offer"       // Oferta extendida
    StageHired      = "hired"       // ¡Contratado! 🎉
    StageRejected   = "rejected"    // Rechazado en cualquier stage
)
```

**RN-APP-002: Transiciones de Stage Permitidas**
```
applied → screening
screening → technical | rejected
technical → offer | rejected
offer → hired | rejected
hired → [FINAL STATE]
rejected → [FINAL STATE]
```

**RN-APP-003: Timestamps Automáticos**
- `AppliedAt`: Se setea en creación (auto now())
- `RejectedAt`: Se setea cuando stage → rejected
- `HiredAt`: Se setea cuando stage → hired

**RN-APP-004: Rating System**
- Rating: 1-5 estrellas (nullable)
- Puede cambiar en cualquier momento
- Usado para ranking de candidatos

**RN-APP-005: Company ID Redundante**
- Applications tiene CompanyID (redundante con Job.CompanyID)
- Razón: Optimizar queries de "todas las aplicaciones de mi empresa"
- Se propaga automáticamente desde Job al crear Application

---

## 5. Pipeline de Candidatos

### 5.1 Flujo Visual del Pipeline

```
┌─────────────┐
│   APPLIED   │ ← Candidato aplica (manual o portal público)
└──────┬──────┘
       │
       v
┌─────────────┐
│  SCREENING  │ ← Recruiter hace screening telefónico
└──────┬──────┘   Rating inicial (1-5)
       │          Notas: "Buena comunicación"
       │
       ├─────────────────────────┐
       │                         │
       v                         v
┌─────────────┐          ┌─────────────┐
│  TECHNICAL  │          │  REJECTED   │ ← No pasó screening
└──────┬──────┘          └─────────────┘
       │                 RejectedAt = NOW
       │ ← Evaluación técnica:
       │   - Code challenge (GitHub)
       │   - Live coding
       │   - System design
       │
       ├─────────────────────────┐
       │                         │
       v                         v
┌─────────────┐          ┌─────────────┐
│    OFFER    │          │  REJECTED   │ ← No pasó técnica
└──────┬──────┘          └─────────────┘
       │
       │ ← Oferta extendida
       │   - Salary proposal
       │   - Benefits
       │
       ├─────────────────────────┐
       │                         │
       v                         v
┌─────────────┐          ┌─────────────┐
│    HIRED    │          │  REJECTED   │ ← Rechazó oferta
└─────────────┘          └─────────────┘
HiredAt = NOW            RejectedAt = NOW
```

### 5.2 Métricas del Pipeline

**Métricas a trackear (futuro analytics)**:
- **Conversion Rate**: % applied → hired
- **Bottleneck Detection**: Stage donde más candidatos se estancan
- **Time-to-Hire**: Días promedio de applied → hired
- **Stage Duration**: Tiempo promedio en cada stage
- **Rejection Reasons**: Por qué rechazamos (futuro: tags)

### 5.3 Automatizaciones del Pipeline (Roadmap)

1. **Auto-rejection de inactivos** (Fase 2)
   - Si candidato en `screening` por >30 días → auto-move a `rejected`
   
2. **Notificaciones automáticas** (Fase 2)
   - Candidato cambia stage → Email automático
   - Recruiter asignado → Slack notification

3. **AI Scoring** (Fase 3)
   - Análisis de CV con IA → Score 0-100
   - Matching job description vs resume

---

## 6. Red Dvra - Marketplace de Talento

### 6.1 Concepto y Valor Único

**Red Dvra** es una red curada de desarrolladores LATAM pre-evaluados técnicamente, integrada directamente en el ATS.

**Diferencia con competencia:**
```
LinkedIn Recruiter: $99/mes → Base de datos sin curación
Mercor/Turing: Fee 20-30% → Solo marketplace, sin herramienta
Dvra: $49-149/mes ATS + Fee 15-20% → Tool + Talent
```

**Value Proposition para Empresas:**
1. Ya usan Dvra ATS para proceso interno
2. Cuando necesitan talento externo → Buscan en Red Dvra
3. Candidatos pre-evaluados (ahorro de 15+ horas screening)
4. Fee solo por contratación exitosa (risk-free)
5. Timezone y cultura LATAM compatible

**Value Proposition para Candidatos:**
1. Validación técnica objetiva (GitHub, challenges)
2. Visibilidad ante empresas que ya usan la plataforma
3. Proceso de aplicación simplificado
4. Oportunidades remotas internacional
5. Feedback constructivo de evaluaciones
6. **100% voluntario** - Opt-in, no spam, control total de tus datos

### 6.2 Modelo de Datos Marketplace

```go
// Candidato de la Red Dvra (diferente de Candidate interno)
type NetworkCandidate struct {
    gorm.Model
    
    // Datos personales
    Email         string `gorm:"type:varchar(255);uniqueIndex;not null"`
    FirstName     string `gorm:"type:varchar(100);not null"`
    LastName      string `gorm:"type:varchar(100);not null"`
    Phone         string
    Country       string `gorm:"type:varchar(50)"`
    City          string
    Timezone      string // "America/Bogota", "America/Mexico_City"
    
    // Perfil profesional
    Title         string // "Senior Backend Developer"
    Bio           string `gorm:"type:text"`
    YearsExp      int
    SeniorityLevel string // junior, mid, senior, staff
    
    // Links
    LinkedinURL   string `gorm:"type:text"`
    GithubURL     string `gorm:"type:text"`
    PortfolioURL  string `gorm:"type:text"`
    ResumeURL     string `gorm:"type:text"` // S3
    
    // Evaluación Dvra
    DvraScore         int     // 0-100 (score global)
    CodeQualityScore  int     // 0-100 (GitHub analysis)
    ChallengeScore    int     // 0-100 (coding challenge)
    CommunicationScore int    // 0-100 (entrevista)
    EvaluatedAt       *time.Time
    EvaluatedBy       *uint   // Admin/Evaluator UserID
    
    // Estado
    Status            string  // pending, evaluating, approved, featured, inactive
    IsFeatured        bool    // Top candidatos destacados
    AvailabilityStatus string // available, interviewing, hired, not_looking
    HourlyRate        *int    // USD/hora (para contractors)
    PreferredRemote   bool
    
    // Metadata
    SourceChannel     string  // landing_page, referral, bootcamp, linkedin
    ReferredBy        *uint   // NetworkCandidateID del referidor
    
    // Relaciones
    Skills            []NetworkCandidateSkill `gorm:"foreignKey:CandidateID"`
    Interests         []NetworkCandidateInterest `gorm:"foreignKey:CandidateID"`
    Applications      []NetworkApplication `gorm:"foreignKey:CandidateID"`
}

// Skills con niveles de proficiency
type NetworkCandidateSkill struct {
    gorm.Model
    CandidateID uint   `gorm:"not null;index"`
    SkillName   string `gorm:"type:varchar(50)"` // "React", "Node.js", "AWS"
    Category    string // frontend, backend, devops, mobile, data
    Proficiency string // beginner, intermediate, advanced, expert
    YearsExp    int
    
    // Auto-detectado de GitHub
    IsVerified  bool   // Si se verificó via GitHub repos
}

// Intereses/Preferencias
type NetworkCandidateInterest struct {
    gorm.Model
    CandidateID uint `gorm:"not null;index"`
    
    // Preferencias de trabajo
    RemoteOnly        bool
    WillingToRelocate bool
    PreferredCountries string `gorm:"type:text"` // JSON array
    PreferredIndustries string `gorm:"type:text"` // ["fintech", "healthtech"]
    
    // Compensación
    MinSalaryExpected int // USD anual
    MaxSalaryExpected int
}

// Interacción entre Company y NetworkCandidate
type NetworkApplication struct {
    gorm.Model
    CompanyID     uint   `gorm:"not null;index"`
    JobID         *uint  // Nullable: puede buscar sin job específico
    CandidateID   uint   `gorm:"not null;index"`
    
    // Flujo
    Status        string // interested, contacted, interviewing, offer, hired, rejected
    InitiatedBy   string // company, candidate
    
    // Communication
    FirstContactAt *time.Time
    InterviewedAt  *time.Time
    OfferExtendedAt *time.Time
    HiredAt        *time.Time
    RejectedAt     *time.Time
    RejectionReason string
    
    // Fee tracking
    FeeAmount     *int    // USD (15-20% del salary)
    FeeStatus     string  // pending, invoiced, paid
    InvoicedAt    *time.Time
    PaidAt        *time.Time
    
    // Relaciones
    Candidate     *NetworkCandidate `gorm:"foreignKey:CandidateID"`
    Company       *Company `gorm:"foreignKey:CompanyID"`
    Job           *Job `gorm:"foreignKey:JobID"`
}

// Evaluaciones técnicas de la Red
type NetworkEvaluation struct {
    gorm.Model
    CandidateID uint `gorm:"not null;index"`
    
    // Tipo de evaluación
    EvaluationType string // github_analysis, coding_challenge, live_interview
    
    // GitHub Analysis (automático)
    GithubReposAnalyzed int
    TotalCommits        int
    LanguagesUsed       string `gorm:"type:text"` // JSON
    TopRepoStars        int
    CodeQualityScore    int
    
    // Coding Challenge
    ChallengeID         *uint
    ChallengeName       string
    CompletedAt         *time.Time
    TimeSpentMinutes    int
    TestsPassed         int
    TestsTotal          int
    CodeReviewNotes     string `gorm:"type:text"`
    
    // Live Interview
    InterviewerID       *uint
    InterviewDuration   int // minutos
    TechnicalScore      int // 0-100
    CommunicationScore  int // 0-100
    InterviewNotes      string `gorm:"type:text"`
    
    // Score final
    FinalScore          int // 0-100
    Recommendation      string // strong_yes, yes, maybe, no
}
```

### 6.3 Sourcing de Talento - Sistema de Estados

> **⚠️ CRÍTICO: Legal, Ético y Escalable**

El sourcing de talento para Red Dvra debe cumplir con:
- **GDPR** (Europa) - Consentimiento explícito
- **LGPD** (Brasil) - Privacidad y derecho al olvido
- **CAN-SPAM Act** (USA/Internacional) - Anti-spam
- **Ética profesional** - Opt-in voluntario, no spam

#### Sistema de Estados del Talento

```go
type TalentStatus string

const (
    // NO VISIBLE para empresas
    StatusProspect    TalentStatus = "prospect"     // Identificado vía scraping/research
    StatusInvited     TalentStatus = "invited"      // Invitación enviada, esperando respuesta
    
    // VISIBLE para empresas (con consentimiento)
    StatusActive      TalentStatus = "active"       // Registrado, profile completo, opted-in
    StatusFeatured    TalentStatus = "featured"     // Score >85, top 20%
    
    // Finales
    StatusInactive    TalentStatus = "inactive"     // Se dio de baja voluntariamente
    StatusBlacklisted TalentStatus = "blacklisted"  // No contactar nunca (solicitado)
)

type NetworkCandidate struct {
    gorm.Model
    
    // ============================================
    // KEYS DE DEDUPLICACIÓN (Índices únicos)
    // ============================================
    Email           string `gorm:"type:varchar(255);uniqueIndex;not null"`
    GithubUsername  string `gorm:"type:varchar(100);uniqueIndex"` // Opcional pero único
    
    // ============================================
    // ESTADO Y TRACKING
    // ============================================
    Status          TalentStatus `gorm:"type:varchar(50);not null;default:'prospect';index"`
    
    // Timestamps de transición
    ProspectedAt    *time.Time   // Cuando lo identificamos (scraping/referral)
    InvitedAt       *time.Time   // Cuando enviamos invitación
    RegisteredAt    *time.Time   // Cuando completó registro voluntario
    LastContactAt   *time.Time   // Último email/interacción
    
    // ============================================
    // CONSENTIMIENTO (CRÍTICO LEGAL)
    // ============================================
    OptedIn         bool `gorm:"default:false;not null"` // Aceptó ser visible
    OptedInAt       *time.Time
    OptOutReason    string // Si se dio de baja, el motivo
    
    // GDPR/LGPD compliance
    ConsentIP       string // IP desde donde dio consentimiento
    ConsentSource   string // "landing_page", "github_oauth", "bootcamp_partner"
    
    // ============================================
    // PERFIL (se enriquece según status)
    // ============================================
    FirstName       string `gorm:"type:varchar(100)"`
    LastName        string `gorm:"type:varchar(100)"`
    // ... resto de campos normales
    
    // Source tracking
    SourceChannel   string // landing_page, github_scraping, bootcamp_partnership, referral
    ReferredBy      *uint  `gorm:"index"` // NetworkCandidateID del referidor
    
    // Anti-spam
    InvitationsSent int `gorm:"default:0"` // Max 3 invitaciones lifetime
    
    // Relaciones
    Skills          []NetworkCandidateSkill `gorm:"foreignKey:CandidateID"`
}
```

#### Transiciones de Estado Permitidas

```
PROSPECT (internal only)
    │
    ├──> INVITED (sent invitation)
    │       │
    │       ├──> ACTIVE (completed registration + opted-in) ✅
    │       └──> BLACKLISTED (clicked "never contact")
    │
    └──> ACTIVE (direct registration via landing page) ✅
            │
            ├──> FEATURED (score >85, auto-promoted)
            ├──> INACTIVE (user opted-out)
            └──> BLACKLISTED (requested removal)
```

**Regla Crítica de Visibilidad:**
```go
// SOLO mostrar a empresas si:
func (nc *NetworkCandidate) IsVisibleToCompanies() bool {
    return (nc.Status == StatusActive || nc.Status == StatusFeatured) &&
           nc.OptedIn == true &&
           nc.DeletedAt == nil
}
```

---

### 6.4 Estrategias de Sourcing (Por Año)

#### **AÑO 1: Opt-In Puro (100% Legal, 100% Ético)**

**Canales Principales:**

**1. Landing Page con GitHub OAuth (Primary)**
```
Landing: talento.dvra.com
    ↓
CTA: "Conectar con GitHub" (OAuth)
    ↓
User autoriza acceso
    ↓
Auto-generamos perfil de sus repos públicos:
    - Languages (weighted by LoC)
    - Top repos, stars, commits
    - Contribution activity
    ↓
User completa datos faltantes:
    - Nombre completo
    - Email de contacto
    - Seniority level
    - Disponibilidad
    ↓
Checkbox: "Acepto ser visible para empresas" (OptedIn = true)
    ↓
Status: ACTIVE → Visible para empresas Business+
```

**Ventajas:**
- ✅ 100% consentimiento explícito
- ✅ Datos actualizados (usuario los confirma)
- ✅ GitHub verification automática
- ✅ Zero riesgo legal

**Implementación:**
```go
// Endpoint: POST /api/v1/network/register-via-github
func (h *NetworkHandler) RegisterViaGitHub(c *gin.Context) {
    // 1. OAuth callback con GitHub token
    githubToken := c.Query("code")
    githubUser := fetchGitHubUser(githubToken)
    
    // 2. Buscar si ya existe (deduplicación)
    var candidate NetworkCandidate
    result := h.db.Where("github_username = ?", githubUser.Login).First(&candidate)
    
    if result.Error == gorm.ErrRecordNotFound {
        // Nuevo candidato
        candidate = NetworkCandidate{
            Email:          githubUser.Email,
            GithubUsername: githubUser.Login,
            GithubURL:      githubUser.HTMLURL,
            Status:         StatusProspect, // Aún no opted-in
            ProspectedAt:   timePtr(time.Now()),
            SourceChannel:  "github_oauth",
            ConsentIP:      c.ClientIP(),
            ConsentSource:  "landing_page",
        }
        
        // Analizar repos y extraer skills
        repos := fetchGitHubRepos(githubToken)
        candidate.Skills = extractSkillsFromRepos(repos)
        
        h.db.Create(&candidate)
    }
    
    // 3. Redirigir a formulario de completar perfil
    c.Redirect(302, "/complete-profile?token="+generateToken(candidate.ID))
}

// Endpoint: POST /api/v1/network/complete-profile
func (h *NetworkHandler) CompleteProfile(c *gin.Context) {
    var dto CompleteProfileDTO
    c.BindJSON(&dto)
    
    candidate := getCandidateByToken(dto.Token)
    
    // Actualizar datos
    candidate.FirstName = dto.FirstName
    candidate.LastName = dto.LastName
    candidate.Phone = dto.Phone
    candidate.Title = dto.Title
    candidate.SeniorityLevel = dto.SeniorityLevel
    
    // CRÍTICO: Consentimiento explícito
    candidate.OptedIn = dto.AcceptTerms // Checkbox "Acepto ser visible"
    candidate.OptedInAt = timePtr(time.Now())
    candidate.Status = StatusActive
    candidate.RegisteredAt = timePtr(time.Now())
    
    h.db.Save(&candidate)
    
    // Analytics event
    trackEvent("candidate_registered", candidate.ID)
    
    c.JSON(200, gin.H{
        "message": "¡Bienvenido a Red Dvra!",
        "profile_url": "/profile/" + candidate.ID,
    })
}
```

**2. Partnerships con Bootcamps (High Quality Pipeline)**

**Bootcamps Target:**
- Platzi (Colombia/México/LATAM)
- Henry (Argentina/LATAM)
- Coderhouse (Argentina)
- Digital House (Argentina/México)
- Laboratoria (LATAM - mujeres tech)

**Deal Structure:**
```
Bootcamp comparte lista de egresados (con permiso previo)
    ↓
Envío email desde bootcamp (not Dvra):
    "Como egresado de [Bootcamp], tienes acceso a Red Dvra"
    ↓
Link con token pre-autorizado
    ↓
Completan perfil → Status: ACTIVE directo
```

**Benefit para Bootcamp:**
- Value add para sus estudiantes (job placement)
- Posible revenue share si hay hire (5-10% del fee)

**Meta Año 1:**
- 3 partnerships activos
- 50-100 candidatos/bootcamp/año
- Total: 150-300 candidatos de alta calidad

**3. Referral Program (Viral Growth)**

```go
type Referral struct {
    gorm.Model
    ReferrerID      uint   // Candidato que refiere
    ReferredEmail   string
    ReferredID      *uint  // Cuando completa registro
    Status          string // pending, registered, approved, rewarded
    RewardAmount    int    // $50 si ambos aprobados
    RewardPaidAt    *time.Time
}
```

**Mecánica:**
- Candidato activo invita amigos (via email o link único)
- Amigo se registra con link de referral
- Cuando AMBOS son aprobados (score >70) → $50 c/u (gift card Amazon)
- Límite: Max 10 referrals/candidato (anti-abuse)

**Viral Coefficient Target:**
- Si 30% de candidatos refieren 1 persona que convierte
- k = 0.30 (sub-viral, pero ayuda)

**Meta Año 1 via Referrals:** 50-80 candidatos

---

#### **AÑO 2: Opt-In + Outreach Inteligente**

Una vez validado el modelo, expandir con:

**4. Scraping Light + Personalized Outreach**

```go
// Background job: Identificar prospects en GitHub
func (s *ScrapingService) IdentifyProspects(criteria ScrapingCriteria) {
    // Buscar en GitHub API (rate limited):
    // - Repos con 50+ stars en Go/Python/JS
    // - Contributors activos (10+ commits/mes)
    // - Location: LATAM countries
    
    users := searchGitHubUsers(criteria)
    
    for _, user := range users {
        // Deduplicación: No contactar si ya existe
        exists := s.db.Where("github_username = ?", user.Login).First(&existing)
        if exists == nil {
            continue // Ya lo tenemos
        }
        
        // Crear como PROSPECT (internal only)
        prospect := NetworkCandidate{
            Email:          user.Email, // GitHub public email
            GithubUsername: user.Login,
            GithubURL:      user.HTMLURL,
            Status:         StatusProspect,
            ProspectedAt:   timePtr(time.Now()),
            SourceChannel:  "github_scraping",
        }
        
        s.db.Create(&prospect)
        
        // Queue para enviar invitación (no inmediato)
        queueInvitation(prospect.ID, delayDays=7)
    }
}

// Envío de invitación personalizada
func (s *OutreachService) SendInvitation(candidateID uint) error {
    candidate, _ := s.GetByID(candidateID)
    
    // Validaciones anti-spam
    if candidate.InvitationsSent >= 3 {
        return errors.New("max invitations reached")
    }
    
    if candidate.Status == StatusBlacklisted {
        return errors.New("candidate blacklisted")
    }
    
    // Generar magic link con token
    token := generateMagicToken(candidate.ID)
    inviteURL := "https://talento.dvra.com/invite/" + token
    
    // Email personalizado (NO spam)
    email := EmailTemplate{
        To:      candidate.Email,
        Subject: "Tu perfil de GitHub llamó nuestra atención - Red Dvra",
        Body: fmt.Sprintf(`
Hola!

Vimos tu perfil de GitHub (%s) y nos impresionó tu trabajo en %s.

Red Dvra es una red curada de desarrolladores LATAM donde empresas tech 
buscan talento pre-evaluado. 

¿Te interesa aparecer en nuestra red? Solo toma 2 minutos:
%s

Preview de lo que vimos:
- %d repos públicos
- Languages: %s
- %d commits en los últimos 6 meses

Si no te interesa, ignora este email o haz click aquí para no recibir más: %s

Saludos,
Equipo Dvra

---
Este es un email one-time. No enviaremos más sin tu consentimiento.
        `, candidate.GithubURL, topRepo, inviteURL, repoCount, langs, commits, optOutURL),
    }
    
    // Enviar
    err := s.emailService.Send(email)
    if err != nil {
        return err
    }
    
    // Update status y tracking
    candidate.Status = StatusInvited
    candidate.InvitedAt = timePtr(time.Now())
    candidate.InvitationsSent++
    candidate.LastContactAt = timePtr(time.Now())
    
    return s.db.Save(candidate).Error
}

// Opt-out permanente (MUST HAVE para compliance)
func (s *OutreachService) Blacklist(token string) error {
    candidate := getCandidateByToken(token)
    
    candidate.Status = StatusBlacklisted
    candidate.OptOutReason = "opted_out_from_email"
    
    s.db.Save(candidate)
    
    log.Info("Candidate blacklisted", "email", candidate.Email)
    return nil
}
```

**Reglas Estrictas de Outreach:**
- ✅ Max 3 invitaciones lifetime por candidato
- ✅ Spacing de 30 días entre invitaciones
- ✅ Opt-out link visible en cada email
- ✅ One-click unsubscribe (StatusBlacklisted)
- ✅ NO mostrar prospects a empresas (solo Invited/Active)

**Meta Año 2 via Outreach:** 300-500 candidatos adicionales

---

#### **AÑO 3: Partnerships at Scale + Community**

**5. Tech Communities & Events**

- Sponsor en meetups (JS LATAM, Python Argentina, etc.)
- Stand en hackathons: "Registra tu perfil, gana visibilidad"
- Discord/Slack communities: Canal de #jobs powered by Dvra
- Conferencias: JSConf Argentina, PyCon Colombia, etc.

**6. Content Marketing (SEO)**

- Blog: "Salarios tech LATAM 2027", "Remote work guide"
- SEO: Ranking #1 para "trabajos remotos desarrollador LATAM"
- YouTube: Canal con tips de carrera tech
- Podcast: Entrevistas a devs LATAM exitosos

**Meta Año 3:** 2,000+ candidatos activos, mix de canales

---

### 6.5 Reglas de Negocio del Marketplace

**RN-NETWORK-001: Registro de Candidato (Actualizado)**

**Vías de Registro:**
1. **Landing page directa** (talento.dvra.com)
   - GitHub OAuth (primary)
   - Formulario manual (secondary)
   - Status inicial: `prospect` → Completar perfil → `active`

2. **Partnership bootcamp**
   - Email con token pre-autorizado
   - Status inicial: `invited` → Registro → `active`

3. **Referral de otro candidato**
   - Link único con tracking
   - Status inicial: `invited` → Registro → `active`

4. **Outreach (Año 2+)**
   - Email de invitación personalizada
   - Status inicial: `prospect` → `invited` → Registro → `active`

**Validación Común (todas las vías):**
- Email único (deduplicación)
- GitHub público válido (si proveen)
- Checkbox consentimiento explícito: "Acepto ser visible para empresas"
- IP y source tracking para compliance

**Email de Confirmación:**
```
Subject: "¡Bienvenido a Red Dvra!"

Hola [Name],

Tu perfil está siendo revisado por nuestro equipo (2-5 días).

Mientras tanto:
- Completa tu perfil al 100%
- Conecta tu GitHub para auto-verification
- Completa un coding challenge (opcional, pero +20 score)

Una vez aprobado, serás visible para 150+ empresas tech.

[Completar Perfil]
```

---

### 6.6 Implementación de Deduplicación

```go
// Service: Crear o actualizar candidato sin duplicados
func (s *NetworkCandidateService) CreateOrUpdate(dto CreateCandidateDTO) (*NetworkCandidate, error) {
    var candidate NetworkCandidate
    
    // Buscar por email O github_username (cualquiera identifica único)
    result := s.db.Where("email = ? OR github_username = ?", 
        dto.Email, dto.GithubUsername).First(&candidate)
    
    if result.Error == gorm.ErrRecordNotFound {
        // NO EXISTE → Crear nuevo
        candidate = NetworkCandidate{
            Email:          dto.Email,
            GithubUsername: dto.GithubUsername,
            FirstName:      dto.FirstName,
            LastName:       dto.LastName,
            Status:         StatusProspect,
            ProspectedAt:   timePtr(time.Now()),
            SourceChannel:  dto.Source,
            OptedIn:        false, // Default: no consentimiento aún
        }
        
        if err := s.db.Create(&candidate).Error; err != nil {
            // Handle unique constraint violation
            if strings.Contains(err.Error(), "duplicate key") {
                // Race condition: otro goroutine creó el registro
                return s.CreateOrUpdate(dto) // Retry
            }
            return nil, err
        }
        
        log.Info("Nuevo candidato creado", 
            "id", candidate.ID,
            "email", candidate.Email, 
            "status", candidate.Status,
            "source", candidate.SourceChannel,
        )
        
        return &candidate, nil
    }
    
    // YA EXISTE → Evaluar si actualizar
    log.Info("Candidato existente encontrado", 
        "id", candidate.ID,
        "current_status", candidate.Status,
    )
    
    // Regla: Solo actualizar si está en estados tempranos
    if candidate.Status == StatusProspect || candidate.Status == StatusInvited {
        // Enriquecer datos, pero NO sobreescribir info existente
        updated := false
        
        if candidate.FirstName == "" && dto.FirstName != "" {
            candidate.FirstName = dto.FirstName
            updated = true
        }
        
        if candidate.LastName == "" && dto.LastName != "" {
            candidate.LastName = dto.LastName
            updated = true
        }
        
        if candidate.GithubUsername == "" && dto.GithubUsername != "" {
            candidate.GithubUsername = dto.GithubUsername
            updated = true
        }
        
        // Actualizar source solo si es más específico
        if dto.Source != "" && candidate.SourceChannel == "unknown" {
            candidate.SourceChannel = dto.Source
            updated = true
        }
        
        if updated {
            candidate.UpdatedAt = time.Now()
            if err := s.db.Save(&candidate).Error; err != nil {
                return nil, err
            }
            log.Info("Candidato actualizado", "id", candidate.ID)
        }
        
        return &candidate, nil
    }
    
    // Si ya es ACTIVE o FEATURED → NO TOCAR (usuario tiene control total)
    if candidate.Status == StatusActive || candidate.Status == StatusFeatured {
        log.Warn("Intento de actualizar candidato activo ignorado",
            "id", candidate.ID,
            "status", candidate.Status,
        )
        return &candidate, nil
    }
    
    // Si está BLACKLISTED → Rechazar operación
    if candidate.Status == StatusBlacklisted {
        return nil, errors.New("candidato en blacklist, no contactar")
    }
    
    return &candidate, nil
}

// Helper: Transicionar estado con validaciones
func (s *NetworkCandidateService) TransitionStatus(
    candidateID uint, 
    newStatus TalentStatus,
    metadata map[string]interface{},
) error {
    candidate, err := s.GetByID(candidateID)
    if err != nil {
        return err
    }
    
    // Validar transición permitida
    if !isValidTransition(candidate.Status, newStatus) {
        return fmt.Errorf(
            "transición inválida: %s → %s no permitida",
            candidate.Status, newStatus,
        )
    }
    
    oldStatus := candidate.Status
    candidate.Status = newStatus
    
    // Setear timestamps según transición
    switch newStatus {
    case StatusInvited:
        candidate.InvitedAt = timePtr(time.Now())
        candidate.InvitationsSent++
    case StatusActive:
        candidate.RegisteredAt = timePtr(time.Now())
        candidate.OptedIn = true // Mandatory para Active
        candidate.OptedInAt = timePtr(time.Now())
    case StatusFeatured:
        // Auto-promoción cuando score >85
    case StatusBlacklisted:
        reason, _ := metadata["reason"].(string)
        candidate.OptOutReason = reason
        candidate.OptedIn = false
    }
    
    if err := s.db.Save(candidate).Error; err != nil {
        return err
    }
    
    // Log audit trail
    log.Info("Status transition",
        "candidate_id", candidateID,
        "old_status", oldStatus,
        "new_status", newStatus,
        "metadata", metadata,
    )
    
    // Trigger side effects (emails, webhooks, etc.)
    s.handleStatusChange(candidate, oldStatus, newStatus)
    
    return nil
}

// Validar transiciones permitidas
func isValidTransition(from, to TalentStatus) bool {
    validTransitions := map[TalentStatus][]TalentStatus{
        StatusProspect: {StatusInvited, StatusActive, StatusBlacklisted},
        StatusInvited:  {StatusActive, StatusBlacklisted},
        StatusActive:   {StatusFeatured, StatusInactive, StatusBlacklisted},
        StatusFeatured: {StatusActive, StatusInactive, StatusBlacklisted},
        StatusInactive: {StatusActive}, // Puede reactivarse
        StatusBlacklisted: {}, // Estado final, no sale
    }
    
    allowed, exists := validTransitions[from]
    if !exists {
        return false
    }
    
    for _, status := range allowed {
        if status == to {
            return true
        }
    }
    
    return false
}
```

---

### 6.7 Compliance y Privacidad (GDPR/LGPD)

#### **Data Subject Rights (Derechos del Usuario)**

```go
// Right to Access (GDPR Art. 15)
func (s *NetworkCandidateService) ExportPersonalData(candidateID uint) (*DataExport, error) {
    candidate, _ := s.GetByID(candidateID)
    
    // Recopilar TODOS los datos del candidato
    export := DataExport{
        PersonalInfo: candidate,
        Skills:       candidate.Skills,
        Applications: candidate.Applications,
        Evaluations:  s.getEvaluations(candidateID),
        ActivityLog:  s.getActivityLog(candidateID),
        GeneratedAt:  time.Now(),
    }
    
    // Generar PDF o JSON
    return &export, nil
}

// Right to Erasure / "Right to be Forgotten" (GDPR Art. 17)
func (s *NetworkCandidateService) DeletePersonalData(candidateID uint) error {
    candidate, _ := s.GetByID(candidateID)
    
    // Validar que no haya procesos activos
    hasActiveApplications := s.db.Where(
        "candidate_id = ? AND status IN (?)", 
        candidateID, 
        []string{"interviewing", "offer"},
    ).First(&NetworkApplication{}).Error == nil
    
    if hasActiveApplications {
        return errors.New("no se puede eliminar: tiene aplicaciones activas")
    }
    
    // Anonimizar datos (GDPR permite mantener datos agregados)
    candidate.Email = fmt.Sprintf("deleted_%d@anonymized.com", candidateID)
    candidate.FirstName = "DELETED"
    candidate.LastName = "USER"
    candidate.Phone = ""
    candidate.LinkedinURL = ""
    candidate.GithubURL = ""
    candidate.PortfolioURL = ""
    candidate.ResumeURL = ""
    candidate.Bio = ""
    candidate.Status = StatusBlacklisted
    candidate.OptOutReason = "gdpr_deletion_request"
    
    // Soft delete
    candidate.DeletedAt = gorm.DeletedAt{Time: time.Now(), Valid: true}
    
    if err := s.db.Save(candidate).Error; err != nil {
        return err
    }
    
    log.Info("Candidato anonimizado (GDPR deletion)", "id", candidateID)
    
    // Eliminar archivos de S3
    s.s3Service.DeleteFile(candidate.ResumeURL)
    
    return nil
}

// Right to Rectification (GDPR Art. 16)
func (s *NetworkCandidateService) UpdatePersonalData(
    candidateID uint, 
    updates map[string]interface{},
) error {
    // Usuario puede actualizar sus propios datos en cualquier momento
    return s.db.Model(&NetworkCandidate{}).
        Where("id = ?", candidateID).
        Updates(updates).Error
}

// Right to Object (GDPR Art. 21) - Opt-out de marketing
func (s *NetworkCandidateService) OptOutMarketing(candidateID uint) error {
    return s.TransitionStatus(candidateID, StatusBlacklisted, map[string]interface{}{
        "reason": "opted_out_marketing",
    })
}
```

#### **Consent Management**

```go
// Registrar consentimiento explícito
type Consent struct {
    gorm.Model
    CandidateID   uint      `gorm:"not null;index"`
    ConsentType   string    // "profile_visibility", "email_marketing", "data_processing"
    Granted       bool      `gorm:"not null"`
    GrantedAt     time.Time
    RevokedAt     *time.Time
    IP            string    // IP desde donde se otorgó
    UserAgent     string    // Browser/device info
    ConsentText   string    `gorm:"type:text"` // Texto exacto que aceptó
    Version       string    // Versión de T&C
}

// Verificar si tiene consentimiento válido
func (s *NetworkCandidateService) HasValidConsent(
    candidateID uint, 
    consentType string,
) bool {
    var consent Consent
    err := s.db.Where(
        "candidate_id = ? AND consent_type = ? AND granted = true AND revoked_at IS NULL",
        candidateID, consentType,
    ).First(&consent).Error
    
    return err == nil
}

// Revocar consentimiento
func (s *NetworkCandidateService) RevokeConsent(
    candidateID uint,
    consentType string,
) error {
    return s.db.Model(&Consent{}).
        Where("candidate_id = ? AND consent_type = ?", candidateID, consentType).
        Updates(map[string]interface{}{
            "granted":    false,
            "revoked_at": time.Now(),
        }).Error
}
```

#### **Data Retention Policy**

```go
// Política de retención de datos
const (
    RetentionProspects  = 180 * 24 * time.Hour // 6 meses
    RetentionInvited    = 365 * 24 * time.Hour // 1 año
    RetentionInactive   = 730 * 24 * time.Hour // 2 años
    RetentionBlacklisted = 0                     // Mantener indefinidamente (legal)
)

// Background job: Limpiar datos antiguos
func (s *NetworkCandidateService) CleanupStaleData() error {
    now := time.Now()
    
    // Eliminar prospects muy viejos (nunca completaron registro)
    s.db.Where(
        "status = ? AND prospected_at < ?",
        StatusProspect, now.Add(-RetentionProspects),
    ).Delete(&NetworkCandidate{})
    
    // Eliminar invited sin respuesta
    s.db.Where(
        "status = ? AND invited_at < ?",
        StatusInvited, now.Add(-RetentionInvited),
    ).Delete(&NetworkCandidate{})
    
    // Anonimizar inactive muy antiguos
    var inactiveCandidates []NetworkCandidate
    s.db.Where(
        "status = ? AND updated_at < ?",
        StatusInactive, now.Add(-RetentionInactive),
    ).Find(&inactiveCandidates)
    
    for _, candidate := range inactiveCandidates {
        s.DeletePersonalData(candidate.ID) // Anonimizar
    }
    
    log.Info("Cleanup stale data completed")
    return nil
}
```

---

### 6.8 Query Performance & Indexes

```sql
-- Índices críticos para búsqueda de empresas
CREATE INDEX idx_network_candidates_search ON network_candidates(status, opted_in, deleted_at)
WHERE status IN ('active', 'featured') AND opted_in = true AND deleted_at IS NULL;

-- Índice para deduplicación rápida
CREATE UNIQUE INDEX idx_network_candidates_email_lower ON network_candidates(LOWER(email))
WHERE deleted_at IS NULL;

CREATE UNIQUE INDEX idx_network_candidates_github_lower ON network_candidates(LOWER(github_username))
WHERE deleted_at IS NULL AND github_username IS NOT NULL;

-- Índice compuesto para búsqueda por skills
CREATE INDEX idx_network_candidate_skills_search ON network_candidate_skills(skill_name, proficiency)
WHERE deleted_at IS NULL;

-- Índice para auditoría y compliance
CREATE INDEX idx_network_candidates_opted_in_at ON network_candidates(opted_in_at DESC)
WHERE opted_in = true;
```

**Query optimizada para empresas:**
```go
// Solo mostrar candidatos ACTIVOS con consentimiento
func (s *NetworkCandidateService) SearchForCompanies(filters SearchFilters) ([]NetworkCandidate, error) {
    query := s.db.Preload("Skills").
        Where("status IN (?)", []TalentStatus{StatusActive, StatusFeatured}).
        Where("opted_in = ?", true).
        Where("deleted_at IS NULL")
    
    // Filtros opcionales
    if filters.SeniorityLevel != "" {
        query = query.Where("seniority_level = ?", filters.SeniorityLevel)
    }
    
    if filters.Country != "" {
        query = query.Where("country = ?", filters.Country)
    }
    
    if filters.MinScore > 0 {
        query = query.Where("dvra_score >= ?", filters.MinScore)
    }
    
    // Filtro por skills (JOIN con tabla skills)
    if len(filters.RequiredSkills) > 0 {
        query = query.Joins(
            "INNER JOIN network_candidate_skills ON network_candidate_skills.candidate_id = network_candidates.id",
        ).Where("network_candidate_skills.skill_name IN (?)", filters.RequiredSkills)
    }
    
    // Ordenar: Featured primero, luego por score
    query = query.Order("is_featured DESC, dvra_score DESC")
    
    var candidates []NetworkCandidate
    err := query.Limit(50).Find(&candidates).Error
    
    return candidates, err
}
```

---

**RN-NETWORK-002: Proceso de Evaluación (Actualizado)**

**Año 1 (Manual - 2-3/semana):**
```
1. Admin revisa perfil → Status: evaluating
2. GitHub analysis manual:
   - Revisar repos principales
   - Calidad de código, tests, READMEs
   - Actividad reciente (commits last 6 months)
3. Enviar coding challenge (1-2 horas)
4. Revisar solución → Score 0-100
5. Decision:
   - Score >70 → Status: approved
   - Score >85 → IsFeatured: true
   - Score <70 → Status: rejected (con feedback)
```

**Año 2+ (Semi-automático):**
```
1. GitHub API analysis automático
2. Auto-send coding challenge
3. Auto-score con tests
4. Admin solo revisa Featured candidates
```

**RN-NETWORK-003: Sistema de Scoring**

```go
func CalculateDvraScore(candidate *NetworkCandidate) int {
    weights := map[string]float64{
        "github_activity": 0.30,
        "code_quality": 0.35,
        "challenge_score": 0.25,
        "communication": 0.10,
    }
    
    githubScore := analyzeGitHub(candidate.GithubURL)
    codeScore := getLatestCodeReview(candidate.ID)
    challengeScore := getChallengeScore(candidate.ID)
    commScore := getCommunicationScore(candidate.ID)
    
    finalScore := 
        githubScore * weights["github_activity"] +
        codeScore * weights["code_quality"] +
        challengeScore * weights["challenge_score"] +
        commScore * weights["communication"]
    
    return int(finalScore)
}
```

**RN-NETWORK-004: Featured Candidates**
- Top 20% de candidatos (score >85)
- Aparecen primero en búsqueda
- Badge "Featured" visible
- Perfil más detallado (video intro opcional)

**RN-NETWORK-005: Fees por Contratación**

**Año 1-2:**
```
Fee fijo: $3,500 USD por hire
- Simple, predecible
- Fácil de comunicar
- Menos fricción legal
```

**Año 3+ (con volumen):**
```
Fee variable: 15% annual salary (min $3,000)
Ejemplos:
- Junior ($30k/año): $3,000 (mínimo)
- Mid ($50k/año): $7,500
- Senior ($80k/año): $12,000
```

**Condiciones de Fee:**
- Solo se cobra si hire es exitoso (30 días de permanencia mínima)
- Refund 100% si candidato deja en primeros 90 días
- Tracking automático en NetworkApplication.FeeStatus

### 6.4 Flujo de Matching Company ↔ Candidate

**Búsqueda desde ATS (Recruiter):**

```
1. Recruiter en Job "Backend Developer - Go"
2. Click botón "Buscar en Red Dvra" 
   (solo visible en plan Business+)
3. Vista modal/página:
   
   ┌────────────────────────────────────────┐
   │  🔍 Red Dvra - 147 candidatos         │
   ├────────────────────────────────────────┤
   │  Filtros:                              │
   │  ☑️ Featured Only                      │
   │  Skills: [Go] [PostgreSQL] [Docker]    │
   │  Seniority: [Mid] [Senior]             │
   │  Country: [🇲🇽] [🇨🇴] [🇦🇷]            │
   │  Availability: [Available Now]         │
   │                                        │
   │  Resultados (12):                      │
   │  ────────────────────────────────────  │
   │  ⭐ Carlos Mendez - Senior Backend     │
   │     Score: 92/100 | 🇲🇽 Mexico City   │
   │     Go, PostgreSQL, AWS | 6 years     │
   │     [Ver Perfil] [Contactar]          │
   │  ────────────────────────────────────  │
   │  ⭐ Ana Rodriguez - Mid Backend        │
   │     Score: 88/100 | 🇨🇴 Bogotá        │
   │     Go, Redis, Kubernetes | 4 years   │
   │     [Ver Perfil] [Contactar]          │
   └────────────────────────────────────────┘
```

4. Click "Ver Perfil" → Modal con:
   - Bio completa
   - GitHub activity graph
   - Top repos destacados
   - Skills con proficiency
   - Dvra Score breakdown
   - Coding challenge results
   - Availability & rate

5. Click "Contactar":
   ```
   Sistema crea NetworkApplication:
   - Status: interested
   - InitiatedBy: company
   - Envía notificación a candidato
   ```

6. Candidato recibe email:
   ```
   Subject: "Azentic Sys está interesada en tu perfil"
   
   Hola Carlos,
   
   Azentic Sys (Startup fintech en Colombia) 
   está revisando tu perfil para:
   "Senior Backend Developer - Go"
   
   Detalles:
   - Remoto 100%
   - $60-80k USD/año
   - Startup serie A, 30 personas
   
   [Ver Job Completo] [Responder Interés]
   ```

7. Candidato acepta → Status: contacted
8. Recruiter coordina entrevista dentro del ATS
9. Proceso normal (Application pipeline)
10. Al marcar "Hired" → Sistema pregunta:
    ```
    "¿Este candidato viene de Red Dvra?"
    ☑️ Sí (se aplicará fee de $3,500)
    ☐ No (candidato externo)
    
    [Confirmar Contratación]
    ```

11. Sistema actualiza:
    ```
    NetworkApplication.Status = hired
    NetworkApplication.HiredAt = NOW()
    NetworkApplication.FeeAmount = 3500
    NetworkApplication.FeeStatus = pending
    ```

12. Admin de Dvra recibe notificación → Envía invoice
13. Empresa paga → FeeStatus = paid

### 6.5 Matching Algorithm (Año 3 - ML)

**Tu expertise en grafos aplicado:**

```go
// Algoritmo de matching avanzado
type MatchingEngine struct {
    graph *Graph // Grafo de relaciones
}

// Nodos del grafo
type Node struct {
    ID   string
    Type string // "job", "candidate", "skill", "company"
    Data interface{}
}

// Aristas con pesos
type Edge struct {
    From   string
    To     string
    Weight float64
    Type   string // "requires", "has_skill", "prefers"
}

func (m *MatchingEngine) CalculateMatch(job *Job, candidate *NetworkCandidate) float64 {
    // Factores de matching
    skillMatch := m.calculateSkillOverlap(job, candidate)
    seniorityMatch := m.calculateSeniorityFit(job, candidate)
    locationMatch := m.calculateLocationCompatibility(job, candidate)
    salaryMatch := m.calculateSalaryFit(job, candidate)
    
    // Pesos adaptativos basados en prioridad de empresa
    weights := m.getCompanyPreferences(job.CompanyID)
    
    matchScore := 
        skillMatch * weights.Skills +
        seniorityMatch * weights.Seniority +
        locationMatch * weights.Location +
        salaryMatch * weights.Salary
    
    return matchScore
}

// Grafos para encontrar "candidatos similares contratados"
func (m *MatchingEngine) FindSimilarHires(candidateID uint) []NetworkCandidate {
    // Encuentra candidatos con skillset similar que fueron hired
    // Usa algoritmo de random walk / PageRank en subgrafo
    similar := m.graph.RandomWalk(candidateID, maxSteps=5, filter="hired")
    return similar
}
```

### 6.6 Métricas del Marketplace

**KPIs Críticos:**

**Lado Candidatos (Supply):**
- Total candidatos registrados
- Candidatos aprobados / Featured
- Tasa de aprobación (% approved de pending)
- Tiempo promedio de evaluación
- NPS candidatos (experiencia evaluación)

**Lado Empresas (Demand):**
- % empresas usando Red Dvra activamente
- Búsquedas realizadas / mes
- Contactos iniciados / mes
- Conversion rate: Contacto → Interview → Hire
- Tiempo promedio hire (días desde contacto)
- Fee revenue / empresa / mes

**Network Effects:**
- Ratio supply/demand (candidatos/empresas activas)
- Matching success rate
- Repeat usage (empresas que contratan 2+ veces)

**Targets Año 2:**
```
Q1: 
- 400 candidatos registrados
- 100 aprobados, 30 Featured
- 10 empresas usando Red activamente
- 5 hires → $17,500 fees

Q4:
- 1,000 candidatos
- 300 aprobados, 80 Featured
- 40 empresas usando Red
- 30 hires → $105,000 fees
```

---

## 7. Evaluación Técnica

### 6.1 Integración GitHub/GitLab (FASE 1 - Crítica)

**Feature Estrella de TechRecruit**: Evaluar código real del candidato

**Implementación**:
```go
type TechnicalEvaluation struct {
    gorm.Model
    ApplicationID uint   `gorm:"not null;index"`
    CompanyID     uint   `gorm:"not null"`
    
    // GitHub/GitLab data
    RepoURL       string `gorm:"type:text"`
    CommitHash    string
    PullRequestURL string
    
    // Evaluation
    CodeQualityScore    int    // 0-100
    TestCoveragePercent float64
    LinesOfCode         int
    ReviewerUserID      *uint  // Quién revisó
    ReviewNotes         string `gorm:"type:text"`
    ReviewedAt          *time.Time
    
    Status string // pending, in_review, approved, rejected
}
```

**Flujo**:
1. Recruiter envía code challenge al candidato
2. Candidato sube solución a GitHub (repo público/privado con access)
3. Recruiter pega URL del repo en TechRecruit
4. Sistema hace análisis automático:
   - Lenguaje detectado
   - Complejidad ciclomática
   - Test coverage (si hay tests)
   - Code smells básicos
5. Recruiter revisa manualmente y deja notas
6. Decision: Approved → Mover a `offer` | Rejected → `rejected`

**APIs a integrar**:
- GitHub API v3: Repos, commits, file contents
- GitLab API: Similar
- SonarQube (opcional): Análisis de calidad de código

### 6.2 Code Challenge Templates (Fase 2)

```go
type ChallengeTemplate struct {
    gorm.Model
    CompanyID     uint
    Name          string // "Backend API Challenge"
    Description   string // Markdown con instrucciones
    Language      string // go, python, javascript, etc.
    DifficultyLevel string // junior, mid, senior
    EstimatedTime int    // minutos
    IsPublic      bool   // Compartir con comunidad
}
```

**Challenges pre-built**:
- "Build a REST API with authentication"
- "Implement rate limiter"
- "Design URL shortener"
- "Solve algorithm problem (LeetCode style)"

---

## 8. Pricing y Límites por Tier

### 8.1 Tabla de Features por Plan

| Feature | Free/Trial | Professional | Business | Enterprise |
|---------|-----------|--------------|----------|-----------|
| **Precio Anual** | $0 | $588/año ($49/mes) | $1,788/año ($149/mes) | $4,788/año ($399/mes) |
| **Precio Mensual** | $0 | $49/mes | $149/mes | $399/mes |
| **Descuento Anual** | - | 17% | 17% | 17% |
| | | | | |
| **ATS CORE** | | | | |
| Posiciones activas | 3 | 10 | Ilimitadas | Ilimitadas |
| Usuarios | 1 | 3 | Ilimitados | Ilimitados |
| Candidatos/mes | 50 | 200 | Ilimitados | Ilimitados |
| Storage (CVs) | 500MB | 10GB | 50GB | Custom |
| Pipeline Kanban | ✅ | ✅ | ✅ | ✅ |
| CV Parsing | Básico | Avanzado | Avanzado | Avanzado |
| Email Templates | 3 | Ilimitados | Ilimitados | Ilimitados |
| Analytics | Básico | Avanzado | Avanzado | Custom |
| Portal Público Jobs | ❌ | ✅ | ✅ | ✅ |
| | | | | |
| **TECH RECRUITING** | | | | |
| GitHub Integration | ❌ | ✅ | ✅ | ✅ |
| GitLab Integration | ❌ | ❌ | ✅ | ✅ |
| Code Challenges | ❌ | 5/mes | Ilimitados | Ilimitados |
| Technical Scoring | ❌ | ✅ | ✅ | ✅ |
| | | | | |
| **RED DVRA (MARKETPLACE)** | | | | |
| **Acceso a Red Dvra** | ❌ | ❌ | ✅ | ✅ |
| Búsqueda Candidatos | ❌ | ❌ | Featured Only | Todos |
| Candidatos Contactables | 0 | 0 | Ilimitados | Ilimitados |
| **Fee por Hire** | - | - | **$3,500 flat** | **$3,500 flat** |
| Match Score AI | ❌ | ❌ | Básico | Avanzado |
| Priority Support Matching | ❌ | ❌ | ❌ | ✅ |
| | | | | |
| **INTEGRACIONES** | | | | |
| Google Calendar | ✅ | ✅ | ✅ | ✅ |
| Calendly/Cal.com | ❌ | ✅ | ✅ | ✅ |
| Slack Notifications | ❌ | ❌ | ✅ | ✅ |
| API Access | ❌ | ❌ | ✅ (100 req/min) | ✅ (Custom) |
| Webhooks | ❌ | ❌ | 5 | Ilimitados |
| Custom Domain | ❌ | ❌ | ❌ | ✅ |
| SSO (SAML) | ❌ | ❌ | ❌ | ✅ |
| | | | | |
| **SOPORTE** | | | | |
| Tipo | Community | Email | Priority Email | Dedicated Slack |
| SLA Response | - | 48h | 24h | 4h |
| Onboarding | Self-service | Video call (30min) | 1-on-1 (1h) | Custom |
| Success Manager | ❌ | ❌ | ❌ | ✅ |
| SLA Uptime | - | - | 99.5% | 99.9% |

### 8.2 Propuesta de Valor por Tier

**Free/Trial (14 días Professional):**
```
Perfecto para: Validar la herramienta
Limitación clave: Sin Red Dvra, 3 jobs, 50 candidatos/mes
CTA: "Prueba gratis 14 días, no requiere tarjeta"
```

**Professional ($49/mes):**
```
Perfecto para: Startups pequeñas (10-30 empleados)
Sweet spot: ATS completo + GitHub integration
Limitación: No acceso a Red Dvra
CTA: "Todo lo que necesitas para reclutamiento interno"
```

**Business ($149/mes) - RECOMMENDED:**
```
Perfecto para: Startups en crecimiento (30-100 empleados)
Diferenciador: ACCESO A RED DVRA ⭐
Value prop: "$149/mes + solo pagas $3,500 si contratas"
ROI pitch: "1 hire de la Red = ROI de 2 años de la herramienta"
CTA: "ATS + Acceso a 500+ devs pre-evaluados"
```

**Enterprise ($399/mes):**
```
Perfecto para: Empresas 100+ empleados, agencias
Diferenciador: SSO, API, Dedicated support
Value prop: "Herramienta enterprise + red de talento"
CTA: "Hablemos de tus necesidades"
```

### 8.3 Estrategia de Upsell

**Journey de Conversión:**

```
1. Free Trial (14 días Professional)
   ↓
2. Conversion a Professional ($49/mes)
   - Pain: "Necesito más de 10 jobs"
   - Upsell: Business por +$100/mes
   ↓
3. Discovery de Red Dvra (en Business)
   - Tooltip: "🎉 Ahora puedes buscar en Red Dvra"
   - Tutorial: "Mira cómo funciona" (video 2 min)
   ↓
4. Primera búsqueda en Red Dvra
   - Email follow-up: "¿Encontraste candidatos interesantes?"
   ↓
5. Primer hire de Red Dvra
   - Invoice automático: $3,500
   - NPS: "¿Qué tal la experiencia?"
   - Cross-sell: "¿Necesitas más hires? Tenemos 500+ candidatos"
```

**Incentivos de Retención:**

```
- Descuento anual: 17% off (2 meses gratis)
- Loyalty bonus: Hire 3+ de Red Dvra → 10% descuento en fees
- Referral: Trae otra empresa → $100 crédito o 1 mes gratis
- Enterprise upgrade: Si contratan 5+ al año → Dedicated CSM gratis
```

### 8.4 Enforcement de Límites (Middleware)

```go
// middleware/tier_limits.go
func EnforceTierLimits() gin.HandlerFunc {
    return func(c *gin.Context) {
        companyID := c.GetUint("company_id")
        company := getCompany(companyID)
        action := c.Request.Method + " " + c.Request.URL.Path
        
        // Feature gating por plan
        switch {
        // Red Dvra: Solo Business+
        case strings.Contains(action, "/network-candidates"):
            if !isPlanBusinessOrHigher(company.PlanTier) {
                c.JSON(403, gin.H{
                    "error": "Red Dvra requires Business or Enterprise plan",
                    "current_plan": company.PlanTier,
                    "upgrade_url": "/pricing",
                    "message": "Unlock 500+ pre-evaluated developers",
                })
                c.Abort()
                return
            }
        
        // Crear Job: Validar límite
        case action == "POST /api/v1/jobs":
            activeJobs := countActiveJobs(companyID)
            limit := getJobLimit(company.PlanTier)
            
            if activeJobs >= limit {
                c.JSON(403, gin.H{
                    "error": "Job limit reached",
                    "limit": limit,
                    "current": activeJobs,
                    "upgrade_to": getNextTier(company.PlanTier),
                    "upgrade_url": "/pricing",
                })
                c.Abort()
                return
            }
        
        // Candidatos por mes
        case action == "POST /api/v1/candidates":
            candidatesThisMonth := countCandidatesThisMonth(companyID)
            limit := getCandidateLimit(company.PlanTier)
            
            if candidatesThisMonth >= limit {
                c.JSON(403, gin.H{
                    "error": "Monthly candidate limit reached",
                    "limit": limit,
                    "resets_in": getDaysUntilNextMonth(),
                    "upgrade_url": "/pricing",
                })
                c.Abort()
                return
            }
        }
        
        c.Next()
    }
}

func getJobLimit(planTier string) int {
    limits := map[string]int{
        "free": 3,
        "professional": 10,
        "business": 999999,
        "enterprise": 999999,
    }
    return limits[planTier]
}

func getCandidateLimit(planTier string) int {
    limits := map[string]int{
        "free": 50,
        "professional": 200,
        "business": 999999,
        "enterprise": 999999,
    }
    return limits[planTier]
}

func isPlanBusinessOrHigher(planTier string) bool {
    return planTier == "business" || planTier == "enterprise"
}
```

### 8.5 Modelo de Fees Marketplace

**Año 1-2: Fee Fijo Simple**
```go
const (
    NetworkHireFeeFlat = 3500 // USD
)

// Al marcar hire de NetworkApplication
func calculateNetworkFee(application *NetworkApplication) int {
    return NetworkHireFeeFlat
}
```

**Año 3+: Fee Variable (con volumen)**
```go
const (
    NetworkHireFeePercent = 0.15  // 15%
    NetworkHireFeMin      = 3000  // USD mínimo
)

func calculateNetworkFee(application *NetworkApplication, annualSalary int) int {
    percentageFee := int(float64(annualSalary) * NetworkHireFeePercent)
    if percentageFee < NetworkHireFeeMin {
        return NetworkHireFeeMin
    }
    return percentageFee
}

// Ejemplos:
// Junior $30k/año → $3,000 (mínimo)
// Mid $50k/año → $7,500
// Senior $80k/año → $12,000
// Staff $120k/año → $18,000
```

**Garantía y Refund Policy:**
```go
type FeeGuarantee struct {
    CandidateLeftBefore30Days  bool // 100% refund
    CandidateLeftBefore90Days  bool // 50% refund o replacement gratis
    CandidateLeftAfter90Days   bool // Sin refund (hire exitoso)
}

// Tracking automático
func checkFeeGuarantee(application *NetworkApplication) *FeeGuarantee {
    daysSinceHire := time.Since(application.HiredAt).Hours() / 24
    
    return &FeeGuarantee{
        CandidateLeftBefore30Days: daysSinceHire < 30,
        CandidateLeftBefore90Days: daysSinceHire < 90,
        CandidateLeftAfter90Days: daysSinceHire >= 90,
    }
}
```

### 8.6 Trial y Conversión

**RN-TRIAL-001: Duración y Scope**
- Trial: 14 días del plan **Professional** (no Business)
- Razón: Validar ATS core sin revelar Red Dvra gratis
- Después del trial → Downgrade a Free (soft landing)
- No se bloquea acceso, solo features premium

**RN-TRIAL-002: Conversion Tactics**

**Día 3 de Trial:**
```
Email: "¿Cómo va tu experiencia con Dvra?"
CTA: Agendar demo 15 min
```

**Día 7 de Trial:**
```
In-app banner: "7 días restantes de trial"
+ Tooltip en botón "Red Dvra": 
  "Upgrade a Business para acceder a 500+ devs pre-evaluados"
```

**Día 12 de Trial:**
```
Email: "Tu trial termina en 2 días"
Offer: "Convierte hoy → 20% descuento primer mes"
```

**Día 14 - Final:**
```
Modal: "Tu trial ha terminado"
Options:
- [Upgrade a Professional - $49/mes]
- [Upgrade a Business - $149/mes] ← RECOMMENDED
- [Continuar con Free]

Highlight Business: 
"💡 1 hire de Red Dvra paga 2 años de la herramienta"
```

**RN-TRIAL-003: Feature Gating UX**

Cuando usuario Free/Professional intenta acceder a Red Dvra:
```
┌────────────────────────────────────────┐
│  🔒 Red Dvra (Requires Business Plan) │
├────────────────────────────────────────┤
│                                        │
│  Unlock access to 500+ pre-evaluated  │
│  developers from LATAM                 │
│                                        │
│  ✓ Technical screening completed       │
│  ✓ GitHub/GitLab verified             │
│  ✓ Coding challenges passed            │
│  ✓ Only pay $3,500 when you hire      │
│                                        │
│  [Upgrade to Business - $149/mes]     │
│  [Learn More]                          │
└────────────────────────────────────────┘
```

---

## 9. Flujos de Usuario Críticos

### 8.1 Onboarding de Nueva Empresa

**Flujo Happy Path**:
```
1. Usuario visita techrecruit.app
2. Click "Start Free Trial" → Sign up form
3. Ingresa:
   - Email
   - Password
   - Company Name (auto-genera slug)
   - Timezone
4. Sistema crea:
   - Company (PlanTier=free, TrialEndsAt=+14 días)
   - User (email, password_hash)
   - Membership (UserID, CompanyID, Role=admin, Status=active)
5. Redirige a Dashboard
6. Tutorial interactivo:
   - "Crea tu primer Job"
   - "Agrega un candidato"
   - "Configura integración GitHub"
7. Email de bienvenida con recursos
```

### 8.2 Flujo de Reclutamiento End-to-End

**Caso: Startup contrata Backend Developer**

```
1. Admin crea Job "Backend Developer - Go"
   - Title, Description, Location
   - Status = draft

2. Admin publica Job
   - Status: draft → active
   - Asigna Recruiter (Juan)

3. Candidatos aplican (3 maneras):
   a) Portal público: jobs.techrecruit.app/azentic-sys/backend-dev
   b) Recruiter agrega manualmente (sourcing LinkedIn)
   c) Email parsing: Forward CV to jobs@techrecruit.app

4. Recruiter (Juan) hace screening
   - Revisa CVs en Dashboard
   - Llama candidato → Notas en sistema
   - Rating inicial: 4/5
   - Move stage: applied → screening

5. Evaluación técnica
   - Juan envía code challenge (template pre-built)
   - Candidato sube solución a GitHub
   - Juan pega URL en TechRecruit
   - Sistema analiza código automáticamente
   - Juan revisa y aprueba
   - Move stage: screening → technical

6. Oferta y cierre
   - Hiring Manager revisa candidato
   - Admin extiende oferta
   - Move stage: technical → offer
   - Candidato acepta 🎉
   - Move stage: offer → hired
   - Sistema setea HiredAt = NOW
   - Job.Status = closed (opcional si solo 1 vacante)

7. Analytics
   - Time-to-hire: 15 días (vs industry avg 45 días)
   - Conversion rate: 8% (6 hired de 75 applied)
   - Source más efectivo: LinkedIn (50% de hired)
```

### 8.3 Flujo de Invitación de Usuario

```
1. Admin va a Settings → Team
2. Click "Invite User"
3. Ingresa:
   - Email: maria@example.com
   - Role: Recruiter
4. Sistema:
   - Busca User con ese email
     - Si existe: Crea Membership con status=pending
     - Si NO existe: Crea User + Membership pending
   - Envía email de invitación con token
5. María recibe email
   - Click link con token
   - Si no tiene cuenta: Sign up
   - Si tiene cuenta: Login
6. Al autenticarse, Membership.Status = active
7. María ve Dashboard con datos de la empresa
```

---

## 9. Integraciones Estratégicas

### 9.1 Prioridad 1 (MVP - Q1 2025)

#### GitHub/GitLab
**Value Prop**: Evaluar código real del candidato
**Implementación**:
- OAuth flow para conectar cuenta
- Webhooks para detectar nuevos commits
- GraphQL API para análisis de repos

#### Email (SendGrid/AWS SES)
**Use Cases**:
- Bienvenida
- Notificaciones de stage change
- Invitaciones de usuario
- Weekly digest (candidatos nuevos)

#### Storage (AWS S3)
**Use Cases**:
- CVs (PDF, DOCX)
- Cover letters
- Documentos de candidato

### 9.2 Prioridad 2 (Q2 2025)

#### Calendly/Cal.com
**Value Prop**: Scheduling de entrevistas sin fricción
**Implementación**:
- Embed widget en perfil de candidato
- Auto-sync con Google Calendar
- Envío automático de invites

#### Job Boards (LinkedIn, GetOnBoard, GetonBrd)
**Value Prop**: Publicar jobs en múltiples plataformas
**Implementación**:
- API de job posting
- Sync bidireccional de aplicaciones

### 9.3 Prioridad 3 (Q3-Q4 2025)

#### AI/LLM (OpenAI GPT-4)
**Use Cases**:
- Resume parsing inteligente
- Generación de job descriptions
- Scoring automático de fit candidato-job
- Sugerencias de preguntas de entrevista

#### Slack/Microsoft Teams
**Use Cases**:
- Notificaciones en tiempo real
- Comandos slash (/techrecuit search backend)
- Aprobaciones de workflow

#### Video Interview (Zoom, Google Meet API)
**Use Cases**:
- Scheduling integrado
- Recording y transcription
- AI analysis de entrevistas

---

## 12. Roadmap de Implementación

> **Estrategia de 3 Años: De ATS a Plataforma Híbrida**  
> 90% ATS → 60/40 Split → 50/50 Full Híbrido

---

### 🏗️ Fase 0: Fundación (COMPLETADO ✅)

**Duración**: Completado  
**Esfuerzo**: Backend + arquitectura básica

**Logros:**
- [x] Arquitectura multi-tenant con CompanyID scoping
- [x] Modelos de datos core (Company, User, Job, Candidate, Application, Membership, Role)
- [x] CRUD completo de entidades principales
- [x] JWT authentication service con claims personalizados
- [x] Docker Compose setup (PostgreSQL 16)
- [x] Go 1.24.0 + Gin + GORM stack

**Estado Técnico**: Backend foundation sólido, listo para features de producto.

---

## AÑO 1: VALIDACIÓN ATS (90% ATS + 10% Marketplace Prep)

### 📅 Q1 2025: MVP y Beta (Enero - Marzo)

**Objetivo**: Lanzar MVP funcional, conseguir 5 clientes beta pagantes  
**Revenue Goal**: $500-1,000 MRR  
**Esfuerzo**: 95% ATS, 5% marketplace prep

#### Mes 1 - Enero: Autenticación y Autorización

**Backend:**
- [ ] Middleware RequireAuth() completo con JWT validation
- [ ] Middleware RequireRole(minLevel) para permisos
- [ ] Endpoints auth:
  - `POST /api/v1/auth/register` - Registro de empresa
  - `POST /api/v1/auth/login` - Login con email/password
  - `POST /api/v1/auth/refresh` - Refresh token
  - `GET /api/v1/auth/me` - Usuario actual
- [ ] Password reset flow con email token
- [ ] Email verification (SendGrid basic)

**Frontend (mínimo):**
- [ ] Login/Register páginas
- [ ] Protected routes con middleware
- [ ] Error handling UX

**Goal**: Autenticación segura funcionando end-to-end

#### Mes 2 - Febrero: Features Core ATS

**Jobs & Applications:**
- [ ] Portal público: `jobs.dvra.app/:company/:slug`
  - Vista de job individual
  - Formulario de aplicación público
  - Thank you page
- [ ] CRUD jobs con validaciones de negocio
- [ ] Pipeline kanban view (drag & drop)
- [ ] Move application entre stages con validaciones

**File Management:**
- [ ] Integración AWS S3 para CVs
- [ ] CV upload con validación (PDF, DOCX, max 5MB)
- [ ] Presigned URLs para descarga segura

**Notifications:**
- [ ] Integración SendGrid/AWS SES
- [ ] Email templates básicos:
  - Bienvenida
  - Application received (candidato)
  - Stage changed
  - Password reset

**Goal**: Flujo básico de reclutamiento funcionando

#### Mes 3 - Marzo: Monetización y Polish

**Pricing & Limits:**
- [ ] Tier enforcement middleware
  - Job limits (free=3, professional=10)
  - Candidate limits por mes
  - Storage limits
- [ ] Feature gates (Red Dvra para Business+)
- [ ] In-app upgrade prompts

**Payments:**
- [ ] Integración Stripe (primary) + MercadoPago (LATAM)
- [ ] Subscription management
- [ ] Invoice generation automático
- [ ] Trial de 14 días del plan Professional

**UX Polish:**
- [ ] Dashboard con métricas:
  - Jobs activos
  - Candidatos por stage
  - Time-to-hire promedio
- [ ] Settings de empresa:
  - Logo upload (S3)
  - Timezone config
  - Plan upgrade UI
- [ ] User invitation flow
- [ ] Onboarding tutorial interactivo (Product Tour)

**Marketing:**
- [ ] Landing page marketing (dvra.app)
- [ ] Pricing page
- [ ] Docs básicos

**Marketplace Prep (10% tiempo):**
- [ ] Landing page `talento.dvra.com` (static, info only)
- [ ] Formulario de registro candidatos (guarda en DB, no procesa aún)
- [ ] Email auto-respuesta: "Tu perfil está en revisión"

**🎯 Q1 Goal**: 5-10 clientes beta, $500-1k MRR, **fundación para escala**

---

### 📅 Q2 2025: Lanzamiento Público (Abril - Junio)

**Objetivo**: Lanzar públicamente, llegar a 15-20 clientes  
**Revenue Goal**: $2,000-3,000 MRR  
**Esfuerzo**: 90% ATS, 10% marketplace prep

#### Features de Producto:

**Tech Recruiting (Diferenciador #1):**
- [ ] Integración GitHub OAuth
- [ ] Pegar repo URL → Auto-analysis:
  - Lenguaje detectado
  - LoC, commits, contributors
  - Detección básica de tests
- [ ] Manual code review con scoring (0-100)
- [ ] Code challenge templates (3-5 pre-built)

**Analytics & Reporting:**
- [ ] Funnel de conversión por job
- [ ] Time-to-hire por job y promedio
- [ ] Source tracking effectiveness
- [ ] Export reports (CSV)

**Collaboration:**
- [ ] Comments en applications
- [ ] @mentions de teammates
- [ ] Activity log (audit trail básico)

**UX Improvements:**
- [ ] Bulk actions (rechazar múltiples candidatos)
- [ ] Tags customizables por empresa
- [ ] Filtros avanzados (date range, source, rating)
- [ ] Search global (candidatos, jobs)

**Marketplace Prep (10% tiempo):**
- [ ] Proceso manual de evaluación de candidatos:
  - Admin panel interno para revisar registros
  - Scoring manual (GitHub review, challenge envío)
  - Approval/rejection con feedback email
- [ ] Meta: Evaluar **2-3 candidatos/semana** manualmente
- [ ] Target: **50 candidatos registrados**, 15-20 evaluados

**🎯 Q2 Goal**: 15-20 clientes, $2-3k MRR, primeros candidatos evaluados en Red

---

### 📅 Q3 2025: Escala ATS (Julio - Septiembre)

**Objetivo**: Escalar a 30-35 clientes, features enterprise  
**Revenue Goal**: $4,000-5,000 MRR  
**Esfuerzo**: 85% ATS, 15% marketplace prep

#### Features Enterprise:

**Integraciones:**
- [ ] Google Calendar sync
- [ ] Calendly/Cal.com embed
- [ ] Slack notifications:
  - New application
  - Stage changes
  - Hired notifications
- [ ] Job boards posting (LinkedIn, GetOnBoard)

**Team Features:**
- [ ] Hiring Manager role refinement
- [ ] Per-job permission controls
- [ ] Team performance dashboard

**AI/Automation:**
- [ ] CV parsing con AWS Textract
  - Auto-extract: name, email, phone, experience
  - Auto-populate candidate fields
- [ ] Email automation:
  - Auto-send rejection emails
  - Auto-reminder para interviews
  - Weekly digest para recruiters

**API (Beta):**
- [ ] REST API pública (read-only para empezar)
- [ ] Documentación con Swagger
- [ ] API keys management

**Marketplace Prep (15% tiempo):**
- [ ] Sistema de tags/skills para NetworkCandidates
- [ ] Scoring algorithm básico (manual inputs)
- [ ] Perfil público de candidato (preview para empresas Business)
- [ ] Meta: **100 candidatos registrados**, 30-40 evaluados
- [ ] Featured candidates (top 20%, score >85)

**🎯 Q3 Goal**: 30-35 clientes, $4-5k MRR, **100 candidatos en Red**

---

### 📅 Q4 2025: Cierre Año 1 (Octubre - Diciembre)

**Objetivo**: 50 clientes, preparar integración marketplace  
**Revenue Goal**: $5,000 MRR ($60k ARR)  
**Esfuerzo**: 80% ATS features, **20% marketplace integration prep**

#### ATS Advanced Features:

**Compliance & Security:**
- [ ] GDPR compliance tools:
  - Candidate data export
  - Right to be forgotten
  - Consent tracking
- [ ] Audit logs completos
- [ ] 2FA para admin users

**Advanced Analytics:**
- [ ] Hiring manager effectiveness
- [ ] Recruiter performance comparison
- [ ] Predictive time-to-fill
- [ ] Cost-per-hire tracking

**Customization:**
- [ ] Custom pipeline stages
- [ ] Custom fields para candidatos
- [ ] Email template editor
- [ ] Branding completo (custom domain para Enterprise)

**Marketplace Integration Prep (20% tiempo - CRÍTICO):**

**Backend:**
- [ ] **Modelos marketplace completos**:
  - NetworkCandidate
  - NetworkCandidateSkill
  - NetworkApplication
  - NetworkEvaluation
- [ ] Migrations para tablas marketplace
- [ ] CRUD completo NetworkCandidate (internal admin)

**Frontend (Internal Admin Only):**
- [ ] Panel interno de evaluación
- [ ] Lista de pending candidates
- [ ] Evaluation workflow:
  - GitHub analysis manual → Score
  - Challenge assignment
  - Approval/rejection
- [ ] Dashboard de métricas:
  - Candidates pending review
  - Approval rate
  - Average evaluation time

**Red Dvra Growth:**
- [ ] Marketing activo para desarrolladores:
  - Ads en LinkedIn LATAM
  - Partnerships con bootcamps
  - Twitter/X presence
- [ ] Meta: **200 candidatos registrados**, **50-60 evaluados y aprobados**
- [ ] Featured candidates visible en preview mode

**Pricing Prep:**
- [ ] Research precio de fees (validar $3,500 flat vs 15%)
- [ ] Legal: Términos de servicio para marketplace
- [ ] Contracts templates para fees

**🎯 Q4 Goal**: **50 clientes**, **$5k MRR** ($60k ARR), **200 candidatos en Red**, infraestructura marketplace lista

**🏁 Fin Año 1 Metrics:**
```
✅ 50 empresas pagantes (40 Professional, 8 Business, 2 Enterprise)
✅ $5,000 MRR = $60,000 ARR
✅ 200 candidatos registrados en Red Dvra
✅ 50-60 candidatos aprobados (Featured: 15-20)
✅ 0 hires del marketplace (aún no lanzado públicamente)
✅ Churn: <5% mensual
✅ NPS: >40
✅ Break-even: Mes 9-10
```

---

## AÑO 2: INTEGRACIÓN MARKETPLACE (60% ATS + 40% Marketplace)

### 📅 Q1 2026: Launch Red Dvra para Empresas (Enero - Marzo)

**Objetivo**: Integrar marketplace en ATS, validar primer hire  
**Revenue Goal**: $10k MRR ($6k SaaS + $4k fees iniciales)  
**Esfuerzo**: **50% marketplace, 50% ATS**

#### Marketplace Launch (Feature Completo):

**Vista "Red Dvra" en ATS:**
- [ ] Nueva sección en sidebar (solo Business+ plan)
- [ ] Búsqueda y filtros de NetworkCandidates:
  - Skills (multi-select)
  - Seniority level
  - Country/timezone
  - Availability status
  - Score range
  - Featured only toggle
- [ ] Perfil detallado de candidato:
  - Bio, experience, skills
  - GitHub activity visualization
  - DvraScore breakdown
  - Coding challenge results
  - Availability & hourly rate
- [ ] Botón "Contactar" → Create NetworkApplication

**NetworkApplication Workflow:**
- [ ] Status: interested → contacted → interviewing → offer → hired/rejected
- [ ] Email automático a candidato (company interested)
- [ ] Candidate dashboard (simple):
  - Ver empresas interesadas
  - Accept/decline interest
  - Track interview status
- [ ] Fee tracking:
  - Marcar "Hired from Red Dvra" checkbox
  - Auto-create invoice ($3,500 flat fee)
  - Payment tracking (pending → paid)

**Tier Enforcement:**
- [ ] Feature gate completo:
  - Professional plan → No access (upgrade prompt)
  - Business plan → Access Featured only
  - Enterprise plan → Access all candidates
- [ ] Upgrade flow desde paywall

**Candidate Experience:**
- [ ] Portal público `talento.dvra.com` (mejorado):
  - Registro completo
  - Upload CV
  - GitHub connect
  - Skills self-assessment
- [ ] Email notifications para candidatos:
  - Empresa interesada
  - Interview invitations
  - Hired notification

**Legal & Contracts:**
- [ ] Términos de servicio marketplace
- [ ] Fee agreement templates
- [ ] Refund policy implementation (90 días)

**🎯 Q1 Goal**: Marketplace 100% funcional, **primeros 3-5 hires** de Red Dvra (**$10-17k fees**), **75 clientes SaaS**

---

### 📅 Q2 2026: Optimización y Escala Marketplace (Abril - Junio)

**Objetivo**: Escalar hires del marketplace, automatizar evaluación  
**Revenue Goal**: $15k MRR ($10k SaaS + $5k fees)  
**Esfuerzo**: 40% ATS, 60% marketplace

#### Auto-Evaluation (Game Changer):

**GitHub Auto-Analysis:**
- [ ] Integración GitHub API (automático para cada candidate)
- [ ] Metrics recolectados:
  - Total repos, commits (last 6 months)
  - Languages used (weighted by LoC)
  - Stars, forks, watchers en top repos
  - Contribution frequency
  - Code review participation
- [ ] Auto-score GitHub activity (0-100)

**Automated Challenges:**
- [ ] Challenge auto-send al registro
- [ ] Time-tracking integrado
- [ ] Auto-run tests on submission
- [ ] Pass/fail scoring automático
- [ ] Manual review solo para edge cases

**Scoring Automation:**
- [ ] DvraScore calculation automático:
  - GitHub: 30%
  - Code quality: 35%
  - Challenge: 25%
  - Communication: 10% (still manual)
- [ ] Auto-approve candidates score >70
- [ ] Auto-feature candidates score >85
- [ ] Human review solo para borderline (60-70)

**Efficiency Gains:**
- Manual: 2-3 candidates/week → **Auto: 20-30 candidates/week**

**Marketplace Improvements:**
- [ ] Match score calculation (basic):
  - Skill overlap %
  - Seniority fit
  - Timezone compatibility
- [ ] Sort candidates by match score
- [ ] "Similar candidates" suggestions
- [ ] Candidate video intro (optional, 1-2 min)

**Growth:**
- [ ] Red Dvra candidate marketing:
  - SEO content (salarios, remote jobs)
  - LinkedIn/Twitter ads targeted
  - University partnerships
- [ ] Referral program:
  - Candidato refiere candidato → $50 gift card si ambos aprueban
  - Empresa refiere empresa → 1 mes gratis SaaS

**🎯 Q2 Goal**: **10-15 hires** marketplace (**$35-50k fees acumulados**), **500 candidatos** en Red, **100 clientes SaaS**, **$15k MRR**

---

### 📅 Q3 2026: Network Effects (Julio - Septiembre)

**Objetivo**: Escalar ambos lados del marketplace  
**Revenue Goal**: $18k MRR ($12k SaaS + $6k fees)  
**Esfuerzo**: 35% ATS, 65% marketplace

#### Two-Sided Growth:

**Supply (Candidatos):**
- [ ] Gamification:
  - Badges (GitHub verified, Top 10%, Fast responder)
  - Leaderboard (score ranking)
  - Profile completion %
- [ ] Candidate dashboard mejorado:
  - Analytics: profile views, interests
  - Interview prep resources
  - Salary benchmarks LATAM
- [ ] Community building:
  - Blog con tips de carrera
  - Newsletter semanal con jobs destacados
- [ ] Target: **1,000 candidatos registrados**, 300 aprobados

**Demand (Empresas):**
- [ ] Upsell agresivo Professional → Business:
  - In-app tooltips: "¿Buscas senior backend? Tenemos 50 en Red Dvra"
  - Webinars: "Cómo contratar remoto LATAM"
  - Case studies de hires exitosos
- [ ] Enterprise sales team (1 SDR):
  - Outbound a startups Series A/B
  - Pitch: "ATS + talento pre-evaluado = ROI instantáneo"
- [ ] Target: **150 clientes SaaS**, 40% en Business+

**Matching Algorithm (v1):**
- [ ] Factores de scoring:
  - Skill overlap (TF-IDF de job description vs candidate skills)
  - Seniority match (junior/mid/senior mapping)
  - Location preference
  - Salary range fit
- [ ] Proactive matching notifications:
  - "3 nuevos candidatos match para tu job Backend Senior"
  - Weekly digest de top matches

**Conversion Optimization:**
- [ ] A/B testing en Red Dvra UI
- [ ] Funnel analytics: Interest → Contacted → Interviewed → Hired
- [ ] Bottleneck identification y fixes

**🎯 Q3 Goal**: **20-25 hires** (**$70-87k fees acumulados**), **1,000 candidatos**, **150 clientes**, **$18k MRR**

---

### 📅 Q4 2026: Cierre Año 2 (Octubre - Diciembre)

**Objetivo**: $20k+ MRR, marketplace como revenue driver validado  
**Revenue Goal**: $20k MRR ($14k SaaS + $6k fees)  
**Esfuerzo**: 30% ATS, 70% marketplace

#### Marketplace Maturity:

**Advanced Features:**
- [ ] Saved searches para recruiters
- [ ] Auto-alerts: "Nuevo candidato Go senior en Colombia"
- [ ] Candidate shortlists compartidas en equipo
- [ ] Interview scheduling integrado (Calendly direct book)

**Trust & Safety:**
- [ ] Background checks integration (Checkr)
- [ ] Identity verification para top candidates
- [ ] Dispute resolution process
- [ ] Refund automation (candidato deja <90 días)

**Analytics & Insights:**
- [ ] Marketplace dashboard para admins:
  - Conversion funnel
  - Average time-to-hire (Red vs normal)
  - Revenue por cliente Business+
  - Most requested skills
- [ ] Candidate analytics:
  - Profile views
  - Interest rate %
  - Interview-to-hire ratio

**Scaling Evaluation:**
- [ ] Batch processing de registros (50+ candidates/week)
- [ ] Distributed evaluation (multiple evaluators)
- [ ] Quality control sampling (10% manual review)

**Enterprise Features:**
- [ ] SSO (SAML) para Enterprise plan
- [ ] Multi-workspace para agencias
- [ ] Dedicated account manager
- [ ] Custom SLA (24/7 support)

**Growth Marketing:**
- [ ] Content marketing:
  - Blog: "State of LATAM Tech Salaries 2026"
  - Ebook: "Remote Hiring Playbook"
- [ ] PR: Publicar casos de éxito en TechCrunch LATAM
- [ ] Partnerships: Integrar con Platzi, Coderhouse

**🎯 Q4 Goal**: **$20-22k MRR**, **150+ clientes SaaS**, **30 hires** marketplace en el trimestre (**$216k ARR total**)

**🏁 Fin Año 2 Metrics:**
```
✅ 150 empresas pagantes (70 Professional, 60 Business, 20 Enterprise)
✅ $18,000 MRR SaaS = $216k ARR SaaS
✅ ~50-60 hires anuales del marketplace = ~$175-210k fees
✅ ARR Total Blended: ~$390-420k
✅ 1,000+ candidatos en Red, 300-350 aprobados
✅ 80+ Featured candidates
✅ LTV:CAC mejora a 9:1 (sticky marketplace)
✅ Churn baja a 3% (Business+ churn <2%)
```

---

## AÑO 3: FULL HÍBRIDO (50% ATS + 50% Marketplace)

### 📅 Q1 2027: ML Matching & Scale (Enero - Marzo)

**Objetivo**: Matching algorithm avanzado, escala agresiva  
**Revenue Goal**: $55k MRR ($35k SaaS + $20k fees)  
**Esfuerzo**: **50/50 split**

#### ML Matching Algorithm (Tu Expertise en Grafos):

**Graph-Based Matching:**
- [ ] Construir grafo de conocimiento:
  - Nodos: Job, Candidate, Skill, Company, Industry
  - Aristas: requires, has_skill, hired_for, similar_to
  - Pesos: based en historical hire success
- [ ] Algoritmos:
  - Random walk para "candidates similares contratados"
  - PageRank para "skills más demandados"
  - Collaborative filtering (empresa X contrató Y → recomendar Z)
- [ ] Features:
  - Score 0-100 de match (no solo skills, sino cultural fit, seniority progression)
  - Explainability: "Match 92% porque: Skills 90%, Timezone 100%, Salary fit 85%"
  - Top 10 candidates auto-ranked por match

**Predictive Analytics:**
- [ ] Time-to-hire predictor por job characteristics
- [ ] Candidate drop-off prediction (re-engagement emails)
- [ ] Salary negotiation recommendations

**Automation:**
- [ ] Auto-suggest candidates para nuevos jobs
- [ ] Auto-invite top matches cuando candidate se aprueba
- [ ] Smart notifications (solo high-match candidates)

#### ATS Innovation:

**AI Features:**
- [ ] GPT-4 integration:
  - Auto-generate job descriptions
  - Interview questions suggestions
  - Candidate summary generation (from CV)
- [ ] Resume screening automation:
  - Auto-score resume vs job description
  - Flag potential red flags
  - Highlight key qualifications

**Mobile App (React Native):**
- [ ] Recruiter mobile app:
  - Review candidates on-the-go
  - Move stages
  - Push notifications
- [ ] Candidate mobile app:
  - Profile management
  - Interview reminders
  - Chat con recruiters (futuro)

**Growth:**
- [ ] International expansion (México, Argentina, Brasil)
- [ ] Multi-language support (ES, PT, EN)
- [ ] Localized pricing

**🎯 Q1 Goal**: **$55k MRR**, **250 clientes**, **25 hires** marketplace en trimestre

---

### 📅 Q2 2027: Platform Ecosystem (Abril - Junio)

**Objetivo**: Convertir en plataforma, ecosystem play  
**Revenue Goal**: $65k MRR  
**Esfuerzo**: 50/50

#### API & Integrations:

**Public API (Full):**
- [ ] RESTful API completa (read + write)
- [ ] Webhooks:
  - application.created
  - application.stage_changed
  - candidate.hired
  - network_application.hired
- [ ] Rate limiting: 1000 req/min (Business), unlimited (Enterprise)
- [ ] Zapier integration

**Partner Ecosystem:**
- [ ] App marketplace:
  - Background checks (Checkr)
  - Video interviews (Loom, Vidyard)
  - Coding challenges (HackerRank, Codility)
  - Contract management (DocuSign)
- [ ] Revenue share con partners (10-20%)

**Developer Portal:**
- [ ] Docs interactivos (Postman collections)
- [ ] SDKs (Python, Node.js, Go)
- [ ] Sandbox environment

#### Advanced Marketplace:

**Candidate Features:**
- [ ] Portfolio showcase (Figma embeds, GitHub repos destacados)
- [ ] Salary expectations visibility (anonymous ranges)
- [ ] Availability calendar (days available for interview)
- [ ] Candidate referrals tracking

**Pricing Evolution:**
- [ ] Transición a fee variable (15% salary, min $3,000):
  - Junior: $3,000 flat
  - Mid: $6-8k (15% of $40-50k)
  - Senior: $10-15k (15% of $70-100k)
- [ ] Incentivos para repeat customers:
  - 3+ hires/año → 10% descuento en fees
  - Enterprise plan → 12% fee (vs 15%)

**Community:**
- [ ] Foro de desarrolladores (networking)
- [ ] Virtual events (webinars, tech talks)
- [ ] Dvra Summit (anual conference)

**🎯 Q2 Goal**: **$65k MRR**, **300 clientes**, **30 hires** marketplace

---

### 📅 Q3 2027: Enterprise Scale (Julio - Septiembre)

**Objetivo**: Dominar segmento enterprise, revenue optimization  
**Revenue Goal**: $70k MRR  

#### Enterprise Focus:

**Features:**
- [ ] Multi-company support (holdings con múltiples subsidiarias)
- [ ] Advanced permissions (custom roles)
- [ ] Compliance reports (EEO, diversity metrics)
- [ ] SLA guarantees (99.9% uptime, 4h response)
- [ ] Dedicated CSM + Slack channel

**Sales:**
- [ ] Enterprise sales team (2-3 SDRs + 1 AE)
- [ ] Custom pricing (>$500/mes)
- [ ] Annual contracts con descuentos (20% off)

**Marketplace Premium:**
- [ ] White-glove service:
  - Dedicated talent scout
  - Custom candidate sourcing
  - Interview coordination full-service
- [ ] Priority access a Featured candidates
- [ ] Exclusive candidates (no mostrar a otros)

**🎯 Q3 Goal**: **$70k MRR**, **350 clientes** (30+ Enterprise), **35 hires**

---

### 📅 Q4 2027: Cierre Año 3 - $900k ARR (Octubre - Diciembre)

**Objetivo**: $75k+ MRR, milestone $900k ARR  
**Revenue Goal**: $75k MRR ($50k SaaS + $25k fees)  

#### Final Push:

**Product:**
- [ ] Voice of customer features (top 10 requested)
- [ ] Performance optimization (sub-second page loads)
- [ ] A11y compliance (WCAG 2.1 AA)
- [ ] Security audit (SOC 2 Type II)

**Growth:**
- [ ] Referral program explosion:
  - Empresa refiere empresa → $500 credits
  - Virality loop
- [ ] Content dominance:
  - SEO #1 para "ATS LATAM", "contratar remoto LATAM"
  - YouTube channel con tutorials
- [ ] PR push: "La startup que cambió el reclutamiento tech en LATAM"

**Marketplace Excellence:**
- [ ] 2,000+ candidatos activos
- [ ] 500+ aprobados, 150+ Featured
- [ ] Average time-to-hire: 12 días (vs industry 45)
- [ ] 95% satisfaction rate empresas
- [ ] 90% satisfaction rate candidatos

**🎯 Q4 Goal**: **$75k MRR** = **$900k ARR**, **400 clientes**, **100+ hires** marketplace anuales

**🏁 Fin Año 3 Metrics - ÉXITO:**
```
✅ 400 empresas pagantes
   - 180 Professional ($49) = $8,820/mes
   - 170 Business ($149) = $25,330/mes
   - 50 Enterprise ($399+) = $19,950/mes
   = ~$54k MRR SaaS = $648k ARR SaaS

✅ Marketplace Revenue:
   - 100 hires/año × $3,500-12,000 avg = $650k fees
   - ~$54k/mes promedio = $25k fees MRR efectivo

✅ Blended MRR: $75k = $900k ARR 🎉

✅ 2,000+ candidatos en Red Dvra
✅ 500-600 candidatos aprobados
✅ 150+ Featured candidates

✅ Unit Economics:
   - LTV: $6,600 (ARPU $200/mes, 33 meses lifetime, 3% churn)
   - CAC: $500 (with paid ads + sales)
   - LTV:CAC = 13:1 💰

✅ Gross Margin Blended: 85%
   - SaaS: 70% margin
   - Marketplace: 95% margin

✅ Team size: 15-20 personas
✅ Runway: 12+ meses
✅ Valuation: $5-7M (based on ARR multiple)
```

---

## Más Allá de Año 3: Visión 2028+

**Expansión Geográfica:**
- Brasil (Portuguese version)
- USA (contratar LATAM talent para USA companies - mercado masivo)
- Europa (timezone LATAM compatible con España)

**Vertical Expansion:**
- Non-tech roles (sales, marketing, customer success)
- Tech-adjacent (product managers, designers)

**Product Lines:**
- Dvra Talent (marketplace B2C - candidatos buscan directamente)
- Dvra API (sell access to network - other platforms pay)
- Dvra Analytics (standalone product - hiring intelligence)

**Exit Options:**
- Acquisition por ATS grande (Greenhouse, Lever) → $20-50M
- Acquisition por marketplace (Mercor, Turing) → $30-70M
- IPO path si llegamos a $10M+ ARR (largo plazo)

---

## Apéndice A: Schemas de Base de Datos

### A.1 Migrations ATS (Año 1)

```sql
-- Agregar campos para tier limits tracking
ALTER TABLE companies 
ADD COLUMN active_jobs_count INT DEFAULT 0,
ADD COLUMN candidates_this_month INT DEFAULT 0,
ADD COLUMN storage_used_mb INT DEFAULT 0;

-- Índices compuestos críticos
CREATE INDEX idx_candidates_company_email ON candidates(company_id, email);
CREATE INDEX idx_applications_company_stage ON applications(company_id, stage);
CREATE INDEX idx_jobs_company_status ON jobs(company_id, status);

-- Tabla de evaluaciones técnicas (GitHub integration)
CREATE TABLE technical_evaluations (
    id BIGSERIAL PRIMARY KEY,
    application_id BIGINT NOT NULL REFERENCES applications(id),
    company_id BIGINT NOT NULL REFERENCES companies(id),
    
    -- GitHub/GitLab data
    repo_url TEXT,
    commit_hash VARCHAR(40),
    pull_request_url TEXT,
    
    -- Evaluation
    code_quality_score INT CHECK (code_quality_score >= 0 AND code_quality_score <= 100),
    test_coverage_percent DECIMAL(5,2),
    lines_of_code INT,
    reviewer_user_id BIGINT REFERENCES users(id),
    review_notes TEXT,
    reviewed_at TIMESTAMP,
    
    status VARCHAR(50) DEFAULT 'pending', -- pending, in_review, approved, rejected
    
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP,
    
    INDEX idx_tech_eval_application (application_id),
    INDEX idx_tech_eval_company (company_id)
);

-- Tabla de templates de challenges
CREATE TABLE challenge_templates (
    id BIGSERIAL PRIMARY KEY,
    company_id BIGINT REFERENCES companies(id), -- NULL = public template
    
    name VARCHAR(255) NOT NULL,
    description TEXT,
    language VARCHAR(50), -- go, python, javascript, etc.
    difficulty_level VARCHAR(20), -- junior, mid, senior
    estimated_time_minutes INT,
    
    is_public BOOLEAN DEFAULT FALSE,
    usage_count INT DEFAULT 0,
    
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP,
    
    INDEX idx_challenge_company (company_id),
    INDEX idx_challenge_language (language),
    INDEX idx_challenge_difficulty (difficulty_level)
);
```

---

### A.2 Migrations Marketplace (Año 2 - Crítico)

```sql
-- ============================================
-- NETWORK CANDIDATES (Red Dvra)
-- ============================================

CREATE TABLE network_candidates (
    id BIGSERIAL PRIMARY KEY,
    
    -- Datos personales
    email VARCHAR(255) NOT NULL UNIQUE,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    phone VARCHAR(50),
    country VARCHAR(50),
    city VARCHAR(100),
    timezone VARCHAR(100), -- "America/Bogota"
    
    -- Perfil profesional
    title VARCHAR(255), -- "Senior Backend Developer"
    bio TEXT,
    years_exp INT,
    seniority_level VARCHAR(20), -- junior, mid, senior, staff
    
    -- Links
    linkedin_url TEXT,
    github_url TEXT,
    portfolio_url TEXT,
    resume_url TEXT, -- S3 path
    
    -- Evaluación Dvra (scoring automático + manual)
    dvra_score INT CHECK (dvra_score >= 0 AND dvra_score <= 100),
    code_quality_score INT CHECK (code_quality_score >= 0 AND code_quality_score <= 100),
    challenge_score INT CHECK (challenge_score >= 0 AND challenge_score <= 100),
    communication_score INT CHECK (communication_score >= 0 AND communication_score <= 100),
    evaluated_at TIMESTAMP,
    evaluated_by BIGINT REFERENCES users(id), -- Admin/Evaluator UserID
    
    -- Estado
    status VARCHAR(50) DEFAULT 'pending', -- pending, evaluating, approved, featured, inactive, rejected
    is_featured BOOLEAN DEFAULT FALSE,
    availability_status VARCHAR(50) DEFAULT 'available', -- available, interviewing, hired, not_looking
    hourly_rate INT, -- USD/hora (contractors)
    preferred_remote BOOLEAN DEFAULT TRUE,
    
    -- Metadata
    source_channel VARCHAR(100), -- landing_page, referral, bootcamp, linkedin
    referred_by BIGINT REFERENCES network_candidates(id), -- Referral tracking
    
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP,
    
    -- Indices críticos para búsqueda
    INDEX idx_network_candidate_email (email),
    INDEX idx_network_candidate_status (status),
    INDEX idx_network_candidate_featured (is_featured),
    INDEX idx_network_candidate_availability (availability_status),
    INDEX idx_network_candidate_seniority (seniority_level),
    INDEX idx_network_candidate_country (country),
    INDEX idx_network_candidate_score (dvra_score DESC)
);

-- ============================================
-- NETWORK CANDIDATE SKILLS
-- ============================================

CREATE TABLE network_candidate_skills (
    id BIGSERIAL PRIMARY KEY,
    candidate_id BIGINT NOT NULL REFERENCES network_candidates(id) ON DELETE CASCADE,
    
    skill_name VARCHAR(50) NOT NULL, -- "React", "Node.js", "AWS"
    category VARCHAR(50), -- frontend, backend, devops, mobile, data
    proficiency VARCHAR(20), -- beginner, intermediate, advanced, expert
    years_exp INT,
    
    -- Auto-detectado de GitHub
    is_verified BOOLEAN DEFAULT FALSE,
    
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP,
    
    INDEX idx_skill_candidate (candidate_id),
    INDEX idx_skill_name (skill_name),
    INDEX idx_skill_category (category),
    
    UNIQUE (candidate_id, skill_name) -- Un skill por candidato (update si existe)
);

-- ============================================
-- NETWORK CANDIDATE INTERESTS
-- ============================================

CREATE TABLE network_candidate_interests (
    id BIGSERIAL PRIMARY KEY,
    candidate_id BIGINT NOT NULL REFERENCES network_candidates(id) ON DELETE CASCADE,
    
    -- Preferencias de trabajo
    remote_only BOOLEAN DEFAULT FALSE,
    willing_to_relocate BOOLEAN DEFAULT FALSE,
    preferred_countries TEXT, -- JSON array: ["USA", "Canada", "Spain"]
    preferred_industries TEXT, -- JSON array: ["fintech", "healthtech"]
    
    -- Compensación
    min_salary_expected INT, -- USD anual
    max_salary_expected INT,
    
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP,
    
    INDEX idx_interest_candidate (candidate_id),
    
    UNIQUE (candidate_id) -- Solo un registro por candidato
);

-- ============================================
-- NETWORK APPLICATIONS (Company ↔ NetworkCandidate)
-- ============================================

CREATE TABLE network_applications (
    id BIGSERIAL PRIMARY KEY,
    company_id BIGINT NOT NULL REFERENCES companies(id),
    job_id BIGINT REFERENCES jobs(id), -- Nullable: búsqueda sin job específico
    candidate_id BIGINT NOT NULL REFERENCES network_candidates(id),
    
    -- Flujo del marketplace
    status VARCHAR(50) DEFAULT 'interested', -- interested, contacted, interviewing, offer, hired, rejected
    initiated_by VARCHAR(20), -- company, candidate
    
    -- Communication timestamps
    first_contact_at TIMESTAMP,
    interviewed_at TIMESTAMP,
    offer_extended_at TIMESTAMP,
    hired_at TIMESTAMP,
    rejected_at TIMESTAMP,
    rejection_reason TEXT,
    
    -- Fee tracking (CRÍTICO para revenue)
    fee_amount INT, -- USD: $3,500 flat (Año 1-2) o 15% salary (Año 3+)
    fee_status VARCHAR(50) DEFAULT 'pending', -- pending, invoiced, paid, refunded
    invoiced_at TIMESTAMP,
    paid_at TIMESTAMP,
    
    -- Notes y metadata
    notes TEXT,
    
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP,
    
    -- Indices para analytics
    INDEX idx_net_app_company (company_id),
    INDEX idx_net_app_candidate (candidate_id),
    INDEX idx_net_app_job (job_id),
    INDEX idx_net_app_status (status),
    INDEX idx_net_app_fee_status (fee_status),
    INDEX idx_net_app_hired (hired_at),
    
    UNIQUE (company_id, candidate_id, job_id) -- No duplicar aplicación al mismo job
);

-- ============================================
-- NETWORK EVALUATIONS (Proceso de evaluación)
-- ============================================

CREATE TABLE network_evaluations (
    id BIGSERIAL PRIMARY KEY,
    candidate_id BIGINT NOT NULL REFERENCES network_candidates(id) ON DELETE CASCADE,
    
    -- Tipo de evaluación
    evaluation_type VARCHAR(50) NOT NULL, -- github_analysis, coding_challenge, live_interview
    
    -- GitHub Analysis (automático)
    github_repos_analyzed INT,
    total_commits INT,
    languages_used TEXT, -- JSON: {"Go": 5000, "Python": 2000}
    top_repo_stars INT,
    code_quality_score INT CHECK (code_quality_score >= 0 AND code_quality_score <= 100),
    
    -- Coding Challenge
    challenge_id BIGINT REFERENCES challenge_templates(id),
    challenge_name VARCHAR(255),
    completed_at TIMESTAMP,
    time_spent_minutes INT,
    tests_passed INT,
    tests_total INT,
    code_review_notes TEXT,
    
    -- Live Interview
    interviewer_id BIGINT REFERENCES users(id),
    interview_duration INT, -- minutos
    technical_score INT CHECK (technical_score >= 0 AND technical_score <= 100),
    communication_score INT CHECK (communication_score >= 0 AND communication_score <= 100),
    interview_notes TEXT,
    
    -- Score final y recomendación
    final_score INT CHECK (final_score >= 0 AND final_score <= 100),
    recommendation VARCHAR(50), -- strong_yes, yes, maybe, no
    
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP,
    
    INDEX idx_eval_candidate (candidate_id),
    INDEX idx_eval_type (evaluation_type),
    INDEX idx_eval_score (final_score DESC)
);

-- ============================================
-- VIEWS ÚTILES
-- ============================================

-- Vista: Candidatos aprobados con sus skills
CREATE VIEW v_approved_candidates AS
SELECT 
    nc.id,
    nc.email,
    nc.first_name,
    nc.last_name,
    nc.title,
    nc.seniority_level,
    nc.country,
    nc.dvra_score,
    nc.is_featured,
    nc.availability_status,
    STRING_AGG(DISTINCT ncs.skill_name, ', ') AS skills,
    nc.created_at,
    nc.updated_at
FROM network_candidates nc
LEFT JOIN network_candidate_skills ncs ON nc.id = ncs.candidate_id
WHERE nc.status = 'approved' 
  AND nc.deleted_at IS NULL
  AND nc.availability_status IN ('available', 'interviewing')
GROUP BY nc.id;

-- Vista: Revenue por empresa del marketplace
CREATE VIEW v_marketplace_revenue_by_company AS
SELECT 
    c.id AS company_id,
    c.name AS company_name,
    COUNT(na.id) FILTER (WHERE na.status = 'hired') AS total_hires,
    SUM(na.fee_amount) FILTER (WHERE na.fee_status = 'paid') AS total_fees_paid,
    SUM(na.fee_amount) FILTER (WHERE na.fee_status = 'pending') AS pending_fees,
    MAX(na.hired_at) AS last_hire_date
FROM companies c
LEFT JOIN network_applications na ON c.id = na.company_id
GROUP BY c.id, c.name
HAVING COUNT(na.id) FILTER (WHERE na.status = 'hired') > 0
ORDER BY total_fees_paid DESC;

-- Vista: Funnel del marketplace
CREATE VIEW v_marketplace_funnel AS
SELECT 
    COUNT(*) FILTER (WHERE status = 'interested') AS interested,
    COUNT(*) FILTER (WHERE status = 'contacted') AS contacted,
    COUNT(*) FILTER (WHERE status = 'interviewing') AS interviewing,
    COUNT(*) FILTER (WHERE status = 'offer') AS offer_extended,
    COUNT(*) FILTER (WHERE status = 'hired') AS hired,
    COUNT(*) FILTER (WHERE status = 'rejected') AS rejected,
    ROUND(
        COUNT(*) FILTER (WHERE status = 'hired')::NUMERIC / 
        NULLIF(COUNT(*) FILTER (WHERE status IN ('interested', 'contacted', 'interviewing', 'offer', 'hired')), 0) * 100, 
        2
    ) AS conversion_rate_percent
FROM network_applications
WHERE deleted_at IS NULL;
```

---

### A.3 Funciones Útiles (PostgreSQL)

```sql
-- Función: Calcular Dvra Score automático
CREATE OR REPLACE FUNCTION calculate_dvra_score(
    p_candidate_id BIGINT
) RETURNS INT AS $$
DECLARE
    v_github_score INT;
    v_code_score INT;
    v_challenge_score INT;
    v_comm_score INT;
    v_final_score INT;
BEGIN
    SELECT 
        COALESCE(AVG(code_quality_score), 0) INTO v_github_score
    FROM network_evaluations
    WHERE candidate_id = p_candidate_id
      AND evaluation_type = 'github_analysis';
    
    SELECT 
        COALESCE(code_quality_score, 0),
        COALESCE(challenge_score, 0),
        COALESCE(communication_score, 0)
    INTO v_code_score, v_challenge_score, v_comm_score
    FROM network_candidates
    WHERE id = p_candidate_id;
    
    -- Pesos: GitHub 30%, Code 35%, Challenge 25%, Comm 10%
    v_final_score := ROUND(
        v_github_score * 0.30 + 
        v_code_score * 0.35 + 
        v_challenge_score * 0.25 + 
        v_comm_score * 0.10
    );
    
    RETURN v_final_score;
END;
$$ LANGUAGE plpgsql;

-- Trigger: Auto-update dvra_score cuando se completa evaluación
CREATE OR REPLACE FUNCTION update_candidate_score_on_evaluation()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE network_candidates
    SET dvra_score = calculate_dvra_score(NEW.candidate_id),
        updated_at = NOW()
    WHERE id = NEW.candidate_id;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_update_score_after_evaluation
AFTER INSERT OR UPDATE ON network_evaluations
FOR EACH ROW
EXECUTE FUNCTION update_candidate_score_on_evaluation();

-- Función: Calcular match score entre job y candidate
CREATE OR REPLACE FUNCTION calculate_match_score(
    p_job_id BIGINT,
    p_candidate_id BIGINT
) RETURNS DECIMAL AS $$
DECLARE
    v_skill_match DECIMAL;
    v_seniority_match DECIMAL;
    v_location_match DECIMAL;
    v_final_score DECIMAL;
BEGIN
    -- Skill overlap (simplified - en prod usar TF-IDF)
    SELECT 
        COUNT(DISTINCT ncs.skill_name)::DECIMAL / 
        NULLIF(
            (SELECT COUNT(*) FROM job_skills WHERE job_id = p_job_id),
            0
        ) * 100
    INTO v_skill_match
    FROM network_candidate_skills ncs
    WHERE ncs.candidate_id = p_candidate_id
      AND EXISTS (
          SELECT 1 FROM job_skills js
          WHERE js.job_id = p_job_id
            AND js.skill_name = ncs.skill_name
      );
    
    v_skill_match := COALESCE(v_skill_match, 0);
    
    -- Seniority match (exact match = 100%, adjacent = 70%, other = 30%)
    -- TODO: Implementar lógica completa
    v_seniority_match := 80; -- Placeholder
    
    -- Location compatibility (timezone overlap)
    v_location_match := 90; -- Placeholder
    
    -- Weighted score: Skills 60%, Seniority 30%, Location 10%
    v_final_score := 
        v_skill_match * 0.60 + 
        v_seniority_match * 0.30 + 
        v_location_match * 0.10;
    
    RETURN ROUND(v_final_score, 2);
END;
$$ LANGUAGE plpgsql;
```

---

### A.4 Queries Analytics Críticos

```sql
-- Dashboard Principal: Métricas clave
SELECT 
    (SELECT COUNT(*) FROM companies WHERE plan_tier != 'free' AND deleted_at IS NULL) AS active_paying_companies,
    (SELECT COUNT(*) FROM network_candidates WHERE status = 'approved' AND deleted_at IS NULL) AS approved_candidates,
    (SELECT COUNT(*) FROM network_candidates WHERE is_featured = TRUE AND deleted_at IS NULL) AS featured_candidates,
    (SELECT COUNT(*) FROM network_applications WHERE status = 'hired') AS total_marketplace_hires,
    (SELECT SUM(fee_amount) FROM network_applications WHERE fee_status = 'paid') AS total_fees_collected,
    (SELECT SUM(fee_amount) FROM network_applications WHERE fee_status = 'pending') AS pending_fees_to_collect;

-- Top Candidates por Match Score (para un job específico)
SELECT 
    nc.id,
    nc.first_name || ' ' || nc.last_name AS full_name,
    nc.title,
    nc.seniority_level,
    nc.country,
    nc.dvra_score,
    nc.is_featured,
    STRING_AGG(DISTINCT ncs.skill_name, ', ') AS skills,
    calculate_match_score(:job_id, nc.id) AS match_score
FROM network_candidates nc
LEFT JOIN network_candidate_skills ncs ON nc.id = ncs.candidate_id
WHERE nc.status = 'approved'
  AND nc.availability_status = 'available'
  AND nc.deleted_at IS NULL
GROUP BY nc.id, nc.first_name, nc.last_name, nc.title, nc.seniority_level, nc.country, nc.dvra_score, nc.is_featured
ORDER BY match_score DESC, nc.dvra_score DESC
LIMIT 20;

-- Conversion Rate por Company (Marketplace)
SELECT 
    c.name AS company_name,
    c.plan_tier,
    COUNT(na.id) AS total_interactions,
    COUNT(*) FILTER (WHERE na.status = 'interested') AS interested,
    COUNT(*) FILTER (WHERE na.status = 'contacted') AS contacted,
    COUNT(*) FILTER (WHERE na.status = 'hired') AS hired,
    ROUND(
        COUNT(*) FILTER (WHERE na.status = 'hired')::NUMERIC / 
        NULLIF(COUNT(na.id), 0) * 100,
        2
    ) AS conversion_rate_percent,
    SUM(na.fee_amount) FILTER (WHERE na.fee_status = 'paid') AS total_revenue
FROM companies c
INNER JOIN network_applications na ON c.id = na.company_id
WHERE c.plan_tier IN ('business', 'enterprise')
  AND c.deleted_at IS NULL
GROUP BY c.id, c.name, c.plan_tier
HAVING COUNT(*) FILTER (WHERE na.status = 'hired') > 0
ORDER BY total_revenue DESC;

-- Skills más demandados en el marketplace
SELECT 
    ncs.skill_name,
    ncs.category,
    COUNT(DISTINCT na.candidate_id) AS times_hired_with_this_skill,
    AVG(nc.dvra_score) AS avg_score_of_candidates_with_skill,
    COUNT(DISTINCT nc.id) AS total_candidates_with_skill
FROM network_candidate_skills ncs
INNER JOIN network_candidates nc ON ncs.candidate_id = nc.id
LEFT JOIN network_applications na ON nc.id = na.candidate_id AND na.status = 'hired'
WHERE nc.status = 'approved'
GROUP BY ncs.skill_name, ncs.category
ORDER BY times_hired_with_this_skill DESC, total_candidates_with_skill DESC
LIMIT 30;

-- Revenue Forecast (próximos 3 meses basado en pipeline)
SELECT 
    DATE_TRUNC('month', na.first_contact_at + INTERVAL '30 days') AS estimated_close_month,
    COUNT(*) FILTER (WHERE na.status IN ('interviewing', 'offer')) AS opportunities_in_pipeline,
    SUM(na.fee_amount) FILTER (WHERE na.status IN ('interviewing', 'offer')) AS potential_revenue,
    -- Asumir 30% close rate histórico
    SUM(na.fee_amount * 0.30) FILTER (WHERE na.status IN ('interviewing', 'offer')) AS forecasted_revenue
FROM network_applications na
WHERE na.status IN ('interviewing', 'offer')
  AND na.first_contact_at >= NOW() - INTERVAL '90 days'
  AND na.deleted_at IS NULL
GROUP BY estimated_close_month
ORDER BY estimated_close_month;
```

---

## Apéndice B: Checklist de Go-Live

### B.1 Pre-Launch (Día -7)

**Seguridad:**
- [ ] SSL/TLS certificates configurados (Let's Encrypt)
- [ ] Secrets en variables de entorno (no hardcoded)
- [ ] Rate limiting en endpoints públicos
- [ ] CORS configurado correctamente
- [ ] SQL injection prevention (prepared statements)
- [ ] XSS protection (sanitize inputs)

**Performance:**
- [ ] Database indexes creados (ver A.2)
- [ ] Connection pooling configurado
- [ ] CDN para assets estáticos (CloudFlare)
- [ ] Image optimization (S3 + CloudFront)
- [ ] Query optimization (EXPLAIN ANALYZE críticos)

**Monitoring:**
- [ ] Logs centralizados (CloudWatch o Papertrail)
- [ ] Error tracking (Sentry)
- [ ] Uptime monitoring (UptimeRobot)
- [ ] Performance monitoring (New Relic o DataDog)

**Backups:**
- [ ] Automated database backups (daily)
- [ ] S3 versioning enabled
- [ ] Disaster recovery plan documentado

### B.2 Launch Day (Día 0)

- [ ] Deploy a producción (AWS Lightsail/EC2)
- [ ] DNS apuntando correctamente
- [ ] Smoke tests en prod
- [ ] Launch blog post + social media
- [ ] Monitoring alerts configurados
- [ ] On-call rotation definido

### B.3 Post-Launch (Día +1 a +7)

- [ ] Daily health checks
- [ ] User feedback monitoring
- [ ] Performance metrics review
- [ ] Bug triage meetings
- [ ] Hotfix deployment plan listo

---

## Apéndice C: Contacto y Recursos

**Equipo Fundador:**
- CEO/CTO: [Tu nombre]
- Email: founders@dvra.app
- LinkedIn: [Tu perfil]

**Documentación Técnica:**
- API Docs: docs.dvra.app
- GitHub: github.com/dvra-team/dvra-api
- Postman Collection: [Link]

**Inversión & Partnerships:**
- Pitch Deck: deck.dvra.app
- Email: investors@dvra.app

**Support:**
- Help Center: help.dvra.app
- Status Page: status.dvra.app
- Community Slack: community.dvra.app

---

**FIN DEL DOCUMENTO**

> **Versión 2.0** - Diciembre 8, 2025  
> **Próxima Revisión**: Q1 2026 (post-MVP launch)
    id BIGSERIAL PRIMARY KEY,
    company_id BIGINT REFERENCES companies(id),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    language VARCHAR(50),
    difficulty_level VARCHAR(20),
    estimated_time_minutes INT,
    is_public BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Tabla de activity logs (audit)
CREATE TABLE activity_logs (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id),
    company_id BIGINT REFERENCES companies(id),
    action VARCHAR(100) NOT NULL, -- 'job.created', 'candidate.moved', etc
    entity_type VARCHAR(50), -- 'job', 'candidate', 'application'
    entity_id BIGINT,
    metadata JSONB,
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);
```

---

## Apéndice B: Decisiones Técnicas Clave

### B.1 ¿Por qué Go?

1. **Performance**: Maneja 10k req/s con recursos mínimos
2. **Concurrencia**: Goroutines ideales para background jobs
3. **Deploy simple**: Binary único, sin runtime dependencies
4. **Ecosistema**: Gin, GORM, excelente para APIs REST
5. **Costo AWS**: Menos memoria = EC2 t3.micro suficiente

### B.2 ¿Por qué PostgreSQL?

1. **JSONB**: Metadata flexible (configs por empresa)
2. **Full-text search**: Búsqueda de candidatos por skills
3. **Mature**: 25+ años, probado en producción
4. **RDS managed**: Backups, replication automático en AWS

### B.3 ¿Por qué Multi-Tenant con CompanyID?

**Alternativas consideradas**:
1. ❌ Database per tenant: No escala (manage 1000 DBs)
2. ❌ Schema per tenant: Mismo problema
3. ✅ Shared database + CompanyID: Escala a millones de filas

**Riesgos mitigados**:
- Scoped queries automáticos (middleware)
- Indices optimizados (company_id + other_field)
- Test exhaustivos de data leakage

---

## Apéndice C: Glossario

- **ATS**: Applicant Tracking System
- **Tenant**: Empresa cliente en arquitectura multi-tenant
- **Stage**: Etapa del pipeline de candidato
- **Source**: Canal de donde viene el candidato
- **Rating**: Calificación 1-5 estrellas
- **Scoped Query**: Query filtrado por CompanyID
- **JWT**: JSON Web Token (autenticación)
- **RBAC**: Role-Based Access Control
- **MRR**: Monthly Recurring Revenue
- **ARR**: Annual Recurring Revenue
- **CAC**: Customer Acquisition Cost
- **LTV**: Lifetime Value

---

## Changelog

- **v1.0** (2025-12-08): Documento inicial basado en análisis TechRecruit

---

**Documento vivo**: Este documento se actualiza conforme el producto evoluciona.  
**Última revisión**: Diciembre 8, 2025  
**Próxima revisión**: Enero 15, 2026
