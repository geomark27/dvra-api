package tenant

import "gorm.io/gorm"

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

// Apply añade el filtro tenant a una consulta. El acceso global debe estar
// representado explícitamente por GlobalScope.
func Apply(db *gorm.DB, scope Scope, column string) *gorm.DB {
	if scope.Global {
		return db
	}
	return db.Where(column+" = ?", scope.CompanyID)
}
