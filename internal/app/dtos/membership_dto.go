package dtos

import "time"

// CreateMembershipDTO represents the data needed to create a membership
type CreateMembershipDTO struct {
	UserID    uint   `json:"user_id" validate:"required,min=1"`
	CompanyID *uint  `json:"company_id,omitempty"` // Nullable para SuperAdmin
	Role      string `json:"role" validate:"required,oneof=superadmin admin recruiter user"`
	Status    string `json:"status" validate:"omitempty,oneof=active inactive pending"`
	IsDefault bool   `json:"is_default"`
	InvitedBy *uint  `json:"invited_by,omitempty"`
}

// UpdateMembershipDTO represents the data needed to update a membership
type UpdateMembershipDTO struct {
	Role      *string    `json:"role,omitempty" validate:"omitempty,oneof=superadmin admin recruiter user"`
	Status    *string    `json:"status,omitempty" validate:"omitempty,oneof=active inactive pending"`
	IsDefault *bool      `json:"is_default,omitempty"`
	JoinedAt  *time.Time `json:"joined_at,omitempty"`
}
