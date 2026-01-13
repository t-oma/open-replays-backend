package usecase

import (
	"context"
	"log"

	"open-replays/api/internal/api/repository/interfaces"
	"open-replays/api/internal/api/stringutil"
)

// VideoProcessor is a service for processing videos.
type VideoProcessor struct {
	thumbnailService interfaces.ThumbnailService
	vStorage         interfaces.StorageService
	tStorage         interfaces.StorageService
	jobQueue         chan ProcessingJob
}

// ProcessingJob is a job for processing a video.
type ProcessingJob struct {
	VideoFilename string
	VideoExt      string
}

// NewVideoProcessor creates a new VideoProcessor.
func NewVideoProcessor(
	thumbnailService interfaces.ThumbnailService,
	vStorage interfaces.StorageService,
	tStorage interfaces.StorageService,
	workers int,
) *VideoProcessor {
	p := &VideoProcessor{
		thumbnailService: thumbnailService,
		vStorage:         vStorage,
		tStorage:         tStorage,
		jobQueue:         make(chan ProcessingJob, 100),
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
			log.Printf("Failed to process video %s: %v", job.VideoFilename, err)
		}
	}
}

func (p *VideoProcessor) processVideo(ctx context.Context, job ProcessingJob) error {
	videoFullFilename := job.VideoFilename + job.VideoExt
	videoPath := p.vStorage.GetPath(videoFullFilename)

	thumbnailFullFilename := stringutil.TrimExt(job.VideoFilename) + ".jpg"
	thumbnailPath := p.tStorage.GetPath(thumbnailFullFilename)

	if err := p.thumbnailService.Generate(ctx, videoPath, thumbnailPath); err != nil {
		return err
	}

	return nil
}
