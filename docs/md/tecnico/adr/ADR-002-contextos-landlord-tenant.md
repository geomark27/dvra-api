# ADR-002: Separación de contextos landlord y tenant

- **Estado:** Aceptado
- **Fecha:** 2026-07-15
- **Tipo:** Decisión arquitectónica
- **Reemplaza:** Ninguno
- **Reemplazado por:** Ninguno

## Decisión

La identidad de una persona vive en `users`. El acceso landlord vive en
`platform_memberships`; el acceso a clientes vive exclusivamente en
`memberships`, donde `company_id` es obligatorio y `(user_id, company_id)` es
único para membresías activas.

Los JWT declaran un `context_type` explícito:

- `platform`: incluye `membership_id` y un rol de plataforma, sin `company_id`.
- `tenant`: incluye `membership_id`, `company_id` y `tenant_role_id`; no lleva
  el nombre textual del rol como autoridad.
- `identity`: usuario autenticado todavía sin contexto autorizado.

Los permisos landlord (`platform.*`) y tenant (`jobs.*`, `candidates.*`, etc.)
se resuelven con matrices distintas. Un token platform no puede consumir rutas
tenant y un token tenant no puede consumir `/api/v1/platform/*`.

## Rutas y cambio de contexto

- `POST /api/v1/auth/switch-company` emite access y refresh token del tenant.
- `POST /api/v1/auth/switch-platform` vuelve al contexto landlord si existe una
  `platform_membership` activa.
- `GET /api/v1/auth/me` devuelve contexto activo, rol y permisos efectivos.
- La administración landlord vive bajo `/api/v1/platform/*`.

## Roles tenant

Los perfiles incorporados son `admin`, `recruiter`, `hiring_manager` e
`interviewer`. Se persisten en `tenant_roles`; las acciones atómicas viven en
`permissions` y la relación efectiva en `tenant_role_permissions`.
`memberships.tenant_role_id` es una FK obligatoria: las APIs de usuarios y
membresías reciben este ID y nunca un nombre de rol libre.

Las claves y grants iniciales se versionan con el release y se sincronizan por
migración. En runtime, el middleware valida la membership activa, toma su rol
vigente desde BD y consulta los permisos relacionales. Por ello, suspender una
membership o cambiar su rol tiene efecto inmediato aun cuando exista un access
token anterior.

`admin` recibe siempre el catálogo tenant completo, incluidas las acciones que
se incorporen en módulos posteriores. Esta regla no cruza el límite landlord:
las rutas `/api/v1/platform/*` continúan exigiendo contexto y permisos platform.

`hiring_manager` se restringe en repositorio a jobs asignados y a los candidatos
y postulaciones derivados de esos jobs. `interviewer` conserva solo lectura
básica y no recibe acceso a candidatos o postulaciones mientras no exista un
modelo explícito de asignación de entrevistas y scorecards.

## Operación

Después de `make db-migrate`, el primer acceso landlord puede otorgarse a una
identidad existente con:

```bash
go run ./cmd/console platform-admin --email admin@example.com
```

La migración convierte memberships legacy con `company_id IS NULL` y rol
`superadmin` en `platform_memberships`. También convierte `user` a
`interviewer`, rellena `tenant_role_id`, crea sus FKs y elimina la columna
textual `memberships.role`. Los tokens emitidos antes de esta separación dejan
de ser válidos y los usuarios deben iniciar sesión nuevamente.
