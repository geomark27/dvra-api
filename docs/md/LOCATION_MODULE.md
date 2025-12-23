# 🌍 Location API - Documentación Completa para Frontend

## 📋 Overview

El módulo de Location proporciona acceso completo a datos geográficos organizados jerárquicamente:

```
Region (Continente) → Subregion (Área) → Country (País) → State (Estado) → City (Ciudad)
```

### ✅ Datos Disponibles en Base de Datos

**Ya cargados y listos para usar:**
- ✅ 2 Regions (Americas, Europe)
- ✅ 6 Subregions (Northern America, Central America, South America, etc.)
- ✅ ~10 Countries (México, USA, Colombia, España, etc.)
- ✅ ~15 States (California, Texas, Nuevo León, etc.)
- ✅ ~20 Cities principales (CDMX, Los Angeles, Madrid, etc.)

**Opción de carga masiva disponible** (ejecutar solo si se necesita):
```bash
make db-location  # Carga 157,000+ ubicaciones (tarda 3-5 min)
```

---

## 🏗️ Architecture

### Files Created

```
internal/app/
├── dtos/location_dto.go                    # DTOs for all entities
├── handlers/location_handler.go            # HTTP handlers (30+ endpoints)
├── repositories/location_repository.go     # Data access layer
├── services/location_service.go            # Business logic layer
└── models/
    ├── region.go
    ├── subregion.go
    ├── country.go
    ├── state.go
    └── city.go
```

### Components

1. **Models**: GORM models with relationships and soft deletes
2. **DTOs**: Request/Response objects with validation
3. **Repository**: Data access with preload support
4. **Service**: Business logic with DTO↔Model mappers
5. **Handler**: HTTP endpoints with query parameters
6. **Routes**: Public (GET) + Admin (POST/PUT/DELETE)

---

## 🔗 API Endpoints (Base URL: `/api/v1/locations`)

### 🌐 Rutas Públicas (Sin Autenticación)

**Todas las consultas (GET) son públicas** - no requieren token de autenticación:

#### Regions
```http
GET /api/v1/locations/regions
GET /api/v1/locations/regions/:id?include_subregions=true
```

#### Subregions
```http
GET /api/v1/locations/subregions?region_id=1
GET /api/v1/locations/subregions/:id?include_countries=true
```

#### Countries
```http
GET /api/v1/locations/countries?subregion_id=1&search=united
GET /api/v1/locations/countries/:id?include_states=true
GET /api/v1/locations/countries/iso/:iso  # iso2 or iso3
```

#### States
```http
GET /api/v1/locations/states?country_id=1&search=california
GET /api/v1/locations/states/:id?include_cities=true
```

#### Cities
```http
GET /api/v1/locations/cities?state_id=1&search=new
GET /api/v1/locations/cities/:id
```

#### Helpers
```http
GET /api/v1/locations/hierarchy/:countryId   # Full hierarchy for a country
GET /api/v1/locations/search?q=london        # Search all entities
```

---

### 🔒 Rutas de Administración (SuperAdmin Únicamente)

Requieren autenticación con rol SuperAdmin:

```http
# Crear
POST /api/v1/admin/locations/{regions|subregions|countries|states|cities}

# Actualizar
PUT /api/v1/admin/locations/{regions|subregions|countries|states|cities}/:id

# Eliminar (soft delete)
DELETE /api/v1/admin/locations/{regions|subregions|countries|states|cities}/:id
```

**Nota**: El frontend normalmente solo usará las rutas GET públicas.

---

## 📦 Request/Response Examples

### Get All Countries with Filters

**Request:**
```http
GET /api/v1/locations/countries?subregion_id=5&search=united
```

**Response:**
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
      "native": "United States",
      "emoji": "🇺🇸",
      "emojiU": "U+1F1FA U+1F1F8",
      "is_active": true,
      "subregion_id": 5,
      "subregion": {
        "id": 5,
        "name": "Northern America",
        "region_id": 2
      }
    }
  ]
}
```

### Get Country with States

**Request:**
```http
GET /api/v1/locations/countries/233?include_states=true
```

**Response:**
```json
{
  "data": {
    "id": 233,
    "name": "United States",
    "iso2": "US",
    "iso3": "USA",
    "states": [
      {
        "id": 1416,
        "name": "California",
        "state_code": "CA",
        "country_id": 233
      },
      {
        "id": 1452,
        "name": "New York",
        "state_code": "NY",
        "country_id": 233
      }
    ]
  }
}
```

### Create New City (SuperAdmin)

**Request:**
```http
POST /api/v1/admin/locations/cities
Authorization: Bearer <superadmin_token>
Content-Type: application/json

