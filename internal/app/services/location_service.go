package services

import (
	"dvra-api/internal/app/dtos"
	"dvra-api/internal/app/models"
	"dvra-api/internal/app/repositories"
	"errors"
	"strings"
)

type LocationService interface {
	// Regions
	GetAllRegions(includeSubregions bool) ([]dtos.RegionDTO, error)
	GetRegionByID(id uint, includeSubregions bool) (*dtos.RegionDTO, error)
	CreateRegion(dto dtos.CreateRegionDTO) (*dtos.RegionDTO, error)
	UpdateRegion(id uint, dto dtos.UpdateRegionDTO) (*dtos.RegionDTO, error)
	DeleteRegion(id uint) error

	// Subregions
	GetAllSubregions(regionID *uint) ([]dtos.SubregionDTO, error)
	GetSubregionByID(id uint, includeCountries bool) (*dtos.SubregionDTO, error)
	CreateSubregion(dto dtos.CreateSubregionDTO) (*dtos.SubregionDTO, error)
	UpdateSubregion(id uint, dto dtos.UpdateSubregionDTO) (*dtos.SubregionDTO, error)
	DeleteSubregion(id uint) error

	// Countries
	GetAllCountries(subregionID *uint, search string) ([]dtos.CountryDTO, error)
	GetCountryByID(id uint, includeStates bool) (*dtos.CountryDTO, error)
	GetCountryByISO(iso string) (*dtos.CountryDTO, error)
	CreateCountry(dto dtos.CreateCountryDTO) (*dtos.CountryDTO, error)
	UpdateCountry(id uint, dto dtos.UpdateCountryDTO) (*dtos.CountryDTO, error)
	DeleteCountry(id uint) error

	// States
	GetAllStates(countryID *uint, search string) ([]dtos.StateDTO, error)
	GetStateByID(id uint, includeCities bool) (*dtos.StateDTO, error)
	CreateState(dto dtos.CreateStateDTO) (*dtos.StateDTO, error)
	UpdateState(id uint, dto dtos.UpdateStateDTO) (*dtos.StateDTO, error)
	DeleteState(id uint) error

	// Cities
	GetAllCities(stateID *uint, search string) ([]dtos.CityDTO, error)
	GetCityByID(id uint) (*dtos.CityDTO, error)
	CreateCity(dto dtos.CreateCityDTO) (*dtos.CityDTO, error)
	UpdateCity(id uint, dto dtos.UpdateCityDTO) (*dtos.CityDTO, error)
	DeleteCity(id uint) error

	// Helpers
	GetLocationHierarchy(countryID uint) (*dtos.CountryDTO, error)
	SearchLocations(search string) (*dtos.LocationSearchDTO, error)
}

type locationService struct {
	repo repositories.LocationRepository
}

func NewLocationService(repo repositories.LocationRepository) LocationService {
	return &locationService{repo: repo}
}

// ==================== REGIONS ====================

func (s *locationService) GetAllRegions(includeSubregions bool) ([]dtos.RegionDTO, error) {
	regions, err := s.repo.GetAllRegions(includeSubregions)
	if err != nil {
		return nil, err
	}

	result := make([]dtos.RegionDTO, len(regions))
	for i, region := range regions {
		result[i] = s.regionToDTO(&region, includeSubregions)
	}
	return result, nil
}

func (s *locationService) GetRegionByID(id uint, includeSubregions bool) (*dtos.RegionDTO, error) {
	region, err := s.repo.GetRegionByID(id, includeSubregions)
	if err != nil {
		return nil, err
	}
	if region == nil {
		return nil, errors.New("region not found")
	}

	dto := s.regionToDTO(region, includeSubregions)
	return &dto, nil
}

func (s *locationService) CreateRegion(dto dtos.CreateRegionDTO) (*dtos.RegionDTO, error) {
	region := &models.Region{
		Name:     dto.Name,
		IsActive: true,
	}

	if err := s.repo.CreateRegion(region); err != nil {
		return nil, err
	}

	result := s.regionToDTO(region, false)
	return &result, nil
}

