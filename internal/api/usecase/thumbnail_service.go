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

	"open-replays/api/internal/api/repository/interfaces"
)

type FFmpegThumbnailService struct {
	storage interfaces.StorageService
}

func NewFFmpegThumbnailService(storage interfaces.StorageService) *FFmpegThumbnailService {
	return &FFmpegThumbnailService{storage: storage}
}

func (s *FFmpegThumbnailService) Generate(ctx context.Context, videoKey string, thumbnailKey string) error {
	videoTempFile, err := s.downloadToTemp(ctx, videoKey)
	if err != nil {
		return fmt.Errorf("failed to download video: %w", err)
	}
	defer os.Remove(videoTempFile)

	thumbnailTempFile := filepath.Join(os.TempDir(), fmt.Sprintf("thumb-%s", filepath.Base(thumbnailKey)))
	defer os.Remove(thumbnailTempFile)

	cmd := exec.CommandContext(ctx, "ffmpeg",
		"-i", videoTempFile,
		"-ss", "00:00:01", // first second
		"-vframes", "1",
		"-vf", "scale=320:-1", // scale to 320px width
		"-y", // overwrite
		thumbnailTempFile,
	)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ffmpeg failed: %w", err)
	}

	if err := s.uploadFromFile(ctx, thumbnailTempFile, thumbnailKey); err != nil {
		return fmt.Errorf("failed to upload thumbnail: %w", err)
	}

	return nil
}

func (s *FFmpegThumbnailService) GetDuration(ctx context.Context, videoKey string) (int, error) {
	videoTempFile, err := s.downloadToTemp(ctx, videoKey)
	if err != nil {
		return 0, fmt.Errorf("failed to download video: %w", err)
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

func (s *FFmpegThumbnailService) downloadToTemp(ctx context.Context, key string) (string, error) {
	reader, err := s.storage.Open(ctx, key)
	if err != nil {
		return "", err
	}
	defer reader.Close()

	tempFile, err := os.CreateTemp("", fmt.Sprintf("video-*%s", filepath.Ext(key)))
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer tempFile.Close()

	if _, err := io.Copy(tempFile, reader); err != nil {
		os.Remove(tempFile.Name())
		return "", fmt.Errorf("failed to copy to temp file: %w", err)
	}

	return tempFile.Name(), nil
}

func (s *FFmpegThumbnailService) uploadFromFile(ctx context.Context, filePath string, key string) error {
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
