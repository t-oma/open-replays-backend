// Package repoiface defines repository interfaces.
package repoiface

import (
	"context"
	"io"
	"mime/multipart"
	"time"
)

// StorageService provides an interface for file storage operations (local storage or remote).
type StorageService interface {
	// Save saves a file from multipart.FileHeader to storage.
	Save(ctx context.Context, file *multipart.FileHeader, key string) error

	// SaveReader saves data from io.Reader to storage (for generated files).
	SaveReader(ctx context.Context, reader io.Reader, key string, contentType string) error

	// Delete deletes a file from storage.
	Delete(ctx context.Context, key string) error

	// Exists checks if a file exists in storage.
	Exists(ctx context.Context, key string) (bool, error)

	// GetURL returns the public URL for a file.
	GetURL(key string) string

	// GetSignedURL returns a signed URL with a limited expiration time (for private files).
	GetSignedURL(ctx context.Context, key string, expiry time.Duration) (string, error)

	// Open opens a file for reading (for internal operations like ffmpeg).
	Open(ctx context.Context, key string) (io.ReadCloser, error)
}
