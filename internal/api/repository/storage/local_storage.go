// Package storage provides storage implementations for file operations.
package storage

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"open-replays/internal/api/repository/repoiface"
)

// LocalStorage is a storage for files on local disk.
type LocalStorage struct {
	baseDir   string
	publicURL string // for example "http://localhost:8080/media"
}

// Compile-time interface compliance check.
var _ repoiface.StorageService = (*LocalStorage)(nil)

// NewLocalStorage creates a new LocalStorage.
func NewLocalStorage(baseDir, publicURL string) *LocalStorage {
	return &LocalStorage{
		baseDir:   baseDir,
		publicURL: publicURL,
	}
}

// Save saves a file from multipart.FileHeader to storage.
func (s *LocalStorage) Save(ctx context.Context, file *multipart.FileHeader, key string) error {
	src, err := file.Open()
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	defer src.Close()

	return s.SaveReader(ctx, src, key, file.Header.Get("Content-Type"))
}

// SaveReader saves data from an io.Reader to storage.
func (s *LocalStorage) SaveReader(
	_ context.Context,
	reader io.Reader,
	key string,
	_ string,
) error {
	fullPath := filepath.Join(s.baseDir, key)
	dir := filepath.Dir(fullPath)

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create directory: %w", err)
	}

	dst, err := os.Create(fullPath)
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	defer dst.Close()

	if _, err = io.Copy(dst, reader); err != nil {
		return fmt.Errorf("save file: %w", err)
	}

	return nil
}

// GetURL returns the public URL for a file.
func (s *LocalStorage) GetURL(key string) string {
	return fmt.Sprintf("%s/%s", s.publicURL, key)
}

// GetSignedURL returns a signed URL for a file (for local storage, returns public URL).
func (s *LocalStorage) GetSignedURL(
	_ context.Context,
	key string,
	_ time.Duration,
) (string, error) {
	// For local storage we just return the public URL
	return s.GetURL(key), nil
}

// Open opens a file for reading.
func (s *LocalStorage) Open(_ context.Context, key string) (io.ReadCloser, error) {
	fullPath := filepath.Join(s.baseDir, key)
	return os.Open(fullPath)
}

// Exists checks if a file exists in storage.
func (s *LocalStorage) Exists(_ context.Context, key string) (bool, error) {
	fullPath := filepath.Join(s.baseDir, key)
	_, err := os.Stat(fullPath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// Delete deletes a file from storage.
func (s *LocalStorage) Delete(_ context.Context, key string) error {
	fullPath := filepath.Join(s.baseDir, key)
	return os.Remove(fullPath)
}
