package utils

import (
	"errors"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
)

const MaxFileSize = 25 * 1024 * 1024 // 25 MB

// Allowed MIME types
var AllowedMimeTypes = map[string]bool{
	// Images
	"image/jpeg":    true,
	"image/jpg":     true,
	"image/png":     true,
	"image/gif":     true,
	"image/webp":    true,
	"image/svg+xml": true,
	// Documents
	"application/pdf": true,
	"application/msword": true,
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
	"text/plain": true,
}

// Allowed file extensions (as fallback)
var AllowedExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
	".webp": true,
	".svg":  true,
	".pdf":  true,
	".doc":  true,
	".docx": true,
	".txt":  true,
}

// ValidateFile checks if the uploaded file meets size and format requirements
func ValidateFile(file multipart.File, header *multipart.FileHeader) error {
	// Check file size
	if header.Size > MaxFileSize {
		return errors.New("file size exceeds 25MB limit")
	}

	if header.Size == 0 {
		return errors.New("file is empty")
	}

	// Check file extension
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !AllowedExtensions[ext] {
		return errors.New("file type not allowed")
	}

	// Detect MIME type from content
	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil {
		return errors.New("failed to read file")
	}

	// Reset file pointer
	if _, err := file.Seek(0, 0); err != nil {
		return errors.New("failed to reset file pointer")
	}

	contentType := http.DetectContentType(buffer[:n])

	// Special handling for some types
	if ext == ".svg" {
		contentType = "image/svg+xml"
	} else if ext == ".docx" {
		contentType = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	} else if ext == ".doc" {
		contentType = "application/msword"
	}

	// Check if MIME type is allowed
	if !AllowedMimeTypes[contentType] {
		return errors.New("file type not allowed")
	}

	return nil
}
