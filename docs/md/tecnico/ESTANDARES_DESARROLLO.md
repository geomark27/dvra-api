# Estándares de Desarrollo — Dvra API

> Fuente de verdad de cómo se escribe, estructura, verifica y documenta el código en Dvra.
> Toda tarea se valida contra este documento **a la entrada** (traducción HU → HT) y **a la salida** (Definición de Done).
> Las skills de Claude `hu-a-ht` y `validar-estandares` automatizan esa validación (ver §8).
> El resumen operativo para agentes está en `AGENTS.md` (raíz del repo).

---

## 0. Cómo se usa este documento

El ciclo de cualquier tarea:

```
HU (Historia de Usuario, funcional)
   │  skill: hu-a-ht  → analiza con contexto del proyecto + estos estándares
   ▼
HT (Historia Técnica): a qué módulo pertenece, qué archivos tocar, estándares aplicables,
   criterios de aceptación, consideraciones multi-tenant/permisos/tests/docs
   │  desarrollo
   ▼
Código + tests + docs
   │  skill: validar-estandares  → valida diff contra §9 (Definición de Done)
   ▼
Tarea lista (cumple estándares de entrada y salida)
```

Este doc es **descriptivo de lo que el proyecto ya hace bien** + **prescriptivo de las brechas a cerrar**. Cuando una regla aún no se cumple en todo el código, se marca con ⚠️ y la dirección a seguir.

**El proyecto está migrando de arquitectura por capas a arquitectura por dominios.** El §3.7 (mapa de migración) dice qué norma aplica a cada dominio mientras dura la transición. Consúltalo antes de decidir dónde va un archivo nuevo.

---

## 1. Stack y versiones

- **Go 1.26.1** (ver `go.mod`). No usar features sobre esa versión sin actualizar `go.mod`.
- **Gin** (HTTP), **GORM** (ORM, PostgreSQL), **swag** (Swagger).
- Migraciones de esquema por **GORM AutoMigrate**: todo modelo nuevo va en `internal/database/models_all.go` (`AllModels`). Las migraciones de **datos** son versionadas, en `internal/database/migrations/data_migrations.go`.

---

## 2. Cómo se escribe el código

### 2.1 Formato y análisis estático (obligatorio)
- `gofmt` sin excepciones. Un archivo sin formatear no entra.
- `go vet ./...` debe pasar limpio.
- `go build ./...` debe compilar.
- `golangci-lint` limpio (incluye `depguard`, que impone las fronteras entre módulos).
- Indentación con **tabs** (lo que produce `gofmt`); nunca espacios.

### 2.2 Naming
- Paquetes: minúscula, sin guiones bajos (`permissions`, `authctx`, `transport`).
- Archivos: `snake_case` por entidad (`staffing_client_service.go`).
- Interfaces: `XService`, `XRepository`; implementación privada `xService`/`xRepository`; constructor `NewXService(...)`.
- **Evitar el *stutter***: dentro del paquete `service` del módulo `staffing`, el tipo es `StaffingClientService`, no `StaffingStaffingClientService`. El paquete ya da contexto.
- Exportar solo lo necesario; el resto en minúscula.

### 2.3 Manejo de errores
- **Repository**: devuelve `(nil, nil)` cuando el registro no existe (no propaga `gorm.ErrRecordNotFound` salvo casos puntuales). Errores de BD se propagan tal cual.
- **Service**: contiene la lógica y los mensajes de negocio. Usa `internal/shared/apperr`: `apperr.NotFound`, `Conflict`, `BadRequest`, `Forbidden`, `Unauthorized`.
- **Transport**: traduce el error a HTTP con `apperr.StatusCode(err)`. No mete lógica de negocio.
- Nunca tragar errores en silencio (`_ =`), salvo conversión de IDs ya validados por la ruta.

### 2.4 Validación de entrada
- Usar **tags `binding`** de Gin (se ejecutan con `ShouldBindJSON`/`ShouldBindQuery`).
- **No usar tags `validate`**: no hay validador go-playground conectado, así que no se ejecutan. (Las que existían fueron migradas; el repo está limpio de ellas y debe seguir así.)
- Validaciones de negocio (pertenencia, estado, reglas) van en el **service**, no en el DTO.

### 2.5 Idioma
- **Identificadores** (tipos, funciones, campos): inglés.
- **Comentarios y documentación**: español.
- **Mensajes de error de la API**: inglés (consistente con lo existente, p. ej. `"Access denied"`, `"No company context"`).

