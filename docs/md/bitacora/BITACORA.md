# Bitácora de Actividades — Dvra API

> Registro cronológico de implementaciones, auditorías y decisiones técnicas.
> Consolida los antiguos `MULTI_COMPANY_IMPLEMENTATION.md`, `SECURITY_AUDIT.md` y `SUPERADMIN_IMPLEMENTATION.md`.
> Las entradas más recientes van al inicio. Para detalle vigente de cada módulo, ver `docs/md/tecnico/`.

---

## 2026-08-28 — Migración modular completada: `billing`, `iam`, `recruitment`, `dashboard`, `publicapi`

**Contexto:** cierre de la migración por dominios iniciada con el ADR-001. En esta tanda entraron
los cinco dominios restantes, incluidos los dos más grandes (`iam`, con toda la superficie de
autenticación, y `recruitment`, el núcleo del ATS) y los dos agregadores.

**Qué se hizo:**
- **`billing`**: planes, suscripciones y configuración de empresa. Absorbió `internal/shared/billing`
  (hallazgo P10, ahora cerrado del todo). Expone `Plans()`, que satisface directamente
  `middleware.FeatureChecker`: así llega el gating por plan a las rutas premium.
- **`iam`**: auth, JWT, usuarios, empresas, membresías y RBAC relacional. `jwtService` y
  `rbacService` ganaron `Session()` y `ResolveRole()` nativos, con lo que **se eliminó el adaptador
  temporal que el composition root arrastraba desde la Fase A**.
- **`recruitment`**: departamentos, vacantes, candidatos y pipeline. Es el módulo mayor.
- **`dashboard`**: **eliminados sus 15 usos del global `database.DB`**, último reducto del estado
  global. Además `.Table("jobs")` pasó a `.Model(&models.Job{})`, para que el nombre de la tabla
  venga de `TableName()` y no de una cadena.
- **`publicapi`** (career page): agregador puro, sin datos propios. Consume cinco puertos hacia iam
  y recruitment.
- **Retrofit de `staffing`**: sus DTOs salieron de `internal/app/dtos` al módulo; los filtros de
  consulta pasaron a `domain/` (no son formato de cable); y **por fin tiene tests** — 5 casos sobre
  su regla de negocio real: una colocación solo nace de una postulación en etapa `hired`, del mismo
  `company_id` tanto en el cliente final como en la postulación, y solo una vez por postulación.
- **`internal/app/` ELIMINADO.** Los modelos GORM pasaron a `internal/shared/models`: el nombre
  `app/models` mentía sobre una capa que ya no existe.
- `registerRoutes` quedó reducido a su esqueleto: arma los tres grupos y deja que cada módulo monte
  lo suyo. Ya no contiene un solo endpoint de dominio.
- **ADR-004** con las decisiones que la ejecución obligó a concretar y que el ADR-001 no cubría.
- `depguard` actualizado a las fronteras reales del código final.

**Decisiones clave:**
- **Los DTOs viven en `service/`, no en `transport/`**, porque el service los devuelve y ponerlos en
  transport crea un ciclo. La regla general —`transport` importa `service`, nunca al revés— quedó
  en §3.5 y en el ADR-004.
- **`staffing` y `recruitment` se construyen en dos fases.** Cada uno necesita algo del otro: no es
  un ciclo de importación (cada uno declara su puerto) pero sí de construcción. Se resuelve con
  `BindApplications` tras construir ambos, en lugar de inicialización perezosa, para que un fallo
  se vea al arrancar y no en la primera petición.
- **Un adaptador donde los tipos coinciden en forma pero no en nombre**: `ResumeUpload` existe en
  recruitment y en publicapi, y Go no los unifica. Traducir en el composition root es el precio de
  que ninguno dependa de los tipos del otro.
- **Cada módulo declara sus propias variables de error.** `ErrPlanNotFound` y `ErrCompanyNotFound`
  están duplicadas a propósito, con el mismo mensaje: un módulo no toma prestadas las de otro. Se
  comprobó antes que nadie las comparaba por identidad.

**Verificado:**
- `TestRouteContract` → **121 rutas idénticas** tras mover los diez dominios
- `gofmt`, `go vet ./...`, `go build ./...`, `go build ./cmd/...`, `go test ./...` → limpios
  (15 paquetes con tests)
- `make swagger`: **61 paths idénticos**; las 49 definiciones son las mismas, solo renombradas al
  paquete de su módulo
- Fronteras comprobadas por grep: ningún módulo importa `internal/app`; `shared` no importa
  `modules`; `platform` (salvo el composition root) no importa `modules`; `shared/tenant` no
  importa gorm

**Pendientes:**
- [ ] ⚠️ `golangci-lint` local sigue sin poder ejecutarse (compilado con go1.24, el proyecto apunta
      a go1.26.1). **Las reglas `depguard` están sin verificar por la herramienta**; se comprobaron
      a mano por grep. Confirmar en CI o actualizando el binario.
- [ ] ⚠️ `go run ./cmd/console migrate --fresh --seed` sigue sin ejecutarse: no hay PostgreSQL ni
      Docker en este entorno. Arrastrado desde la migración de `catalog`.
- [ ] `PublishJobWithinLimit` y `AssignPlanToCompany` conservan `*gorm.DB` en el service por
      necesitar transacción con bloqueo de fila (excepción documentada en ADR-004)
- [ ] Empezar las épicas de beta (E0-E9) del plan de desarrollo

**Referencia vigente:** `tecnico/adr/ADR-004-plantilla-de-modulo.md`,
`tecnico/ESTANDARES_DESARROLLO.md` §3

---

## 2026-08-28 — Migración del módulo `catalog`

**Contexto:** cuarto dominio migrado, y el primero con dos particularidades: tiene **scoping
multi-tenant real** (valores propios de cada empresa y overrides sobre los globales) y tiene
**consumidores internos** — `internal/shared/catalog` lo usaban 5 sitios: los seeders de planes y
empresas, y los services de `plan`, `company` y `auth`.

**Qué se hizo:**
- `internal/modules/catalog/` con la plantilla de módulo. Sus 11 rutas (7 de tenant y 4 de
  landlord) salieron de `server/routes.go`. Es el primer módulo que usa a la vez los grupos
  **protegido y landlord**, con lo que `httpx.Router` queda validado en sus tres grupos entre este
  módulo y `platformsettings`.
- **Hallazgo P10 cerrado**: `internal/shared/catalog` no era una utilidad transversal, era lógica
  de este dominio con consultas GORM propias. Se movió a
  `internal/modules/catalog/repository/resolver.go` y el paquete `shared/catalog` desapareció.
- La interfaz del repositorio pasó al consumidor (`service/ports.go`), §3.3 regla 3.
- **Tests nuevos** (13 casos, incluidos subtests). Los importantes son los de **aislamiento
  multi-tenant**, obligatorios en este módulo por ser el primero con datos de empresa: que la
  empresa A no pueda editar ni borrar un valor de la empresa B, ni uno global, ni uno inexistente
  — y que en los tres casos responda `404` y no `403`, porque un `403` confirmaría que el valor
  existe. Verificado por mutación: al quitar el guard, los tests fallan.

**Decisiones clave:**
- El resolver sigue siendo funciones sin estado que reciben la conexión, y sus 5 consumidores solo
  cambiaron la ruta del import. Es deliberado y **transitorio**: esos consumidores pertenecen a
  `billing` e `iam`, que aún no son módulos y pasan el `*gorm.DB` que ya tienen. Cuando migren,
  declararán su puerto mínimo y recibirán el resolver inyectado.

