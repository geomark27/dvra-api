package service

import (
	"net/http"
	"testing"

	"dvra-api/internal/shared/apperr"
	"dvra-api/internal/shared/models"
)

// fakeRepo implementa el puerto sin base de datos.
type fakeRepo struct {
	SystemValueRepository // embebido: lo que el test no ejerce queda en nil

	values     map[uint]*models.SystemValue
	categories map[uint]*models.SystemValueCategory

	updated  *models.SystemValue
	deleted  []uint
	override *models.CompanySystemValueOverride
	created  *models.SystemValue
}

func (f *fakeRepo) GetByID(id uint) (*models.SystemValue, error) { return f.values[id], nil }

func (f *fakeRepo) GetCategoryByID(id uint) (*models.SystemValueCategory, error) {
	return f.categories[id], nil
}

func (f *fakeRepo) Update(v *models.SystemValue) error { f.updated = v; return nil }
func (f *fakeRepo) Delete(id uint) error               { f.deleted = append(f.deleted, id); return nil }
func (f *fakeRepo) Create(v *models.SystemValue) error { f.created = v; return nil }

func (f *fakeRepo) UpsertOverride(o *models.CompanySystemValueOverride) error {
	f.override = o
	return nil
}

func uptr(v uint) *uint { return &v }

const (
	empresaA = uint(10)
	empresaB = uint(20)
)

// ---------------------------------------------------------------------------
// AISLAMIENTO MULTI-TENANT (ESTANDARES §4). Es el test que importa: sin él, un
// tenant podría editar o borrar el catálogo de otro conociendo un id.
// ---------------------------------------------------------------------------

func repoConValores() *fakeRepo {
	return &fakeRepo{values: map[uint]*models.SystemValue{
		1: {BaseModel: models.BaseModel{ID: 1}, CompanyID: uptr(empresaA), Label: "de A"},
		2: {BaseModel: models.BaseModel{ID: 2}, CompanyID: uptr(empresaB), Label: "de B"},
		3: {BaseModel: models.BaseModel{ID: 3}, CompanyID: nil, Label: "global"},
	}}
}

func TestUpdateTenantNoAlcanzaValoresAjenos(t *testing.T) {
	casos := []struct {
		nombre string
		id     uint
	}{
		{"valor de otra empresa", 2},
		{"valor global", 3},
		{"valor inexistente", 99},
	}

	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			repo := repoConValores()
			svc := NewSystemValueService(repo)

			nuevaEtiqueta := "secuestrado"
			_, err := svc.UpdateTenant(empresaA, caso.id, UpdateSystemValueDTO{Label: &nuevaEtiqueta})
			if err == nil {
				t.Fatal("se esperaba error: la empresa A no debe alcanzar este valor")
			}
			// 404 y no 403: un 403 confirmaría que el valor existe.
			if code := apperr.StatusCode(err); code != http.StatusNotFound {
				t.Fatalf("status %d, se esperaba 404", code)
			}
			if repo.updated != nil {
				t.Fatal("no debía modificarse nada")
			}
		})
	}
}

func TestDeleteTenantNoAlcanzaValoresAjenos(t *testing.T) {
	for _, id := range []uint{2, 3, 99} {
		repo := repoConValores()
		svc := NewSystemValueService(repo)

		if err := svc.DeleteTenant(empresaA, id); err == nil {
			t.Fatalf("id %d: se esperaba error", id)
		}
		if len(repo.deleted) != 0 {
			t.Fatalf("id %d: no debía borrarse nada", id)
		}
	}
}

func TestTenantSiEditaSuPropioValor(t *testing.T) {
	repo := repoConValores()
	svc := NewSystemValueService(repo)

	nuevaEtiqueta := "Prueba cultural"
	if _, err := svc.UpdateTenant(empresaA, 1, UpdateSystemValueDTO{Label: &nuevaEtiqueta}); err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if repo.updated == nil || repo.updated.Label != nuevaEtiqueta {
		t.Fatal("no se aplicó el cambio sobre el valor propio")
	}
}

// ---------------------------------------------------------------------------
// EXTENSIBILIDAD DEL CATÁLOGO (ADR-003)
// ---------------------------------------------------------------------------

