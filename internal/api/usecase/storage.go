package usecase

import (
	"context"
	"mime/multipart"
)

type StorageService interface {
	Save(ctx context.Context, file *multipart.FileHeader, filename string) error
	Delete(ctx context.Context, filename string) error
	GetPath(filename string) string
}
