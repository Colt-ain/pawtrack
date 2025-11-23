package models

import "time"

// ConsultantAccess represents access rights for a consultant to a specific dog
type ConsultantAccess struct {
	ID           uint       `json:"id" gorm:"primaryKey"`
	ConsultantID uint       `json:"consultant_id" gorm:"not null;index"`
	Consultant   *User      `json:"consultant,omitempty" gorm:"foreignKey:ConsultantID"`
	DogID        uint       `json:"dog_id" gorm:"not null;index"`
	Dog          *Dog       `json:"dog,omitempty" gorm:"foreignKey:DogID"`
	GrantedAt    time.Time  `json:"granted_at" gorm:"not null"`
	RevokedAt    *time.Time `json:"revoked_at,omitempty"`
}

// TableName specifies the table name for Consultant Access
func (ConsultantAccess) TableName() string {
	return "consultant_access"
}
