package domain

import "time"

// Video is a video domain model.
type Video struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Filename    string    `json:"filename"`
	Extension   string    `json:"extension"`
	Thumbnail   string    `json:"thumbnail"`
	Duration    int       `json:"duration"`
	Views       int       `json:"views"`
	UploadedAt  time.Time `json:"uploadedAt"`
}
