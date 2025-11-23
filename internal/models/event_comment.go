package models

import "time"

// EventComment represents a comment on an event
type EventComment struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	EventID   uint      `json:"event_id" gorm:"not null;index"`
	Event     *Event    `json:"event,omitempty" gorm:"foreignKey:EventID"`
	UserID    uint      `json:"user_id" gorm:"not null;index"`
	User      *User     `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Content   string    `json:"content" gorm:"type:text;not null"` // Markdown content
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