func (s *locationService) UpdateRegion(id uint, dto dtos.UpdateRegionDTO) (*dtos.RegionDTO, error) {
	region, err := s.repo.GetRegionByID(id, false)
	if err != nil {
		return nil, err
	}
	if region == nil {
		return nil, errors.New("region not found")
	}

	if dto.Name != nil {
		region.Name = *dto.Name
	}
	if dto.IsActive != nil {
		region.IsActive = *dto.IsActive
	}

	if err := s.repo.UpdateRegion(region); err != nil {
		return nil, err
	}

	result := s.regionToDTO(region, false)
	return &result, nil
}

func (s *locationService) DeleteRegion(id uint) error {
	return s.repo.DeleteRegion(id)
}

// ==================== SUBREGIONS ====================

func (s *locationService) GetAllSubregions(regionID *uint) ([]dtos.SubregionDTO, error) {
	subregions, err := s.repo.GetAllSubregions(regionID)
	if err != nil {
		return nil, err
	}

	result := make([]dtos.SubregionDTO, len(subregions))
	for i, subregion := range subregions {
		result[i] = s.subregionToDTO(&subregion, false, true)
	}
	return result, nil
}

func (s *locationService) GetSubregionByID(id uint, includeCountries bool) (*dtos.SubregionDTO, error) {
	subregion, err := s.repo.GetSubregionByID(id, includeCountries)
	if err != nil {
		return nil, err
	}
	if subregion == nil {
		return nil, errors.New("subregion not found")
	}

	dto := s.subregionToDTO(subregion, includeCountries, true)
	return &dto, nil
}

func (s *locationService) CreateSubregion(dto dtos.CreateSubregionDTO) (*dtos.SubregionDTO, error) {
	// Verificar que la región existe
	region, err := s.repo.GetRegionByID(dto.RegionID, false)
	if err != nil || region == nil {
		return nil, errors.New("region not found")
	}

	subregion := &models.Subregion{
		Name:     dto.Name,
		RegionID: dto.RegionID,
		IsActive: true,
	}

	if err := s.repo.CreateSubregion(subregion); err != nil {
		return nil, err
	}

	// Recargar con relaciones
	subregion, _ = s.repo.GetSubregionByID(subregion.ID, false)
	result := s.subregionToDTO(subregion, false, true)
	return &result, nil
}

func (s *locationService) UpdateSubregion(id uint, dto dtos.UpdateSubregionDTO) (*dtos.SubregionDTO, error) {
	subregion, err := s.repo.GetSubregionByID(id, false)
	if err != nil {
		return nil, err
	}
	if subregion == nil {
		return nil, errors.New("subregion not found")
	}

	if dto.Name != nil {
		subregion.Name = *dto.Name
	}
	if dto.RegionID != nil {
		subregion.RegionID = *dto.RegionID
	}
	if dto.IsActive != nil {
		subregion.IsActive = *dto.IsActive
	}

	if err := s.repo.UpdateSubregion(subregion); err != nil {
		return nil, err
	}

	subregion, _ = s.repo.GetSubregionByID(id, false)
	result := s.subregionToDTO(subregion, false, true)
	return &result, nil
}

func (s *locationService) DeleteSubregion(id uint) error {
	return s.repo.DeleteSubregion(id)
}

// ==================== COUNTRIES ====================

func (s *locationService) GetAllCountries(subregionID *uint, search string) ([]dtos.CountryDTO, error) {
	countries, err := s.repo.GetAllCountries(subregionID, search)
	if err != nil {
		return nil, err
	}

	result := make([]dtos.CountryDTO, len(countries))
	for i, country := range countries {
		result[i] = s.countryToDTO(&country, false, true)
	}
	return result, nil
}

