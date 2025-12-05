package seeders

import (
	"gorm.io/gorm"
)

type Seeder interface {
	Run(db *gorm.DB) error
}

var AllSeeders = []Seeder{
	&RoleSeeder{},       // 1. Primero roles del sistema
	&UserSeeder{},       // 2. Luego usuarios (SuperAdmin y Admin de empresa)
	&SuperAdminSeeder{}, // 3. Asignar rol SuperAdmin al usuario modo dios (ANTES de companies)
	&CompanySeeder{},    // 4. Luego empresas y membership del admin
}
