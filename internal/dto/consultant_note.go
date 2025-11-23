package dto

import "time"

// CreateNoteRequest for creating a new consultant note
type CreateNoteRequest struct {
	DogID   uint   `json:"dog_id" binding:"required"`
	Title   string `json:"title" binding:"required,max=255"`
	Content string `json:"content" binding:"required"`
}

// UpdateNoteRequest for updating a consultant note
type UpdateNoteRequest struct {
	Title   string `json:"title" binding:"max=255"`
	Content string `json:"content"`
}

// NoteFilterParams for filtering and sorting notes
type NoteFilterParams struct {
	Search   string `form:"search"`
	DogID    uint   `form:"dog_id"`
	OwnerID  uint   `form:"owner_id"`
	FromDate string `form:"from_date"` // RFC3339 format
	ToDate   string `form:"to_date"`   // RFC3339 format
	SortBy   string `form:"sort_by,default=created_at"` // created_at, updated_at, dog_name, owner_name
	Order    string `form:"order,default=desc"`         // asc, desc
	Page     int    `form:"page,default=1" binding:"min=1"`
	PageSize int    `form:"page_size,default=20" binding:"min=1,max=100"`
}

// NoteResponse for returning note details
type NoteResponse struct {
	ID           uint      `json:"id"`
	ConsultantID uint      `json:"consultant_id"`
	DogID        uint      `json:"dog_id"`
	DogName      string    `json:"dog_name"`
	OwnerID      uint      `json:"owner_id"`
	OwnerName    string    `json:"owner_name"`
	Title        string    `json:"title"`
	Content      string    `json:"content"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// NoteListResponse for paginated note list
type NoteListResponse struct {
	Notes      []NoteResponse `json:"notes"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
	TotalCount int64          `json:"total_count"`
	TotalPages int            `json:"total_pages"`
}
