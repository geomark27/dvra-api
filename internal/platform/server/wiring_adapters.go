package server

import (
	"context"
	"errors"

	publicapiservice "dvra-api/internal/modules/publicapi/service"
	recruitmentrepo "dvra-api/internal/modules/recruitment/repository"
	recruitmentservice "dvra-api/internal/modules/recruitment/service"
	staffingdomain "dvra-api/internal/modules/staffing/domain"
	"dvra-api/internal/shared/models"
	"dvra-api/internal/shared/tenant"

	"gorm.io/gorm"
)

// staffingAppFinder adapta el repositorio de applications (módulo recruitment) al
// puerto staffingdomain.ApplicationFinder que el módulo staffing necesita. Vive en
// el composition root: ni staffing ni recruitment se importan mutuamente.
type staffingAppFinder struct {
	repo *recruitmentrepo.ApplicationRepository
}

func (a staffingAppFinder) FindByID(id uint) (*staffingdomain.HiredApplication, error) {
	app, err := a.repo.GetByID(id, tenant.GlobalScope())
	if err != nil || app == nil {
		return nil, err
	}
	return &staffingdomain.HiredApplication{
		ID:          app.ID,
		CompanyID:   app.CompanyID,
		CandidateID: app.CandidateID,
		JobID:       app.JobID,
		Stage:       app.Stage,
	}, nil
}

// gormPinger adapta *gorm.DB al puerto health.Pinger. Vive en el composition
// root para que la sonda no conozca la persistencia.
type gormPinger struct{ db *gorm.DB }

func (p gormPinger) Ping(ctx context.Context) error {
	if p.db == nil {
		return errors.New("database is not configured")
	}
	sqlDB, err := p.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.PingContext(ctx)
}

// resumeStore adapta el service de candidatos de recruitment al puerto
// ResumeStore de publicapi.
//
// Hace falta un adaptador y no basta la satisfacción estructural porque el
// parámetro es un struct con nombre propio en cada módulo: Go no unifica
// recruitment.ResumeUpload con publicapi.ResumeUpload aunque tengan los mismos
// campos. Que cada módulo tenga el suyo es justamente lo que evita que uno
// dependa de los tipos del otro.
type resumeStore struct {
	svc recruitmentservice.CandidateService
}

func (r resumeStore) StoreResume(id uint, scope tenant.Scope, upload publicapiservice.ResumeUpload) (*models.StoredFile, error) {
	return r.svc.StoreResume(id, scope, recruitmentservice.ResumeUpload{
		Reader:       upload.Reader,
		OriginalName: upload.OriginalName,
		Size:         upload.Size,
	})
}
