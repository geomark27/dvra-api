# ADR-002: Separación de contextos landlord y tenant

**Estado:** Aceptado  
**Fecha:** 2026-07-15

## Decisión

La identidad de una persona vive en `users`. El acceso landlord vive en
`platform_memberships`; el acceso a clientes vive exclusivamente en
`memberships`, donde `company_id` es obligatorio y `(user_id, company_id)` es
único para membresías activas.

Los JWT declaran un `context_type` explícito:

- `platform`: incluye `membership_id` y un rol de plataforma, sin `company_id`.
- `tenant`: incluye `membership_id`, `company_id` y el rol dentro de ese tenant.
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

Los perfiles `admin`, `recruiter`, `hiring_manager` y `user` son roles de
producto definidos en código. No se mantiene un catálogo `roles` parcialmente
conectado. Si se habilitan roles personalizados, se introducirán
`tenant_roles`, `permissions` y `role_permissions` como una migración completa,
sin mezclar esa configuración con roles landlord.

`hiring_manager` se restringe en repositorio a jobs asignados y a los candidatos
y postulaciones derivados de esos jobs. `user` no recibe acceso a candidatos o
postulaciones mientras no exista un modelo explícito de asignación para ese rol.

## Operación

Después de `make db-migrate`, el primer acceso landlord puede otorgarse a una
identidad existente con:

```bash
go run ./cmd/console platform-admin --email admin@example.com
```

La migración convierte memberships legacy con `company_id IS NULL` y rol
`superadmin` en `platform_memberships`. Los tokens emitidos antes de esta
separación dejan de ser válidos y los usuarios deben iniciar sesión nuevamente.