func (s *locationService) GetCountryByID(id uint, includeStates bool) (*dtos.CountryDTO, error) {
	country, err := s.repo.GetCountryByID(id, includeStates)
	if err != nil {
		return nil, err
	}
	if country == nil {
		return nil, errors.New("country not found")
	}

	dto := s.countryToDTO(country, includeStates, true)
	return &dto, nil
}

func (s *locationService) GetCountryByISO(iso string) (*dtos.CountryDTO, error) {
	country, err := s.repo.GetCountryByISO(strings.ToUpper(iso))
	if err != nil {
		return nil, err
	}
	if country == nil {
		return nil, errors.New("country not found")
	}

	dto := s.countryToDTO(country, false, true)
	return &dto, nil
}

func (s *locationService) CreateCountry(dto dtos.CreateCountryDTO) (*dtos.CountryDTO, error) {
	country := &models.Country{
		Name:        dto.Name,
		Iso2:        strings.ToUpper(dto.Iso2),
		Iso3:        strings.ToUpper(dto.Iso3),
		NumericCode: dto.NumericCode,
		PhoneCode:   dto.PhoneCode,
		Timezones:   dto.Timezones,
		SubregionID: dto.SubregionID,
		IsActive:    true,
	}

	if err := s.repo.CreateCountry(country); err != nil {
		return nil, err
	}

	country, _ = s.repo.GetCountryByID(country.ID, false)
	result := s.countryToDTO(country, false, true)
	return &result, nil
}

func (s *locationService) UpdateCountry(id uint, dto dtos.UpdateCountryDTO) (*dtos.CountryDTO, error) {
	country, err := s.repo.GetCountryByID(id, false)
	if err != nil {
		return nil, err
	}
	if country == nil {
		return nil, errors.New("country not found")
	}

	if dto.Name != nil {
		country.Name = *dto.Name
	}
	if dto.Iso2 != nil {
		country.Iso2 = strings.ToUpper(*dto.Iso2)
	}
	if dto.Iso3 != nil {
		country.Iso3 = strings.ToUpper(*dto.Iso3)
	}
	if dto.NumericCode != nil {
		country.NumericCode = *dto.NumericCode
	}
	if dto.PhoneCode != nil {
		country.PhoneCode = *dto.PhoneCode
	}
	if dto.Timezones != nil {
		country.Timezones = *dto.Timezones
	}
	if dto.SubregionID != nil {
		country.SubregionID = dto.SubregionID
	}
	if dto.IsActive != nil {
		country.IsActive = *dto.IsActive
	}

	if err := s.repo.UpdateCountry(country); err != nil {
		return nil, err
	}

	country, _ = s.repo.GetCountryByID(id, false)
	result := s.countryToDTO(country, false, true)
	return &result, nil
}

func (s *locationService) DeleteCountry(id uint) error {
	return s.repo.DeleteCountry(id)
}

// ==================== STATES ====================

func (s *locationService) GetAllStates(countryID *uint, search string) ([]dtos.StateDTO, error) {
	states, err := s.repo.GetAllStates(countryID, search)
	if err != nil {
		return nil, err
	}

	result := make([]dtos.StateDTO, len(states))
	for i, state := range states {
		result[i] = s.stateToDTO(&state, false, true)
	}
	return result, nil
}

func (s *locationService) GetStateByID(id uint, includeCities bool) (*dtos.StateDTO, error) {
	state, err := s.repo.GetStateByID(id, includeCities)
	if err != nil {
		return nil, err
	}
	if state == nil {
		return nil, errors.New("state not found")
	}

	dto := s.stateToDTO(state, includeCities, true)
	return &dto, nil
}

func (s *locationService) CreateState(dto dtos.CreateStateDTO) (*dtos.StateDTO, error) {
	state := &models.State{
		Name:        dto.Name,
		CountryID:   dto.CountryID,
		CountryCode: strings.ToUpper(dto.CountryCode),
		IsActive:    true,
	}

	if err := s.repo.CreateState(state); err != nil {
		return nil, err
	}

	state, _ = s.repo.GetStateByID(state.ID, false)
	result := s.stateToDTO(state, false, true)
	return &result, nil
}