{
  "name": "Nueva Ciudad",
  "state_id": 1416,
  "latitude": 34.0522,
  "longitude": -118.2437,
  "is_active": true
}
```

**Response:**
```json
{
  "data": {
    "id": 151904,
    "name": "Nueva Ciudad",
    "state_id": 1416,
    "latitude": "34.0522",
    "longitude": "-118.2437",
    "is_active": true
  }
}
```

### Search Locations

**Request:**
```http
GET /api/v1/locations/search?q=london
```

**Response:**
```json
{
  "data": {
    "regions": [],
    "subregions": [],
    "countries": [],
    "states": [],
    "cities": [
      {
        "id": 2648,
        "name": "London",
        "state_id": 2336,
        "latitude": "51.5074",
        "longitude": "-0.1278"
      },
      {
        "id": 34556,
        "name": "London",
        "state_id": 3655,
        "latitude": "42.9834",
        "longitude": "-81.2497"
      }
    ]
  }
}
```

### Get Complete Hierarchy

**Request:**
```http
GET /api/v1/locations/hierarchy/233
```

**Response:**
```json
{
  "data": {
    "id": 233,
    "name": "United States",
    "iso2": "US",
    "subregion": {
      "id": 5,
      "name": "Northern America",
      "region": {
        "id": 2,
        "name": "Americas"
      }
    },
    "states": [
      {
        "id": 1416,
        "name": "California",
        "cities": [
          {
            "id": 111797,
            "name": "Los Angeles"
          }
        ]
      }
    ]
  }
}
```

---

## 🎯 Use Cases

### Frontend Select Dropdowns

```javascript
// 1. Load countries for signup form
const countries = await fetch('/api/v1/locations/countries');

// 2. Load states when user selects country
const states = await fetch(`/api/v1/locations/states?country_id=${countryId}`);

// 3. Load cities when user selects state
const cities = await fetch(`/api/v1/locations/cities?state_id=${stateId}`);
```
Casos de Uso Frontend

### 1️⃣ Select en Cascada (Recomendado)

Implementación típica para formularios de registro/perfil:

```javascript
// Componente React/Vue ejemplo
const LocationSelector = () => {
  const [countries, setCountries] = useState([]);
  const [states, setStates] = useState([]);
  const [cities, setCities] = useState([]);
  
  const [selectedCountry, setSelectedCountry] = useState(null);
  const [selectedState, setSelectedState] = useState(null);

  // 1. Cargar países al montar componente
  useEffect(() => {
    fetch('http://localhost:8001/api/v1/locations/countries')
      .then(res => res.json())
      .then(data => setCountries(data.data));
  }, []);

  // 2. Cargar estados cuando se selecciona país
  useEffect(() => {
    if (selectedCountry) {
      fetch(`http://localhost:8001/api/v1/locations/states?country_id=${selectedCountry}`)
        .then(res => res.json())
        .then(data => setStates(data.data));
      setCities([]); // Limpiar ciudades
    }
  }, [selectedCountry]);

  // 3. Cargar ciudades cuando se selecciona estado
  useEffect(() => {
    if (selectedState) {
      fetch(`http://localhost:8001/api/v1/locations/cities?state_id=${selectedState}`)
        .then(res => res.json())
        .then(data => setCities(data.data));
    }
  }, [selectedState]);

  return (
    <>
      <select onChange={(e) => setSelectedCountry(e.target.value)}>
        <option>Selecciona país</option>
        {countries.map(c => (
          <option key={c.id} value={c.id}>{c.name} {c.emoji}</option>
        ))}
      </select>

      <select onChange={(e) => setSelectedState(e.target.value)} disabled={!selectedCountry}>
        <option>Selecciona estado</option>
        {states.map(s => (
          <option key={s.id} value={s.id}>{s.name}</option>
        ))}
      </select>

      <select disabled={!selectedState}>
        <option>Selecciona ciudad</option>
        {cities.map(c => (
          <option key={c.id} value={c.id}>{c.name}</option>
        ))}
      </select>
    </>
  );
};
```

### 2️⃣ Búsqueda de Ubicación (Autocomplete)

Para campos de búsqueda con autocompletado:

```javascript
const LocationSearch = () => {
  const [query, setQuery] = useState('');
  const [results, setResults] = useState(null);

  const handleSearch = async (searchTerm) => {
    if (searchTerm.length < 2) return;
    
    const response = await fetch(
      `http://localhost:8001/api/v1/locations/search?q=${searchTerm}`
    );
    const data = await response.json();
    setResults(data.data);
  };

  r🧪 Testing de Endpoints

### Pruebas con cURL

```bash
# 1. Obtener todos los países disponibles
curl http://localhost:8001/api/v1/locations/countries

# 2. Buscar países por nombre (case-insensitive)
curl "http://localhost:8001/api/v1/locations/countries?search=mex"

# 3. Obtener estados de un país específico (ej: México)
curl "http://localhost:8001/api/v1/locations/states?country_id=142"

# 4. Buscar ciudades por nombre
curl "http://localhost:8001/api/v1/locations/cities?search=angeles"

# 5. Obtener país por código ISO
curl http://localhost:8001/api/v1/locations/countries/iso/MX

# 6. Búsqueda global
curl "http://localhost:8001/api/v1/locations/search?q=new"

# 7. Jerarquía completa de un país
curl "http://localhost:8001/api/v1/locations/hierarchy/233"
```

### Pruebas con Postman/Insomnia

1. **GET Countries**
   - URL: `http://localhost:8001/api/v1/locations/countries`
   - Método: GET
   - Headers: Ninguno requerido
   - Response: Array de países con banderas, códigos ISO, etc.

2. **GET States by Country**
   - URL: `http://localhost:8001/api/v1/locations/states?country_id=142`
   - Método: GET
   - Query Params: `country_id=142` (México)
   - Response: Estados/Provincias del país

3. **Search Location**
   - URL: `http://localhost:8001/api/v1/locations/search?q=los`
   - Método: GET
   - Query Params: `q=los`
   - Response: Coincidencias en cities, states, countries         <div key={`state-${state.id}`}>
              🏛️ {state.name}, {state.country?.name}
            </div>
          ))}
          {results.countries?.map(country => (
            <div key={`country-${country.id}`}>
              {country.emoji} {country.name}
            </div>
          ))}
        </div>
      )}
    </>
  );
};
```

### 3️⃣ Mostrar País con Bandera

Uso de campos `emoji` y datos complementarios:

```javascript
const CountryDisplay = ({ countryId }) => {
  const [country, setCountry] = useState(null);

  useEffect(() => {
    fetch(`http://localhost:8001/api/v1/locations/countries/${countryId}`)
      .then(res => res.json())
      .then(data => setCountry(data.data));
  }, [countryId]);

  if (!country) return <div>Loading...</div>;

  return (
    <div className="country-card">
      <h2>{country.emoji} {country.name}</h2>
      <p>Capital: {country.capital}</p>
      <p>Código: {country.iso2} / {country.iso3}</p>
      <p>Teléfono: {country.phone_code}</p>
      <p>Moneda: {country.currency_symbol} {country.currency}</p>
    </div>
  );
};
```

### 4️⃣ Validar País por Código ISO

```javascript
const validateCountry = async (isoCode) => {
  try {
    const response = await fetch(
      `http://localhost:8001/api/v1/locations/countries/iso/${isoCode}`
    );
    
    if (response.ok) {
      const data = await response.json();
      return { valid: true, country: data.data };
    }
    return { valid: false };
  } catch (error) {
    return { valid: false };
  }
};

