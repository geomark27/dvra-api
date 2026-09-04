// Package tenantgorm aplica el scope de tenant sobre una consulta GORM.
//
// Vive separado de shared/tenant para que el núcleo (el tipo Scope y sus
// reglas) no arrastre gorm: la lógica de negocio de cualquier módulo puede
// razonar sobre un Scope sin conocer la persistencia.
package tenantgorm

import (
	"dvra-api/internal/shared/tenant"

	"gorm.io/gorm"
)

// Apply añade el filtro tenant a una consulta. El acceso global debe estar
// representado explícitamente por tenant.GlobalScope.
func Apply(db *gorm.DB, scope tenant.Scope, column string) *gorm.DB {
	if scope.Global {
		return db
	}
	return db.Where(column+" = ?", scope.CompanyID)
}
