# 001 — Plan de la fase 0 y estructura de la bitácora

**Fecha:** 2026-10-24
**Rama:** docs/fase-0-plan → develop
**Fase del roadmap:** 0

## 1. Solicitud

Tras el análisis del proyecto: planificar una purga (lo que no aporta al proyecto, la documentación
sobrante) y las mejoras por fases, y reorganizar la bitácora para que cada fase tenga su carpeta
con las entradas de sus tareas. Estructura tomada de la documentación de `draco-hq` (roadmap con
estado, entradas numeradas por fase, plantilla fija).

## 2. Motivación

| Problema | Daño real |
|---|---|
| La bitácora era un solo archivo de 769 líneas, con entradas nuevas al inicio y la última del 2026-08-28 | No refleja el trabajo posterior (la migración modular de 189 archivos y el borrado del CI no tienen entrada) y no se puede enlazar una tarea concreta |
| No hay roadmap: el análisis dejó ~35 hallazgos sin orden ni criterio de cierre | Se corre el riesgo de corregir al azar y de no saber cuándo una fase está terminada |
| ~11.800 líneas de `.md`, con documentos duplicados u obsoletos que citan `internal/app` | Quien lee la documentación obtiene rutas y precios falsos |
| Sin CI, sin lint ejecutable y sin hooks instalados (el commit `21c3ccd` borró el CI) | Nada impide romper las reglas de `AGENTS.md` al purgar o refactorizar |

## 3. Diseño técnico

- **Roadmap** en `docs/md/roadmap.md`: fases 0 a 5, cada una con actividades y criterio de cierre
  medible. Orden elegido: primero los controles automáticos (0), después lo que expone datos (1),
  y solo entonces las refactorizaciones (2 a 4), para refactorizar con la red puesta.
- **Bitácora por carpetas `fase-N/`, numeración global.** Alternativa descartada: numerar
  `0.2.1` por actividad; se rompe si una tarea cambia de fase. Una tarea grande se parte en
  entradas consecutivas (`003a`, `003b`); alternativa descartada: subcarpetas, sin anexos que las
  justifiquen.
- **La HT vive en la entrada de la bitácora** (secciones Solicitud, Motivación y Diseño técnico),
  así el plan y lo hecho quedan juntos. Por eso `docs/md/historias/` se elimina en la actividad 0.3.
- **`BITACORA.md` congelado** como `00-historico.md`, sin partirlo por hitos: es historia, no se
  edita. Solo se le quitó la sección «Formato para nuevas entradas», que pasa al `README.md`.
- **No se reescribe el historial de git.** `main` y `develop` están publicados; el mensaje del commit `91c653a`
  («credit card management2222») se aclara en una entrada retroactiva (actividad 0.5), no con
  *force-push*.
- **Fuera de alcance:** funcionalidades de producto, como la tabla de precios aprobada el 2026-08-26.

## 4. Archivos

- `docs/md/roadmap.md` — nuevo: visión, estado y actividades de las fases 0 a 5.
- `docs/md/bitacora/README.md` — nuevo: índice, plantilla y reglas.
- `docs/md/bitacora/00-historico.md` — antes `BITACORA.md` (`git mv`), con aviso de congelado y sin
  la sección de formato.
- `docs/md/bitacora/fase-0/001-plan-fase-0.md` — esta entrada.

## 5. Cómo verificar

```bash
git diff --stat develop -- docs/md
git log --follow --oneline docs/md/bitacora/00-historico.md   # conserva el historial
```

Solo hay cambios en `docs/md/`: no se tocó código y no se borró ningún documento.

## 6. Pendientes

- Las referencias a `BITACORA.md` en `AGENTS.md`, `ESTANDARES_DESARROLLO.md` (§7 y §9), las skills
  `hu-a-ht` y `validar-estandares`, `docs/md/README.md` y ADR-001 siguen apuntando al archivo
  antiguo: se actualizan en la actividad 0.3.
- Decisiones abiertas: destino de `PLAN_DE_NEGOCIO_ANO1.md` y si la bitácora de dvra-web adopta la
  misma estructura.
- `check.sh` (como el de draco) que valide que carpeta y fase coincidan: se evalúa en 0.3.
