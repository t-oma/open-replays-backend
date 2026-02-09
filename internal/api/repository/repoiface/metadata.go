// Package repoiface defines repository interfaces.
package repoiface

import "context"

// MetadataService provides video metadata extraction operations.
type MetadataService interface {
	// GenerateThumbnail creates a thumbnail image from a video file.
	GenerateThumbnail(ctx context.Context, videoKey string, thumbnailKey string) error
	// GetDuration extracts the duration of a video in seconds.
	GetDuration(ctx context.Context, videoKey string) (int, error)
}
