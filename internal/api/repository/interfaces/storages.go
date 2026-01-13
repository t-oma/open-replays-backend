package interfaces

import (
	"context"
	"mime/multipart"
)

// StorageService service for files management
type StorageService interface {
	Save(ctx context.Context, file *multipart.FileHeader, filename string) error
	Delete(ctx context.Context, filename string) error
	GetPath(filename string) string
}
