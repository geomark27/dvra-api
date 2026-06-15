# Estándares de Desarrollo — Dvra API

> Fuente de verdad de cómo se escribe, estructura, verifica y documenta el código en Dvra.
> Toda tarea se valida contra este documento **a la entrada** (traducción HU → HT) y **a la salida** (Definición de Done).
> Las skills de Claude `hu-a-ht` y `validar-estandares` automatizan esa validación (ver §8).

---

## 0. Cómo se usa este documento

El ciclo de cualquier tarea:

```
HU (Historia de Usuario, funcional)
   │  skill: hu-a-ht  → analiza con contexto del proyecto + estos estándares
   ▼
HT (Historia Técnica): qué capas/archivos tocar, estándares aplicables,
   criterios de aceptación, consideraciones multi-tenant/permisos/tests/docs
   │  desarrollo
   ▼
Código + tests + docs
   │  skill: validar-estandares  → valida diff contra §9 (Definición de Done)
   ▼
Tarea lista (cumple estándares de entrada y salida)
```

Este doc es **descriptivo de lo que el proyecto ya hace bien** + **prescriptivo de las brechas a cerrar**. Cuando una regla aún no se cumple en todo el código, se marca con ⚠️ y la dirección a seguir.

---

## 1. Stack y versiones

- **Go 1.24.0** (ver `go.mod`). No usar features sobre esa versión sin actualizar `go.mod`.
- **Gin** (HTTP), **GORM** (ORM, PostgreSQL), **swag** (Swagger), framework interno **Loom**.
- Migraciones por **GORM AutoMigrate** (no SQL manual): todo modelo nuevo va en `internal/database/models_all.go` (`AllModels`).

---

## 2. Cómo se escribe el código

### 2.1 Formato y análisis estático (obligatorio)
- `gofmt` sin excepciones. Un archivo sin formatear no entra.
- `go vet ./...` debe pasar limpio.
- `go build ./...` debe compilar.
- Indentación con **tabs** (lo que produce `gofmt`); nunca espacios.

### 2.2 Naming
- Paquetes: minúscula, sin guiones bajos (`permissions`, `authctx`).
- Archivos: `snake_case` por entidad (`staffing_client_service.go`).
- Interfaces de capa: `XService`, `XRepository`; implementación privada `xService`/`xRepository`; constructor `NewXService(...)`.
- Exportar solo lo necesario; el resto en minúscula.

### 2.3 Manejo de errores (patrón vigente)
- **Repository**: devuelve `(nil, nil)` cuando el registro no existe (no propaga `gorm.ErrRecordNotFound` salvo casos puntuales). Errores de BD se propagan tal cual.
- **Service**: contiene la lógica y los mensajes de negocio. Para "no encontrado" devuelve `fmt.Errorf("x not found")`. Errores reutilizables como variables `ErrXNotFound` cuando se comparan en varios lugares.
- **Handler**: traduce el error a HTTP. No mete lógica de negocio.
- Nunca tragar errores en silencio (`_ =`), salvo conversión de IDs ya validados por la ruta.

### 2.4 Validación de entrada
- Usar **tags `binding`** de Gin (se ejecutan con `ShouldBindJSON`/`ShouldBindQuery`).
- ⚠️ Las tags `validate` que aparecen en algunos DTOs **no se ejecutan** (no hay validador go-playground conectado). No depender de ellas. Convención: migrar a `binding` o documentar explícitamente si se añade un validador.
- Validaciones de negocio (pertenencia, estado, reglas) van en el **service**, no en el DTO.

### 2.5 Idioma
- **Identificadores** (tipos, funciones, campos): inglés.
- **Comentarios y documentación**: español.
- **Mensajes de error de la API**: inglés (consistente con lo existente, p. ej. `"Access denied"`, `"No company context"`).

### 2.6 Contrato de respuesta HTTP (envelope)
- Éxito: `gin.H{"status": "success", "data": ...}` (listados incluyen `"count"`).
- Error: `gin.H{"error": "<mensaje>"}` con el status HTTP correcto.
- Códigos: `400` binding/validación, `401` sin sesión, `403` sin permiso/sin contexto de empresa/acceso cross-tenant denegado, `404` no encontrado, `201` creación, `200` lectura/actualización/borrado.

