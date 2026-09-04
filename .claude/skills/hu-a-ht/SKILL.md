---
name: hu-a-ht
description: Traduce una Historia de Usuario (HU) funcional a una Historia Técnica (HT) para el proyecto Dvra, validando los estándares de desarrollo ANTES de programar. Úsala al iniciar cualquier tarea que venga descrita como HU/requerimiento funcional ("tengo esta tarea", "esta HU", "hay que implementar X"). Es la puerta de ENTRADA del ciclo de tarea.
---

# Skill: HU → HT (puerta de entrada)

Convierte un requerimiento funcional (HU) en una Historia Técnica (HT) accionable y alineada con los estándares de Dvra. NO escribe código de la feature; produce el plan técnico que se validará a la salida con la skill `validar-estandares`.

## Procedimiento

1. **Lee la fuente de verdad.** Abre `docs/md/tecnico/ESTANDARES_DESARROLLO.md` completo. Es la base de toda decisión técnica. Si no existe, detente y avísalo. (`AGENTS.md` en la raíz es el resumen operativo de las reglas de modularidad.)

2. **Obtén la HU.** Tómala de los argumentos o del mensaje del usuario. Si no hay HU o es ambigua (falta el "qué" o el "para quién"), pide la HU o las 1-3 aclaraciones mínimas antes de continuar.

3. **Determina el módulo y su estado.** *Antes que nada:* ¿a qué dominio pertenece la HU — `iam`, `recruitment`, `staffing`, `billing`, `catalog`, `locations`, `dashboard`, `public`, o uno nuevo? Consulta el **mapa de migración (§3.7)** y decide dónde vive el código:
   - Dominio **migrado** o **dominio nuevo** → `internal/modules/<dominio>/{transport,service,repository}` + `module.go` + `ports.go` (§3.2).
   - Dominio **pendiente** → estructura por capas del **Anexo A** (`internal/app/...`).

   Si la HU cruza dos dominios, di explícitamente cuál es el dueño del dato y qué necesita el otro; eso será un puerto en `ports.go` del consumidor (§3.3 regla 3), nunca un import del `repository/` ajeno.

4. **Levanta contexto real del proyecto** (no de memoria): identifica módulos/patrones existentes relevantes a la HU (entidad similar, rutas, permisos, modelo de datos). Usa búsqueda en el repo; cita rutas `archivo:línea`. La referencia canónica de un módulo bien formado es `internal/modules/staffing/`.

5. **Produce la HT** con esta estructura exacta:
   - **Módulo y estado**: a qué dominio pertenece y qué estructura le corresponde según §3.7.
   - **Resumen funcional**: qué y para quién, en 2-3 líneas.
   - **Criterios de aceptación**: lista verificable (incluye casos de error).
   - **Archivos** a crear/modificar, con la ruta real que corresponda (§3.2 o Anexo A).
   - **Modelo de datos**: entidades, campos, relaciones, índices, migración (`AllModels`), `BaseModel`+`TableName()`. Sin claves foráneas nuevas entre dominios (§3.3 regla 4).
   - **Fronteras**: qué necesita este módulo de otros y qué puerto lo cubre; qué expone a los demás.
   - **Multi-tenancy (§4)**: cómo se aplica la regla de oro `company_id` (list/get/create/update/delete) y qué integridad cruzada hay que validar en el service. Indica los helpers a usar (`authctx.RequireScope`, `scope.Guard`). Marca explícitamente el riesgo de fuga cross-tenant.
   - **Autorización y entitlements (§5)**: permisos nuevos (`modulo.recurso.accion`) + grants por rol; si es módulo premium, flag de plan + `RequireFeature`.
   - **Contrato HTTP (§2.6)**: envelope, códigos, si la colección se pagina o no.
   - **Validación y DTOs (§2.4/§3.5)**: tags `binding` (nunca `validate`), qué se fuerza desde el token, DTOs de respuesta en el `transport/` del módulo.
   - **Plan de pruebas (§6)**: qué services/validaciones llevan tests + el test de aislamiento tenant obligatorio.
   - **Documentación (§7)**: endpoints a anotar en Swagger (con DTO concreto) + entrada de bitácora + ADR si la decisión condiciona código futuro.
   - **Riesgos y decisiones abiertas**: lo que requiere confirmación del usuario antes de codear.

6. **Cierra con las decisiones abiertas.** Si alguna decisión cambia el modelo de datos o el contrato, pregúntala (usa AskUserQuestion) antes de dar la HT por cerrada.

7. **Guarda la HT** en `docs/md/historias/HT-<slug>.md` (crea la carpeta si no existe) y entrégala también en la conversación.

## Reglas
- No programes la feature en esta skill; solo el plan técnico.
- Toda afirmación sobre el código actual va con ruta `archivo:línea`, no de memoria.
- **Nunca propongas añadir un dominio nuevo a `internal/app/`**: los dominios nuevos nacen siempre como módulo.
- No propongas `ports/`, `adapters/` ni una clase por caso de uso (§3.3 regla 5).
- Si la HU implica datos sensibles reales (tarifas, contratos, datos personales), recuerda usar datos ficticios para diseño/pruebas.
