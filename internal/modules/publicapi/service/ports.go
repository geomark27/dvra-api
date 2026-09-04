package service

import (
	"io"

	"dvra-api/internal/shared/models"
	"dvra-api/internal/shared/tenant"
)

// Puertos del módulo publicapi.
//
// Es un agregador: no tiene datos propios. Sirve la career page pública
// componiendo lo que exponen iam (la empresa) y recruitment (vacantes,
// candidatos, aplicaciones). Por eso todos sus puertos son cross-módulo y se
// declaran aquí, en el consumidor (§3.3 regla 3).
//
// Al ser público y sin sesión, el scope no viene de un token: lo fija este
// service a partir de la empresa resuelta por el slug.

// CompanyReader resuelve la empresa dueña de la career page. Puerto hacia iam.
type CompanyReader interface {
	GetBySlug(slug string) (*models.Company, error)
}

// JobReader expone las vacantes publicadas. Puerto hacia recruitment.
type JobReader interface {
	GetByID(id uint, scope tenant.Scope) (*models.Job, error)
	GetPublishedByCompanyID(companyID uint) ([]models.Job, error)
}

// CandidateWriter da de alta o actualiza al candidato que se postula. Puerto
// hacia recruitment.
type CandidateWriter interface {
	GetByEmail(email string, companyID uint) (*models.Candidate, error)
	Create(candidate *models.Candidate) (*models.Candidate, error)
	Update(candidate *models.Candidate) (*models.Candidate, error)
}

// ApplicationWriter registra la postulación, evitando duplicados. Puerto hacia
// recruitment.
type ApplicationWriter interface {
	GetByCandidateAndJob(candidateID, jobID uint) (*models.Application, error)
	Create(application *models.Application) (*models.Application, error)
}

// ResumeUpload es el CV que llega en la postulación pública.
type ResumeUpload struct {
	Reader       io.Reader
	OriginalName string
	Size         int64
}

// ResumeStore guarda el CV aplicando la cuota del plan. Puerto hacia
// recruitment, cuyo service ya implementa esa regla; publicapi no la duplica.
type ResumeStore interface {
	StoreResume(id uint, scope tenant.Scope, upload ResumeUpload) (*models.StoredFile, error)
}
