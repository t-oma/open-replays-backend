package usecase

import (
	"context"
	"fmt"
	"log"

	"open-replays/api/internal/api/repository/interfaces"
)

// VideoProcessor is a service for processing videos.
type VideoProcessor struct {
	metadataService interfaces.MetadataService
	repo            interfaces.VideosRepository
	storage         interfaces.StorageService
	jobQueue        chan ProcessingJob
}

// ProcessingJob is a job for processing a video.
type ProcessingJob struct {
	VideoID       string
	VideoFilename string
	VideoExt      string
}

// NewVideoProcessor creates a new VideoProcessor.
func NewVideoProcessor(
	thumbnailService interfaces.MetadataService,
	repo interfaces.VideosRepository,
	storage interfaces.StorageService,
	workers int,
) *VideoProcessor {
	p := &VideoProcessor{
		metadataService: thumbnailService,
		repo:            repo,
		storage:         storage,
		jobQueue:        make(chan ProcessingJob, 100),
	}

	for range workers {
		go p.worker()
	}

	return p
}

// Enqueue enqueues a video for processing
func (p *VideoProcessor) Enqueue(job ProcessingJob) {
	p.jobQueue <- job
}

func (p *VideoProcessor) worker() {
	for job := range p.jobQueue {
		if err := p.processVideo(context.Background(), job); err != nil {
			log.Printf("Failed to process video %s: %v", job.VideoID, err)
		}
	}
}

func (p *VideoProcessor) processVideo(ctx context.Context, job ProcessingJob) error {
	videoKey := fmt.Sprintf("videos/%s%s", job.VideoFilename, job.VideoExt)
	thumbnailKey := fmt.Sprintf("thumbnails/%s.jpg", job.VideoFilename)

	// Generate thumbnail
	if err := p.metadataService.Generate(ctx, videoKey, thumbnailKey); err != nil {
		return fmt.Errorf("failed to generate thumbnail: %w", err)
	}

	// Get video duration
	duration, err := p.metadataService.GetDuration(ctx, videoKey)
	if err != nil {
		log.Printf("Failed to extract duration for video %s: %v", job.VideoID, err)
		duration = 0
	}

	// 3. Оновлюємо БД
	thumbnailURL := p.storage.GetURL(thumbnailKey)
	return p.repo.UpdateVideoMetadata(ctx, job.VideoID, thumbnailURL, duration)
}
