// Package usecase contains business logic and use case implementations.
package usecase

import (
	"context"
	"fmt"
	"log"
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
	if err != nil {
		return nil, err
	}
	video.CreateURLs(s.storage.GetURL(""))

	return video, nil
}

// Upload uploads a video.
func (s *VideosService) Upload(ctx context.Context, params UploadParams) (*domain.Video, error) {
	ext := filepath.Ext(params.File.Filename)
	if !s.allowedExtensions[ext] {
		return nil, domain.ErrInvalidFileType
	}
	if params.File.Size > s.maxFileSize {
		return nil, domain.ErrFileTooLarge
	}
	if params.Title == "" {
		return nil, domain.ErrValidation
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
		return nil, fmt.Errorf("create video record: %w", err)
	}

	videoKey := savedVideo.GetVideoKey() // "videos/{id}{ext}"

	if err = s.storage.Save(ctx, params.File, videoKey); err != nil {
		// Rollback: delete video record
		_ = s.repo.Delete(ctx, savedVideo.ID)
		return nil, fmt.Errorf("save file: %w", err)
	}

	if s.shouldGenerateThumbnail(params.Thumbnail) {
		s.processor.Enqueue(ProcessingJob{
			VideoID:  savedVideo.ID,
			VideoExt: ext,
		})
	} else {
		thumbnailKey := savedVideo.GetThumbnailKey()
		if err = s.storage.Save(ctx, params.Thumbnail, thumbnailKey); err != nil {
			log.Printf("save thumbnail: %v", err)
		} else {
			// Update thumbnail path in DB
			_ = s.repo.UpdateVideoMetadata(ctx, savedVideo.ID, 0)
		}
	}

	savedVideo.CreateURLs(s.storage.GetURL(""))

	return savedVideo, nil
}

// Delete deletes a video.
func (s *VideosService) Delete(ctx context.Context, params DeleteParams) error {
	video, err := s.GetByID(ctx, GetByIDParams(params))
	if err != nil {
		return err
	}

	videoKey := video.GetVideoKey()
	thumbnailKey := video.GetThumbnailKey()

	if err = s.storage.Delete(ctx, videoKey); err != nil {
		log.Printf("delete video file %s: %v", videoKey, err)
	}

	if video.ThumbnailURL != "" {
		if err = s.storage.Delete(ctx, thumbnailKey); err != nil {
			log.Printf("delete thumbnail %s: %v", thumbnailKey, err)
		}
	}

	return s.repo.Delete(ctx, video.ID)
}

func (s *VideosService) shouldGenerateThumbnail(thumbnail *multipart.FileHeader) bool {
	return thumbnail == nil || thumbnail.Size == 0
}