**Bugs corregidos:**
- Ninguno en producción, pero el test destapó un supuesto equivocado **mío** sobre el dominio: dar
  de alta un valor de tenant exige `IsExtensible` **y** `AllowTenantValues`, no solo el segundo.
  Son cosas distintas: una categoría puede ser extensible solo por el landlord. El test se corrigió
  y ahora cubre las tres combinaciones que deben rechazarse.

**Verificado:**
- `TestRouteContract` → **121 rutas idénticas**
- `gofmt`, `go vet ./...`, `go build ./...`, `go test ./...` → limpios (13 paquetes con tests)
- `make swagger`: 61 paths **y** 49 definiciones idénticas — este módulo no movió ningún tipo
  referenciado desde las anotaciones
- El cuerpo del resolver es byte a byte el mismo que en `shared/catalog`: solo cambió el paquete
- Fronteras: el módulo no importa `app/*`; su `service` no conoce gin ni gorm

**Pendientes:**
- [ ] ⚠️ **`go run ./cmd/console migrate --fresh --seed` NO se pudo ejecutar**: no hay PostgreSQL
      accesible ni Docker en este entorno. Era el criterio de aceptación de esta fase por el
      cambio de import en los seeders. El riesgo real es bajo (las llamadas se resuelven en
      compilación, el build pasa y el cuerpo del resolver no cambió), pero **queda por confirmar
      en un entorno con base**.
- [ ] Siguiente dominio: `billing` (planes, suscripciones, company settings)

**Referencia vigente:** `tecnico/adr/ADR-003-catalogos-modulares-y-suscripciones.md`,
`tecnico/ESTANDARES_DESARROLLO.md` §3.7

---

## 2026-08-28 — Migración de `platformsettings` y `health`

**Contexto:** segundo y tercer dominio tras la Fase A. `platformsettings` es la configuración
global del SaaS (marca, contacto, URLs, textos legales, días de prueba y plan por defecto): un
singleton que se sirve **público** al frontend antes del login y solo edita el SuperAdmin. Es el
primer módulo que usa dos grupos de rutas distintos, lo que valida el diseño de `httpx.Router`.

**Qué se hizo:**
- `internal/modules/platformsettings/` con la plantilla de módulo. Se movió con `git mv`
  (repositorio, service, DTOs y handler), y sus 3 rutas salieron de `server/routes.go`.
- **Deuda P7 cerrada en este dominio**: el service dependía del *struct concreto* del
  repositorio, lo que impedía probarlo sin base. Ahora depende de un puerto declarado en
  `service/ports.go` (§3.3 regla 3).
- `health` **no** se llevó a `modules/`: es infraestructura, no dominio de negocio. Vive en
  `internal/platform/health/`.
- **Violación de capas corregida en `health`**: el handler recibía `*gorm.DB` crudo y hacía el
  ping él mismo. Ahora depende del puerto `Pinger`, con el adaptador `gormPinger` en el
  composition root.
- Se separó el significado de las dos sondas: `/health` (vida) **ya no depende de la base**. Antes
  sí, y eso hacía que un problema de base provocara reinicios del proceso, que es justo lo que no
  hay que hacer. `/health/ready` (disponibilidad) sí la comprueba, con timeout de 2 s.
- **Tests nuevos**: 4 en `platformsettings/service` (campos de marca vacíos rechazados, días de
  prueba negativos, actualización parcial que conserva lo no enviado y registra al autor,
  propagación de errores) y 2 en `health` (que la sonda de vida sobrevive a la base caída y que la
  de disponibilidad refleja su estado). Verificados por mutación.

**Decisiones clave:**
- El módulo se llama `platformsettings` y **no** `platform`, para no confundirse con
  `internal/platform/`, que es infraestructura (config, servidor, almacenamiento, httpx) y no un
  dominio de negocio. El mapa de migración (§3.7) se actualizó con ese nombre.

**Verificado:**
- `TestRouteContract` → **121 rutas idénticas**: el contrato HTTP no cambió
- `gofmt`, `go vet ./...`, `go build ./...`, `go test ./...` → limpios (12 paquetes con tests)
- `make swagger`: 61 paths idénticos; solo cambian de paquete los 3 DTOs de platform settings
- Fronteras comprobadas: ningún módulo importa `app/*`; sus `service` no conocen gin ni gorm;
  `health` no conoce gorm
- `internal/app/{handlers,services,repositories,dtos}` ya no contienen nada de locations,
  platform settings ni health

**Pendientes:**
- [ ] Siguiente dominio: `catalog` (system-values), luego `billing`
- [ ] `PlatformSettingsService` sigue siendo un struct concreto: quien lo consuma desde otro
      módulo debe declarar su propio puerto mínimo, no depender del tipo

**Referencia vigente:** `tecnico/ESTANDARES_DESARROLLO.md` §3.7

---

## 2026-08-27 — Migración del módulo `locations`

**Contexto:** primer dominio migrado tras la Fase A. Se eligió `locations` por ser el más simple
de los pendientes: catálogo de referencia global, solo lectura, sin scoping por empresa y sin
dependencias salientes. Sirve para calibrar la plantilla de módulo con volumen real (1.787 líneas)
y riesgo bajo.

**Qué se hizo:**
- Movido a `internal/modules/locations/` con `git mv` (se conserva la historia):
  `repository/`, `service/`, `transport/` y `module.go`.
- **Eliminados los 30 usos del global `database.DB`**: `NewLocationRepository()` pasa a recibir
  `*gorm.DB` inyectado, con lo que el repositorio queda mockeable.
- La interfaz del repositorio se movió al **consumidor** (`service/ports.go`), según §3.3 regla 3.
- El módulo se registra por `httpx.Module`; sus 14 rutas salieron de `server/routes.go`.
- **Tests nuevos** (`service/location_service_test.go`, 6 casos) sobre la lógica real:
  normalización de ISO a mayúsculas en consulta y en alta, `404` en no encontrado, integridad
  referencial al crear una subregión, y propagación de errores del repositorio. Se verificaron por
  mutación: introducir el bug hace fallar el test.
- No lleva test de aislamiento tenant **a propósito**: el catálogo es global, no pertenece a
  ninguna empresa. Queda anotado en `module.go`.

**Decisiones clave:**
- Los DTOs quedan en `service/` y no en `transport/`: el service los devuelve, así que ponerlos en
  `transport` crearía un ciclo (`transport → service → transport`). Se documentó la regla general
  en §3.5, con `transport` importando a `service` y nunca al revés.

**Bugs corregidos:**
- **`make swagger` no reproducía el `swagger.json` commiteado**: al target le faltaban
  `--parseDependency --parseInternal`, así que regenerar renombraba **todas** las definiciones
  (`dvra-api_internal_app_dtos.X` → `dtos.X`). Corregido en el Makefile. Sin esto, cualquier
  regeneración habría producido un diff global espurio.

**Verificado:**
- `TestRouteContract` → **121 rutas idénticas**: el contrato HTTP no cambió
- `gofmt`, `go vet ./...`, `go build ./...`, `go test ./...` → limpios
- `make swagger`: 61 paths idénticos; el único cambio son los 16 DTOs de locations moviéndose de
  `app_dtos` a `modules_locations_service`
- Fronteras comprobadas: el módulo no importa `app/*`, su `service` no conoce gin ni gorm, su
  `transport` no conoce gorm

