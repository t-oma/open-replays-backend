// Package repoiface defines repository interfaces.
package repoiface

import (
	"context"

	"open-replays/internal/api/domain"
)

// VideosRepository defines the interface for video storage operations.
type VideosRepository interface {
	// List returns all videos from the repository.
	List(ctx context.Context) ([]domain.Video, error)
	// GetByID returns a video by its ID.
	GetByID(ctx context.Context, id string) (*domain.Video, error)
	// Create creates a new video record in the repository.
	Create(ctx context.Context, video domain.Video) (*domain.Video, error)
	// UpdateVideoMetadata updates the duration metadata for a video.
	UpdateVideoMetadata(ctx context.Context, id string, duration int) error
	// Delete removes a video from the repository by its ID.
	Delete(ctx context.Context, filename string) error
}
