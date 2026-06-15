# Dvra — Plan de Negocio

> **Plan operativo y financiero — Año 1 (2026) + visión multi-año**
> Consolidado: Junio 2026 | Reemplaza a: `BUSINESS_PLAN_YEAR1.md`
>
> 📌 Documentos relacionados:
> - Reglas y modelo de negocio → [01_LOGICA_DE_NEGOCIO.md](./01_LOGICA_DE_NEGOCIO.md)
> - Flujo de uso de la aplicación → [03_FLUJO_DE_USO.md](./03_FLUJO_DE_USO.md)
> - Detalle técnico del API → [04_DOCUMENTACION_TECNICA_API.md](./04_DOCUMENTACION_TECNICA_API.md)

---

## 1. Resumen Ejecutivo

### Misión Año 1
Validar el modelo SaaS ATS en el mercado LATAM con **50 empresas pagando** y **$60,000 USD ARR** al cierre de Q4 2026.

### Foco estratégico
- **Producto:** MVP funcional con features core del ATS.
- **Mercado:** startups tech de 11–50 empleados en Colombia, México y Argentina.
- **Revenue:** 100% SaaS (el marketplace se pospone al Año 2).
- **Operación:** bootstrapped, sin equity, crecimiento orgánico.

### Estado del producto
- Backend Go funcional (API multi-tenant con auth JWT, pipeline de candidatos, dashboard, career page pública, planes, ubicaciones).
- Base de datos PostgreSQL estructurada y con seeders.
- Pendiente: frontend, billing (Stripe), emails transaccionales, enforcement de límites de plan.

### Meta final Año 1 (Diciembre 2026)
- **50+ empresas activas** | **$60,000 ARR** ($5,000 MRR en diciembre)
- **Churn < 10% mensual** | **NPS > 40** | **Product-market fit validado**

---

## 2. Objetivos del Año 1

### 2.1 Revenue

| Trimestre | MRR target | Empresas activas |
|---|---|---|
| Q1 2026 | $500 | 5 |
| Q2 2026 | $1,500 | 15 |
| Q3 2026 | $3,000 | 30 |
| Q4 2026 | $5,000 | 50 |

**Asunciones:** ~84% en plan Professional, ~16% en Business/superior; churn promedio 8% mensual; upgrade rate 5% mensual.

### 2.2 Producto (por trimestre)

| Q | Entregables |
|---|---|
| **Q1** | Frontend React (dashboard + jobs + candidates), pipeline visual Kanban, email notifications, onboarding wizard, beta privada con 5 empresas, billing Stripe |
| **Q2** | Collaborative hiring, interview scheduling, rating/scoring, reportes básicos, Chrome extension (import LinkedIn) |
| **Q3** | Email templates personalizables, career page embeddable, GitHub OAuth, export CSV/Excel, mobile-responsive |
| **Q4** | Slack notifications, Zapier, sourcing tools, analytics avanzado, fundación del marketplace (modelo NetworkCandidate) |

### 2.3 Marketing

| Q | Visitantes/mes objetivo |
|---|---|
| Q1 | 500 |
| Q2 | 2,000 |
| Q3 | 5,000 |
| Q4 | 10,000 |

Conversión objetivo: 3% visitante → trial; 25% trial → pago.

### 2.4 Operación
- Onboarding < 15 minutos | Uptime > 99.5% | API p95 < 200ms
- Soporte: primera respuesta < 4 horas | 100% de features core documentadas

---

## 3. Timeline Mensual (Año 1)

