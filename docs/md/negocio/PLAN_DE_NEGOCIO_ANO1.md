# Dvra - Plan de Negocios Año 1

> **Plan Operativo y Financiero - Primeros 12 Meses**  
> Versión: 1.0 | Última actualización: Diciembre 8, 2025

---

## 📋 Tabla de Contenidos

1. [Resumen Ejecutivo](#resumen-ejecutivo)
2. [Objetivos del Año 1](#objetivos-del-año-1)
3. [Timeline Mensual](#timeline-mensual)
4. [Estrategia de Go-to-Market](#estrategia-de-go-to-market)
5. [Modelo Financiero](#modelo-financiero)
6. [Adquisición de Clientes](#adquisición-de-clientes)
7. [Plan de Producto](#plan-de-producto)
8. [Equipo y Recursos](#equipo-y-recursos)
9. [Métricas y KPIs](#métricas-y-kpis)
10. [Riesgos y Mitigación](#riesgos-y-mitigación)

---

## 1. Resumen Ejecutivo

### Misión Año 1
Validar el modelo SaaS ATS en el mercado LATAM con **50 empresas pagando** y generar **$60,000 USD ARR** al finalizar Q4 2026.

### Foco Estratégico
- ✅ **Producto**: MVP funcional con features core del ATS
- ✅ **Mercado**: Startups tech de 11-50 empleados en Colombia, México, Argentina
- ✅ **Revenue**: 100% SaaS (marketplace se pospone para Año 2)
- ✅ **Operación**: Bootstrapped, sin equity, crecimiento orgánico

### Estado Actual (Diciembre 2025)
- Backend Go funcional (60% completado)
- Base de datos PostgreSQL estructurada
- Autenticación JWT implementada
- Sin frontend (próximo sprint)
- Sin clientes aún

### Meta Final Año 1 (Diciembre 2026)
- **50 empresas activas** (42 en plan Professional, 8 en Business)
- **$60,000 ARR** ($5,000 MRR en diciembre)
- **Churn < 10%** mensual
- **NPS > 40**
- **Product-market fit validado**

---

## 2. Objetivos del Año 1

### 2.1 Objetivos de Revenue

| Trimestre | MRR Target | Empresas Activas | ARR Acumulado |
|-----------|------------|------------------|---------------|
| **Q1 2026** | $500 | 5 empresas | $1,500 |
| **Q2 2026** | $1,500 | 15 empresas | $10,500 |
| **Q3 2026** | $3,000 | 30 empresas | $28,500 |
| **Q4 2026** | $5,000 | 50 empresas | $60,000 |

**Asunciones:**
- 84% en plan Professional ($49/mes)
- 16% en plan Business ($149/mes)
- Churn promedio: 8% mensual
- Upgrade rate: 5% mensual (Professional → Business)

### 2.2 Objetivos de Producto

**Q1 (Ene - Mar 2026)**
- ✅ Frontend React completado (dashboard + jobs + candidates)
- ✅ Sistema de aplicaciones con pipeline visual (Kanban)
- ✅ Email notifications funcionales
- ✅ Onboarding flow optimizado
- ✅ Beta privada lanzada con 5 empresas

**Q2 (Abr - Jun 2026)**
- ✅ Collaborative hiring (múltiples recruiters)
- ✅ Interview scheduling integrado (Calendly-style)
- ✅ Candidate scoring/rating system
- ✅ Reportes básicos de pipeline
- ✅ Chrome extension para importar de LinkedIn

**Q3 (Jul - Sep 2026)**
- ✅ Email templates personalizables
- ✅ Career page embeddable (widget)
- ✅ Integración básica con GitHub (OAuth)
- ✅ Exportación de datos (CSV/Excel)
- ✅ Mobile-responsive optimizado

**Q4 (Oct - Dic 2026)**
- ✅ Slack notifications
- ✅ Zapier integration (basic)
- ✅ Candidate sourcing tools (extensiones)
- ✅ Analytics dashboard avanzado
- ✅ Preparación para marketplace (Año 2)

### 2.3 Objetivos de Marketing

- **Q1**: 500 visitantes/mes al sitio web
- **Q2**: 2,000 visitantes/mes
- **Q3**: 5,000 visitantes/mes
- **Q4**: 10,000 visitantes/mes

**Tasa de conversión objetivo**: 3% (visitante → trial) → 25% (trial → pago)

### 2.4 Objetivos de Operación

- Tiempo de onboarding: < 15 minutos
- Uptime: > 99.5%
- Response time API: < 200ms (p95)
- Soporte: < 4 horas respuesta inicial
- Documentación: 100% features core documentadas

---

## 3. Timeline Mensual

### Mes 1: Enero 2026 - "Foundation"

**Prioridad: Completar MVP funcional**

**Desarrollo (120 horas)**
- Frontend: Dashboard básico + Jobs CRUD
- Candidates: Crear, editar, ver perfil
- Applications: Submit application flow
- Auth: Login, signup, password reset
- Deploy: AWS Lightsail inicial

**Marketing (20 horas)**
- Landing page con value proposition clara
- Setup analytics (Google Analytics, Hotjar)
- Crear perfiles en redes sociales
- LinkedIn company page

**Sales (10 horas)**
- Lista de 100 empresas target (scraping LinkedIn)
- Outreach template v1
- Pricing page publicada

**Resultado esperado:**
- ✅ MVP deployado en producción
- ✅ 2 empresas en beta privada
- ✅ $0 MRR

---

### Mes 2: Febrero 2026 - "First Paying Customers"

**Prioridad: Cerrar primeros 5 clientes pagando**

**Desarrollo (100 horas)**
- Pipeline visual (drag & drop Kanban)
- Email notifications (application received, status change)
- Billing integration (Stripe)
- Onboarding wizard

**Marketing (30 horas)**
- 10 posts en LinkedIn (founder content)
- Outreach directo a 50 empresas
- Crear demo video (Loom, 3 minutos)
- SEO básico (keywords research)

**Sales (30 horas)**
- 20 demos/calls con prospects
- Cerrar 5 clientes pagando
- Recopilar feedback de beta users
- Iterar pricing si es necesario

**Resultado esperado:**
- ✅ 5 empresas pagando
- ✅ $245 MRR (5 × $49)
- ✅ 10 trials activos

---

### Mes 3: Marzo 2026 - "Product-Market Fit Signals"

**Prioridad: Validar retención y scaling de adquisición**

**Desarrollo (80 horas)**
- Collaborative hiring (invitar team members)
- Candidate rating/scoring
- Reportes básicos (pipeline conversion)
- Performance optimizations

**Marketing (40 horas)**
- Content marketing: 2 blog posts técnicos
- LinkedIn ads experiment ($200 budget)
- Partnerships outreach (5 aceleradoras, 3 comunidades tech)
- Testimonial videos (2 clientes)

**Sales (40 horas)**
- 30 demos/calls
- Cerrar 8 empresas más (total: 13)
- Implementar feedback loop
- Crear case study inicial

**Resultado esperado:**
- ✅ 13 empresas pagando
- ✅ $637 MRR
- ✅ Churn < 10%
- ✅ Al menos 1 upgrade a Business plan

---

### Mes 4: Abril 2026 - "Scale Acquisition"

**Prioridad: Doblar adquisición mensual**

**Desarrollo (70 horas)**
- Interview scheduling (calendar integration)
- Chrome extension v1 (LinkedIn import)
- Email templates system
- API documentation (para futuras integraciones)

**Marketing (50 horas)**
- Content: 3 blog posts + 1 video tutorial
- LinkedIn ads scaling ($500/mes)
- Webinar: "Cómo contratar developers remotos"
- Guest posting en blogs tech LATAM

**Sales (40 horas)**
- 40 demos/calls
- Cerrar 10 empresas más (total: 23)
- Implementar referral program (beta)
- Optimizar demo script

**Resultado esperado:**
- ✅ 23 empresas pagando
- ✅ $1,323 MRR
- ✅ 2-3 referrals orgánicos

---

### Mes 5: Mayo 2026 - "Feature Expansion"

**Prioridad: Agregar features que desbloqueen upgrades**

**Desarrollo (80 horas)**
- Career page embeddable (widget JavaScript)
- Advanced search & filters
- Bulk actions (emails, status changes)
- Mobile optimization

**Marketing (50 horas)**
- Content: 3 blog posts
- LinkedIn organic (15 posts/mes)
- Partnerships activadas (2 aceleradoras)
- Case study #2 publicado

**Sales (40 horas)**
- 35 demos/calls
- Cerrar 8 empresas más (total: 31)
- Upsell a 2 clientes existentes
- Implementar NPS survey

**Resultado esperado:**
- ✅ 31 empresas pagando
- ✅ $1,820 MRR
- ✅ 3 clientes en Business plan
- ✅ NPS > 35

---

### Mes 6: Junio 2026 - "Mid-Year Review"

**Prioridad: Consolidar product-market fit**

**Desarrollo (60 horas)**
- GitHub OAuth integration
- Slack notifications
- Data export (CSV/Excel)
- Bug fixes & polish

**Marketing (60 horas)**
- Content: 4 blog posts
- Webinar #2: "ATS para startups"
- LinkedIn ads optimization
- Email marketing setup (Mailchimp)

**Sales (50 horas)**
- 40 demos/calls
- Cerrar 10 empresas más (total: 41)
- Win-back campaign (churned users)
- Mid-year pricing review

**Resultado esperado:**
- ✅ 41 empresas pagando
- ✅ $2,450 MRR
- ✅ $14,700 ARR acumulado (50% del target anual)
- ✅ Churn estabilizado < 8%

---

### Mes 7: Julio 2026 - "Automation & Scaling"

**Prioridad: Reducir carga operativa, automatizar procesos**

**Desarrollo (70 horas)**
- Zapier integration (triggers: new candidate, new application)
- Automated email sequences
- Custom fields system
- Permissions granulares

**Marketing (50 horas)**
- Content: 3 blog posts + 1 comparison guide (vs competitors)
- LinkedIn ads ($700/mes)
- Podcast tour (3 entrevistas en podcasts tech LATAM)
- SEO optimization (backlinks campaign)

**Sales (40 horas)**
- 35 demos/calls
- Cerrar 7 empresas más (total: 48)
- Self-service trial optimizado
- Automated onboarding emails

**Resultado esperado:**
- ✅ 48 empresas pagando
- ✅ $3,150 MRR
- ✅ 20% trials convierten sin demo

---

### Mes 8: Agosto 2026 - "Community Building"

**Prioridad: Crear comunidad de usuarios activos**

**Desarrollo (60 horas)**
- Analytics dashboard v2
- Candidate sourcing tools
- API rate limiting & monitoring
- Performance improvements

**Marketing (60 horas)**
- Lanzar comunidad Slack/Discord
- Content: 4 blog posts
- Webinar #3 con invitado especial
- User-generated content campaign

**Sales (40 horas)**
- 30 demos/calls
- Cerrar 5 empresas más (total: 53)
- Community-led growth experiments
- Implementar customer success check-ins

**Resultado esperado:**
- ✅ 53 empresas pagando
- ✅ $3,600 MRR
- ✅ 100 miembros en comunidad
- ✅ 5+ testimonials públicos

---

### Mes 9: Septiembre 2026 - "Upsell & Expansion"

**Prioridad: Maximizar revenue per customer**

**Desarrollo (70 horas)**
- Advanced reporting (custom reports)
- Compliance features (GDPR exports)
- Integración API abierta (documentación)
- White-label options (para Enterprise futuro)

**Marketing (50 horas)**
- Content: 3 blog posts técnicos
- Comparison SEO pages (vs competidores)
- LinkedIn ads ($800/mes)
- Case studies video (3 clientes)

**Sales (50 horas)**
- 30 demos/calls
- Cerrar 4 empresas más (total: 57)
- Upsell campaign (Professional → Business)
- Annual plan discounts (2 meses gratis)

**Resultado esperado:**
- ✅ 57 empresas pagando
- ✅ $4,200 MRR
- ✅ 8 clientes en Business plan
- ✅ 3 clientes en plan anual

---

### Mes 10: Octubre 2026 - "Retention Focus"

**Prioridad: Reducir churn, aumentar engagement**

**Desarrollo (50 horas)**
- In-app messaging/chat
- Feature usage analytics (internal)
- Onboarding improvements (based on data)
- Quick wins features (user requests)

**Marketing (50 horas)**
- Content: 4 blog posts
- Webinar #4: Advanced ATS workflows
- Email nurture sequences
- Retargeting ads setup

**Sales (60 horas)**
- 25 demos/calls
- Cerrar 3 empresas más (total: 60)
- Customer health scoring
- Proactive churn prevention

**Resultado esperado:**
- ✅ 60 empresas pagando
- ✅ $4,500 MRR
- ✅ Churn < 6%
- ✅ Engagement score > 75%

---

### Mes 11: Noviembre 2026 - "Prepare for Scale"

**Prioridad: Preparar infraestructura para Año 2**

**Desarrollo (80 horas)**
- Database optimization (indexes, queries)
- Marketplace foundation (NetworkCandidate model)
- Infrastructure scaling (ECS Fargate)
- Security audit

**Marketing (50 horas)**
- Content: 3 blog posts + Año 2 teaser
- LinkedIn ads ($1,000/mes)
- Year-in-review campaign
- Customer appreciation program

**Sales (40 horas)**
- 25 demos/calls
- Cerrar 3 empresas más (total: 63)
- Plan anual push (Black Friday promo)
- Feedback sessions para roadmap Año 2

**Resultado esperado:**
- ✅ 63 empresas pagando
- ✅ $4,800 MRR
- ✅ Infrastructure ready para 200 empresas
- ✅ Roadmap Año 2 definido

---

### Mes 12: Diciembre 2026 - "Close the Year Strong"

**Prioridad: Alcanzar target y celebrar wins**

**Desarrollo (40 horas)**
- Bug fixes & polish
- Year-end features (annual reports para clientes)
- Marketplace beta preview
- Documentation updates

**Marketing (60 horas)**
- Content: 2 blog posts + Year in Review
- End-of-year promotions
- Customer success stories compilation
- 2027 announcements teaser

**Sales (50 horas)**
- 20 demos/calls
- Cerrar 2 empresas más (total: 65)
- Renewal campaign (annual plans)
- Thank you campaign (regalos, swag)

**Resultado esperado:**
- ✅ **65 empresas pagando** (30% por encima del target)
- ✅ **$5,200 MRR**
- ✅ **$62,400 ARR** (104% del target)
- ✅ **Churn < 5%**
- ✅ **NPS > 50**

---

## 4. Estrategia de Go-to-Market

### 4.1 Segmento Target (ICP - Ideal Customer Profile)

**Perfil de empresa:**
- **Industria**: Tech startups, software houses, agencias digitales
- **Tamaño**: 11-50 empleados
- **Geografía**: Colombia, México, Argentina (top prioridad)
- **Etapa**: Serie A, bootstrapped con revenue
- **Pain point**: Están usando Excel/Notion o herramientas genéricas, necesitan profesionalizar hiring

**Buyer persona:**
- **Título**: Head of People, HR Manager, CEO/CTO (en startups pequeñas)
- **Edad**: 28-40 años
- **Tech-savvy**: Sí, cómodos con SaaS tools
- **Budget authority**: Sí o influencia directa
- **Motivación**: Contratar más rápido, organizar pipeline, colaboración con equipo

### 4.2 Value Proposition

**Para startups tech de 11-50 empleados:**

*"Dvra es el ATS diseñado para startups tech en LATAM que te ayuda a contratar más rápido, organizar tu pipeline de candidatos y colaborar con tu equipo, sin pagar miles de dólares al mes."*

**Diferenciadores clave:**
1. **Pricing LATAM-friendly**: Desde $49/mes (vs $300+/mes de competencia US)
2. **Setup en minutos**: No requiere implementación compleja
3. **Built for tech hiring**: GitHub integration, coding challenges, tech-focused
4. **Spanish-first**: Producto, soporte y content en español
5. **Remote-first ready**: Features para equipos distribuidos

### 4.3 Canales de Adquisición

**Trimestre 1-2: Manual & Direct**
- **LinkedIn outreach** (founder-led sales): 50 mensajes/semana
- **Demo calls**: 15-20/mes
- **Referrals**: Programa de referidos (1 mes gratis)
- **Partnerships**: Aceleradoras (Platzi, Y Combinator alumni, 500 Startups LATAM)

**Trimestre 3-4: Content & Paid**
- **Content marketing**: 3-4 blog posts/mes (SEO)
- **LinkedIn ads**: $500-1,000/mes (targeting HR/founders)
- **Webinars**: 1/mes (educational, soft-sell)
- **Community**: Slack/Discord para users

**Año 2: Automated & Scaled**
- **SEO**: Posicionar en "ATS para startups", "software de reclutamiento"
- **Affiliate program**: HR consultants, coaches
- **Product-led growth**: Self-service trial optimizado
- **Marketplace**: Network effect empieza a generar inbound

### 4.4 Sales Process

**Paso 1: Lead Generation**
- Fuentes: LinkedIn outreach, website inbound, referrals, partnerships
- Calificación: BANT (Budget, Authority, Need, Timeline)
- Tool: HubSpot CRM (free tier)

**Paso 2: Demo/Trial**
- Opción A: 14-day free trial (self-service)
- Opción B: Demo call de 30 minutos + trial extendido
- Follow-up: Email sequence automatizado (día 1, 3, 7, 13)

**Paso 3: Conversion**
- Onboarding call (20 minutos)
- Success criteria definition
- Assign customer success contact

**Paso 4: Retention**
- Check-in semanal primera semana
- Check-in mensual después
- NPS survey trimestral
- Feature announcements regulares

**Sales Cycle Target:**
- **Discovery → Trial**: 0-3 días
- **Trial → Paid**: 7-14 días
- **Total**: 10-17 días promedio

---

## 5. Modelo Financiero

### 5.1 Proyección de Ingresos Año 1

| Mes | Empresas Nuevas | Empresas Activas | Churn | MRR | ARR Acumulado |
|-----|-----------------|------------------|-------|-----|---------------|
| Ene | 2 | 2 | 0 | $98 | $98 |
| Feb | 5 | 7 | 0 | $343 | $441 |
| Mar | 8 | 15 | 1 | $735 | $1,176 |
| Abr | 10 | 24 | 1 | $1,274 | $2,450 |
| May | 8 | 31 | 1 | $1,820 | $4,270 |
| Jun | 10 | 41 | 0 | $2,450 | $6,720 |
| Jul | 7 | 48 | 0 | $3,150 | $9,870 |
| Ago | 5 | 53 | 0 | $3,600 | $13,470 |
| Sep | 4 | 57 | 0 | $4,200 | $17,670 |
| Oct | 3 | 60 | 0 | $4,500 | $22,170 |
| Nov | 3 | 63 | 0 | $4,800 | $26,970 |
| Dic | 2 | 65 | 0 | $5,200 | $32,170 |
| **TOTAL** | **67** | **65** | **3** | **$5,200** | **$32,170** |

**Nota:** ARR acumulado es la suma del MRR de cada mes (no multiplica por 12). ARR final = $5,200 × 12 = **$62,400**.

**Mix de planes (final del año):**
- Professional ($49/mes): 55 empresas = $2,695/mes
- Business ($149/mes): 10 empresas = $1,490/mes
- **Total MRR**: $4,185/mes
- **ARR proyectado**: $50,220

**Upside scenario (+20%):**
- 78 empresas activas
- $6,240 MRR
- **$74,880 ARR**

### 5.2 Estructura de Costos

**Costos Fijos Mensuales (promedio Año 1):**

| Concepto | Costo/mes | Anual |
|----------|-----------|-------|
| **Infraestructura** | | |
| AWS (Lightsail → ECS) | $100 | $1,200 |
| Database (RDS) | $50 | $600 |
| S3 Storage | $30 | $360 |
| CDN (CloudFront) | $20 | $240 |
| SendGrid (emails) | $15 | $180 |
| **Subtotal Infra** | **$215** | **$2,580** |
| | | |
| **Software & Tools** | | |
| Stripe fees (2.9% + $0.30) | ~$150 | $1,800 |
| Domain + SSL | $3 | $36 |
| GitHub | $0 | $0 |
| Sentry (monitoring) | $0 | $0 |
| HubSpot CRM | $0 | $0 |
| **Subtotal Tools** | **$153** | **$1,836** |
| | | |
| **Marketing** | | |
| LinkedIn ads | $600 | $7,200 |
| Content (freelancers) | $200 | $2,400 |
| Tools (Mailchimp, etc) | $50 | $600 |
| **Subtotal Marketing** | **$850** | **$10,200** |
| | | |
| **Operational** | | |
| Soporte (outsourced) | $300 | $3,600 |
| Legal/Accounting | $100 | $1,200 |
| **Subtotal Ops** | **$400** | **$4,800** |
| | | |
| **TOTAL COSTOS FIJOS** | **$1,618** | **$19,416** |

**Costos Variables:**
- Payment processing: 3% del MRR
- Customer success time: 2 horas/mes por empresa (asumido incluido en founder time)

**COGS (Cost of Goods Sold):**
- Infraestructura + Stripe fees ≈ $368/mes
- **COGS/Revenue ratio**: ~7% (muy saludable para SaaS)

### 5.3 P&L Proyectado Año 1

| Concepto | Q1 | Q2 | Q3 | Q4 | **Total Año 1** |
|----------|----|----|----|----|-----------------|
| **Revenue** | $1,176 | $7,368 | $11,820 | $14,500 | **$34,864** |
| | | | | | |
| **Costos** | | | | | |
| Infraestructura | $645 | $645 | $645 | $645 | $2,580 |
| Software/Tools | $459 | $459 | $459 | $459 | $1,836 |
| Marketing | $1,200 | $2,400 | $3,000 | $3,600 | $10,200 |
| Operacional | $1,200 | $1,200 | $1,200 | $1,200 | $4,800 |
| **Total Costos** | **$3,504** | **$4,704** | **$5,304** | **$5,904** | **$19,416** |
| | | | | | |
| **EBITDA** | **-$2,328** | **$2,664** | **$6,516** | **$8,596** | **$15,448** |
| **Margen** | -198% | 36% | 55% | 59% | **44%** |

**Observaciones:**
- ✅ Break-even alcanzado en **Mes 5** (Mayo 2026)
- ✅ Cash-positive desde Q2
- ✅ Margen final del año: 44% (excelente para SaaS early-stage)
- ✅ Runway necesario: $5,000 USD (cubre losses de Q1)

### 5.4 Unit Economics

**CAC (Customer Acquisition Cost):**
- Marketing + Sales time / Nuevos clientes
- Año 1 promedio: $10,200 / 67 empresas = **$152 por empresa**

**LTV (Lifetime Value):**
- ARPU: $65/mes (promedio Professional + Business)
- Churn: 7% mensual
- Lifetime: 1 / 0.07 = 14.3 meses
- **LTV**: $65 × 14.3 = **$930**

**LTV/CAC Ratio:**
- $930 / $152 = **6.1x** (Target: >3x, excelente)

**Payback Period:**
- CAC / (ARPU × Gross Margin)
- $152 / ($65 × 0.93) = **2.5 meses** (Target: <12 meses, excelente)

**Proyección mejorada Año 2:**
- CAC optimizado: $100 (más eficiencia)
- Churn reducido: 5% mensual → LTV: $1,300
- LTV/CAC: **13x**

---

## 6. Adquisición de Clientes

### 6.1 Canales y Tácticas Detalladas

#### **Canal 1: LinkedIn Outreach (Founder-led)**

**Estrategia:**
- Identificar HR managers, founders, CTOs en startups tech LATAM
- Mensaje personalizado (no plantilla genérica)
- Ofrecer valor primero (contenido, insights)
- Soft ask para demo

**Playbook:**
1. Scraping: Sales Navigator search (100 leads/semana)
2. Enriquecer: Visitar perfil, ver contenido que publican
3. Connect: Request con nota personalizada
4. Nurture: Comentar en sus posts (3-5 días)
5. Outreach: DM con propuesta de valor
6. Follow-up: 2 follow-ups espaciados (3 y 7 días)

**Template ejemplo:**

```
Hola [Nombre] 👋

Vi que [empresa] está creciendo rápido (felicitaciones por [milestone reciente]).

Estoy construyendo Dvra, un ATS pensado para startups tech en LATAM 
que están profesionalizando su hiring pero no quieren pagar $300+/mes 
por herramientas gringas.

¿Te interesaría ver un demo de 15 minutos? Sin compromiso, 
y si no te sirve, te comparto algunos recursos gratis sobre hiring 
que quizás te sean útiles.

Saludos,
[Tu nombre]
```

**Volumen esperado:**
- Mensajes enviados: 200/mes
- Acceptance rate: 40% (80 aceptan)
- Response rate: 15% (30 responden)
- Demo booked: 10% (20 demos)
- Conversión: 25% (5 clientes)

**Costo:** $0 (solo tiempo founder)

---

#### **Canal 2: Content Marketing (SEO Blog)**

**Estrategia:**
- Posicionar en keywords de long-tail con buyer intent
- Resolver problemas reales (no solo hablar del producto)
- Link building orgánico (partnerships, guest posts)

**Keywords target (ejemplos):**
- "software de reclutamiento para startups" (150 búsquedas/mes)
- "ats economico latinoamerica" (80 búsquedas/mes)
- "como organizar proceso de seleccion startup" (200 búsquedas/mes)
- "mejores ats para pequeñas empresas" (120 búsquedas/mes)

**Content calendar (ejemplo Q2):**

| Semana | Tema | Keyword | CTA |
|--------|------|---------|-----|
| 1 | "10 errores al contratar developers remotos" | hiring developers remote | eBook download |
| 2 | "ATS vs Excel: ¿Cuándo es momento de cambiar?" | ats vs excel | Free trial |
| 3 | "Cómo reducir tu time-to-hire en 50%" | reducir time to hire | Demo |
| 4 | "Guía: Entrevistas técnicas efectivas" | como hacer entrevistas tecnicas | Template download |

**Distribución:**
- Blog post publicado
- LinkedIn post (orgánico)
- Newsletter (si tenemos lista)
- Repost en comunidades (Slack, Discord)
- Guest post version (con backlink)

**Volumen esperado (mes 6+):**
- 3-4 posts/mes
- 2,000 visitantes orgánicos/mes (crece lento, pero compounding)
- Conversión: 2% → 40 trials/mes
- 25% trial→paid → 10 clientes/mes

**Costo:** $200/mes (freelancer para writing)

---

#### **Canal 3: LinkedIn Ads (Paid)**

**Estrategia:**
- Targeting muy específico (job titles, industries, company size)
- Múltiples creatives (test A/B)
- Lead magnet + demo booking

**Targeting:**
- **Geografía**: Colombia, México, Argentina, Chile
- **Job titles**: HR Manager, People Ops, Head of Talent, CEO (company size 11-50)
- **Industries**: Computer Software, Internet, Tech
- **Company size**: 11-50, 51-200

**Ad formats:**
1. **Sponsored content** (feed posts)
   - Imagen + copy corto
   - CTA: "Download Free Guide" o "Book Demo"
   
2. **Message ads** (InMail)
   - Mensaje directo a inbox
   - Personalizado con variables

**Creative examples:**

*Ad 1: Pain point*
> "¿Sigues usando Excel para trackear candidatos? 🤯  
> Startups tech están cambiando a Dvra y reduciendo su time-to-hire en 50%.  
> Prueba gratis 14 días → [link]"

*Ad 2: Social proof*
> "20+ startups en LATAM ya confían en Dvra para profesionalizar su hiring.  
> Desde $49/mes. Sin contratos anuales.  
> Agenda demo → [link]"

**Budget allocation:**
- Q1: $200/mes (test)
- Q2: $500/mes (scale lo que funciona)
- Q3-Q4: $800-1,000/mes

**Volumen esperado (promedio Q2-Q4):**
- Budget: $700/mes
- CPM: $30 (LinkedIn LATAM)
- Impressions: 23,000/mes
- CTR: 1.5% → 345 clicks
- Conversión landing: 20% → 69 trials
- Trial→paid: 25% → **17 clientes/mes**

**CAC:** $700 / 17 = **$41 por cliente** (excelente)

---

#### **Canal 4: Partnerships (Aceleradoras & Comunidades)**

**Estrategia:**
- Alianzas con aceleradoras que tienen cohorts regulares de startups
- Ofrecer descuentos/trials extendidos a sus portfolio companies
- Co-marketing (webinars, workshops, contenido)

**Target partners:**
- **Aceleradoras**: Platzi Startups, 500 Startups LATAM, NXTP Ventures, Mountain Nazca
- **Comunidades**: Startup Chile, Geek Girls LatAm, tech meetups locales
- **Influencers HR**: LinkedIn influencers LATAM con audiencia HR/founders

**Value proposition para partner:**
- Herramienta útil para sus startups (sin costo para el partner)
- Commission: 20% recurring por cada referral (opcional)
- Co-branding opportunities

**Playbook:**
1. Identificar partner (research)
2. Warm intro si es posible (LinkedIn, mutual connections)
3. Pitch: "Queremos ayudar a tus startups a profesionalizar hiring"
4. Propuesta: Descuento exclusivo (2 meses gratis) + co-webinar
5. Onboarding: Landing page custom con código promo
6. Seguimiento: Monthly report de usage

**Volumen esperado:**
- 5 partnerships activos (año 1)
- 10 startups/partner/año promedio
- **50 clientes de partnerships**

**Costo:** $0 directo (tiempo de outreach + comisiones si aplica)

---

#### **Canal 5: Referral Program**

**Estrategia:**
- Incentivar a clientes actuales a referir otras startups
- Reward: 1 mes gratis para quien refiere Y para referido
- Trackear con códigos únicos

**Mecánica:**
1. Cliente invita a colega (email o link único)
2. Referido se registra y convierte a pago
3. Ambos reciben 1 mes gratis (crédito en cuenta)

**Promoción:**
- In-app banner
- Email reminder mensual
- Testimonial requests include referral ask

**Volumen esperado:**
- Referral rate: 15% de clientes (optimista para product con buena experiencia)
- 65 clientes × 15% = **10 referrals/año**

**Costo:** $49 × 2 × 10 = $980 en créditos (pero retención alta)

---

### 6.2 Funnel de Conversión

```
🌐 Website Visitors: 30,000/año
         ↓ (10% → trial signup)
🧪 Trial Signups: 3,000/año
         ↓ (25% → paid)
💳 Paying Customers: 750/año
         ↓ (- churn 7%/mes promedio)
✅ Active Customers (end of year): 65
```

**Métricas clave:**
- **Visitor → Trial**: 10% (industry standard: 2-5%, tenemos ventaja por pricing)
- **Trial → Paid**: 25% (industry standard: 15-25%)
- **Churn**: 7% mensual promedio (alto para Año 1, mejora en Año 2)

**Optimizaciones para mejorar funnel:**
- Landing page A/B testing (headlines, CTAs)
- Onboarding wizard mejorado (time-to-value < 10 min)
- Email sequences automáticos durante trial
- Exit intent popups (ofrecer extensión de trial)
- Demo calls para trials que no activan en 3 días

---

## 7. Plan de Producto

### 7.1 Roadmap Features (priorizado)

**Q1 2026: MVP Core**
| Feature | Esfuerzo | Impacto | Prioridad |
|---------|----------|---------|-----------|
| Jobs CRUD (crear, editar, publicar) | 40h | 🔥🔥🔥 | P0 |
| Candidates database | 40h | 🔥🔥🔥 | P0 |
| Application pipeline (stages) | 60h | 🔥🔥🔥 | P0 |
| Email notifications (basic) | 30h | 🔥🔥 | P0 |
| Dashboard (overview) | 20h | 🔥🔥 | P1 |
| Team collaboration (comments) | 30h | 🔥🔥 | P1 |
| **Total Q1** | **220h** | | |

**Q2 2026: Collaboration & Efficiency**
| Feature | Esfuerzo | Impacto | Prioridad |
|---------|----------|---------|-----------|
| Interview scheduling | 50h | 🔥🔥🔥 | P0 |
| Chrome extension (LinkedIn import) | 40h | 🔥🔥🔥 | P0 |
| Candidate rating/scoring | 30h | 🔥🔥 | P1 |
| Bulk actions | 25h | 🔥🔥 | P1 |
| Reportes (pipeline analytics) | 35h | 🔥🔥 | P1 |
| **Total Q2** | **180h** | | |

**Q3 2026: Growth & Integrations**
| Feature | Esfuerzo | Impacto | Prioridad |
|---------|----------|---------|-----------|
| Career page (embeddable) | 50h | 🔥🔥🔥 | P0 |
| Email templates (customizable) | 30h | 🔥🔥 | P0 |
| GitHub OAuth integration | 40h | 🔥🔥🔥 | P0 |
| Slack notifications | 25h | 🔥🔥 | P1 |
| Mobile optimization | 40h | 🔥🔥 | P1 |
| **Total Q3** | **185h** | | |

**Q4 2026: Scale & Advanced Features**
| Feature | Esfuerzo | Impacto | Prioridad |
|---------|----------|---------|-----------|
| Zapier integration | 60h | 🔥🔥🔥 | P0 |
| Advanced analytics dashboard | 50h | 🔥🔥 | P0 |
| Candidate sourcing tools | 40h | 🔥🔥 | P1 |
| API pública (docs) | 45h | 🔥 | P2 |
| Marketplace foundation | 60h | 🔥🔥🔥 | P0 |
| **Total Q4** | **255h** | | |

**Total desarrollo Año 1: 840 horas (~21 semanas full-time)**

### 7.2 Decisiones de Producto

**Build vs Buy:**
- ✅ **Build**: Core ATS features (diferenciador)
- 🛒 **Buy/Integrate**: Scheduling (Calendly API), payments (Stripe), emails (SendGrid)

**Tech Stack:**
- Frontend: React + TypeScript + TailwindCSS
- Backend: Go (ya existe)
- Database: PostgreSQL (ya configurada)
- Hosting: AWS Lightsail → ECS

**Design Philosophy:**
- Simplicity > Features (menos es más)
- Speed (time-to-value < 10 minutos)
- Mobile-friendly desde día 1
- Accesibilidad (WCAG AA)

---

## 8. Equipo y Recursos

### 8.1 Estructura Año 1

**Founder/CEO (tú):**
- Desarrollo (60% tiempo)
- Sales (20% tiempo)
- Strategy/Operations (20% tiempo)

**Contractors/Freelancers:**
- Frontend developer (Q1-Q2): 200 horas @ $30/h = $6,000
- Content writer (Q2-Q4): $200/mes = $1,800/año
- Customer support (Q3-Q4): $300/mes = $1,800/año

**Total labor costs (beyond founder):** $9,600

**Founder compensation:**
- Año 1: $0 (ramen profitability, bootstrapped)
- Año 2: $3,000/mes una vez llegues a $10k MRR

### 8.2 Time Allocation (Founder)

| Actividad | Q1 | Q2 | Q3 | Q4 |
|-----------|----|----|----|----|
| **Desarrollo** | 80% | 60% | 50% | 40% |
| **Sales/Marketing** | 15% | 30% | 35% | 40% |
| **Operations** | 5% | 10% | 15% | 20% |

**Horas/semana:** 60 horas (startup mode)

### 8.3 Skills Necesarias

**Debe tener (tú o contratar):**
- ✅ Backend development (Go) - TÚ
- ✅ Frontend development (React) - CONTRATAR Q1
- ✅ DevOps básico (AWS, Docker) - TÚ
- ✅ Sales (outreach, demos, closing) - TÚ

**Nice to have:**
- 📊 Data analysis (SQL, analytics)
- 🎨 Design (Figma, UI/UX)
- 📝 Content writing (SEO, copywriting)
- 📢 Marketing (ads, growth hacking)

**Learning path:**
- Mes 1-2: Founder-led sales (aprender por hacer)
- Mes 3-4: Content marketing básico
- Mes 5-6: LinkedIn ads optimization
- Mes 7-12: Customer success & retention

---

## 9. Métricas y KPIs

### 9.1 North Star Metric

**MRR (Monthly Recurring Revenue)**

Por qué: Indicador único que captura crecimiento, retención y health del negocio.

**Target Año 1:** $5,000 MRR (Diciembre 2026)

### 9.2 KPIs por Categoría

**Growth Metrics:**
| Métrica | Q1 | Q2 | Q3 | Q4 | Meta Anual |
|---------|----|----|----|----|------------|
| Nuevas empresas/mes | 5 | 9 | 5 | 3 | 67 total |
| MRR growth rate | 250% | 100% | 40% | 20% | 5,000% |
| Total empresas activas | 15 | 41 | 57 | 65 | 65 |

**Product Metrics:**
| Métrica | Target Q4 |
|---------|-----------|
| Weekly Active Users (WAU) | 70% de usuarios |
| Average session time | 12 minutos |
| Feature adoption rate | 60% usan ≥3 features |
| Time-to-first-value | < 10 minutos |

**Sales & Marketing:**
| Métrica | Target Q4 |
|---------|-----------|
| Website visitors/mes | 10,000 |
| Trial signups/mes | 150 |
| Trial → Paid conversion | 25% |
| CAC | < $150 |
| Sales cycle length | 10-14 días |

**Customer Success:**
| Métrica | Target Q4 |
|---------|-----------|
| Monthly churn rate | < 6% |
| NPS | > 50 |
| Customer health score | > 75/100 |
| Support ticket resolution | < 24h |

**Financial:**
| Métrica | Target Año 1 |
|---------|--------------|
| ARR | $60,000 |
| Gross margin | > 90% |
| CAC payback | < 3 meses |
| LTV/CAC | > 6x |
| Burn rate | Break-even Q2 |

### 9.3 Dashboard de Seguimiento

**Herramientas:**
- Analytics: Google Analytics + Mixpanel
- CRM: HubSpot (free)
- Financial: Excel/Google Sheets (inicial) → Stripe Dashboard
- Product: Custom dashboard en Dvra admin

**Revisión:**
- **Diaria**: MRR, nuevos signups, churn events
- **Semanal**: Pipeline de ventas, demos booked, conversiones
- **Mensual**: P&L, cohort analysis, NPS
- **Trimestral**: Strategic review, roadmap ajustes

---

## 10. Riesgos y Mitigación

### 10.1 Riesgos Top 5

**Riesgo 1: No alcanzar product-market fit**
- **Probabilidad**: Media (40%)
- **Impacto**: Crítico
- **Señales**: Churn >15%, NPS <20, feedback negativo recurrente
- **Mitigación**:
  - Hablar con 50+ prospects antes de codear features grandes
  - Weekly user interviews (Q1-Q2)
  - Pivotar rápido si métricas no mejoran en 3 meses
  - Plan B: Nicho más específico (ej: solo agencias de developers)

**Riesgo 2: Competencia agresiva (Greenhouse, Lever, Gusto entran a LATAM)**
- **Probabilidad**: Media (30%)
- **Impacto**: Alto
- **Mitigación**:
  - Construir moat con marketplace (Año 2)
  - Brand en español-first muy fuerte
  - Pricing LATAM siempre 3-5x más barato
  - Community lock-in (red de usuarios)
  - Partnerships exclusivos con aceleradoras

**Riesgo 3: Crecimiento muy lento (no llegamos a 50 empresas)**
- **Probabilidad**: Media-Alta (50%)
- **Impacto**: Medio (no mata negocio, pero retrasa)
- **Señales**: No cumplir target Q1 (5 empresas)
- **Mitigación**:
  - Pricing experiment (bajar a $29/mes temporalmente)
  - Freemium tier (gratis hasta 5 candidatos)
  - Partnerships más agresivos
  - Contratar SDR part-time (si hay budget)
  - Content marketing más agresivo

**Riesgo 4: Churn alto (>12% mensual)**
- **Probabilidad**: Media (35%)
- **Impacto**: Crítico (imposibilita crecimiento)
- **Señales**: Churn >10% en Q1-Q2
- **Mitigación**:
  - Onboarding mejorado (white-glove primeros 50 clientes)
  - Feature usage tracking (alertas de low engagement)
  - Proactive customer success check-ins
  - Quarterly business reviews (QBRs)
  - Annual plan incentives (2 meses gratis)

**Riesgo 5: Problemas técnicos/downtime**
- **Probabilidad**: Media (40%)
- **Impacto**: Alto (pérdida de confianza)
- **Mitigación**:
  - Monitoring desde día 1 (Sentry, CloudWatch)
  - Uptime target: 99.5% (4 horas downtime/mes permitido)
  - DB backups daily
  - Incident response playbook
  - Status page público (transparencia)

### 10.2 Contingency Plans

**Si no llegamos a $60k ARR:**
- Extender runway con freelancing (founder)
- Reducir marketing spend 50%
- Focus en retention vs acquisition
- Evaluar raising pre-seed ($100k) si hay tracción (>20 empresas)

**Si logramos superar targets:**
- Acelerar roadmap marketplace
- Contratar 1 SDR full-time
- Expandir a Brasil (Portuguese version)
- Considerar YC/500 Startups application

---

## Conclusión: 3 Keys to Success Año 1

### 1️⃣ **Execution Speed**
- Ship fast, iterate faster
- No perfectionism: 80% es suficiente para lanzar
- Weekly releases (viernes deploy)

### 2️⃣ **Customer Obsession**
- Hablar con users 5+ veces/semana
- Support response < 4 horas
- Feature decisions driven por feedback real, no assumptions

### 3️⃣ **Focus Brutal**
- Decir NO a distracciones (conferencias, partnerships no-core, features vanity)
- Priorizar ruthlessly: MRR es el único objetivo que importa
- No marketplace hasta tener 50 empresas pagando (disciplina)

---

**"Build something people want, charge for it, and don't run out of money."**

¡Vamos con todo! 🚀

---

**FIN BUSINESS PLAN AÑO 1**

> Versión 1.0 | Diciembre 8, 2025
