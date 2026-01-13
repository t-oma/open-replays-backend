package usecase

import (
	"context"
	"fmt"
	"path/filepath"
	"sort"
	"time"

	"open-replays/api/internal/api/domain"
	"open-replays/api/internal/api/repository/interfaces"
	"open-replays/api/internal/api/stringutil"
)

const (
	MEGABYTE    = 1024 * 1024
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
	repo             interfaces.VideosRepository
	videoStorage     interfaces.StorageService
	thumbnailStorage interfaces.StorageService
	processor        *VideoProcessor
}

// NewVideosService creates a new VideosService.
func NewVideosService(
	repo interfaces.VideosRepository,
	videoStorage interfaces.StorageService,
	thumbnailStorage interfaces.StorageService,
	processor *VideoProcessor,
) *VideosService {
	return &VideosService{
		repo:             repo,
		videoStorage:     videoStorage,
		thumbnailStorage: thumbnailStorage,
		processor:        processor,
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
	if err := s.videoStorage.Save(ctx, params.File, uniqueFilename+ext); err != nil {
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
		_ = s.videoStorage.Delete(ctx, uniqueFilename+ext)
		// _ = s.thumbnailStorage.Delete(ctx, thumbnailFilename)
		return domain.Video{}, fmt.Errorf("failed to save video metadata: %w", err)
	}

	// 8. Generate thumbnail
	thumbnailFilename := uniqueFilename + ".jpg"
	if params.Thumbnail.Filename == "" && params.Thumbnail.Size == 0 {
		s.processor.Enqueue(ProcessingJob{
			VideoFilename: savedVideo.Filename,
			VideoExt:      ext,
		})
	} else {
		if err := s.thumbnailStorage.Save(ctx, params.Thumbnail, thumbnailFilename); err != nil {
			// Rollback: delete video metadata and file if failed to save thumbnail
			_ = s.videoStorage.Delete(ctx, uniqueFilename+ext)
			_ = s.repo.Delete(ctx, savedVideo.Filename)
			return domain.Video{}, fmt.Errorf("failed to save thumbnail: %w", err)
		}
	}

	return savedVideo, nil
}

// Delete deletes a video.
func (s *VideosService) Delete(ctx context.Context, params DeleteParams) error {
	err := s.videoStorage.Delete(ctx, params.Filename)
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

	files, _ := filepath.Glob(s.videoStorage.GetPath(video.Filename + "." + video.Extension))
	if len(files) == 0 {
		return "", domain.ErrFileNotFound
	}

	return files[0], nil
}
