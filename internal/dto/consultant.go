package dto

import "time"

// UpdateProfileRequest for updating consultant profile
type UpdateProfileRequest struct {
	Description string `json:"description" binding:"max=1000"`
	Services    string `json:"services" binding:"max=1000"`
	Breeds      string `json:"breeds" binding:"max=1000"`
	Location    string `json:"location" binding:"max=255"`
	Surname     string `json:"surname" binding:"max=255"`
}

// ConsultantSearchRequest for searching consultants
type ConsultantSearchRequest struct {
	Query    string `form:"query"`    // Search in name, surname, description
	Services string `form:"services"` // Filter by service
	Breeds   string `form:"breeds"`   // Filter by breed
	Location string `form:"location"` // Filter by location
	Page     int    `form:"page,default=1" binding:"min=1"`
	PageSize int    `form:"page_size,default=20" binding:"min=1,max=100"`
}

// ConsultantProfileResponse for returning consultant details
type ConsultantProfileResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Surname     string `json:"surname"`
	Description string `json:"description"`
	Services    string `json:"services"`
	Breeds      string `json:"breeds"`
	Location    string `json:"location"`
}

// CreateInviteRequest for inviting a consultant
type CreateInviteRequest struct {
	DogID uint `json:"dog_id" binding:"required"`
}

// InviteResponse for returning invite details
type InviteResponse struct {
	ID           uint      `json:"id"`
	Token        string    `json:"token,omitempty"` // Only returned to owner
	Status       string    `json:"status"`
	ConsultantID uint      `json:"consultant_id"`
	DogID        uint      `json:"dog_id"`
	ExpiresAt    time.Time `json:"expires_at"`
}
