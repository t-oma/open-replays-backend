// Package domain contains domain models and business logic.
package domain

import (
	"errors"
	"fmt"
)

var (
	// ErrVideoNotFound represents a video not found error.
	ErrVideoNotFound = errors.New("video not found")
	// ErrInternal represents an internal error.
	ErrInternal = errors.New("internal server error")
)

// FieldError represents a field error with field details.
type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// MultiFieldError holds multiple field errors.
type MultiFieldError struct {
	Errors []FieldError `json:"errors"`
}

// Error implements the error interface.
func (v MultiFieldError) Error() string {
	if len(v.Errors) == 0 {
		return "field validation failed"
	}
	return "field validation failed: " + v.Errors[0].Message
}

// FileSizeError represents a file size validation error.
type FileSizeError struct {
	MaxSize      int64   `json:"maxSizeBytes"`
	MaxSizeMB    int     `json:"maxSizeMb"`
	ActualSize   int64   `json:"actualSizeBytes"`
	ActualSizeMB float64 `json:"actualSizeMb"`
}

// Error implements the error interface.
func (f FileSizeError) Error() string {
	return fmt.Sprintf("file too large: max %d MB, got %.2f MB", f.MaxSizeMB, f.ActualSizeMB)
}

// FileTypeError represents a file type validation error.
type FileTypeError struct {
	AllowedTypes []string `json:"allowedTypes"`
	ActualType   string   `json:"actualType"`
	Filename     string   `json:"filename"`
}

// Error implements the error interface.
func (f FileTypeError) Error() string {
	return fmt.Sprintf("invalid file type: %s. Allowed types: %v", f.ActualType, f.AllowedTypes)
}

// IsNotFound checks if an error is a not found error.
func IsNotFound(err error) bool {
	return errors.Is(err, ErrVideoNotFound)
}

// IsMultiFieldError checks if an error is a multi field error.
func IsMultiFieldError(err error) (*MultiFieldError, bool) {
	var v MultiFieldError
	if errors.As(err, &v) {
		return &v, true
	}
	return nil, false
}

// IsFileSizeError checks if an error is a file size error.
func IsFileSizeError(err error) (*FileSizeError, bool) {
	var f FileSizeError
	if errors.As(err, &f) {
		return &f, true
	}
	return nil, false
}

// IsFileTypeError checks if an error is a file type error.
func IsFileTypeError(err error) (*FileTypeError, bool) {
	var f FileTypeError
	if errors.As(err, &f) {
		return &f, true
	}
	return nil, false
}
