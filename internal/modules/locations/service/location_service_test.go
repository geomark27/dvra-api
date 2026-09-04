package service

import (
	"errors"
	"net/http"
	"testing"

	"dvra-api/internal/shared/apperr"
	"dvra-api/internal/shared/models"
)

// fakeRepo implementa el puerto LocationRepository con lo justo para probar la
// lógica del service sin base de datos. Los métodos no usados devuelven cero.
type fakeRepo struct {
	LocationRepository // embebido: deja en nil lo que el test no ejerce

	regions   map[uint]*models.Region
	countries map[string]*models.Country

	createdCountry   *models.Country
	createdSubregion *models.Subregion
}

func (f *fakeRepo) GetRegionByID(id uint, _ bool) (*models.Region, error) {
	return f.regions[id], nil
}

func (f *fakeRepo) GetCountryByISO(iso string) (*models.Country, error) {
	return f.countries[iso], nil
}

func (f *fakeRepo) CreateCountry(country *models.Country) error {
	f.createdCountry = country
	return nil
}

func (f *fakeRepo) CreateSubregion(subregion *models.Subregion) error {
	f.createdSubregion = subregion
	return nil
}

func (f *fakeRepo) GetSubregionByID(uint, bool) (*models.Subregion, error) {
	return &models.Subregion{}, nil
}

func (f *fakeRepo) GetCountryByID(uint, bool) (*models.Country, error) {
	return &models.Country{}, nil
}

// El catálogo se consulta por ISO en mayúsculas: quien llama puede mandar "ec"
// y debe encontrar "EC".
func TestGetCountryByISONormalizaMayusculas(t *testing.T) {
	repo := &fakeRepo{countries: map[string]*models.Country{
		"EC": {Name: "Ecuador", Iso2: "EC"},
	}}
	svc := NewLocationService(repo)

	country, err := svc.GetCountryByISO("ec")
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if country.Name != "Ecuador" {
		t.Fatalf("devolvió %q, se esperaba Ecuador", country.Name)
	}
}

func TestGetCountryByISODevuelveNotFound(t *testing.T) {
	svc := NewLocationService(&fakeRepo{countries: map[string]*models.Country{}})

	_, err := svc.GetCountryByISO("zz")
	if err == nil {
		t.Fatal("se esperaba error para un ISO inexistente")
	}
	if code := apperr.StatusCode(err); code != http.StatusNotFound {
		t.Fatalf("status %d, se esperaba 404", code)
	}
}

// Los códigos ISO se persisten siempre en mayúsculas, venga como venga la
// entrada: de lo contrario la búsqueda por ISO no encontraría el registro.
func TestCreateCountryNormalizaCodigosISO(t *testing.T) {
	repo := &fakeRepo{}
	svc := NewLocationService(repo)

	if _, err := svc.CreateCountry(CreateCountryDTO{
		Name: "Ecuador",
		Iso2: "ec",
		Iso3: "ecu",
	}); err != nil {
		t.Fatalf("error inesperado: %v", err)
	}

	if repo.createdCountry == nil {
		t.Fatal("no se persistió el país")
	}
	if repo.createdCountry.Iso2 != "EC" {
		t.Errorf("Iso2 = %q, se esperaba EC", repo.createdCountry.Iso2)
	}
	if repo.createdCountry.Iso3 != "ECU" {
		t.Errorf("Iso3 = %q, se esperaba ECU", repo.createdCountry.Iso3)
	}
}

// Integridad referencial: no se crea una subregión colgando de una región que
// no existe.
func TestCreateSubregionExigeRegionExistente(t *testing.T) {
	repo := &fakeRepo{regions: map[uint]*models.Region{}}
	svc := NewLocationService(repo)

	_, err := svc.CreateSubregion(CreateSubregionDTO{Name: "Andina", RegionID: 99})
	if err == nil {
		t.Fatal("se esperaba error al referenciar una región inexistente")
	}
	if code := apperr.StatusCode(err); code != http.StatusNotFound {
		t.Fatalf("status %d, se esperaba 404", code)
	}
	if repo.createdSubregion != nil {
		t.Fatal("no debía persistirse la subregión")
	}
}

func TestCreateSubregionConRegionValida(t *testing.T) {
	repo := &fakeRepo{regions: map[uint]*models.Region{
		1: {BaseModel: models.BaseModel{ID: 1}, Name: "América"},
	}}
	svc := NewLocationService(repo)

	if _, err := svc.CreateSubregion(CreateSubregionDTO{Name: "Andina", RegionID: 1}); err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if repo.createdSubregion == nil {
		t.Fatal("no se persistió la subregión")
	}
	if repo.createdSubregion.RegionID != 1 {
		t.Errorf("RegionID = %d, se esperaba 1", repo.createdSubregion.RegionID)
	}
	if !repo.createdSubregion.IsActive {
		t.Error("una subregión nueva nace activa")
	}
}

// Un fallo de la base se propaga tal cual: no se disfraza de "no encontrado",
// que llevaría al cliente a creer que el recurso no existe.
func TestErrorDeRepositorioSePropaga(t *testing.T) {
	boom := errors.New("connection refused")
	repo := &failingRepo{err: boom}
	svc := NewLocationService(repo)

	_, err := svc.GetAllRegions(false)
	if !errors.Is(err, boom) {
		t.Fatalf("se esperaba el error del repositorio, se obtuvo %v", err)
	}
}

type failingRepo struct {
	LocationRepository
	err error
}

func (f *failingRepo) GetAllRegions(bool) ([]models.Region, error) { return nil, f.err }