### 2.6 Contrato de respuesta HTTP (envelope)

Un único emisor: `internal/platform/httpx`. No construir `gin.H` a mano en un handler.

**Éxito, recurso único** — `data` es el payload directo:
```json
{ "status": "success", "data": { "id": 12, "title": "..." } }
```

**Éxito, colección** — `data` es siempre el array desnudo (no `data.jobs`):
```json
{
  "status": "success",
  "data": [ { "id": 12 } ],
  "meta": { "page": 1, "per_page": 20, "total": 137, "total_pages": 7 }
}
```

**Error** — `error` es el mensaje legible; `code` es un slug estable para lógica de UI; `errors` solo aparece en validación:
```json
{
  "status": "error",
  "error": "Los datos enviados no son válidos",
  "code": "validation_error",
  "errors": { "title": ["El título es obligatorio"] }
}
```

Códigos `code` vigentes: `validation_error`, `not_found`, `forbidden`, `unauthorized`, `conflict`, `feature_not_included`, `plan_limit_reached`, `cross_tenant_denied`, `internal_error`.

Status HTTP: `400` binding/validación · `401` sin sesión · `403` sin permiso, sin contexto de empresa o acceso cross-tenant denegado · `404` no encontrado · `409` conflicto o límite de plan · `201` creación · `200` lectura/actualización/borrado.

**Paginación:** `page` (default 1) y `per_page` (default 20, máx 100). Se pagina **solo lo no acotado**: jobs, candidates, applications, users, memberships, staffing-clients, placements, companies. **No** se paginan las colecciones acotadas: departments, system-values, tenant-roles, plans, locations — devuelven `data` sin `meta`.

---

## 3. Cómo se estructura (arquitectura por dominios)

Decidido en `adr/ADR-001-arquitectura-modular.md`. Referencia canónica: `internal/modules/staffing/`.

### 3.1 Las dos mitades

```
internal/
├── modules/<dominio>/     ← NEGOCIO. Consume la plataforma.
├── platform/              ← Transversal + landlord. NUNCA importa modules.
└── shared/                ← Tipos y utilidades sin lógica de dominio.
```

Dominios: `iam` (auth, users, companies, memberships, roles), `recruitment` (jobs, candidates, applications, departments), `staffing`, `billing` (plans, suscripciones, settings), `catalog` (system-values), `locations`, `dashboard`, `public`.

### 3.2 Anatomía de un módulo

```
internal/modules/<dominio>/
├── module.go        # Wiring + RegisterRoutes. Solo lo toca el composition root.
├── ports.go         # Interfaces que ESTE módulo necesita de otros.
├── transport/       # Handlers, DTOs de request/response, rutas.
├── service/         # Lógica de negocio y validaciones.
└── repository/      # Acceso a datos con GORM.
```

Reglas duras de dirección:
- **`transport`** hace bind, resuelve el scope, llama al service y mapea a DTO. **No** accede a GORM ni a `database.DB`.
- **`service`** tiene la lógica y las validaciones. **No** conoce `gin.Context` ni GORM.
- **`repository`** solo accede a datos. **No** conoce HTTP ni reglas de negocio.
- **Un paquete solo existe cuando tiene contenido real.** Nada de módulos ni carpetas vacías.

### 3.3 Fronteras entre módulos

1. **`platform` nunca importa `modules`.**
2. **Un módulo no importa el `repository/` ni el `transport/` de otro.** Solo su `service/`.
3. **La interfaz se define en el consumidor.** `recruitment/ports.go` declara lo que necesita; el composition root inyecta la implementación; el proveedor no se entera. Precedente en el código: `job_service.go` ya declara así `staffingClientReader` y `departmentReader`.
4. **Sin claves foráneas nuevas entre dominios.** ID suelto + validación en el service dueño.
5. **No introducir `ports/`, `adapters/`, `usecase/`** ni una clase por caso de uso.

Verificado por `depguard` en `.golangci.yml`. Los ciclos los prohíbe el compilador.

### 3.4 Modelos

Viven **compartidos** en `internal/shared/models` (decisión deliberada del ADR-001 regla 2: separarlos por módulo genera ciclos por las relaciones bidireccionales tipo `StaffingClient ↔ Job`). No es deuda pendiente.

