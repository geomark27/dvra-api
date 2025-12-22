# 🌍 Sistema de Ubicaciones - Guía de Uso

## 📊 Estructura Simplificada

```
Region (Continente)
  └── Subregion (Área geográfica)
        └── Country (País)
              └── State (Estado/Provincia)
                    └── City (Ciudad)
```

## 🔗 Relaciones Configuradas

- **Region** → tiene muchas **Subregions**
- **Subregion** → pertenece a **Region**, tiene muchos **Countries**
- **Country** → pertenece a **Subregion**
- **State** → pertenece a **Country** (si existe)
- **City** → pertenece a **State** (si existe)

## ✅ Setup Rápido

```bash
# 1. Crear tablas y seeders básicos (7 segundos)
make fresh

# 2. (Opcional) Si necesitas datos de ubicación masivos
make db-location  # Tarda 3-5 minutos, carga 150k+ registros
```

## 💡 Ejemplos de Uso en Go

### 1. Crear Regiones y Subregiones

```go
// Crear una región
region := models.Region{
    Name:     "Americas",
    IsActive: true,
}
db.Create(&region)

// Crear subregión asociada
subregion := models.Subregion{
    Name:     "South America",
    RegionID: region.ID,
    IsActive: true,
}
db.Create(&subregion)
```

### 2. Crear País con Relación

```go
// Crear país
country := models.Country{
    Name:        "Mexico",
    Iso2:        "MX",
    Iso3:        "MEX",
    NumericCode: "484",
    PhoneCode:   "+52",
    SubregionID: &subregionID,  // Asociar a subregión
    IsActive:    true,
}
db.Create(&country)
```

### 3. Consultar con Relaciones

```go
// Obtener país con su subregión y región
var country models.Country
db.Preload("Subregion.Region").
   Where("iso2 = ?", "MX").
   First(&country)

fmt.Printf("%s - %s - %s", 
    country.Name,                    // "Mexico"
    country.Subregion.Name,          // "Central America"
    country.Subregion.Region.Name)   // "Americas"
```

### 4. Obtener Todos los Países de una Región

```go
var region models.Region
db.Preload("Subregions.Countries").
   Where("name = ?", "Americas").
   First(&region)

for _, subregion := range region.Subregions {
    for _, country := range subregion.Countries {
        fmt.Printf("%s (%s)\n", country.Name, country.Iso2)
    }
}
```

### 5. Búsqueda de Países

```go
// Por código ISO
var country models.Country
db.Where("iso2 = ?", "US").First(&country)

// Por nombre (búsqueda parcial)
var countries []models.Country
db.Where("name ILIKE ?", "%united%").Find(&countries)

// Países activos de una subregión
var countries []models.Country
db.Where("subregion_id = ? AND is_active = ?", subregionID, true).
   Find(&countries)
```

## 📋 Campos Disponibles

### Region
- `ID` - ID auto-incremental
- `Name` - Nombre de la región (ej: "Europe", "Asia")
- `IsActive` - Estado activo/inactivo
- `CreatedAt`, `UpdatedAt`, `DeletedAt` - Timestamps

### Subregion
- `ID` - ID auto-incremental
- `Name` - Nombre de la subregión (ej: "Western Europe")
- `RegionID` - ID de la región padre
- `IsActive` - Estado activo/inactivo
- `CreatedAt`, `UpdatedAt`, `DeletedAt` - Timestamps

### Country
- `ID` - ID auto-incremental
- `Name` - Nombre del país
- `Iso2` - Código ISO de 2 letras (único) ej: "MX"
- `Iso3` - Código ISO de 3 letras (único) ej: "MEX"
- `NumericCode` - Código numérico
- `PhoneCode` - Código telefónico (ej: "+52")
- `Timezones` - Zonas horarias (texto)
- `SubregionID` - ID de la subregión
- `IsActive` - Estado activo/inactivo
- `CreatedAt`, `UpdatedAt`, `DeletedAt` - Timestamps

## 🎯 Casos de Uso

### 1. Dropdown en Cascada

