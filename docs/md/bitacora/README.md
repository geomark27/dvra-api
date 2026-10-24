# Bitácora — Dvra API

Una entrada por actividad: cada rama (`feature/`, `fix/`, `docs/`, `chore/`) que llega a `develop`
por PR (`develop` llega a `main` con `make promote`). La entrada viaja **en el mismo commit** que el cambio que describe y vive en la carpeta de
su fase del roadmap (`fase-N/`).

La historia y el porqué de cada decisión van aquí, no en comentarios del código. Si una decisión
condiciona código futuro, además se registra como ADR (`../tecnico/adr/`): la bitácora dice *qué se
hizo*; el ADR dice *por qué se decidió así y qué se descartó*.

Roadmap y criterio de cierre de cada fase: [`../roadmap.md`](../roadmap.md).

## Histórico

[`00-historico.md`](./00-historico.md): bitácora anterior, de 2025-12 a 2026-08-28. Está congelada.

## Índice

La numeración es **global y correlativa**: no se reinicia en cada fase, porque las entradas se
citan por número. Una tarea grande se parte en varias entradas consecutivas (`003a`, `003b`) que
se enlazan entre sí; no se usan subcarpetas.

### Fase 0 — Purga y orden

| # | Actividad | Rama | Fecha | Resumen |
|---|---|---|---|---|
| [001](./fase-0/001-plan-fase-0.md) | Plan de la fase 0 y estructura de la bitácora | `docs/fase-0-plan` | 2026-10-24 | Roadmap por fases, bitácora por carpetas, histórico congelado |
| [002](./fase-0/002-makefile-flujo-git.md) | Makefile con el flujo de ramas y PR | `docs/fase-0-plan` | 2026-10-24 | `make start/commit/push/pr/merge/finish/ship` sobre `develop`, más `make promote`; `push` ya no hace `git add .` |

## Formato de cada entrada

Archivo `fase-N/NNN-<tema>.md`: `N` es la fase del roadmap, `NNN` el número correlativo global de
tres dígitos y el tema, el de la rama. La primera entrada de una fase crea su carpeta y su sección
en el índice. Encabezado y secciones obligatorias:

```markdown
# NNN — Título corto

**Fecha:** AAAA-MM-DD
**Rama:** feature/<tema> → develop
**Fase del roadmap:** N   (la misma N de la carpeta)

## 1. Solicitud
Qué se pidió, en lenguaje natural (la HU o el hallazgo que lo origina).

## 2. Motivación
| Problema | Daño real |
|---|---|
| Qué está mal | Qué pasa si no se corrige |

## 3. Diseño técnico
Módulo(s) afectado(s), decisiones tomadas y sus razones, incluidas las alternativas descartadas.

## 4. Archivos
Cada archivo tocado con una línea sobre su cambio.

## 5. Cómo verificar
Comandos concretos (`make check`, tests, `curl`).

## 6. Pendientes
Lo que quedó fuera o se descubrió en el camino. Omitir si no hay.

## 7. Contraparte
Entrada equivalente en dvra-web. Omitir si el cambio no toca el frontend.
```

## Reglas

- La fase del encabezado y la carpeta coinciden.
- Una entrada de una actividad ya fusionada en `develop` **no se reescribe**. Solo admite notas
  aclaratorias al final, con fecha. Lo nuevo va en la entrada de la actividad nueva.
- Cada fase abre con una entrada de **plan** y cierra con una de **cierre** que verifica su
  criterio de cierre.
- Sin contraseñas, tokens, URLs de instancias ni datos reales de personas o empresas.
