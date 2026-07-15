package dtos

// CreateUserDTO representa los datos para crear un usuario
type CreateUserDTO struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	Role      string `json:"role" binding:"required,oneof=admin recruiter hiring_manager user"`
	CompanyID uint   `json:"-"`
	InvitedBy uint   `json:"-"`
}

// UpdateUserDTO representa los datos para actualizar un usuario
type UpdateUserDTO struct {
	FirstName *string `json:"first_name,omitempty" binding:"omitempty,min=1"`
	LastName  *string `json:"last_name,omitempty" binding:"omitempty,min=1"`
	Email     *string `json:"email,omitempty" binding:"omitempty,email"`
}

// UserResponseDTO representa la respuesta de un usuario (sin datos sensibles)
type UserResponseDTO struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}
