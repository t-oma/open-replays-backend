package domain

import "errors"

var (
	ErrValidation      = errors.New("validation error")
	ErrFileNotFound    = errors.New("file not found")
	ErrNotFound        = errors.New("not found")
	ErrInvalidFileType = errors.New("invalid file type: only .mp4, .webm, .mov allowed")
	ErrFileTooLarge    = errors.New("file too large: max 100MB")
)
