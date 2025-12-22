# 📖 Swagger/OpenAPI Documentation - Guía de Uso

## ✅ Configuración Completada

Tu proyecto ya tiene Swagger integrado y listo para usar.

### Archivos Configurados:

1. **[cmd/dvra-api/main.go](../cmd/dvra-api/main.go)**: Anotaciones principales de la API
2. **[internal/platform/server/server.go](../internal/platform/server/server.go)**: Importación de docs generados
3. **[internal/platform/server/routes.go](../internal/platform/server/routes.go)**: Ruta `/swagger/*`
4. **[docs/](../docs/)**: Carpeta con documentación generada (swagger.json, swagger.yaml, docs.go)

---

## 🚀 Cómo Usar

### 1. Generar Documentación

Cada vez que agregues o modifiques endpoints, regenera la documentación:

```bash
make swagger
```

O manualmente:

```bash
~/go/bin/swag init -g cmd/dvra-api/main.go -o docs
```

### 2. Iniciar el Servidor

```bash
make run
# o
go run cmd/dvra-api/main.go
```

### 3. Acceder a Swagger UI

Abre tu navegador en:

```
http://localhost:8000/swagger/index.html
```

Verás la interfaz interactiva de Swagger con todos tus endpoints documentados.

---

## 📝 Cómo Documentar Endpoints

### Anotaciones en Handlers

Agrega comentarios antes de cada función handler siguiendo el formato Swaggo:

#### Ejemplo Básico (GET sin auth):

```go
// GetAllRegions godoc
// @Summary      Obtener todas las regiones
// @Description  Retorna la lista completa de regiones activas
// @Tags         Locations
// @Accept       json
// @Produce      json
// @Param        include_subregions  query  bool  false  "Incluir subregiones"
// @Success      200  {array}   dtos.RegionDTO
// @Failure      500  {object}  map[string]interface{}
// @Router       /locations/regions [get]
func (h *LocationHandler) GetAllRegions(c *gin.Context) {
    // ...
}
```

#### Ejemplo con Autenticación (POST):

```go
// CreateRegion godoc
// @Summary      Crear nueva región
// @Description  Crea una región (solo SuperAdmin)
// @Tags         Locations
// @Accept       json
// @Produce      json
// @Param        region  body      dtos.CreateRegionDTO  true  "Datos de la región"
// @Success      201     {object}  dtos.RegionDTO
// @Failure      400     {object}  map[string]interface{}
// @Failure      401     {object}  map[string]interface{}
// @Failure      500     {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /admin/locations/regions [post]
func (h *LocationHandler) CreateRegion(c *gin.Context) {
    // ...
}
```

#### Ejemplo con Path Parameter:

```go
// GetCountryByID godoc
// @Summary      Obtener país por ID
// @Description  Retorna información detallada de un país
// @Tags         Locations
// @Accept       json
// @Produce      json
// @Param        id                path   int   true   "ID del país"
// @Param        include_states    query  bool  false  "Incluir estados"
// @Success      200  {object}  dtos.CountryDTO
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /locations/countries/{id} [get]
func (h *LocationHandler) GetCountryByID(c *gin.Context) {
    // ...
}
```

---

## 🏷️ Anotaciones Disponibles

### Metadata del Endpoint:

| Anotación | Descripción | Ejemplo |
|-----------|-------------|---------|
| `@Summary` | Título corto del endpoint | `@Summary Get user by ID` |
| `@Description` | Descripción detallada | `@Description Returns detailed user info` |
| `@Tags` | Agrupa endpoints | `@Tags Users` |
| `@Accept` | Content-Type aceptado | `@Accept json` |
| `@Produce` | Content-Type de respuesta | `@Produce json` |

### Parámetros:

| Tipo | Formato | Ejemplo |
|------|---------|---------|
| **Path** | `@Param name path type required "description"` | `@Param id path int true "User ID"` |
| **Query** | `@Param name query type required "description"` | `@Param search query string false "Search term"` |
| **Body** | `@Param name body type required "description"` | `@Param user body dtos.CreateUserDTO true "User data"` |
| **Header** | `@Param name header type required "description"` | `@Param Authorization header string true "Bearer token"` |

### Respuestas:

```go
// @Success 200 {object}  dtos.UserDTO             "Success"
// @Success 200 {array}   dtos.UserDTO             "List of users"
// @Failure 400 {object}  map[string]interface{}   "Bad Request"
// @Failure 401 {object}  map[string]interface{}   "Unauthorized"
// @Failure 404 {object}  map[string]interface{}   "Not Found"
// @Failure 500 {object}  map[string]interface{}   "Internal Server Error"
```

### Seguridad:

```go
// @Security BearerAuth  // Requiere JWT token
```

### Router:

```go
// @Router /api/v1/users [get]
// @Router /api/v1/users/{id} [get]
// @Router /api/v1/users [post]
// @Router /api/v1/users/{id} [put]
// @Router /api/v1/users/{id} [delete]
```

---

## 🎯 Ejemplos Prácticos

### GET con Filtros

