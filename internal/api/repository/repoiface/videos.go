// Package repoiface defines repository interfaces.
package repoiface

import (
	"context"

	"open-replays/internal/api/domain"
)

// VideosRepository defines the interface for video storage operations.
type VideosRepository interface {
	List(ctx context.Context) ([]domain.Video, error)
	GetByID(ctx context.Context, id string) (*domain.Video, error)
	Create(ctx context.Context, video domain.Video) (*domain.Video, error)
	UpdateVideoMetadata(ctx context.Context, id string, duration int) error
	Delete(ctx context.Context, filename string) error
}
