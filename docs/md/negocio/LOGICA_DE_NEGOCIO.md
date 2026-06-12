# Dvra - Lógica de Negocio

> **Documento Maestro de Reglas de Negocio y Estrategia**  
> Versión: 3.0 | Última actualización: Diciembre 8, 2025

---

## 📋 Tabla de Contenidos

1. [Resumen Ejecutivo](#resumen-ejecutivo)
2. [Modelo de Negocio Dual](#modelo-de-negocio-dual)
3. [Sistema de Roles y Permisos](#sistema-de-roles-y-permisos)
4. [Reglas de Negocio Fundamentales](#reglas-de-negocio-fundamentales)
5. [Pipeline de Candidatos](#pipeline-de-candidatos)
6. [Red Dvra - Marketplace](#red-dvra-marketplace)
7. [Pricing y Límites por Tier](#pricing-y-límites-por-tier)
8. [Flujos de Usuario Críticos](#flujos-de-usuario-críticos)
9. [Integraciones Estratégicas](#integraciones-estratégicas)
10. [Roadmap de Implementación](#roadmap-de-implementación)

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

### Diferenciador Único

**Dvra = ATS + Marketplace**

1. Empresa usa Dvra ATS ($49-399/mes)
2. Gestiona proceso de reclutamiento interno
3. **Cuando necesita talento externo**: "Buscar en Red Dvra"
4. Acceso a candidatos pre-evaluados
5. **Fee solo por contratación exitosa**: $3,500 USD flat

**Ventaja Competitiva Sostenible:**
- Empresas ya pagan por el ATS (ingresos recurrentes)
- Marketplace genera ingresos adicionales sin costo marginal
- Network effects: más empresas → más candidatos → más valor
- Sticky: difícil cambiar cuando tienes proceso + candidatos en una plataforma

**Competencia:**
- **Mercor/Turing**: Solo marketplace, sin herramienta
- **Greenhouse/Lever**: Solo ATS ($5k-15k/año), sin red de candidatos
- **LinkedIn**: Tablón sin curación ni herramientas
- **Dvra**: ATS económico + Marketplace = Propuesta única

---

## 2. Modelo de Negocio Dual

### 2.1 Revenue Stream 1: SaaS Subscriptions

**Pricing Mensual:**

| Tier | Precio/mes | Target |
|------|------------|--------|
| **Professional** | $49 | Startups 10-30 empleados |
| **Business** | $149 | Empresas 30-100 empleados |
| **Enterprise** | $399 | Corporativos 100+ empleados |

**Modelo de crecimiento:**
- **Año 1**: 50 empresas → $60k ARR (84% Professional, 16% Business)
- **Año 2**: 150 empresas → $216k ARR (70% Pro, 25% Bus, 5% Ent)
- **Año 3**: 400 empresas → $2.1M ARR (marketplace activo)

### 2.2 Revenue Stream 2: Marketplace Fees

**Modelo de Fee:**
- **Flat fee**: $3,500 USD por contratación exitosa
- **Cuando se cobra**: Solo cuando empresa contrata y candidato completa 90 días
- **Quién paga**: La empresa contratante
- **Incluye**: Acceso a perfil, facilitar intro, garantía de 90 días

**Proyecciones Marketplace:**

| Año | Empresas con acceso | Hires esperados | Revenue fees |
|-----|---------------------|-----------------|--------------|
| Año 1 | 0 | 0 | $0 |
| Año 2 | 150 | 60 hires | $210,000 |
| Año 3 | 400 | 250 hires | $875,000 |

**Supuestos conservadores:**
- Año 2: 40% de empresas usan marketplace al menos 1 vez
- Average hires/empresa/año: 1
- Fee efectivo después de refunds: $3,500 (sin negociación)

### 2.3 Total Revenue Projection

| Métrica | Año 1 | Año 2 | Año 3 |
|---------|-------|-------|-------|
| **SaaS ARR** | $60,000 | $216,000 | $840,000 |
| **Marketplace Revenue** | $0 | $210,000 | $875,000 |
| **Total Revenue** | **$60,000** | **$426,000** | **$1,715,000** |
| **SaaS %** | 100% | 51% | 49% |
| **Marketplace %** | 0% | 49% | 51% |

**Observación clave**: En Año 3, marketplace supera ingresos SaaS, validando modelo dual.

---

## 3. Sistema de Roles y Permisos

### 3.1 Jerarquía de Roles

```
SuperAdmin (Nivel 100)
    └── Gestiona toda la plataforma, todas las empresas
    
Admin (Nivel 50)
    └── Dueño/manager de la empresa
    └── Permisos completos dentro de su empresa
    
Recruiter (Nivel 30)
    └── Gestiona jobs y candidatos
    └── No puede invitar users ni cambiar billing
    
Hiring Manager (Nivel 20)
    └── Ve candidatos de sus jobs asignados
    └── Puede comentar y calificar
    └── No puede crear jobs
    
User (Nivel 10)
    └── Solo lectura
    └── Ve reportes básicos
```

### 3.2 Matriz de Permisos

> **Implementación:** esta matriz está aplicada en código en `internal/shared/permissions/`
> (un archivo por módulo) y se exige por endpoint vía el middleware `RequirePermission`
> en `routes.go`. Las celdas "Solo asignados / Solo sus jobs" requieren chequeo a nivel
> de recurso, pendiente de la tabla de asignaciones (ver RN-MEMB-007).

| Acción | Super | Admin | Recruiter | Hiring Mgr | User |
|--------|-------|-------|-----------|------------|------|
| **Company** |
| Ver company settings | ✅ | ✅ | ❌ | ❌ | ❌ |
| Editar company | ✅ | ✅ | ❌ | ❌ | ❌ |
| Ver billing | ✅ | ✅ | ❌ | ❌ | ❌ |
| Cambiar plan | ✅ | ✅ | ❌ | ❌ | ❌ |
| **Users & Teams** |
| Invitar usuarios | ✅ | ✅ | ❌ | ❌ | ❌ |
| Asignar roles | ✅ | ✅ | ❌ | ❌ | ❌ |
| Remover usuarios | ✅ | ✅ | ❌ | ❌ | ❌ |
| Ver team members | ✅ | ✅ | ✅ | ✅ | ✅ |
| **Jobs** |
| Crear job | ✅ | ✅ | ✅ | ❌ | ❌ |
| Editar job | ✅ | ✅ | ✅ | Solo asignados | ❌ |
| Publicar/cerrar job | ✅ | ✅ | ✅ | ❌ | ❌ |
| Ver todos los jobs | ✅ | ✅ | ✅ | Solo asignados | ✅ |
| **Candidates** |
| Crear candidato | ✅ | ✅ | ✅ | ❌ | ❌ |
| Editar candidato | ✅ | ✅ | ✅ | ❌ | ❌ |
| Ver todos candidatos | ✅ | ✅ | ✅ | Solo sus jobs | Solo sus jobs |
| **Applications** |
| Cambiar stage | ✅ | ✅ | ✅ | Solo sus jobs | ❌ |
| Calificar (rating) | ✅ | ✅ | ✅ | ✅ | ❌ |
| Agregar notas | ✅ | ✅ | ✅ | ✅ | ❌ |
| Ver aplicaciones | ✅ | ✅ | ✅ | Solo sus jobs | Solo sus jobs |
| **Red Dvra (Marketplace)** |
| Buscar candidatos red | ✅ | ✅ | ✅ | ❌ | ❌ |
| Contactar candidato | ✅ | ✅ | Solo con aprobación | ❌ | ❌ |
| Ver analytics | ✅ | ✅ | ✅ | ✅ | ✅ |

### 3.3 Reglas de Negocio de Membresías

**RN-MEMB-001: Usuario puede pertenecer a múltiples empresas**
- Un email puede tener N membresías (N companies)
- Cada membresía tiene su propio rol
- Ejemplo: juan@email.com puede ser Admin en CompanyA y Recruiter en CompanyB

**RN-MEMB-002: Membresía por defecto**
- Usuario marca 1 empresa como "default" para login inicial
- Si no tiene default, se usa la primera membresía creada

**RN-MEMB-003: SuperAdmin es especial**
- SuperAdmin NO tiene CompanyID (NULL en memberships)
- Puede acceder a cualquier empresa via admin panel
- No aparece en listado de team members de empresas

**RN-MEMB-004: Creación de Memberships (MVP - Restringido)**
- **MVP:** Solo SuperAdmin puede crear memberships (POST /admin/memberships)
- Admin de empresa NO puede agregar usuarios existentes a su empresa
- Razón: Evitar manipulación cross-company de usuarios
- Admin puede: ver, actualizar roles, eliminar memberships de su empresa
- **Fase 2:** Sistema de invitaciones por email para que admins puedan invitar de forma segura
- Workaround temporal: Admin crea nuevos usuarios, SuperAdmin asigna usuarios existentes

**RN-MEMB-005: Lifecycle de membresía**
- Estados: `pending` (invitado), `active`, `suspended`, `removed`
- Usuario invitado recibe email con magic link
- Tras aceptar: `pending` → `active`
- Admin puede suspender: `active` → `suspended` (no puede hacer login)
- Admin puede remover: cualquier estado → `removed` (soft delete)

### 3.4 Perfiles Internos vs Externos (Roadmap)

> **Estado:** diseño aprobado, pendiente de implementación. Habilita el caso de uso
> de reclutadores freelance/agencias que dan servicio a una o varias empresas cliente.

**Concepto.** El rol define *qué puede hacer* un miembro (admin, recruiter, hiring_manager);
el tipo de miembro define *hasta dónde llega*. Una empresa (ej. su departamento de RRHH)
puede invitar a un reclutador externo que trabaja solo sobre las vacantes que le asignen,
sin acceso al resto de la operación. El modelo multi-empresa actual ya lo soporta de base:
el freelancer es un usuario con N memberships (una por cliente) y usa `switch-company`
para moverse entre ellas.

**RN-MEMB-006: Tipo de miembro (interno / externo)**
- Cada membership tiene un `member_type`: `internal` (default) | `external`
- `member_type` es independiente del rol: un `recruiter` puede ser interno o externo
- El tipo lo define quien invita y solo un Admin puede cambiarlo

**RN-MEMB-007: Alcance del colaborador externo**
- Solo ve y opera las vacantes que tiene asignadas (tabla de asignación vacante ↔ recruiter)
- No accede a: configuración de empresa, billing, plan, gestión de usuarios, exportación de datos
- Visibilidad de datos sensibles del candidato (salario esperado, notas internas): configurable por empresa
- Dentro de sus vacantes asignadas opera igual que un recruiter interno:
  crear candidatos, mover el pipeline, calificar, comentar

**Matriz de permisos (recruiter interno vs externo):**

| Capacidad | Admin | Recruiter interno | Recruiter externo |
|-----------|-------|-------------------|-------------------|
| Configuración de empresa, plan, billing | ✅ | ❌ | ❌ |
| Invitar/gestionar usuarios | ✅ | ❌ | ❌ |
| Ver todas las vacantes | ✅ | ✅ | Solo asignadas |
| Crear/editar candidatos y aplicaciones | ✅ | ✅ | Solo en sus vacantes |
| Mover candidatos en el pipeline | ✅ | ✅ | Solo en sus vacantes |
| Ver datos sensibles del candidato | ✅ | ✅ | Configurable |
| Exportar datos | ✅ | ✅ | ❌ |

**RN-MEMB-008: Atribución y reporting por recruiter**
- Todo cambio de etapa en el pipeline registra `(user_id, from_stage, to_stage, timestamp)`
  (historial de pipeline); jobs, candidates y applications registran quién los creó
- Sobre ese historial se construyen reportes por recruiter (interno o externo):
  - Candidatos sourceados y aplicaciones gestionadas
  - Tiempo promedio por etapa y time-to-hire
  - Tasa de conversión a contratado
  - Comparativa interno vs externo (¿rinde el freelancer que la empresa paga?)
- **Monetización:** reportes por recruiter como feature de tiers altos (Professional/Enterprise)

**Orden de implementación recomendado** (cada paso funciona solo y habilita el siguiente):
1. Sistema de invitaciones por email (Fase 2 ya planificada; el modelo `Membership`
   ya tiene `invited_by`, `invited_at`, `joined_at`)
2. `member_type` en Membership + tabla de asignación vacante ↔ recruiter
3. Historial de pipeline con atribución por usuario
4. Reportes por recruiter sobre ese historial

---

## 4. Reglas de Negocio Fundamentales

### 4.1 Multi-Tenancy

**RN-TENANT-001: Aislamiento Completo**
- Cada empresa (Company) es un tenant independiente
- Datos 100% aislados: CompanyA NO puede ver candidatos/jobs de CompanyB
- Implementación: Todas las queries llevan `WHERE company_id = ?`

**RN-TENANT-002: Facturación Independiente**
- Cada empresa tiene su propia suscripción
- Billing cycle independiente
- No hay "org level" billing (al menos en Año 1-2)

**RN-TENANT-003: Límites por Tenant**
- Límites (jobs activos, candidatos, storage) son por empresa
- No se comparten entre empresas del mismo usuario

### 4.2 Gestión de Jobs

**RN-JOB-001: Estados de Job**
- `draft`: Creado pero no publicado
- `published`: Publicado (visible en career page si habilitado)
- `closed`: Cerrado (no acepta aplicaciones)
- `archived`: Archivado (solo lectura histórica)

**RN-JOB-002: Ownership de Jobs**
- `AssignedRecruiterID`: Recruiter responsable (opcional)
- `HiringManagerID`: Manager que toma decisión final (opcional)
- Si no hay asignados, cualquier recruiter puede gestionar

**RN-JOB-003: Límites de Jobs Activos**
- Se aplica a jobs en estado `published`
- `draft` y `closed` no cuentan para el límite
- Al alcanzar límite: No puede publicar más jobs (debe upgradear o cerrar existentes)

**RN-JOB-004: Jobs No Se Eliminan**
- Soft delete: `deleted_at IS NOT NULL`
- Razón: Mantener historial de applications
- Hard delete solo por SuperAdmin en casos extremos

### 4.3 Gestión de Candidatos

**RN-CAND-001: Unicidad de Candidato**
- Único por email dentro de la misma empresa
- Constraint: `UNIQUE(company_id, email)`
- Email puede existir en múltiples empresas (diferentes candidatos)

**RN-CAND-002: Source Tracking**
- `Source`: De dónde vino (linkedin, referral, career_page, manual)
- `SourceDetails`: Info adicional (ej: "Referido por Juan Pérez")
- Importante para analytics de canales de adquisición

**RN-CAND-003: Datos Obligatorios**
- Mínimo requerido: Email, FirstName, LastName
- Recomendado: ResumeURL (PDF en S3)
- Opcional pero valioso: GithubURL (para evaluación técnica)

**RN-CAND-004: Deduplicación en Red Dvra**
- Email y GitHub username son únicos globalmente en NetworkCandidates
- Si candidato ya existe: Enriquecer perfil, no duplicar
- Evita spam y múltiples perfiles del mismo dev

### 4.4 Gestión de Applications

**RN-APP-001: Pipeline Stages**
- `applied`: Candidato aplicó
- `screening`: Screening telefónico/inicial
- `technical`: Evaluación técnica
- `offer`: Oferta extendida
- `hired`: ¡Contratado! 🎉
- `rejected`: Rechazado en cualquier stage

**RN-APP-002: Transiciones Permitidas**
```
applied → screening
screening → technical | rejected
technical → offer | rejected
offer → hired | rejected
hired → [ESTADO FINAL]
rejected → [ESTADO FINAL]
```

**RN-APP-003: Timestamps Automáticos**
- `AppliedAt`: Se setea en creación (auto now())
- `RejectedAt`: Se setea cuando stage → rejected
- `HiredAt`: Se setea cuando stage → hired

**RN-APP-004: Rating System**
- Rating: 1-5 estrellas (nullable)
- Puede cambiar en cualquier momento
- Usado para ranking interno de candidatos

**RN-APP-005: Un Candidato Puede Aplicar a Múltiples Jobs**
- Mismo candidato puede tener N applications (diferentes jobs)
- Cada application es independiente (puede estar en diferentes stages)

---

## 5. Pipeline de Candidatos

### 5.1 Flujo Visual del Pipeline

```
APPLIED
  ↓
  → Candidato aplica (manual o portal público)
  
SCREENING
  ↓
  → Recruiter hace screening telefónico
  → Rating inicial (1-5)
  → Notas: "Buena comunicación, años exp coinciden"
  
  ├─→ TECHNICAL (pasa)
  └─→ REJECTED (no pasa)
  
TECHNICAL
  ↓
  → Evaluación técnica:
     - Code challenge (GitHub)
     - Live coding interview
     - System design
  
  ├─→ OFFER (pasa)
  └─→ REJECTED (no pasa)
  
OFFER
  ↓
  → Oferta extendida
  → Salary proposal, benefits, start date
  
  ├─→ HIRED (acepta)
  └─→ REJECTED (rechaza o empresa retracta)
  
HIRED / REJECTED
  → Estados finales
```

### 5.2 Métricas del Pipeline

**Métricas a trackear (analytics futuro):**
- **Conversion Rate**: % applied → hired
- **Bottleneck Detection**: Stage donde más candidatos se estancan
- **Time-to-Hire**: Días promedio de applied → hired
- **Stage Duration**: Tiempo promedio en cada stage
- **Rejection Reasons**: Por qué rechazamos (futuro: categorías/tags)

**Benchmarks industria:**
- Conversion rate saludable: 3-5% (1 hire por cada 20-30 aplicaciones)
- Time-to-hire promedio: 30-45 días para tech roles
- Bottleneck común: Technical stage (70% rechazados aquí)

### 5.3 Automatizaciones del Pipeline

**Fase 2 (Año 2):**
- Auto-rejection de inactivos: Si candidato en `screening` por >30 días → auto-move a `rejected`
- Notificaciones automáticas: Candidato cambia stage → Email automático personalizado
- Recruiter asignado → Slack notification

**Fase 3 (Año 3):**
- AI Scoring: Análisis de CV con IA → Score 0-100
- Matching: Job description vs resume → % match
- Predictive analytics: Probabilidad de hire basado en datos históricos

---

## 6. Red Dvra - Marketplace

### 6.1 Concepto y Valor Único

**Red Dvra** es una red curada de desarrolladores LATAM pre-evaluados técnicamente, integrada directamente en el ATS.

**Diferencia con competencia:**
- **LinkedIn Recruiter**: $99/mes → Base de datos sin curación
- **Mercor/Turing**: Fee 20-30% → Solo marketplace, sin herramienta
- **Dvra**: $49-399/mes ATS + $3,500 fee → Tool + Talent en un solo lugar

**Value Proposition para Empresas:**
1. Ya usan Dvra ATS para proceso interno
2. Cuando necesitan talento externo → Buscan en Red Dvra (mismo dashboard)
3. Candidatos pre-evaluados (ahorro de 15+ horas screening)
4. Fee solo por contratación exitosa (risk-free)
5. Timezone LATAM, inglés fluido, cultura remota

**Value Proposition para Candidatos:**
1. Validación técnica objetiva (GitHub analysis, challenges)
2. Visibilidad ante empresas que contratan activamente
3. Proceso de aplicación simplificado (un perfil, múltiples oportunidades)
4. Oportunidades remotas internacionales
5. Feedback constructivo de evaluaciones
6. **100% voluntario** - Opt-in, no spam, control total de datos (GDPR/LGPD compliant)

### 6.2 State Machine del Talento

**Estados de NetworkCandidate:**

```
prospect → invited → registered → evaluated → approved → featured
                                              ↓
                                         rejected
                                              ↓
                                         blacklisted (casos extremos)
```

**Transiciones:**

| Estado | Descripción | Siguiente Estados Válidos |
|--------|-------------|---------------------------|
| **prospect** | Identificado (scraping, referral), no contactado aún | invited, rejected |
| **invited** | Email de invitación enviado | registered, prospect (si no responde en 30 días) |
| **registered** | Completó registro, perfil creado | evaluated |
| **evaluated** | Evaluación técnica completada | approved, rejected |
| **approved** | Pasa evaluación, entra a red | featured, inactive |
| **featured** | Top 10% talento (destacado en búsquedas) | approved, inactive |
| **inactive** | No disponible temporalmente | approved, featured |
| **rejected** | No pasa evaluación o no cumple criterios | [FINAL] |
| **blacklisted** | Violó términos (fake profile, spam) | [FINAL] |

### 6.3 Proceso de Sourcing (Año 2)

**Q1-Q2 2027: Sourcing Inicial (200 candidatos)**

**Estrategia opt-in orgánico:**
1. Landing page pública: "Únete a la Red Dvra"
2. Formulario simple: Nombre, email, GitHub username, LinkedIn
3. Auto-invitación: Email de bienvenida con link de registro
4. Self-serve profile creation
5. Technical evaluation (2 opciones):
   - Opción A: Conectar GitHub (análisis automático de repos)
   - Opción B: Completar coding challenge (HackerRank-style, 1 hora)

**Fuentes de candidatos Año 2:**
- Referrals de candidatos actuales (incentivo: $100 si su referido es hired)
- Partnerships con bootcamps LATAM (Platzi, CoderHouse, etc)
- Content marketing (blog posts técnicos → CTA "Únete a Red")
- LinkedIn outreach a devs activos en comunidades tech
- Developer events (hackathons, meetups) con booth

**Target realista Año 2:**
- Q1: 50 candidatos aprobados
- Q2: 100 adicionales (150 total)
- Q3: 150 adicionales (300 total)
- Q4: 200 adicionales (500 total en red)

### 6.4 Evaluación Técnica

**Criterios de Evaluación (4 dimensiones):**

1. **Code Quality Score (0-100)**
   - GitHub analysis: Repos públicos, commits, PRs, code reviews
   - Métricas: Test coverage, documentation, code complexity
   - Peso: 30%

2. **Challenge Score (0-100)**
   - Coding challenge técnico (2 problemas, 90 minutos)
   - Evaluado automáticamente + revisión manual
   - Peso: 40%

3. **Communication Score (0-100)**
   - Entrevista async (video recording, 3 preguntas)
   - Evaluación de inglés, claridad, pensamiento estructurado
   - Peso: 20%

4. **Experience Score (0-100)**
   - Años de experiencia × proyectos relevantes
   - Referencias verificables (LinkedIn endorsements)
   - Peso: 10%

**Dvra Score Final:**
- Fórmula: `(CodeQuality * 0.3) + (Challenge * 0.4) + (Communication * 0.2) + (Experience * 0.1)`
- **Aprobado**: ≥ 70 puntos
- **Featured**: ≥ 85 puntos (top 10-15%)

**Proceso de Evaluación:**
1. Candidato completa registro → Status: `registered`
2. Conecta GitHub + completa challenge → Status: `evaluated`
3. Evaluador interno revisa (1-2 días) → Status: `approved` o `rejected`
4. Si aprobado y score ≥85 → `featured`
5. Notificación por email con feedback constructivo

### 6.5 Flujo de Matching (Empresa ↔ Candidato)

**Paso 1: Búsqueda**
- Empresa busca en Red Dvra (filtros: skills, seniority, country, availability)
- Ve perfiles resumidos (nombre, título, skills, Dvra Score, hourly rate)
- Perfiles featured aparecen primero

**Paso 2: Interés**
- Empresa marca candidato como "Interested"
- Crea NetworkApplication (status: `interested`)
- Sistema envía notificación a candidato: "Empresa X está interesada en tu perfil"

**Paso 3: Introducción**
- Candidato revisa empresa y job
- Acepta intro → Status: `contacted`
- Dvra facilita intro via email (doble opt-in)
- Empresa y candidato conectan directamente

**Paso 4: Proceso de Hiring**
- Empresa lleva proceso (entrevistas, offer)
- Candidato puede usar Dvra para trackear (opcional)
- Estados: `interviewing` → `offer` → `hired` o `rejected`

**Paso 5: Fee**
- Si candidato es hired y completa 90 días
- Dvra envía invoice a empresa: $3,500 USD
- Métodos de pago: Stripe, wire transfer
- Garantía: Si candidato renuncia antes de 90 días → Reemplazo gratis o 50% refund

### 6.6 Reglas de Anti-Spam

**RN-SPAM-001: Límite de Invitaciones**
- Max 3 invitaciones por candidato (lifetime)
- Si no responde a 3 invites → Status: `inactive`, no contactar más
- Razón: Evitar spam, respetar inbox

**RN-SPAM-002: Cooling Period**
- Si candidato rechaza intro → Empresa no puede contactar por 90 días
- Evita harassment de múltiples empresas

**RN-SPAM-003: Opt-out Permanente**
- Candidato puede opt-out cualquier momento
- Perfil se marca como `opted_out = true`
- No aparece en búsquedas futuras
- Data se mantiene 30 días (compliance) luego soft-delete

### 6.7 Compliance GDPR/LGPD

**RN-GDPR-001: Consentimiento Explícito**
- Candidato debe aceptar términos al registrarse
- Checkbox explícito: "Acepto que empresas vean mi perfil y me contacten"
- Timestamp + IP guardados como prueba

**RN-GDPR-002: Derecho al Olvido**
- Candidato puede solicitar eliminación de datos
- Proceso: Email a privacy@dvra.app
- Plazo: 30 días para eliminar
- Implementación: Soft delete → Hard delete tras 30 días

**RN-GDPR-003: Portabilidad de Datos**
- Candidato puede exportar sus datos (JSON)
- Incluye: Perfil, evaluaciones, aplicaciones, historial de contactos

**RN-GDPR-004: Transparencia**
- Candidato ve quién vio su perfil (log de vistas)
- Notificación cada vez que empresa marca "Interested"

### 6.8 Políticas de Retención de Datos

**NetworkCandidates (Red Dvra):**
- `prospect` (no contactado): Eliminar después de 6 meses
- `invited` (no responde): Eliminar después de 1 año
- `registered`/`approved`: Mantener mientras esté activo
- `inactive` (>2 años sin actividad): Soft delete automático
- `rejected`: Mantener 90 días, luego eliminar
- `blacklisted`: Mantener indefinidamente (prevención fraude)

**Candidates internos (ATS):**
- Empresa es dueña de sus datos
- Dvra no elimina automáticamente
- Si empresa cancela suscripción:
  - 30 días grace period para exportar
  - Luego soft delete
  - Hard delete tras 1 año

---

## 7. Pricing y Límites por Tier

### 7.1 Planes SaaS

| Feature | Free | Professional | Business | Enterprise |
|---------|------|-------------|----------|-----------|
| **Precio** | $0 | $49/mes | $149/mes | $399/mes |
| **Target** | Testing | Startups 10-30 | Scale-ups 30-100 | Corporativos 100+ |
| | | | | |
| **Jobs activos** | 2 | 10 | 50 | Ilimitados |
| **Candidatos/mes** | 50 | 500 | 2,000 | Ilimitados |
| **Storage** | 1 GB | 10 GB | 50 GB | 200 GB |
| **Team members** | 2 | 5 | 20 | Ilimitados |
| | | | | |
| **Features Core** | | | | |
| Jobs & candidates | ✅ | ✅ | ✅ | ✅ |
| Application pipeline | ✅ | ✅ | ✅ | ✅ |
| Email notifications | ✅ | ✅ | ✅ | ✅ |
| Basic analytics | ✅ | ✅ | ✅ | ✅ |
| Career page | ❌ | ✅ | ✅ | ✅ |
| Chrome extension | ❌ | ✅ | ✅ | ✅ |
| | | | | |
| **Collaboration** | | | | |
| Comments & notes | ✅ | ✅ | ✅ | ✅ |
| @mentions | ❌ | ✅ | ✅ | ✅ |
| Interview scheduling | ❌ | ✅ | ✅ | ✅ |
| | | | | |
| **Integraciones** | | | | |
| GitHub OAuth | ❌ | ✅ | ✅ | ✅ |
| Google Calendar | ❌ | ✅ | ✅ | ✅ |
| Slack | ❌ | ❌ | ✅ | ✅ |
| Zapier | ❌ | ❌ | ✅ | ✅ |
| API access | ❌ | ❌ | ❌ | ✅ |
| | | | | |
| **Advanced** | | | | |
| Custom workflows | ❌ | ❌ | ✅ | ✅ |
| Advanced analytics | ❌ | ❌ | ✅ | ✅ |
| HRIS integrations | ❌ | ❌ | ❌ | ✅ |
| SSO (SAML) | ❌ | ❌ | ❌ | ✅ |
| Dedicated support | ❌ | ❌ | ❌ | ✅ |
| | | | | |
| **Red Dvra** | | | | |
| Acceso marketplace | ❌ | ✅ | ✅ | ✅ |
| Búsqueda básica | ❌ | ✅ | ✅ | ✅ |
| Featured profiles | ❌ | ❌ | ✅ | ✅ |
| Búsqueda avanzada | ❌ | ❌ | ✅ | ✅ |
| Dedicated sourcer | ❌ | ❌ | ❌ | ✅ |

### 7.2 Enforcement de Límites

**RN-LIMIT-001: Jobs Activos**
- Count: `SELECT COUNT(*) WHERE company_id = ? AND status = 'published' AND deleted_at IS NULL`
- Al intentar publicar job: Check vs límite del tier
- Si excede: Error + mensaje "Upgrade to Business plan to publish more jobs"

**RN-LIMIT-002: Candidatos por Mes**
- Count: `SELECT COUNT(*) WHERE company_id = ? AND EXTRACT(MONTH FROM created_at) = CURRENT_MONTH`
- Al intentar crear candidato: Check vs límite
- Si excede: Error + mensaje "You've reached your monthly candidate limit"
- Reset: Primer día de cada mes

**RN-LIMIT-003: Storage**
- Tracking: `companies.storage_used_mb` (actualizado en cada file upload)
- Al subir CV: Check filesize + current storage vs límite
- Si excede: Error + mensaje "Storage limit reached. Delete old files or upgrade."

**RN-LIMIT-004: Team Members**
- Count: `SELECT COUNT(*) FROM memberships WHERE company_id = ? AND status = 'active'`
- Al invitar usuario: Check vs límite
- Si excede: Prompt para upgrade

**RN-LIMIT-005: Soft Limits vs Hard Limits**
- **Hard limits**: Jobs, Storage (bloquea acción)
- **Soft limits**: Candidatos/mes, Team members (warning + permite exceder 10%, luego bloquea)
- Razón: UX, evitar fricción en momentos críticos

### 7.3 Lógica de Upgrades

**RN-UPGRADE-001: Upgrade Inmediato**
- Usuario clickea "Upgrade to Business"
- Stripe subscription se actualiza
- Límites se actualizan inmediatamente
- Proration: Se cobra diferencia prorrateada del mes

**RN-UPGRADE-002: Downgrade al Final del Ciclo**
- Usuario clickea "Downgrade to Professional"
- Cambio se agenda para fin del billing cycle
- Hasta entonces: Mantiene features del plan actual
- Email de confirmación enviado

**RN-UPGRADE-003: Exceso al Downgrade**
- Si empresa tiene 15 jobs publicados y downgradea a Professional (límite: 10)
- Opciones:
  - A: Empresa cierra 5 jobs antes de downgrade
  - B: Sistema auto-cierra los 5 jobs más antiguos (con warning)
- Implementación: Opción A (empresa decide)

---

## 8. Flujos de Usuario Críticos

### 8.1 Onboarding de Nueva Empresa

**Paso 1: Signup**
- Usuario va a /signup
- Completa: Email, Password, FirstName, LastName
- Se crea: User + Company (nombre temporal: "FirstName's Company")
- Se crea: Membership (role: admin, status: active)
- Auto-login + redirect a onboarding wizard

**Paso 2: Company Setup (Wizard)**
- Screen 1: Nombre de empresa, logo (opcional), industry
- Screen 2: Company size, location, timezone
- Screen 3: Seleccionar plan (14-day trial gratis en todos)
- Screen 4: Payment method (Stripe) - solo si no es trial

**Paso 3: First Job Creation (Tutorial)**
- Guided tour: "Crea tu primer job posting"
- Pre-populate ejemplo: "Senior Backend Developer"
- Usuario edita y publica

**Paso 4: Invite Team**
- Prompt: "Invita a tu equipo" (opcional, puede skip)
- Ingresar emails + roles
- Invitations enviadas

**Paso 5: Onboarding Complete**
- Redirect a dashboard
- Show quick tips (tooltips)
- Marcar `company.onboarded_at = NOW()`

**Tiempo estimado:** 5-10 minutos

### 8.2 Aplicación de Candidato (Career Page)

**Paso 1: Descubrimiento**
- Candidato llega a career page: `dvra.app/c/company-slug`
- Ve lista de jobs publicados
- Filtra por: Location, Remote, Employment Type

**Paso 2: Ver Job**
- Click en job → Página detalle
- Muestra: Título, descripción, requirements, benefits, salary (si público)
- CTA: "Apply Now"

**Paso 3: Application Form**
- Campos: FirstName, LastName, Email, Phone (optional)
- Upload resume (PDF, max 5MB)
- LinkedIn URL (optional)
- GitHub URL (optional)
- Cover letter (optional, textarea)

**Paso 4: Submit**
- Validación: Email válido, resume uploaded
- Se crea: Candidate (si no existe) + Application (stage: applied)
- Email de confirmación a candidato
- Email de notificación a recruiter asignado

**Paso 5: Thank You**
- Página de agradecimiento
- Mensaje: "Gracias por aplicar. Revisaremos tu perfil y te contactaremos pronto."
- Tracking code (futuro: para candidate portal)

**Tiempo estimado:** 3-5 minutos

### 8.3 Proceso de Screening (Recruiter)

**Paso 1: Revisar Aplicaciones**
- Recruiter abre dashboard
- Ve pipeline con stages
- Filtra: `stage = applied`, ordenado por `applied_at DESC`

**Paso 2: Abrir Perfil**
- Click en candidato
- Ve: Resume (PDF viewer), LinkedIn, GitHub
- Lee cover letter y notas previas (si las hay)

**Paso 3: Screening Call**
- Recruiter agenda llamada (Google Calendar integration)
- Durante llamada: Toma notas en app
- Evalúa: Experience, communication, culture fit

**Paso 4: Decision**
- Rating: 3-5 estrellas → Move a `screening` stage
- Rating: 1-2 estrellas → Move a `rejected` stage
- Agrega nota: "Excelente comunicación, 5 años exp en Go"

**Paso 5: Notificación**
- Email automático a candidato:
  - Si avanza: "¡Buenas noticias! Avanzas a evaluación técnica"
  - Si rechazado: "Gracias por tu interés. En esta ocasión..." (template personalizable)

### 8.4 Hiring Decision

**Paso 1: Final Interview**
- Candidato llega a stage `offer`
- Hiring Manager + Recruiter + Team hacen entrevista final
- Todos agregan notas y ratings

**Paso 2: Offer Preparation**
- Recruiter marca "Prepare Offer"
- Llena: Salary, benefits, start date, contract details
- Guarda como draft

**Paso 3: Approval (si aplica)**
- Si empresa tiene workflow: Offer requiere aprobación de Admin
- Admin recibe notificación
- Aprueba o rechaza con feedback

**Paso 4: Extend Offer**
- Recruiter envía offer letter (PDF generado o manual upload)
- Email a candidato con link de aceptación
- Stage → `offer`

**Paso 5: Acceptance**
- Candidato acepta o rechaza
- Si acepta: Stage → `hired`, `hired_at = NOW()`
- Si rechaza: Stage → `rejected`, nota automática "Candidate declined offer"

**Paso 6: Post-Hire**
- Celebración interna (confetti animation 🎉)
- Export candidate data a HRIS (futuro: integration)
- Job status puede cambiar a `closed` (si solo 1 posición)

---

## 9. Integraciones Estratégicas

### 9.1 Integraciones Año 1 (MVP)

**Stripe (Payments)**
- Para: Billing, subscriptions, invoicing
- Uso: Upgrade/downgrade planes, marketplace fees
- Implementación: Stripe Checkout + Webhooks

**SendGrid (Email)**
- Para: Transactional emails (invitations, notifications, application confirmations)
- Uso: Templates personalizables, tracking
- Volumen estimado: 10k emails/mes (Año 1)

**AWS S3 (Storage)**
- Para: CVs/resumes, company logos
- Uso: Secure upload, CDN via CloudFront
- Bucket structure: `dvra-resumes/{company_id}/{candidate_id}/resume.pdf`

**Google Calendar (Scheduling)**
- Para: Interview scheduling
- Uso: OAuth, crear eventos, invitar candidatos
- Alternativa: Calendly embed

### 9.2 Integraciones Año 2 (Growth)

**GitHub OAuth**
- Para: Candidate sourcing, technical evaluation
- Uso: Import profile, analyze repos, contribution graph
- Permissions: `read:user`, `repo` (public repos only)

**LinkedIn OAuth**
- Para: Import candidate profile, enrich data
- Uso: Auto-fill application form
- Permissions: `r_liteprofile`, `r_emailaddress`

**Slack**
- Para: Notificaciones de equipo
- Uso: New application → Post en canal #recruiting
- Configuración: Webhook URL por empresa

**Zapier**
- Para: Conectar con 1000+ apps
- Uso: Triggers (new candidate, application stage change)
- Implementación: Zapier Platform API

### 9.3 Integraciones Año 3 (Enterprise)

**BambooHR / Gusto (HRIS)**
- Para: Sync hired candidates a sistema de RRHH
- Uso: Auto-create employee record
- Datos: FirstName, LastName, Email, StartDate, Position

**SSO (SAML)**
- Para: Enterprise single sign-on
- Proveedores: Okta, Auth0, Azure AD
- Beneficio: Seguridad, compliance

**API Pública**
- Para: Custom integrations, partners
- Uso: Access jobs, candidates, applications programmatically
- Rate limits: 100 req/min (Professional), 500 req/min (Enterprise)

---

## 10. Roadmap de Implementación

### Fase 0: Fundación (Completado ✅)
- Arquitectura multi-tenant
- Modelos de datos core
- Sistema de roles básico
- CRUD entidades principales
- JWT authentication

### Fase 1: MVP ATS (Q1-Q2 2026)
**Objetivo:** Primeros 20 clientes pagando

**Q1 (Ene-Mar):**
- Frontend React completado
- Application pipeline visual
- Email notifications
- Onboarding wizard
- Career page básico
- Billing integration (Stripe)

**Q2 (Abr-Jun):**
- Interview scheduling
- Chrome extension (LinkedIn import)
- Candidate rating system
- Analytics básico
- Mobile-responsive

**Success criteria:**
- 20 empresas activas
- $1,500 MRR
- Churn < 10%
- NPS > 30

### Fase 2: Growth & Retention (Q3-Q4 2026)
**Objetivo:** Escalar a 50 empresas, preparar marketplace

**Q3 (Jul-Sep):**
- Slack integration
- Email templates personalizables
- GitHub OAuth (preview)
- Reportes avanzados
- Perfiles internos/externos + invitaciones (ver sección 3.4): habilita recruiters freelance/agencias
- Reportes por recruiter (atribución de pipeline, ver RN-MEMB-008) como feature Professional/Enterprise
- Customer success tools (health scoring)

**Q4 (Oct-Dic):**
- Zapier integration
- Career page avanzado (SEO, ATS API)
- NetworkCandidate model (foundation)
- Sourcing landing page
- API documentation (internal)

**Success criteria:**
- 50 empresas activas
- $5,000 MRR
- Churn < 8%
- NPS > 40

### Fase 3: Marketplace Launch (Q1-Q2 2027)
**Objetivo:** Lanzar Red Dvra con primeros 200 candidatos aprobados

**Q1:**
- Sourcing process completado
- Technical evaluation pipeline
- Candidate registration flow
- GitHub analyzer (auto-scoring)
- Admin panel para evaluators

**Q2:**
- Marketplace search (empresa busca candidatos)
- Matching algorithm v1
- NetworkApplication flow
- Fee invoicing automático
- Garantía de 90 días

**Success criteria:**
- 200 candidatos aprobados en red
- 10 hires via marketplace
- $35,000 en marketplace fees
- 150 empresas activas en ATS

### Fase 4: Scale (Q3-Q4 2027)
**Objetivo:** Network effects, escalar marketplace

**Q3:**
- AI-powered matching
- Referral program para candidatos
- Partnerships con bootcamps
- Featured profiles (premium candidatos)
- Analytics de marketplace

**Q4:**
- Multi-idioma (Portuguese para Brasil)
- Enterprise features (SSO, HRIS integrations)
- API pública beta
- Mobile app (candidate-facing)

**Success criteria:**
- 500 candidatos en red
- 60 hires via marketplace
- $210,000 en marketplace fees
- 200 empresas activas

### Fase 5: Dominación LATAM (2028+)
- Expansion a Brasil (mercado 10x Colombia)
- Verticales especializados (DevOps, Data, Mobile)
- White-label ATS para RPOs
- Dvra Academy (bootcamp propio)
- **Exit o Series A**: $10M+ valuación

---

## Conclusión: Claves del Éxito

### 1. Execution Speed
- Ship rápido, iterar más rápido
- Weekly releases
- No perfeccionismo: 80% es suficiente

### 2. Customer Obsession
- Hablar con users 10+ veces/semana
- Respuesta soporte < 4 horas
- Feature decisions data-driven

### 3. Focus Brutal
- Año 1: 100% ATS, 0% marketplace
- Año 2: 70% ATS, 30% marketplace foundation
- Año 3: 50/50

### 4. Unit Economics Sólidos
- CAC < $150
- LTV/CAC > 6x
- Churn < 8%
- Gross margin > 85%

### 5. Network Effects
- Más empresas → Más candidatos aplican
- Más candidatos en red → Más valor para empresas
- Más hires exitosos → Más referrals

---

**"Construir algo que la gente quiera, cobrar por ello, y no quedarse sin dinero."**

---

**FIN LÓGICA DE NEGOCIO**

> Versión 3.0 | Diciembre 8, 2025  
> Documentación técnica: Ver [`../tecnico/ARQUITECTURA.md`](../tecnico/ARQUITECTURA.md)  
> Plan operativo Año 1: Ver [`PLAN_DE_NEGOCIO_ANO1.md`](./PLAN_DE_NEGOCIO_ANO1.md)