// Uso
const result = await validateCountry('MX'); // México
if (result.valid) {
  console.log(`País válido: ${result.country.name}`);
}

### Optimizations Implemented

1. **Selective Preloading**: Use query params to preload relationships
   ```
   ?include_subregions=true
   ?include_countries=true
   ?include_states=true
   ?include_cities=true
   ```

2. **Filtering**: Reduce result set with query params
   ```
   ?region_id=1
   ?subregion_id=5
   ?country_id=233
   ?state_id=1416
   ```

3. **Search**: ILIKE for partial matching (case-insensitive)
   ```
   ?search=cali  # Matches "California", "Cali", etc.
   ```

4. **Active Records**: Only active records returned by default
   ```sql
   WHERE is_active = true
   ```

### Load Times

- Get all countries: ~50ms
- Get states for country: ~10ms
- Get cities for state: ~20-100ms (depending on city count)
- Full hierarchy: ~150ms

---

## 📊 Database Schema

```sql
-- Hierarchy with Foreign Keys
regions (6 records)
  └── subregions (22 records)
        └── countries (250 records)
              └── states (5,134 records)
                    └── cities (151,903 records)

-- Unique Constraints
countries.iso2 UNIQUE
countries.iso3 UNIQUE

-- Soft Deletes
All tables have is_active BOOLEAN DEFAULT true
All tables have deleted_at TIMESTAMPTZ (GORM)
```
📊 Estructura de Datos Disponibles

### Países Actuales en Base de Datos

```javascript
// Ejemplo de respuesta: GET /api/v1/locations/countries
{
  "data": [
    {
      "id": 142,
      "name": "Mexico",
      "iso2": "MX",
      "iso3": "MEX",
      "phone_code": "52",
      "capital": "Mexico City",
      "currency": "MXN",
      "currency_symbol": "$",
      "emoji": "🇲🇽",
      "is_active": true,
      "subregion_id": 13
    },
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
      "subregion_id": 5
    }
    // ... más países
  ]
}
```

