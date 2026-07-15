package dtos

import "time"

// CreateMembershipDTO represents the data needed to create a membership
type CreateMembershipDTO struct {
	UserID    uint   `json:"user_id" binding:"required,min=1"`
	CompanyID *uint  `json:"company_id"` // Pointer allows reporting a missing required field.
	Role      string `json:"role" binding:"required,oneof=admin recruiter hiring_manager user"`
	Status    string `json:"status" binding:"omitempty,oneof=active inactive pending suspended removed"`
	IsDefault bool   `json:"is_default"`
	InvitedBy *uint  `json:"invited_by,omitempty"`
}

// UpdateMembershipDTO represents the data needed to update a membership
type UpdateMembershipDTO struct {
	Role      *string    `json:"role,omitempty" binding:"omitempty,oneof=admin recruiter hiring_manager user"`
	Status    *string    `json:"status,omitempty" binding:"omitempty,oneof=active inactive pending suspended removed"`
	IsDefault *bool      `json:"is_default,omitempty"`
	JoinedAt  *time.Time `json:"joined_at,omitempty"`
}
