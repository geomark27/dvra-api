package seeders

import (
	"gorm.io/gorm"
)

type Seeder interface {
	Run(db *gorm.DB) error
}

var AllSeeders = []Seeder{
	&RoleSeeder{},       // 1. Primero roles del sistema
	&PlanSeeder{},       // 2. Planes de suscripción
	&UserSeeder{},       // 3. Luego usuarios (SuperAdmin y Admin de empresa)
	&SuperAdminSeeder{}, // 4. Asignar rol SuperAdmin al usuario modo dios (ANTES de companies)
	&CompanySeeder{},    // 5. Luego empresas y membership del admin
}
