package v1

import (
	"mime/multipart"

	"open-replays/api/internal/api/domain"
)

type APIResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
	Code    int    `json:"code,omitempty"`
}

type VideoResponse struct {
	domain.Video
	FullFilename string `json:"fullFilename"`
	Thumbnail    string `json:"thumbnail"`
}

type ListVideosRequest struct{}

type GetVideoRequest struct {
	Filename string `uri:"filename"`
}

type UploadVideoRequest struct {
	Title       string                `form:"title"`
	Description string                `form:"description"`
	Video       *multipart.FileHeader `form:"video"`
	Thumbnail   *multipart.FileHeader `form:"thumbnail"`
}

type DeleteVideoRequest struct {
	Filename string `uri:"filename"`
}

type WatchVideoRequest struct {
	Filename string `uri:"filename"`
}
