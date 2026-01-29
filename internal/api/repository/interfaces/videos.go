package interfaces

import (
	"context"

	"open-replays/api/internal/api/domain"
)

type VideosRepository interface {
	List(ctx context.Context) ([]domain.Video, error)
	GetByID(ctx context.Context, id string) (domain.Video, error)
	Create(ctx context.Context, video domain.Video) (domain.Video, error)
	Delete(ctx context.Context, filename string) error
}
