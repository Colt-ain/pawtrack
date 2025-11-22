package dto

import "github.com/you/pawtrack/internal/models"

// CreateUserRequest for user registration
type CreateUserRequest struct {
	Name     string            `json:"name" binding:"required,min=2,max=255" example:"John Doe"`
	Email    string            `json:"email" binding:"required,email,max=255" example:"john@example.com"`
	Password string            `json:"password" binding:"required,min=6,max=100" example:"securepassword123"`
	Role     models.UserRole   `json:"role" binding:"omitempty,oneof=owner consultant" example:"owner"`
}

// UpdateUserRequest for updating a user
type UpdateUserRequest struct {
	Name     string `json:"name" binding:"omitempty,min=2,max=255" example:"John Doe"`
	Email    string `json:"email" binding:"omitempty,email,max=255" example:"john@example.com"`
	Password string `json:"password" binding:"omitempty,min=6,max=100" example:"newpassword123"`
}
