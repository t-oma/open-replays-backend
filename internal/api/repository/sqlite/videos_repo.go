// Package sqlite provides SQLite repository implementations.
package sqlite

import (
	"context"
	"database/sql"

	gonanoid "github.com/matoous/go-nanoid/v2"

	"open-replays/internal/api/db/sqlc"
	"open-replays/internal/api/domain"
	"open-replays/internal/api/repository/repoiface"
)

// VideosRepo is a repository for videos.
type VideosRepo struct {
	q *sqlc.Queries
}

// Compile-time interface compliance check.
var _ repoiface.VideosRepository = (*VideosRepo)(nil)

// NewVideosRepo creates a new VideosRepo.
func NewVideosRepo(db *sql.DB) *VideosRepo {
	return &VideosRepo{q: sqlc.New(db)}
}

// List lists all videos.
func (r *VideosRepo) List(ctx context.Context) ([]domain.Video, error) {
	rows, err := r.q.ListVideos(ctx)
	if err != nil {
		return nil, err
	}

	out := make([]domain.Video, 0, len(rows))
	for _, row := range rows {
		out = append(out, domain.Video{
			ID:          row.ID,
			Title:       row.Title,
			Description: row.Description,
			Extension:   row.Extension,
			Duration:    int(row.Duration),
			Views:       int(row.Views),
			UploadedAt:  row.UploadedAt,
		})
	}

	return out, nil
}

// GetByID gets a video by ID.
func (r *VideosRepo) GetByID(ctx context.Context, id string) (*domain.Video, error) {
	row, err := r.q.GetVideo(ctx, id)
	if err != nil {
		return nil, err
	}

	return &domain.Video{
		ID:          row.ID,
		Extension:   row.Extension,
		Title:       row.Title,
		Description: row.Description,
		Duration:    int(row.Duration),
		Views:       int(row.Views),
		UploadedAt:  row.UploadedAt,
	}, nil
}

// Create uploads a video.
func (r *VideosRepo) Create(ctx context.Context, video domain.Video) (*domain.Video, error) {
	id, err := gonanoid.New()
	if err != nil {
		return nil, err
	}

	row, err := r.q.CreateVideo(ctx, sqlc.CreateVideoParams{
		ID:          id,
		Title:       video.Title,
		Description: video.Description,
		Extension:   video.Extension,
		Duration:    int64(video.Duration),
		UploadedAt:  video.UploadedAt,
	})
	if err != nil {
		return nil, err
	}

	return &domain.Video{
		ID:          row.ID,
		Title:       row.Title,
		Description: row.Description,
		Extension:   row.Extension,
		Duration:    int(row.Duration),
		Views:       int(row.Views),
		UploadedAt:  row.UploadedAt,
	}, nil
}

// UpdateVideoMetadata updates duration for a video.
func (r *VideosRepo) UpdateVideoMetadata(ctx context.Context, id string, duration int) error {
	return r.q.UpdateVideoMetadata(ctx, sqlc.UpdateVideoMetadataParams{
		ID:       id,
		Duration: int64(duration),
	})
}

// Delete deletes a video.
func (r *VideosRepo) Delete(ctx context.Context, filename string) error {
	return r.q.DeleteVideo(ctx, filename)
}
