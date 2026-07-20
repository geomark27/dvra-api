# Ajustes frontend: catálogos, suscripciones y settings

- **Fecha:** 2026-07-20
- **Backend:** Implementado
- **Frontend:** Pendiente de aplicar
- **Referencia arquitectónica:** [ADR-003](../tecnico/adr/ADR-003-catalogos-modulares-y-suscripciones.md)

## Objetivo

Adaptar el frontend al catálogo modular, la relación empresa → suscripción →
plan y la configuración personalizada de cada tenant. No se debe adelantar
todavía la migración de estados operativos de recruitment y staffing.

## 1. Cliente de catálogos tenant

Agregar los siguientes métodos:

```http
GET /api/v1/system-values?module={module}
GET /api/v1/system-values/modules/{module}/{category}
```

Ejemplos:

```http
GET /api/v1/system-values?module=recruitment
GET /api/v1/system-values/modules/recruitment/application_stage
GET /api/v1/system-values/modules/billing/currency
GET /api/v1/system-values/modules/billing/billing_cycle
GET /api/v1/system-values/modules/billing/support_level
```

La ruta anterior sigue disponible temporalmente:

```http
GET /api/v1/system-values/{category}
```

Las implementaciones nuevas deben usar la ruta con `modules`.

Query keys recomendadas:

```ts
["tenant", "system-values", "categories", module]
["tenant", "system-values", module, category]
```

Tipo base:

```ts
interface SystemValue {
  id: number;
  category_id: number;
  company_id?: number;
  value: string;
  label: string;
  description?: string;
  display_order: number;
  is_system: boolean;
  is_active: boolean;
}
```

El identificador técnico se devuelve como `value`, aunque internamente el
backend lo persiste como `key`.

## 2. Administración de catálogos tenant

Permisos:

```text
catalogs.view
catalogs.manage
```

Crear un valor propio de la empresa:

```http
POST /api/v1/system-values
```

```json
{
  "category_id": 12,
  "value": "headhunter",
  "label": "Headhunter",
  "description": "Fuente personalizada",
  "display_order": 20
}
```

No enviar `company_id`; el backend lo obtiene del JWT.

Editar un valor creado por el tenant:

```http
PUT /api/v1/system-values/{id}
```

```json
{
  "label": "Nueva etiqueta",
  "description": "Descripción actualizada",
  "display_order": 10,
  "is_active": true
}
```

Personalizar un valor global:

```http
PUT /api/v1/system-values/{id}/override
```

```json
{
  "custom_label": "Nombre personalizado",
  "display_order": 5,
  "is_enabled": true
}
```

Después de crear, editar, eliminar o personalizar:

```ts
queryClient.invalidateQueries({
  queryKey: ["tenant", "system-values"],
});
```

Manejar `403` cuando una categoría no permita valores tenant u overrides.

## 3. Configuración personalizada de empresa

Crear una vista tenant de configuración utilizando:

```http
GET /api/v1/settings
PUT /api/v1/settings
```

Permisos:

```text
settings.view
settings.update
```

```ts
interface CompanySettings {
  id: number;
  company_id: number;
  display_name?: string;
  logo_url?: string;
  primary_color: string;
  secondary_color: string;
  locale: string;
  career_page_title?: string;
  career_page_description?: string;
  career_page_enabled: boolean;
  custom_domain?: string;
  reply_to_email?: string;
}
```

Ejemplo de actualización parcial:

```json
{
  "display_name": "Azentic Talent",
  "logo_url": "https://example.com/logo.png",
  "primary_color": "#2563eb",
  "secondary_color": "#0f172a",
  "locale": "es",
  "career_page_title": "Trabaja con nosotros",
  "career_page_description": "Conoce nuestras vacantes",
  "career_page_enabled": true,
  "custom_domain": "empleos.empresa.com",
  "reply_to_email": "talento@empresa.com"
}
```

Query e invalidación:

```ts
const settingsKey = ["tenant", "settings"];

queryClient.invalidateQueries({ queryKey: settingsKey });
```

Los colores deben enviarse en formato hexadecimal, por ejemplo `#2563eb`.

## 4. Plan y suscripción del tenant

Agregar la sección **Plan y uso** dentro de configuración:

```http
GET /api/v1/subscription
GET /api/v1/subscription/usage
```

Permiso:

```text
billing.view
```

Query keys:

```ts
["tenant", "subscription"]
["tenant", "subscription", "usage"]
```

```ts
interface CompanySubscription {
  id: number;
  company_id: number;
  plan_id: number;
  status_value_id: number;
  currency_value_id: number;
  billing_cycle_value_id: number;
  price_snapshot: number;
  started_at: string;
  current_period_start: string;
  current_period_end?: string;
  trial_ends_at?: string;
  ended_at?: string;
  cancel_at_period_end: boolean;
  plan?: Plan;
  status?: SystemValue;
  currency?: SystemValue;
  billing_cycle?: SystemValue;
}

interface SubscriptionUsage {
  users: number;
  jobs: number;
  candidates: number;
  applications: number;
  storage_bytes: number;
}
```

Mostrar barras de consumo comparando el uso contra:

