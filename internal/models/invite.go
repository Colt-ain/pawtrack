package models

import "time"

// InviteStatus represents the status of an invitation
type InviteStatus string

const (
	InvitePending  InviteStatus = "pending"
	InviteAccepted InviteStatus = "accepted"
	InviteRejected InviteStatus = "rejected"
)

// Invite represents an invitation for a consultant to manage a dog
type Invite struct {
	ID           uint         `json:"id" gorm:"primaryKey"`
	OwnerID      uint         `json:"owner_id" gorm:"not null"`
	Owner        *User        `json:"owner,omitempty" gorm:"foreignKey:OwnerID"`
	ConsultantID uint         `json:"consultant_id" gorm:"not null"`
	Consultant   *User        `json:"consultant,omitempty" gorm:"foreignKey:ConsultantID"`
	DogID        uint         `json:"dog_id" gorm:"not null"`
	Dog          *Dog         `json:"dog,omitempty" gorm:"foreignKey:DogID"`
	Token        string       `json:"-" gorm:"uniqueIndex;not null;size:255"`
	Status       InviteStatus `json:"status" gorm:"default:'pending';size:20"`
	CreatedAt    time.Time    `json:"created_at"`
	ExpiresAt    time.Time    `json:"expires_at"`
}
