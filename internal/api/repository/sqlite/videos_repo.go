package sqlite

import (
	"context"
	"database/sql"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"open-replays/api/internal/api/db/sqlc"
	"open-replays/api/internal/api/domain"
)

// VideosRepo is a repository for videos.
type VideosRepo struct {
	q *sqlc.Queries
}

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
			Filename:    row.Filename,
			Extension:   row.Extension,
			Duration:    int(row.Duration),
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
		Filename:    row.Filename,
		Extension:   row.Extension,
		Title:       row.Title,
		Description: row.Description,
		Thumbnail:   row.Thumbnail,
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
		Filename:    video.Filename,
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
		Filename:    row.Filename,
		Extension:   row.Extension,
		Duration:    int(row.Duration),
		UploadedAt:  row.UploadedAt,
	}, nil
}

// UpdateVideoMetadata updates thumbnail path and duration for a video.
func (r *VideosRepo) UpdateVideoMetadata(ctx context.Context, id string, thumbnailURL string, duration int) error {
	return r.q.UpdateVideoMetadata(ctx, sqlc.UpdateVideoMetadataParams{
		ID:        id,
		Thumbnail: thumbnailURL,
		Duration:  int64(duration),
	})
}

// Delete deletes a video.
func (r *VideosRepo) Delete(ctx context.Context, filename string) error {
	return r.q.DeleteVideo(ctx, filename)
}