| Mes | Tema | Foco | Resultado esperado |
|---|---|---|---|
| **Ene** | Foundation | MVP deployado (AWS), landing page, lista de 100 empresas target | 2 betas, $0 MRR |
| **Feb** | First Paying Customers | Kanban drag&drop, Stripe, onboarding wizard; 20 demos | 5 clientes, ~$245 MRR |
| **Mar** | PMF Signals | Collaborative hiring, rating, reportes; partnerships, 2 testimonials | 13 clientes, ~$637 MRR, churn <10% |
| **Abr** | Scale Acquisition | Scheduling, Chrome ext, LinkedIn ads $500/mes, webinar | 23 clientes, ~$1,323 MRR |
| **May** | Feature Expansion | Career page widget, búsqueda avanzada, bulk actions; NPS survey | 31 clientes, ~$1,820 MRR, NPS >35 |
| **Jun** | Mid-Year Review | GitHub OAuth, Slack, exports; win-back de churned | 41 clientes, ~$2,450 MRR (50% del target anual) |
| **Jul** | Automation | Zapier, email sequences, custom fields; podcast tour | 48 clientes, ~$3,150 MRR |
| **Ago** | Community | Comunidad Slack/Discord, analytics v2; customer success check-ins | 53 clientes, ~$3,600 MRR, 100 miembros |
| **Sep** | Upsell & Expansion | Reportes custom, compliance/GDPR exports, API docs; planes anuales (2 meses gratis) | 57 clientes, ~$4,200 MRR |
| **Oct** | Retention | In-app chat, usage analytics, churn prevention proactiva | 60 clientes, ~$4,500 MRR, churn <6% |
| **Nov** | Prepare for Scale | DB optimization, fundación marketplace, security audit, infra ECS | 63 clientes, ~$4,800 MRR |
| **Dic** | Close Strong | Polish, annual reports, marketplace beta preview, renewals | **65 clientes, $5,200 MRR → $62,400 ARR** |

---

## 4. Go-to-Market

### 4.1 ICP (Ideal Customer Profile)

**Empresa:** startup tech / software house / agencia digital, 11–50 empleados, Colombia–México–Argentina, Serie A o bootstrapped con revenue. Pain: usan Excel/Notion y necesitan profesionalizar el hiring.

**Buyer persona:** Head of People / HR Manager / CEO–CTO (en startups pequeñas), 28–40 años, tech-savvy, con autoridad de presupuesto. Motivación: contratar más rápido, organizar el pipeline, colaborar en equipo.

### 4.2 Value proposition

> *"Dvra es el ATS diseñado para startups tech en LATAM que te ayuda a contratar más rápido, organizar tu pipeline de candidatos y colaborar con tu equipo, sin pagar miles de dólares al mes."*

**Diferenciadores:** pricing LATAM-friendly (3–5x más barato que herramientas US) · setup en minutos · built for tech hiring (GitHub, coding challenges) · español-first · remote-first ready.

### 4.3 Canales de adquisición

| Canal | Cuándo | Mecánica | Volumen esperado | CAC |
|---|---|---|---|---|
| **LinkedIn outreach (founder-led)** | Q1–Q2 | 200 mensajes/mes personalizados → 20 demos → 5 clientes/mes | 5/mes | $0 (tiempo) |
| **Content marketing / SEO** | Q2+ | 3–4 posts/mes en keywords long-tail ("ats económico latinoamérica", "ats vs excel") | ~10 clientes/mes (mes 6+) | $200/mes freelance |
| **LinkedIn Ads** | Q2–Q4 | Targeting HR/founders 11–50 LATAM; $200→$1,000/mes | ~17 clientes/mes | ~$41 |
| **Partnerships** | Q1+ | Aceleradoras y comunidades (Platzi, 500 LATAM, NXTP): descuentos exclusivos + co-webinars, 20% comisión recurrente opcional | ~50 clientes/año | $0 directo |
| **Referral program** | Q2+ | 1 mes gratis para referente y referido | ~10/año | $980/año en créditos |

### 4.4 Proceso de ventas

1. **Lead gen** — outreach, inbound, referrals, partnerships. Calificación BANT. CRM: HubSpot free.
2. **Demo/Trial** — 14-day free trial self-service, o demo call de 30 min + trial extendido. Email sequence días 1/3/7/13.
3. **Conversión** — onboarding call de 20 min, definición de success criteria.
4. **Retención** — check-in semanal la primera semana, mensual después; NPS trimestral.

**Sales cycle objetivo:** 10–17 días (discovery → trial: 0–3 días; trial → pago: 7–14 días).

### 4.5 Funnel anual

```
Visitantes web: 30,000/año
  ↓ 10% → Trials: 3,000
  ↓ 25% → Clientes pagando (acumulado bruto)
  − churn ~7%/mes → 65 empresas activas al cierre
```

