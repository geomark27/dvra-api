package server

import (
	"dvra-api/internal/app/repositories"
	staffingdomain "dvra-api/internal/modules/staffing/domain"
)

// staffingAppFinder adapta el repositorio de applications (módulo recruitment) al
// puerto staffingdomain.ApplicationFinder que el módulo staffing necesita. Vive en
// el composition root: ni staffing ni recruitment se importan mutuamente.
type staffingAppFinder struct {
	repo repositories.ApplicationRepository
}

func (a staffingAppFinder) FindByID(id uint) (*staffingdomain.HiredApplication, error) {
	app, err := a.repo.GetByID(id)
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
