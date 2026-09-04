# AGENTS.md — Dvra API

> Contrato de trabajo para cualquier agente o desarrollador que toque este repositorio.
> La fuente de verdad detallada es `docs/md/tecnico/ESTANDARES_DESARROLLO.md`. Este archivo es el
> resumen operativo: **qué forma tiene un módulo, qué está permitido y qué está prohibido.**

---

## 1. Qué es Dvra

ATS SaaS multi-tenant para empresas de tecnología en LATAM, con módulo de staffing/outsourcing
(gestionar vacantes propias y de clientes finales en la misma herramienta).

Stack: **Go 1.26** · Gin · GORM · PostgreSQL · swag (Swagger) · JWT.

`company_id` es la **única** frontera de aislamiento entre empresas. Cualquier otra dimensión
(`staffing_client_id`, por ejemplo) es un filtro interno, no una frontera.

---

## 2. Arquitectura macro

Dos mitades, con la dependencia en un solo sentido:

```
internal/
├── modules/<dominio>/     ← NEGOCIO. Consume la plataforma.
├── platform/              ← Transversal + landlord. NUNCA importa modules.
└── shared/                ← Tipos y utilidades sin lógica de dominio.
```

**Dominios:** `iam`, `recruitment`, `staffing`, `billing`, `catalog`, `locations`,
`platformsettings`, `dashboard`, `publicapi`. Más `platform/health`, que es infraestructura.

---

## 3. Anatomía de un módulo

```
internal/modules/<dominio>/
├── module.go        # Wiring del módulo + RegisterRoutes. Solo lo toca el composition root.
├── ports.go         # Interfaces que ESTE módulo necesita de otros (definidas por el consumidor).
├── transport/       # Handlers HTTP, DTOs de request/response, rutas.
├── service/         # Lógica de negocio y validaciones.
└── repository/      # Acceso a datos con GORM.
```

Qué va en cada uno:

| Paquete | Contiene | Nunca contiene |
|---|---|---|
| `transport/` | Bind de entrada, scoping, llamada al service, mapeo a DTO, anotaciones Swagger | Lógica de negocio, GORM, `database.DB` |
| `service/` | Reglas de negocio, validaciones, transacciones, integridad cruzada | `gin.Context`, HTTP, GORM |
| `repository/` | Queries, filtros, joins | Reglas de negocio, HTTP |
| `ports.go` | Interfaces mínimas hacia otros módulos | Implementaciones |

**Dirección interna:** `transport` importa `service`, nunca al revés. Por eso los tipos que
aparecen en la firma del service (sus DTOs) viven en `service/`, y `transport/` solo guarda lo
puramente HTTP. La interfaz del repositorio la declara el **service** (su consumidor), no el
paquete `repository`.

Referencias: `internal/modules/locations/` (el más simple), `internal/modules/staffing/` (con
puertos cross-módulo) y `internal/modules/publicapi/` (agregador puro, sin datos propios).

**Un paquete solo existe cuando tiene contenido real.** No crear módulos ni paquetes vacíos.

---

## 4. Las 7 reglas de modularidad

1. **`platform` nunca importa `modules`.** Los módulos consumen la plataforma, jamás al revés.
2. **Un módulo no importa el `repository/` ni el `transport/` de otro.** Solo su `service/`.
3. **La interfaz se define en el consumidor.** Si `recruitment` necesita datos de `staffing`,
   declara la interfaz mínima en `recruitment/ports.go` y el composition root le inyecta la
   implementación. `staffing` no se entera. (En Go el acoplamiento es estructural: no hace falta
   que el proveedor importe nada.)
4. **Sin claves foráneas nuevas entre dominios.** ID suelto + validación en el service dueño del
   dato.
5. **Los modelos GORM son compartidos** y viven en `internal/shared/models`. Es una decisión deliberada
   (ADR-001 regla 2): separarlos genera ciclos de importación por las relaciones bidireccionales.
   No es deuda técnica pendiente.
6. **No introducir `ports/`, `adapters/`, `usecase/` ni una clase por caso de uso.** Para una app
   mayormente CRUD sobre GORM, el boilerplate no compensa.
7. **Sin estado global.** Los repositorios reciben `*gorm.DB` inyectado, nunca el global
   `database.DB`.

Las reglas 1, 2 y 6 las verifica `depguard` en `.golangci.yml`. Los ciclos de importación los
prohíbe el compilador de Go.

---

## 5. Multi-tenancy — regla de oro

Para todo recurso de tenant, sin excepción:

1. **Listar/leer:** filtrar por el `company_id` del token.
2. **Crear:** ignorar cualquier `company_id` del body y forzar el del token.
3. **Leer/actualizar/eliminar por ID:** verificar pertenencia (`403` si no).
4. **Integridad cruzada:** al referenciar otra entidad, el service valida que pertenezca al mismo
   `company_id`.

Usar los helpers, no reimplementar: `authctx.RequireScope(c)`, `authctx.RequireCompanyScope(c)`,
`tenant.Scope.Guard(companyID)`.

---

## 6. Autorización — dos ejes ortogonales

- **`RequirePermission(rbac, perm)`** → ¿el **rol** puede? Permisos con formato
  `modulo.recurso.accion`, declarados en `internal/shared/permissions/` (un archivo por dominio).
- **`RequireFeature(plans, feature)`** → ¿el **plan** de la empresa lo incluye?

Un módulo premium se gatea con `RequireFeature`; las acciones dentro de él, además, con
`RequirePermission`.

---

## 7. Estrategia de desarrollo

Trabajar en **vertical slices**, no por capas horizontales:

```
regla de negocio → service → repository → transport → test
```

Priorizar código ejecutable sobre andamiaje vacío.

---

## 8. Qué NO hacer

- Crear un paquete global `services/`, `handlers/` o `repositories/` para todo Dvra.
- Compartir DTOs entre módulos: cada uno declara los suyos, aunque tengan los mismos campos.
- Usar el modelo GORM crudo como respuesta HTTP.
- Crear claves foráneas entre dominios.
- Meter lógica de negocio en handlers o en repositorios.
- Añadir dependencias sin justificar la necesidad.
- Implementar funcionalidad futura solo porque aparece en la documentación (la "Red Dvra" no
  existe y no se construye todavía).
- Saltarse la entrada de bitácora: es parte de la Definición de Done, no es opcional.

---

## 9. Antes de dar una tarea por terminada

```bash
make check        # gofmt + vet + lint + test + build
make swagger      # si tocaste endpoints
```

Y el checklist completo de la Definición de Done en `ESTANDARES_DESARROLLO.md` §9.

El ciclo de tarea es **HU → HT → código → validación**, con las skills `hu-a-ht` (entrada) y
`validar-estandares` (salida).
