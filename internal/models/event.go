package models

import "time"

// Event represents a pet event
// Examples of type: walk, feed, meds, training, vet, note
type Event struct {
	ID        uint      `json:"id" gorm:"primaryKey" example:"1"`
	Type      string    `json:"type" gorm:"index;size:32;not null" example:"walk"`
	Note      string    `json:"note" gorm:"size:500" example:"morning walk"`
	At        time.Time `json:"at" gorm:"index;not null" example:"2025-11-22T10:00:00Z"`
	CreatedAt time.Time `json:"created_at" example:"2025-11-22T10:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2025-11-22T10:00:00Z"`
}