- Embeber **`BaseModel`** (nunca `gorm.Model` ni campos `ID/CreatedAt/...` inline).
- Definir `TableName()`.
- Tags `gorm` + `json`. Índices compuestos pensados para las queries reales.
- Relaciones con `foreignKey` y `omitempty` en JSON.

### 3.5 DTOs

**Viven dentro de su módulo**, nunca en un paquete compartido: son formato de cable, y un módulo jamás usa los DTOs de otro.

**Dónde exactamente, dentro del módulo:** `transport` importa `service`, nunca al revés. De ahí sale la regla, que además evita ciclos de importación:

- **Si el tipo aparece en la firma del service** (lo recibe o lo devuelve), vive en `service/`: es parte del contrato del service, no del transporte. Ejemplo: `locations/service/dto_location.go`.
- **Si es puramente HTTP** (solo lo toca el handler), vive en `transport/`.

Ponerlos siempre en `transport/` no es posible: un service que devuelve sus DTOs tendría que importar `transport`, que a su vez lo importa a él.

- Por entidad: `CreateXDTO`, `UpdateXDTO` (punteros para parcial), `XResponseDTO` + `ToXResponse`/`ToXResponseList`.
- Nombre con prefijo de entidad (`JobResponseDTO`, no `ResponseDTO`): evita colisiones en el Swagger entre paquetes `transport` distintos.
- Los endpoints exponen **DTOs de respuesta**, nunca el modelo crudo.
- `CompanyID` se **fuerza desde el token**; nunca se confía en el body (ver §4).
- Los **filtros y criterios de consulta** que consume el service no son DTOs: van en `ports.go` o junto al service.

### 3.6 Inyección de dependencias

Manual y centralizada en `internal/platform/server/`. Cada módulo expone `New(db, deps)` y `RegisterRoutes(httpx.Router)`; el composition root los construye y los registra. Coste de un módulo nuevo: una línea.

`httpx.Router` da acceso a los tres grupos de rutas (`Public`, `Protected`, `Platform`) porque un mismo módulo puede servir a los tres contextos.

### 3.7 Mapa de migración (estado vigente)

**La migración terminó** (2026-08-28): `internal/app/` ya no existe. Todo el código nuevo va en `internal/modules/<dominio>/`, y un dominio nuevo nace siempre como módulo.


| Dominio | Estado | Ubicación |
|---|---|---|
| `iam` (auth, users, companies, memberships, RBAC) | ✅ migrado | `internal/modules/iam/` |
| `recruitment` (jobs, candidates, applications, departments) | ✅ migrado | `internal/modules/recruitment/` |
| `staffing` | ✅ migrado | `internal/modules/staffing/` |
| `billing` (plans, suscripciones, settings) | ✅ migrado | `internal/modules/billing/` |
| `catalog` (system-values) | ✅ migrado | `internal/modules/catalog/` |
| `locations` | ✅ migrado | `internal/modules/locations/` |
| `platformsettings` | ✅ migrado | `internal/modules/platformsettings/` |
| `dashboard` | ✅ migrado | `internal/modules/dashboard/` |
| `publicapi` (career page) | ✅ migrado | `internal/modules/publicapi/` |
| `health` | ✅ migrado | `internal/platform/health/` (infraestructura, no dominio) |

> Mantener esta tabla al día es parte de cada PR que añada un dominio.

---

## 4. Multi-tenancy — REGLA DE ORO

> **`company_id` es la ÚNICA frontera de aislamiento entre tenants.** Cualquier otra dimensión (p. ej. `staffing_client_id`) es un filtro interno, **no** una frontera.

Reglas obligatorias para todo recurso de tenant:
1. **Listar / leer**: filtrar por el `company_id` del token. SuperAdmin (sin `company_id`) tiene acceso global.
2. **Crear**: ignorar cualquier `company_id` del body y **forzar el del token**.
3. **Leer/actualizar/eliminar por ID**: verificar que `recurso.CompanyID == company_id` del token (si no, `403`).
4. **Integridad cruzada**: al referenciar otra entidad (p. ej. asignar `Job.StaffingClientID`, o crear un `Placement` desde una `Application`), el service valida que la entidad referenciada pertenezca al **mismo** `company_id`.

**Usa los helpers, no reimplementes el patrón:**

```go
scope, ok := authctx.RequireScope(c)        // 403 + abort ya emitidos si !ok
if !ok { return }
// ...y en el service:
if err := scope.Guard(recurso.CompanyID); err != nil { return nil, err }
```

