package seeders

import (
	"gorm.io/gorm"
)

type Seeder interface {
	Run(db *gorm.DB) error
}

var AllSeeders = []Seeder{
	&PlatformSettingsSeeder{}, // 0. Configuración global de la plataforma (primero)
	&RoleSeeder{},             // 1. Primero roles del sistema
	&PlanSeeder{},             // 2. Planes de suscripción
	&SystemValueSeeder{},      // 3. Valores del sistema (catálogos)
	&UserSeeder{},             // 4. Usuarios (Admin de empresa)
	&CompanySeeder{},          // 5. Empresas y membership del admin
}
