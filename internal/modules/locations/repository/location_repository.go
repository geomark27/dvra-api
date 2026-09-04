package repository

import (
	"dvra-api/internal/shared/models"

	"gorm.io/gorm"
)

type LocationRepository struct{ db *gorm.DB }

// NewLocationRepository construye el adaptador de persistencia del catálogo
// geográfico. Recibe la conexión inyectada: sin estado global (ADR-001 regla 4).
func NewLocationRepository(db *gorm.DB) *LocationRepository {
	return &LocationRepository{db: db}
}

// ==================== REGIONS ====================

func (r *LocationRepository) GetAllRegions(includeSubregions bool) ([]models.Region, error) {
	var regions []models.Region
	query := r.db.Where("is_active = ?", true)

	if includeSubregions {
		query = query.Preload("Subregions", "is_active = ?", true)
	}

	if err := query.Order("name ASC").Find(&regions).Error; err != nil {
		return nil, err
	}
	return regions, nil
}

func (r *LocationRepository) GetRegionByID(id uint, includeSubregions bool) (*models.Region, error) {
	var region models.Region
	query := r.db

	if includeSubregions {
		query = query.Preload("Subregions", "is_active = ?", true)
	}

	if err := query.First(&region, id).Error; err != nil {
		return nil, err
	}
	return &region, nil
}

func (r *LocationRepository) CreateRegion(region *models.Region) error {
	return r.db.Create(region).Error
}

func (r *LocationRepository) UpdateRegion(region *models.Region) error {
	return r.db.Save(region).Error
}

func (r *LocationRepository) DeleteRegion(id uint) error {
	return r.db.Delete(&models.Region{}, id).Error
}

// ==================== SUBREGIONS ====================

func (r *LocationRepository) GetAllSubregions(regionID *uint) ([]models.Subregion, error) {
	var subregions []models.Subregion
	query := r.db.Where("is_active = ?", true).Preload("Region")

	if regionID != nil {
		query = query.Where("region_id = ?", *regionID)
	}

	if err := query.Order("name ASC").Find(&subregions).Error; err != nil {
		return nil, err
	}
	return subregions, nil
}

func (r *LocationRepository) GetSubregionByID(id uint, includeCountries bool) (*models.Subregion, error) {
	var subregion models.Subregion
	query := r.db.Preload("Region")

	if includeCountries {
		query = query.Preload("Countries", "is_active = ?", true)
	}

	if err := query.First(&subregion, id).Error; err != nil {
		return nil, err
	}
	return &subregion, nil
}

func (r *LocationRepository) CreateSubregion(subregion *models.Subregion) error {
	return r.db.Create(subregion).Error
}

func (r *LocationRepository) UpdateSubregion(subregion *models.Subregion) error {
	return r.db.Save(subregion).Error
}

func (r *LocationRepository) DeleteSubregion(id uint) error {
	return r.db.Delete(&models.Subregion{}, id).Error
}

// ==================== COUNTRIES ====================

func (r *LocationRepository) GetAllCountries(subregionID *uint, search string) ([]models.Country, error) {
	var countries []models.Country
	query := r.db.Where("is_active = ?", true).Preload("Subregion.Region")

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

func (r *LocationRepository) GetCountryByID(id uint, includeStates bool) (*models.Country, error) {
	var country models.Country
	query := r.db.Preload("Subregion.Region")

	if includeStates {
		query = query.Preload("States", "is_active = ?", true)
	}

	if err := query.First(&country, id).Error; err != nil {
		return nil, err
	}
	return &country, nil
}

func (r *LocationRepository) GetCountryByISO(iso string) (*models.Country, error) {
	var country models.Country
	query := r.db.Preload("Subregion.Region")

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

func (r *LocationRepository) CreateCountry(country *models.Country) error {
	return r.db.Create(country).Error
}

func (r *LocationRepository) UpdateCountry(country *models.Country) error {
	return r.db.Save(country).Error
}

func (r *LocationRepository) DeleteCountry(id uint) error {
	return r.db.Delete(&models.Country{}, id).Error
}

// ==================== STATES ====================

func (r *LocationRepository) GetAllStates(countryID *uint, search string) ([]models.State, error) {
	var states []models.State
	query := r.db.Where("is_active = ?", true).Preload("Country")

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

func (r *LocationRepository) GetStateByID(id uint, includeCities bool) (*models.State, error) {
	var state models.State
	query := r.db.Preload("Country.Subregion.Region")

	if includeCities {
		query = query.Preload("Cities", "is_active = ?", true)
	}

	if err := query.First(&state, id).Error; err != nil {
		return nil, err
	}
	return &state, nil
}

func (r *LocationRepository) CreateState(state *models.State) error {
	return r.db.Create(state).Error
}

func (r *LocationRepository) UpdateState(state *models.State) error {
	return r.db.Save(state).Error
}

func (r *LocationRepository) DeleteState(id uint) error {
	return r.db.Delete(&models.State{}, id).Error
}

// ==================== CITIES ====================

func (r *LocationRepository) GetAllCities(stateID *uint, search string) ([]models.City, error) {
	var cities []models.City
	query := r.db.Where("is_active = ?", true).Preload("State.Country")

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

func (r *LocationRepository) GetCityByID(id uint) (*models.City, error) {
	var city models.City
	query := r.db.Preload("State.Country.Subregion.Region")

	if err := query.First(&city, id).Error; err != nil {
		return nil, err
	}
	return &city, nil
}

func (r *LocationRepository) CreateCity(city *models.City) error {
	return r.db.Create(city).Error
}

func (r *LocationRepository) UpdateCity(city *models.City) error {
	return r.db.Save(city).Error
}

func (r *LocationRepository) DeleteCity(id uint) error {
	return r.db.Delete(&models.City{}, id).Error
}

// ==================== HELPERS ====================

func (r *LocationRepository) GetLocationHierarchy(countryID uint) (*models.Country, error) {
	var country models.Country
	if err := r.db.
		Preload("Subregion.Region").
		Preload("States", "is_active = ?", true).
		Preload("States.Cities", "is_active = ?", true).
		First(&country, countryID).Error; err != nil {
		return nil, err
	}
	return &country, nil
}

func (r *LocationRepository) SearchLocations(search string) (map[string]interface{}, error) {
	results := make(map[string]interface{})

	var countries []models.Country
	r.db.Where("is_active = ? AND (name ILIKE ? OR iso2 ILIKE ? OR iso3 ILIKE ?)",
		true, "%"+search+"%", "%"+search+"%", "%"+search+"%").
		Preload("Subregion").Limit(10).Find(&countries)
	results["countries"] = countries

	var states []models.State
	r.db.Where("is_active = ? AND name ILIKE ?", true, "%"+search+"%").
		Preload("Country").Limit(10).Find(&states)
	results["states"] = states

	var cities []models.City
	r.db.Where("is_active = ? AND name ILIKE ?", true, "%"+search+"%").
		Preload("State.Country").Limit(10).Find(&cities)
	results["cities"] = cities

	return results, nil
}
