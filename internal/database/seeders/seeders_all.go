package seeders

import "gorm.io/gorm"

type Seeder interface{ Run(db *gorm.DB) error }

type RegisteredSeeder struct {
	Name   string
	Module string
	Seeder Seeder
}

var AllSeeders = []RegisteredSeeder{
	{Name: "identity-catalogs", Module: "identity", Seeder: &SystemValueSeeder{Module: "identity"}},
	{Name: "recruitment-catalogs", Module: "recruitment", Seeder: &SystemValueSeeder{Module: "recruitment"}},
	{Name: "staffing-catalogs", Module: "staffing", Seeder: &SystemValueSeeder{Module: "staffing"}},
	{Name: "billing-catalogs", Module: "billing", Seeder: &SystemValueSeeder{Module: "billing"}},
	{Name: "plans", Module: "billing", Seeder: &PlanSeeder{}},
	{Name: "platform-settings", Module: "core", Seeder: &PlatformSettingsSeeder{}},
	{Name: "demo-users", Module: "demo", Seeder: &UserSeeder{}},
	{Name: "demo-companies", Module: "demo", Seeder: &CompanySeeder{}},
}