**Operación:**
- `swag` debe ser la versión de `go.mod` (v1.16.6). Ojo: esa versión **reporta mal su propio
  número** (`const Version = "v1.16.4"` upstream), así que `swag --version` no sirve para
  verificarla; instálala con `go install github.com/swaggo/swag/cmd/swag@v1.16.6`.

**Pendientes:**
- [ ] Siguiente dominio: `platform` (settings + health), luego `catalog`
- [ ] Fijar la versión de `swag` en el target `install-tools` para evitar la deriva

**Referencia vigente:** `tecnico/ESTANDARES_DESARROLLO.md` §3.5 y §3.7

---

## 2026-08-27 — Fase A: plataforma de migración modular

**Contexto:** `internal/shared/middleware` importaba `internal/app/services`, así que la cadena
`staffing/transport → shared/middleware → app/services` metía **todos** los servicios en el grafo
de cualquier módulo. Mientras eso existiera, ningún módulo podía ser autocontenido y cada fase de
migración habría repetido el mismo trabajo. Esta fase no migra ningún dominio: construye el
terreno.

**Qué se hizo:**
- **Puertos de middleware** (`shared/middleware/ports.go`): `Authenticator`, `PermissionChecker`,
  `FeatureChecker`, definidos en el consumidor. `TenantRBACService.Can` y
  `PlanService.CompanyHasFeature` los satisfacen tal cual, sin adaptador. Solo `Authenticator`
  necesita uno, temporal, en el composition root: desaparece cuando se migre `iam`.
- **Tipos neutrales** (`shared/identity/session.go`): `SessionClaims`, `TenantRole` y los errores
  de token. `services.ErrExpiredToken` pasa a ser un alias del de `identity`, de modo que el
  código existente sigue compilando.
- `authctx.SetSession` es ahora el **único escritor** de las keys del contexto; `AuthMiddleware`
  dejó de usar literales.
- **Helpers de scoping** (`authctx/scope.go` + `tenant.Scope.Guard`): `RequireScope`,
  `RequireCompanyScope`, `RequireCompanyID`. Son el reemplazo del patrón `IsSuperAdmin` repetido a
  mano, que se irá aplicando al migrar `iam` y `billing`.
- **`shared/tenant` ya no importa gorm**: `Apply` se movió al subpaquete `tenant/tenantgorm`, así
  el núcleo de la frontera de datos puede usarse desde la lógica de negocio de cualquier módulo.
- **`platform/httpx`**: `Router` (tres grupos de rutas + RBAC + Plans) y la interfaz `Module`.
  `staffing` ya se registra por ese contrato.
- **Composition root**: `registerRoutes` pasa de **22 parámetros posicionales a 3**; el cableado
  vive en `buildDeps` (`server/deps.go`), con una lista `[]httpx.Module` para los dominios ya
  migrados.
- **`depguard`** en `.golangci.yml` con 4 reglas: `platform !→ modules` (salvo el composition
  root), fronteras de módulo, `shared !→ services`, y `tenant` sin gorm. Es el equivalente en Go a
  las pruebas ArchUnit de Terra.
- **Golden de rutas** (`server/routes_test.go` + `testdata/routes.golden`, **121 rutas**): congela
  el contrato HTTP. Se verificó que detecta altas y bajas de rutas.
- **Código muerto eliminado**: `database/scoped_db.go` y `middleware/cors_middleware.go` (0 usos
  ambos), `RequireCompany()`, `OptionalAuth()`. `go mod tidy` quitó `gin-contrib/cors` y
  `uber/mock`.
- Tests de `shared/middleware` reescritos con un `Authenticator` falso: el paquete ya no depende de
  `app/services` **ni siquiera en sus tests**. Se añadieron los casos de token expirado, token
  desconocido, header mal formado y rol revocado.

**Decisiones clave:**
- El adaptador `authenticator` vive en el composition root y es explícitamente temporal.
- `platform/httpx` va en `platform/` y no en `platform/server` para evitar el ciclo
  `server → modules → server`.
- El golden de rutas **no se toca durante la migración**: si una fase lo cambia, es un bug de esa
  fase, no un cambio legítimo.

**Verificado:**
- `gofmt -l internal/ cmd/` → vacío
- `go vet ./...` → limpio
- `go build ./...` → OK
- `go test ./...` → todo verde, incluido el nuevo `TestRouteContract`
- Las 4 reglas de frontera comprobadas por grep sobre el código real
- ⚠️ `golangci-lint` local no ejecutable: está compilado con go1.24 y el proyecto apunta a
  go1.26.1. El YAML se validó por separado; **falta confirmar `depguard` en CI**.

**Pendientes:**
- [ ] Actualizar `golangci-lint` local a una versión compilada con Go ≥1.26 y correr `make lint`
- [ ] Migrar los dominios pendientes (§3.7), empezando por `locations`
- [ ] Aplicar los helpers de scoping a los ~15 `IsSuperAdmin` restantes al migrar `iam` y `billing`
- [ ] ADR-004 con las decisiones que esta fase añade al ADR-001

**Referencia vigente:** `tecnico/ESTANDARES_DESARROLLO.md` §3.3 y §3.7, `AGENTS.md`

---

## 2026-08-26 — Gobernanza de la migración modular y plan de beta

**Contexto:** `ESTANDARES_DESARROLLO.md` §3 seguía describiendo la arquitectura **por capas** como
norma vigente, en contradicción con ADR-001 (aceptado, con `staffing` ya migrado). Como las skills
`hu-a-ht` y `validar-estandares` leen ese documento como fuente de verdad, todo el código nuevo se
generaba y validaba en la arquitectura que se está abandonando. Además no existía convención de
commits (4 de 34 usaban Conventional Commits) ni bitácora en el frontend.

Se tomó como referencia el patrón de `terra-api-purchasing`, que resultó **más simple** que el
hexagonal-lite que prescribía el ADR-001: dos mitades `modules/` y `platform/`, sin `ports/`,
`adapters/` ni `usecase/`, y con las fronteras verificadas por análisis estático en lugar de
artefactos separados.

**Qué se hizo:**
- `AGENTS.md` en la raíz: resumen operativo con la anatomía de módulo y las 7 reglas de
  modularidad.
- `ESTANDARES_DESARROLLO.md` §3 reescrito con la arquitectura por dominios
  (`modules/<dominio>/{transport,service,repository}` + `module.go` + `ports.go`). La estructura
  por capas pasa al **Anexo A**, marcada como legacy en retirada.
- Nuevo **§3.7 Mapa de migración**: tabla dominio → migrado/pendiente que indica dónde va el
  código nuevo mientras dura la transición. Un dominio nuevo nace siempre como módulo.
- §2.6 reescrito con el envelope definitivo (`data` como payload directo, `meta` de paginación,
  errores con `code` y `errors` por campo) y las reglas de qué colecciones se paginan.
- §4 corregido: la brecha de scoping manual son ~15 usos de `IsSuperAdmin` en `iam` y `billing`,
  no 41 — los handlers de jobs, candidates y applications ya migraron a `authctx.TenantScope`. Se
  documentan los helpers a usar.
- §2.4 corregido: las tags `validate:` inertes ya no existen en el repo (138 tags `binding:`
  activas); deja de ser una brecha y pasa a ser una prohibición.
- §6 actualizado: la cobertura ya no es "solo `permissions_test.go`"; se añade el test de
  aislamiento tenant como obligatorio por módulo.
