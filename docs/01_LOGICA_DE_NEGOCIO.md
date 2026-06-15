# Dvra — Lógica de Negocio

> **Documento maestro de reglas de negocio**
> Consolidado: Junio 2026 | Reemplaza a: `BUSINESS_LOGIC.md`, secciones de negocio de `PLANS_MODULE.md`, `CLIENT_FLUJO_COMPLETO.md` y `SUPERADMIN_FLUJO_COMPLETO.md`
>
> 📌 Documentos relacionados:
> - Plan operativo y financiero → [02_PLAN_DE_NEGOCIO.md](./02_PLAN_DE_NEGOCIO.md)
> - Flujo de uso de la aplicación → [03_FLUJO_DE_USO.md](./03_FLUJO_DE_USO.md)
> - Detalle técnico del API → [04_DOCUMENTACION_TECNICA_API.md](./04_DOCUMENTACION_TECNICA_API.md)

---

## 1. Visión del Producto

**Dvra** es una plataforma híbrida que combina dos negocios complementarios:

1. **ATS (Applicant Tracking System)** — Software SaaS para que empresas gestionen su proceso de reclutamiento interno (jobs, candidatos, pipeline de aplicaciones).
2. **Red Dvra (Marketplace de Talento)** — Red curada de desarrolladores LATAM pre-evaluados técnicamente, accesible desde el mismo ATS. *(Fase futura — Año 2+)*

### El problema que resuelve

**Para empresas:**
- Herramientas enterprise prohibitivamente caras ($15k–50k/año).
- Alternativas económicas genéricas, sin especialización tech.
- Pierden 25+ horas semanales con herramientas dispersas (Excel, Notion, email).
- Sin acceso a talento tech pre-evaluado de LATAM.

**Para desarrolladores:**
- LinkedIn saturado sin curación real.
- No existe validación técnica objetiva de skills.
- Procesos de aplicación repetitivos y manuales.
- Falta de visibilidad del talento LATAM ante empresas internacionales.

### Diferenciador único

```
Empresa usa Dvra ATS (suscripción mensual)
    → Gestiona su reclutamiento interno
    → Cuando necesita talento externo: "Buscar en Red Dvra"
    → Accede a candidatos pre-evaluados
    → Paga fee SOLO por contratación exitosa: $3,500 USD flat
```

**Ventaja competitiva sostenible:**
- Las empresas ya pagan por el ATS (ingreso recurrente); el marketplace genera ingresos adicionales sin costo marginal.
- Network effects: más empresas → más candidatos → más valor.
- Sticky: difícil cambiar de plataforma cuando el proceso + los candidatos viven en un solo lugar.

**Competencia:**

| Competidor | Qué ofrece | Limitación |
|---|---|---|
| Mercor / Turing | Solo marketplace (fee 20–30%) | Sin herramienta de gestión |
| Greenhouse / Lever | Solo ATS ($5k–15k/año) | Sin red de candidatos, caro para LATAM |
| LinkedIn Recruiter | Base de datos ($99/mes) | Sin curación ni herramientas |
| **Dvra** | **ATS económico + Marketplace** | **Propuesta única** |

---

## 2. Modelo de Negocio Dual

### 2.1 Revenue Stream 1: Suscripciones SaaS

El ATS se vende por suscripción mensual con planes escalonados (ver §7 para límites detallados).

> ⚠️ **Nota de coherencia:** la visión comercial maneja precios objetivo de $49 / $149 / $399. Los planes **actualmente implementados y sembrados en el sistema** son: Free $0, Starter **$39.99**, Professional **$79.99**, Enterprise **$159.99**. La fuente de verdad operativa es la base de datos (`plan_seeder.go`); los precios comerciales del plan de negocio son metas de pricing a validar.

### 2.2 Revenue Stream 2: Fees de Marketplace (Año 2+)

