// Package v1 provides HTTP handlers and DTOs for API version 1.
package v1

import (
	"mime/multipart"
)

// APIResponse is the standard response format for all API endpoints.
type APIResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

// ListVideosRequest is the request for listing videos.
type ListVideosRequest struct{}

// GetVideoRequest is the request for getting a video by ID.
type GetVideoRequest struct {
	ID string `uri:"id"`
}

// UploadVideoRequest is the request for uploading a video.
type UploadVideoRequest struct {
	Title       string                `form:"title"`
	Description string                `form:"description"`
	Video       *multipart.FileHeader `form:"video"`
	Thumbnail   *multipart.FileHeader `form:"thumbnail"`
}

// DeleteVideoRequest is the request for deleting a video.
type DeleteVideoRequest struct {
	ID string `uri:"id"`
}

// WatchVideoRequest is the request for watching a video.
type WatchVideoRequest struct {
	ID string `uri:"id"`
}

// VideoSummaryDTO is a summary representation of a video for list views.
type VideoSummaryDTO struct {
	ID           string `json:"id"`
	Title        string `json:"title"`
	ThumbnailURL string `json:"thumbnailUrl"`
	Duration     int    `json:"duration"`
	UploadedAt   string `json:"uploadedAt"`
}

// VideoDetailDTO is a detailed representation of a video.
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

// AuthorDTO represents an author of a video or comment.
type AuthorDTO struct {
	ID int64 `json:"id"`
	// Username  string `json:"username"`
	// AvatarURL string `json:"avatarUrl,omitempty"`
}

// CommentDTO represents a comment on a video.
type CommentDTO struct {
	ID int64 `json:"id"`
	// Author    *AuthorDTO
	// Text      string `json:"text"`
	// CreatedAt string `json:"createdAt"`
}

// PaginationDTO represents pagination information.
type PaginationDTO struct {
	// Page     int  `json:"page"`
	// PageSize int  `json:"pageSize"`
	// Total    int  `json:"total"`
	// HasNext  bool `json:"hasNext"`
}
