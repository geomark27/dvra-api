# Auditoría de Estándares — Dvra API

> Evaluación del estado del código frente a `ESTANDARES_DESARROLLO.md`.
> **Fecha original:** 2026-06-14 · **Base:** commit del fundamento de estándares (`b1b3a96`).
> **Actualización:** 2026-06-20 — correcciones al módulo staffing/placement (ver bitácora).
> Esta es una foto en el tiempo; los números pueden cambiar con cada release.

---

## Resumen ejecutivo

| Área | Estado | Nota |
|---|---|---|
| Formato / análisis estático | 🟢 | `gofmt`, `vet`, `golangci-lint`, build y tests pasan |
| Arquitectura por capas | 🟢 | Consistente: handler → service → repository con interfaces y DI manual |
| Convenciones documentadas | 🟢 | `BaseModel`, `authctx`, `apperr`, matriz de permisos, bitácora |
| Códigos HTTP en errores | 🟢 | Módulo staffing/placement corregido (2026-06-20); patrón `apperr` disponible para el resto |
| **Aislamiento multi-tenant** | 🔴 | Correcto en concepto pero **manual y repetido**; alto riesgo de fuga por omisión. Módulo staffing mejorado (companyID en service), pero el resto sigue sin centralizar |
| **Validación de entrada** | 🟠 | 47 tags `validate` en 6 DTOs **no se ejecutan**; reglas silenciosamente inactivas |
| **Cobertura de tests** | 🟠 | 1 archivo de test para 16 services |
| Consistencia (contexto/repos) | 🟡 | Dos patrones coexisten (acceso a contexto y a BD) |
| Sincronía de docs generados | 🟡 | Swagger se regenera a mano; sin verificación en CI |

🔴 crítico · 🟠 alto · 🟡 medio · 🟢 conforme

---

## Hallazgos priorizados

### 🔴 C1 — Aislamiento multi-tenant manual y repetido
- **Evidencia:** 41 usos de `IsSuperAdmin` en handlers; el bloque de scoping por `company_id` se repite por endpoint (`application_handler.go` 21 ocurrencias, `job_handler.go` 15, `candidate_handler.go` 13, `membership_handler.go` 11, …). `internal/database/scoped_db.go` (`ScopedDB`) existe pero casi no se usa.
- **Riesgo:** cada copia es una oportunidad de olvido → **fuga de datos entre tenants** (un usuario de la empresa A ve/edita datos de la empresa B). Es un riesgo de seguridad, no de estilo.
- **Estándar:** §4 (regla de oro).
- **Recomendación:** centralizar el scoping en un helper/middleware reutilizable (o adoptar `ScopedDB`/GORM scopes), de modo que el filtrado por `company_id` sea el **default** y no algo que cada handler deba recordar. → **Fase 4 del roadmap de estándares**.

### 🟠 A1 — Validaciones `validate` que nunca se ejecutan
- **Evidencia:** 47 tags `validate:"..."` en 6 DTOs (`user_dto.go`, `company_dto.go`, `job_dto.go`, `candidate_dto.go`, `membership_dto.go`, `application_dto.go`). Gin con `ShouldBindJSON/Query` solo procesa tags **`binding`**; las `validate` se ignoran porque no hay validador go-playground conectado.
- **Riesgo:** reglas que se creen activas no lo están. Ej.: `CreateJobDTO.Title` declara `validate:"required,min=3,max=255"` pero **no** es `binding:"required"`, así que se puede crear un job con título vacío o fuera de rango; emails y `oneof` tampoco se validan al bindear.
- **Estándar:** §2.4.
- **Recomendación:** decidir UNA vía y aplicarla a todos los DTOs: (a) migrar `validate` → `binding` (más simple, sin dependencias), o (b) conectar `go-playground/validator` y ejecutarlo explícitamente. Documentar la elección en §2.4.

### 🟠 A2 — Cobertura de tests casi nula
- **Evidencia:** 16 archivos de service, **1** archivo de test en todo el repo (`permissions_test.go`).
- **Riesgo:** sin red de seguridad para la lógica de negocio (validaciones de integridad, scoping, transiciones de estado). Refactors —incluido el de C1— se vuelven arriesgados.
- **Estándar:** §6.
- **Recomendación:** priorizar tests de services con validaciones: `placement_service` (etapa `hired` + integridad cross-tenant), `job_service` (validación de staffing client), `plan_service` (`CompanyHasFeature`), `auth_service`. Prerrequisito práctico: ver M2.

### 🟡 M1 — Acceso directo al contexto en vez de `authctx`
- **Evidencia:** 7 handlers leen `c.Get("company_id")` directo en lugar de `authctx.CompanyID(c)` (11 ya usan `authctx`). Ya estaba registrado como pendiente en la bitácora (2026-06-12).
- **Riesgo:** type assertion manual (`.(uint)`) propensa a panic, y se evade el único punto tipado de acceso al contexto.
- **Estándar:** §2 / §4.
- **Recomendación:** refactor mecánico a `authctx.*`. Conviene hacerlo junto con C1 (ambos tocan los mismos handlers).

### 🟡 M2 — Dos patrones de repositorio
- **Evidencia:** 10 repos usan el global `database.DB`; 3 reciben `db *gorm.DB` inyectado.
- **Riesgo:** inconsistencia + el global **no es mockeable**, lo que bloquea tests unitarios de repos/services que dependen de él (alimenta A2).
- **Estándar:** §3.2.
- **Recomendación:** unificar hacia `db` inyectado. Habilita tests con SQLite en memoria. Migración incremental, repo por repo.

### 🟡 B1 — Docs generados sin verificación de sincronía
- **Evidencia:** Swagger (`docs/`) se regenera a mano con `make swagger`; nada en CI verifica que esté al día respecto a las anotaciones.
- **Riesgo:** la doc publicada se desincroniza del código sin que nadie lo note.
- **Estándar:** §7.
- **Recomendación:** paso opcional en CI que corra `swag init` y falle si `git diff` en `docs/` no está limpio.

---

## Lo que ya cumple bien (no tocar)
- **Arquitectura por capas** consistente en todos los módulos, con interfaces y DI.
- **Convenciones explícitas**: `BaseModel`, `authctx`, `apperr`, matriz de permisos por módulo, separación DTO/modelo.
- **Multi-tenant conceptualmente correcto**: el modelo de datos y el diseño son sólidos; falta *centralizar* la aplicación (C1), no rediseñar.
- **Bitácora** de decisiones al día.
- **Enforcement automático** recién incorporado (linter + CI + pre-commit) con baseline verde.
- **Módulo staffing/placement** corregido (2026-06-20): errores HTTP correctos, sin doble lookup, sin código muerto en repositorios, modelos de BD con constraints reales, DTOs sin filtración de modelos crudos.

---

## Plan de remediación sugerido (por impacto/riesgo)

1. **A1 — Validación** (acotado, alto valor): unificar `validate`→`binding`. Cierra un hueco de datos silencioso.
2. **C1 — Centralizar scoping** (Fase 4): el de mayor riesgo; idealmente estrenando las skills `hu-a-ht` → `validar-estandares`.
3. **M1 — `authctx`** : refactor mecánico, hacerlo junto con C1.
4. **M2 — Inyección de `db`** : desbloquea tests.
5. **A2 — Tests** : empezar por los services con validaciones, una vez M2 avance.
6. **B1 — Check de Swagger en CI** : barato, evita desincronía.

> Nota: B2 (deuda de `gofmt` y 4 `errcheck`) ya se cerró durante la incorporación del enforcement.
