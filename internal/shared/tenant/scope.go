// Package tenant define la frontera de datos de una operación autenticada.
//
// No importa gorm a propósito: así puede usarse desde la lógica de negocio de
// cualquier módulo sin arrastrar la persistencia. El filtro sobre una consulta
// vive en el subpaquete tenantgorm.
package tenant

import "dvra-api/internal/shared/apperr"

// Scope representa la frontera de datos de una operación autenticada.
type Scope struct {
	CompanyID uint
	Global    bool
	UserID    uint
	Role      string
}

func ForCompany(companyID uint) Scope {
	return Scope{CompanyID: companyID}
}

func ForActor(companyID, userID uint, role string) Scope {
	return Scope{CompanyID: companyID, UserID: userID, Role: role}
}

func GlobalScope() Scope {
	return Scope{Global: true}
}

func (s Scope) Valid() bool {
	return s.Global || s.CompanyID > 0
}

func (s Scope) Allows(companyID uint) bool {
	return s.Global || (s.CompanyID > 0 && s.CompanyID == companyID)
}

// Guard es la comprobación de pertenencia del §4.3 del estándar: devuelve
// Forbidden si el scope no alcanza a companyID.
//
// Se usa en el service, una sola vez por recurso, justo tras leerlo:
//
//	if err := scope.Guard(job.CompanyID); err != nil {
//		return nil, err
//	}
func (s Scope) Guard(companyID uint) error {
	if !s.Allows(companyID) {
		return apperr.Forbidden("Access denied")
	}
	return nil
}
