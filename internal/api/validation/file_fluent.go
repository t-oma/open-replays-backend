package validation

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"slices"
	"strings"

	"open-replays/internal/api/response"
)

// FileValidator validates file uploads using fluent API.
type FileValidator struct {
	field          string
	value          *multipart.FileHeader
	errors         []response.ValidationError
	skipValidation bool
}

// File creates a new file validator.
func File(field string, file *multipart.FileHeader) *FileValidator {
	return &FileValidator{
		field:  field,
		value:  file,
		errors: make([]response.ValidationError, 0),
	}
}

// When conditionally applies validation only if the condition is true.
func (v *FileValidator) When(condition bool) *FileValidator {
	v.skipValidation = !condition
	return v
}

// Optional skips validation if file is nil or empty.
func (v *FileValidator) Optional() *FileValidator {
	return v.When(v.value != nil && v.value.Size > 0)
}

// Required checks if file is provided.
func (v *FileValidator) Required() *FileValidator {
	if v.skipValidation {
		return v
	}
	if v.value == nil || v.value.Size <= 0 {
		v.errors = append(v.errors, response.ValidationError{
			Field:   v.field,
			Message: v.field + " is required",
			Tag:     TagRequired,
		})
	}
	return v
}

// MinSize checks minimum file size.
func (v *FileValidator) MinSize(minLen int64) *FileValidator {
	if v.skipValidation {
		return v
	}
	if v.value != nil && v.value.Size < minLen {
		v.errors = append(v.errors, response.ValidationError{
			Field:   v.field,
			Message: v.field + " is too small",
			Tag:     TagMinSize,
			Details: response.FileSizeDetails{
				MinSizeBytes:    minLen,
				ActualSizeBytes: v.value.Size,
			},
		})
	}
	return v
}

// MaxSize checks maximum file size.
func (v *FileValidator) MaxSize(maxLen int64) *FileValidator {
	if v.skipValidation {
		return v
	}
	if v.value != nil && v.value.Size > maxLen {
		v.errors = append(v.errors, response.ValidationError{
			Field:   v.field,
			Message: fmt.Sprintf("%s must not exceed %s", v.field, formatSize(maxLen)),
			Tag:     TagMaxSize,
			Details: response.FileSizeDetails{
				MaxSizeBytes:    maxLen,
				ActualSizeBytes: v.value.Size,
			},
		})
	}
	return v
}

// AllowedExts checks if file extension is allowed.
func (v *FileValidator) AllowedExts(exts []string) *FileValidator {
	if v.skipValidation {
		return v
	}
	if v.value == nil {
		return v
	}

	ext := strings.ToLower(filepath.Ext(v.value.Filename))
	if slices.Contains(exts, ext) {
		return v
	}

	v.errors = append(v.errors, response.ValidationError{
		Field:   v.field,
		Message: v.field + " has invalid format",
		Tag:     TagInvalidFormat,
		Details: response.FileTypeDetails{
			AllowedTypes: exts,
			ActualType:   ext,
			Filename:     v.value.Filename,
		},
	})
	return v
}

// Errors returns collected errors.
func (v *FileValidator) Errors() []response.ValidationError {
	return v.errors
}

// IsValid returns true if no errors.
func (v *FileValidator) IsValid() bool {
	return len(v.errors) == 0
}

// formatSize formats bytes to human readable string.
func formatSize(bytes int64) string {
	switch {
	case bytes >= GB:
		return fmt.Sprintf("%.2f GB", float64(bytes)/GB)
	case bytes >= MB:
		return fmt.Sprintf("%.2f MB", float64(bytes)/MB)
	case bytes >= KB:
		return fmt.Sprintf("%.2f KB", float64(bytes)/KB)
	default:
		return fmt.Sprintf("%d bytes", bytes)
	}
}
