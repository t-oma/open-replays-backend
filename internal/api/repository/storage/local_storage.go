package storage

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

type LocalStorage struct {
	uploadDir string
}

func NewLocalStorage(uploadDir string) *LocalStorage {
	return &LocalStorage{uploadDir: uploadDir}
}

func (s *LocalStorage) Save(ctx context.Context, fileHeader *multipart.FileHeader, filename string) error {
	src, err := fileHeader.Open()
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	// Створюємо директорію якщо не існує
	if err := os.MkdirAll(s.uploadDir, 0o755); err != nil {
		return fmt.Errorf("failed to create upload dir: %w", err)
	}

	// Створюємо файл на диску
	dst, err := os.Create(filepath.Join(s.uploadDir, filename))
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	// Копіюємо вміст
	if _, err := io.Copy(dst, src); err != nil {
		return fmt.Errorf("failed to save file: %w", err)
	}

	return nil
}

func (s *LocalStorage) Delete(ctx context.Context, filename string) error {
	path := filepath.Join(s.uploadDir, filename)
	return os.Remove(path)
}

func (s *LocalStorage) GetPath(filename string) string {
	return filepath.Join(s.uploadDir, filename)
}