- `authctx.RequireScope(c)` — resuelve el scope (global para SuperAdmin, de empresa para tenant).
- `authctx.RequireCompanyScope(c)` — exige contexto de empresa; rechaza platform/SuperAdmin.
- `authctx.RequireCompanyID(c)` — para forzar el `CompanyID` en creación.
- `tenant.Scope.Guard(companyID)` — comprobación de pertenencia (§4.3).

⚠️ **Brecha en cierre**: quedan ~15 usos de `IsSuperAdmin` con scoping manual, concentrados en `company_handler`, `membership_handler`, `user_handler` y `plan_handler` (los handlers de `jobs`, `candidates` y `applications` ya migraron a `authctx.TenantScope`). Se cierran al migrar `iam` y `billing`. **Código nuevo: siempre con helpers.** Un `c.Get("company_id")` con literal o un `IsSuperAdmin` nuevo fuera de `internal/shared/` es un hallazgo bloqueante en CI.

---

## 5. Autorización y entitlements

Dos capas **ortogonales** — una acción puede requerir ambas:
- **`RequirePermission(rbac, perm)`** → ¿el **rol** del usuario puede? Formato `modulo.recurso.accion`. Matriz en `internal/shared/permissions/` (un archivo por dominio, grants en `init()`), más RBAC relacional por tenant (`tenant_roles`/`permissions`), con efecto inmediato al cambiar un rol.
- **`RequireFeature(plans, feature)`** → ¿el **plan** de la empresa lo incluye? Flags en `Plan` (`CanUseX`) + `Plan.HasFeature(...)` + `CompanyHasFeature(...)`. SuperAdmin pasa siempre.

Regla: módulos premium se gatean con `RequireFeature`; las acciones dentro de ellos, además, con `RequirePermission`.

Los middlewares dependen de interfaces mínimas definidas en `internal/shared/middleware` (`Authenticator`, `PermissionChecker`, `FeatureChecker`), no de paquetes de servicios concretos. Eso es lo que permite que un módulo sea autocontenido.

---

## 6. Cómo se verifica (tests)

- **Obligatorio**: todo **service con lógica de negocio o validaciones** lleva tests unitarios (caminos felices + cada validación que devuelve error), con mocks de los puertos, sin BD.
- **Obligatorio**: todo módulo lleva un test de **aislamiento tenant** (que un scope de otra empresa no alcance sus datos).
- **Obligatorio**: cambios en la matriz de permisos se cubren en `permissions_test.go`.
- **Obligatorio en migración**: ningún dominio se mergea sin test de service + test de aislamiento.
- **Recomendado**: repos contra PostgreSQL de test cuando la query tenga lógica (filtros, joins, visibilidad por rol).
- Los tests corren en CI junto a `gofmt`, `go vet`, `go build` y `golangci-lint`.

⚠️ Cobertura actual: ~14 archivos de test, por debajo del 10%. Meta de corto plazo: ≥30% en la capa de servicios, y **100% de los caminos de límites de plan y ciclo de suscripción**.

---

## 7. Cómo se documenta

- **Swagger**: todo endpoint nuevo lleva anotaciones `godoc` (`@Summary`, `@Tags`, `@Param`, `@Success/@Failure`, `@Security BearerAuth`, `@Router`). Documentar el DTO concreto, nunca `map[string]interface{}`. Tras añadirlas: `make swagger`.
- **Bitácora** (`docs/md/bitacora/BITACORA.md`): una entrada por feature o decisión técnica, con el formato definido al final de ese archivo. Entradas nuevas al inicio. Si el cambio toca también el frontend, rellenar el campo `Contraparte`.
- **ADR** (`docs/md/tecnico/adr/`): cuando la decisión condiciona código futuro. La bitácora dice *qué se hizo*; el ADR dice *por qué se decidió así y qué se descartó*.
- **Docs técnicos vigentes** por módulo en `docs/md/tecnico/`.

### 7.1 Commits y ramas

Conventional Commits con **scope = módulo**, para que `git log --oneline` diga qué dominio se movió:

```
<tipo>(<scope>): <resumen en imperativo, minúscula, sin punto, ≤72 caracteres>

<cuerpo opcional: qué y por qué, no cómo>
```