- Skills `hu-a-ht` y `validar-estandares` actualizadas: la primera pregunta el módulo antes que
  nada y emite rutas según §3.7; la segunda añade checks de frontera entre módulos.
- Formato de bitácora formalizado con las secciones que ya se usaban de facto (`Contexto`,
  `Decisiones clave`, `Bugs corregidos`, `Verificado`, `Operación`) y el campo `Contraparte`.

**Decisiones clave:**
- La migración modular se ejecuta **intercalada** con el camino a beta, congelando cada dominio en
  las semanas en que las épicas de producto lo tocan. Cero refactor durante la entrada de los
  primeros clientes.
- El aislamiento por **base de datos por empresa** se pospone a v2 (post-beta): hoy hay 31 de 64
  `Preload` cruzando la frontera landlord/tenant y claves foráneas entre dominios que no
  sobrevivirían a bases separadas. La migración modular es su prerrequisito, no una tarea
  paralela.
- Los modelos GORM siguen compartidos (ADR-001 regla 2). Es deliberado, no deuda.

**Pendientes:**
- [ ] Convención de commits: hook `commit-msg` + protección de `main` en los targets del Makefile
- [ ] Consolidar los dos árboles de documentación (`docs/01..04` → `docs/md/_archivo/`)
- [ ] Fase A: desacoplar `shared/middleware` de `app/services`, `platform/httpx`, helpers de
      scoping, `depguard`, golden de rutas
- [ ] Migrar los 8 dominios pendientes (§3.7)
- [ ] ADR-004 con las decisiones que este plan añade al ADR-001

**Referencia vigente:** `tecnico/ESTANDARES_DESARROLLO.md` §3 y §3.7,
`tecnico/adr/ADR-001-arquitectura-modular.md`, `AGENTS.md`

**Contraparte:** dvra-web — entrada `2026-08-26 — Bitácora, estándares y alcance de beta`

---

## 2026-07-20 — Catálogos modulares, suscripciones y archivo de ADR

Se reemplazó el catálogo plano por la jerarquía módulo → categoría → valor,
con valores globales, valores tenant y overrides aislados por empresa. Se
incorporaron suscripciones relacionales empresa → plan, settings propios del
tenant y ejecución de migraciones/seeders por módulo.

Los Architecture Decision Records se separaron de la documentación técnica
operativa y ahora viven en `docs/md/tecnico/adr/`. Su `README.md` actúa como
índice histórico con fecha y estado. La bitácora conserva el hito de
implementación; los ADR conservan el contexto y la justificación de las
decisiones.

**Referencia vigente:** `docs/md/tecnico/adr/ADR-003-catalogos-modulares-y-suscripciones.md`.

## 2026-07-20 — Estructura organizacional para el ATS

Se incorporó `departments` con jerarquía autorreferenciada, responsable
opcional, lifecycle activo/inactivo y aislamiento por empresa. Las vacantes
aceptan y filtran `department_id`. Se añadieron permisos `departments.*`, CRUD,
vista plana y árbol, constraints case-insensitive, pruebas y Swagger. Los jobs
tenant y públicos incluyen `department.path`, calculado mediante una consulta
recursiva sin exponer endpoints privados en la career page.

## 2026-07-17 — RBAC tenant relacional y perfil `interviewer`

El perfil genérico `user` se renombró a `interviewer`. Las memberships
existentes se migran de forma idempotente a una FK `tenant_role_id`; la columna
textual `role` se elimina.

Se incorporaron `tenant_roles`, `permissions` y `tenant_role_permissions`. Los
grants iniciales siguen versionados en código, pero middleware y servicios
resuelven rol y permisos desde BD. Crear o editar usuarios/memberships requiere
`tenant_role_id`; `GET /api/v1/tenant-roles` expone las opciones activas. El
middleware relee la membership vigente para que cambios de rol y suspensiones
tengan efecto inmediato.

El perfil `admin` se sincroniza con el catálogo tenant completo y una migración
adicional repara instalaciones que ya hubieran sembrado grants incompletos. La
regla no concede permisos landlord ni omite los límites de features del plan.

**Referencia vigente:** `internal/shared/permissions/`,
`internal/database/migrations/data_migrations.go` y ADR-002.

## 2026-07-15 — Separación landlord/tenant y contextos de autorización

Se reemplazó la membresía ambigua de SuperAdmin con `company_id = NULL` por
`platform_memberships`. Las memberships tenant requieren company, son únicas por
usuario/tenant y conservan exclusivamente los perfiles `admin`, `recruiter`,
`hiring_manager` y `user`.

JWT access/refresh incluyen `context_type` y `membership_id`. Login selecciona
un contexto activo, refresh lo conserva y valida contra BD, `switch-company`
rota ambos tokens y `switch-platform` permite volver al landlord. Los permisos
`platform.*` se evalúan separadamente de los permisos tenant y las rutas
landlord viven en `/api/v1/platform/*`.

Se agregó migración para perfiles legacy, comando `console platform-admin`,
respuesta de contexto en `/auth/me` y filtros por recurso para que
`hiring_manager` solo vea jobs asignados y sus candidatos/postulaciones.

**Referencia vigente:** `docs/md/tecnico/adr/ADR-002-contextos-landlord-tenant.md`.

## 2026-07-15 — Remediación de seguridad, tenancy, JWT y storage

**Qué se hizo:** se separó la validación de access/refresh JWT; se bloqueó el rol reservado `superadmin` en memberships de empresa; `POST /users` ahora crea credenciales y membresía transaccional; y las mutaciones de usuarios quedaron tenant-scoped. Jobs usa únicamente `draft/published/closed`, aplica `max_jobs` al publicar y Applications valida Job/Candidate contra el mismo tenant.

Se introdujo `tenant.Scope` fail-closed y se migraron repositorios ATS críticos a `*gorm.DB` inyectado. Los CVs usan filesystem abstraído, nombres generados, firmas de archivo, límite de 10 MB, cuota del plan, metadatos en `stored_files` y descarga autenticada. Docker monta un volumen persistente.

**Operación:** `console migrate` ejecuta migraciones idempotentes que convierten jobs `active` a `published` y memberships tenant con rol `superadmin` a `admin`. Producción rechaza secretos JWT inseguros; seeders demo requieren `SEED_DEMO_DATA=true`; readiness hace ping real a PostgreSQL y el servidor cierra con graceful shutdown.

**Verificado:** tests unitarios de JWT, memberships, tenant scope, configuración y filesystem; `go test ./...`, `go vet ./...` y `go build ./...`.

**Pendientes:** límites soft de candidatos/usuarios, rate limiting, invitaciones por email y backend S3.

**Referencia vigente:** `docs/04_DOCUMENTACION_TECNICA_API.md`, `internal/shared/tenant/`, `internal/platform/storage/`.

## 2026-06-20 — Piloto monolito modular + hexagonal-lite: módulo `staffing`

**Contexto:** primer paso del ADR-001. Se migró `staffing` de la organización por capa técnica (`internal/app/{models,dtos,...}`) a un módulo autocontenido en `internal/modules/staffing/`, como piloto para validar el patrón antes de replicarlo.

**Estructura resultante:**
- `domain/ports.go` — interfaces (StaffingClientRepository, PlacementRepository) + el puerto cross-módulo `ApplicationFinder`/`HiredApplication`. No importa gin ni gorm.
- `service/` — casos de uso (StaffingClientService, PlacementService). Dependen solo de `domain`, `models`, `apperr`.
- `repository/` — adaptadores GORM con `*gorm.DB` **inyectado** (ya no usan el global `database.DB` → testeables).
- `transport/` — handlers HTTP + `routes.go` (adaptador de entrada).
- `module.go` — wiring del módulo + `RegisterRoutes`.

