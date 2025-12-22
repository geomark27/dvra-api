package dtos

// Response DTOs
type RegionDTO struct {
	ID         uint           `json:"id"`
	Name       string         `json:"name"`
	IsActive   bool           `json:"is_active"`
	Subregions []SubregionDTO `json:"subregions,omitempty"`
}

type SubregionDTO struct {
	ID        uint         `json:"id"`
	Name      string       `json:"name"`
	RegionID  uint         `json:"region_id"`
	IsActive  bool         `json:"is_active"`
	Region    *RegionDTO   `json:"region,omitempty"`
	Countries []CountryDTO `json:"countries,omitempty"`
}

type CountryDTO struct {
	ID          uint          `json:"id"`
	Name        string        `json:"name"`
	Iso2        string        `json:"iso2"`
	Iso3        string        `json:"iso3"`
	NumericCode string        `json:"numeric_code"`
	PhoneCode   string        `json:"phone_code"`
	Timezones   string        `json:"timezones,omitempty"`
	SubregionID *uint         `json:"subregion_id,omitempty"`
	IsActive    bool          `json:"is_active"`
	Subregion   *SubregionDTO `json:"subregion,omitempty"`
	States      []StateDTO    `json:"states,omitempty"`
}

type StateDTO struct {
	ID          uint        `json:"id"`
	Name        string      `json:"name"`
	CountryID   uint        `json:"country_id"`
	CountryCode string      `json:"country_code"`
	IsActive    bool        `json:"is_active"`
	Country     *CountryDTO `json:"country,omitempty"`
	Cities      []CityDTO   `json:"cities,omitempty"`
}

type CityDTO struct {
	ID       uint      `json:"id"`
	Name     string    `json:"name"`
	StateID  uint      `json:"state_id"`
	IsActive bool      `json:"is_active"`
	State    *StateDTO `json:"state,omitempty"`
}

// Create DTOs
type CreateRegionDTO struct {
	Name string `json:"name" binding:"required,min=2,max=255"`
}

type CreateSubregionDTO struct {
	Name     string `json:"name" binding:"required,min=2,max=255"`
	RegionID uint   `json:"region_id" binding:"required"`
}

type CreateCountryDTO struct {
	Name        string `json:"name" binding:"required,min=2,max=100"`
	Iso2        string `json:"iso2" binding:"required,len=2"`
	Iso3        string `json:"iso3" binding:"required,len=3"`
	NumericCode string `json:"numeric_code" binding:"required,len=3"`
	PhoneCode   string `json:"phone_code" binding:"required"`
	Timezones   string `json:"timezones"`
	SubregionID *uint  `json:"subregion_id"`
}

type CreateStateDTO struct {
	Name        string `json:"name" binding:"required,min=2,max=255"`
	CountryID   uint   `json:"country_id" binding:"required"`
	CountryCode string `json:"country_code" binding:"required,len=3"`
}

type CreateCityDTO struct {
	Name    string `json:"name" binding:"required,min=2,max=255"`
	StateID uint   `json:"state_id" binding:"required"`
}

// Update DTOs
type UpdateRegionDTO struct {
	Name     *string `json:"name,omitempty" binding:"omitempty,min=2,max=255"`
	IsActive *bool   `json:"is_active,omitempty"`
}

type UpdateSubregionDTO struct {
	Name     *string `json:"name,omitempty" binding:"omitempty,min=2,max=255"`
	RegionID *uint   `json:"region_id,omitempty"`
	IsActive *bool   `json:"is_active,omitempty"`
}

type UpdateCountryDTO struct {
	Name        *string `json:"name,omitempty" binding:"omitempty,min=2,max=100"`
	Iso2        *string `json:"iso2,omitempty" binding:"omitempty,len=2"`
	Iso3        *string `json:"iso3,omitempty" binding:"omitempty,len=3"`
	NumericCode *string `json:"numeric_code,omitempty" binding:"omitempty,len=3"`
	PhoneCode   *string `json:"phone_code,omitempty"`
	Timezones   *string `json:"timezones,omitempty"`
	SubregionID *uint   `json:"subregion_id,omitempty"`
	IsActive    *bool   `json:"is_active,omitempty"`
}

type UpdateStateDTO struct {
	Name        *string `json:"name,omitempty" binding:"omitempty,min=2,max=255"`
	CountryID   *uint   `json:"country_id,omitempty"`
	CountryCode *string `json:"country_code,omitempty" binding:"omitempty,len=3"`
	IsActive    *bool   `json:"is_active,omitempty"`
}

type UpdateCityDTO struct {
	Name     *string `json:"name,omitempty" binding:"omitempty,min=2,max=255"`
	StateID  *uint   `json:"state_id,omitempty"`
	IsActive *bool   `json:"is_active,omitempty"`
}

// Search Response
type LocationSearchDTO struct {
	Countries []CountryDTO `json:"countries"`
	States    []StateDTO   `json:"states"`
	Cities    []CityDTO    `json:"cities"`
}
