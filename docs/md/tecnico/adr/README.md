# Registros de decisiones arquitectónicas

Esta carpeta conserva las decisiones que definen la arquitectura de DVRA. Un
ADR explica el contexto, la decisión y sus consecuencias; no sustituye Swagger,
los manuales técnicos ni la bitácora de implementaciones.

## Índice

| ADR | Decisión | Fecha | Estado |
|---|---|---|---|
| [ADR-001](ADR-001-arquitectura-modular.md) | Monolito modular con hexagonal-lite | 2026-06-20 | Aceptado |
| [ADR-002](ADR-002-contextos-landlord-tenant.md) | Separación de contextos landlord y tenant | 2026-07-15 | Aceptado |
| [ADR-003](ADR-003-catalogos-modulares-y-suscripciones.md) | Catálogos modulares y suscripciones | 2026-07-20 | Aceptado |

## Estados posibles

- **Propuesto:** todavía está en evaluación.
- **Aceptado:** es la decisión vigente.
- **Rechazado:** fue evaluado, pero no adoptado.
- **Reemplazado:** otro ADR posterior ocupa su lugar.
- **Obsoleto:** ya no aplica y no tiene reemplazo directo.

Los ADR aceptados no se reescriben para ocultar decisiones anteriores. Cuando
una decisión cambia, se crea un ADR nuevo y se enlazan ambos documentos mediante
los campos `Reemplaza` y `Reemplazado por`.
| [ADR-004](ADR-004-plantilla-de-modulo.md) | Plantilla de módulo y fronteras entre dominios | Aceptado | 2026-08-28 |
