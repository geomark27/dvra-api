# Flujo landlord (SuperAdmin)

> Documento operativo vigente desde 2026-07-15. La decisión arquitectónica se
> detalla en [ADR-002](adr/ADR-002-contextos-landlord-tenant.md).

## Modelo de acceso

El SuperAdmin es un operador del SaaS, no un integrante de una empresa cliente.
Por eso su acceso se almacena en `platform_memberships`:

| Dato | Valor |
|---|---|
| Identidad | Registro global en `users` |
| Perfil landlord | `platform_memberships` con rol `superadmin` |
| Empresa | No aplica; no existe `company_id` en el perfil landlord |
| Contexto JWT | `context_type: "platform"` |
| Permisos | Matriz `platform.*` independiente de los roles tenant |

Una misma identidad puede tener un perfil landlord y, además, memberships en
empresas. Cada token representa exactamente uno de esos contextos; nunca mezcla
privilegios landlord y tenant.

## Aprovisionamiento

Primero debe existir la identidad en `users`. El acceso inicial se concede con:

```bash
go run ./cmd/console platform-admin --email admin@example.com
```

El comando es idempotente: crea o reactiva el perfil `superadmin`. No crea una
membership tenant ni asigna una empresa ficticia.

## Autenticación y contextos

1. `POST /api/v1/auth/login` selecciona el perfil platform activo como contexto
   inicial cuando el usuario dispone de él.
2. `POST /api/v1/auth/switch-company` valida una membership activa y emite un
   access token y un refresh token tenant nuevos.
3. `POST /api/v1/auth/switch-platform` valida el perfil landlord activo y emite
   un access token y un refresh token platform nuevos.
4. `POST /api/v1/auth/refresh` revalida en base de datos el perfil del contexto
   actual. Un perfil suspendido o eliminado no puede renovar tokens.
5. `GET /api/v1/auth/me` devuelve `active_context`, accesos disponibles y los
   permisos efectivos del contexto actual.

Ejemplo conceptual de claims platform:

```json
{
  "user_id": 1,
  "context_type": "platform",
  "membership_id": 4,
  "role": "superadmin"
}
```

Un token platform no contiene `company_id` y es rechazado por las rutas tenant.
Un token tenant es rechazado por las rutas `/api/v1/platform/*`.

## API landlord

Todas estas rutas requieren `Authorization: Bearer <access_token>` con contexto
`platform` y el permiso `platform.*` correspondiente:

| Recurso | Rutas |
|---|---|
| Empresas | `GET/POST /api/v1/platform/companies`, `GET/PUT/DELETE /api/v1/platform/companies/:id` |
| Identidades | `GET /api/v1/platform/users`, `GET /api/v1/platform/users/:id` |
| Accesos tenant | CRUD en `/api/v1/platform/memberships` |
| Planes | CRUD, toggle y assign en `/api/v1/platform/plans` |
| Configuración | `GET/PUT /api/v1/platform/settings` |

La especificación exacta de request y response se publica en
`/swagger/index.html` al ejecutar la API.

## Límites de seguridad

- El rol `superadmin` está reservado y no puede asignarse a `memberships`.
- Toda membership tenant requiere `company_id` y solo puede existir una activa
  por `(user_id, company_id)`.
- Los handlers tenant no conceden acceso global implícito al SuperAdmin.
- Los roles tenant (`admin`, `recruiter`, `hiring_manager`, `interviewer`) se
  asignan mediante `memberships.tenant_role_id`; sus grants se resuelven en
  `tenant_role_permissions`, separados de los permisos landlord.
- `hiring_manager` solo ve los jobs asignados y candidatos/aplicaciones ligados
  a esos jobs.
- Eliminar un usuario desde un tenant elimina su membership en esa empresa; no
  elimina la identidad global ni sus accesos a otros tenants.

## Migración desde el esquema anterior

`go run ./cmd/console migrate` transforma las memberships legacy con
`company_id IS NULL` y rol `superadmin` en perfiles `platform_memberships`.
Luego exige `memberships.company_id NOT NULL`, migra el rol textual hacia
`tenant_role_id`, elimina la columna antigua y crea el índice único parcial de
membresías activas. La migración se detiene ante roles globales desconocidos o
duplicados para evitar una conversión silenciosa incorrecta.

Los JWT anteriores a esta separación no tienen `context_type`, por lo que dejan
de ser válidos. Después del despliegue todos los usuarios deben iniciar sesión
nuevamente.
