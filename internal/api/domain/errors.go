// Package domain contains domain models and business logic.
package domain

import "errors"

var (
	// ErrValidation represents a validation error.
	ErrValidation = errors.New("validation error")
	// ErrFileNotFound represents a file not found error.
	ErrFileNotFound = errors.New("file not found")
	// ErrVideoNotFound represents a video not found error.
	ErrVideoNotFound = errors.New("video not found")
	// ErrNotFound represents a not found error.
	ErrNotFound = errors.New("not found")
	// ErrInvalidFileType represents an invalid file type error.
	ErrInvalidFileType = errors.New("invalid file type: only .mp4, .webm, .mov allowed")
	// ErrFileTooLarge represents a file too large error.
	ErrFileTooLarge = errors.New("file too large: max 100MB")
)