**Decisiones clave (del ADR-001):**
- **Modelos compartidos.** `models.StaffingClient`/`Placement` y los DTOs **se quedaron** en `internal/app/{models,dtos}`. Separarlos provocaría ciclos de importación por las relaciones GORM mutuas (`StaffingClient↔Job`, `Placement↔Application`). Es la 1ª pasada pragmática.
- **Puertos definidos por el consumidor para romper ciclos cross-módulo:**
  - `staffing` necesita datos de `Application` (recruitment) → define `domain.ApplicationFinder`; el composition root (`platform/server/wiring_adapters.go`) inyecta un adaptador sobre `ApplicationRepository`.
  - `recruitment` (job_service) necesita validar un cliente final → ahora depende de una interfaz local `staffingClientReader`; el composition root le pasa `staffingModule.ClientRepo`.
  - Resultado: **ningún módulo importa al otro**. Grafo acíclico verificado por `go build`. El único consumidor de `internal/modules/staffing` es el composition root.
- **Acople aceptado:** el módulo importa `app/services` solo por el *tipo* `PlanService` (entitlement, para `RequireFeature`). Podría convertirse en otro puerto del consumidor más adelante.

**Archivos eliminados** (reemplazados por el módulo): `services/{staffing_client,placement}_service.go`, `repositories/{staffing_client,placement}_repository.go`, `handlers/{staffing_client,placement}_handler.go`.

**Verificado:** `go build ./...`, `go vet ./...`, `gofmt -l` (limpio), `go test ./...` y build de ambos binarios pasan. Sin cambio de comportamiento ni de contrato de API (mismas rutas, mismos códigos).

**Pendientes / siguiente paso (del ADR-001):**
- [ ] Tests del módulo `staffing` (ya es trivial: repos con `db` inyectado + services mockeables vía puertos).
- [ ] Replicar el patrón al resto de dominios (recruitment, iam, billing, platform), uno por uno con build/tests verdes.

**Referencia vigente:** `docs/md/tecnico/adr/ADR-001-arquitectura-modular.md`, `internal/modules/staffing/`

---

## 2026-06-20 — Unificación de `auth` y `plan` a `apperr`

**Contexto:** en la pasada anterior se dejaron `auth` y `plan` fuera porque ya mapeaban códigos correctamente vía errores centinela + ladders `if err == services.ErrX`. Para uniformar el manejo de errores en TODO el dominio, ahora también usan `apperr`.

**Qué se hizo:**
- **Nuevo constructor `apperr.Unauthorized` (401)** para los errores de credenciales de auth.
- **Centinelas convertidos a valores `*apperr.AppError`** (en `plan_service.go` y `auth_service.go`): siguen siendo variables exportadas comparables con `errors.Is`/`==` (compatibilidad total), pero ahora llevan su código HTTP. Mapeo: `ErrPlanNotFound`/`ErrUserNotFound`/`ErrCompanyNotFound` → 404; `ErrPlanSlugExists`/`ErrPlanInUse`/`ErrEmailExists` → 409; `ErrInvalidCredentials`/`ErrInvalidPassword` → 401; `ErrNoMembership` → 403; `ErrInvalidPlanData` → 400. Inline: `account is inactive` → 403, `invalid refresh token` → 401, `plan is not active` → 400.
- **Handlers simplificados** (`plan_handler`, `auth_handler`): se eliminaron los ladders `if err == services.ErrX { ... }` (≈12 bloques) a favor de `apperr.StatusCode(err)`. Menos código y misma (o más correcta) semántica de códigos.
- **Preservado por seguridad/config:**
  - `RefreshToken` mantiene su **401 genérico** para cualquier fallo (decisión deliberada: el frontend trata 401 en refresh como "cerrar sesión"; un 403/404 podría no disparar ese flujo).
  - `free plan is not available` sigue como error plano → 500 (es un problema de configuración del servidor, no del cliente).
  - Los mensajes de credenciales siguen siendo genéricos ("invalid email or password") — no revelan si el email existe.

**Verificado:** `go build ./...`, `go vet ./...`, `gofmt -l` (limpio), `go test ./...` y build de ambos binarios (`dvra-api`, `console`) pasan. Resultado: 14 services y 14 handlers usan `apperr` (solo quedan fuera `dashboard`/`health`, sin errores de negocio).

**Referencia vigente:** `internal/shared/apperr/errors.go`, `internal/app/services/{auth,plan}_service.go`, `internal/app/handlers/{auth,plan,platform_settings}_handler.go`

---

## 2026-06-20 — Adopción de `apperr` en todo el dominio (códigos HTTP correctos)

**Contexto:** tras introducir el paquete `internal/shared/apperr` en el módulo staffing, se extendió el mismo patrón al resto de la API. El problema transversal: los services devolvían errores genéricos (`fmt.Errorf`/`errors.New`) y los handlers respondían **HTTP 500 para errores que son del cliente** (no encontrado, duplicado, validación). Además había dos antipatrones: detección de tipo de error por comparación de strings (`err.Error() == "company not found"`, `strings.Contains(err, "already applied")`) y `strconv.ParseUint` ignorando el error de parseo.

**Qué se hizo:**
- **Services convertidos a `apperr`** (`NotFound`/`Conflict`/`BadRequest`/`Forbidden`): `job`, `candidate`, `company`, `user`, `application`, `membership`, `public`, `platform_settings`, `system_value`, `location`. Los errores de negocio ahora llevan su código HTTP; los fallos internos reales (wraps con `%w`: creación de directorios/candidato/aplicación) se mantienen como errores planos → siguen siendo 500.
- **Handlers actualizados** para usar `apperr.StatusCode(err)` en lugar de códigos hardcodeados tras llamadas al service: `job`, `candidate`, `company`, `user`, `application`, `membership`, `public`, `system_value`, `location` (+ los `staffing`/`placement` ya hechos). Resultado: slug/email duplicado → 409, recurso inexistente → 404, validación → 400.
- **Eliminado el matching frágil de strings**: `company_handler` y `user_handler` ya no comparan `err.Error() == "..."`; `public_handler` ya no usa `strings.Contains` para clasificar el error de postulación.
- **`strconv.ParseUint` validado** en los handlers que lo ignoraban (`job` ×5, `candidate` ×4, `application` ×5): ID no numérico → 400 en vez de un 404 confuso.
- **Mensaje genérico preservado donde corresponde**: en el endpoint público `ApplyToJob`, los errores internos (500) siguen devolviendo "Failed to submit application" para no filtrar detalles; solo los errores de negocio exponen su mensaje.

**Decisiones de alcance (lo que NO se tocó y por qué):**
- **`auth` y `plan` se dejaron intactos**: ya mapeaban códigos correctamente vía errores centinela (`errors.Is`/`==`). Además `auth` usa mensajes genéricos deliberados ("Invalid email or password") por seguridad; reescribirlos arriesgaba filtrar si un email existe. No tenían el bug del 500.
- **No se refactorizó el scoping multi-tenant** (`c.Get("company_id")` → `authctx`, ítems **C1/M1** de la auditoría) ni se eliminaron los double-lookups de los handlers genéricos: son cambios mayores y ortogonales a "aplicar apperr".

