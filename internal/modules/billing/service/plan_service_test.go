package service

import (
	"errors"
	"net/http"
	"testing"

	"dvra-api/internal/shared/apperr"
	"dvra-api/internal/shared/models"
)

type fakePlanRepo struct{ PlanRepository }

type fakeCompanies struct {
	companies map[uint]*models.Company
	err       error
}

func (f fakeCompanies) GetByID(id uint) (*models.Company, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.companies[id], nil
}

type fakeSubscriptions struct {
	plan *models.Plan
	err  error
}

func (f fakeSubscriptions) CurrentPlan(uint) (*models.Plan, error) { return f.plan, f.err }

const empresa = uint(7)

func svcCon(plan *models.Plan, subErr error) PlanService {
	return NewPlanService(
		fakePlanRepo{},
		fakeCompanies{companies: map[uint]*models.Company{empresa: {BaseModel: models.BaseModel{ID: empresa}}}},
		fakeSubscriptions{plan: plan, err: subErr},
		nil, // db: estas rutas no abren transacción
	)
}

// El gating por plan de todas las rutas premium sale de aquí.
func TestCompanyHasFeatureSigueLosFlagsDelPlan(t *testing.T) {
	plan := &models.Plan{CanUseStaffing: true, CanExportData: false}

	if ok, err := svcCon(plan, nil).CompanyHasFeature(empresa, "staffing"); err != nil || !ok {
		t.Fatalf("staffing debía estar habilitado (ok=%v err=%v)", ok, err)
	}
	if ok, err := svcCon(plan, nil).CompanyHasFeature(empresa, "export_data"); err != nil || ok {
		t.Fatalf("export_data NO debía estar habilitado (ok=%v err=%v)", ok, err)
	}
	// Un feature desconocido nunca se concede por defecto.
	if ok, err := svcCon(plan, nil).CompanyHasFeature(empresa, "inventado"); err != nil || ok {
		t.Fatalf("un feature desconocido no se concede (ok=%v err=%v)", ok, err)
	}
}

// Sin suscripción vigente no hay features. El repositorio solo considera
// vigentes los estados trialing, active y past_due: una suscripción cancelada o
// vencida devuelve nil y aquí se traduce en denegar.
func TestCompanyHasFeatureNiegaSinSuscripcionVigente(t *testing.T) {
	ok, err := svcCon(nil, nil).CompanyHasFeature(empresa, "staffing")
	if ok {
		t.Fatal("sin suscripción vigente no debe concederse ningún feature")
	}
	if code := apperr.StatusCode(err); code != http.StatusNotFound {
		t.Fatalf("status %d, se esperaba 404", code)
	}
}

func TestCompanyHasFeatureExigeEmpresaExistente(t *testing.T) {
	svc := NewPlanService(fakePlanRepo{}, fakeCompanies{companies: map[uint]*models.Company{}},
		fakeSubscriptions{}, nil)

	ok, err := svc.CompanyHasFeature(999, "staffing")
	if ok {
		t.Fatal("no debe concederse nada para una empresa inexistente")
	}
	if code := apperr.StatusCode(err); code != http.StatusNotFound {
		t.Fatalf("status %d, se esperaba 404", code)
	}
}

// Un fallo de infraestructura no puede degradarse en "no tienes el feature":
// eso convertiría una caída en una denegación silenciosa de servicio de pago.
func TestCompanyHasFeaturePropagaErroresDeInfraestructura(t *testing.T) {
	boom := errors.New("connection refused")

	if _, err := svcCon(nil, boom).CompanyHasFeature(empresa, "staffing"); !errors.Is(err, boom) {
		t.Fatalf("se esperaba el error de la suscripción, se obtuvo %v", err)
	}

	svc := NewPlanService(fakePlanRepo{}, fakeCompanies{err: boom}, fakeSubscriptions{}, nil)
	if _, err := svc.CompanyHasFeature(empresa, "staffing"); !errors.Is(err, boom) {
		t.Fatalf("se esperaba el error de empresas, se obtuvo %v", err)
	}
}
