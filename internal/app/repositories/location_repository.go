package repositories

import (
	"dvra-api/internal/app/models"
	"dvra-api/internal/database"

	"gorm.io/gorm"
)

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

type locationRepository struct{}

func NewLocationRepository() LocationRepository {
	return &locationRepository{}
}

// ==================== REGIONS ====================

func (r *locationRepository) GetAllRegions(includeSubregions bool) ([]models.Region, error) {
	var regions []models.Region
	query := database.DB.Where("is_active = ?", true)

	if includeSubregions {
		query = query.Preload("Subregions", "is_active = ?", true)
	}

	if err := query.Order("name ASC").Find(&regions).Error; err != nil {
		return nil, err
	}
	return regions, nil
}

func (r *locationRepository) GetRegionByID(id uint, includeSubregions bool) (*models.Region, error) {
	var region models.Region
	query := database.DB

	if includeSubregions {
		query = query.Preload("Subregions", "is_active = ?", true)
	}

	if err := query.First(&region, id).Error; err != nil {
		return nil, err
	}
	return &region, nil
}

func (r *locationRepository) CreateRegion(region *models.Region) error {
	return database.DB.Create(region).Error
}

func (r *locationRepository) UpdateRegion(region *models.Region) error {
	return database.DB.Save(region).Error
}

func (r *locationRepository) DeleteRegion(id uint) error {
	return database.DB.Delete(&models.Region{}, id).Error
}

// ==================== SUBREGIONS ====================

func (r *locationRepository) GetAllSubregions(regionID *uint) ([]models.Subregion, error) {
	var subregions []models.Subregion
	query := database.DB.Where("is_active = ?", true).Preload("Region")

	if regionID != nil {
		query = query.Where("region_id = ?", *regionID)
	}

	if err := query.Order("name ASC").Find(&subregions).Error; err != nil {
		return nil, err
	}
	return subregions, nil
}

func (r *locationRepository) GetSubregionByID(id uint, includeCountries bool) (*models.Subregion, error) {
	var subregion models.Subregion
	query := database.DB.Preload("Region")

	if includeCountries {
		query = query.Preload("Countries", "is_active = ?", true)
	}

	if err := query.First(&subregion, id).Error; err != nil {
		return nil, err
	}
	return &subregion, nil
}

func (r *locationRepository) CreateSubregion(subregion *models.Subregion) error {
	return database.DB.Create(subregion).Error
}

func (r *locationRepository) UpdateSubregion(subregion *models.Subregion) error {
	return database.DB.Save(subregion).Error
}

func (r *locationRepository) DeleteSubregion(id uint) error {
	return database.DB.Delete(&models.Subregion{}, id).Error
}

// ==================== COUNTRIES ====================

