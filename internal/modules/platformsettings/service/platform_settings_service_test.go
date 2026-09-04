package service

import (
	"errors"
	"net/http"
	"testing"

	"dvra-api/internal/shared/apperr"
	"dvra-api/internal/shared/models"
)

// fakeRepo implementa el puerto sin base de datos.
type fakeRepo struct {
	settings *models.PlatformSettings
	saved    *models.PlatformSettings
	getErr   error
}

func (f *fakeRepo) Get() (*models.PlatformSettings, error) {
	if f.getErr != nil {
		return nil, f.getErr
	}
	return f.settings, nil
}

func (f *fakeRepo) GetCached() (*models.PlatformSettings, error) { return f.Get() }

func (f *fakeRepo) Update(s *models.PlatformSettings) error {
	f.saved = s
	return nil
}

func baseSettings() *models.PlatformSettings {
	return &models.PlatformSettings{
		PlatformName:     "DVRA ATS",
		PlatformShort:    "DVRA",
		SupportEmail:     "support@dvra.app",
		DefaultTrialDays: 14,
	}
}

func ptr[T any](v T) *T { return &v }

// Los campos que la marca necesita no pueden quedarse vacíos: la ruta pública
// los sirve al frontend antes del login, así que un nombre vacío deja la
// pantalla de acceso sin identidad.
func TestUpdateRechazaCamposDeMarcaVacios(t *testing.T) {
	casos := []struct {
		nombre string
		dto    UpdatePlatformSettingsDTO
	}{
		{"nombre", UpdatePlatformSettingsDTO{PlatformName: ptr("")}},
		{"nombre corto", UpdatePlatformSettingsDTO{PlatformShort: ptr("")}},
		{"email de soporte", UpdatePlatformSettingsDTO{SupportEmail: ptr("")}},
	}

	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			repo := &fakeRepo{settings: baseSettings()}
			svc := NewPlatformSettingsService(repo)

			_, err := svc.Update(&caso.dto, 1)
			if err == nil {
				t.Fatal("se esperaba error con el campo vacío")
			}
			if code := apperr.StatusCode(err); code != http.StatusBadRequest {
				t.Fatalf("status %d, se esperaba 400", code)
			}
			if repo.saved != nil {
				t.Fatal("no debía persistirse nada")
			}
		})
	}
}

func TestUpdateRechazaDiasDePruebaNegativos(t *testing.T) {
	repo := &fakeRepo{settings: baseSettings()}
	svc := NewPlatformSettingsService(repo)

	_, err := svc.Update(&UpdatePlatformSettingsDTO{DefaultTrialDays: ptr(-1)}, 1)
	if err == nil {
		t.Fatal("se esperaba error con días de prueba negativos")
	}
	if code := apperr.StatusCode(err); code != http.StatusBadRequest {
		t.Fatalf("status %d, se esperaba 400", code)
	}
	if repo.saved != nil {
		t.Fatal("no debía persistirse nada")
	}
}

// Una actualización parcial solo toca lo que viene: el resto se conserva.
func TestUpdateEsParcialYRegistraAutor(t *testing.T) {
	repo := &fakeRepo{settings: baseSettings()}
	svc := NewPlatformSettingsService(repo)

	if _, err := svc.Update(&UpdatePlatformSettingsDTO{PlatformName: ptr("Dvra")}, 42); err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if repo.saved == nil {
		t.Fatal("no se persistió la configuración")
	}
	if repo.saved.PlatformName != "Dvra" {
		t.Errorf("PlatformName = %q, se esperaba Dvra", repo.saved.PlatformName)
	}
	if repo.saved.SupportEmail != "support@dvra.app" {
		t.Errorf("se pisó un campo que no venía en el DTO: %q", repo.saved.SupportEmail)
	}
	if repo.saved.UpdatedByID == nil || *repo.saved.UpdatedByID != 42 {
		t.Error("no se registró quién hizo el cambio")
	}
}

func TestErrorDeRepositorioSePropaga(t *testing.T) {
	boom := errors.New("connection refused")
	svc := NewPlatformSettingsService(&fakeRepo{getErr: boom})

	if _, err := svc.Update(&UpdatePlatformSettingsDTO{PlatformName: ptr("Dvra")}, 1); !errors.Is(err, boom) {
		t.Fatalf("se esperaba el error del repositorio, se obtuvo %v", err)
	}
}
