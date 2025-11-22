package models

import "time"

// Dog represents a dog
type Dog struct {
	ID        uint      `json:"id" gorm:"primaryKey" example:"1"`
	OwnerID   uint      `json:"owner_id" gorm:"not null;index" example:"1"`
	Owner     *User     `json:"owner,omitempty" gorm:"foreignKey:OwnerID"`
	Name      string    `json:"name" gorm:"size:255;not null" example:"Rex"`
	Breed     string    `json:"breed" gorm:"size:255" example:"Golden Retriever"`
	BirthDate time.Time `json:"birth_date" example:"2020-01-01T00:00:00Z"`
	CreatedAt time.Time `json:"created_at" example:"2025-11-22T10:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2025-11-22T10:00:00Z"`
}
