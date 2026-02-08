// Package httperr provides HTTP error handling with structured error codes.
package httperr

import (
	"errors"
	"net/http"
)

// ErrorCode represents a machine-readable error code.
type ErrorCode string

// Common error codes.
const (
	// Validation errors.
	ErrCodeValidation     ErrorCode = "VALIDATION_ERROR"
	ErrCodeInvalidRequest ErrorCode = "INVALID_REQUEST"

	// File upload errors.
	ErrCodeFileTooLarge    ErrorCode = "FILE_TOO_LARGE"
	ErrCodeInvalidFileType ErrorCode = "INVALID_FILE_TYPE"
	ErrCodeFileNotFound    ErrorCode = "FILE_NOT_FOUND"

	// Resource errors.
	ErrCodeVideoNotFound ErrorCode = "VIDEO_NOT_FOUND"
	ErrCodeVideoExists   ErrorCode = "VIDEO_ALREADY_EXISTS"
	ErrCodeDeleteFailed  ErrorCode = "DELETE_FAILED"

	// Server errors.
	ErrCodeInternal ErrorCode = "INTERNAL_ERROR"
	ErrCodeDatabase ErrorCode = "DATABASE_ERROR"
	ErrCodeStorage  ErrorCode = "STORAGE_ERROR"
)

// AppError represents a structured application error.
type AppError struct {
	Code       ErrorCode `json:"code"`
	Message    string    `json:"message"`
	Details    any       `json:"details,omitempty"`
	HTTPStatus int
	Err        error
}

// Error implements the error interface.
func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Message
}

// Unwrap allows errors.Is and errors.As to work with AppError.
func (e *AppError) Unwrap() error {
	return e.Err
}

// New creates a new AppError.
func New(code ErrorCode, message string, httpStatus int) *AppError {
	return &AppError{
		Code:       code,
		Message:    message,
		HTTPStatus: httpStatus,
	}
}

// Wrap wraps an existing error with additional context.
func (e *AppError) Wrap(err error) *AppError {
	e.Err = err
	return e
}

// WithDetails adds details to the error.
func (e *AppError) WithDetails(details any) *AppError {
	e.Details = details
	return e
}

// Predefined errors for common use cases.
var (
	// Validation.
	ErrValidation     = New(ErrCodeValidation, "Validation failed", http.StatusBadRequest)
	ErrInvalidRequest = New(ErrCodeInvalidRequest, "Invalid request", http.StatusBadRequest)

	// File errors.
	ErrFileTooLarge    = New(ErrCodeFileTooLarge, "File too large", http.StatusBadRequest)
	ErrInvalidFileType = New(ErrCodeInvalidFileType, "Invalid file type", http.StatusBadRequest)

	// Video errors.
	ErrVideoNotFound = New(ErrCodeVideoNotFound, "Video not found", http.StatusNotFound)

	// Server errors.
	ErrInternal = New(ErrCodeInternal, "Internal server error", http.StatusInternalServerError)
	ErrDatabase = New(ErrCodeDatabase, "Database error", http.StatusInternalServerError)
	ErrStorage  = New(ErrCodeStorage, "Storage error", http.StatusInternalServerError)
)

// IsAppError checks if an error is an AppError.
func IsAppError(err error) (*AppError, bool) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr, true
	}
	return nil, false
}

// MapError maps domain errors to AppError.
func MapError(err error) *AppError {
	if appErr, ok := IsAppError(err); ok {
		return appErr
	}

	// Map domain errors
	switch {
	case errors.Is(err, ErrVideoNotFound):
		return ErrVideoNotFound.Wrap(err)
	case errors.Is(err, ErrFileTooLarge):
		return ErrFileTooLarge.Wrap(err)
	case errors.Is(err, ErrInvalidFileType):
		return ErrInvalidFileType.Wrap(err)
	case errors.Is(err, ErrValidation):
		return ErrValidation.Wrap(err)
	default:
		return ErrInternal.Wrap(err)
	}
}
