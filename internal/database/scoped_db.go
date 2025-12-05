package database

import "gorm.io/gorm"

type ScopedDB struct {
    db        *gorm.DB
    companyID *uint
    userID    uint
}

func NewScopedDB(db *gorm.DB, companyID *uint, userID uint) *ScopedDB {
	return &ScopedDB{
		db:        db,
		companyID: companyID,
		userID:    userID,
	}
}

func (s *ScopedDB) WithCompanyScope() *gorm.DB  {
	if s.companyID == nil {
		return s.db
	}
	
	return s.db.Where("company_id = ?", &s.companyID)
}

