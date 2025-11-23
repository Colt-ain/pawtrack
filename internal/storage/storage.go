package storage

import (
	"io"
	"time"
)

// FileStorage defines the interface for file storage operations
type FileStorage interface {
	// Upload uploads a file and returns the URL
	Upload(file io.Reader, filename string, contentType string) (string, error)

	// Delete removes a file from storage
	Delete(fileURL string) error

	// GetSignedURL generates a temporary signed URL for private files
	GetSignedURL(fileURL string, expiryDuration time.Duration) (string, error)
}
