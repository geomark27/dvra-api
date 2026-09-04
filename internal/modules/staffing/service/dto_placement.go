package service

import (
	"time"

	"dvra-api/internal/shared/models"
)

// PlacementResponseDTO representa la colocación en las respuestas de la API
type PlacementResponseDTO struct {
	ID               uint                       `json:"id"`
	CreatedAt        time.Time                  `json:"created_at"`
	UpdatedAt        time.Time                  `json:"updated_at"`
	CompanyID        uint                       `json:"company_id"`
	StaffingClientID uint                       `json:"staffing_client_id"`
	CandidateID      uint                       `json:"candidate_id"`
	JobID            uint                       `json:"job_id"`
	ApplicationID    uint                       `json:"application_id"`
	StartDate        *time.Time                 `json:"start_date,omitempty"`
	EndDate          *time.Time                 `json:"end_date,omitempty"`
	ContractType     string                     `json:"contract_type,omitempty"`
	Position         string                     `json:"position,omitempty"`
	BillRateAmount   *float64                   `json:"bill_rate_amount,omitempty"`
	BillRateCurrency string                     `json:"bill_rate_currency,omitempty"`
	BillRateType     string                     `json:"bill_rate_type,omitempty"`
	PayRateAmount    *float64                   `json:"pay_rate_amount,omitempty"`
	Status           string                     `json:"status"`
	Notes            string                     `json:"notes,omitempty"`
	StaffingClient   *StaffingClientResponseDTO `json:"staffing_client,omitempty"`
	Candidate        *PlacementCandidateDTO     `json:"candidate,omitempty"`
}

// CreatePlacementDTO representa los datos para crear una colocación.
// La colocación deriva de una Application: CandidateID y JobID se copian de ella,
// no se aceptan en el body. CompanyID se fuerza desde el token en el handler.
type CreatePlacementDTO struct {
	StaffingClientID uint       `json:"staffing_client_id" binding:"required,min=1"`
	ApplicationID    uint       `json:"application_id" binding:"required,min=1"`
	StartDate        *time.Time `json:"start_date,omitempty"`
	EndDate          *time.Time `json:"end_date,omitempty"`
	ContractType     string     `json:"contract_type,omitempty" binding:"omitempty,oneof=outsourcing staffing project"`
	Position         string     `json:"position,omitempty"`
	BillRateAmount   *float64   `json:"bill_rate_amount,omitempty"`
	BillRateCurrency string     `json:"bill_rate_currency,omitempty" binding:"omitempty,len=3"`
	BillRateType     string     `json:"bill_rate_type,omitempty" binding:"omitempty,oneof=hourly monthly"`
	PayRateAmount    *float64   `json:"pay_rate_amount,omitempty"`
	Status           string     `json:"status,omitempty" binding:"omitempty,oneof=active ended suspended"`
	Notes            string     `json:"notes,omitempty"`
}

// UpdatePlacementDTO representa la actualización parcial de una colocación.
// El origen (Application/Candidate/StaffingClient) es inmutable; solo cambian
// los datos de contrato, billing y estado.
type UpdatePlacementDTO struct {
	StartDate        *time.Time `json:"start_date,omitempty"`
	EndDate          *time.Time `json:"end_date,omitempty"`
	ContractType     *string    `json:"contract_type,omitempty" binding:"omitempty,oneof=outsourcing staffing project"`
	Position         *string    `json:"position,omitempty"`
	BillRateAmount   *float64   `json:"bill_rate_amount,omitempty"`
	BillRateCurrency *string    `json:"bill_rate_currency,omitempty" binding:"omitempty,len=3"`
	BillRateType     *string    `json:"bill_rate_type,omitempty" binding:"omitempty,oneof=hourly monthly"`
	PayRateAmount    *float64   `json:"pay_rate_amount,omitempty"`
	Status           *string    `json:"status,omitempty" binding:"omitempty,oneof=active ended suspended"`
	Notes            *string    `json:"notes,omitempty"`
}

// ToPlacementResponse convierte un modelo Placement a su DTO de respuesta
func ToPlacementResponse(p *models.Placement) PlacementResponseDTO {
	dto := PlacementResponseDTO{
		ID:               p.ID,
		CreatedAt:        p.CreatedAt,
		UpdatedAt:        p.UpdatedAt,
		CompanyID:        p.CompanyID,
		StaffingClientID: p.StaffingClientID,
		CandidateID:      p.CandidateID,
		JobID:            p.JobID,
		ApplicationID:    p.ApplicationID,
		StartDate:        p.StartDate,
		EndDate:          p.EndDate,
		ContractType:     p.ContractType,
		Position:         p.Position,
		BillRateAmount:   p.BillRateAmount,
		BillRateCurrency: p.BillRateCurrency,
		BillRateType:     p.BillRateType,
		PayRateAmount:    p.PayRateAmount,
		Status:           p.Status,
		Notes:            p.Notes,
	}
	if p.StaffingClient != nil {
		r := ToStaffingClientResponse(p.StaffingClient)
		dto.StaffingClient = &r
	}
	if p.Candidate != nil {
		r := toPlacementCandidate(p.Candidate)
		dto.Candidate = &r
	}
	return dto
}

// ToPlacementResponseList convierte un slice de Placement a DTOs
func ToPlacementResponseList(placements []models.Placement) []PlacementResponseDTO {
	result := make([]PlacementResponseDTO, len(placements))
	for i := range placements {
		result[i] = ToPlacementResponse(&placements[i])
	}
	return result
}

// PlacementCandidateDTO es la vista del candidato que staffing incluye en una
// colocación.
//
// Es propia del módulo y no el DTO de recruitment: un módulo no expone los
// contratos de otro (ESTANDARES §3.3 regla 2). Los campos y sus etiquetas JSON
// se mantienen idénticos para no alterar la respuesta.
type PlacementCandidateDTO struct {
	ID          uint      `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	CompanyID   uint      `json:"company_id"`
	Email       string    `json:"email"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Phone       string    `json:"phone,omitempty"`
	ResumeURL   string    `json:"resume_url,omitempty"`
	GithubURL   string    `json:"github_url,omitempty"`
	LinkedinURL string    `json:"linkedin_url,omitempty"`
	Source      string    `json:"source,omitempty"`
}

func toPlacementCandidate(c *models.Candidate) PlacementCandidateDTO {
	return PlacementCandidateDTO{
		ID:          c.ID,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
		CompanyID:   c.CompanyID,
		Email:       c.Email,
		FirstName:   c.FirstName,
		LastName:    c.LastName,
		Phone:       c.Phone,
		ResumeURL:   c.ResumeURL,
		GithubURL:   c.GithubURL,
		LinkedinURL: c.LinkedinURL,
		Source:      c.Source,
	}
}
