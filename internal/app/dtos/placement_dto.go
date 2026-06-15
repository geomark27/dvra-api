package dtos

import (
	"time"

	"dvra-api/internal/app/models"
)

// PlacementFilters representa los filtros para listar colocaciones
type PlacementFilters struct {
	Status           string `form:"status"`
	StaffingClientID *uint  `form:"staffing_client_id"`
	CandidateID      *uint  `form:"candidate_id"`
}

// PlacementResponseDTO representa la colocación en las respuestas de la API
type PlacementResponseDTO struct {
	ID               uint                   `json:"id"`
	CreatedAt        time.Time              `json:"created_at"`
	UpdatedAt        time.Time              `json:"updated_at"`
	CompanyID        uint                   `json:"company_id"`
	StaffingClientID uint                   `json:"staffing_client_id"`
	CandidateID      uint                   `json:"candidate_id"`
	JobID            *uint                  `json:"job_id,omitempty"`
	ApplicationID    *uint                  `json:"application_id,omitempty"`
	StartDate        *time.Time             `json:"start_date,omitempty"`
	EndDate          *time.Time             `json:"end_date,omitempty"`
	ContractType     string                 `json:"contract_type,omitempty"`
	Position         string                 `json:"position,omitempty"`
	BillRateAmount   *float64               `json:"bill_rate_amount,omitempty"`
	BillRateCurrency string                 `json:"bill_rate_currency,omitempty"`
	BillRateType     string                 `json:"bill_rate_type,omitempty"`
	PayRateAmount    *float64               `json:"pay_rate_amount,omitempty"`
	Status           string                 `json:"status"`
	Notes            string                 `json:"notes,omitempty"`
	StaffingClient   *models.StaffingClient `json:"staffing_client,omitempty"`
	Candidate        *models.Candidate      `json:"candidate,omitempty"`
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
	return PlacementResponseDTO{
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
		StaffingClient:   p.StaffingClient,
		Candidate:        p.Candidate,
	}
}

// ToPlacementResponseList convierte un slice de Placement a DTOs
func ToPlacementResponseList(placements []models.Placement) []PlacementResponseDTO {
	result := make([]PlacementResponseDTO, len(placements))
	for i := range placements {
		result[i] = ToPlacementResponse(&placements[i])
	}
	return result
}
