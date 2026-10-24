# 002 — Makefile con el flujo de ramas y PR

**Fecha:** 2026-10-24
**Rama:** docs/fase-0-plan → develop
**Fase del roadmap:** 0

## 1. Solicitud

Mejorar el `Makefile` con un flujo de git como el de `draco-hq` (rama, push, PR, merge, cierre),
en la misma rama que el plan de la fase 0.

## 2. Motivación

| Problema | Daño real |
|---|---|
| `make push m=...` y `make sync m=...` hacían `git add .` + `git commit` + `git push` en un solo paso | Se pueden subir archivos no deseados (`.env`, binarios) y el mensaje no pasa por el hook si este no está instalado; el commit «credit card management2222» salió de ese hábito |
| No había forma de crear la rama de la tarea ni el PR desde el `Makefile`, y el último commit se hizo directo en `main` | Se salta la revisión y la regla «una entrada de bitácora por rama» no se puede cumplir |
| `sync` hacía `pull` + commit + push, y `pull` dependía de que la rama ya existiera en origin | Dos objetivos con nombres engañosos |

## 3. Diseño técnico

- **Base `develop`** (variable `BASE`): todo desarrollo entra a `develop` por PR y `develop` llega
  a `main` con `make promote`. `develop` se creó en origin en el mismo commit que `main`
  (`21c3ccd`). Objetivos nuevos, como en draco: `sync`, `start`, `push`, `pr`, `merge`, `finish`,
  `push-pr`, `push-pr-merge`, `merge-finish`, `ship`, `promote` y `prune`. `status` se mantiene.
- **`start name=<scope>-<slug> [type=feat|fix|docs|refactor|chore|test|perf|build|ci]`:** valida
  el nombre, exige árbol limpio, actualiza `develop` con fast-forward y crea `<type>/<name>`. Los
  tipos son los que acepta el hook `commit-msg`; la rama sigue la convención `feat/<scope>-<slug>`
  de los mensajes de error del `Makefile` anterior. `push`, `commit`, `pr` y `merge` rechazan
  estar en `main` o en `develop`.
- **`commit m='tipo(scope): resumen'`:** objetivo nuevo que no existe en draco. Comitea **solo lo
  staged**; `all=1` agrega todo. Alternativa descartada: conservar `git add .` por defecto.
- **`push`** solo sube la rama y fija su upstream (`force=1` usa `--force-with-lease`). Ya no
  commitea.
- **`merge`** (`scripts/merge.sh`) fusiona en `develop` el PR **con merge commit**, tras confirmar (`yes=1` no
  pregunta); conserva los commits convencionales de la rama y no usa `--delete-branch` porque
  `finish` borra la rama local.
- **Se retiran** `pull` (lo cubre `sync`) y el `push`/`sync` con mensaje. Cambio de comportamiento
  deliberado.
- **`promote`** (`scripts/promote.sh`, portado de draco): PR `develop` → `main` con merge commit,
  confirma antes, deja `main` y `develop` locales al día. No se porta `release`: genera binarios
  por plataforma y dvra es un servicio.
- **No se instala el hook automáticamente.** `make hooks-install` sigue siendo un paso aparte
  (actividad 0.1).

## 4. Archivos

- `Makefile` — `SHELL := /bin/bash`, variables `REMOTE`/`BASE`/`type`, bloque de git reescrito,
  `.PHONY` y la ayuda actualizados.
- `scripts/merge.sh` — nuevo: fusión del PR de la rama actual contra `develop`.
- `scripts/promote.sh` — nuevo: `develop` → `main`.

## 5. Cómo verificar

```bash
make help                    # sección Git con los objetivos nuevos
make start name=Mal          # rechaza: nombre en mayúsculas
make start name=x type=zzz   # rechaza: tipo inválido
make start name=x            # rechaza si hay cambios sin commitear
make push BRANCH=develop     # rechaza: no se sube develop
make commit                  # rechaza: falta m=
bash -n scripts/merge.sh scripts/promote.sh
```

Se ejecutaron estas comprobaciones, todas con el resultado indicado. No se ejecutó el flujo
completo (`push`, `pr`, `merge`, `ship`, `finish`, `promote`): requiere subir la rama y abrir un PR reales.

## 6. Pendientes

- Probar el flujo completo con esta misma rama (`make push-pr`) cuando se decida abrir el PR.
- `AGENTS.md`, `ESTANDARES_DESARROLLO.md` y las skills todavía no mencionan `make start` /
  `make ship`: actividad 0.3.
