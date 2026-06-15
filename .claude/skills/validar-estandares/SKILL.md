---
name: validar-estandares
description: Valida que el código de una tarea cumpla los estándares de desarrollo de Dvra (Definición de Done) en el proyecto dvra-api. Úsala al TERMINAR una tarea, antes de commit/PR, o cuando el usuario pida "valida los estándares", "revisa si cumple", "está listo para commit". Es la puerta de SALIDA del ciclo de tarea.
---

# Skill: Validar estándares (puerta de salida)

Revisa el cambio actual contra la Definición de Done de Dvra y reporta incumplimientos con ubicación y corrección. Complementa (no reemplaza) a `/code-review` y `/security-review`.

## Procedimiento

1. **Lee la fuente de verdad.** Abre `docs/md/tecnico/ESTANDARES_DESARROLLO.md`, en especial §4 (multi-tenant), §5 (autorización), §9 (Definición de Done).

2. **Determina el alcance del cambio.** `git status` y `git diff` (working tree y/o contra la rama base). Lista archivos tocados por capa.

3. **Checks mecánicos** (ejecútalos y reporta resultado real):
   - `gofmt -l <archivos>` → debe salir vacío.
   - `go vet ./...` → limpio.
   - `go build ./...` → compila.
   - `go test ./...` → pasa (señala si la lógica nueva no tiene tests).

4. **Revisión contra el checklist §9** — por cada ítem marca ✅/❌ con evidencia (`archivo:línea`):
   - Capas respetadas (§3.1): sin `database.DB`/GORM en handlers; sin `gin.Context`/HTTP en services o repos.
   - **Multi-tenant (§4)** — el de mayor riesgo: en cada endpoint nuevo de recurso de tenant verifica scoping en list/get, `company_id` forzado en create, pertenencia verificada en get/update/delete, e integridad cruzada validada en el service. Cualquier query de datos de tenant sin filtro por `company_id` es un hallazgo crítico.
   - Autorización (§5): `RequirePermission` en las rutas y `RequireFeature` si es módulo premium.
   - Modelo: `BaseModel` + `TableName()` + registrado en `AllModels` si es nuevo.
   - DTOs de respuesta (no modelo crudo) y validación con `binding` (no `validate`).
   - Envelope y códigos HTTP (§2.6).
   - Tests de la lógica nueva (§6).
   - Swagger anotado + recordatorio de `make swagger` (§7).
   - Entrada en la bitácora (§7).

5. **Reporta** un resumen claro: lista de ✅ y ❌, ordenando los ❌ por severidad (crítico: fuga cross-tenant / autorización ausente; mayor: capa violada / sin tests; menor: estilo/docs). Para cada ❌ da la corrección concreta.

6. **Ofrece aplicar las correcciones** (`--fix` mental): si el usuario acepta, corrige los hallazgos menores/mayores directamente; los críticos confírmalos antes de tocar.

## Reglas
- No declares "cumple" sin haber corrido los checks mecánicos y citado evidencia.
- Si falta la entrada de bitácora o las anotaciones Swagger, es un ❌ (no opcional).
- Sé específico: cada hallazgo con `archivo:línea` y el fragmento del estándar que incumple.
