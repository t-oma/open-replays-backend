package interfaces

import "context"

// MetadataService service responsible for thumbnail generation and extraction
type MetadataService interface {
	Generate(ctx context.Context, videoKey string, thumbnailKey string) error
	GetDuration(ctx context.Context, videoKey string) (int, error)
}