**Nota de comportamiento:** los handlers que ya devolvían `err.Error()` en 500 (la mayoría) lo siguen haciendo; ahora también lo hacen `company`/`user` (antes usaban mensajes genéricos). Es consistente con la norma del código, pero expone el texto del error en fallos 5xx — pendiente de endurecer si se requiere.

**Verificado:** `go build ./...`, `go vet ./...`, `gofmt -l` (limpio) y `go test ./...` (permissions) pasan.

**Referencia vigente:** `internal/shared/apperr/errors.go` y todos los `*_service.go` / `*_handler.go` listados arriba.

---

## 2026-06-20 — Correcciones al módulo Staffing / Placement (auditoría de calidad)

**Contexto:** revisión completa del módulo staffing/placement (modelos, DTOs, repositorios, servicios, handlers y rutas) para corregir bugs reales, inconsistencias de diseño y malos códigos HTTP detectados en auditoría. El código compilaba pero tenía errores de comportamiento silenciosos.

**Bugs corregidos:**

- **Error de compilación encubierto:** `placement_service.go` tenía un método `GetAll` con cuerpo vacío que no compilaría. Se eliminó.
- **Código HTTP incorrecto en errores de negocio:** `CreateStaffingClient` y `UpdateStaffingClient` devolvían HTTP 500 cuando el slug ya existía (conflicto de datos). Ahora devuelven 409 Conflict. El 500 indica fallo del servidor; el 409 indica que el cliente envió datos que chocan con datos existentes.
- **ID inválido en URL no se validaba:** si alguien llamaba `GET /staffing-clients/abc`, el sistema ignoraba el error de parseo, buscaba el registro con ID=0 y devolvía un 404 sin explicación. Ahora devuelve 400 Bad Request con mensaje claro, en los 6 handlers afectados (get, update, delete de ambos módulos).
- **Placement duplicado permitido:** se podía llamar `POST /placements` dos veces con la misma `application_id` y se creaban dos registros. Ahora el service verifica unicidad y devuelve 409 si ya existe un placement para esa application.
- **Método fantasma en repositorio:** `PlacementRepository` declaraba `GetAll` y `GetByCandidate` en su interface pero ningún service ni handler los llamaba. Código muerto en la interface confunde y obliga a las implementaciones futuras a cumplir contratos inútiles. Ambos fueron eliminados de la interface y sus implementaciones.
- **`deleted_at IS NULL` manual y redundante:** `GetAll` del repositorio filtraba explícitamente `deleted_at IS NULL`, algo que GORM ya hace solo con soft delete. Era redundante e inconsistente con el resto del repo.

**Inconsistencias de diseño corregidas:**

- **Doble consulta a BD en Update y Delete:** los handlers de `UpdateStaffingClient`, `DeleteStaffingClient`, `UpdatePlacement` y `DeletePlacement` hacían dos consultas separadas — una en el handler para verificar el tenant, y otra dentro del service para traer el registro. Ahora el service recibe el `companyID` directamente y hace una sola consulta que verifica todo. Se agregó la convención: `companyID = 0` significa SuperAdmin (omite validación de tenant).
- **Modelos crudos en el DTO de respuesta:** `PlacementResponseDTO` incluía `*models.StaffingClient` y `*models.Candidate` directamente — los modelos de base de datos expuestos en la API. Esto puede filtrar campos internos no deseados. Ahora usa `*StaffingClientResponseDTO` y `*CandidateResponseDTO` (DTOs de presentación).
- **`JobID` y `ApplicationID` nullable sin razón:** en el modelo `Placement`, ambos campos eran punteros (`*uint`, nullable en BD) aunque la lógica de negocio siempre los rellena al copiarlos de una `Application` (que los tiene como `not null`). Cambiados a `uint not null`. Esto fuerza la constraint en BD vía AutoMigrate y simplifica el código del service (ya no necesita `&jobID`).

**Nuevo paquete:**

- **`internal/shared/apperr`**: tipos de error con código HTTP asociado (`NotFound`, `Conflict`, `BadRequest`, `Forbidden` y `StatusCode(err)`). Permite que los services señalen el tipo de error y los handlers elijan el status code correcto sin inspeccionar strings de mensajes. Resuelve el patrón repetido de "todos los errores del service → 500".

**Versión:**

- Swagger `@version` en `main.go` sincronizado a `1.4.0` (estaba en `1.2.0`; la ruta raíz ya devolvía `v1.4.0`).

**Verificado:** `go build ./...` y `go vet ./...` pasan sin errores ni advertencias.

**Pendientes que siguen abiertos (de auditoría anterior):**
- [ ] Regenerar Swagger: `make swagger`
- [ ] BD ya poblada: planes Professional/Enterprise necesitan `can_use_staffing=true` manual en BD existente
- [ ] Tests del módulo staffing (ahora hay más superficies que probar: unicidad de application_id, apperr codes)
- [ ] Ítem de auditoría **A1** sigue abierto: 47 tags `validate` que no se ejecutan

**Referencia vigente:** `internal/shared/apperr/errors.go`, `internal/app/models/placement.go`, `internal/app/repositories/placement_repository.go`, `internal/app/services/{staffing_client,placement}_service.go`, `internal/app/handlers/{staffing_client,placement}_handler.go`, `internal/app/dtos/placement_dto.go`

---

## 2026-06-14 — Estándares de Desarrollo: documento + skills + enforcement + auditoría

**Contexto:** formalizar los estándares de Dvra (cómo se escribe y estructura el código) no solo por escrito, sino **operacionalizados** como skills de Claude que validan a la entrada y a la salida de cada tarea.

**Qué se hizo:**
- **Documento fuente de verdad** `docs/md/tecnico/ESTANDARES_DESARROLLO.md`: stack, cómo se escribe/estructura/verifica/documenta el código, **regla de oro multi-tenant** (`company_id` única frontera), autorización vs entitlement, proceso **HU → HT** y **Definición de Done**.
- **Skills** (`.claude/skills/`): `hu-a-ht` (puerta de entrada: traduce Historia de Usuario funcional a Historia Técnica validando estándares) y `validar-estandares` (puerta de salida: valida el diff contra la Definición de Done). Ambas anclan al documento como fuente de verdad.
- **Enforcement automático**: `.golangci.yml` (v1.64.8: errcheck, govet, ineffassign, staticcheck, unused, gofmt — sin `misspell` por falsos positivos en español, sin `goimports` local-prefixes por churn), `.editorconfig`, CI en `.github/workflows/ci.yml` (gofmt+vet+build+test+lint), hook `scripts/git-hooks/pre-commit`, y targets `make check`/`fmt-check`/`hooks-install` (+ `swag` en `install-tools`).
- **Baseline verde**: se cerró deuda mínima pre-existente — 4 `errcheck` (`CloseDB` en CLI, `UpdateLastLogin` best-effort) y 4 archivos sin `gofmt`. `make check` pasa entero.
- **Auditoría de brechas** `docs/md/tecnico/AUDITORIA_ESTANDARES.md`: hallazgos priorizados con evidencia numérica.

**Pendientes (del plan de remediación de la auditoría):**
- [ ] **A1** — Unificar validación: 47 tags `validate` en 6 DTOs no se ejecutan; migrar a `binding`.
- [ ] **C1 (Fase 4)** — Centralizar el scoping multi-tenant (41 usos de `IsSuperAdmin` repetidos); riesgo de fuga cross-tenant. Candidato a estrenar las skills.
- [ ] **M1** — Migrar 7 handlers de `c.Get("company_id")` a `authctx`.
- [ ] **M2** — Unificar repos hacia `db` inyectado (10 usan global `database.DB`); habilita tests.
- [ ] **A2** — Tests de services con validaciones (hoy 1 test / 16 services).
- [ ] **B1** — Check en CI de que Swagger esté regenerado.