```text
plan.max_users
plan.max_jobs
plan.max_candidates
plan.max_applications
plan.max_storage_gb
```

Cuando un límite sea `-1`, mostrar **Ilimitado**. El tenant no puede cambiar
su plan directamente desde esta vista. Puede mostrarse un CTA informativo para
solicitar una mejora.

## 5. Formulario landlord de planes

Los formularios ya no deben enviar textos para moneda, ciclo o soporte.

Contrato anterior:

```json
{
  "currency": "USD",
  "billing_cycle": "monthly",
  "support_level": "priority"
}
```

Contrato vigente:

```json
{
  "currency_value_id": 10,
  "billing_cycle_value_id": 12,
  "support_level_value_id": 15
}
```

En contexto landlord, cargar las opciones desde:

```http
GET /api/v1/platform/system-values
```

No usar las rutas tenant de catálogos con un token platform. Filtrar por:

```ts
value.category.module.key === "billing"
value.category.key === "currency"
value.category.key === "billing_cycle"
value.category.key === "support_level"
```

```ts
interface Plan {
  id: number;
  name: string;
  slug: string;
  description: string;
  price: number;

  currency: string;
  currency_value_id: number;
  billing_cycle: string;
  billing_cycle_value_id: number;
  support_level: string;
  support_level_value_id: number;

  is_active: boolean;
  is_public: boolean;
  trial_days: number;
  display_order: number;

  max_users: number;
  max_jobs: number;
  max_candidates: number;
  max_applications: number;
  max_storage_gb: number;

  can_export_data: boolean;
  can_use_custom_brand: boolean;
  can_use_api: boolean;
  can_use_integrations: boolean;
  can_use_staffing: boolean;
}
```

## 6. Asignación de plan desde landlord

```http
POST /api/v1/platform/plans/assign
```

```json
{
  "company_id": 4,
  "plan_id": 2
}
```

No asociar el plan modificando directamente `company.plan_tier`.

Después de asignar:

```ts
queryClient.invalidateQueries({ queryKey: ["platform", "companies"] });
queryClient.invalidateQueries({ queryKey: ["platform", "plans"] });
```

Cuando se cambie posteriormente al contexto tenant, limpiar también:

```ts
["tenant", "subscription"]
["tenant", "subscription", "usage"]
```

## 7. Configuración global del landlord

En platform settings, reemplazar el campo editable:

```text
default_plan_tier
```

por:

```text
default_plan_id
```

Ejemplo:

```json
{
  "default_trial_days": 14,
  "default_plan_id": 2
}
```

La respuesta conserva `default_plan_tier` como dato informativo, pero las
actualizaciones deben enviar `default_plan_id`.

## 8. Actualización de empresa tenant

En `PUT /api/v1/companies/{id}`, el tenant no debe enviar ni mostrar como
editables:

```text
plan_tier
trial_ends_at
```

El backend los ignora en contexto tenant. En la creación landlord de una
empresa, `plan_tier` continúa temporalmente como slug por compatibilidad. Los
cambios posteriores de plan deben usar `/api/v1/platform/plans/assign`.

## 9. Guards y menú

Incorporar los permisos entregados por `/api/v1/auth/me`:

```text
catalogs.view
catalogs.manage
settings.view
settings.update
billing.view
platform.catalogs.manage
```

Menú tenant recomendado:

```text
Configuración
├── Empresa
├── Apariencia
├── Career page
├── Catálogos
└── Plan y uso
```

Reglas:

- **Catálogos:** requiere `catalogs.view`.
- Editar catálogos requiere `catalogs.manage`.
- **Configuración:** requiere `settings.view`.
- Guardar configuración requiere `settings.update`.
- **Plan y uso:** requiere `billing.view`.
- El administrador tenant recibe estos permisos automáticamente.

## 10. Cambio de empresa y caché

Después de cambiar el tenant activo, limpiar como mínimo:

```ts
queryClient.removeQueries({ queryKey: ["tenant"] });
```

Nunca reutilizar settings, catálogos, suscripción o usage del tenant anterior.

## 11. Compatibilidad temporal

Todavía no cambiar a IDs los DTO de estados operativos de vacantes,
aplicaciones, candidatos o staffing. Por ahora continúan como claves de texto:

```text
job.status
job.location_type
application.stage
candidate.source
staffingClient.status
placement.status
```

La migración de esos campos se entregará como un cambio posterior del backend.

## Checklist de aceptación

- [ ] Los catálogos nuevos usan la ruta `/system-values/modules/...`.
- [ ] Ninguna petición tenant envía `company_id` para catálogos o settings.
- [ ] El administrador tenant puede abrir Configuración, Catálogos y Plan y uso.
- [ ] Cambiar de empresa limpia toda la caché tenant.
- [ ] El formulario de planes envía IDs de moneda, ciclo y soporte.
- [ ] El landlord asigna planes usando `company_id` y `plan_id`.
- [ ] El tenant no puede editar su plan ni las fechas de trial.
- [ ] Platform settings guarda `default_plan_id`.
- [ ] Los errores `401`, `403` y `409` muestran mensajes apropiados.
- [ ] ESLint, TypeScript y build final pasan correctamente.