- **Flat fee:** $3,500 USD por contratación exitosa.
- **Cuándo se cobra:** solo cuando la empresa contrata Y el candidato completa 90 días.
- **Quién paga:** la empresa contratante.
- **Incluye:** acceso al perfil, facilitación de la introducción, garantía de 90 días.
- **Garantía:** si el candidato renuncia antes de 90 días → reemplazo gratis o 50% de refund.

### 2.3 Proyección combinada

| Métrica | Año 1 | Año 2 | Año 3 |
|---|---|---|---|
| SaaS ARR | $60,000 | $216,000 | $840,000 |
| Marketplace Revenue | $0 | $210,000 | $875,000 |
| **Total** | **$60,000** | **$426,000** | **$1,715,000** |

**Observación clave:** en el Año 3 el marketplace supera al SaaS, validando el modelo dual. El detalle financiero del Año 1 está en [02_PLAN_DE_NEGOCIO.md](./02_PLAN_DE_NEGOCIO.md).

---

## 3. Sistema de Roles y Permisos

### 3.1 Jerarquía

```
SuperAdmin (Nivel 100)  → Gestiona toda la plataforma, todas las empresas.
                          ⚠️ Operativamente SEPARADO de este API (servicio aparte).
Admin (Nivel 50)        → Dueño/manager de la empresa. Permisos completos en SU empresa.
Recruiter (Nivel 30)    → Gestiona jobs, candidatos y pipeline. No gestiona usuarios ni billing.
Hiring Manager (Nivel 20) → Ve y califica candidatos de SUS jobs asignados. No crea jobs.
User (Nivel 10)         → Solo lectura, reportes básicos.
```

### 3.2 Matriz de permisos

| Acción | SuperAdmin | Admin | Recruiter | Hiring Mgr | User |
|---|---|---|---|---|---|
| **Empresa** |
| Ver/editar configuración de empresa | ✅ | ✅ | ❌ | ❌ | ❌ |
| Ver billing / cambiar plan | ✅ | Solo ver | ❌ | ❌ | ❌ |
| **Usuarios y equipo** |
| Ver team members | ✅ | ✅ | ✅ | ✅ | ✅ |
| Crear usuarios / asignar roles / remover | ✅ | ✅ | ❌ | ❌ | ❌ |
| **Jobs** |
| Ver jobs | ✅ todos | ✅ | ✅ | Solo asignados | ✅ |
| Crear / editar / publicar / cerrar | — | ✅ | ✅ | Editar solo asignados | ❌ |
| **Candidatos** |
| Ver candidatos | ✅ todos | ✅ | ✅ | Solo de sus jobs | Solo de sus jobs |
| Crear / editar | — | ✅ | ✅ | ❌ | ❌ |
| **Aplicaciones** |
| Cambiar stage | — | ✅ | ✅ | Solo de sus jobs | ❌ |
| Calificar (rating) / agregar notas | — | ✅ | ✅ | ✅ | ❌ |
| **Red Dvra (futuro)** |
| Buscar candidatos en la red | — | ✅ | ✅ | ❌ | ❌ |
| Contactar candidato | — | ✅ | Con aprobación | ❌ | ❌ |

> El SuperAdmin tiene lectura global pero **no crea** jobs/candidatos/aplicaciones porque no tiene contexto de empresa (`company_id = NULL`).

### 3.3 Reglas de membresías

- **RN-MEMB-001 — Multi-empresa:** un usuario (email) puede tener N membresías en N empresas, cada una con su propio rol. Ej.: Admin en CompanyA y Recruiter en CompanyB.
- **RN-MEMB-002 — Membresía por defecto:** el usuario marca 1 empresa como `is_default` para el login inicial; si no hay default se usa la primera membresía.
- **RN-MEMB-003 — SuperAdmin es especial:** su membresía tiene `company_id = NULL`. No aparece en listados de team members de ninguna empresa.
- **RN-MEMB-004 — Creación restringida (MVP):** solo SuperAdmin puede asignar usuarios *existentes* a empresas (evita manipulación cross-company). El Admin de empresa puede crear usuarios *nuevos* en su empresa, ver/actualizar roles y eliminar membresías propias. **Fase 2:** sistema de invitaciones por email.
- **RN-MEMB-005 — Lifecycle:** estados `pending` (invitado) → `active` → `suspended` / `removed`. Suspendido no puede hacer login; removido es soft delete.

