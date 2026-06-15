---
name: hu-a-ht
description: Traduce una Historia de Usuario (HU) funcional a una Historia Técnica (HT) para el proyecto Dvra, validando los estándares de desarrollo ANTES de programar. Úsala al iniciar cualquier tarea que venga descrita como HU/requerimiento funcional ("tengo esta tarea", "esta HU", "hay que implementar X"). Es la puerta de ENTRADA del ciclo de tarea.
---

# Skill: HU → HT (puerta de entrada)

Convierte un requerimiento funcional (HU) en una Historia Técnica (HT) accionable y alineada con los estándares de Dvra. NO escribe código de la feature; produce el plan técnico que se validará a la salida con la skill `validar-estandares`.

## Procedimiento

1. **Lee la fuente de verdad.** Abre `docs/md/tecnico/ESTANDARES_DESARROLLO.md` completo. Es la base de toda decisión técnica. Si no existe, detente y avísalo.

2. **Obtén la HU.** Tómala de los argumentos o del mensaje del usuario. Si no hay HU o es ambigua (falta el "qué" o el "para quién"), pide la HU o las 1-3 aclaraciones mínimas antes de continuar.

3. **Levanta contexto real del proyecto** (no de memoria): identifica módulos/patrones existentes relevantes a la HU (entidad similar, capas, rutas, permisos, modelo de datos). Usa búsqueda en el repo; cita rutas `archivo:línea`.

4. **Produce la HT** con esta estructura exacta:
   - **Resumen funcional**: qué y para quién, en 2-3 líneas.
   - **Criterios de aceptación**: lista verificable (incluye casos de error).
   - **Capas y archivos** a crear/modificar (según §3.5 del estándar): `models`, `dtos`, `repository`, `service`, `handler`, `routes.go`, `server.go`, `models_all.go`, seeders.
   - **Modelo de datos**: entidades, campos, relaciones, índices, migración (`AllModels`), `BaseModel`+`TableName()`.
   - **Multi-tenancy (§4)**: cómo se aplica la regla de oro `company_id` (list/get/create/update/delete) y qué integridad cruzada hay que validar en el service. Marca explícitamente el riesgo de fuga cross-tenant.
   - **Autorización y entitlements (§5)**: permisos nuevos (`modulo.accion`) + grants por rol; si es módulo premium, flag de plan + `RequireFeature`.
   - **Validación y DTOs (§2.4/§3.4)**: tags `binding`, qué se fuerza desde el token, DTOs de respuesta.
   - **Plan de pruebas (§6)**: qué services/validaciones llevan tests.
   - **Documentación (§7)**: endpoints a anotar en Swagger + entrada de bitácora.
   - **Riesgos y decisiones abiertas**: lo que requiere confirmación del usuario antes de codear.

5. **Cierra con las decisiones abiertas.** Si alguna decisión cambia el modelo de datos o el contrato, pregúntala (usa AskUserQuestion) antes de dar la HT por cerrada.

6. **Ofrece guardar la HT.** Por defecto entrégala en la conversación; ofrece guardarla en `docs/md/historias/HT-<slug>.md` si el usuario quiere trazabilidad. (La segmentación por módulo es evolución futura, no la impongas ahora.)

## Reglas
- No programes la feature en esta skill; solo el plan técnico.
- Toda afirmación sobre el código actual va con ruta `archivo:línea`, no de memoria.
- Si la HU implica datos sensibles reales (tarifas, contratos, datos personales), recuerda usar datos ficticios para diseño/pruebas.
