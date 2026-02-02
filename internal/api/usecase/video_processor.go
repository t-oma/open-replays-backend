package usecase

import (
	"context"
	"fmt"
	"log"

	"open-replays/internal/api/repository/repoiface"
)

// VideoProcessor is a service for processing videos.
type VideoProcessor struct {
	metadataService repoiface.MetadataService
	repo            repoiface.VideosRepository
	storage         repoiface.StorageService
	jobQueue        chan ProcessingJob
}

// ProcessingJob is a job for processing a video.
type ProcessingJob struct {
	VideoID  string
	VideoExt string
}

// NewVideoProcessor creates a new VideoProcessor.
func NewVideoProcessor(
	thumbnailService repoiface.MetadataService,
	repo repoiface.VideosRepository,
	storage repoiface.StorageService,
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

// Enqueue enqueues a video for processing.
func (p *VideoProcessor) Enqueue(job ProcessingJob) {
	p.jobQueue <- job
}

func (p *VideoProcessor) worker() {
	for job := range p.jobQueue {
		if err := p.processVideo(context.Background(), job); err != nil {
			log.Printf("process video %s: %v", job.VideoID, err)
		}
	}
}

func (p *VideoProcessor) processVideo(ctx context.Context, job ProcessingJob) error {
	videoKey := fmt.Sprintf("videos/%s%s", job.VideoID, job.VideoExt)
	thumbnailKey := fmt.Sprintf("thumbnails/%s.jpg", job.VideoID)

	if err := p.metadataService.GenerateThumbnail(ctx, videoKey, thumbnailKey); err != nil {
		return fmt.Errorf("generate thumbnail: %w", err)
	}

	duration, err := p.metadataService.GetDuration(ctx, videoKey)
	if err != nil {
		log.Printf("extract duration for video %s: %v", job.VideoID, err)
		duration = 0
	}

	return p.repo.UpdateVideoMetadata(ctx, job.VideoID, duration)
}
