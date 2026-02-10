package usecase

import "mime/multipart"

// ListParams is a params for listing videos.
type ListParams struct {
	Page     int
	PageSize int
}

// GetByIDParams is a params for getting a video by filename.
type GetByIDParams struct {
	ID string
}

// UploadParams is a params for uploading a video.
type UploadParams struct {
	File        *multipart.FileHeader
	Thumbnail   *multipart.FileHeader
	Title       string
	Description string
	Ext         string
}

// DeleteParams is a params for deleting a video.
type DeleteParams struct {
	ID string
}

// WatchParams is a params for watching a video.
type WatchParams struct {
	ID string
}
