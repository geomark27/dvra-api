# Historias Técnicas (HT)

Traducciones de Historias de Usuario a plan técnico, producidas por la skill `hu-a-ht` antes de
programar. Son la **puerta de entrada** del ciclo de tarea; la de salida es `validar-estandares`.

## Convención

- Un archivo por historia: `HT-<slug>.md` (ej. `HT-limites-de-plan.md`).
- El slug describe la funcionalidad, no el ticket.
- La estructura de cada HT la define `../tecnico/ESTANDARES_DESARROLLO.md` §8.1.

## Qué debe llevar siempre

Toda HT empieza identificando **a qué módulo pertenece** y consultando el mapa de migración
(`ESTANDARES_DESARROLLO.md` §3.7), porque de ahí sale la ruta real de los archivos: módulo si el
dominio ya está migrado o es nuevo, Anexo A si sigue pendiente.

Una HT no se da por cerrada con decisiones abiertas que cambien el modelo de datos o el contrato
HTTP: esas se confirman antes de codear.