func (r *locationRepository) GetAllCountries(subregionID *uint, search string) ([]models.Country, error) {
	var countries []models.Country
	query := database.DB.Where("is_active = ?", true).Preload("Subregion.Region")

	if subregionID != nil {
		query = query.Where("subregion_id = ?", *subregionID)
	}

	if search != "" {
		query = query.Where("name ILIKE ? OR iso2 ILIKE ? OR iso3 ILIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Order("name ASC").Find(&countries).Error; err != nil {
		return nil, err
	}
	return countries, nil
}

func (r *locationRepository) GetCountryByID(id uint, includeStates bool) (*models.Country, error) {
	var country models.Country
	query := database.DB.Preload("Subregion.Region")

	if includeStates {
		query = query.Preload("States", "is_active = ?", true)
	}

	if err := query.First(&country, id).Error; err != nil {
		return nil, err
	}
	return &country, nil
}

func (r *locationRepository) GetCountryByISO(iso string) (*models.Country, error) {
	var country models.Country
	query := database.DB.Preload("Subregion.Region")

	if len(iso) == 2 {
		query = query.Where("UPPER(iso2) = ?", iso)
	} else if len(iso) == 3 {
		query = query.Where("UPPER(iso3) = ?", iso)
	} else {
		return nil, gorm.ErrRecordNotFound
	}

	if err := query.First(&country).Error; err != nil {
		return nil, err
	}
	return &country, nil
}

func (r *locationRepository) CreateCountry(country *models.Country) error {
	return database.DB.Create(country).Error
}

func (r *locationRepository) UpdateCountry(country *models.Country) error {
	return database.DB.Save(country).Error
}

func (r *locationRepository) DeleteCountry(id uint) error {
	return database.DB.Delete(&models.Country{}, id).Error
}

// ==================== STATES ====================

func (r *locationRepository) GetAllStates(countryID *uint, search string) ([]models.State, error) {
	var states []models.State
	query := database.DB.Where("is_active = ?", true).Preload("Country")

	if countryID != nil {
		query = query.Where("country_id = ?", *countryID)
	}

	if search != "" {
		query = query.Where("name ILIKE ?", "%"+search+"%")
	}

	if err := query.Order("name ASC").Find(&states).Error; err != nil {
		return nil, err
	}
	return states, nil
}

func (r *locationRepository) GetStateByID(id uint, includeCities bool) (*models.State, error) {
	var state models.State
	query := database.DB.Preload("Country.Subregion.Region")

	if includeCities {
		query = query.Preload("Cities", "is_active = ?", true)
	}

	if err := query.First(&state, id).Error; err != nil {
		return nil, err
	}
	return &state, nil
}

func (r *locationRepository) CreateState(state *models.State) error {
	return database.DB.Create(state).Error
}

func (r *locationRepository) UpdateState(state *models.State) error {
	return database.DB.Save(state).Error
}

func (r *locationRepository) DeleteState(id uint) error {
	return database.DB.Delete(&models.State{}, id).Error
}

// ==================== CITIES ====================

func (r *locationRepository) GetAllCities(stateID *uint, search string) ([]models.City, error) {
	var cities []models.City
	query := database.DB.Where("is_active = ?", true).Preload("State.Country")

	if stateID != nil {
		query = query.Where("state_id = ?", *stateID)
	}

	if search != "" {
		query = query.Where("name ILIKE ?", "%"+search+"%")
	}

	if err := query.Order("name ASC").Find(&cities).Error; err != nil {
		return nil, err
	}
	return cities, nil
}

func (r *locationRepository) GetCityByID(id uint) (*models.City, error) {
	var city models.City
	query := database.DB.Preload("State.Country.Subregion.Region")

	if err := query.First(&city, id).Error; err != nil {
		return nil, err
	}
	return &city, nil
}

func (r *locationRepository) CreateCity(city *models.City) error {
	return database.DB.Create(city).Error
}

func (r *locationRepository) UpdateCity(city *models.City) error {
	return database.DB.Save(city).Error
}

func (r *locationRepository) DeleteCity(id uint) error {
	return database.DB.Delete(&models.City{}, id).Error
}

// ==================== HELPERS ====================

func (r *locationRepository) GetLocationHierarchy(countryID uint) (*models.Country, error) {
	var country models.Country
	if err := database.DB.
		Preload("Subregion.Region").
		Preload("States", "is_active = ?", true).
		Preload("States.Cities", "is_active = ?", true).
		First(&country, countryID).Error; err != nil {
		return nil, err
	}
	return &country, nil
}

func (r *locationRepository) SearchLocations(search string) (map[string]interface{}, error) {
	results := make(map[string]interface{})

	var countries []models.Country
	database.DB.Where("is_active = ? AND (name ILIKE ? OR iso2 ILIKE ? OR iso3 ILIKE ?)",
		true, "%"+search+"%", "%"+search+"%", "%"+search+"%").
		Preload("Subregion").Limit(10).Find(&countries)
	results["countries"] = countries

	var states []models.State
	database.DB.Where("is_active = ? AND name ILIKE ?", true, "%"+search+"%").
		Preload("Country").Limit(10).Find(&states)
	results["states"] = states

	var cities []models.City
	database.DB.Where("is_active = ? AND name ILIKE ?", true, "%"+search+"%").
		Preload("State.Country").Limit(10).Find(&cities)
	results["cities"] = cities

	return results, nil
}
