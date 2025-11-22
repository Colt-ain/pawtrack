package models

import "time"

// UserRole represents user role in the system
type UserRole string

const (
	RoleOwner      UserRole = "owner"
	RoleConsultant UserRole = "consultant"
	RoleAdmin      UserRole = "admin"
)

// User represents a user
type User struct {
	ID           uint      `json:"id" gorm:"primaryKey" example:"1"`
	Name         string    `json:"name" gorm:"size:255;not null" example:"John Doe"`
	Email        string    `json:"email" gorm:"size:255;uniqueIndex;not null" example:"john@example.com"`
	PasswordHash string    `json:"-" gorm:"size:255;not null"`
	Role         UserRole  `json:"role" gorm:"type:varchar(20);not null;default:'owner'" example:"owner"`
	CreatedAt    time.Time `json:"created_at" example:"2025-11-22T10:00:00Z"`
	UpdatedAt    time.Time `json:"updated_at" example:"2025-11-22T10:00:00Z"`
}
