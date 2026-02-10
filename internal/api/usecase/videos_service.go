// Package usecase contains business logic and use case implementations.
package usecase

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"mime/multipart"
	"time"

	"open-replays/internal/api/domain"
	"open-replays/internal/api/repository/repoiface"
)

// VideosService is a service for videos.
type VideosService struct {
	repo      repoiface.VideosRepository
	storage   repoiface.StorageService
	processor *VideoProcessor
}

// NewVideosService creates a new VideosService.
func NewVideosService(
	repo repoiface.VideosRepository,
	storage repoiface.StorageService,
	processor *VideoProcessor,
) *VideosService {
	return &VideosService{
		repo:      repo,
		storage:   storage,
		processor: processor,
	}
}

// List lists all videos.
func (s *VideosService) List(ctx context.Context, params ListParams) ([]domain.Video, error) {
	offset := (params.Page - 1) * params.PageSize
	limit := params.PageSize

	videos, err := s.repo.List(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	for i := range len(videos) {
		videos[i].CreateURLs(s.storage.GetURL(""))
	}

	return videos, nil
}

// Count returns a count of videos.
func (s *VideosService) Count(ctx context.Context) (int64, error) {
	count, err := s.repo.Count(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, domain.ErrVideoNotFound
	}
	if err != nil {
		return 0, err
	}

	return count, nil
}

// GetByID gets a video by ID.
func (s *VideosService) GetByID(ctx context.Context, params GetByIDParams) (*domain.Video, error) {
	video, err := s.repo.GetByID(ctx, params.ID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrVideoNotFound
	}
	if err != nil {
		return nil, err
	}
	video.CreateURLs(s.storage.GetURL(""))

	return video, nil
}

// Upload uploads a video.
func (s *VideosService) Upload(ctx context.Context, params UploadParams) (*domain.Video, error) {
	video := domain.Video{
		Title:       params.Title,
		Description: params.Description,
		Extension:   params.Ext,
		Duration:    0,
		Views:       0,
		UploadedAt:  time.Now(),
	}

	savedVideo, err := s.repo.Create(ctx, video)
	if err != nil {
		return nil, domain.ErrInternal
	}

	videoKey := savedVideo.GetVideoKey() // "videos/{id}{ext}"

	if err = s.storage.Save(ctx, params.File, videoKey); err != nil {
		// Rollback: delete video record
		_ = s.repo.Delete(ctx, savedVideo.ID)
		return nil, domain.ErrInternal
	}

	if s.shouldGenerateThumbnail(params.Thumbnail) {
		s.processor.Enqueue(ProcessingJob{
			VideoID:  savedVideo.ID,
			VideoExt: params.Ext,
		})
		slog.Info("thumbnail generation queued", "video_id", savedVideo.ID)
	} else {
		thumbnailKey := savedVideo.GetThumbnailKey()
		if err = s.storage.Save(ctx, params.Thumbnail, thumbnailKey); err != nil {
			slog.Error("save thumbnail", "video_id", savedVideo.ID, "error", err)
		} else {
			// Update thumbnail path in DB
			_ = s.repo.UpdateVideoMetadata(ctx, savedVideo.ID, 0)
			slog.Info("thumbnail uploaded", "video_id", savedVideo.ID)
		}
	}

	savedVideo.CreateURLs(s.storage.GetURL(""))

	return savedVideo, nil
}

// Delete deletes a video.
func (s *VideosService) Delete(ctx context.Context, params DeleteParams) error {
	video, err := s.GetByID(ctx, GetByIDParams(params))
	if errors.Is(err, sql.ErrNoRows) {
		return domain.ErrVideoNotFound
	}
	if err != nil {
		return err
	}

	videoKey := video.GetVideoKey()
	thumbnailKey := video.GetThumbnailKey()

	if err = s.storage.Delete(ctx, videoKey); err != nil {
		slog.Error("delete video file", "key", videoKey, "error", err)
	}

	if video.ThumbnailURL != "" {
		if err = s.storage.Delete(ctx, thumbnailKey); err != nil {
			slog.Error("delete thumbnail", "key", thumbnailKey, "error", err)
		}
	}

	return s.repo.Delete(ctx, video.ID)
}

func (s *VideosService) shouldGenerateThumbnail(thumbnail *multipart.FileHeader) bool {
	return thumbnail == nil || thumbnail.Size == 0
}