---

## 4. Reglas de Negocio Fundamentales

### 4.1 Multi-tenancy

- **RN-TENANT-001 — Aislamiento completo:** cada empresa es un tenant independiente. CompanyA jamás ve datos de CompanyB. Toda query lleva `WHERE company_id = ?`.
- **RN-TENANT-002 — Facturación independiente:** cada empresa tiene su propia suscripción y ciclo de billing. No hay billing a nivel "organización" (Año 1–2).
- **RN-TENANT-003 — Límites por tenant:** los límites del plan (jobs activos, candidatos, storage) son por empresa, no se comparten entre empresas del mismo usuario.

### 4.2 Jobs

- **RN-JOB-001 — Estados:** `draft` (creado, no publicado) → `published` (visible en career page) → `closed` (no acepta aplicaciones). `archived` previsto a futuro.
- **RN-JOB-002 — Ownership:** `AssignedRecruiter` (responsable, opcional) y `HiringManager` (decisión final, opcional). Sin asignados, cualquier recruiter puede gestionar.
- **RN-JOB-003 — Límite de jobs activos:** aplica a jobs `published` (`draft` y `closed` no cuentan). Al alcanzar el límite del plan no se pueden publicar más jobs (upgrade o cerrar existentes).
- **RN-JOB-004 — Jobs no se eliminan:** soft delete siempre, para conservar el historial de aplicaciones. Hard delete solo SuperAdmin en casos extremos.

### 4.3 Candidatos

- **RN-CAND-001 — Unicidad:** único por email **dentro de la misma empresa** (`UNIQUE(company_id, email)`). El mismo email puede existir como candidato en varias empresas.
- **RN-CAND-002 — Source tracking:** `source` (linkedin, website, referral, job_board, direct, other) + detalles. Clave para analítica de canales.
- **RN-CAND-003 — Datos mínimos:** Email, FirstName, LastName. Recomendado: CV (resume). Valioso: GithubURL para evaluación técnica.
- **RN-CAND-004 — Deduplicación en Red Dvra (futuro):** email y GitHub username únicos globalmente en la red; si ya existe, se enriquece el perfil en vez de duplicar.

### 4.4 Aplicaciones (pipeline)

- **RN-APP-001 — Stages:** `applied` → `screening` → `technical` → `interview` → `offer` → `hired` | `rejected`.
- **RN-APP-002 — Transiciones permitidas:**

```
applied   → screening
screening → technical | rejected
technical → interview | offer | rejected
interview → offer | rejected
offer     → hired | rejected
hired / rejected → ESTADOS FINALES
```

- **RN-APP-003 — Timestamps automáticos:** `applied_at` al crear; `rejected_at` al pasar a rejected; `hired_at` al pasar a hired.
- **RN-APP-004 — Rating:** 1–5 estrellas (nullable), modificable en cualquier momento. Usado para ranking interno.
- **RN-APP-005 — Múltiples aplicaciones:** un candidato puede aplicar a N jobs; cada aplicación es independiente y avanza por su propio pipeline.

---

## 5. Pipeline de Candidatos

### 5.1 Flujo

```
APPLIED    → Candidato aplica (manual o portal público /public/jobs/:id/apply)
SCREENING  → Screening telefónico/inicial; rating inicial + notas
TECHNICAL  → Code challenge, live coding, system design
INTERVIEW  → Entrevista final con hiring manager / equipo
OFFER      → Propuesta salarial, beneficios, fecha de inicio
HIRED 🎉 / REJECTED → Estados finales
```

### 5.2 Métricas del pipeline

