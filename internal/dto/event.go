package dto

import "time"

// CreateEventRequest for creating an event
type CreateEventRequest struct {
	Type string     `json:"type" binding:"required,min=2,max=32" example:"feed"`
	Note string     `json:"note" example:"dinner"`
	At   *time.Time `json:"at" example:"2025-11-22T18:00:00Z"` // if not specified, use now()
}