func (s *locationService) UpdateState(id uint, dto dtos.UpdateStateDTO) (*dtos.StateDTO, error) {
	state, err := s.repo.GetStateByID(id, false)
	if err != nil {
		return nil, err
	}
	if state == nil {
		return nil, errors.New("state not found")
	}

	if dto.Name != nil {
		state.Name = *dto.Name
	}
	if dto.CountryID != nil {
		state.CountryID = *dto.CountryID
	}
	if dto.CountryCode != nil {
		state.CountryCode = strings.ToUpper(*dto.CountryCode)
	}
	if dto.IsActive != nil {
		state.IsActive = *dto.IsActive
	}

	if err := s.repo.UpdateState(state); err != nil {
		return nil, err
	}

	state, _ = s.repo.GetStateByID(id, false)
	result := s.stateToDTO(state, false, true)
	return &result, nil
}

func (s *locationService) DeleteState(id uint) error {
	return s.repo.DeleteState(id)
}

// ==================== CITIES ====================

func (s *locationService) GetAllCities(stateID *uint, search string) ([]dtos.CityDTO, error) {
	cities, err := s.repo.GetAllCities(stateID, search)
	if err != nil {
		return nil, err
	}

	result := make([]dtos.CityDTO, len(cities))
	for i, city := range cities {
		result[i] = s.cityToDTO(&city, true)
	}
	return result, nil
}

func (s *locationService) GetCityByID(id uint) (*dtos.CityDTO, error) {
	city, err := s.repo.GetCityByID(id)
	if err != nil {
		return nil, err
	}
	if city == nil {
		return nil, errors.New("city not found")
	}

	dto := s.cityToDTO(city, true)
	return &dto, nil
}

func (s *locationService) CreateCity(dto dtos.CreateCityDTO) (*dtos.CityDTO, error) {
	city := &models.City{
		Name:     dto.Name,
		StateID:  dto.StateID,
		IsActive: true,
	}

	if err := s.repo.CreateCity(city); err != nil {
		return nil, err
	}

	city, _ = s.repo.GetCityByID(city.ID)
	result := s.cityToDTO(city, true)
	return &result, nil
}

func (s *locationService) UpdateCity(id uint, dto dtos.UpdateCityDTO) (*dtos.CityDTO, error) {
	city, err := s.repo.GetCityByID(id)
	if err != nil {
		return nil, err
	}
	if city == nil {
		return nil, errors.New("city not found")
	}

	if dto.Name != nil {
		city.Name = *dto.Name
	}
	if dto.StateID != nil {
		city.StateID = *dto.StateID
	}
	if dto.IsActive != nil {
		city.IsActive = *dto.IsActive
	}

	if err := s.repo.UpdateCity(city); err != nil {
		return nil, err
	}

	city, _ = s.repo.GetCityByID(id)
	result := s.cityToDTO(city, true)
	return &result, nil
}

func (s *locationService) DeleteCity(id uint) error {
	return s.repo.DeleteCity(id)
}

// ==================== HELPERS ====================

func (s *locationService) GetLocationHierarchy(countryID uint) (*dtos.CountryDTO, error) {
	country, err := s.repo.GetLocationHierarchy(countryID)
	if err != nil {
		return nil, err
	}

	dto := s.countryToDTO(country, true, true)
	return &dto, nil
}

func (s *locationService) SearchLocations(search string) (*dtos.LocationSearchDTO, error) {
	results, err := s.repo.SearchLocations(search)
	if err != nil {
		return nil, err
	}

	dto := &dtos.LocationSearchDTO{
		Countries: []dtos.CountryDTO{},
		States:    []dtos.StateDTO{},
		Cities:    []dtos.CityDTO{},
	}

	if countries, ok := results["countries"].([]models.Country); ok {
		for _, c := range countries {
			dto.Countries = append(dto.Countries, s.countryToDTO(&c, false, true))
		}
	}

	if states, ok := results["states"].([]models.State); ok {
		for _, st := range states {
			dto.States = append(dto.States, s.stateToDTO(&st, false, true))
		}
	}

	if cities, ok := results["cities"].([]models.City); ok {
		for _, city := range cities {
			dto.Cities = append(dto.Cities, s.cityToDTO(&city, true))
		}
	}

	return dto, nil
}

