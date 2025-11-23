package dto

import "time"

// CreateEventRequest for creating a new event
type CreateEventRequest struct {
	DogID *uint     `json:"dog_id" example:"1"`
	Type  string    `json:"type" binding:"required,max=50" example:"walk"`
	Note  string    `json:"note" binding:"max=255" example:"morning walk"`
	At    *time.Time `json:"at" example:"2025-11-22T10:00:00Z"`
} // if not specified, use now()