---

## 5. Modelo Financiero

### 5.1 Proyección de ingresos (mensual)

| Mes | Nuevas | Activas | MRR |
|---|---|---|---|
| Ene | 2 | 2 | $98 |
| Feb | 5 | 7 | $343 |
| Mar | 8 | 15 | $735 |
| Abr | 10 | 24 | $1,274 |
| May | 8 | 31 | $1,820 |
| Jun | 10 | 41 | $2,450 |
| Jul | 7 | 48 | $3,150 |
| Ago | 5 | 53 | $3,600 |
| Sep | 4 | 57 | $4,200 |
| Oct | 3 | 60 | $4,500 |
| Nov | 3 | 63 | $4,800 |
| Dic | 2 | 65 | $5,200 |

**ARR proyectado al cierre:** $5,200 × 12 = **$62,400** (104% del target). Upside +20%: 78 empresas / $74,880 ARR.

### 5.2 Estructura de costos (promedio mensual Año 1)

| Categoría | Costo/mes | Anual |
|---|---|---|
| **Infraestructura** (AWS, RDS, S3, CDN, SendGrid) | $215 | $2,580 |
| **Software/Tools** (Stripe fees, dominio; Sentry/HubSpot free) | $153 | $1,836 |
| **Marketing** (LinkedIn ads, content freelance, tools) | $850 | $10,200 |
| **Operacional** (soporte outsourced, legal/contabilidad) | $400 | $4,800 |
| **TOTAL** | **$1,618** | **$19,416** |

**COGS/Revenue:** ~7% (muy saludable para SaaS).

### 5.3 P&L proyectado

| | Q1 | Q2 | Q3 | Q4 | Año 1 |
|---|---|---|---|---|---|
| Revenue | $1,176 | $7,368 | $11,820 | $14,500 | **$34,864** |
| Costos | $3,504 | $4,704 | $5,304 | $5,904 | $19,416 |
| **EBITDA** | −$2,328 | $2,664 | $6,516 | $8,596 | **$15,448** |
| Margen | −198% | 36% | 55% | 59% | **44%** |

- ✅ Break-even en **Mes 5 (Mayo 2026)**; cash-positive desde Q2.
- ✅ Runway necesario: **$5,000 USD** (cubre pérdidas de Q1).

### 5.4 Unit economics

| Métrica | Valor Año 1 | Target | Proyección Año 2 |
|---|---|---|---|
| CAC | $152 | < $150 | $100 |
| ARPU | $65/mes | — | — |
| Lifetime (churn 7%) | 14.3 meses | — | 20 meses (churn 5%) |
| **LTV** | **$930** | — | $1,300 |
| **LTV/CAC** | **6.1x** | > 3x | 13x |
| Payback | 2.5 meses | < 12 meses | — |

---

## 6. Plan de Producto (esfuerzo estimado)

| Q | Features clave (P0) | Horas |
|---|---|---|
| Q1 | Jobs CRUD, candidates DB, pipeline stages, email notifications, dashboard | 220h |
| Q2 | Interview scheduling, Chrome extension, rating, bulk actions, reportes | 180h |
| Q3 | Career page embeddable, email templates, GitHub OAuth, Slack, mobile | 185h |
| Q4 | Zapier, analytics avanzado, sourcing tools, API pública (docs), marketplace foundation | 255h |
| **Total** | | **840h (~21 semanas full-time)** |

**Decisiones build vs buy:** build el core ATS (diferenciador); comprar/integrar scheduling (Calendly API), pagos (Stripe), emails (SendGrid).

**Tech stack:** Frontend React + TypeScript + TailwindCSS · Backend Go (existente) · PostgreSQL · AWS Lightsail → ECS.

**Filosofía de diseño:** simplicidad > features · time-to-value < 10 min · mobile-friendly desde día 1 · accesibilidad WCAG AA.

---

## 7. Equipo y Recursos

### Estructura Año 1