- **Tipos:** `feat`, `fix`, `refactor`, `perf`, `docs`, `test`, `build`, `ci`, `chore`, `style`.
- **Scopes:** `iam`, `recruitment`, `staffing`, `billing`, `catalog`, `locations`, `platform`, `dashboard`, `public`, `httpx`, `shared`, `docs`, `infra`, `deps`.
- Un `!` antes de los dos puntos marca cambio incompatible.
- **Ramas:** `feat/<scope>-<slug>`, `fix/<scope>-<slug>`, `docs/<slug>`. Nada de commits directos a `main`.
- Cuando el cambio toca también el frontend, la rama lleva **el mismo nombre en ambos repos**.

Se impone con el hook `commit-msg` (`make hooks-install`); los targets `push`/`sync` del Makefile rechazan operar sobre `main`.

---

## 8. Proceso de tarea: HU → HT (entrada) y validación (salida)

### 8.1 Entrada — traducir HU a HT (skill `hu-a-ht`)
Dada una Historia de Usuario, producir una **Historia Técnica** que incluya:
- **Módulo al que pertenece** y su estado en el mapa de migración (§3.7).
- **Resumen funcional** y criterios de aceptación.
- **Archivos** a crear/modificar (según §3.2 si es módulo, Anexo A si es dominio pendiente).
- **Estándares aplicables**: multi-tenant (§4), permisos/entitlements (§5), DTOs/validación (§2.4/§3.5), errores (§2.3), envelope y paginación (§2.6).
- **Modelo de datos**: entidades, relaciones, índices, migración (`AllModels`).
- **Plan de pruebas** (§6) y **documentación** requerida (§7).
- **Riesgos** (especialmente fuga cross-tenant) y decisiones abiertas a confirmar.

Las HT se guardan en `docs/md/historias/HT-<slug>.md`.

### 8.2 Salida — validar cumplimiento (skill `validar-estandares`)
Al terminar, validar el diff contra la **Definición de Done** (§9) y reportar incumplimientos con ubicación y corrección sugerida.

---

## 9. Definición de Done (checklist de salida)

- [ ] `make check` limpio (gofmt, vet, lint con `depguard`, test, build).
- [ ] El código está en el módulo correcto según el mapa de migración (§3.7).
- [ ] Fronteras respetadas (§3.3): sin importar `repository/` ni `transport/` de otro módulo; sin `platform → modules`.
- [ ] Capas respetadas (§3.2): sin GORM en transport, sin `gin.Context` en services o repos.
- [ ] Multi-tenant completo (§4): scoping en list/get, `company_id` forzado en create, pertenencia verificada en get/update/delete, integridad cruzada en el service. Con helpers, no a mano.
- [ ] Autorización correcta (§5): `RequirePermission` y, si aplica, `RequireFeature`.
- [ ] Modelo con `BaseModel` + `TableName()` + en `AllModels` si es nuevo.
- [ ] DTOs de respuesta en el `transport/` del módulo (no modelo crudo) y validación con `binding`.
- [ ] Envelope, códigos HTTP y paginación correctos (§2.6), emitidos vía `httpx`.
- [ ] Tests de la lógica nueva + test de aislamiento tenant (§6).
- [ ] Anotaciones Swagger con DTO concreto + `make swagger` (§7).
- [ ] Entrada en la bitácora (§7).

---

## 10. Roadmap de evolución

- **Completado** (ago-2026): migración por dominios (§3.7). `internal/app/` eliminado.
- **Corto plazo**: cerrar la brecha de scoping manual del §4; observabilidad (logs estructurados con `slog` y request ID, métricas); rate limiting.
- **v2 (post-beta)**: evaluar aislamiento por **base de datos por empresa** (patrón landlord + BD por tenant). Requiere que la migración modular esté completa, porque hoy hay relaciones GORM que cruzan la frontera landlord/tenant y no sobrevivirían a bases separadas.
- Estos cambios se registran como ADR y entrada de bitácora **antes** de ejecutarse.

---

## Anexo A — Arquitectura por capas (histórico)

Hasta agosto de 2026 el proyecto se organizaba por capa técnica en
`internal/app/{handlers,services,repositories,dtos,models}`. Esa estructura **ya
no existe**: los diez dominios se migraron a módulos y los modelos GORM pasaron
a `internal/shared/models`.

Se documenta solo para leer commits y documentos anteriores a esa fecha. No es
una alternativa válida para código nuevo.
