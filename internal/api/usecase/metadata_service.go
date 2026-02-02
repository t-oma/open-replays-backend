// Package usecase contains business logic and use case implementations.
package usecase

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"open-replays/internal/api/repository/repoiface"
)

// MetadataService handles video metadata operations.
type MetadataService struct {
	storage repoiface.StorageService
}

// Compile-time interface compliance check.
var _ repoiface.MetadataService = (*MetadataService)(nil)

// NewMetadataService creates a new MetadataService.
func NewMetadataService(storage repoiface.StorageService) *MetadataService {
	return &MetadataService{storage: storage}
}

// GenerateThumbnail generates a thumbnail from a video at the 1-second mark.
func (s *MetadataService) GenerateThumbnail(
	ctx context.Context,
	videoKey string,
	thumbnailKey string,
) error {
	videoTempFile, err := s.downloadToTemp(ctx, videoKey)
	if err != nil {
		return fmt.Errorf("download video: %w", err)
	}
	defer os.Remove(videoTempFile)

	thumbnailTempFile := filepath.Join(
		os.TempDir(),
		"thumb-"+filepath.Base(thumbnailKey),
	)
	defer os.Remove(thumbnailTempFile)

	cmd := exec.CommandContext(ctx, "ffmpeg",
		"-i", videoTempFile,
		"-ss", "00:00:01", // first second
		"-vframes", "1",
		"-vf", "scale=320:-1", // scale to 320px width
		"-y", // overwrite
		thumbnailTempFile,
	)

	if err = cmd.Run(); err != nil {
		return fmt.Errorf("ffmpeg failed: %w", err)
	}

	if err = s.uploadFromFile(ctx, thumbnailTempFile, thumbnailKey); err != nil {
		return fmt.Errorf("upload thumbnail: %w", err)
	}

	return nil
}

// GetDuration extracts the duration of a video in seconds.
func (s *MetadataService) GetDuration(ctx context.Context, videoKey string) (int, error) {
	videoTempFile, err := s.downloadToTemp(ctx, videoKey)
	if err != nil {
		return 0, fmt.Errorf("download video: %w", err)
	}
	defer os.Remove(videoTempFile)

	cmd := exec.CommandContext(ctx, "ffprobe",
		"-v", "error",
		"-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1",
		videoTempFile,
	)

	output, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("ffprobe failed: %w", err)
	}

	durationStr := strings.TrimSpace(string(output))
	duration, err := strconv.ParseFloat(durationStr, 64)
	if err != nil {
		return 0, err
	}

	return int(duration), nil
}

func (s *MetadataService) downloadToTemp(ctx context.Context, key string) (string, error) {
	reader, err := s.storage.Open(ctx, key)
	if err != nil {
		return "", err
	}
	defer reader.Close()

	tempFile, err := os.CreateTemp("", fmt.Sprintf("video-*%s", filepath.Ext(key)))
	if err != nil {
		return "", fmt.Errorf("create temp file: %w", err)
	}
	defer tempFile.Close()

	if _, err = io.Copy(tempFile, reader); err != nil {
		os.Remove(tempFile.Name())
		return "", fmt.Errorf("copy to temp file: %w", err)
	}

	return tempFile.Name(), nil
}

func (s *MetadataService) uploadFromFile(ctx context.Context, filePath string, key string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	contentType := "image/jpeg"
	if filepath.Ext(key) == ".png" {
		contentType = "image/png"
	}

	return s.storage.SaveReader(ctx, file, key, contentType)
}
