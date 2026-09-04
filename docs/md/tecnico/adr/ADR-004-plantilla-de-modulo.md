# ADR-004 — Plantilla de módulo y fronteras entre dominios

- **Estado:** Aceptado
- **Fecha:** 2026-08-28
- **Tipo:** Decisión arquitectónica
- **Complementa:** ADR-001 (monolito modular)

Este ADR recoge las decisiones que la ejecución del ADR-001 obligó a concretar.
El ADR-001 fijó el rumbo; aquí está la forma exacta que tomó.

## Contexto

El ADR-001 describía módulos con `domain/`, `service/`, `repository/`,
`transport/` y `module.go`, con los DTOs compartidos en `internal/app/dtos`. Al
migrar los diez dominios aparecieron tres problemas que esa descripción no
resolvía, y una referencia externa que simplificaba el diseño.

La referencia es `terra-api-purchasing` (Spring Boot, del mismo autor): un
monolito modular con las mitades `modules/` y `platform/`, fronteras verificadas
por análisis estático, y una prohibición explícita de `ports`, `adapters` y una
clase por caso de uso. Su patrón resultó más simple que el hexagonal-lite y
suficiente para una aplicación mayormente CRUD.

## Decisiones

### 1. Tres subpaquetes, no cuatro

```
internal/modules/<dominio>/
├── module.go     # wiring + RegisterRoutes
├── ports.go      # interfaces hacia OTROS módulos (solo si las hay)
├── transport/    # handlers, rutas
├── service/      # lógica de negocio, DTOs, puertos de repositorio
└── repository/   # GORM
```

Se descarta `domain/` como paquete separado. En Go, un paquete por capa
conceptual multiplica los ciclos de importación sin aportar aislamiento real:
`internal/` ya es una barrera del compilador. `staffing` conserva su `domain/`
por ser el piloto y porque ahí viven sus tipos cross-módulo.

### 2. `transport` importa `service`, nunca al revés

De esta regla se deduce dónde va cada tipo, y es lo que evita los ciclos:

- Los tipos que aparecen en la firma del service (sus DTOs de entrada y salida)
  viven en `service/`.
- `transport/` solo guarda lo puramente HTTP.
- El `repository/` referencia los tipos del puerto que implementa, que son del
  service. No es una inversión: quien implementa una interfaz referencia los
  tipos de su firma.

Poner los DTOs siempre en `transport/` es imposible: un service que devuelve sus
DTOs tendría que importar `transport`, que a su vez lo importa a él.

### 3. La interfaz del repositorio la declara el service

No el paquete que la implementa. En Go la satisfacción es estructural, así que
`repository` no necesita importar `service` para satisfacer su puerto. Esto hace
que cada service se pueda probar sin base de datos, que es lo que permitió cubrir
la lógica de negocio de los diez módulos.

### 4. Los módulos que se necesitan mutuamente se construyen en dos fases

`recruitment` necesita los clientes finales de `staffing`, y `staffing` necesita
las postulaciones de `recruitment`. No es un ciclo de importación —cada uno
declara su puerto— pero sí un ciclo de construcción.

Se resuelve construyendo primero lo que solo depende de la base y cableando el
resto después (`BindApplications`). Se prefiere a la inicialización perezosa
porque el fallo, si lo hay, ocurre al arrancar y no en la primera petición.

### 5. Un adaptador cuando los tipos coinciden en forma pero no en nombre

La satisfacción estructural de Go no unifica dos structs con el mismo contenido
y distinto nombre. Cuando dos módulos definen un tipo equivalente
(`ResumeUpload`), el composition root traduce. Es el precio de que ninguno
dependa de los tipos del otro, y se paga en un solo sitio.

### 6. Los modelos GORM viven en `internal/shared/models`

Se mantiene la regla 2 del ADR-001 (modelos compartidos), pero se renombra:
`internal/app/models` mentía sobre una capa que ya no existe.

### 7. Las fronteras se verifican, no se confían

`depguard` en `.golangci.yml` impone: `platform` no importa `modules` (salvo el
composition root), ningún módulo importa `internal/app`, `shared` no importa
`modules`, y `shared/tenant` no importa gorm. Es el equivalente en Go a las
pruebas ArchUnit de la referencia. Los ciclos los prohíbe el compilador gratis.

Además, `TestRouteContract` congela el conjunto de rutas expuestas: mover un
handler de paquete nunca debe cambiarlo.

## Consecuencias

**A favor:** cada dominio se lee y se prueba solo; el composition root es el
único punto que conoce a todos; añadir un módulo cuesta una línea; las reglas
fallan en CI en lugar de erosionarse.

**En contra:** algún adaptador de traducción en el composition root; dos módulos
declaran errores con el mismo mensaje en vez de compartir la variable — es
deliberado, un módulo no toma prestadas las variables de error de otro.

**Lo que sigue abierto:** `PublishJobWithinLimit` y la asignación de plan usan
`*gorm.DB` en el service porque necesitan una transacción con bloqueo de fila.
Meterlas tras un puerto exigiría una abstracción de unidad de trabajo que el
ADR-001 descarta; se acepta como excepción documentada.

## Referencia

Módulo más simple: `internal/modules/locations/`. Módulo con puertos
cross-módulo: `internal/modules/staffing/`. Agregador puro:
`internal/modules/publicapi/`.
