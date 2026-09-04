package transport

import (
	"dvra-api/internal/modules/locations/service"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes monta el catálogo geográfico bajo /locations.
//
// Todas las rutas son públicas y de solo lectura: son datos de referencia
// (región → subregión → país → estado → ciudad) que el formulario de una
// vacante o una career page necesitan antes de haber iniciado sesión. No hay
// scoping por empresa porque el catálogo no pertenece a ninguna.
func RegisterRoutes(rg *gin.RouterGroup, svc service.LocationService) {
	h := NewLocationHandler(svc)

	locations := rg.Group("/locations")
	{
		// Regions
		locations.GET("/regions", h.GetAllRegions)
		locations.GET("/regions/:id", h.GetRegionByID)

		// Subregions
		locations.GET("/subregions", h.GetAllSubregions)
		locations.GET("/subregions/:id", h.GetSubregionByID)

		// Countries
		locations.GET("/countries", h.GetAllCountries)
		locations.GET("/countries/:id", h.GetCountryByID)
		locations.GET("/countries/iso/:iso", h.GetCountryByISO)

		// States
		locations.GET("/states", h.GetAllStates)
		locations.GET("/states/:id", h.GetStateByID)

		// Cities
		locations.GET("/cities", h.GetAllCities)
		locations.GET("/cities/:id", h.GetCityByID)

		// Helpers
		locations.GET("/hierarchy/:id", h.GetLocationHierarchy)
		locations.GET("/search", h.SearchLocations)
	}
}
