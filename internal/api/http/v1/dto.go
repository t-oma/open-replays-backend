package v1

import (
	"mime/multipart"
)

type APIResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

type ListVideosRequest struct{}

type GetVideoRequest struct {
	ID string `uri:"id"`
}

type UploadVideoRequest struct {
	Title       string                `form:"title"`
	Description string                `form:"description"`
	Video       *multipart.FileHeader `form:"video"`
	Thumbnail   *multipart.FileHeader `form:"thumbnail"`
}

type DeleteVideoRequest struct {
	ID string `uri:"id"`
}

type WatchVideoRequest struct {
	ID string `uri:"id"`
}

type VideoSummaryDTO struct {
	ID           string `json:"id"`
	Title        string `json:"title"`
	ThumbnailURL string `json:"thumbnailUrl"`
	Duration     int    `json:"duration"`
	UploadedAt   string `json:"uploadedAt"`
}

type VideoDetailDTO struct {
	ID           string       `json:"id"`
	Title        string       `json:"title"`
	Description  string       `json:"description"`
	ThumbnailURL string       `json:"thumbnailUrl"`
	VideoURL     string       `json:"videoUrl"`
	UploadedAt   string       `json:"uploadedAt"`
	Duration     int          `json:"duration"`
	Views        int          `json:"views"`
	Author       *AuthorDTO   `json:"author,omitempty"`
	Comments     []CommentDTO `json:"comments,omitempty"`
}

type AuthorDTO struct {
	ID int64 `json:"id"`
	// Username  string `json:"username"`
	// AvatarURL string `json:"avatarUrl,omitempty"`
}

type CommentDTO struct {
	ID int64 `json:"id"`
	// Author    *AuthorDTO
	// Text      string `json:"text"`
	// CreatedAt string `json:"createdAt"`
}

type PaginationDTO struct {
	// Page     int  `json:"page"`
	// PageSize int  `json:"pageSize"`
	// Total    int  `json:"total"`
	// HasNext  bool `json:"hasNext"`
}