**Referencia vigente:** `docs/md/tecnico/ESTANDARES_DESARROLLO.md`, `docs/md/tecnico/AUDITORIA_ESTANDARES.md`, `.claude/skills/`

---

## 2026-06-14 — Módulo Staffing / Outsourcing (StaffingClient + Placement)

**Contexto:** habilitar a Dvra para firmas de staffing/outsourcing (que gestionan talento para clientes finales), no solo empresas que reclutan para sí mismas. Caso motivador: una firma que presta personal a un cliente final. Se modeló como segmento objetivo sin romper el multi-tenancy existente.

**Decisión de arquitectura:** el cliente final NO es un tenant. `company_id` sigue siendo la ÚNICA frontera de aislamiento; el cliente final es una dimensión interna del tenant.

**Qué se hizo:**
- Nuevos modelos `StaffingClient` (cliente final dentro del tenant) y `Placement` (colocación de un candidato en un cliente final, con datos de contrato y billing). Reescritos para embeber `BaseModel` (antes usaban `gorm.Model`, sin tags JSON ni `TableName()`) y registrados en `AllModels` para AutoMigrate (antes estaban huérfanos: las tablas nunca se creaban).
- `Job.StaffingClientID *uint` opcional: `nil` = reclutamiento propio; con valor = modo staffing. Filtro `GET /jobs?staffing_client_id=` y exposición en respuestas.
- `Placement` **deriva de una `Application` en etapa `hired`**: `candidate_id`/`job_id` se copian de ella (no se confían del body). El service valida integridad cross-tenant: application y cliente final deben pertenecer al mismo `company_id`. Igual validación al asignar `Job.StaffingClientID`.
- CRUD completo de ambos (dto/repo/service/handler + rutas `/api/v1/staffing-clients` y `/api/v1/placements`), con scoping manual por `company_id` (mismo patrón que jobs) y slug único por empresa validado en service (no índice único, para no chocar con soft delete).
- **Entitlement por plan:** nuevo flag `Plan.CanUseStaffing` + `Plan.HasFeature("staffing")` + método `PlanService.CompanyHasFeature()`. Nuevo middleware `RequireFeature(planService, "staffing")` (ortogonal a `RequirePermission`: la acción exige plan que lo incluya Y rol que lo permita). Seeder: Professional y Enterprise lo traen activo.
- Nuevos permisos en `permissions/staffing.go` (`staffing_clients.*`, `placements.*`) con grants por rol (delete solo admin, igual que jobs; `user` no ve placements por el billing).
- Anotaciones Swagger en los 10 endpoints nuevos.
- Verificado con `go build ./...`, `go vet ./...` y tests (`permissions` pasa).

**Pendientes:**
- [ ] Regenerar Swagger: `make swagger` (requiere instalar `swag`; no estaba en el entorno al implementar)
- [ ] BD ya poblada: `SeedPlans` omite planes existentes, así que las filas Professional/Enterprise NO obtienen `can_use_staffing=true` solas (AutoMigrate agrega la columna con default `false`). Actualizar vía `PUT /plans/:id` o SQL. En BD fresca (`make fresh`) quedan correctas.
- [ ] `PUT /jobs/:id` no permite desasignar el cliente (mandar `staffing_client_id: null` se interpreta como "no enviado") — requiere flag explícito si se necesita.
- [ ] Tests unitarios del flujo staffing (validaciones de integridad y entitlement)
- [ ] Evaluar control field-level sobre pay/bill rate (datos sensibles) por rol

**Referencia vigente:** `internal/app/models/{staffing_client,placement}.go`, `internal/shared/permissions/staffing.go`, `internal/shared/middleware/feature_middleware.go`

---

## 2026-06-12 — Sistema de Permisos por Acción (RBAC en código)

**Qué se hizo:**
- Nuevo paquete `internal/shared/permissions`: matriz rol → permisos en código, organizada por módulos (jobs.go, candidates.go, applications.go, users.go, companies.go, memberships.go, dashboard.go), análogo a un seeder de permisos por módulo pero verificado en compile-time. Punto único de consulta: `permissions.Can(role, permission)` — si a futuro se necesitan roles personalizados por empresa, solo se cambia la implementación de `Can()` a BD con caché, sin tocar rutas ni constantes.
- Nuevo middleware `RequirePermission(perm)` aplicado a TODAS las rutas protegidas en `routes.go` (estilo declarativo, como `can:` de Laravel). Antes ningún endpoint distinguía roles dentro de la empresa: un rol `user` (solo lectura) podía crear jobs o mover el pipeline.
- Permisos no asignados a ningún rol = exclusivos de SuperAdmin (companies.create/delete, memberships.create según RN-MEMB-004).
- Nuevo paquete `internal/shared/authctx`: accesores tipados al contexto (`Role`, `IsSuperAdmin`, `UserID`, `CompanyID`). Refactor de los 34 `if role == "superadmin"` dispersos en 8 handlers → `authctx.IsSuperAdmin(c)`; el literal "superadmin" ya no existe fuera de `permissions`/`authctx`.
- Fix: `getRoleLevel` no incluía `superadmin` (habría recibido nivel 0 con `RequireRole`); ahora nivel 100.
- `GET /auth/me` ahora devuelve `role` y `permissions` de la empresa activa para que el frontend muestre/oculte acciones (mismo contrato que con Spatie).
- Primeros tests unitarios del proyecto: `permissions_test.go` valida la matriz por rol.

**Pendientes:**
- [ ] Chequeos a nivel de recurso (hiring_manager edita solo sus jobs asignados) — requiere tabla de asignaciones (RN-MEMB-007)
- [ ] Refactorizar `c.Get("company_id")` en handlers a `authctx.CompanyID(c)`
- [ ] Cuando exista invitaciones (Fase 2): otorgar `memberships.invite` a admin

**Referencia vigente:** `negocio/LOGICA_DE_NEGOCIO.md` sección 3.2 (matriz) y `internal/shared/permissions/`

---

## 2026-06-12 — BaseModel unificado + diseño de perfiles internos/externos

**Qué se hizo:**
- Nuevo `models.BaseModel` (`internal/app/models/base.go`): equivalente a `gorm.Model` pero con tags JSON en snake_case (`id`, `created_at`, `updated_at`, `deleted_at`).
- Refactor de TODOS los modelos para embeber `BaseModel`: los 6 que duplicaban los campos inline (User, Company, Job, Candidate, Application, Membership) y los que usaban `gorm.Model` (Plan, Role, ubicaciones, SystemValue, PlatformSettings). No se usó `gorm.Model` directo porque no tiene tags JSON y los primeros 6 se serializan directamente en las respuestas de la API.
- `BaseModel` genera las mismas columnas que `gorm.Model`, así que no hay impacto en BD; excepción: `platform_settings` gana `created_at` y `deleted_at` (AutoMigrate las agrega; el soft delete es inerte porque el singleton nunca se elimina).
- Convención a futuro: todo modelo nuevo embebe `BaseModel`, nunca `gorm.Model` ni campos inline.
- Verificado con `go build ./...`, `go vet ./...` y tests.
- Documentado el diseño de **perfiles internos vs externos** (recruiters freelance/agencias) en `negocio/LOGICA_DE_NEGOCIO.md` sección 3.4: RN-MEMB-006 (`member_type`), RN-MEMB-007 (alcance del externo), RN-MEMB-008 (atribución y reporting por recruiter), y su incorporación al roadmap Fase 2.

