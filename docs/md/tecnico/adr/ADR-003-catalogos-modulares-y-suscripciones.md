# ADR-003: catálogos modulares y suscripciones

- **Estado:** Aceptado
- **Fecha:** 2026-07-20
- **Tipo:** Decisión arquitectónica
- **Reemplaza:** Ninguno
- **Reemplazado por:** Ninguno

## Decisión

DVRA usa un catálogo relacional único, organizado como:

`system_value_modules` → `system_value_categories` → `system_values`.

Un valor con `company_id IS NULL` es global. Un valor con `company_id` pertenece exclusivamente a ese tenant. El API nunca recibe `company_id` para operaciones tenant: lo toma del JWT activo.

Las categorías declaran si permiten valores tenant y si una empresa puede personalizar la etiqueta, orden o visibilidad de un valor global. Esas personalizaciones viven en `company_system_value_overrides`; el registro global no se modifica.

Los planes se relacionan con sus catálogos mediante IDs. La asociación empresa-plan vive en `company_subscriptions.plan_id`, conserva historial y snapshots de precio, moneda y ciclo. `companies.plan_tier` deja de ser persistente. La configuración global usa `platform_settings.default_plan_id` y cada tenant posee una fila en `company_settings`.

## Módulos y comandos

Los módulos registrados son `core`, `identity`, `recruitment`, `staffing` y `billing`. Los datos demo son un módulo explícito separado.

```bash
go run ./cmd/console migrate --module recruitment
go run ./cmd/console seed --module recruitment
go run ./cmd/console migrate --module billing --seed
go run ./cmd/console seed --module demo
```

`--fresh` no puede combinarse con `--module`, porque eliminar parcialmente el esquema invalidaría el historial de migraciones de los demás módulos.

## Contrato tenant

- `GET /api/v1/system-values?module=recruitment`
- `GET /api/v1/system-values/modules/recruitment/application_stage`
- `POST /api/v1/system-values` para categorías extensibles.
- `PUT /api/v1/system-values/:id/override` para personalizaciones permitidas.
- `GET|PUT /api/v1/settings` para personalización de la empresa.
- `GET /api/v1/subscription` para el plan y estado actuales.
- `GET /api/v1/subscription/usage` para consumo de límites.

La ruta antigua `GET /api/v1/system-values/:category` se mantiene temporalmente, pero no debe usarse para nuevas integraciones porque no identifica el módulo.

## Contrato landlord

- CRUD global: `/api/v1/platform/system-values`.
- CRUD de planes: `/api/v1/platform/plans` usando `currency_value_id`, `billing_cycle_value_id` y `support_level_value_id`.
- Asignación: `POST /api/v1/platform/plans/assign` con `company_id` y `plan_id`.

Solo el landlord asigna planes. Un tenant no puede modificar plan ni fechas de prueba mediante la actualización de empresa.
