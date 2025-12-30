package usecase

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"open-replays/api/internal/api/domain"
	"open-replays/api/internal/api/repository/interfaces"
	"open-replays/api/internal/api/stringutil"
)

const (
	MEGABYTE    = 1024 * 1024
	MaxFileSize = 100 * MEGABYTE
)

var AllowedExtensions = map[string]bool{
	".mp4":  true,
	".webm": true,
	".mov":  true,
}

// VideosService is a service for videos.
type VideosService struct {
	repo    interfaces.VideosRepository
	storage StorageService
}

// NewVideosService creates a new VideosService.
func NewVideosService(repo interfaces.VideosRepository, storage StorageService) *VideosService {
	return &VideosService{
		repo:    repo,
		storage: storage,
	}
}

// List lists all videos.
func (s *VideosService) List(ctx context.Context) ([]domain.Video, error) {
	return s.repo.List(ctx)
}

// GetByFilename gets a video by filename.
func (s *VideosService) GetByFilename(ctx context.Context, params GetByFilenameParams) (domain.Video, error) {
	return s.repo.GetByFilename(ctx, params.Filename)
}

// Upload uploads a video.
func (s *VideosService) Upload(ctx context.Context, params UploadParams) (domain.Video, error) {
	// 1. Extension validation
	ext := filepath.Ext(params.File.Filename)
	if !AllowedExtensions[ext] {
		return domain.Video{}, domain.ErrInvalidFileType
	}

	// 2. Size validation
	if params.File.Size > MaxFileSize {
		return domain.Video{}, domain.ErrFileTooLarge
	}

	// 3. Bisness rules validation (title required)
	if params.Title == "" {
		return domain.Video{}, domain.ErrValidation
	}

	// 4. Unique filename generation
	uploadTime := time.Now()
	baseFilename := stringutil.TrimExt(params.File.Filename)
	uniqueFilename := fmt.Sprintf("%s-%d", baseFilename, uploadTime.Unix())

	// 5. Save file via storage service
	if err := s.storage.Save(ctx, params.File, uniqueFilename+ext); err != nil {
		return domain.Video{}, fmt.Errorf("failed to save file: %w", err)
	}

	// 6. Create domain model
	video := domain.Video{
		Title:       params.Title,
		Description: params.Description,
		Filename:    uniqueFilename,
		Extension:   ext[1:],
		UploadedAt:  uploadTime,
	}

	// 7. Save video metadata
	savedVideo, err := s.repo.Create(ctx, video)
	if err != nil {
		// Rollback: delete file if failed to save video metadata
		_ = s.storage.Delete(ctx, uniqueFilename+ext)
		return domain.Video{}, fmt.Errorf("failed to save video metadata: %w", err)
	}

	return savedVideo, nil
}

// Delete deletes a video.
func (s *VideosService) Delete(ctx context.Context, params DeleteParams) error {
	filepath := filepath.Join("uploads/videos", params.Filename)

	err := os.Remove(filepath)
	if err != nil {
		return err
	}

	return s.repo.Delete(ctx, params.Filename)
}

// Watch watch a video
func (s *VideosService) Watch(ctx context.Context, params WatchParams) (string, error) {
	if params.Filename == "" {
		return "", domain.ErrFileNotFound
	}

	getByFilenameParams := GetByFilenameParams{Filename: params.Filename}
	video, err := s.GetByFilename(ctx, getByFilenameParams)
	if err != nil {
		return "", err
	}

	files, _ := filepath.Glob("uploads/videos/" + video.Filename + "." + video.Extension)
	if len(files) == 0 {
		return "", domain.ErrFileNotFound
	}

	return files[0], nil
}
