# ADR-001 — Monolito modular con hexagonal-lite

> **Estado:** Aceptado · **Fecha:** 2026-06-20 · **Piloto:** módulo `staffing`
> Architecture Decision Record. Registra la decisión, no la implementación detallada.

## Contexto

El proyecto está organizado **por capa técnica**: `internal/app/{models,dtos,repositories,services,handlers}`, con interfaces e inyección de dependencias manual. Funciona, pero a medida que crece:

- Cada dominio (auth, recruitment, staffing, billing…) está esparcido en 5 carpetas.
- No hay fronteras que impidan que un módulo toque las tablas/servicios de otro.
- La lógica de negocio está acoplada a infraestructura: los modelos GORM **son** las entidades, y `plan_service` importa `gorm` directamente.

Aún no hay frontend, pero esto es **interno** — el frontend nunca ve la estructura de paquetes Go. La decisión es de mantenibilidad, no de contrato de API.

## Decisión

Adoptar un **monolito modular** organizado por dominio, con disciplina **hexagonal-lite** (puertos y adaptadores) dentro de cada módulo, **sin** caer en la ceremonia de Clean Architecture textbook.

```
internal/modules/<dominio>/
  domain/       # puertos (interfaces) + tipos cross-módulo. No importa gin ni gorm.
  service/      # casos de uso. Depende solo de domain.
  repository/   # adaptador de salida (impl GORM, *gorm.DB inyectado).
  transport/    # adaptador de entrada (handlers HTTP + rutas).
  module.go     # wiring del módulo (paquete raíz).
```

Dominios objetivo: `iam` (auth, users, companies, memberships), `recruitment` (jobs, candidates, applications), `staffing`, `billing` (plans), `platform` (locations, system_values, settings, dashboard), `public`.

### Reglas de la decisión

1. **Dirección de dependencias hacia `domain`.** `transport` y `repository` conocen `domain`; `domain` no conoce a nadie. `gin` y `gorm` viven solo en los adaptadores.
2. **Los modelos se quedan compartidos (1ª pasada).** Las entidades GORM (`models.*`) permanecen en `internal/app/models`. Separarlas por módulo provoca **ciclos de importación** (p. ej. `StaffingClient ↔ Job`, `Placement ↔ Application`) porque las relaciones GORM se referencian mutuamente. Se difiere hasta que un módulo gane lógica que lo justifique (entonces se reemplazan punteros cross-módulo por IDs).
3. **Dependencias entre módulos vía puertos definidos por el consumidor.** Un módulo NO importa a otro. Define la interfaz mínima que necesita y el *composition root* le inyecta un adaptador. Ejemplo: `staffing` necesita datos de una `Application` (de `recruitment`); define `domain.ApplicationFinder` y `server.go` le pasa un adaptador sobre el repo de recruitment. El grafo queda **acíclico**.
4. **Sin estado global.** Los repos reciben `*gorm.DB` inyectado (mockeable), no el global `database.DB`.
5. **DTOs compartidos (lite).** Los DTOs de request/response viven en `internal/app/dtos` y los usan tanto `service` como `transport`. Es un compromiso "lite"; se pueden internalizar como input structs por módulo si un módulo lo amerita.

## Consecuencias

**A favor:** fronteras de módulo claras; lógica de negocio testeable sin BD ni HTTP; imposible saltarse el service de otro módulo; base para extraer microservicios en el futuro si alguna vez se necesita.

**En contra:** un poco más de wiring en el *composition root*; los puertos cross-módulo requieren un adaptador de mapeo (p. ej. `ApplicationFinder`).

**Lo que NO hacemos (y por qué):** Clean Architecture textbook (entidad pura + modelo de persistencia + mappers por todos lados). Para una app mayormente CRUD sobre GORM y equipo chico, el boilerplate no compensa. Go favorece lo pragmático.

## Secuencia

1. Piloto: migrar `staffing` (este ADR lo acompaña). ← **hecho**
2. Establecer red de tests (inyección de `db` ya queda lista en el módulo migrado).
3. Replicar el patrón al resto de dominios, uno por uno, con build/tests verdes entre cada uno.

## Referencia

Piloto implementado en `internal/modules/staffing/`. Ver `BITACORA.md` (entrada 2026-06-20).
