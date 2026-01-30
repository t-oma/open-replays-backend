package interfaces

import (
	"context"
	"io"
	"mime/multipart"
	"time"
)

// StorageService provides an interface for file storage operations (local storage or remote).
type StorageService interface {
	// Save saves file with multipart.FileHeader
	Save(ctx context.Context, file *multipart.FileHeader, key string) error

	// SaveReader saves data from io.Reader (for generated files)
	SaveReader(ctx context.Context, reader io.Reader, key string, contentType string) error

	// Detele deletes file
	Delete(ctx context.Context, key string) error

	// Exists checks if file exists
	Exists(ctx context.Context, key string) (bool, error)

	// GetURL return public URL for file access
	GetURL(key string) string

	// GetSignedURL returns a signed URL with a limited expiration time (for private files)
	GetSignedURL(ctx context.Context, key string, expiry time.Duration) (string, error)

	// Open opens file for reading (for internal operations like ffmpeg)
	Open(ctx context.Context, key string) (io.ReadCloser, error)
}