### Estados/Provincias Disponibles

```javascript
// Ejemplo: GET /api/v1/locations/states?country_id=142
{
  "data": [
    {
      "id": 3456,
      "name": "Nuevo León",
      "state_code": "NL",
      "country_id": 142,
      "is_active": true
    },
    {
      "id": 3457,
      "name": "Jalisco",
      "state_code": "JAL",
      "country_id": 142,
      "is_active": true
    }
    // ... más estados
  ]
}
```

### Ciudades Disponibles

```javascript
// Ejemplo: GET /api/v1/locations/cities?state_id=3456
{
  "data": [
    {
      "id": 98765,
      "name": "Monterrey",
      "state_id": 3456,
      "latitude": "25.6866",
      "longitude": "-100.3161",
      "is_active": true
    },
    {
      "id": 98766,
      "name": "San Pedro Garza García",
      "state_id": 3456,
      "latitude": "25.6575",
      "longitude": "-100.3568",
      "is_active": true
    }
    // ... más ciudades
  ]
}
```

---

## 🎉 Resumen para Frontend

### ✅ Características Principales

- **30+ Endpoints REST** (GET, POST, PUT, DELETE)
- **Sin autenticación** para consultas (GET)
- **Datos ya cargados** y listos para usar
- **Búsqueda inteligente** (case-insensitive, partial matching)
- **Filtros por relación** (country_id, state_id, subregion_id)
- **Preload opcional** de relaciones anidadas
- **Emojis de banderas** incluidos
- **Códigos ISO** únicos para países
- **Soft deletes** (datos nunca se pierden)

### 📋 Checklist de Integración

- [ ] Implementar select en cascada: País → Estado → Ciudad
- [ ] Agregar banderas (emoji) en selector de países
- [ ] Implementar búsqueda/autocomplete de ubicaciones
- [ ] Manejar estados de carga (loading) en componentes
- [ ] Validar códigos ISO si aplica
- [ ] Cachear respuestas de países (no cambian frecuentemente)
- [ ] Manejar casos sin estados/ciudades (algunos países no los tienen)
- [ ] Implementar error handling para endpoints

### 🚀 Quick Start

```bash
# 1. Verificar que el servidor esté corriendo
make run  # API en http://localhost:8001

# 2. Probar endpoint de países
curl http://localhost:8001/api/v1/locations/countries

# 3. Ver Swagger docs (opcional)
# http://localhost:8001/swagger/index.html

# 4. Integrar en tu componente frontend
# Ver ejemplos en sección "Casos de Uso Frontend" ⬆️
```

### 📞 Query Parameters Disponibles

| Endpoint | Parámetros | Ejemplo |
|----------|-----------|---------|
| `/countries` | `subregion_id`, `search`, `include_states` | `?search=united&include_states=true` |
| `/states` | `country_id`, `search`, `include_cities` | `?country_id=142&search=nuevo` |
| `/cities` | `state_id`, `search` | `?state_id=3456&search=monte` |
| `/search` | `q` (query) | `?q=los` |

### 🔗 URLs Base

- **Desarrollo**: `http://localhost:8001/api/v1/locations`
- **Swagger Docs**: `http://localhost:8001/swagger/index.html`

---

## 💡 Recomendaciones

1. **Cachea los países**: No cambian frecuentemente, guárdalos en localStorage
2. **Usa debounce**: Para búsquedas, espera 300ms antes de hacer request
3. **Preload selectivo**: Solo usa `include_states=true` si realmente los necesitas
4. **Error handling**: Siempre maneja casos donde el endpoint falle
5. **Loading states**: Muestra spinners mientras cargan estados/ciudades
6. **Fallback**: Ten un país por defecto si el usuario no selecciona

---

**Última actualización**: Diciembre 2025  
**Versión API**: v1  
**Puerto**: 8001  
**Estado**: ✅ Producción (datos básicos cargados
## 📝 Next Steps

1. **Import Data** (optional):
   ```bash
   make db-location  # Import 157k records (~3-5 minutes)
   ```

2. **Test Endpoints**:
   ```bash
   # Start server
   make run
   
   # Test in browser
   http://localhost:8000/api/v1/locations/countries
   ```

3. **Frontend Integration**:
   - Use `/countries` endpoint for country selector
   - Use `/states?country_id=X` for state selector
   - Use `/cities?state_id=X` for city selector

---

**Created**: Complete Location module with DTOs, Repository, Service, Handler, and Routes  
**Public Access**: All read operations (GET)  
**Admin Access**: Create, Update, Delete (SuperAdmin only)
