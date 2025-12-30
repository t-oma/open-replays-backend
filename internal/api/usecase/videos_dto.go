package usecase

import "mime/multipart"

// ListParams is a params for listing videos.
type ListParams struct{}

// GetByFilenameParams is a params for getting a video by filename.
type GetByFilenameParams struct {
	Filename string
}

// UploadParams is a params for uploading a video.
type UploadParams struct {
	File        *multipart.FileHeader
	Title       string
	Description string
}

// DeleteParams is a params for deleting a video.
type DeleteParams struct {
	Filename string
}

// WatchParams is a params for watching a video.
type WatchParams struct {
	Filename string
}
