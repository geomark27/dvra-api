// Package publicapi es el punto de ensamblaje de la career page pública: el
// portal donde un candidato ve las vacantes de una empresa y se postula sin
// tener cuenta.
//
// Es un agregador: no tiene datos ni repositorios propios. Todo lo compone a
// partir de puertos hacia iam y recruitment.
package publicapi

import (
	"dvra-api/internal/modules/publicapi/service"
	"dvra-api/internal/modules/publicapi/transport"
	"dvra-api/internal/platform/httpx"
)

// Deps son los puertos cross-módulo que este agregador necesita.
type Deps struct {
	Companies    service.CompanyReader
	Jobs         service.JobReader
	Candidates   service.CandidateWriter
	Applications service.ApplicationWriter
	Resumes      service.ResumeStore
}

// Module agrupa las dependencias ya cableadas del módulo.
type Module struct {
	svc service.PublicService
}

// New construye el módulo a partir de sus puertos.
func New(deps Deps) *Module {
	return &Module{
		svc: service.NewPublicService(deps.Companies, deps.Jobs, deps.Candidates, deps.Applications, deps.Resumes),
	}
}

// Name identifica al módulo en logs y en la auditoría de rutas.
func (m *Module) Name() string { return "publicapi" }

// RegisterRoutes monta las rutas públicas.
func (m *Module) RegisterRoutes(r httpx.Router) {
	transport.RegisterRoutes(r.Public, m.svc)
}
