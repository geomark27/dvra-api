# Documentación Dvra API

> Índice general de la documentación del proyecto.
> Última reorganización: Julio 2026.

## 📂 Estructura

```
docs/md/
├── README.md                          ← Este índice
├── negocio/                           ← Modelo, lógica y plan de negocio
│   ├── LOGICA_DE_NEGOCIO.md
│   └── PLAN_DE_NEGOCIO_ANO1.md
├── tecnico/                           ← Detalles técnicos y desarrollo
│   ├── ESTANDARES_DESARROLLO.md       ← FUENTE DE VERDAD del cómo se programa
│   ├── AUDITORIA_ESTANDARES.md
│   ├── ARQUITECTURA.md
│   ├── adr/                           ← Registro histórico de decisiones
│   │   ├── README.md
│   │   ├── ADR-001-arquitectura-modular.md
│   │   ├── ADR-002-contextos-landlord-tenant.md
│   │   └── ADR-003-catalogos-modulares-y-suscripciones.md
│   ├── API_ENDPOINTS.md
│   ├── FLUJO_CLIENTE.md
│   ├── FLUJO_SUPERADMIN.md
│   ├── MODULO_PLANES.md
│   └── MODULO_UBICACIONES.md
├── historias/                         ← Historias Técnicas (HT) de cada tarea
│   └── README.md
├── frontend/                          ← Contratos e indicaciones para el frontend
│   ├── README.md
│   └── AJUSTES_CATALOGOS_SUSCRIPCIONES_SETTINGS.md
└── bitacora/                          ← Registro cronológico de actividades
    └── BITACORA.md
```

> El resumen operativo para agentes y desarrolladores está en **`AGENTS.md`** (raíz del repo):
> anatomía de módulo y las 7 reglas de modularidad.

## 💼 Negocio

| Documento | Contenido | Úsalo para |
|-----------|-----------|------------|
| [LOGICA_DE_NEGOCIO.md](negocio/LOGICA_DE_NEGOCIO.md) | Visión del producto, modelo dual (ATS + Marketplace Red Dvra), reglas de negocio (RN-*), roles y permisos, pipeline de candidatos, pricing por tier, roadmap | Entender **qué** hace el producto y **por qué** |
| [PLAN_DE_NEGOCIO_ANO1.md](negocio/PLAN_DE_NEGOCIO_ANO1.md) | Plan operativo y financiero del primer año: timeline mensual, go-to-market, modelo financiero, KPIs, riesgos | Decisiones de priorización y metas comerciales |

## 🔧 Técnico

| Documento | Contenido | Úsalo para |
|-----------|-----------|------------|
| [ESTANDARES_DESARROLLO.md](tecnico/ESTANDARES_DESARROLLO.md) | **Fuente de verdad**: cómo se escribe, estructura, verifica y documenta el código. Arquitectura por dominios, mapa de migración (§3.7), envelope HTTP (§2.6), multi-tenancy (§4), Definición de Done (§9) | **Antes de escribir cualquier línea.** Las skills `hu-a-ht` y `validar-estandares` lo leen como fuente de verdad |
| [AUDITORIA_ESTANDARES.md](tecnico/AUDITORIA_ESTANDARES.md) | Brechas del código frente al estándar, con severidad y evidencia | Priorizar deuda técnica |
| [ARQUITECTURA.md](tecnico/ARQUITECTURA.md) | Tech stack (Go/Gin/GORM/PostgreSQL), multi-tenancy, modelos de datos, seguridad JWT, infraestructura AWS | Entender la estructura del sistema antes de tocar código |
| [Índice de ADR](tecnico/adr/README.md) | Decisiones arquitectónicas aceptadas, reemplazadas u obsoletas | Entender por qué existe la arquitectura actual |
| [API_ENDPOINTS.md](tecnico/API_ENDPOINTS.md) | **Referencia completa de todos los endpoints** con requests/responses de ejemplo | Consulta diaria durante el desarrollo e integración frontend |
| [FLUJO_CLIENTE.md](tecnico/FLUJO_CLIENTE.md) | Flujo completo del cliente (admin de empresa): registro, multi-empresa, permisos por rol, aislamiento multi-tenant, casos de uso | Entender el recorrido y permisos de un usuario de empresa |
| [FLUJO_SUPERADMIN.md](tecnico/FLUJO_SUPERADMIN.md) | Flujo completo del SuperAdmin: autenticación global, gestión de empresas/planes/memberships, analytics, restricciones | Entender la operación administrativa de la plataforma |
| [MODULO_PLANES.md](tecnico/MODULO_PLANES.md) | Módulo de planes de suscripción: arquitectura, endpoints, modelo de datos, validaciones, seeder | Trabajar con planes, límites y asignación a empresas |
| [MODULO_UBICACIONES.md](tecnico/MODULO_UBICACIONES.md) | Módulo geográfico (Region→Subregion→Country→State→City): endpoints públicos, CRUD admin, seeders, integración frontend | Trabajar con datos de ubicación |

## 🖥️ Frontend

| Documento | Contenido | Úsalo para |
|-----------|-----------|------------|
| [Índice frontend](frontend/README.md) | Índice de contratos e indicaciones pendientes para el cliente web | Localizar los cambios que debe implementar el frontend |
| [AJUSTES_CATALOGOS_SUSCRIPCIONES_SETTINGS.md](frontend/AJUSTES_CATALOGOS_SUSCRIPCIONES_SETTINGS.md) | Endpoints, permisos, tipos, caché y checklist para catálogos modulares, configuración de empresa, planes y suscripciones | Entregar al equipo frontend como guía de implementación |

## 📓 Bitácora

| Documento | Contenido |
|-----------|-----------|
| [BITACORA.md](bitacora/BITACORA.md) | Registro cronológico de implementaciones y auditorías: multi-empresa, auditoría de seguridad multi-tenancy, separación SuperAdmin, módulo de planes, módulo de ubicaciones. Incluye pendientes de cada hito. |

## ⚠️ Notas importantes

- **Precios:** la fuente de verdad de lo que está **desplegado** es el seeder (`internal/database/seeders/plan_seeder.go`): Free $0, Starter $39.99, Professional $79.99, Enterprise $159.99.

  La tabla **aprobada** (2026-08-26), aún pendiente de implementar, es: Free $0 · Starter $49 · Growth $99 · **Agency $149** · Enterprise a convenir (`is_public: false`). Precio **por empresa, no por asiento**; `can_use_staffing` pasa de Professional a Agency. Los precios de `negocio/` ($49/$149/$399) quedan obsoletos.

  ⚠️ Al implementarla: `plan_seeder.go` **omite los planes que ya existen**, así que sobre una BD poblada el seeder no cambia nada. Hace falta una migración de datos en `internal/database/migrations/data_migrations.go`. El `PriceSnapshot` de las suscripciones vivas conserva el precio anterior, así que el grandfathering es automático.
- **Credenciales por defecto del SuperAdmin:** deben cambiarse en producción (ver bitácora).
- Al completar un hito relevante, agregar una entrada al inicio de `bitacora/BITACORA.md` con el formato indicado al final de ese archivo.
