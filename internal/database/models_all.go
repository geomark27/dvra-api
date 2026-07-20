package database

import (
	"fmt"
	"sort"

	"dvra-api/internal/app/models"
)

const (
	ModuleCore        = "core"
	ModuleIdentity    = "identity"
	ModuleRecruitment = "recruitment"
	ModuleStaffing    = "staffing"
	ModuleBilling     = "billing"
)

var ModuleModels = map[string][]interface{}{
	ModuleCore: {
		&models.SystemValueModule{}, &models.SystemValueCategory{}, &models.SystemValue{}, &models.CompanySystemValueOverride{},
		&models.Region{}, &models.Subregion{}, &models.Country{}, &models.State{}, &models.City{}, &models.StoredFile{}, &models.PlatformSettings{},
	},
	ModuleIdentity: {
		&models.TenantRole{}, &models.Permission{}, &models.TenantRolePermission{}, &models.User{}, &models.PlatformMembership{}, &models.Company{}, &models.Membership{},
	},
	ModuleRecruitment: {
		&models.Department{}, &models.Job{}, &models.Candidate{}, &models.Application{},
	},
	ModuleStaffing: {&models.StaffingClient{}, &models.Placement{}},
	ModuleBilling:  {&models.Plan{}, &models.CompanySubscription{}, &models.CompanySettings{}},
}

// AllModels keeps dependency order for the normal full migration.
var AllModels = modelsFor(ModuleCore, ModuleIdentity, ModuleBilling, ModuleRecruitment, ModuleStaffing)

func modelsFor(modules ...string) []interface{} {
	var result []interface{}
	for _, module := range modules {
		result = append(result, ModuleModels[module]...)
	}
	return result
}

func ModelsForModule(module string) ([]interface{}, error) {
	if _, ok := ModuleModels[module]; !ok {
		return nil, fmt.Errorf("unknown module %q (valid: %v)", module, Modules())
	}
	if module == ModuleCore {
		return modelsFor(ModuleCore), nil
	}
	// Every bounded context depends on the catalog infrastructure, but not on
	// the schema of unrelated business modules.
	return modelsFor(ModuleCore, module), nil
}

func Modules() []string {
	result := make([]string, 0, len(ModuleModels))
	for module := range ModuleModels {
		result = append(result, module)
	}
	sort.Strings(result)
	return result
}
