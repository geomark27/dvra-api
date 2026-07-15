package seeders

import (
	"gorm.io/gorm"
)

type Seeder interface {
	Run(db *gorm.DB) error
}

var AllSeeders = []Seeder{
	&PlatformSettingsSeeder{}, // 0. Configuración global de la plataforma (primero)
	&PlanSeeder{},             // 1. Planes de suscripción
	&SystemValueSeeder{},      // 2. Valores del sistema (catálogos)
	&UserSeeder{},             // 3. Usuarios y perfil landlord demo
	&CompanySeeder{},          // 4. Empresas y membership del admin
}
