package storage

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

// LocalStorage is a storage for files on local disk
type LocalStorage struct {
	baseDir   string
	publicURL string // for example "http://localhost:8080/media"
}

// NewLocalStorage creates a new LocalStorage
func NewLocalStorage(baseDir, publicURL string) *LocalStorage {
	return &LocalStorage{
		baseDir:   baseDir,
		publicURL: publicURL,
	}
}

func (s *LocalStorage) Save(ctx context.Context, file *multipart.FileHeader, key string) error {
	src, err := file.Open()
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	return s.SaveReader(ctx, src, key, file.Header.Get("Content-Type"))
}

func (s *LocalStorage) SaveReader(ctx context.Context, reader io.Reader, key string, contentType string) error {
	fullPath := filepath.Join(s.baseDir, key)
	dir := filepath.Dir(fullPath)

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("failed to create dir: %w", err)
	}

	dst, err := os.Create(fullPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, reader); err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}

	return nil
}

func (s *LocalStorage) GetURL(key string) string {
	return fmt.Sprintf("%s/%s", s.publicURL, key)
}

func (s *LocalStorage) GetSignedURL(ctx context.Context, key string, expiry time.Duration) (string, error) {
	// For local storage we just return the public URL
	return s.GetURL(key), nil
}

func (s *LocalStorage) Open(ctx context.Context, key string) (io.ReadCloser, error) {
	fullPath := filepath.Join(s.baseDir, key)
	return os.Open(fullPath)
}

func (s *LocalStorage) Exists(ctx context.Context, key string) (bool, error) {
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

func (s *LocalStorage) Delete(ctx context.Context, key string) error {
	fullPath := filepath.Join(s.baseDir, key)
	return os.Remove(fullPath)
}