```go
// GetAllCountries godoc
// @Summary      Listar países
// @Description  Obtiene países con filtros opcionales
// @Tags         Locations
// @Accept       json
// @Produce      json
// @Param        subregion_id  query  int     false  "Filtrar por subregión"
// @Param        search        query  string  false  "Buscar por nombre o código ISO"
// @Success      200  {array}   dtos.CountryDTO
// @Failure      500  {object}  map[string]interface{}
// @Router       /locations/countries [get]
```

### POST con Auth

```go
// CreateUser godoc
// @Summary      Crear usuario
// @Description  Crea un nuevo usuario en la empresa actual
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        user  body  dtos.CreateUserDTO  true  "Datos del usuario"
// @Success      201   {object}  dtos.UserDTO
// @Failure      400   {object}  map[string]interface{}  "Validación fallida"
// @Failure      401   {object}  map[string]interface{}  "No autenticado"
// @Failure      500   {object}  map[string]interface{}  "Error del servidor"
// @Security     BearerAuth
// @Router       /users [post]
```

### PUT con Path Param

```go
// UpdateCity godoc
// @Summary      Actualizar ciudad
// @Description  Actualiza los datos de una ciudad (SuperAdmin)
// @Tags         Locations
// @Accept       json
// @Produce      json
// @Param        id    path  int                  true  "ID de la ciudad"
// @Param        city  body  dtos.UpdateCityDTO   true  "Datos a actualizar"
// @Success      200   {object}  dtos.CityDTO
// @Failure      400   {object}  map[string]interface{}
// @Failure      404   {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /admin/locations/cities/{id} [put]
```

### DELETE

```go
// DeleteCountry godoc
// @Summary      Eliminar país
// @Description  Elimina un país (soft delete)
// @Tags         Locations
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "ID del país"
// @Success      204
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Security     BearerAuth
// @Router       /admin/locations/countries/{id} [delete]
```

---

## 🔐 Probar Endpoints con Autenticación

### 1. Login para obtener token:

En Swagger UI:
1. Busca el endpoint `POST /api/v1/auth/login`
2. Click en "Try it out"
3. Ingresa credenciales:
   ```json
   {
     "email": "admin@example.com",
     "password": "password123"
   }
   ```
4. Copia el `access_token` de la respuesta

### 2. Autorizar requests:

1. Click en el botón **"Authorize"** (🔒) en la parte superior derecha
2. Ingresa: `Bearer <tu_token_aqui>`
3. Click "Authorize"
4. Ahora puedes ejecutar endpoints protegidos

---

## 📊 DTOs en Swagger

Tus DTOs se documentan automáticamente si usas tags JSON y binding:

```go
type CreateCountryDTO struct {
    Name         string  `json:"name" binding:"required,min=2,max=100"`
    ISO2         string  `json:"iso2" binding:"required,len=2"`
    ISO3         string  `json:"iso3" binding:"required,len=3"`
    PhoneCode    string  `json:"phone_code" binding:"required"`
    Capital      string  `json:"capital"`
    Currency     string  `json:"currency"`
    SubregionID  *uint   `json:"subregion_id"`
    IsActive     bool    `json:"is_active"`
}
```

Swagger detecta:
- Nombres de campos (`json` tag)
- Si son requeridos (`binding:"required"`)
- Validaciones (`min`, `max`, `len`)
- Tipos de datos (string, int, bool, etc.)

---

## 🎨 Personalizar Documentación General

Edita [cmd/dvra-api/main.go](../cmd/dvra-api/main.go):

```go
// @title           DVRA API
// @version         1.2.0
// @description     API para sistema de reclutamiento y gestión de candidatos
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  support@dvra.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8000
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
```

Luego regenera: `make swagger`

---

## 🔄 Workflow Recomendado

1. **Crea/Modifica un Handler**
2. **Agrega anotaciones Swagger** (comentarios antes de la función)
3. **Regenera documentación**: `make swagger`
4. **Verifica en Swagger UI**: http://localhost:8000/swagger/index.html
5. **Prueba el endpoint** directamente desde Swagger

---

## 📁 Archivos Generados

Después de `make swagger`:

```
docs/
├── docs.go          # Documentación en formato Go (importado por server.go)
├── swagger.json     # Spec OpenAPI 3.0 en JSON
└── swagger.yaml     # Spec OpenAPI 3.0 en YAML
```

**Nota**: Estos archivos se regeneran automáticamente. No los edites manualmente.

---

## 🎉 Resumen

✅ **Swagger instalado y configurado**  
✅ **Ruta**: http://localhost:8000/swagger/index.html  
✅ **Comando**: `make swagger` para regenerar  
✅ **Autenticación**: Bearer token configurado  
✅ **DTOs**: Auto-documentados desde structs  

---

## 📚 Recursos

- [Swaggo Documentation](https://github.com/swaggo/swag)
- [Swagger Annotations](https://github.com/swaggo/swag#declarative-comments-format)
- [OpenAPI Specification](https://swagger.io/specification/)

---

**¡Tu API ahora tiene documentación interactiva completa!** 🎊
