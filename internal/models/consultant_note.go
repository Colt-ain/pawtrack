package models

import "time"

// ConsultantNote represents a note created by a consultant about a dog
type ConsultantNote struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	ConsultantID uint      `json:"consultant_id" gorm:"not null;index"`
	Consultant   *User     `json:"consultant,omitempty" gorm:"foreignKey:ConsultantID"`
	DogID        uint      `json:"dog_id" gorm:"not null;index"`
	Dog          *Dog      `json:"dog,omitempty" gorm:"foreignKey:DogID"`
	Title        string    `json:"title" gorm:"size:255;not null"`
	Content      string    `json:"content" gorm:"type:text;not null"` // Markdown content
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
