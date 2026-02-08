// Package domain contains domain models and business logic.
package domain

import (
	"errors"
)

var (
	// ErrVideoNotFound represents a video not found error.
	ErrVideoNotFound = errors.New("video not found")
	// ErrInternal represents an internal error.
	ErrInternal = errors.New("internal server error")
)

// IsNotFound checks if an error is a not found error.
func IsNotFound(err error) bool {
	return errors.Is(err, ErrVideoNotFound)
}
