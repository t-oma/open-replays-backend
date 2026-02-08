// Package usecase contains business logic and use case implementations.
package usecase

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"mime/multipart"
	"path/filepath"
	"sort"
	"time"

	"open-replays/internal/api/domain"
	"open-replays/internal/api/repository/repoiface"
)

// VideosService is a service for videos.
type VideosService struct {
	repo              repoiface.VideosRepository
	storage           repoiface.StorageService
	processor         *VideoProcessor
	maxFileSize       int64
	allowedExtensions map[string]bool
}

// NewVideosService creates a new VideosService.
func NewVideosService(
	repo repoiface.VideosRepository,
	storage repoiface.StorageService,
	processor *VideoProcessor,
	maxFileSize int64,
	allowedExtensions []string,
) *VideosService {
	extMap := make(map[string]bool, len(allowedExtensions))
	for _, ext := range allowedExtensions {
		extMap[ext] = true
	}

	return &VideosService{
		repo:              repo,
		storage:           storage,
		processor:         processor,
		maxFileSize:       maxFileSize,
		allowedExtensions: extMap,
	}
}

// List lists all videos.
func (s *VideosService) List(ctx context.Context) ([]domain.Video, error) {
	videos, err := s.repo.List(ctx)
	if err != nil {
		return nil, err
	}

	sort.Slice(videos, func(i, j int) bool {
		return videos[i].UploadedAt.After(videos[j].UploadedAt)
	})

	for i := range len(videos) {
		videos[i].CreateURLs(s.storage.GetURL(""))
	}

	return videos, nil
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
	var fieldsValidationErr domain.MultiFieldError

	if params.Title == "" {
		fieldsValidationErr.Errors = append(fieldsValidationErr.Errors, domain.FieldError{
			Field:   "title",
			Message: "title is required",
		})
	}
	if params.File == nil {
		fieldsValidationErr.Errors = append(fieldsValidationErr.Errors, domain.FieldError{
			Field:   "video",
			Message: "video is required",
		})
		// If video file not found then we can't go further
		return nil, fieldsValidationErr
	}

	ext := filepath.Ext(params.File.Filename)
	if !s.allowedExtensions[ext] {
		return nil, domain.FileTypeError{
			AllowedTypes: s.getAllowedExtensionsList(),
			ActualType:   ext,
			Filename:     params.File.Filename,
		}
	}

	if params.File.Size > s.maxFileSize {
		return nil, domain.FileSizeError{
			MaxSize:      s.maxFileSize,
			MaxSizeMB:    int(s.maxFileSize / (1024 * 1024)),
			ActualSize:   params.File.Size,
			ActualSizeMB: float64(params.File.Size) / (1024 * 1024),
		}
	}

	if len(fieldsValidationErr.Errors) > 0 {
		return nil, fieldsValidationErr
	}

	video := domain.Video{
		Title:       params.Title,
		Description: params.Description,
		Extension:   ext,
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
			VideoExt: ext,
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

// getAllowedExtensionsList returns a list of allowed file extensions.
func (s *VideosService) getAllowedExtensionsList() []string {
	exts := make([]string, 0, len(s.allowedExtensions))
	for ext := range s.allowedExtensions {
		exts = append(exts, ext)
	}
	return exts
}