| Métrica | Definición | Benchmark industria |
|---|---|---|
| Conversion rate | % applied → hired | 3–5% saludable (1 hire por 20–30 aplicaciones) |
| Time-to-hire | Días promedio applied → hired | 30–45 días en roles tech |
| Bottleneck | Stage donde más candidatos se estancan | Típicamente `technical` (~70% de rechazos) |
| Stage duration | Tiempo promedio por stage | — |

El endpoint `GET /dashboard/stats` ya entrega conversion rate, time-to-hire promedio, distribución por stage, tendencias 30 días, top jobs y fuentes de candidatos.

### 5.3 Automatizaciones (roadmap)

- **Fase 2 (Año 2):** auto-rejection de inactivos (>30 días en screening), emails automáticos al cambiar de stage, notificación Slack al recruiter asignado.
- **Fase 3 (Año 3):** AI scoring de CV (0–100), matching job↔resume (% match), analítica predictiva de probabilidad de hire.

---

## 6. Red Dvra — Marketplace (Año 2+)

> ⚠️ **Estado:** diseño de negocio definido; **no implementado en el API actual**. Los modelos `NetworkCandidate` / `NetworkApplication` son parte del roadmap (fundación prevista para Q4 del Año 1).

### 6.1 Propuesta de valor

**Para empresas:** mismo dashboard del ATS; candidatos pre-evaluados (ahorro 15+ horas de screening); fee solo por éxito; timezone LATAM, inglés fluido, cultura remota.

**Para candidatos:** validación técnica objetiva (GitHub + challenges); visibilidad ante empresas que contratan activamente; un solo perfil para múltiples oportunidades; feedback constructivo; **100% voluntario** (opt-in, control total de datos, GDPR/LGPD).

### 6.2 State machine del talento

```
prospect → invited → registered → evaluated → approved → featured
                                       ↓
                                   rejected → [FINAL]
                                   blacklisted → [FINAL] (fraude/spam)
approved/featured ↔ inactive (no disponible temporalmente)
```

| Estado | Descripción |
|---|---|
| `prospect` | Identificado (scraping/referral), no contactado |
| `invited` | Email de invitación enviado (regresa a prospect si no responde en 30 días) |
| `registered` | Completó registro, perfil creado |
| `evaluated` | Evaluación técnica completada |
| `approved` | Pasó evaluación (Dvra Score ≥ 70), entra a la red |
| `featured` | Top 10–15% (Dvra Score ≥ 85), destacado en búsquedas |
| `rejected` | No pasó evaluación |
| `blacklisted` | Violó términos (perfil falso, spam) — se conserva indefinidamente para prevención de fraude |

### 6.3 Evaluación técnica — Dvra Score

| Dimensión | Peso | Cómo se mide |
|---|---|---|
| Challenge Score | 40% | Coding challenge (2 problemas, 90 min), auto + revisión manual |
| Code Quality Score | 30% | Análisis de GitHub: repos, commits, PRs, tests, documentación |
| Communication Score | 20% | Entrevista async en video (3 preguntas): inglés, claridad |
| Experience Score | 10% | Años de experiencia × proyectos relevantes, referencias |

`DvraScore = CodeQuality*0.3 + Challenge*0.4 + Communication*0.2 + Experience*0.1`
**Aprobado:** ≥ 70 | **Featured:** ≥ 85.

### 6.4 Flujo de matching empresa ↔ candidato

1. **Búsqueda** — filtros: skills, seniority, país, disponibilidad. Featured aparecen primero.
2. **Interés** — empresa marca "Interested" → se crea NetworkApplication → notificación al candidato.
3. **Introducción** — doble opt-in: el candidato acepta → Dvra facilita la intro por email.
4. **Proceso de hiring** — la empresa conduce entrevistas/oferta; estados: `interviewing` → `offer` → `hired`/`rejected`.
5. **Fee** — hired + 90 días completados → invoice de $3,500 (Stripe / wire).

### 6.5 Reglas anti-spam

