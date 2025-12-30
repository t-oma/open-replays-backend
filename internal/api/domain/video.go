package domain

import "time"

type Video struct {
	Filename    string    `json:"filename"`
	Extension   string    `json:"extension"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	UploadedAt  time.Time `json:"uploadedAt"`
}
