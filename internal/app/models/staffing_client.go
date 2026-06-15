package models

// StaffingClient representa al cliente final (empresa usuaria) gestionado por un
// tenant en modo staffing/outsourcing. NO es un tenant: vive dentro del tenant
// dueño (CompanyID) y es solo una dimensión de agrupación/billing. El aislamiento
// entre tenants sigue siendo company_id; staffing_client_id es un filtro interno.
type StaffingClient struct {
	BaseModel

	CompanyID uint `gorm:"not null;index:idx_staffing_clients_company_slug,priority:1" json:"company_id"` // tenant dueño (ej. firma de staffing)

	Name         string `gorm:"type:varchar(255);not null" json:"name"`
	Slug         string `gorm:"type:varchar(100);not null;index:idx_staffing_clients_company_slug,priority:2" json:"slug"` // único por company_id (validado en service)
	Industry     string `gorm:"type:varchar(100)" json:"industry,omitempty"`
	Website      string `gorm:"type:text" json:"website,omitempty"`
	LogoURL      string `gorm:"type:text" json:"logo_url,omitempty"`
	ContactName  string `gorm:"type:varchar(255)" json:"contact_name,omitempty"`
	ContactEmail string `gorm:"type:varchar(255)" json:"contact_email,omitempty"`
	ContactPhone string `gorm:"type:varchar(50)" json:"contact_phone,omitempty"`
	Status       string `gorm:"type:varchar(50);default:'active'" json:"status"` // active, inactive, prospect
	Notes        string `gorm:"type:text" json:"notes,omitempty"`

	// Relaciones
	Company    *Company    `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Jobs       []Job       `gorm:"foreignKey:StaffingClientID" json:"jobs,omitempty"`
	Placements []Placement `gorm:"foreignKey:StaffingClientID" json:"placements,omitempty"`
}

// TableName overrides the table name (optional)
func (StaffingClient) TableName() string {
	return "staffing_clients"
}
