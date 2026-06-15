package models

import "time"

// Placement representa la colocación de un candidato en un cliente final
// (StaffingClient) dentro del modo staffing. Nace de una Application que llegó a
// la etapa de contratación: ApplicationID es el origen y CandidateID/JobID se
// copian de ella para facilitar consultas. CompanyID es el tenant dueño y la
// frontera de aislamiento; StaffingClientID es solo un filtro interno.
type Placement struct {
	BaseModel

	CompanyID        uint  `gorm:"not null;index:idx_placements_company_status,priority:1" json:"company_id"` // tenant (firma staffing)
	StaffingClientID uint  `gorm:"not null;index" json:"staffing_client_id"`                                  // cliente final
	CandidateID      uint  `gorm:"not null;index" json:"candidate_id"`
	JobID            *uint `gorm:"index" json:"job_id,omitempty"`         // vacante origen (copiada de la Application)
	ApplicationID    *uint `gorm:"index" json:"application_id,omitempty"` // Application de origen

	// Contrato
	StartDate    *time.Time `json:"start_date,omitempty"`
	EndDate      *time.Time `json:"end_date,omitempty"`                              // nil = indefinido
	ContractType string     `gorm:"type:varchar(50)" json:"contract_type,omitempty"` // outsourcing, staffing, project
	Position     string     `gorm:"type:varchar(255)" json:"position,omitempty"`

	// Billing
	BillRateAmount   *float64 `gorm:"type:decimal(12,2)" json:"bill_rate_amount,omitempty"`
	BillRateCurrency string   `gorm:"type:varchar(3);default:'USD'" json:"bill_rate_currency,omitempty"`
	BillRateType     string   `gorm:"type:varchar(20)" json:"bill_rate_type,omitempty"` // hourly, monthly
	PayRateAmount    *float64 `gorm:"type:decimal(12,2)" json:"pay_rate_amount,omitempty"`

	Status string `gorm:"type:varchar(50);default:'active';index:idx_placements_company_status,priority:2" json:"status"` // active, ended, suspended
	Notes  string `gorm:"type:text" json:"notes,omitempty"`

	// Relaciones
	Company        *Company        `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	StaffingClient *StaffingClient `gorm:"foreignKey:StaffingClientID" json:"staffing_client,omitempty"`
	Candidate      *Candidate      `gorm:"foreignKey:CandidateID" json:"candidate,omitempty"`
	Job            *Job            `gorm:"foreignKey:JobID" json:"job,omitempty"`
	Application    *Application    `gorm:"foreignKey:ApplicationID" json:"application,omitempty"`
}

// TableName overrides the table name (optional)
func (Placement) TableName() string {
	return "placements"
}
