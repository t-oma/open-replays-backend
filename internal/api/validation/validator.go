package validation

import (
	"mime/multipart"

	"open-replays/internal/api/response"
)

type VideoValidator interface {
	ValidateUpload(req UploadVideoRequest) []response.ValidationError
	ValidateTitle(title string) []response.ValidationError
	ValidateVideoFile(file *multipart.FileHeader) []response.ValidationError
	ValidateThumbnailFile(file *multipart.FileHeader) []response.ValidationError
}

type UploadVideoRequest interface {
	GetTitle() string
	GetDescription() string
	GetVideo() *multipart.FileHeader
	GetThumbnail() *multipart.FileHeader
}
