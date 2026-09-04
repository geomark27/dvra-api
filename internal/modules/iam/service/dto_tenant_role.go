package service

import "dvra-api/internal/shared/models"

type TenantRoleResponse struct {
	ID          uint   `json:"id"`
	Key         string `json:"key"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	IsSystem    bool   `json:"is_system"`
}

func ToTenantRoleResponse(role *models.TenantRole) *TenantRoleResponse {
	if role == nil {
		return nil
	}
	return &TenantRoleResponse{
		ID: role.ID, Key: role.Key, Name: role.Name,
		Description: role.Description, IsSystem: role.IsSystem,
	}
}
