# Módulo de Ubicaciones (Locations)

> Referencia técnica del módulo de datos geográficos.
> Consolida los antiguos `LOCATION_MODULE.md` y `LOCATION_USAGE.md`.

---

## 1. Descripción General

El módulo de ubicaciones provee datos geográficos organizados jerárquicamente:

```
Region (Continente)
  └── Subregion (Área geográfica)
        └── Country (País)
              └── State (Estado/Provincia)
                    └── City (Ciudad)
```

**Relaciones:**
- **Region** → tiene muchas **Subregions**
- **Subregion** → pertenece a una Region, tiene muchos **Countries**
- **Country** → pertenece a una Subregion (`SubregionID` nullable)
- **State** → pertenece a un Country
- **City** → pertenece a un State

---

## 2. Setup

```bash
# 1. Crear tablas y seeders básicos (~7 segundos)
make fresh

# 2. (Opcional) Carga masiva de ubicaciones — 157,000+ registros (3-5 min)
make db-location
```

**Datos básicos incluidos en `make fresh`:**
- 2 Regions (Americas, Europe)
- 6 Subregions
- ~10 Countries (México, USA, Colombia, España, etc.)
- ~15 States y ~20 Cities principales

**Datos con carga masiva (`make db-location`):**

| Tabla | Registros |
|-------|-----------|
| regions | 6 |
| subregions | 22 |
| countries | 250 |
| states | 5,134 |
| cities | 151,903 |

---

## 3. Arquitectura

```
internal/app/
├── dtos/location_dto.go                    # DTOs de request/response
├── handlers/location_handler.go            # Handlers HTTP (30+ endpoints)
├── repositories/location_repository.go     # Acceso a datos (con preloads)
├── services/location_service.go            # Lógica de negocio y mappers
└── models/
    ├── region.go
    ├── subregion.go
    ├── country.go
    ├── state.go
    └── city.go
```

Todos los modelos usan `gorm.Model` (soft deletes incluido) y campo `IsActive` para deshabilitar sin borrar. Los códigos `iso2` e `iso3` de Country tienen índice único.

---

## 4. Endpoints

### 4.1 Rutas Públicas (GET, sin autenticación)

Base URL: `/api/v1/locations`

```http
# Regions
GET /api/v1/locations/regions
GET /api/v1/locations/regions/:id?include_subregions=true

# Subregions
GET /api/v1/locations/subregions?region_id=1
GET /api/v1/locations/subregions/:id?include_countries=true

# Countries
GET /api/v1/locations/countries?subregion_id=1&search=united
GET /api/v1/locations/countries/:id?include_states=true
GET /api/v1/locations/countries/iso/:iso          # acepta iso2 o iso3

# States
GET /api/v1/locations/states?country_id=1&search=california
GET /api/v1/locations/states/:id?include_cities=true

# Cities
GET /api/v1/locations/cities?state_id=1&search=new
GET /api/v1/locations/cities/:id

# Helpers
GET /api/v1/locations/hierarchy/:countryId        # Jerarquía completa de un país
GET /api/v1/locations/search?q=london             # Búsqueda global en todas las entidades
```

### 4.2 Rutas de Administración (solo SuperAdmin)

```http
POST   /api/v1/admin/locations/{regions|subregions|countries|states|cities}
PUT    /api/v1/admin/locations/{regions|subregions|countries|states|cities}/:id
DELETE /api/v1/admin/locations/{regions|subregions|countries|states|cities}/:id   # soft delete
```

### 4.3 Query Parameters

| Endpoint | Parámetros | Ejemplo |
|----------|-----------|---------|
| `/countries` | `subregion_id`, `search`, `include_states` | `?search=united&include_states=true` |
| `/states` | `country_id`, `search`, `include_cities` | `?country_id=142&search=nuevo` |
| `/cities` | `state_id`, `search` | `?state_id=3456&search=monte` |
| `/search` | `q` | `?q=los` |

**Comportamiento:**
- `search` usa `ILIKE` (case-insensitive, coincidencia parcial).
- Por defecto solo se devuelven registros con `is_active = true`.
- Los preloads (`include_*`) son opcionales para no cargar relaciones innecesarias.

---

## 5. Ejemplos de Request/Response

### Obtener países con filtros

```http
GET /api/v1/locations/countries?subregion_id=5&search=united
```

```json
{
  "data": [
    {
      "id": 233,
      "name": "United States",
      "iso2": "US",
      "iso3": "USA",
      "phone_code": "1",
      "capital": "Washington",
      "currency": "USD",
      "currency_symbol": "$",
      "emoji": "🇺🇸",
      "is_active": true,
      "subregion_id": 5,
      "subregion": { "id": 5, "name": "Northern America", "region_id": 2 }
    }
  ]
}
```

### Jerarquía completa de un país

```http
GET /api/v1/locations/hierarchy/233
```

```json
{
  "data": {
    "id": 233,
    "name": "United States",
    "iso2": "US",
    "subregion": {
      "id": 5,
      "name": "Northern America",
      "region": { "id": 2, "name": "Americas" }
    },
    "states": [
      {
        "id": 1416,
        "name": "California",
        "cities": [{ "id": 111797, "name": "Los Angeles" }]
      }
    ]
  }
}
```

### Crear ciudad (SuperAdmin)

```http
POST /api/v1/admin/locations/cities
Authorization: Bearer <superadmin_token>

{
  "name": "Nueva Ciudad",
  "state_id": 1416,
  "latitude": 34.0522,
  "longitude": -118.2437,
  "is_active": true
}
```

---

## 6. Uso desde Go (interno)

### Consultar con relaciones

```go
// País con su subregión y región
var country models.Country
db.Preload("Subregion.Region").
   Where("iso2 = ?", "MX").
   First(&country)

// Todos los países de una región
var region models.Region
db.Preload("Subregions.Countries").
   Where("name = ?", "Americas").
   First(&region)
```

### Validar país por ISO

```go
func ValidateCountryISO(iso2 string) (bool, error) {
    var count int64
    err := db.Model(&models.Country{}).
        Where("iso2 = ? AND is_active = ?", iso2, true).
        Count(&count).Error
    return count > 0, err
}
```

### Asociar país a una empresa

```go
var country models.Country
if err := db.First(&country, registration.CountryID).Error; err != nil {
    return errors.New("país no válido")
}
company := models.Company{
    Name:      registration.Name,
    CountryID: &registration.CountryID,
}
return db.Create(&company).Error
```

---

## 7. Integración Frontend

### Select en cascada (patrón recomendado)

```javascript
// 1. Cargar países al montar el componente
const countries = await fetch('/api/v1/locations/countries');

// 2. Cargar estados cuando se selecciona país
const states = await fetch(`/api/v1/locations/states?country_id=${countryId}`);

// 3. Cargar ciudades cuando se selecciona estado
const cities = await fetch(`/api/v1/locations/cities?state_id=${stateId}`);
```

### Recomendaciones

1. **Cachear países** (localStorage): cambian muy poco.
2. **Debounce de 300ms** en búsquedas/autocomplete.
3. **Preload selectivo**: usar `include_states=true` solo si se necesita.
4. **Manejar países sin estados/ciudades**: algunos no los tienen.
5. **Loading states y error handling** en cada nivel del select.
6. El campo `emoji` trae la bandera del país lista para mostrar.

---

## 8. Performance

| Operación | Tiempo aproximado |
|-----------|-------------------|
| Listar todos los países | ~50ms |
| Estados de un país | ~10ms |
| Ciudades de un estado | ~20-100ms |
| Jerarquía completa | ~150ms |