- **Founder/CEO:** desarrollo (60%), ventas (20%), estrategia/ops (20%). ~60 h/semana. Compensación Año 1: $0 (bootstrapped); Año 2: $3,000/mes al llegar a $10k MRR.
- **Contractors:** frontend dev Q1–Q2 (200h @ $30/h = $6,000), content writer Q2–Q4 ($200/mes), soporte Q3–Q4 ($300/mes). **Total labor externo: $9,600.**

### Time allocation del founder

| Actividad | Q1 | Q2 | Q3 | Q4 |
|---|---|---|---|---|
| Desarrollo | 80% | 60% | 50% | 40% |
| Sales/Marketing | 15% | 30% | 35% | 40% |
| Operaciones | 5% | 10% | 15% | 20% |

---

## 8. Métricas y KPIs

### North Star Metric: **MRR** — Target Año 1: $5,000 (Diciembre 2026)

| Categoría | Métrica | Target Q4 |
|---|---|---|
| **Growth** | Total empresas activas | 65 |
| | MRR growth rate | 20% mensual |
| **Producto** | WAU | 70% de usuarios |
| | Time-to-first-value | < 10 min |
| | Feature adoption | 60% usan ≥3 features |
| **Sales/Mkt** | Trial signups/mes | 150 |
| | Trial → paid | 25% |
| | CAC | < $150 |
| **Customer Success** | Churn mensual | < 6% |
| | NPS | > 50 |
| | Resolución de tickets | < 24h |
| **Financial** | Gross margin | > 90% |
| | CAC payback | < 3 meses |
| | LTV/CAC | > 6x |

**Cadencia de revisión:** diaria (MRR, signups, churn) · semanal (pipeline de ventas, demos) · mensual (P&L, cohorts, NPS) · trimestral (strategic review, roadmap).

**Herramientas:** Google Analytics + Mixpanel · HubSpot (free) · Stripe Dashboard · dashboard interno de Dvra.

---

## 9. Riesgos y Mitigación

| # | Riesgo | Prob. | Impacto | Mitigación |
|---|---|---|---|---|
| 1 | No alcanzar product-market fit | 40% | Crítico | 50+ entrevistas antes de features grandes; weekly user interviews; pivotar en 3 meses si no mejora; plan B: nicho más específico |
| 2 | Competencia entra a LATAM (Greenhouse/Lever) | 30% | Alto | Moat con marketplace (Año 2); brand español-first; pricing 3–5x más barato; community lock-in; partnerships exclusivos |
| 3 | Crecimiento lento (< 50 empresas) | 50% | Medio | Experimento de pricing ($29/mes temporal); freemium; partnerships agresivos; SDR part-time |
| 4 | Churn alto (> 12%/mes) | 35% | Crítico | Onboarding white-glove primeros 50; usage tracking con alertas; QBRs; incentivos de plan anual |
| 5 | Problemas técnicos / downtime | 40% | Alto | Monitoring desde día 1 (Sentry/CloudWatch); uptime 99.5%; backups diarios; status page público |

### Planes de contingencia

- **Si no llegamos a $60k ARR:** extender runway con freelancing; reducir marketing 50%; foco en retención; evaluar pre-seed ($100k) si hay tracción (>20 empresas).
- **Si superamos targets:** acelerar marketplace; contratar SDR full-time; expandir a Brasil; aplicar a YC/500.

---

## 10. Visión Multi-Año

| Año | Empresas | SaaS ARR | Marketplace | Total revenue |
|---|---|---|---|---|
| 2026 (Año 1) | 50–65 | $60k | $0 | **$60k** |
| 2027 (Año 2) | 150 | $216k | $210k (60 hires) | **$426k** |
| 2028 (Año 3) | 400 | $840k | $875k (250 hires) | **$1.7M** |

**2028+:** expansión a Brasil (mercado 10x Colombia), verticales especializados, white-label para RPOs, Dvra Academy. **Exit o Serie A: valuación $10M+.**

---

## Conclusión: 3 claves del Año 1

1. **Execution speed** — ship fast, weekly releases, sin perfeccionismo.
2. **Customer obsession** — 5+ conversaciones con usuarios/semana, soporte < 4h, decisiones por feedback real.
3. **Focus brutal** — MRR es el único objetivo; NO marketplace hasta tener 50 empresas pagando.

> *"Build something people want, charge for it, and don't run out of money."*