---

## 3. Cómo se estructura (arquitectura)

### 3.1 Capas y dirección de dependencias
```
handler  →  service  →  repository  →  models / database
   │           │
  DTOs ◄───────┘
```
Reglas duras:
- El **handler** solo hace: bind de entrada, auth/scoping (contexto), llamar al service, mapear a DTO de respuesta. **No** accede a `database.DB` ni a GORM.
- El **service** tiene la lógica de negocio y las validaciones. **No** conoce `gin.Context` ni HTTP.
- El **repository** solo accede a datos (GORM). **No** conoce HTTP ni reglas de negocio.
- Los **models** no dependen de capas superiores.
- Dependencias siempre hacia abajo; nunca un repo importando services/handlers.

### 3.2 Inyección de dependencias
- Manual y centralizada en `internal/platform/server/server.go`: se construyen repos → services → handlers y se pasan a `registerRoutes`.
- Services y repos se exponen por **interface**; los handlers reciben la interface, no la implementación.

### 3.3 Modelos
- Embeber **`BaseModel`** (nunca `gorm.Model` ni campos `ID/CreatedAt/...` inline) — da tags JSON en snake_case.
- Definir `TableName()`.
- Campos con tags `gorm` + `json`. Índices compuestos pensados para las queries reales (p. ej. `idx_jobs_company_status`).
- Relaciones declaradas con `foreignKey` y `omitempty` en JSON.

### 3.4 DTOs
- Por entidad: `CreateXDTO`, `UpdateXDTO` (punteros para parcial), `XResponseDTO` + `ToXResponse`/`ToXResponseList`.
- Los endpoints exponen **DTOs de respuesta**, no el modelo crudo (evita filtrar relaciones/campos sensibles por accidente).
- `CompanyID` en `CreateXDTO` se **fuerza desde el token** en el handler; nunca se confía en el body (ver §4).

### 3.5 Convención de archivos por entidad
Una entidad nueva = un archivo por capa, mismo nombre base:
`models/x.go`, `dtos/x_dto.go`, `repositories/x_repository.go`, `services/x_service.go`, `handlers/x_handler.go`, y registro en `models_all.go` + rutas en `routes.go` + wiring en `server.go`.

---

## 4. Multi-tenancy — REGLA DE ORO

> **`company_id` es la ÚNICA frontera de aislamiento entre tenants.** Cualquier otra dimensión (p. ej. `staffing_client_id`) es un filtro interno, **no** una frontera.

Reglas obligatorias para todo recurso de tenant:
1. **Listar / leer**: filtrar por el `company_id` del token. SuperAdmin (sin `company_id`) tiene acceso global.
2. **Crear**: ignorar cualquier `company_id` del body y **forzar el del token**.
3. **Leer/actualizar/eliminar por ID**: verificar que `recurso.CompanyID == company_id` del token (si no, `403`).
4. **Integridad cruzada**: al referenciar otra entidad (p. ej. asignar `Job.StaffingClientID`, o crear un `Placement` desde una `Application`), el service valida que la entidad referenciada pertenezca al **mismo** `company_id`.

⚠️ **Brecha conocida**: hoy este scoping está implementado **a mano y repetido** en los handlers (41 usos de `IsSuperAdmin`, 9 handlers leyendo `company_id`). Cada copia es una oportunidad de fuga cross-tenant. `internal/database/scoped_db.go` (`ScopedDB`) existe para centralizarlo pero casi no se usa. **Dirección**: centralizar el scoping en un helper/middleware reutilizable (Fase 4 del roadmap de estándares). Mientras tanto, **toda nueva entidad debe replicar el patrón completo de los 4 puntos** sin omitir ninguno.

---

## 5. Autorización y entitlements

Dos capas **ortogonales** — una acción puede requerir ambas:
- **`RequirePermission(perm)`** → ¿el **rol** del usuario puede? Matriz en `internal/shared/permissions/` (un archivo por módulo, grants en `init()`). Consulta única vía `permissions.Can(role, perm)`. Permiso no asignado a ningún rol = exclusivo de SuperAdmin.
- **`RequireFeature(planService, feature)`** → ¿el **plan** de la empresa lo incluye? Flags en `Plan` (`CanUseX`) + `Plan.HasFeature(...)` + `PlanService.CompanyHasFeature(...)`. SuperAdmin pasa siempre.

