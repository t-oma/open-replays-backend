package usecase

import (
	"context"
	"fmt"
	"os/exec"
)

type FFmpegThumbnailService struct{}

func NewFFmpegThumbnailService() *FFmpegThumbnailService {
	return &FFmpegThumbnailService{}
}

func (s *FFmpegThumbnailService) Generate(ctx context.Context, videoPath string, outputPath string) error {
	// ffmpeg -i input.mp4 -ss 00:00:01 -vframes 1 -vf scale=320:-1 output.jpg
	cmd := exec.CommandContext(ctx, "ffmpeg",
		"-i", videoPath,
		"-ss", "00:00:01", // first second
		"-vframes", "1",
		"-vf", "scale=320:-1", // scale to 320px width
		"-y", // overwrite
		outputPath,
	)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ffmpeg failed: %w", err)
	}

	return nil
}
