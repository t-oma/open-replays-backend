package usecase

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"path/filepath"
	"sort"
	"time"

	"open-replays/api/internal/api/domain"
	"open-replays/api/internal/api/repository/interfaces"
	"open-replays/api/internal/api/stringutil"
)

const (
	// MEGABYTE represents one megabyte in bytes.
	MEGABYTE = 1024 * 1024
	// MaxFileSize represents the maximum allowed file size.
	MaxFileSize = 100 * MEGABYTE
)

// AllowedExtensions is a map of all allowed extensions.
var AllowedExtensions = map[string]bool{
	".mp4":  true,
	".webm": true,
	".mov":  true,
}

// VideosService is a service for videos.
type VideosService struct {
	repo      interfaces.VideosRepository
	storage   interfaces.StorageService
	processor *VideoProcessor
}

// NewVideosService creates a new VideosService.
func NewVideosService(
	repo interfaces.VideosRepository,
	storage interfaces.StorageService,
	processor *VideoProcessor,
) *VideosService {
	return &VideosService{
		repo:      repo,
		storage:   storage,
		processor: processor,
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

	return videos, nil
}

// GetByID gets a video by ID.
func (s *VideosService) GetByID(ctx context.Context, params GetByIDParams) (*domain.Video, error) {
	return s.repo.GetByID(ctx, params.ID)
}

// Upload uploads a video.
func (s *VideosService) Upload(ctx context.Context, params UploadParams) (*domain.Video, error) {
	ext := filepath.Ext(params.File.Filename)
	if !AllowedExtensions[ext] {
		return nil, domain.ErrInvalidFileType
	}
	if params.File.Size > MaxFileSize {
		return nil, domain.ErrFileTooLarge
	}
	if params.Title == "" {
		return nil, domain.ErrValidation
	}

	uploadTime := time.Now()
	uniqueFilename := fmt.Sprintf(
		"%s-%d", stringutil.TrimExt(params.File.Filename), uploadTime.Unix(),
	)
	videoPath := fmt.Sprintf("videos/%s%s", uniqueFilename, ext)

	// Save file via storage service
	if err := s.storage.Save(ctx, params.File, videoPath); err != nil {
		return nil, fmt.Errorf("failed to save file: %w", err)
	}

	video := domain.Video{
		Title:       params.Title,
		Description: params.Description,
		Filename:    uniqueFilename,
		Extension:   ext,
		Duration:    0,
		UploadedAt:  uploadTime,
	}

	savedVideo, err := s.repo.Create(ctx, video)
	if err != nil {
		// Rollback: delete file if failed to save video metadata
		if delErr := s.storage.Delete(ctx, videoPath); delErr != nil {
			log.Printf("failed to rollback file deletion: %v", delErr)
		}
		return nil, fmt.Errorf("failed to save video metadata: %w", err)
	}

	thumbnailPath := fmt.Sprintf("thumbnails/%s.jpg", uniqueFilename)
	if s.shouldGenerateThumbnail(params.Thumbnail) {
		s.processor.Enqueue(ProcessingJob{
			VideoID:       savedVideo.ID,
			VideoFilename: savedVideo.Filename,
			VideoExt:      ext,
		})
	} else {
		if err := s.storage.Save(ctx, params.Thumbnail, thumbnailPath); err != nil {
			log.Printf("failed to save thumbnail: %v", err)
		} else {
			// Update thumbnail path in DB
			_ = s.repo.UpdateVideoMetadata(ctx, savedVideo.ID, thumbnailPath, 0)
		}
	}

	return savedVideo, nil
}

// Delete deletes a video.
func (s *VideosService) Delete(ctx context.Context, params DeleteParams) error {
	err := s.storage.Delete(ctx, params.ID)
	if err != nil {
		return err
	}

	return s.repo.Delete(ctx, params.ID)
}

// Watch watch a video
func (s *VideosService) Watch(ctx context.Context, params WatchParams) (io.ReadCloser, *domain.Video, error) {
	if params.ID == "" {
		return nil, nil, domain.ErrVideoNotFound
	}

	video, err := s.GetByID(ctx, GetByIDParams(params))
	if err != nil {
		return nil, nil, err
	}

	videoKey := fmt.Sprintf("videos/%s%s", video.Filename, video.Extension)

	exists, err := s.storage.Exists(ctx, videoKey)
	if err != nil || !exists {
		return nil, nil, domain.ErrFileNotFound
	}

	reader, err := s.storage.Open(ctx, videoKey)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open video file: %w", err)
	}

	return reader, video, nil
}

func (s *VideosService) shouldGenerateThumbnail(thumbnail *multipart.FileHeader) bool {
	return thumbnail == nil || thumbnail.Size == 0
}
