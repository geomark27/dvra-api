package database

import (
	"dvra-api/internal/app/models"
)

var AllModels = []interface{}{
	&models.User{},
	&models.Company{},
	&models.Job{},
	&models.Membership{},
	&models.Candidate{},
	&models.Application{},
}