**Pendientes:**
- [ ] Implementar invitaciones por email (prerequisito, ya era Fase 2)
- [ ] Implementar `member_type` + asignación vacante ↔ recruiter
- [ ] Historial de pipeline con atribución
- [ ] Reportes por recruiter

**Referencia vigente:** `negocio/LOGICA_DE_NEGOCIO.md` (sección 3.4)

---

## 2025-12-22/23 — Módulo de Ubicaciones (Locations)

**Qué se hizo:**
- Módulo completo de datos geográficos: Region → Subregion → Country → State → City.
- Capas completas: models, DTOs, repository, service, handler y rutas.
- Rutas GET públicas (sin auth); CRUD restringido a SuperAdmin (`/admin/locations/*`).
- Seeder básico incluido en `make fresh` (~7s) y carga masiva opcional con `make db-location` (157k+ registros).
- Comando `make fresh` corregido y Makefile actualizado.

**Referencia vigente:** `tecnico/MODULO_UBICACIONES.md`

---

## 2025-12-09 — Módulo de Planes (Plans)

**Qué se hizo:**
- CRUD completo de planes de suscripción gestionado por SuperAdmin.
- Rutas públicas para pricing page (`GET /plans`, `GET /plans/:slug`).
- Asignación de plan a empresa (`POST /admin/plans/assign`) actualizando `company.plan_tier`.
- Seeder con 4 planes: Free ($0), Starter ($39.99), Professional ($79.99), Enterprise ($159.99).
- Validaciones: slug único, no eliminar planes en uso (soft delete), planes validados contra BD (sin strings hardcodeados) al crear empresas.

**Referencia vigente:** `tecnico/MODULO_PLANES.md`

---

## 2025-12-08 — Separación Admin / SuperAdmin

**Qué se hizo:**
- Reorganización de la arquitectura para separar Admin regular (scoped a empresa) de SuperAdmin (acceso global, sin `company_id` en el JWT).
- Nuevos archivos: `dtos/superadmin_dto.go`, `services/admin/superadmin_companies_service.go`, `handlers/admin/superadmin_companies_handler.go`.
- Middleware `RequireSuperAdmin()`: valida JWT + role `superadmin` + ausencia de `company_id`.
- Rutas `/api/v1/admin/*`: listar/crear empresas, cambiar plan, suspender, ver usuarios de cualquier empresa, analytics globales (MRR, churn, totales).
- Creación de empresa + admin + membership en una sola transacción.

**Pendientes registrados en su momento:**
- [ ] Testing automatizado de handlers SuperAdmin
- [ ] Audit logs de acciones del SuperAdmin
- [ ] Impersonation (`POST /admin/impersonate/:user_id`)
- [ ] ⚠️ Cambiar credenciales por defecto del SuperAdmin en producción

**Referencia vigente:** `tecnico/FLUJO_SUPERADMIN.md`

---

## 2025-12-08 — Auditoría de Seguridad Multi-Tenancy

**Problema detectado (crítico):** los endpoints devolvían datos de TODAS las empresas sin filtrar por `company_id`; usuarios de una empresa podían ver, crear, modificar y eliminar recursos de otras.

**Correcciones aplicadas:**
- **Lectura:** filtrado por `company_id` del token en listados de memberships, users, companies, jobs, candidates y applications. SuperAdmin (sin `company_id`) conserva acceso global.
- **Lectura individual:** `GET /:id` valida que el recurso pertenezca a la empresa del token (403 si no).
- **Creación:** `POST /jobs|candidates|applications` ignora el `company_id` del body y fuerza el del token.
- **Actualización/Eliminación:** verifican pertenencia del recurso antes de proceder. `DELETE /companies/:id` solo SuperAdmin.
- **Memberships (MVP):** solo SuperAdmin puede crear memberships (`POST /admin/memberships`); admins de empresa reciben 403.
- Métodos nuevos: `UserRepository.GetByCompanyID()`, `UserService.GetUsersByCompanyID()`.

**Pendientes registrados en su momento:**
- [ ] `PUT /users/:id` y `DELETE /users/:id`: falta verificar memberships
- [ ] Tests unitarios de filtrado por empresa
- [ ] Logging de intentos de acceso cross-company
- [ ] Rate limiting contra enumeración

**Referencia vigente:** `tecnico/FLUJO_CLIENTE.md` (sección de aislamiento multi-tenant)

---

## 2025-12-08 — Sistema Multi-Empresa (v1.3.0)

**Qué se hizo:**
- Registro de empresa con su primer admin en una transacción: `POST /auth/register-company` (Company + User + Membership con role admin e `is_default: true`). `POST /auth/register` quedó deprecated.
- Login devuelve lista de empresas del usuario; el token usa la empresa por defecto.
- `POST /auth/switch-company`: genera nuevo token con distinto `company_id` (y el rol correspondiente) sin re-login.
- `GET /auth/my-companies`: lista de empresas del usuario para el selector de la UI.
- Un usuario puede tener N memberships con roles distintos por empresa (ej. admin en una, recruiter en otra).
- Errores nuevos: `ErrCompanyNotFound`, `ErrNoMembership`.

**Pendientes registrados en su momento:**
- [ ] **Fase 2 (prioritario):** sistema de invitaciones por email (`POST /memberships/invite`, `POST /auth/accept-invite`) para que admins de empresa puedan invitar usuarios sin depender del SuperAdmin
- [ ] Fase 3: onboarding post-registro (wizard, configuración inicial, importación de datos)

**Referencia vigente:** `tecnico/FLUJO_CLIENTE.md` y `tecnico/API_ENDPOINTS.md`

---

## Formato para nuevas entradas

Tres secciones **obligatorias** (`Qué se hizo`, `Pendientes`, `Referencia vigente`) y cinco
opcionales que se usan cuando aportan. Entradas nuevas **al inicio** del archivo.

```markdown
## AAAA-MM-DD — Título corto

**Contexto:** (opcional — 1-3 líneas: qué problema motivó el cambio)

**Qué se hizo:**
- ...

**Decisiones clave:** (opcional — solo si condicionan código futuro; si la decisión
es arquitectónica, aquí va el enlace al ADR, no la argumentación)
- ...

**Bugs corregidos:** (opcional)
- ...

**Verificado:** (opcional — comandos ejecutados y su resultado)
- `make check` → OK

**Operación:** (opcional — migraciones, variables de entorno, pasos de despliegue)
- ...

**Pendientes:**
- [ ] ...

**Referencia vigente:** `tecnico/...`

**Contraparte:** dvra-web — entrada `AAAA-MM-DD — <mismo título>`
  (obligatorio solo si el cambio toca ambos repositorios)
```

### Reglas

- Una entrada por **feature o decisión técnica**, no por commit.
- **Bitácora vs ADR:** la bitácora dice *qué se hizo y cuándo*; el ADR dice *por qué se decidió
  así y qué se descartó*. Si una entrada tiene "Decisiones clave" que restringen código futuro,
  debe enlazar un ADR existente o crear uno.
- Cuando el cambio toca también el frontend, ambas bitácoras llevan una entrada con la **misma
  fecha y el mismo título**, cada una describiendo su lado, y el campo `Contraparte` relleno.
- La bitácora del frontend vive en la ruta espejo: `dvra-web/docs/md/bitacora/BITACORA.md`.
