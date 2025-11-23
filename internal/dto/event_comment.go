package dto

import "time"

// CreateCommentRequest for creating a new event comment
type CreateCommentRequest struct {
	EventID uint   `json:"event_id" binding:"required"`
	Content string `json:"content" binding:"required"`
	AttachmentURL *string `json:"-"` // Set by handler after upload
}

// UpdateCommentRequest for updating an event comment
type UpdateCommentRequest struct {
	Content string `json:"content" binding:"required"`
}

// CommentResponse for returning comment details
type CommentResponse struct {
	ID        uint      `json:"id"`
	EventID   uint      `json:"event_id"`
	UserID    uint      `json:"user_id"`
	UserName  string    `json:"user_name"`
	UserRole  string    `json:"user_role"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CommentListResponse for returning list of comments
type CommentListResponse struct {
	Comments []CommentResponse `json:"comments"`
	Count    int               `json:"count"`
}
