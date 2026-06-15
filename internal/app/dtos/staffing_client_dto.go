package dtos

import (
	"time"

	"dvra-api/internal/app/models"
)

// StaffingClientFilters representa los filtros para listar clientes finales
type StaffingClientFilters struct {
	Status string `form:"status"`
}

// StaffingClientResponseDTO representa el cliente final en las respuestas de la API
type StaffingClientResponseDTO struct {
	ID           uint      `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	CompanyID    uint      `json:"company_id"`
	Name         string    `json:"name"`
	Slug         string    `json:"slug"`
	Industry     string    `json:"industry,omitempty"`
	Website      string    `json:"website,omitempty"`
	LogoURL      string    `json:"logo_url,omitempty"`
	ContactName  string    `json:"contact_name,omitempty"`
	ContactEmail string    `json:"contact_email,omitempty"`
	ContactPhone string    `json:"contact_phone,omitempty"`
	Status       string    `json:"status"`
	Notes        string    `json:"notes,omitempty"`
}

// CreateStaffingClientDTO representa los datos para crear un cliente final.
// CompanyID se fuerza desde el token en el handler (no se confía en el body).
type CreateStaffingClientDTO struct {
	CompanyID    uint   `json:"company_id"`
	Name         string `json:"name" binding:"required,min=2,max=255"`
	Slug         string `json:"slug" binding:"required,min=2,max=100"`
	Industry     string `json:"industry,omitempty"`
	Website      string `json:"website,omitempty"`
	LogoURL      string `json:"logo_url,omitempty"`
	ContactName  string `json:"contact_name,omitempty"`
	ContactEmail string `json:"contact_email,omitempty" binding:"omitempty,email"`
	ContactPhone string `json:"contact_phone,omitempty"`
	Status       string `json:"status,omitempty" binding:"omitempty,oneof=active inactive prospect"`
	Notes        string `json:"notes,omitempty"`
}

// UpdateStaffingClientDTO representa la actualización parcial de un cliente final
type UpdateStaffingClientDTO struct {
	Name         *string `json:"name,omitempty" binding:"omitempty,min=2,max=255"`
	Slug         *string `json:"slug,omitempty" binding:"omitempty,min=2,max=100"`
	Industry     *string `json:"industry,omitempty"`
	Website      *string `json:"website,omitempty"`
	LogoURL      *string `json:"logo_url,omitempty"`
	ContactName  *string `json:"contact_name,omitempty"`
	ContactEmail *string `json:"contact_email,omitempty" binding:"omitempty,email"`
	ContactPhone *string `json:"contact_phone,omitempty"`
	Status       *string `json:"status,omitempty" binding:"omitempty,oneof=active inactive prospect"`
	Notes        *string `json:"notes,omitempty"`
}

// ToStaffingClientResponse convierte un modelo StaffingClient a su DTO de respuesta
func ToStaffingClientResponse(c *models.StaffingClient) StaffingClientResponseDTO {
	return StaffingClientResponseDTO{
		ID:           c.ID,
		CreatedAt:    c.CreatedAt,
		UpdatedAt:    c.UpdatedAt,
		CompanyID:    c.CompanyID,
		Name:         c.Name,
		Slug:         c.Slug,
		Industry:     c.Industry,
		Website:      c.Website,
		LogoURL:      c.LogoURL,
		ContactName:  c.ContactName,
		ContactEmail: c.ContactEmail,
		ContactPhone: c.ContactPhone,
		Status:       c.Status,
		Notes:        c.Notes,
	}
}

// ToStaffingClientResponseList convierte un slice de StaffingClient a DTOs
func ToStaffingClientResponseList(clients []models.StaffingClient) []StaffingClientResponseDTO {
	result := make([]StaffingClientResponseDTO, len(clients))
	for i := range clients {
		result[i] = ToStaffingClientResponse(&clients[i])
	}
	return result
}