```go
// API: GET /api/regions
func GetRegions(c *gin.Context) {
    var regions []models.Region
    db.Where("is_active = ?", true).Find(&regions)
    c.JSON(200, regions)
}

// API: GET /api/regions/:id/subregions
func GetSubregions(c *gin.Context) {
    regionID := c.Param("id")
    var subregions []models.Subregion
    db.Where("region_id = ? AND is_active = ?", regionID, true).
       Find(&subregions)
    c.JSON(200, subregions)
}

// API: GET /api/subregions/:id/countries
func GetCountries(c *gin.Context) {
    subregionID := c.Param("id")
    var countries []models.Country
    db.Where("subregion_id = ? AND is_active = ?", subregionID, true).
       Find(&countries)
    c.JSON(200, countries)
}
```

### 2. Validación de País

```go
func ValidateCountryISO(iso2 string) (bool, error) {
    var count int64
    err := db.Model(&models.Country{}).
        Where("iso2 = ? AND is_active = ?", iso2, true).
        Count(&count).Error
    
    return count > 0, err
}
```

### 3. Formulario de Registro

```go
type CompanyRegistration struct {
    Name      string `json:"name"`
    CountryID uint   `json:"country_id"`
}

func CreateCompany(registration CompanyRegistration) error {
    // Verificar que el país existe
    var country models.Country
    if err := db.First(&country, registration.CountryID).Error; err != nil {
        return errors.New("país no válido")
    }
    
    // Crear empresa con país asociado
    company := models.Company{
        Name:      registration.Name,
        CountryID: &registration.CountryID,
    }
    
    return db.Create(&company).Error
}
```

## 🗄️ Seeders Personalizados

Si quieres cargar datos manualmente:

```go
// internal/database/seeders/location_basic_seeder.go

package seeders

import (
    "dvra-api/internal/app/models"
    "gorm.io/gorm"
    "log"
)

type LocationBasicSeeder struct{}

func (s *LocationBasicSeeder) Run(db *gorm.DB) error {
    log.Println("🌍 Seeding basic location data...")
    
    // Crear Américas
    americas := models.Region{Name: "Americas", IsActive: true}
    db.Create(&americas)
    
    // Subregiones
    northAmerica := models.Subregion{
        Name: "North America", 
        RegionID: americas.ID, 
        IsActive: true,
    }
    db.Create(&northAmerica)
    
    southAmerica := models.Subregion{
        Name: "South America", 
        RegionID: americas.ID, 
        IsActive: true,
    }
    db.Create(&southAmerica)
    
    // Países
    countries := []models.Country{
        {
            Name: "Mexico", 
            Iso2: "MX", 
            Iso3: "MEX", 
            NumericCode: "484",
            PhoneCode: "+52",
            SubregionID: &northAmerica.ID,
            IsActive: true,
        },
        {
            Name: "United States", 
            Iso2: "US", 
            Iso3: "USA", 
            NumericCode: "840",
            PhoneCode: "+1",
            SubregionID: &northAmerica.ID,
            IsActive: true,
        },
        {
            Name: "Colombia", 
            Iso2: "CO", 
            Iso3: "COL", 
            NumericCode: "170",
            PhoneCode: "+57",
            SubregionID: &southAmerica.ID,
            IsActive: true,
        },
    }
    
    for _, country := range countries {
        db.Create(&country)
    }
    
    log.Println("✅ Basic location data seeded")
    return nil
}
```

## 🚀 Ventajas del Enfoque Simplificado

✅ **Rápido**: `make fresh` toma solo 7 segundos  
✅ **Limpio**: Solo las tablas y campos esenciales  
✅ **Flexible**: Agrega datos según necesites  
✅ **Escalable**: Puedes cargar world.sql cuando lo requieras  
✅ **Mantenible**: Relaciones claras y documentadas  

## 📝 Notas

- Los códigos ISO2 e ISO3 tienen índice único
- Las relaciones usan `foreignKey` para integridad referencial
- `SubregionID` es nullable en Country (por si no aplica)
- Todos los modelos usan `gorm.Model` (soft deletes incluido)
- El campo `IsActive` permite deshabilitar sin borrar

---

¡Listo para usar! 🎉
