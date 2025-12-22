package database

import (
	"dvra-api/internal/app/models"
)

var AllModels = []interface{}{
	&models.User{},
	&models.Role{},
	&models.Company{},
	&models.Job{},
	&models.Membership{},
	&models.Candidate{},
	&models.Application{},
	&models.Plan{},
	&models.SystemValue{},
}
