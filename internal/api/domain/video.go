package domain

import (
	"fmt"
	"time"
)

// Video is a video domain model.
type Video struct {
	ID           string    `json:"id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	Extension    string    `json:"extension"`
	Duration     int       `json:"duration"`
	VideoURL     string    `json:"videoUrl"`
	ThumbnailURL string    `json:"thumbnailUrl"`
	Views        int       `json:"views"`
	UploadedAt   time.Time `json:"uploadedAt"`
}

// GetVideoKey returns the key for the video file.
func (v *Video) GetVideoKey() string {
	return fmt.Sprintf("videos/%s%s", v.ID, v.Extension)
}

// GetThumbnailKey returns the key for the thumbnail file.
func (v *Video) GetThumbnailKey() string {
	return fmt.Sprintf("thumbnails/%s.jpg", v.ID)
}

// CreateURLs creates URLs for the video and thumbnail.
func (v *Video) CreateURLs(baseURL string) {
	v.VideoURL = fmt.Sprintf("%s%s", baseURL, v.GetVideoKey())
	v.ThumbnailURL = fmt.Sprintf("%s%s", baseURL, v.GetThumbnailKey())
}
