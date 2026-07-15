package models

import "time"

// SchemaMigration registra migraciones de datos ejecutadas una sola vez.
type SchemaMigration struct {
	ID        string    `gorm:"primaryKey;type:varchar(255)" json:"id"`
	AppliedAt time.Time `gorm:"not null" json:"applied_at"`
}

func (SchemaMigration) TableName() string {
	return "schema_migrations"
}