- **RN-SPAM-001:** máx. 3 invitaciones por candidato (lifetime); sin respuesta → `inactive`, no contactar más.
- **RN-SPAM-002:** si el candidato rechaza una intro, esa empresa no puede contactarlo por 90 días.
- **RN-SPAM-003:** opt-out permanente en cualquier momento; desaparece de búsquedas; datos se retienen 30 días (compliance) y luego soft delete.

### 6.6 Compliance GDPR/LGPD

- **RN-GDPR-001 — Consentimiento explícito:** checkbox al registrarse; timestamp + IP guardados como prueba.
- **RN-GDPR-002 — Derecho al olvido:** solicitud a privacy@dvra.app; 30 días para eliminar (soft → hard delete).
- **RN-GDPR-003 — Portabilidad:** export JSON de perfil, evaluaciones, aplicaciones e historial de contactos.
- **RN-GDPR-004 — Transparencia:** el candidato ve quién vio su perfil y recibe notificación con cada "Interested".

### 6.7 Retención de datos

| Tipo | Política |
|---|---|
| NetworkCandidate `prospect` | Eliminar a los 6 meses |
| NetworkCandidate `invited` sin respuesta | Eliminar al año |
| `registered`/`approved` | Mientras esté activo |
| `inactive` > 2 años | Soft delete automático |
| `rejected` | 90 días y eliminar |
| `blacklisted` | Indefinido (prevención de fraude) |
| Candidatos internos del ATS | La empresa es dueña de sus datos. Si cancela: 30 días de gracia para exportar → soft delete → hard delete al año |

---

## 7. Pricing y Límites por Plan

### 7.1 Planes implementados (fuente de verdad: `plan_seeder.go`)

| | Free | Starter | Professional | Enterprise |
|---|---|---|---|---|
| **Precio/mes** | $0 | $39.99 | $79.99 | $159.99 |
| **Trial** | — | 14 días | 14 días | 30 días |
| Usuarios | 2 | 5 | 15 | Ilimitado |
| Jobs | 3 | 10 | 50 | Ilimitado |
| Candidatos | 50 | 200 | 1,000 | Ilimitado |
| Aplicaciones | 100 | 500 | 5,000 | Ilimitado |
| Storage | 1 GB | 5 GB | 20 GB | Ilimitado |
| Exportar datos | ❌ | ✅ | ✅ | ✅ |
| Custom branding | ❌ | ❌ | ✅ | ✅ |
| API access | ❌ | ❌ | ✅ | ✅ |
| Integraciones | ❌ | ❌ | ✅ | ✅ |
| Soporte | Email | Email | Priority | Dedicated |

Convención: `-1` = ilimitado. Cada plan también define `is_public` (visible en pricing page) e `is_active` (asignable).

### 7.2 Enforcement de límites

> ⚠️ **Estado actual:** los límites están **definidos en el modelo Plan pero aún NO se validan** en los services de creación (jobs, usuarios, candidatos). Es deuda funcional prioritaria. Las reglas objetivo son:

- **RN-LIMIT-001 — Jobs activos (hard):** count de jobs `published` no eliminados vs `max_jobs`; al exceder → error "Upgrade to publish more jobs".
- **RN-LIMIT-002 — Candidatos/mes (soft):** count de candidatos creados en el mes vs `max_candidates`; warning, permite exceder 10%, luego bloquea. Reset el día 1 de cada mes.
- **RN-LIMIT-003 — Storage (hard):** tracking de `storage_used_mb` actualizado en cada upload; bloquear si filesize + usado > límite.
- **RN-LIMIT-004 — Team members (soft):** count de memberships `active` vs `max_users`; prompt de upgrade.
- **RN-LIMIT-005 — Soft vs hard:** hard limits bloquean la acción (jobs, storage); soft limits avisan primero (candidatos, usuarios) para no generar fricción en momentos críticos.

### 7.3 Upgrades y downgrades