Regla: módulos premium se gatean con `RequireFeature`; las acciones dentro de ellos, además, con `RequirePermission`.

---

## 6. Cómo se verifica (tests)

⚠️ **Brecha conocida**: solo existe `permissions_test.go`. Estándar a partir de ahora:
- **Obligatorio**: todo **service con lógica de negocio o validaciones** lleva tests unitarios (caminos felices + cada validación que devuelve error). Ej.: validación de etapa `hired` y de pertenencia cross-tenant en `PlacementService`.
- **Obligatorio**: cambios en la matriz de permisos se cubren en `permissions_test.go`.
- **Recomendado**: repos contra SQLite en memoria o mock cuando la query tenga lógica (filtros, joins).
- Los tests corren en CI junto a `gofmt`, `go vet`, `go build`.

---

## 7. Cómo se documenta

- **Swagger**: todo endpoint nuevo lleva anotaciones `godoc` (`@Summary`, `@Tags`, `@Param`, `@Success/@Failure`, `@Security BearerAuth`, `@Router`) siguiendo el estilo de los handlers existentes. Tras añadirlas: `make swagger` para regenerar `docs/`.
- **Bitácora** (`docs/md/bitacora/BITACORA.md`): una entrada por feature o decisión técnica, formato `## AAAA-MM-DD — Título` con **Qué se hizo / Pendientes / Referencia vigente**, entradas nuevas al inicio.
- **Docs técnicos vigentes** por módulo en `docs/md/tecnico/`.

---

## 8. Proceso de tarea: HU → HT (entrada) y validación (salida)

### 8.1 Entrada — traducir HU a HT (skill `hu-a-ht`)
Dada una Historia de Usuario (lenguaje funcional/común), producir una **Historia Técnica** que incluya:
- **Resumen funcional** y criterios de aceptación.
- **Capas y archivos** a crear/modificar (según §3.5).
- **Estándares aplicables**: multi-tenant (§4), permisos/entitlements (§5), DTOs/validación (§2.4/§3.4), errores (§2.3).
- **Modelo de datos**: entidades, relaciones, índices, migración (`AllModels`).
- **Plan de pruebas** (§6) y **documentación** requerida (§7).
- **Riesgos** (especialmente fuga cross-tenant) y decisiones abiertas a confirmar.

### 8.2 Salida — validar cumplimiento (skill `validar-estandares`)
Al terminar, validar el diff contra la **Definición de Done** (§9) y reportar incumplimientos con ubicación y corrección sugerida.

---

## 9. Definición de Done (checklist de salida)

- [ ] `gofmt`, `go vet ./...`, `go build ./...` limpios.
- [ ] Capas respetadas (§3.1); sin DB en handlers ni HTTP en services/repos.
- [ ] Multi-tenant completo (§4): scoping en list/get, `company_id` forzado en create, pertenencia verificada en get/update/delete, integridad cruzada en el service.
- [ ] Autorización correcta (§5): `RequirePermission` y, si aplica, `RequireFeature`.
- [ ] Modelo con `BaseModel` + `TableName()` + en `AllModels` si es nuevo.
- [ ] DTOs de respuesta (no modelo crudo) y validación con `binding`.
- [ ] Envelope de respuesta y códigos HTTP correctos (§2.6).
- [ ] Tests de la lógica de negocio nueva (§6).
- [ ] Anotaciones Swagger + `make swagger` (§7).
- [ ] Entrada en la bitácora (§7).

---

## 10. Roadmap de evolución (futuro)

Cuando la app crezca y se migre de esta arquitectura simple a una más estable/escalable:
- **Segmentar HU/HT por módulo** (carpetas por dominio) para organizar el backlog técnico.
- Centralizar el scoping multi-tenant (cerrar la brecha del §4) y considerar middleware de tenancy a nivel de query (GORM scopes/callbacks).
- Evaluar separación por dominios/bounded contexts, capa de aplicación explícita, y observabilidad (logs estructurados, métricas, trazas).
- Estos cambios se registran como decisiones en la bitácora antes de ejecutarse.
