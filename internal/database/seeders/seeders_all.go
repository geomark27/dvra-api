package seeders

import (
	"gorm.io/gorm"
)

type Seeder interface {
	Run(db *gorm.DB) error
}

var AllSeeders = []Seeder{
	&RoleSeeder{},        // 1. Primero roles del sistema
	&PlanSeeder{},        // 2. Planes de suscripción
	&SystemValueSeeder{}, // 3. Valores del sistema (catálogos)
	&UserSeeder{},        // 4. Luego usuarios (SuperAdmin y Admin de empresa)
	&SuperAdminSeeder{},  // 5. Asignar rol SuperAdmin al usuario modo dios (ANTES de companies)
	&CompanySeeder{},     // 6. Luego empresas y membership del admin
}
