---
name: validar-estandares
description: Valida que el código de una tarea cumpla los estándares de desarrollo de Dvra (Definición de Done) en el proyecto dvra-api. Úsala al TERMINAR una tarea, antes de commit/PR, o cuando el usuario pida "valida los estándares", "revisa si cumple", "está listo para commit". Es la puerta de SALIDA del ciclo de tarea.
---

# Skill: Validar estándares (puerta de salida)

Revisa el cambio actual contra la Definición de Done de Dvra y reporta incumplimientos con ubicación y corrección. Complementa (no reemplaza) a `/code-review` y `/security-review`.

## Procedimiento

1. **Lee la fuente de verdad.** Abre `docs/md/tecnico/ESTANDARES_DESARROLLO.md`, en especial §3.3 (fronteras), §3.7 (mapa de migración), §4 (multi-tenant), §5 (autorización), §9 (Definición de Done).

2. **Determina el alcance del cambio.** `git status` y `git diff` (working tree y/o contra la rama base). Lista archivos tocados por módulo y capa.

3. **Checks mecánicos** (ejecútalos y reporta resultado real):
   - `make check` → gofmt, vet, lint (incluye `depguard`), test y build. Si no está disponible, ejecuta los pasos por separado.
   - `gofmt -l <archivos>` → debe salir vacío.
   - `go vet ./...` → limpio.
   - `go build ./...` → compila.
   - `go test ./...` → pasa (señala si la lógica nueva no tiene tests).

4. **Checks de frontera** (grep, reporta rutas reales):
   - **Ubicación correcta**: el código nuevo está donde dice el mapa de migración (§3.7). Un dominio nuevo en `internal/app/` es un ❌ mayor.
   - **Sin imports cruzados**: `grep -rn "internal/modules/" internal/modules/<X>/` no debe traer `repository/` ni `transport/` de otro módulo.
   - **`platform` no importa `modules`**: `grep -rn "dvra-api/internal/modules" internal/platform/` debe salir vacío.
   - **Sin estado global**: `database.DB` no aparece en código nuevo.
   - **Sin literales de contexto**: `c.Get("` con string literal fuera de `internal/shared/` es ❌.
   - **Sin `IsSuperAdmin` nuevo** fuera de `internal/shared/`: el scoping va con `authctx.RequireScope` + `scope.Guard`.
   - **Sin tags `validate:`**: el repo está limpio de ellas y debe seguir así.

5. **Revisión contra el checklist §9** — por cada ítem marca ✅/❌ con evidencia (`archivo:línea`):
   - Capas respetadas (§3.2): sin GORM en `transport`; sin `gin.Context`/HTTP en `service` o `repository`.
   - **Multi-tenant (§4)** — el de mayor riesgo: en cada endpoint nuevo de recurso de tenant verifica scoping en list/get, `company_id` forzado en create, pertenencia verificada en get/update/delete, e integridad cruzada validada en el service. Cualquier query de datos de tenant sin filtro por `company_id` es un hallazgo crítico.
   - Autorización (§5): `RequirePermission` en las rutas y `RequireFeature` si es módulo premium.
   - Modelo: `BaseModel` + `TableName()` + registrado en `AllModels` si es nuevo. Sin FK nuevas entre dominios.
   - DTOs de respuesta en el `transport/` del módulo (no modelo crudo, no DTOs compartidos entre módulos) y validación con `binding`.
   - Envelope, códigos HTTP y paginación (§2.6), emitidos vía `httpx` y no con `gin.H` a mano.
   - Tests de la lógica nueva **+ test de aislamiento tenant** del módulo (§6).
   - Swagger anotado con DTO concreto (no `map[string]interface{}`) + recordatorio de `make swagger` (§7).
   - Entrada en la bitácora (§7), con campo `Contraparte` si el cambio toca también el frontend.

6. **Reporta** un resumen claro: lista de ✅ y ❌, ordenando los ❌ por severidad (crítico: fuga cross-tenant / autorización ausente; mayor: frontera de módulo violada / capa violada / sin tests; menor: estilo/docs). Para cada ❌ da la corrección concreta.

7. **Ofrece aplicar las correcciones**: si el usuario acepta, corrige los hallazgos menores/mayores directamente; los críticos confírmalos antes de tocar.

## Reglas
- No declares "cumple" sin haber corrido los checks mecánicos y citado evidencia.
- Si falta la entrada de bitácora o las anotaciones Swagger, es un ❌ (no opcional).
- Sé específico: cada hallazgo con `archivo:línea` y el fragmento del estándar que incumple.
- Si el cambio migra un dominio, verifica que el **mapa de migración (§3.7) quedó actualizado**: no hacerlo es ❌ mayor.
