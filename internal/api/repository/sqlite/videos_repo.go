package sqlite

import (
	"context"
	"database/sql"

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
			Filename:    row.Filename,
			Extension:   row.Extension,
			Title:       row.Title,
			Description: row.Description,
			UploadedAt:  row.UploadedAt,
		})
	}

	return out, nil
}

// GetByFilename gets a video by filename.
func (r *VideosRepo) GetByFilename(ctx context.Context, filename string) (domain.Video, error) {
	row, err := r.q.GetVideo(ctx, filename)
	if err != nil {
		return domain.Video{}, err
	}

	return domain.Video{
		Filename:    row.Filename,
		Extension:   row.Extension,
		Title:       row.Title,
		Description: row.Description,
		UploadedAt:  row.UploadedAt,
	}, nil
}

// Create uploads a video.
func (r *VideosRepo) Create(ctx context.Context, video domain.Video) (domain.Video, error) {
	row, err := r.q.CreateVideo(ctx, sqlc.CreateVideoParams{
		Title:       video.Title,
		Description: video.Description,
		Filename:    video.Filename,
		Extension:   video.Extension,
		UploadedAt:  video.UploadedAt,
	})
	if err != nil {
		return domain.Video{}, err
	}

	return domain.Video{
		Filename:    row.Filename,
		Extension:   row.Extension,
		Title:       row.Title,
		Description: row.Description,
		UploadedAt:  row.UploadedAt,
	}, nil
}

// Delete deletes a video.
func (r *VideosRepo) Delete(ctx context.Context, filename string) error {
	return r.q.DeleteVideo(ctx, filename)
}
