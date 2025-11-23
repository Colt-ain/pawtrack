package models

import "time"

// Event represents a tracking event
// Examples of type: walk, feed, meds, training, vet, note
type Event struct {
	ID        uint      `json:"id" gorm:"primaryKey" example:"1"`
	DogID     *uint     `json:"dog_id,omitempty" gorm:"index" example:"1"`
	Dog       *Dog      `json:"dog,omitempty" gorm:"foreignKey:DogID"`
	Type      string    `json:"type" gorm:"size:50;not null" example:"walk"`
	Note      string    `json:"note" gorm:"size:255" example:"morning walk"`
	At        time.Time `json:"at" gorm:"not null" example:"2025-11-22T10:00:00Z"`
	CreatedAt time.Time `json:"created_at" example:"2025-11-22T10:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2025-11-22T10:00:00Z"`
}
