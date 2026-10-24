# Roadmap — Dvra API

## Visión

Llevar el backend a un estado en el que la arquitectura de `AGENTS.md` **se cumpla y se verifique
sola**, y en el que no queden huecos de seguridad conocidos. Origen: el análisis del proyecto del
2026-10-24 (ver bitácora `fase-0/001-plan-fase-0.md`).

Las funcionalidades de producto (p. ej. la tabla de precios aprobada el 2026-08-26) no forman parte
de este roadmap; van por su propio ciclo HU → HT.

## Estado

| Fase | Contenido | Estado |
|---|---|---|
| 0 | Purga y orden: controles automáticos, documentación, bitácora por fases, residuos de código | 🔄 en curso |
| 1 | Seguridad crítica: lo que hoy expone datos o permite suplantación | ⬜ |
| 2 | Fronteras de arquitectura: que las reglas de `AGENTS.md` se cumplan y se verifiquen | ⬜ |
| 3 | Autenticación y reglas de negocio: endurecimiento | ⬜ |
| 4 | Calidad transversal: respuestas HTTP, logging, dashboard, despliegue | ⬜ |
| 5 | Tests: red de pruebas en lo que no se cubrió antes | ⬜ |

Cada fase abre con una entrada de plan y cierra con una de cierre que verifica su criterio.

---

## Fase 0 — Purga y orden

**Criterio de cierre:** `make check` verde en local y en CI; los hooks instalados; `docs/` sin
archivos obsoletos ni referencias a `internal/app`; bitácora por fases en uso.

| Actividad | Contenido |
|---|---|
| 0.0 Flujo de git en el `Makefile` | `make start/commit/push/pr/merge/finish/ship` sobre `develop` y `make promote` hacia `main` (hecho en la `002`) |
| 0.1 Controles automáticos | Reinstalar `golangci-lint` compatible con Go 1.26; restaurar `.github/workflows/ci.yml` (recuperable de `21c3ccd^`); instalar hooks (`make hooks-install`); trabajar en ramas que nacen de `develop`, no en `develop` ni `main` |
| 0.2 Purga de documentación | Eliminar `docs/01_…` a `04_…`, `tecnico/API_ENDPOINTS.md` y `tecnico/AUDITORIA_ESTANDARES.md`; fusionar `FLUJO_CLIENTE` y `FLUJO_SUPERADMIN`; reducir `ARQUITECTURA`, `MODULO_PLANES` y `MODULO_UBICACIONES`; marcar la «Red Dvra» como futura; decidir el destino de `PLAN_DE_NEGOCIO_ANO1.md`. No se tocan `docs.go`, `swagger.json` ni `swagger.yaml` (los genera `swag`) |
| 0.3 Bitácora por fases | Estructura de carpetas, plantilla y roadmap (hecho en la `001`); actualizar `AGENTS.md`, `ESTANDARES_DESARROLLO.md` §7 y §9, las skills `hu-a-ht` y `validar-estandares`, `docs/md/README.md` y ADR-001; eliminar `docs/md/historias/` (la HT vive en la entrada de la bitácora); retirar la carpeta duplicada de `docs/` raíz |
| 0.4 Residuos de código | `.loom`, `scripts/init-db.sql`, `database.DB` y `GetDB` sin uso, lo que reporte `unused`; rama remota `ADD-1-implementacion-docker` (confirmar antes) |
| 0.5 Documentos que mienten | `README.md` raíz (estructura y enlaces reales); `AGENTS.md` (`ports.go` está en `service/`; la regla 2 hoy no la verifica depguard); entrada retroactiva de la migración modular (`91c653a`, cuyo mensaje «credit card management2222» no corresponde a su contenido) |
| 0.6 Cierre | Entrada de cierre que verifica el criterio |

## Fase 1 — Seguridad crítica

**Criterio de cierre:** cada hallazgo con un test que reproduce el fallo y pasa tras la corrección.

