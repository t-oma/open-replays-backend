package interfaces

import "context"

// ThumbnailService service responsible for thumbnail generation and extraction
type ThumbnailService interface {
	Generate(ctx context.Context, videoPath string, outputPath string) error
}
