package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

type LocalStorage struct {
	uploadDir string
	baseURL   string
}

func NewLocalStorage(uploadDir, baseURL string) *LocalStorage {
	// Create upload directories if they don't exist
	os.MkdirAll(filepath.Join(uploadDir, "events"), 0755)
	os.MkdirAll(filepath.Join(uploadDir, "comments"), 0755)

	return &LocalStorage{
		uploadDir: uploadDir,
		baseURL:   baseURL,
	}
}

func (s *LocalStorage) Upload(file io.Reader, filename string, contentType string) (string, error) {
	// Generate unique filename
	ext := filepath.Ext(filename)
	uniqueFilename := fmt.Sprintf("%s%s", uuid.New().String(), ext)

	// Determine subdirectory based on content type (events vs comments)
	// For now, store all in events directory - can be enhanced later
	subdir := "events"
	fullPath := filepath.Join(s.uploadDir, subdir, uniqueFilename)

	// Create the file
	out, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	// Copy the uploaded file to the destination
	if _, err := io.Copy(out, file); err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	// Return the public URL
	fileURL := fmt.Sprintf("%s/uploads/%s/%s", s.baseURL, subdir, uniqueFilename)
	return fileURL, nil
}

func (s *LocalStorage) Delete(fileURL string) error {
	// Extract filename from URL
	// URL format: http://localhost:8080/uploads/events/{filename}
	// Parse the path and delete the file
	
	// For now, return nil (implement deletion later if needed)
	return nil
}

func (s *LocalStorage) GetSignedURL(fileURL string, expiryDuration time.Duration) (string, error) {
	// Local storage doesn't need signed URLs - files are publicly accessible
	return fileURL, nil
}