| Actividad | Contenido |
|---|---|
| 1.1 Postulación pública | `ApplyToJob` no modifica datos ni CV de un candidato existente; la postulación lleva su propio CV |
| 1.2 Configuración insegura | Validar secretos JWT en todo entorno que no sea `development` o `test`; rechazar `CORS_ALLOWED_ORIGINS=*` con credenciales |
| 1.3 Exposición de datos | DTO en `GET /users/:id` filtrado a la empresa del token; `GET /plans/:slug` sin planes privados o inactivos; usuarios eliminados fuera de `GET /users` |
| 1.4 Entradas | Solo URLs `http(s)` del candidato; validar `assigned_recruiter` y `hiring_manager` como miembros de la empresa |
| 1.5 Career page | Respetar `CareerPageEnabled` en listado, detalle y `/apply` |

## Fase 2 — Fronteras de arquitectura

**Criterio de cierre:** 0 importaciones del `repository/` o `transport/` de otro módulo, y depguard
lo hace cumplir; 0 usos de GORM en `service/`.

| Actividad | Contenido |
|---|---|
| 2.1 depguard | Regla que prohíba `modules/*/repository` y `modules/*/transport` desde otros módulos |
| 2.2 Puertos | Reemplazar las 6 importaciones cruzadas por puertos declarados en el consumidor |
| 2.3 GORM fuera de `service/` | `auth_service`, `plan_service`, `company_settings_service`, `tenant_rbac_service`, `user_service`, `company_service`, `candidate_service` |
| 2.4 `iam` sobre `authctx` | Reemplazar `IsSuperAdmin` y `c.Get("company_id").(uint)` por los helpers de scope |
| 2.5 DTOs | Dejar de devolver modelos GORM crudos |
| 2.6 Suscripción | Centralizar el alta de `CompanySubscription` en el service de billing |

## Fase 3 — Autenticación y reglas de negocio

**Criterio de cierre:** tests de los caminos de límite, revocación y concurrencia.

| Actividad | Contenido |
|---|---|
| 3.1 Rate limiting | Login, registro, refresh y postulación pública |
| 3.2 Tokens | `jti`, rotación y revocación del refresh token; invalidar al cambiar contraseña y al cerrar sesión; revalidar la membresía de plataforma por petición |
| 3.3 Límites de plan | `MaxUsers` y `MaxCandidates` con guard transaccional |
| 3.4 Enumeración de cuentas | Verificar la contraseña antes que el estado; bcrypt ficticio; respuestas genéricas |
| 3.5 Integridad de datos | Índice único en `Placement.ApplicationID`; índices únicos parciales con soft delete (con migración); claves foráneas entre dominios |

## Fase 4 — Calidad transversal

**Criterio de cierre:** una sola forma de responder, de loguear y de configurar.

| Actividad | Contenido |
|---|---|
| 4.1 Respuestas HTTP | Decidir: implementar el envelope de `ESTANDARES_DESARROLLO` §2.6 o cambiar el estándar |
| 4.2 Logging | `slog` inyectado; sale la dependencia `loom-go`; `LOG_LEVEL` efectivo |
| 4.3 Dashboard | Lógica de negocio al service; consultas en paralelo o agregadas |
| 4.4 Despliegue | `alpine` con soporte, `WriteTimeout` acorde con las subidas, Swagger fuera de producción, `@host` correcto, cierre ordenado en `main` |

## Fase 5 — Tests

**Criterio de cierre:** ≥ 30 % en la capa de servicios y 100 % de los caminos de límites de plan y
ciclo de suscripción (meta de `ESTANDARES_DESARROLLO` §6).

| Actividad | Contenido |
|---|---|
| 5.1 Servicios sin test | `auth_service`, `public_service`, `candidate_service`, `company_service`, `tenant_rbac_service` |
| 5.2 Integración | Repositorios contra el Postgres local de `docker-compose` |
| 5.3 Handlers | Pruebas de `iam/transport` y de los permisos por ruta |

> Cada fase incluye los tests de lo que toca; la fase 5 cubre solo lo que no se cubrió antes.
