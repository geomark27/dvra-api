package handlers

import (
	"dvra-api/internal/app/dtos"
	"dvra-api/internal/app/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type LocationHandler struct {
	service services.LocationService
}

func NewLocationHandler(service services.LocationService) *LocationHandler {
	return &LocationHandler{service: service}
}

// ==================== REGIONS ====================

// GetAllRegions godoc
// @Summary Get all regions
// @Tags Locations
// @Param include_subregions query bool false "Include subregions"
// @Success 200 {array} dtos.RegionDTO
// @Router /locations/regions [get]
func (h *LocationHandler) GetAllRegions(c *gin.Context) {
	includeSubregions := c.Query("include_subregions") == "true"

	regions, err := h.service.GetAllRegions(includeSubregions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": regions})
}

// GetRegionByID godoc
// @Summary Get region by ID
// @Tags Locations
// @Param id path int true "Region ID"
// @Param include_subregions query bool false "Include subregions"
// @Success 200 {object} dtos.RegionDTO
// @Router /locations/regions/{id} [get]
func (h *LocationHandler) GetRegionByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	includeSubregions := c.Query("include_subregions") == "true"

	region, err := h.service.GetRegionByID(uint(id), includeSubregions)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": region})
}

// CreateRegion godoc
// @Summary Create region
// @Tags Locations
// @Param region body dtos.CreateRegionDTO true "Region data"
// @Success 201 {object} dtos.RegionDTO
// @Router /locations/regions [post]
func (h *LocationHandler) CreateRegion(c *gin.Context) {
	var dto dtos.CreateRegionDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	region, err := h.service.CreateRegion(dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": region})
}

// UpdateRegion godoc
// @Summary Update region
// @Tags Locations
// @Param id path int true "Region ID"
// @Param region body dtos.UpdateRegionDTO true "Region data"
// @Success 200 {object} dtos.RegionDTO
// @Router /locations/regions/{id} [put]
func (h *LocationHandler) UpdateRegion(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var dto dtos.UpdateRegionDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	region, err := h.service.UpdateRegion(uint(id), dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": region})
}

// DeleteRegion godoc
// @Summary Delete region
// @Tags Locations
// @Param id path int true "Region ID"
// @Success 204
// @Router /locations/regions/{id} [delete]
func (h *LocationHandler) DeleteRegion(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.service.DeleteRegion(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// ==================== SUBREGIONS ====================

// GetAllSubregions godoc
// @Summary Get all subregions
// @Tags Locations
// @Param region_id query int false "Filter by region ID"
// @Success 200 {array} dtos.SubregionDTO
// @Router /locations/subregions [get]
func (h *LocationHandler) GetAllSubregions(c *gin.Context) {
	var regionID *uint
	if regionIDStr := c.Query("region_id"); regionIDStr != "" {
		id, err := strconv.ParseUint(regionIDStr, 10, 32)
		if err == nil {
			idUint := uint(id)
			regionID = &idUint
		}
	}

	subregions, err := h.service.GetAllSubregions(regionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": subregions})
}

// GetSubregionByID godoc
// @Summary Get subregion by ID
// @Tags Locations
// @Param id path int true "Subregion ID"
// @Param include_countries query bool false "Include countries"
// @Success 200 {object} dtos.SubregionDTO
// @Router /locations/subregions/{id} [get]
func (h *LocationHandler) GetSubregionByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	includeCountries := c.Query("include_countries") == "true"

	subregion, err := h.service.GetSubregionByID(uint(id), includeCountries)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": subregion})
}

// CreateSubregion godoc
// @Summary Create subregion
// @Tags Locations
// @Param subregion body dtos.CreateSubregionDTO true "Subregion data"
// @Success 201 {object} dtos.SubregionDTO
// @Router /locations/subregions [post]
func (h *LocationHandler) CreateSubregion(c *gin.Context) {
	var dto dtos.CreateSubregionDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	subregion, err := h.service.CreateSubregion(dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": subregion})
}

// UpdateSubregion godoc
// @Summary Update subregion
// @Tags Locations
// @Param id path int true "Subregion ID"
// @Param subregion body dtos.UpdateSubregionDTO true "Subregion data"
// @Success 200 {object} dtos.SubregionDTO
// @Router /locations/subregions/{id} [put]
func (h *LocationHandler) UpdateSubregion(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var dto dtos.UpdateSubregionDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	subregion, err := h.service.UpdateSubregion(uint(id), dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": subregion})
}

// DeleteSubregion godoc
// @Summary Delete subregion
// @Tags Locations
// @Param id path int true "Subregion ID"
// @Success 204
// @Router /locations/subregions/{id} [delete]
func (h *LocationHandler) DeleteSubregion(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.service.DeleteSubregion(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// ==================== COUNTRIES ====================

// GetAllCountries godoc
// @Summary Get all countries
// @Tags Locations
// @Param subregion_id query int false "Filter by subregion ID"
// @Param search query string false "Search by name or ISO code"
// @Success 200 {array} dtos.CountryDTO
// @Router /locations/countries [get]
func (h *LocationHandler) GetAllCountries(c *gin.Context) {
	var subregionID *uint
	if subregionIDStr := c.Query("subregion_id"); subregionIDStr != "" {
		id, err := strconv.ParseUint(subregionIDStr, 10, 32)
		if err == nil {
			idUint := uint(id)
			subregionID = &idUint
		}
	}

	search := c.Query("search")

	countries, err := h.service.GetAllCountries(subregionID, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": countries})
}

// GetCountryByID godoc
// @Summary Get country by ID
// @Tags Locations
// @Param id path int true "Country ID"
// @Param include_states query bool false "Include states"
// @Success 200 {object} dtos.CountryDTO
// @Router /locations/countries/{id} [get]
func (h *LocationHandler) GetCountryByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	includeStates := c.Query("include_states") == "true"

	country, err := h.service.GetCountryByID(uint(id), includeStates)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": country})
}

// GetCountryByISO godoc
// @Summary Get country by ISO code
// @Tags Locations
// @Param iso path string true "ISO2 or ISO3 code"
// @Success 200 {object} dtos.CountryDTO
// @Router /locations/countries/iso/{iso} [get]
func (h *LocationHandler) GetCountryByISO(c *gin.Context) {
	iso := c.Param("iso")

	country, err := h.service.GetCountryByISO(iso)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": country})
}

// CreateCountry godoc
// @Summary Create country
// @Tags Locations
// @Param country body dtos.CreateCountryDTO true "Country data"
// @Success 201 {object} dtos.CountryDTO
// @Router /locations/countries [post]
func (h *LocationHandler) CreateCountry(c *gin.Context) {
	var dto dtos.CreateCountryDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	country, err := h.service.CreateCountry(dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": country})
}

// UpdateCountry godoc
// @Summary Update country
// @Tags Locations
// @Param id path int true "Country ID"
// @Param country body dtos.UpdateCountryDTO true "Country data"
// @Success 200 {object} dtos.CountryDTO
// @Router /locations/countries/{id} [put]
func (h *LocationHandler) UpdateCountry(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var dto dtos.UpdateCountryDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	country, err := h.service.UpdateCountry(uint(id), dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": country})
}

// DeleteCountry godoc
// @Summary Delete country
// @Tags Locations
// @Param id path int true "Country ID"
// @Success 204
// @Router /locations/countries/{id} [delete]
func (h *LocationHandler) DeleteCountry(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.service.DeleteCountry(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// ==================== STATES ====================

// GetAllStates godoc
// @Summary Get all states
// @Tags Locations
// @Param country_id query int false "Filter by country ID"
// @Param search query string false "Search by name"
// @Success 200 {array} dtos.StateDTO
// @Router /locations/states [get]
func (h *LocationHandler) GetAllStates(c *gin.Context) {
	var countryID *uint
	if countryIDStr := c.Query("country_id"); countryIDStr != "" {
		id, err := strconv.ParseUint(countryIDStr, 10, 32)
		if err == nil {
			idUint := uint(id)
			countryID = &idUint
		}
	}

	search := c.Query("search")

	states, err := h.service.GetAllStates(countryID, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": states})
}

// GetStateByID godoc
// @Summary Get state by ID
// @Tags Locations
// @Param id path int true "State ID"
// @Param include_cities query bool false "Include cities"
// @Success 200 {object} dtos.StateDTO
// @Router /locations/states/{id} [get]
func (h *LocationHandler) GetStateByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	includeCities := c.Query("include_cities") == "true"

	state, err := h.service.GetStateByID(uint(id), includeCities)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": state})
}

// CreateState godoc
// @Summary Create state
// @Tags Locations
// @Param state body dtos.CreateStateDTO true "State data"
// @Success 201 {object} dtos.StateDTO
// @Router /locations/states [post]
func (h *LocationHandler) CreateState(c *gin.Context) {
	var dto dtos.CreateStateDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	state, err := h.service.CreateState(dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": state})
}

// UpdateState godoc
// @Summary Update state
// @Tags Locations
// @Param id path int true "State ID"
// @Param state body dtos.UpdateStateDTO true "State data"
// @Success 200 {object} dtos.StateDTO
// @Router /locations/states/{id} [put]
func (h *LocationHandler) UpdateState(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var dto dtos.UpdateStateDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	state, err := h.service.UpdateState(uint(id), dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": state})
}

// DeleteState godoc
// @Summary Delete state
// @Tags Locations
// @Param id path int true "State ID"
// @Success 204
// @Router /locations/states/{id} [delete]
func (h *LocationHandler) DeleteState(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.service.DeleteState(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// ==================== CITIES ====================

// GetAllCities godoc
// @Summary Get all cities
// @Tags Locations
// @Param state_id query int false "Filter by state ID"
// @Param search query string false "Search by name"
// @Success 200 {array} dtos.CityDTO
// @Router /locations/cities [get]
func (h *LocationHandler) GetAllCities(c *gin.Context) {
	var stateID *uint
	if stateIDStr := c.Query("state_id"); stateIDStr != "" {
		id, err := strconv.ParseUint(stateIDStr, 10, 32)
		if err == nil {
			idUint := uint(id)
			stateID = &idUint
		}
	}

	search := c.Query("search")

	cities, err := h.service.GetAllCities(stateID, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": cities})
}

// GetCityByID godoc
// @Summary Get city by ID
// @Tags Locations
// @Param id path int true "City ID"
// @Success 200 {object} dtos.CityDTO
// @Router /locations/cities/{id} [get]
func (h *LocationHandler) GetCityByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	city, err := h.service.GetCityByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": city})
}

// CreateCity godoc
// @Summary Create city
// @Tags Locations
// @Param city body dtos.CreateCityDTO true "City data"
// @Success 201 {object} dtos.CityDTO
// @Router /locations/cities [post]
func (h *LocationHandler) CreateCity(c *gin.Context) {
	var dto dtos.CreateCityDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	city, err := h.service.CreateCity(dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": city})
}

// UpdateCity godoc
// @Summary Update city
// @Tags Locations
// @Param id path int true "City ID"
// @Param city body dtos.UpdateCityDTO true "City data"
// @Success 200 {object} dtos.CityDTO
// @Router /locations/cities/{id} [put]
func (h *LocationHandler) UpdateCity(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var dto dtos.UpdateCityDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	city, err := h.service.UpdateCity(uint(id), dto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": city})
}

// DeleteCity godoc
// @Summary Delete city
// @Tags Locations
// @Param id path int true "City ID"
// @Success 204
// @Router /locations/cities/{id} [delete]
func (h *LocationHandler) DeleteCity(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if err := h.service.DeleteCity(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// ==================== HELPERS ====================

// GetLocationHierarchy godoc
// @Summary Get complete location hierarchy for a country
// @Tags Locations
// @Param id path int true "Country ID"
// @Success 200 {object} dtos.CountryDTO
// @Router /locations/hierarchy/{id} [get]
func (h *LocationHandler) GetLocationHierarchy(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	country, err := h.service.GetLocationHierarchy(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": country})
}

// SearchLocations godoc
// @Summary Search across all location types
// @Tags Locations
// @Param q query string true "Search query"
// @Success 200 {object} dtos.LocationSearchDTO
// @Router /locations/search [get]
func (h *LocationHandler) SearchLocations(c *gin.Context) {
	search := c.Query("q")
	if search == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search query required"})
		return
	}

	results, err := h.service.SearchLocations(search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": results})
}
