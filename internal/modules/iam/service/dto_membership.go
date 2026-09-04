package service

import "time"

// CreateMembershipDTO represents the data needed to create a membership
type CreateMembershipDTO struct {
	UserID       uint   `json:"user_id" binding:"required,min=1"`
	CompanyID    *uint  `json:"company_id"` // Pointer allows reporting a missing required field.
	TenantRoleID uint   `json:"tenant_role_id" binding:"required,min=1"`
	Status       string `json:"status" binding:"omitempty,oneof=active inactive pending suspended removed"`
	IsDefault    bool   `json:"is_default"`
	InvitedBy    *uint  `json:"invited_by,omitempty"`
}

// UpdateMembershipDTO represents the data needed to update a membership
type UpdateMembershipDTO struct {
	TenantRoleID *uint      `json:"tenant_role_id,omitempty" binding:"omitempty,min=1"`
	Status       *string    `json:"status,omitempty" binding:"omitempty,oneof=active inactive pending suspended removed"`
	IsDefault    *bool      `json:"is_default,omitempty"`
	JoinedAt     *time.Time `json:"joined_at,omitempty"`
}
