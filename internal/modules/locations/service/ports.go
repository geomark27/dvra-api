package service

import "dvra-api/internal/shared/models"

// LocationRepository es el puerto de salida hacia la persistencia del catálogo
// geográfico.
//
// Se declara AQUÍ, en el consumidor, y no junto a su implementación
// (ESTANDARES_DESARROLLO.md §3.3 regla 3): el service define lo que necesita y
// el composition root le inyecta el adaptador. En Go la satisfacción es
// estructural, así que repository/ no necesita importar este paquete.
type LocationRepository interface {
	// Regions
	GetAllRegions(includeSubregions bool) ([]models.Region, error)
	GetRegionByID(id uint, includeSubregions bool) (*models.Region, error)
	CreateRegion(region *models.Region) error
	UpdateRegion(region *models.Region) error
	DeleteRegion(id uint) error

	// Subregions
	GetAllSubregions(regionID *uint) ([]models.Subregion, error)
	GetSubregionByID(id uint, includeCountries bool) (*models.Subregion, error)
	CreateSubregion(subregion *models.Subregion) error
	UpdateSubregion(subregion *models.Subregion) error
	DeleteSubregion(id uint) error

	// Countries
	GetAllCountries(subregionID *uint, search string) ([]models.Country, error)
	GetCountryByID(id uint, includeStates bool) (*models.Country, error)
	GetCountryByISO(iso string) (*models.Country, error)
	CreateCountry(country *models.Country) error
	UpdateCountry(country *models.Country) error
	DeleteCountry(id uint) error

	// States
	GetAllStates(countryID *uint, search string) ([]models.State, error)
	GetStateByID(id uint, includeCities bool) (*models.State, error)
	CreateState(state *models.State) error
	UpdateState(state *models.State) error
	DeleteState(id uint) error

	// Cities
	GetAllCities(stateID *uint, search string) ([]models.City, error)
	GetCityByID(id uint) (*models.City, error)
	CreateCity(city *models.City) error
	UpdateCity(city *models.City) error
	DeleteCity(id uint) error

	// Helpers
	GetLocationHierarchy(countryID uint) (*models.Country, error)
	SearchLocations(search string) (map[string]interface{}, error)
}
