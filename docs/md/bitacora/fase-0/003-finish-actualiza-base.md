# 003 — `make finish` actualiza develop antes de cambiar de rama

**Fecha:** 2026-10-24
**Rama:** fix/infra-finish-actualiza-base → develop
**Fase del roadmap:** 0

## 1. Solicitud

Primer uso real del flujo (`make push-pr`, PR #2 fusionado) y `make finish` falló con
«Error: Debes proporcionar un mensaje» (objetivo `sync`, Makefile:351).

## 2. Motivación

| Problema | Daño real |
|---|---|
| `finish` hacía `git switch develop` y después `make sync`. Si el PR cambió el propio `Makefile` (como el #2), el `develop` local aún era el anterior: `make` volvía a leer el archivo del disco y ejecutaba el `sync` viejo, que exige `m=` | El cierre de la tarea falla justo cuando el PR modifica el `Makefile`; hay que actualizar a mano |

## 3. Diseño técnico

- Causa: dos invocaciones de `make` leen archivos distintos porque entre ellas cambia la rama.
- `finish` ahora actualiza `develop` **antes** de cambiar de rama (`git fetch origin develop:develop`;
  si ya se está en `develop`, ese fetch no aplica y se omite), luego hace `git switch develop` y
  `git pull --ff-only`. Solo después llama a `make prune`, que ya lee el `Makefile` nuevo.
- Alternativa descartada: llamar a `git` con rutas absolutas y no usar sub-`make`; más código para
  el mismo resultado.
- Si `develop` no se puede avanzar con fast-forward, `finish` se detiene en vez de seguir.

## 4. Archivos

- `Makefile` — objetivo `finish`.

## 5. Cómo verificar

```bash
make -n finish          # fetch develop:develop → switch → pull --ff-only → make prune
```

Tras fusionar un PR que cambie el `Makefile`, `make finish` debe terminar sin pedir `m=`. Solo se
comprobó la secuencia de comandos con `-n`; el caso real se prueba con el PR de esta rama.

## 6. Pendientes

- Confirmar `make finish` con el PR de esta misma rama.
