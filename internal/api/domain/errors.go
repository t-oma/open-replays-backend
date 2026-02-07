// Package domain contains domain models and business logic.
package domain

import "errors"

var (
	// ErrValidation represents a validation error.
	ErrValidation = errors.New("validation error")
	// ErrVideoNotFound represents a video not found error.
	ErrVideoNotFound = errors.New("video not found")
	// ErrInvalidFileType represents an invalid file type error.
	ErrInvalidFileType = errors.New("invalid file type: only .mp4, .webm, .mov allowed")
	// ErrFileTooLarge represents a file too large error.
	ErrFileTooLarge = errors.New("file too large: max 100MB")
)

// IsNotFound checks if an error is a not found error.
func IsNotFound(err error) bool {
	return errors.Is(err, ErrVideoNotFound)
}

// IsValidationError checks if an error is a validation error.
func IsValidationError(err error) bool {
	return errors.Is(err, ErrValidation)
}

// IsInvalidFileType checks if an error is an invalid file type error.
func IsInvalidFileType(err error) bool {
	return errors.Is(err, ErrInvalidFileType)
}

// IsFileTooLarge checks if an error is a file too large error.
func IsFileTooLarge(err error) bool {
	return errors.Is(err, ErrFileTooLarge)
}
