# Dvra - Arquitectura Técnica

> **Documentación de Arquitectura y Decisiones Técnicas**  
> Versión: 1.0 | Última actualización: Diciembre 8, 2025

---

## 📋 Tabla de Contenidos

1. [Tech Stack](#tech-stack)
2. [Arquitectura de Sistema](#arquitectura-de-sistema)
3. [Modelos de Datos](#modelos-de-datos)
4. [Schemas de Base de Datos](#schemas-de-base-de-datos)
5. [Servicios y APIs](#servicios-y-apis)
6. [Implementación de Seguridad](#implementación-de-seguridad)
7. [Performance y Escalabilidad](#performance-y-escalabilidad)
8. [Infraestructura](#infraestructura)

---

## 1. Tech Stack

### Backend
- **Lenguaje**: Go 1.24.0
- **Framework Web**: Gin (HTTP router)
- **ORM**: GORM (PostgreSQL)
- **Autenticación**: JWT (golang-jwt/jwt)
- **Validación**: go-playground/validator

### Base de Datos
- **Primary**: PostgreSQL 16
- **Features usadas**:
  - JSONB columns
  - Full-text search
  - Partial indexes
  - Foreign keys con CASCADE

### Infraestructura
- **Hosting**: AWS Lightsail / EC2 (start), ECS (scale)
- **Database**: AWS RDS PostgreSQL (Multi-AZ en prod)
- **Storage**: AWS S3 (CVs, documentos)
- **CDN**: CloudFront (assets estáticos)
- **Email**: SendGrid / AWS SES
- **Monitoring**: AWS CloudWatch + Sentry

### DevOps
- **Containerization**: Docker + Docker Compose
- **CI/CD**: GitHub Actions
- **Secrets**: AWS Secrets Manager
- **Backups**: Automated daily snapshots

---

## 2. Arquitectura de Sistema

### 2.1 Arquitectura de Capas

```
┌─────────────────────────────────────────┐
│           Frontend (Futuro)             │
│         React/Next.js + TypeScript      │
└─────────────────┬───────────────────────┘
                  │ HTTP/REST
┌─────────────────▼───────────────────────┐
│              API Gateway                │
│         Gin Router + Middleware         │
│   (Auth, CORS, Rate Limiting, Logging)  │
└─────────────────┬───────────────────────┘
                  │
┌─────────────────▼───────────────────────┐
│            Handler Layer                │
│   (HTTP Request/Response handling)      │
│   - user_handler.go                     │
│   - job_handler.go                      │
│   - application_handler.go              │
└─────────────────┬───────────────────────┘
                  │
┌─────────────────▼───────────────────────┐
│            Service Layer                │
│       (Business Logic Execution)        │
│   - user_service.go                     │
│   - job_service.go                      │
│   - jwt_service.go                      │
└─────────────────┬───────────────────────┘
                  │
┌─────────────────▼───────────────────────┐
│          Repository Layer               │
│      (Data Access Abstraction)          │
│   - user_repository.go                  │
│   - job_repository.go                   │
└─────────────────┬───────────────────────┘
                  │
┌─────────────────▼───────────────────────┐
│          Database Layer                 │
│      PostgreSQL 16 + GORM ORM           │
└─────────────────────────────────────────┘
```

### 2.2 Multi-Tenancy

**Tipo**: Row-Level Security con CompanyID scoping

```go
// Cada entidad tiene company_id
type Job struct {
    gorm.Model
    CompanyID uint `gorm:"not null;index"`
    // ...
}

// Middleware inyecta company_id en contexto
func TenantScope() gin.HandlerFunc {
    return func(c *gin.Context) {
        companyID := c.GetUint("company_id")
        
        // Scoped DB para queries automáticas
        scopedDB := db.Where("company_id = ?", companyID)
        c.Set("db", scopedDB)
        
        c.Next()
    }
}
```

**Ventajas**:
- ✅ Simple de implementar
- ✅ Cost-effective (single DB)
- ✅ Fácil de escalar inicialmente

**Consideraciones futuras** (si llegamos a 1000+ empresas):
- Sharding por company_id
- Separar empresas Enterprise en DB dedicadas

---

## 3. Modelos de Datos

### 3.1 Entidades Core (ATS)

```go
// Company (Tenant)
type Company struct {
    gorm.Model
    Name            string `gorm:"type:varchar(255);not null"`
    Slug            string `gorm:"type:varchar(100);uniqueIndex;not null"`
    LogoURL         string `gorm:"type:text"`
    Website         string `gorm:"type:text"`
    Industry        string `gorm:"type:varchar(100)"`
    CompanySize     string // "1-10", "11-50", "51-200", "201-500", "500+"
    
    // Subscription
    PlanTier        string `gorm:"type:varchar(50);default:'free'"`
    TrialEndsAt     *time.Time
    SubscriptionID  string // Stripe/MercadoPago subscription ID
    
    // Config
    Timezone        string `gorm:"type:varchar(100);default:'America/Bogota'"`
    
    // Límites (tracking)
    ActiveJobsCount     int `gorm:"default:0"`
    CandidatesThisMonth int `gorm:"default:0"`
    StorageUsedMB       int `gorm:"default:0"`
    
    // Relaciones
    Users        []User        `gorm:"many2many:memberships"`
    Jobs         []Job         `gorm:"foreignKey:CompanyID"`
    Candidates   []Candidate   `gorm:"foreignKey:CompanyID"`
}

// User
type User struct {
    gorm.Model
    Email        string `gorm:"type:varchar(255);uniqueIndex;not null"`
    PasswordHash string `gorm:"type:varchar(255);not null"`
    FirstName    string `gorm:"type:varchar(100)"`
    LastName     string `gorm:"type:varchar(100)"`
    AvatarURL    string `gorm:"type:text"`
    
    // Auth
    IsActive     bool      `gorm:"default:true"`
    LastLoginAt  *time.Time
    
    // Relaciones
    Companies    []Company    `gorm:"many2many:memberships"`
    Memberships  []Membership `gorm:"foreignKey:UserID"`
}

// Membership (User ↔ Company)
type Membership struct {
    gorm.Model
    UserID    uint   `gorm:"not null;index"`
    CompanyID *uint  `gorm:"index"` // NULL = SuperAdmin
    Role      string `gorm:"type:varchar(50);not null;default:'user'"`
    Status    string `gorm:"type:varchar(50);not null;default:'active'"`
    IsDefault bool   `gorm:"default:false"`
    
    // Relaciones
    User      User     `gorm:"foreignKey:UserID"`
    Company   *Company `gorm:"foreignKey:CompanyID"`
}

// Job
type Job struct {
    gorm.Model
    CompanyID         uint   `gorm:"not null;index"`
    Title             string `gorm:"type:varchar(255);not null"`
    Description       string `gorm:"type:text"`
    Location          string `gorm:"type:varchar(255)"`
    RemotePolicy      string // "remote", "hybrid", "onsite"
    EmploymentType    string // "full_time", "part_time", "contract"
    SalaryMin         *int
    SalaryMax         *int
    SalaryCurrency    string `gorm:"type:varchar(10);default:'USD'"`
    
    // Ownership
    AssignedRecruiterID *uint `gorm:"index"`
    HiringManagerID     *uint `gorm:"index"`
    
    // Status
    Status              string `gorm:"type:varchar(50);default:'draft'"`
    PublishedAt         *time.Time
    ClosedAt            *time.Time
    
    // Relaciones
    Company             Company       `gorm:"foreignKey:CompanyID"`
    AssignedRecruiter   *User         `gorm:"foreignKey:AssignedRecruiterID"`
    HiringManager       *User         `gorm:"foreignKey:HiringManagerID"`
    Applications        []Application `gorm:"foreignKey:JobID"`
}

// Candidate
type Candidate struct {
    gorm.Model
    CompanyID      uint   `gorm:"not null;index:idx_candidates_company_email"`
    Email          string `gorm:"type:varchar(255);not null;index:idx_candidates_company_email"`
    FirstName      string `gorm:"type:varchar(100);not null"`
    LastName       string `gorm:"type:varchar(100);not null"`
    Phone          string `gorm:"type:varchar(50)"`
    Location       string `gorm:"type:varchar(255)"`
    
    // Resume & Links
    ResumeURL      string `gorm:"type:text"`
    LinkedinURL    string `gorm:"type:text"`
    GithubURL      string `gorm:"type:text"`
    PortfolioURL   string `gorm:"type:text"`
    
    // Source tracking
    Source         string `gorm:"type:varchar(100)"` // linkedin, referral, etc.
    SourceDetails  string `gorm:"type:text"`
    
    // Relaciones
    Company        Company       `gorm:"foreignKey:CompanyID"`
    Applications   []Application `gorm:"foreignKey:CandidateID"`
}

// Application (candidato aplica a job)
type Application struct {
    gorm.Model
    CompanyID      uint   `gorm:"not null;index"`
    JobID          uint   `gorm:"not null;index"`
    CandidateID    uint   `gorm:"not null;index"`
    
    // Pipeline
    Stage          string `gorm:"type:varchar(50);not null;default:'applied'"`
    Rating         *int   // 1-5 stars
    Notes          string `gorm:"type:text"`
    
    // Timestamps
    AppliedAt      time.Time  `gorm:"not null"`
    RejectedAt     *time.Time
    HiredAt        *time.Time
    
    // Relaciones
    Company        Company   `gorm:"foreignKey:CompanyID"`
    Job            Job       `gorm:"foreignKey:JobID"`
    Candidate      Candidate `gorm:"foreignKey:CandidateID"`
}

// Role
type Role struct {
    gorm.Model
    Name        string `gorm:"type:varchar(50);uniqueIndex;not null"`
    Level       int    `gorm:"not null"` // Para jerarquía
    Description string `gorm:"type:text"`
}
```

### 3.2 Entidades Marketplace (Año 2+)

```go
// NetworkCandidate (Red Dvra)
type NetworkCandidate struct {
    gorm.Model
    
    // Identificación única (deduplicación)
    Email           string `gorm:"type:varchar(255);uniqueIndex;not null"`
    GithubUsername  string `gorm:"type:varchar(100);uniqueIndex"`
    
    // Perfil
    FirstName       string `gorm:"type:varchar(100);not null"`
    LastName        string `gorm:"type:varchar(100);not null"`
    Phone           string `gorm:"type:varchar(50)"`
    Country         string `gorm:"type:varchar(50)"`
    City            string `gorm:"type:varchar(100)"`
    Timezone        string `gorm:"type:varchar(100)"`
    
    Title           string `gorm:"type:varchar(255)"`
    Bio             string `gorm:"type:text"`
    YearsExp        int
    SeniorityLevel  string `gorm:"type:varchar(20)"`
    
    // Links
    LinkedinURL     string `gorm:"type:text"`
    GithubURL       string `gorm:"type:text"`
    PortfolioURL    string `gorm:"type:text"`
    ResumeURL       string `gorm:"type:text"`
    
    // Evaluación
    DvraScore           int `gorm:"check:dvra_score >= 0 AND dvra_score <= 100"`
    CodeQualityScore    int `gorm:"check:code_quality_score >= 0 AND code_quality_score <= 100"`
    ChallengeScore      int `gorm:"check:challenge_score >= 0 AND challenge_score <= 100"`
    CommunicationScore  int `gorm:"check:communication_score >= 0 AND communication_score <= 100"`
    EvaluatedAt         *time.Time
    EvaluatedBy         *uint
    
    // Estado
    Status              string `gorm:"type:varchar(50);not null;default:'prospect';index"`
    IsFeatured          bool   `gorm:"default:false;index"`
    AvailabilityStatus  string `gorm:"type:varchar(50);default:'available'"`
    HourlyRate          *int
    PreferredRemote     bool   `gorm:"default:true"`
    
    // Tracking
    ProspectedAt    *time.Time
    InvitedAt       *time.Time
    RegisteredAt    *time.Time
    LastContactAt   *time.Time
    
    // Consentimiento (GDPR/LGPD)
    OptedIn         bool   `gorm:"default:false;not null"`
    OptedInAt       *time.Time
    OptOutReason    string `gorm:"type:text"`
    ConsentIP       string `gorm:"type:varchar(50)"`
    ConsentSource   string `gorm:"type:varchar(100)"`
    
    // Source
    SourceChannel   string `gorm:"type:varchar(100)"`
    ReferredBy      *uint  `gorm:"index"`
    
    // Anti-spam
    InvitationsSent int `gorm:"default:0"`
    
    // Relaciones
    Skills          []NetworkCandidateSkill    `gorm:"foreignKey:CandidateID"`
    Interests       []NetworkCandidateInterest `gorm:"foreignKey:CandidateID"`
    Applications    []NetworkApplication       `gorm:"foreignKey:CandidateID"`
}

// NetworkCandidateSkill
type NetworkCandidateSkill struct {
    gorm.Model
    CandidateID uint   `gorm:"not null;index"`
    SkillName   string `gorm:"type:varchar(50);not null"`
    Category    string `gorm:"type:varchar(50)"`
    Proficiency string `gorm:"type:varchar(20)"`
    YearsExp    int
    IsVerified  bool   `gorm:"default:false"`
}

// NetworkApplication (Company contacta NetworkCandidate)
type NetworkApplication struct {
    gorm.Model
    CompanyID   uint  `gorm:"not null;index"`
    JobID       *uint `gorm:"index"`
    CandidateID uint  `gorm:"not null;index"`
    
    // Flujo
    Status          string `gorm:"type:varchar(50);default:'interested'"`
    InitiatedBy     string `gorm:"type:varchar(20)"`
    
    // Timestamps
    FirstContactAt  *time.Time
    InterviewedAt   *time.Time
    OfferExtendedAt *time.Time
    HiredAt         *time.Time
    RejectedAt      *time.Time
    RejectionReason string `gorm:"type:text"`
    
    // Fee tracking
    FeeAmount   *int   `gorm:"type:integer"`
    FeeStatus   string `gorm:"type:varchar(50);default:'pending'"`
    InvoicedAt  *time.Time
    PaidAt      *time.Time
    
    // Relaciones
    Candidate   NetworkCandidate `gorm:"foreignKey:CandidateID"`
    Company     Company          `gorm:"foreignKey:CompanyID"`
    Job         *Job             `gorm:"foreignKey:JobID"`
}
```

---

## 4. Schemas de Base de Datos

### 4.1 Migrations ATS (Año 1)

```sql
-- Agregar tracking de límites
ALTER TABLE companies 
ADD COLUMN active_jobs_count INT DEFAULT 0,
ADD COLUMN candidates_this_month INT DEFAULT 0,
ADD COLUMN storage_used_mb INT DEFAULT 0;

-- Índices críticos para performance
CREATE INDEX idx_candidates_company_email 
ON candidates(company_id, email);

CREATE INDEX idx_applications_company_stage 
ON applications(company_id, stage);

CREATE INDEX idx_jobs_company_status 
ON jobs(company_id, status);

-- Evaluaciones técnicas (GitHub integration)
CREATE TABLE technical_evaluations (
    id BIGSERIAL PRIMARY KEY,
    application_id BIGINT NOT NULL REFERENCES applications(id) ON DELETE CASCADE,
    company_id BIGINT NOT NULL REFERENCES companies(id),
    
    repo_url TEXT,
    commit_hash VARCHAR(40),
    pull_request_url TEXT,
    
    code_quality_score INT CHECK (code_quality_score >= 0 AND code_quality_score <= 100),
    test_coverage_percent DECIMAL(5,2),
    lines_of_code INT,
    reviewer_user_id BIGINT REFERENCES users(id),
    review_notes TEXT,
    reviewed_at TIMESTAMP,
    
    status VARCHAR(50) DEFAULT 'pending',
    
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE INDEX idx_tech_eval_application ON technical_evaluations(application_id);
CREATE INDEX idx_tech_eval_company ON technical_evaluations(company_id);

-- Templates de coding challenges
CREATE TABLE challenge_templates (
    id BIGSERIAL PRIMARY KEY,
    company_id BIGINT REFERENCES companies(id),
    
    name VARCHAR(255) NOT NULL,
    description TEXT,
    language VARCHAR(50),
    difficulty_level VARCHAR(20),
    estimated_time_minutes INT,
    
    is_public BOOLEAN DEFAULT FALSE,
    usage_count INT DEFAULT 0,
    
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE INDEX idx_challenge_company ON challenge_templates(company_id);
CREATE INDEX idx_challenge_language ON challenge_templates(language);
```

### 4.2 Migrations Marketplace (Año 2)

```sql
-- Network Candidates (Red Dvra)
CREATE TABLE network_candidates (
    id BIGSERIAL PRIMARY KEY,
    
    -- Deduplicación
    email VARCHAR(255) NOT NULL UNIQUE,
    github_username VARCHAR(100) UNIQUE,
    
    -- Perfil
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    phone VARCHAR(50),
    country VARCHAR(50),
    city VARCHAR(100),
    timezone VARCHAR(100),
    
    title VARCHAR(255),
    bio TEXT,
    years_exp INT,
    seniority_level VARCHAR(20),
    
    -- Links
    linkedin_url TEXT,
    github_url TEXT,
    portfolio_url TEXT,
    resume_url TEXT,
    
    -- Evaluación
    dvra_score INT CHECK (dvra_score >= 0 AND dvra_score <= 100),
    code_quality_score INT CHECK (code_quality_score >= 0 AND code_quality_score <= 100),
    challenge_score INT CHECK (challenge_score >= 0 AND challenge_score <= 100),
    communication_score INT CHECK (communication_score >= 0 AND communication_score <= 100),
    evaluated_at TIMESTAMP,
    evaluated_by BIGINT REFERENCES users(id),
    
    -- Estado
    status VARCHAR(50) NOT NULL DEFAULT 'prospect',
    is_featured BOOLEAN DEFAULT FALSE,
    availability_status VARCHAR(50) DEFAULT 'available',
    hourly_rate INT,
    preferred_remote BOOLEAN DEFAULT TRUE,
    
    -- Tracking
    prospected_at TIMESTAMP,
    invited_at TIMESTAMP,
    registered_at TIMESTAMP,
    last_contact_at TIMESTAMP,
    
    -- Consentimiento (GDPR/LGPD)
    opted_in BOOLEAN NOT NULL DEFAULT FALSE,
    opted_in_at TIMESTAMP,
    opt_out_reason TEXT,
    consent_ip VARCHAR(50),
    consent_source VARCHAR(100),
    
    -- Source
    source_channel VARCHAR(100),
    referred_by BIGINT REFERENCES network_candidates(id),
    
    -- Anti-spam
    invitations_sent INT DEFAULT 0,
    
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP
);

-- Índices críticos
CREATE INDEX idx_network_candidate_email ON network_candidates(email);
CREATE INDEX idx_network_candidate_github ON network_candidates(github_username);
CREATE INDEX idx_network_candidate_status ON network_candidates(status);
CREATE INDEX idx_network_candidate_featured ON network_candidates(is_featured);
CREATE INDEX idx_network_candidate_score ON network_candidates(dvra_score DESC);

-- Índice compuesto para búsqueda de empresas
CREATE INDEX idx_network_candidates_search 
ON network_candidates(status, opted_in, deleted_at)
WHERE status IN ('active', 'featured') AND opted_in = TRUE AND deleted_at IS NULL;

-- Skills
CREATE TABLE network_candidate_skills (
    id BIGSERIAL PRIMARY KEY,
    candidate_id BIGINT NOT NULL REFERENCES network_candidates(id) ON DELETE CASCADE,
    
    skill_name VARCHAR(50) NOT NULL,
    category VARCHAR(50),
    proficiency VARCHAR(20),
    years_exp INT,
    is_verified BOOLEAN DEFAULT FALSE,
    
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP,
    
    UNIQUE(candidate_id, skill_name)
);

CREATE INDEX idx_skill_candidate ON network_candidate_skills(candidate_id);
CREATE INDEX idx_skill_name ON network_candidate_skills(skill_name);
CREATE INDEX idx_skill_category ON network_candidate_skills(category);

-- Network Applications
CREATE TABLE network_applications (
    id BIGSERIAL PRIMARY KEY,
    company_id BIGINT NOT NULL REFERENCES companies(id),
    job_id BIGINT REFERENCES jobs(id),
    candidate_id BIGINT NOT NULL REFERENCES network_candidates(id),
    
    status VARCHAR(50) DEFAULT 'interested',
    initiated_by VARCHAR(20),
    
    first_contact_at TIMESTAMP,
    interviewed_at TIMESTAMP,
    offer_extended_at TIMESTAMP,
    hired_at TIMESTAMP,
    rejected_at TIMESTAMP,
    rejection_reason TEXT,
    
    fee_amount INT,
    fee_status VARCHAR(50) DEFAULT 'pending',
    invoiced_at TIMESTAMP,
    paid_at TIMESTAMP,
    
    notes TEXT,
    
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP,
    
    UNIQUE(company_id, candidate_id, job_id)
);

CREATE INDEX idx_net_app_company ON network_applications(company_id);
CREATE INDEX idx_net_app_candidate ON network_applications(candidate_id);
CREATE INDEX idx_net_app_job ON network_applications(job_id);
CREATE INDEX idx_net_app_status ON network_applications(status);
CREATE INDEX idx_net_app_fee_status ON network_applications(fee_status);
CREATE INDEX idx_net_app_hired ON network_applications(hired_at);
```

### 4.3 Views Útiles

```sql
-- Vista: Candidatos aprobados con skills
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

-- Vista: Revenue del marketplace por empresa
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

## 5. Servicios y APIs

### 5.1 Estructura de Servicios

```
internal/
├── app/
│   ├── handlers/          # HTTP handlers (controllers)
│   │   ├── user_handler.go
│   │   ├── job_handler.go
│   │   ├── application_handler.go
│   │   └── network_handler.go
│   │
│   ├── services/          # Business logic
│   │   ├── user_service.go
│   │   ├── job_service.go
│   │   ├── application_service.go
│   │   ├── jwt_service.go
│   │   └── network_candidate_service.go
│   │
│   ├── repositories/      # Data access
│   │   ├── user_repository.go
│   │   ├── job_repository.go
│   │   └── application_repository.go
│   │
│   └── models/           # Domain models
│       ├── user.go
│       ├── job.go
│       └── application.go
│
├── shared/
│   ├── middleware/       # HTTP middleware
│   │   ├── auth.go
│   │   ├── cors.go
│   │   └── tenant.go
│   │
│   └── utils/           # Helpers
│       ├── pagination.go
│       └── errors.go
│
└── platform/
    ├── config/          # Configuration
    │   └── config.go
    │
    └── database/        # DB setup
        └── database.go
```

### 5.2 Ejemplo de Service

```go
package services

import (
    "errors"
    "time"
    "gorm.io/gorm"
)

type NetworkCandidateService struct {
    db            *gorm.DB
    emailService  EmailService
    s3Service     S3Service
}

func NewNetworkCandidateService(db *gorm.DB) *NetworkCandidateService {
    return &NetworkCandidateService{
        db: db,
    }
}

// CreateOrUpdate - Deduplicación inteligente
func (s *NetworkCandidateService) CreateOrUpdate(dto CreateCandidateDTO) (*NetworkCandidate, error) {
    var candidate NetworkCandidate
    
    result := s.db.Where("email = ? OR github_username = ?", 
        dto.Email, dto.GithubUsername).First(&candidate)
    
    if result.Error == gorm.ErrRecordNotFound {
        // Nuevo candidato
        candidate = NetworkCandidate{
            Email:          dto.Email,
            GithubUsername: dto.GithubUsername,
            FirstName:      dto.FirstName,
            LastName:       dto.LastName,
            Status:         StatusProspect,
            ProspectedAt:   timePtr(time.Now()),
            SourceChannel:  dto.Source,
            OptedIn:        false,
        }
        
        if err := s.db.Create(&candidate).Error; err != nil {
            return nil, err
        }
        
        return &candidate, nil
    }
    
    // Ya existe - enriquecer si está en estado temprano
    if candidate.Status == StatusProspect || candidate.Status == StatusInvited {
        if candidate.FirstName == "" && dto.FirstName != "" {
            candidate.FirstName = dto.FirstName
        }
        if candidate.LastName == "" && dto.LastName != "" {
            candidate.LastName = dto.LastName
        }
        
        s.db.Save(&candidate)
    }
    
    return &candidate, nil
}

// TransitionStatus - State machine con validaciones
func (s *NetworkCandidateService) TransitionStatus(
    candidateID uint, 
    newStatus TalentStatus,
) error {
    candidate, err := s.GetByID(candidateID)
    if err != nil {
        return err
    }
    
    if !isValidTransition(candidate.Status, newStatus) {
        return errors.New("invalid status transition")
    }
    
    candidate.Status = newStatus
    
    switch newStatus {
    case StatusInvited:
        candidate.InvitedAt = timePtr(time.Now())
        candidate.InvitationsSent++
    case StatusActive:
        candidate.RegisteredAt = timePtr(time.Now())
        candidate.OptedIn = true
        candidate.OptedInAt = timePtr(time.Now())
    }
    
    return s.db.Save(candidate).Error
}
```

---

## 6. Implementación de Seguridad

### 6.1 JWT Authentication

```go
package services

import (
    "time"
    "github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
    secretKey string
}

type Claims struct {
    UserID    uint   `json:"user_id"`
    CompanyID uint   `json:"company_id"`
    Role      string `json:"role"`
    jwt.RegisteredClaims
}

func (s *JWTService) GenerateToken(user *User, membership *Membership) (string, error) {
    claims := Claims{
        UserID:    user.ID,
        CompanyID: *membership.CompanyID,
        Role:      membership.Role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(s.secretKey))
}

func (s *JWTService) ValidateToken(tokenString string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        return []byte(s.secretKey), nil
    })
    
    if err != nil {
        return nil, err
    }
    
    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims, nil
    }
    
    return nil, errors.New("invalid token")
}
```

### 6.2 Middleware de Autenticación

```go
package middleware

import (
    "github.com/gin-gonic/gin"
)

func RequireAuth(jwtService *JWTService) gin.HandlerFunc {
    return func(c *gin.Context) {
        tokenString := c.GetHeader("Authorization")
        if tokenString == "" {
            c.JSON(401, gin.H{"error": "authorization header required"})
            c.Abort()
            return
        }
        
        // Remove "Bearer " prefix
        if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
            tokenString = tokenString[7:]
        }
        
        claims, err := jwtService.ValidateToken(tokenString)
        if err != nil {
            c.JSON(401, gin.H{"error": "invalid token"})
            c.Abort()
            return
        }
        
        // Inyectar claims en contexto
        c.Set("user_id", claims.UserID)
        c.Set("company_id", claims.CompanyID)
        c.Set("role", claims.Role)
        
        c.Next()
    }
}

func RequireRole(minLevel int) gin.HandlerFunc {
    return func(c *gin.Context) {
        role := c.GetString("role")
        userLevel := getRoleLevel(role)
        
        if userLevel < minLevel {
            c.JSON(403, gin.H{"error": "insufficient permissions"})
            c.Abort()
            return
        }
        
        c.Next()
    }
}

func getRoleLevel(role string) int {
    levels := map[string]int{
        "superadmin":     100,
        "admin":          50,
        "recruiter":      30,
        "hiring_manager": 20,
        "user":           10,
    }
    return levels[role]
}
```

---

## 7. Performance y Escalabilidad

### 7.1 Database Connection Pool

```go
// config/database.go
func SetupDatabase(config *Config) (*gorm.DB, error) {
    dsn := fmt.Sprintf(
        "host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
        config.DBHost, config.DBUser, config.DBPassword, config.DBName, config.DBPort,
    )
    
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
    })
    
    if err != nil {
        return nil, err
    }
    
    sqlDB, err := db.DB()
    if err != nil {
        return nil, err
    }
    
    // Connection pool settings
    sqlDB.SetMaxIdleConns(10)
    sqlDB.SetMaxOpenConns(100)
    sqlDB.SetConnMaxLifetime(time.Hour)
    
    return db, nil
}
```

### 7.2 Caching Strategy (Futuro)

```go
// Redis para session caching
type CacheService struct {
    redis *redis.Client
}

func (s *CacheService) GetUser(userID uint) (*User, error) {
    key := fmt.Sprintf("user:%d", userID)
    
    // Try cache first
    val, err := s.redis.Get(ctx, key).Result()
    if err == nil {
        var user User
        json.Unmarshal([]byte(val), &user)
        return &user, nil
    }
    
    // Cache miss - fetch from DB
    user := fetchUserFromDB(userID)
    
    // Cache for 1 hour
    data, _ := json.Marshal(user)
    s.redis.Set(ctx, key, data, time.Hour)
    
    return user, nil
}
```

### 7.3 Query Optimization

```go
// Usar preloading para evitar N+1
func (s *JobService) GetJobsWithApplications(companyID uint) ([]Job, error) {
    var jobs []Job
    
    err := s.db.
        Preload("Applications").
        Preload("Applications.Candidate").
        Where("company_id = ?", companyID).
        Find(&jobs).Error
    
    return jobs, err
}

// Usar select para limitar campos
func (s *CandidateService) ListCandidates(companyID uint) ([]Candidate, error) {
    var candidates []Candidate
    
    err := s.db.
        Select("id", "email", "first_name", "last_name", "created_at").
        Where("company_id = ?", companyID).
        Limit(100).
        Find(&candidates).Error
    
    return candidates, err
}
```

---

## 8. Infraestructura

### 8.1 Docker Setup

```dockerfile
# Dockerfile
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o main cmd/dvra-api/main.go

# Runtime image
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]
```

```yaml
# docker-compose.yml
version: '3.8'

services:
  api:
    build: .
    ports:
      - "8080:8080"
    environment:
      - DB_HOST=postgres
      - DB_USER=dvra
      - DB_PASSWORD=secret
      - DB_NAME=dvra_db
      - JWT_SECRET=your-secret-key
    depends_on:
      - postgres

  postgres:
    image: postgres:16-alpine
    ports:
      - "5432:5432"
    environment:
      - POSTGRES_USER=dvra
      - POSTGRES_PASSWORD=secret
      - POSTGRES_DB=dvra_db
    volumes:
      - postgres_data:/var/lib/postgresql/data

volumes:
  postgres_data:
```

### 8.2 AWS Infrastructure (Producción)

```
┌─────────────────────────────────────────┐
│         Route 53 (DNS)                  │
│         api.dvra.app                    │
└────────────────┬────────────────────────┘
                 │
┌────────────────▼────────────────────────┐
│    CloudFront (CDN + SSL)               │
└────────────────┬────────────────────────┘
                 │
┌────────────────▼────────────────────────┐
│  Application Load Balancer              │
│  (Multi-AZ, health checks)              │
└─────┬──────────────────────┬────────────┘
      │                      │
┌─────▼─────┐        ┌──────▼──────┐
│  ECS Task │        │  ECS Task   │
│  (Go API) │        │  (Go API)   │
│  AZ-1     │        │  AZ-2       │
└─────┬─────┘        └──────┬──────┘
      │                     │
      └──────────┬──────────┘
                 │
┌────────────────▼────────────────────────┐
│  RDS PostgreSQL (Multi-AZ)              │
│  Primary (AZ-1) + Standby (AZ-2)        │
└─────────────────────────────────────────┘

┌─────────────────────────────────────────┐
│  S3 Buckets                             │
│  - dvra-resumes (CVs)                   │
│  - dvra-assets (static files)           │
└─────────────────────────────────────────┘
```

**Estimado de costos Año 1:**
- ECS Fargate (2 tasks): ~$50/mes
- RDS db.t3.micro: ~$15/mes
- ALB: ~$20/mes
- S3 + CloudFront: ~$10/mes
- **Total: ~$95/mes**

---

## 9. Monitoring y Observabilidad

### 9.1 Logging

```go
// Structured logging con zap
import "go.uber.org/zap"

logger, _ := zap.NewProduction()
defer logger.Sync()

logger.Info("User created",
    zap.Uint("user_id", user.ID),
    zap.String("email", user.Email),
    zap.String("company_id", company.ID),
)
```

### 9.2 Error Tracking

```go
// Sentry integration
import "github.com/getsentry/sentry-go"

func init() {
    sentry.Init(sentry.ClientOptions{
        Dsn: os.Getenv("SENTRY_DSN"),
        Environment: os.Getenv("ENV"),
    })
}

func HandleError(err error, context map[string]interface{}) {
    sentry.CaptureException(err)
    sentry.ConfigureScope(func(scope *sentry.Scope) {
        scope.SetContext("custom", context)
    })
}
```

### 9.3 Metrics (Futuro)

```go
// Prometheus metrics
import "github.com/prometheus/client_golang/prometheus"

var (
    httpRequestsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total HTTP requests",
        },
        []string{"method", "endpoint", "status"},
    )
)

func init() {
    prometheus.MustRegister(httpRequestsTotal)
}
```

---

## Conclusión

Esta arquitectura está diseñada para:
- ✅ Escalar de 0 a 400 empresas sin rewrite
- ✅ Mantener costos bajos en early stage
- ✅ Soportar marketplace cuando sea necesario
- ✅ Compliance GDPR/LGPD desde día 1

**Next Steps:**
1. Implementar autenticación completa (Q1)
2. Setup CI/CD con GitHub Actions
3. Deploy inicial a AWS Lightsail
4. Monitoring con CloudWatch + Sentry

---

**FIN ARQUITECTURA TÉCNICA**

> Versión 1.0 | Diciembre 8, 2025
