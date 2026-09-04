package models

// StoredFile registra metadatos de archivos privados almacenados fuera de BD.
type StoredFile struct {
	BaseModel
	CompanyID    uint   `gorm:"not null;index:idx_stored_files_company" json:"company_id"`
	CandidateID  uint   `gorm:"not null;uniqueIndex" json:"candidate_id"`
	Key          string `gorm:"type:text;not null;uniqueIndex" json:"-"`
	OriginalName string `gorm:"type:varchar(255);not null" json:"original_name"`
	ContentType  string `gorm:"type:varchar(150);not null" json:"content_type"`
	SizeBytes    int64  `gorm:"not null" json:"size_bytes"`
}

func (StoredFile) TableName() string { return "stored_files" }
