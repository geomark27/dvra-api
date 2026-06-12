# Documentación Dvra API

> Índice general de la documentación del proyecto.
> Última reorganización: Junio 2026.

## 📂 Estructura

```
docs/md/
├── README.md                          ← Este índice
├── negocio/                           ← Modelo, lógica y plan de negocio
│   ├── LOGICA_DE_NEGOCIO.md
│   └── PLAN_DE_NEGOCIO_ANO1.md
├── tecnico/                           ← Detalles técnicos y desarrollo
│   ├── ARQUITECTURA.md
│   ├── API_ENDPOINTS.md
│   ├── FLUJO_CLIENTE.md
│   ├── FLUJO_SUPERADMIN.md
│   ├── MODULO_PLANES.md
│   └── MODULO_UBICACIONES.md
└── bitacora/                          ← Registro cronológico de actividades
    └── BITACORA.md
```

## 💼 Negocio

| Documento | Contenido | Úsalo para |
|-----------|-----------|------------|
| [LOGICA_DE_NEGOCIO.md](negocio/LOGICA_DE_NEGOCIO.md) | Visión del producto, modelo dual (ATS + Marketplace Red Dvra), reglas de negocio (RN-*), roles y permisos, pipeline de candidatos, pricing por tier, roadmap | Entender **qué** hace el producto y **por qué** |
| [PLAN_DE_NEGOCIO_ANO1.md](negocio/PLAN_DE_NEGOCIO_ANO1.md) | Plan operativo y financiero del primer año: timeline mensual, go-to-market, modelo financiero, KPIs, riesgos | Decisiones de priorización y metas comerciales |

## 🔧 Técnico

| Documento | Contenido | Úsalo para |
|-----------|-----------|------------|
| [ARQUITECTURA.md](tecnico/ARQUITECTURA.md) | Tech stack (Go/Gin/GORM/PostgreSQL), arquitectura de capas, multi-tenancy, modelos de datos, seguridad JWT, infraestructura AWS | Entender la estructura del sistema antes de tocar código |
| [API_ENDPOINTS.md](tecnico/API_ENDPOINTS.md) | **Referencia completa de todos los endpoints** con requests/responses de ejemplo | Consulta diaria durante el desarrollo e integración frontend |
| [FLUJO_CLIENTE.md](tecnico/FLUJO_CLIENTE.md) | Flujo completo del cliente (admin de empresa): registro, multi-empresa, permisos por rol, aislamiento multi-tenant, casos de uso | Entender el recorrido y permisos de un usuario de empresa |
| [FLUJO_SUPERADMIN.md](tecnico/FLUJO_SUPERADMIN.md) | Flujo completo del SuperAdmin: autenticación global, gestión de empresas/planes/memberships, analytics, restricciones | Entender la operación administrativa de la plataforma |
| [MODULO_PLANES.md](tecnico/MODULO_PLANES.md) | Módulo de planes de suscripción: arquitectura, endpoints, modelo de datos, validaciones, seeder | Trabajar con planes, límites y asignación a empresas |
| [MODULO_UBICACIONES.md](tecnico/MODULO_UBICACIONES.md) | Módulo geográfico (Region→Subregion→Country→State→City): endpoints públicos, CRUD admin, seeders, integración frontend | Trabajar con datos de ubicación |

## 📓 Bitácora

| Documento | Contenido |
|-----------|-----------|
| [BITACORA.md](bitacora/BITACORA.md) | Registro cronológico de implementaciones y auditorías: multi-empresa, auditoría de seguridad multi-tenancy, separación SuperAdmin, módulo de planes, módulo de ubicaciones. Incluye pendientes de cada hito. |

## ⚠️ Notas importantes

- **Precios:** la fuente de verdad de los planes actuales es el seeder (`internal/database/seeders/plan_seeder.go`): Free $0, Starter $39.99, Professional $79.99, Enterprise $159.99. Los precios en `negocio/` ($49/$149/$399) son el **objetivo comercial** del plan de negocio, aún no reflejado en el código.
- **Credenciales por defecto del SuperAdmin:** deben cambiarse en producción (ver bitácora).
- Al completar un hito relevante, agregar una entrada al inicio de `bitacora/BITACORA.md` con el formato indicado al final de ese archivo.
