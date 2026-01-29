package storage

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

// LocalStorage is a storage for files on local disk.
type LocalStorage struct {
	uploadDir string
}

// NewLocalStorage creates a new LocalStorage.
func NewLocalStorage(uploadDir string) *LocalStorage {
	return &LocalStorage{uploadDir: uploadDir}
}

// Save saves a file to local disk.
func (s *LocalStorage) Save(_ context.Context, fileHeader *multipart.FileHeader, filename string) error {
	src, err := fileHeader.Open()
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer func() {
		_ = src.Close()
	}()

	// Створюємо директорію якщо не існує
	if err := os.MkdirAll(s.uploadDir, 0o755); err != nil {
		return fmt.Errorf("failed to create upload dir: %w", err)
	}

	// Створюємо файл на диску
	dst, err := os.Create(filepath.Join(s.uploadDir, filename))
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer func() {
		_ = dst.Close()
	}()

	// Копіюємо вміст
	if _, err := io.Copy(dst, src); err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}

	return nil
}

// Delete deletes a file from local disk.
func (s *LocalStorage) Delete(_ context.Context, filename string) error {
	path := filepath.Join(s.uploadDir, filename)
	return os.Remove(path)
}

// GetPath returns the path to a file on local disk.
func (s *LocalStorage) GetPath(filename string) string {
	return filepath.Join(s.uploadDir, filename)
}