- **RN-UPGRADE-001 — Upgrade inmediato:** se actualiza la suscripción y los límites al instante; se cobra prorrateado.
- **RN-UPGRADE-002 — Downgrade al final del ciclo:** se agenda para el fin del billing cycle; hasta entonces conserva el plan actual.
- **RN-UPGRADE-003 — Exceso al downgrade:** si la empresa supera los límites del plan destino (ej. 15 jobs publicados → Starter con 10), debe cerrar el excedente antes; el sistema NO auto-cierra (decisión de la empresa).
- En el MVP los cambios de plan los ejecuta el SuperAdmin manualmente (pago offline); la integración Stripe es parte del roadmap Q1.

### 7.4 Suspensión de empresas

- El SuperAdmin puede suspender una empresa (típicamente por falta de pago): `plan_tier = "suspended"`.
- Efecto: los usuarios de esa empresa no pueden hacer login hasta reactivación (asignación de un plan válido).

---

## 8. Integraciones Estratégicas

| Fase | Integración | Para qué |
|---|---|---|
| **Año 1 (MVP)** | Stripe | Billing, suscripciones, marketplace fees |
| | SendGrid / AWS SES | Emails transaccionales (invitaciones, notificaciones, confirmaciones) |
| | AWS S3 + CloudFront | CVs y logos (`dvra-resumes/{company_id}/{candidate_id}/`) |
| | Google Calendar | Agendamiento de entrevistas (OAuth) |
| **Año 2** | GitHub OAuth | Sourcing y evaluación técnica (análisis de repos públicos) |
| | LinkedIn OAuth | Import de perfil, autocompletar aplicaciones |
| | Slack | Notificaciones de equipo (nueva aplicación → canal #recruiting) |
| | Zapier | Conexión con 1000+ apps (triggers: new candidate, stage change) |
| **Año 3 (Enterprise)** | BambooHR / Gusto | Sync de contratados al HRIS |
| | SSO (SAML) | Okta, Auth0, Azure AD |
| | API pública | Rate limits: 100 req/min (Professional), 500 req/min (Enterprise) |

---

## 9. Roadmap de Implementación (resumen)

| Fase | Periodo | Objetivo | Estado |
|---|---|---|---|
| **0 — Fundación** | Completada | Multi-tenant, modelos core, roles, CRUD, JWT, dashboard, career page pública | ✅ |
| **1 — MVP ATS** | Q1–Q2 2026 | Frontend, pipeline visual, emails, onboarding, billing Stripe → 20 clientes / $1,500 MRR | 🔄 |
| **2 — Growth** | Q3–Q4 2026 | Slack, templates, GitHub OAuth, Zapier, fundación marketplace → 50 empresas / $5,000 MRR | ⏳ |
| **3 — Marketplace** | Q1–Q2 2027 | Red Dvra con 200 candidatos aprobados, 10 hires, $35k en fees | ⏳ |
| **4 — Scale** | Q3–Q4 2027 | AI matching, referrals, multi-idioma (PT), enterprise features → 500 candidatos, 60 hires | ⏳ |
| **5 — LATAM** | 2028+ | Brasil, verticales (DevOps/Data/Mobile), white-label, Dvra Academy. Exit o Serie A ($10M+) | ⏳ |

El detalle mensual del Año 1 está en [02_PLAN_DE_NEGOCIO.md](./02_PLAN_DE_NEGOCIO.md).

---

## 10. Claves del Éxito

1. **Execution speed** — ship rápido, iterar más rápido; releases semanales; 80% es suficiente.
2. **Customer obsession** — hablar con usuarios 10+ veces/semana; soporte < 4 horas; decisiones data-driven.
3. **Focus brutal** — Año 1: 100% ATS / 0% marketplace; Año 2: 70/30; Año 3: 50/50.
4. **Unit economics sólidos** — CAC < $150, LTV/CAC > 6x, churn < 8%, gross margin > 85%.
5. **Network effects** — más empresas → más candidatos → más hires → más referrals.

> *"Construir algo que la gente quiera, cobrar por ello, y no quedarse sin dinero."*