// Una categoría cerrada no admite valores propios del tenant: dejar pasar esto
// rompería supuestos del dominio que la usa (p. ej. los estados de suscripción,
// que billing compara contra constantes).
//
// Hacen falta AMBOS flags: IsExtensible marca que la categoría admite crecer, y
// AllowTenantValues que quien la hace crecer puede ser el tenant. Una categoría
// extensible solo por el landlord tiene el primero y no el segundo.
func TestCreateTenantRespetaCategoriaCerrada(t *testing.T) {
	casos := []struct {
		nombre     string
		extensible bool
		permite    bool
	}{
		{"ni extensible ni abierta al tenant", false, false},
		{"extensible pero solo para el landlord", true, false},
		{"abierta al tenant pero no extensible", false, true},
	}

	for _, caso := range casos {
		t.Run(caso.nombre, func(t *testing.T) {
			repo := &fakeRepo{categories: map[uint]*models.SystemValueCategory{
				5: {BaseModel: models.BaseModel{ID: 5}, IsExtensible: caso.extensible, AllowTenantValues: caso.permite},
			}}
			svc := NewSystemValueService(repo)

			_, err := svc.CreateTenant(empresaA, CreateSystemValueDTO{CategoryID: 5, Value: "x", Label: "X"})
			if err == nil {
				t.Fatal("se esperaba error: la categoría no admite valores del tenant")
			}
			if code := apperr.StatusCode(err); code != http.StatusForbidden {
				t.Fatalf("status %d, se esperaba 403", code)
			}
			if repo.created != nil {
				t.Fatal("no debía crearse el valor")
			}
		})
	}
}

func TestCreateTenantFuerzaLaEmpresaDelToken(t *testing.T) {
	repo := &fakeRepo{categories: map[uint]*models.SystemValueCategory{
		5: {BaseModel: models.BaseModel{ID: 5}, IsExtensible: true, AllowTenantValues: true},
	}}
	svc := NewSystemValueService(repo)

	if _, err := svc.CreateTenant(empresaA, CreateSystemValueDTO{CategoryID: 5, Value: "cultural", Label: "Prueba cultural"}); err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if repo.created == nil {
		t.Fatal("no se creó el valor")
	}
	if repo.created.CompanyID == nil || *repo.created.CompanyID != empresaA {
		t.Error("el valor debe quedar atado a la empresa del token")
	}
	if repo.created.IsSystem {
		t.Error("un valor de tenant no es del sistema")
	}
}

// Un override solo aplica sobre un valor global y solo si la categoría lo
// permite.
func TestUpsertOverrideExigeValorGlobalYCategoriaQueLoPermita(t *testing.T) {
	t.Run("sobre un valor de tenant", func(t *testing.T) {
		repo := repoConValores()
		svc := NewSystemValueService(repo)

		if err := svc.UpsertTenantOverride(empresaA, 1, UpsertSystemValueOverrideDTO{}); err == nil {
			t.Fatal("se esperaba error: no se sobrescribe un valor de tenant")
		}
		if repo.override != nil {
			t.Fatal("no debía guardarse el override")
		}
	})

	t.Run("categoría que no lo permite", func(t *testing.T) {
		repo := &fakeRepo{values: map[uint]*models.SystemValue{
			3: {BaseModel: models.BaseModel{ID: 3}, CompanyID: nil,
				Category: &models.SystemValueCategory{AllowTenantOverride: false}},
		}}
		svc := NewSystemValueService(repo)

		err := svc.UpsertTenantOverride(empresaA, 3, UpsertSystemValueOverrideDTO{})
		if code := apperr.StatusCode(err); code != http.StatusForbidden {
			t.Fatalf("status %d, se esperaba 403", code)
		}
		if repo.override != nil {
			t.Fatal("no debía guardarse el override")
		}
	})

	t.Run("caso permitido", func(t *testing.T) {
		repo := &fakeRepo{values: map[uint]*models.SystemValue{
			3: {BaseModel: models.BaseModel{ID: 3}, CompanyID: nil,
				Category: &models.SystemValueCategory{AllowTenantOverride: true}},
		}}
		svc := NewSystemValueService(repo)

		if err := svc.UpsertTenantOverride(empresaA, 3, UpsertSystemValueOverrideDTO{}); err != nil {
			t.Fatalf("error inesperado: %v", err)
		}
		if repo.override == nil || repo.override.CompanyID != empresaA {
			t.Fatal("el override debe quedar atado a la empresa del token")
		}
	})
}
