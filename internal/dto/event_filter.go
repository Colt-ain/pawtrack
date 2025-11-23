package dto

import (
	"time"

	"github.com/you/pawtrack/internal/models"
)

// EventFilterParams contains all possible filters for listing events
type EventFilterParams struct {
	// Date range
	FromDate *time.Time `form:"from_date" time_format:"2006-01-02"`
	ToDate   *time.Time `form:"to_date" time_format:"2006-01-02"`

	// Type filter (comma-separated)
	Types string `form:"types"`

	// Text search (searches in notes and dog name)
	Search string `form:"search"`

	// Dog name exact match
	DogName string `form:"dog_name"`

	// Pagination
	Page     int `form:"page,default=1" binding:"min=1"`
	PageSize int `form:"page_size,default=20" binding:"min=1,max=100"`

	// Sorting
	SortBy    string `form:"sort_by" binding:"omitempty,oneof=created_at type"`
	SortOrder string `form:"sort_order" binding:"omitempty,oneof=asc desc"`

	// Context (filled by handler)
	UserID   uint            `json:"-" form:"-"`
	UserRole models.UserRole `json:"-" form:"-"`
}

// EventListResponse represents paginated event list
type EventListResponse struct {
	Events     []models.Event `json:"events"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
	TotalCount int64          `json:"total_count"`
	TotalPages int            `json:"total_pages"`
}