// ==================== MAPPERS ====================

func (s *locationService) regionToDTO(region *models.Region, includeSubregions bool) dtos.RegionDTO {
	dto := dtos.RegionDTO{
		ID:       region.ID,
		Name:     region.Name,
		IsActive: region.IsActive,
	}

	if includeSubregions && len(region.Subregions) > 0 {
		dto.Subregions = make([]dtos.SubregionDTO, len(region.Subregions))
		for i, sub := range region.Subregions {
			dto.Subregions[i] = s.subregionToDTO(&sub, false, false)
		}
	}

	return dto
}

func (s *locationService) subregionToDTO(subregion *models.Subregion, includeCountries bool, includeRegion bool) dtos.SubregionDTO {
	dto := dtos.SubregionDTO{
		ID:       subregion.ID,
		Name:     subregion.Name,
		RegionID: subregion.RegionID,
		IsActive: subregion.IsActive,
	}

	if includeRegion && subregion.Region.ID > 0 {
		regionDTO := s.regionToDTO(&subregion.Region, false)
		dto.Region = &regionDTO
	}

	if includeCountries && len(subregion.Countries) > 0 {
		dto.Countries = make([]dtos.CountryDTO, len(subregion.Countries))
		for i, c := range subregion.Countries {
			dto.Countries[i] = s.countryToDTO(&c, false, false)
		}
	}

	return dto
}

func (s *locationService) countryToDTO(country *models.Country, includeStates bool, includeSubregion bool) dtos.CountryDTO {
	dto := dtos.CountryDTO{
		ID:          country.ID,
		Name:        country.Name,
		Iso2:        country.Iso2,
		Iso3:        country.Iso3,
		NumericCode: country.NumericCode,
		PhoneCode:   country.PhoneCode,
		Timezones:   country.Timezones,
		SubregionID: country.SubregionID,
		IsActive:    country.IsActive,
	}

	if includeSubregion && country.Subregion != nil && country.Subregion.ID > 0 {
		subDTO := s.subregionToDTO(country.Subregion, false, true)
		dto.Subregion = &subDTO
	}

	if includeStates && len(country.States) > 0 {
		dto.States = make([]dtos.StateDTO, len(country.States))
		for i, st := range country.States {
			dto.States[i] = s.stateToDTO(&st, true, false)
		}
	}

	return dto
}

func (s *locationService) stateToDTO(state *models.State, includeCities bool, includeCountry bool) dtos.StateDTO {
	dto := dtos.StateDTO{
		ID:          state.ID,
		Name:        state.Name,
		CountryID:   state.CountryID,
		CountryCode: state.CountryCode,
		IsActive:    state.IsActive,
	}

	if includeCountry && state.Country.ID > 0 {
		countryDTO := s.countryToDTO(&state.Country, false, true)
		dto.Country = &countryDTO
	}

	if includeCities && len(state.Cities) > 0 {
		dto.Cities = make([]dtos.CityDTO, len(state.Cities))
		for i, city := range state.Cities {
			dto.Cities[i] = s.cityToDTO(&city, false)
		}
	}

	return dto
}

func (s *locationService) cityToDTO(city *models.City, includeState bool) dtos.CityDTO {
	dto := dtos.CityDTO{
		ID:       city.ID,
		Name:     city.Name,
		StateID:  city.StateID,
		IsActive: city.IsActive,
	}

	if includeState && city.State.ID > 0 {
		stateDTO := s.stateToDTO(&city.State, false, true)
		dto.State = &stateDTO
	}

	return dto
}
