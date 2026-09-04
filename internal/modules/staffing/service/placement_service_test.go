package service

import (
	"net/http"
	"testing"

	"dvra-api/internal/modules/staffing/domain"
	"dvra-api/internal/shared/apperr"
	"dvra-api/internal/shared/models"
)

const (
	empresaA = uint(10)
	empresaB = uint(20)
)

type fakePlacementRepo struct {
	domain.PlacementRepository
	existe    bool
	creado    *models.Placement
	createErr error
}

func (f *fakePlacementRepo) ExistsByApplicationID(uint) (bool, error) { return f.existe, nil }
func (f *fakePlacementRepo) Create(p *models.Placement) (*models.Placement, error) {
	f.creado = p
	return p, f.createErr
}

type fakeClientRepo struct {
	domain.StaffingClientRepository
	clients map[uint]*models.StaffingClient
}

func (f *fakeClientRepo) GetByID(id uint) (*models.StaffingClient, error) { return f.clients[id], nil }

type fakeApps struct {
	apps map[uint]*domain.HiredApplication
}

func (f fakeApps) FindByID(id uint) (*domain.HiredApplication, error) { return f.apps[id], nil }

func nuevoServicio(repo *fakePlacementRepo, clienteEmpresa uint, app *domain.HiredApplication) *PlacementService {
	clients := &fakeClientRepo{clients: map[uint]*models.StaffingClient{
		1: {BaseModel: models.BaseModel{ID: 1}, CompanyID: clienteEmpresa},
	}}
	apps := map[uint]*domain.HiredApplication{}
	if app != nil {
		apps[app.ID] = app
	}
	return NewPlacementService(repo, clients, fakeApps{apps: apps})
}

func appEn(etapa string, empresa uint) *domain.HiredApplication {
	return &domain.HiredApplication{ID: 5, CompanyID: empresa, CandidateID: 77, JobID: 88, Stage: etapa}
}

// Una colocación solo nace de una postulación efectivamente contratada. Es la
// regla de negocio del módulo: sin ella se facturaría a un cliente final por
// alguien que nunca se contrató.
func TestCreateExigeEtapaHired(t *testing.T) {
	for _, etapa := range []string{"applied", "screening", "interview", "offer", "rejected"} {
		repo := &fakePlacementRepo{}
		svc := nuevoServicio(repo, empresaA, appEn(etapa, empresaA))

		_, err := svc.Create(empresaA, CreatePlacementDTO{StaffingClientID: 1, ApplicationID: 5})
		if err == nil {
			t.Fatalf("etapa %q: se esperaba error", etapa)
		}
		if code := apperr.StatusCode(err); code != http.StatusBadRequest {
			t.Fatalf("etapa %q: status %d, se esperaba 400", etapa, code)
		}
		if repo.creado != nil {
			t.Fatalf("etapa %q: no debía crearse la colocación", etapa)
		}
	}
}

// Integridad cross-tenant (ESTANDARES §4.4): ni el cliente final ni la
// postulación pueden ser de otra empresa.
func TestCreateValidaPertenenciaCrossTenant(t *testing.T) {
	t.Run("cliente final de otra empresa", func(t *testing.T) {
		repo := &fakePlacementRepo{}
		svc := nuevoServicio(repo, empresaB, appEn("hired", empresaA))

		_, err := svc.Create(empresaA, CreatePlacementDTO{StaffingClientID: 1, ApplicationID: 5})
		if code := apperr.StatusCode(err); code != http.StatusForbidden {
			t.Fatalf("status %d, se esperaba 403", code)
		}
		if repo.creado != nil {
			t.Fatal("no debía crearse la colocación")
		}
	})

	t.Run("postulación de otra empresa", func(t *testing.T) {
		repo := &fakePlacementRepo{}
		svc := nuevoServicio(repo, empresaA, appEn("hired", empresaB))

		_, err := svc.Create(empresaA, CreatePlacementDTO{StaffingClientID: 1, ApplicationID: 5})
		if code := apperr.StatusCode(err); code != http.StatusForbidden {
			t.Fatalf("status %d, se esperaba 403", code)
		}
		if repo.creado != nil {
			t.Fatal("no debía crearse la colocación")
		}
	})
}

// Una postulación se coloca una sola vez: si no, se duplicaría la facturación
// al cliente final.
func TestCreateRechazaColocacionDuplicada(t *testing.T) {
	repo := &fakePlacementRepo{existe: true}
	svc := nuevoServicio(repo, empresaA, appEn("hired", empresaA))

	_, err := svc.Create(empresaA, CreatePlacementDTO{StaffingClientID: 1, ApplicationID: 5})
	if code := apperr.StatusCode(err); code != http.StatusConflict {
		t.Fatalf("status %d, se esperaba 409", code)
	}
	if repo.creado != nil {
		t.Fatal("no debía crearse la colocación")
	}
}

// El caso feliz: candidato y vacante se copian de la postulación, no del body,
// y la empresa se fuerza desde el token.
func TestCreateDerivaLosDatosDeLaPostulacion(t *testing.T) {
	repo := &fakePlacementRepo{}
	svc := nuevoServicio(repo, empresaA, appEn("hired", empresaA))

	if _, err := svc.Create(empresaA, CreatePlacementDTO{StaffingClientID: 1, ApplicationID: 5}); err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if repo.creado == nil {
		t.Fatal("no se creó la colocación")
	}
	if repo.creado.CompanyID != empresaA {
		t.Errorf("CompanyID = %d, se esperaba la empresa del token (%d)", repo.creado.CompanyID, empresaA)
	}
	if repo.creado.CandidateID != 77 || repo.creado.JobID != 88 {
		t.Errorf("candidato/vacante no se copiaron de la postulación: %d/%d", repo.creado.CandidateID, repo.creado.JobID)
	}
	if repo.creado.Status != "active" {
		t.Errorf("una colocación nueva nace activa, no %q", repo.creado.Status)
	}
}

func TestCreateExigeClienteYPostulacionExistentes(t *testing.T) {
	repo := &fakePlacementRepo{}
	svc := nuevoServicio(repo, empresaA, appEn("hired", empresaA))

	if _, err := svc.Create(empresaA, CreatePlacementDTO{StaffingClientID: 99, ApplicationID: 5}); apperr.StatusCode(err) != http.StatusNotFound {
		t.Fatal("cliente final inexistente debía dar 404")
	}
	if _, err := svc.Create(empresaA, CreatePlacementDTO{StaffingClientID: 1, ApplicationID: 99}); apperr.StatusCode(err) != http.StatusNotFound {
		t.Fatal("postulación inexistente debía dar 404")
	}
}
